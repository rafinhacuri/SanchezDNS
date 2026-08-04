# 🏗️ Arquitetura

Esta página descreve **como o SanchezDNS é montado por dentro**: quais peças existem, qual é o papel de cada uma e, principalmente, **por que** o sistema foi organizado dessa forma. É a leitura recomendada antes de mergulhar nas páginas específicas de PowerDNS, Docker e autenticação.

## Visão de alto nível

O SanchezDNS não é um servidor DNS. Ele é um **painel de operação** que fica na frente de um servidor DNS autoritativo (PowerDNS) e transforma tarefas que normalmente seriam feitas via linha de comando, edição de banco ou chamadas HTTP manuais em uma interface web com controle de acesso e auditoria.

Em outras palavras: quem resolve nomes de domínio de verdade é o **PowerDNS**. O SanchezDNS apenas conversa com a **API HTTP** desse PowerDNS para criar zonas, editar registros, ativar DNSSEC e ler estatísticas — e, ao redor disso, adiciona login, permissões por zona, automações (como o PTR reverso) e logs.

```
┌──────────────┐        HTTPS         ┌───────────────────────────────────────────┐
│   Navegador  │  ───────────────▶    │            Container SanchezDNS            │
│  (usuário)   │                      │                                           │
└──────────────┘                      │   Caddy (:80)  ── reverse proxy           │
                                      │     ├── /go/*  ▶ API Go / Gin (:8080)     │
                                      │     └── /*     ▶ Nuxt SSR / Bun (:3000)   │
                                      └───────────────┬───────────────────────────┘
                                                      │
              ┌───────────────┬───────────────┬───────┴────────┬──────────────────┐
              ▼               ▼               ▼                ▼                  ▼
        ┌──────────┐   ┌──────────┐    ┌──────────┐    ┌──────────────┐   ┌──────────────┐
        │ MongoDB  │   │  Redis   │    │ S3/RustFS│    │  PowerDNS    │   │  (Geo/IP)    │
        │ dados    │   │  cache   │    │  fotos   │    │  API :8081   │   │  externo     │
        └──────────┘   └──────────┘    └──────────┘    └──────────────┘   └──────────────┘
```

## As camadas e por que cada uma existe

### Frontend — Nuxt 4 (Vue 3)

A interface é uma aplicação **Nuxt 4** com **Nuxt UI 4**. O Nuxt roda em modo **SSR (Server-Side Rendering)** sobre o runtime **Bun**, o que traz dois benefícios importantes para um painel administrativo:

- A primeira renderização já chega pronta do servidor, então a tela de login e as listagens aparecem rápido, sem "tela branca" enquanto o JavaScript carrega.
- O servidor Nuxt (camada Nitro) funciona como um **intermediário** entre o navegador e a API Go. Ele repassa os cookies e cabeçalhos de proxy (`x-forwarded-for`, `user-agent`, etc.), o que permite ao backend saber o IP real e o navegador do usuário mesmo estando atrás de um proxy.

**Por que Nuxt e não um SPA puro (só Vue no navegador)?** Porque um painel de DNS lida com sessões em cookie `HttpOnly` (que o JavaScript do navegador não pode ler por segurança). Com SSR, o próprio servidor Nuxt encaminha o cookie para a API, e o token de sessão nunca precisa ficar exposto ao código do cliente.

### Backend — Go com Gin

Toda a lógica de negócio vive em um serviço **Go** usando o framework **Gin**. É esse serviço que:

- valida sessões e permissões;
- fala com o **MongoDB** (cadastros, sessões, permissões, logs, solicitações);
- fala com o **Redis** (cache de sessão e de nível de acesso);
- fala com a **API do PowerDNS** (zonas, registros, DNSSEC, estatísticas);
- fala com o **S3** (fotos de perfil).

**Por que Go?** Três razões que importam para um TCC:

1. **Concorrência barata.** Operações como "criar registro A e, em paralelo, garantir o PTR reverso" ou "gravar um log sem travar a resposta ao usuário" usam _goroutines_ (`go logs.Insert(...)`). Isso mantém a resposta rápida.
2. **Binário único, sem runtime.** O Go compila para um executável estático (`CGO_ENABLED=0`), o que deixa a imagem Docker pequena e simples — sem interpretador, sem dependências de sistema.
3. **Tipagem forte na fronteira com o PowerDNS.** Os corpos JSON trocados com a API são modelados como _structs_ Go, o que reduz erros ao montar RRsets e payloads de DNSSEC.

O backend escuta na porta **8080** e **não é exposto diretamente**: só é alcançado através do Caddy, sob o prefixo `/go`.

### Proxy reverso — Caddy

Dentro do container existe um **Caddy** que unifica frontend e backend em uma única porta (**80**):

```
handle_path /go/*  →  reverse_proxy 127.0.0.1:8080   (API Go, prefixo removido)
handle             →  reverse_proxy 127.0.0.1:3000   (Nuxt SSR)
```

**Por que ter um proxy interno em vez de expor duas portas?** Porque assim o navegador enxerga **uma única origem**. O frontend chama `/go/records`, `/go/zones`, `/go/login`, e o Caddy entrega essas requisições ao serviço Go removendo o prefixo `/go`. Do ponto de vista do navegador, tudo vem do mesmo host — o que simplifica cookies (mesma origem, sem CORS) e o deploy (uma porta só para publicar).

Em desenvolvimento, esse mesmo papel é feito pelo `devProxy` do Nuxt (`/go` → `http://localhost:8080`), então o código do frontend é idêntico em dev e em produção.

### MongoDB — a "verdade" do aplicativo

O **MongoDB** guarda tudo o que é do domínio do _aplicativo_ (não do DNS em si):

| Coleção        | Guarda                                                       |
| -------------- | ------------------------------------------------------------ |
| `cadastros`    | contas aprovadas (nome, email, senha em bcrypt, foto, nível) |
| `sessions`     | sessões ativas (SID, email, IP, navegador, última atividade) |
| `users`        | permissões por zona (listas de `leitura` e `escrita`)        |
| `solicitacoes` | pedidos de cadastro aguardando aprovação                     |
| `logs`         | trilha de auditoria das operações                            |

**Por que MongoDB?** Os documentos aqui têm formato flexível e crescem por listas (por exemplo, a lista de emails com permissão de escrita por zona). Um modelo orientado a documentos encaixa bem nesse formato sem exigir tabelas de junção. Vale notar a separação de responsabilidades: **zonas e registros DNS _não_ ficam no MongoDB** — a fonte da verdade do DNS é o próprio PowerDNS. O Mongo só cuida de quem pode fazer o quê e do histórico.

### Redis — cache de sessão e de permissão

O **Redis** é uma camada de **cache** na frente do MongoDB para as duas leituras mais frequentes do sistema: "essa sessão é válida?" e "qual o nível desse usuário?".

A cada requisição protegida o backend precisa responder a essas perguntas. Consultar o MongoDB toda vez seria desperdício, então o Redis guarda a sessão (via `JSON.GET`, chave `sanchezdns:session:<sid>`) e o nível (`sanchezdns:level:<email>`). Se o Redis tiver o dado, responde na hora; se não tiver, o backend cai para o MongoDB e segue funcionando.

O detalhe importante — **por que Redis, e por que a sessão continua sendo do lado do servidor** — está detalhado na página [Autenticação e Sessões](/authentication). Ali explico JWT vs. sessão e stateless vs. stateful, que é o coração desta decisão.

### S3 compatível (RustFS/MinIO) — arquivos

As **fotos de perfil** não vão para o banco nem para o disco do container. Vão para um **storage compatível com o protocolo S3** (no compose, o `rustfs`, que também pode ser MinIO ou a AWS S3 real).

**Por que S3 e não o disco local?** Porque o container do SanchezDNS é **descartável**: pode ser recriado, atualizado ou escalado a qualquer momento. Guardar arquivos dentro dele significaria perdê-los na próxima atualização. Delegar isso a um serviço de objetos mantém o container **sem estado (stateless)** — todo o estado persistente vive em serviços dedicados (Mongo, Redis, S3, PowerDNS).

### PowerDNS — o servidor DNS de verdade

O **PowerDNS Authoritative** é quem responde às consultas DNS na internet. O SanchezDNS o trata como a **fonte da verdade das zonas e registros** e conversa com ele exclusivamente pela API HTTP. Tudo sobre essa integração — configuração do servidor, o que é um RRset, DNSSEC, o formato do SOA — está na página [PowerDNS a fundo](/powerdns).

## O ciclo de vida de uma requisição

Para amarrar as peças, vale seguir o caminho de uma ação típica — **"adicionar um registro A"**:

1. **Navegador → Nuxt.** O usuário preenche o formulário e o frontend faz `PUT /go/records`. A requisição chega ao servidor Nuxt (SSR), que a repassa com o cookie de sessão.
2. **Caddy → Go.** O Caddy vê o prefixo `/go`, remove-o e entrega `PUT /records` à API Go na porta 8080.
3. **Middleware de sessão.** Antes do controlador rodar, o middleware `ValidateSession` lê o cookie `sanchezdns_session_id`, valida a sessão (Redis, com fallback no Mongo) e injeta `email` e `level` no contexto da requisição.
4. **Autorização por zona.** O controlador chama `users.PodeEscrever(ctx, zona, email, level)`: `admin` passa direto; `member` precisa estar na lista `escrita` daquela zona (consulta ao Mongo).
5. **Integração com o PowerDNS.** O backend normaliza o valor do registro, monta o **RRset** e envia um `PATCH` para a zona no PowerDNS. Em seguida, relê a zona e faz um segundo `PATCH` para alinhar os comentários (detalhe explicado em [Zonas e Registros](/zones)).
6. **Automação em segundo plano.** Como o tipo é `A`, o controlador chama `records.InsertReverso`, que procura a melhor zona `in-addr.arpa` e cria o PTR reverso correspondente. Se essa parte falhar, o erro vai para o log e a resposta segue como sucesso — o registro direto já foi gravado.
7. **Auditoria assíncrona.** Uma _goroutine_ grava um log no MongoDB (`insert_record`) sem atrasar a resposta.
8. **Resposta.** O backend devolve `200` com a mensagem de sucesso, que sobe de volta por Caddy → Nuxt → navegador, e a tabela é atualizada.

Esse fluxo mostra a filosofia do projeto: **o painel orquestra, o PowerDNS executa, e cada serviço de apoio (Mongo, Redis, S3) tem uma responsabilidade única e bem delimitada.**

## Por que essa separação vale a pena

- **Cada peça é substituível.** Trocar RustFS por AWS S3, ou Redis por outra instância, não muda o código de negócio — só variáveis de ambiente.
- **O container é sem estado.** Todo dado durável vive fora dele, então atualizar a aplicação é só subir uma imagem nova.
- **A fonte da verdade do DNS é única.** Nunca há risco de o painel e o servidor DNS "discordarem": o painel não guarda cópia das zonas, ele lê e escreve direto no PowerDNS.
- **Segurança concentrada.** Autenticação, permissão e auditoria ficam todas no backend Go; o frontend nunca decide sozinho o que pode ou não ser feito.
