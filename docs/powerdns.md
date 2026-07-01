# 🌐 PowerDNS a fundo

Esta é a página mais importante para entender o coração técnico do SanchezDNS. Aqui explico **o que é o PowerDNS, como configurá-lo, o que são RRsets, como o DNSSEC funciona** e, sobretudo, **como e por que o SanchezDNS usa cada endpoint da API**.

## O que é um servidor DNS autoritativo

O DNS (Domain Name System) é a "lista telefônica" da internet: ele traduz nomes (`exemplo.com`) em endereços (`192.0.2.10`). Existem dois papéis principais:

- **Resolvedor (recursivo):** o servidor que o seu computador pergunta e que sai atrás da resposta pela internet.
- **Autoritativo:** o servidor que **detém a verdade** sobre um domínio. Quando alguém pergunta "qual o IP de `exemplo.com`?", é o servidor autoritativo do `exemplo.com` que dá a resposta final.

O **PowerDNS Authoritative Server** é um servidor do segundo tipo. O SanchezDNS existe justamente para administrar esse servidor autoritativo — criar e editar as zonas e registros que ele serve.

**Por que PowerDNS e não BIND (o servidor DNS mais clássico)?** Porque o PowerDNS foi projetado com duas características que o SanchezDNS depende diretamente:

1. **Uma API HTTP/JSON de primeira classe.** Enquanto o BIND é tradicionalmente configurado por arquivos de zona em texto, o PowerDNS expõe uma API REST completa para criar zonas, editar registros e gerenciar DNSSEC. Sem isso, um painel web como o SanchezDNS teria que manipular arquivos e recarregar o servidor manualmente.
2. **Backends plugáveis.** O PowerDNS pode guardar as zonas em bancos de dados (SQLite, MySQL, PostgreSQL) em vez de arquivos, o que combina bem com automação e DNSSEC.

## Conceitos fundamentais

### Zona

Uma **zona** é o conjunto de registros de um domínio sob a mesma autoridade — por exemplo, tudo abaixo de `exemplo.com`. No PowerDNS, toda zona é identificada pelo nome **com ponto final** (`exemplo.com.`), o que representa a raiz absoluta do DNS. Por isso o código do SanchezDNS sempre normaliza os nomes acrescentando esse ponto:

```go
domain = strings.TrimSuffix(domain, ".")
domainWithDot := domain + "."   // "exemplo.com."
```

### Registro (Resource Record) e o SOA

Cada entrada dentro da zona é um **registro** (RR — Resource Record). Todo registro tem um **nome**, um **tipo** (`A`, `MX`, `TXT`...), um **TTL** (tempo em segundos que pode ser cacheado) e um **conteúdo**.

Um registro é especial e obrigatório em toda zona: o **SOA (Start of Authority)**. Ele descreve os parâmetros administrativos da zona. Seu conteúdo é uma linha com sete campos:

```
mname   rname          serial refresh retry expire minimum
ns1.exemplo.com. admin.exemplo.com. 1 3600 600 604800 3600
```

| Campo     | Significado                                                   | Por que importa                              |
| --------- | ------------------------------------------------------------- | -------------------------------------------- |
| `mname`   | servidor primário da zona                                     | quem é a autoridade principal                |
| `rname`   | email do responsável (o `@` vira `.`)                         | contato administrativo                       |
| `serial`  | número de versão da zona                                      | secundários usam para saber se houve mudança |
| `refresh` | de quanto em quanto tempo um secundário verifica atualizações | controla propagação                          |
| `retry`   | espera antes de tentar de novo se o refresh falhar            | resiliência                                  |
| `expire`  | quando um secundário desiste de dados velhos                  | evita servir dados obsoletos                 |
| `minimum` | TTL do cache negativo (respostas "não existe")                | quanto tempo um "NXDOMAIN" é lembrado        |

No SanchezDNS, o formulário de criação/edição de zona pede exatamente esses campos (`startOfAuthority`, `email`, `refresh`, `retry`, `expire`, `negativeCacheTtl`) e o backend os monta assim:

```go
Content: fmt.Sprintf("%s. %s. 1 %d %d %d %d", soaMname, soaRname, refresh, retry, expire, negativeCacheTtl)
```

Repare que o **serial é gravado como `1`**. Isso é proposital: ao criar a zona, o backend define `soa_edit_api: "DEFAULT"`, o que faz o **próprio PowerDNS reescrever o serial automaticamente** a cada alteração feita pela API (tipicamente no formato de data `AAAAMMDDnn`). Ou seja, o painel não precisa gerenciar versionamento de serial na mão — delega isso ao servidor. É a razão de o valor inicial ser um placeholder.

### RRset — o conceito central da API

Aqui está o conceito que a API do PowerDNS gira em torno e que confunde quem vem do modelo de "arquivo de zona": o **RRset (Resource Record Set)**.

Um **RRset é o conjunto de todos os registros que compartilham o mesmo nome e o mesmo tipo**. Por exemplo, se `www.exemplo.com` aponta para dois IPs:

```
www.exemplo.com.  A  192.0.2.10
www.exemplo.com.  A  192.0.2.11
```

isso **não são dois objetos independentes** para o PowerDNS — é **um único RRset** (`www.exemplo.com` / tipo `A`) que contém dois `records`. O TTL e os comentários pertencem ao RRset como um todo.

**Por que isso importa tanto?** Porque a API do PowerDNS **não deixa você editar um registro isolado**. Você sempre envia o **RRset inteiro** e diz o que fazer com ele através de um `changetype`:

- **`REPLACE`** — substitui todo o conjunto de registros daquele nome+tipo pelo que você mandou (cria se não existir);
- **`DELETE`** — remove o RRset inteiro daquele nome+tipo.

Essa é a razão de o SanchezDNS, para **adicionar** um registro A a um nome que já tem outros, primeiro **ler** os registros existentes, **juntar** o novo e reenviar o conjunto completo com `REPLACE`. Se ele mandasse só o registro novo com `REPLACE`, apagaria os anteriores. Esse cuidado está no `buildRecordsStep1`:

```go
// copia os records que já existem no RRset...
for _, rec := range existingRR.Records {
    recordsStep1 = append(recordsStep1, Record{Content: rec.Content, Disabled: rec.Disabled})
}
// ...e acrescenta o novo
recordsStep1 = append(recordsStep1, Record{Content: res, Disabled: false})
```

A struct que o SanchezDNS envia reflete esse modelo:

```go
type PDNSRRSetChange struct {
    Name       string     // "www.exemplo.com."
    Type       string     // "A"
    TTL        *int       // 3600
    ChangeType string     // "REPLACE" ou "DELETE"
    Records    []Record   // a lista completa de registros do conjunto
    Comments   []Comment  // comentários, alinhados por posição
}
```

### Um RRset de verdade

Para tornar tudo isso concreto, veja o retorno **real** da API do PowerDNS para uma zona de teste. A chamada é a mesma que o SanchezDNS faz para ler uma zona:

```bash
curl -s \
  -H 'X-API-Key: SUA_API_KEY' \
  -H 'Accept: application/json' \
  'http://127.0.0.1:8081/api/v1/servers/localhost/zones/teste.com.'
```

A resposta descreve a zona inteira, e dentro dela o array `rrsets` — cada item é um RRset (um nome + um tipo). Recorte do JSON retornado:

```json
{
  "name": "teste.com.",
  "kind": "Native",
  "dnssec": true,
  "serial": 2026060302,
  "edited_serial": 2026060302,
  "soa_edit_api": "DEFAULT",
  "rrsets": [
    {
      "name": "rafael.teste.com.",
      "type": "A",
      "ttl": 3600,
      "records": [
        { "content": "152.20.120.224", "disabled": false }
      ],
      "comments": [
        { "account": "", "content": "teste", "modified_at": 1780500082 }
      ]
    },
    {
      "name": "teste.com.",
      "type": "HTTPS",
      "ttl": 3600,
      "records": [
        { "content": "3 teste.com. alpn=h2", "disabled": false }
      ],
      "comments": [
        { "account": "", "content": "dscscd", "modified_at": 1771535023 }
      ]
    },
    {
      "name": "teste.com.",
      "type": "SRV",
      "ttl": 3600,
      "records": [
        { "content": "60 4 80 teste.com.", "disabled": false }
      ]
    },
    {
      "name": "teste.com.",
      "type": "MX",
      "ttl": 3600,
      "records": [
        { "content": "50 mo.mo.com.", "disabled": false }
      ]
    },
    {
      "name": "teste.com.",
      "type": "SOA",
      "ttl": 3600,
      "records": [
        { "content": "teste.com. teste.com. 2026060302 3600 600 604800 86400", "disabled": false }
      ]
    }
  ]
}
```

Repare em como esse retorno conecta tudo o que foi explicado até aqui:

- **Cada objeto do array `rrsets` é um RRset** — a combinação única de `name` + `type`. Note que `teste.com.` aparece **quatro vezes**, mas são RRsets **diferentes** (`HTTPS`, `SRV`, `MX`, `SOA`): o que os separa é o `type`. O nome sozinho não identifica um RRset; é o par nome+tipo.
- **O `ttl` e os `comments` pertencem ao RRset como um todo**, não a cada registro. É por isso que, ao editar, o SanchezDNS reenvia o conjunto inteiro (ver o padrão de dois PATCHes adiante) — e por que os comentários são **posicionais**, alinhados à ordem dos `records`.
- **O array `records`** é a lista de valores daquele conjunto. Aqui cada RRset tem só um registro, mas se `rafael.teste.com.` tivesse dois IPs, apareceriam dois objetos dentro do mesmo `records` — um único RRset com dois registros.
- **O `content` já vem no formato normalizado** que o SanchezDNS gravou: o `MX` traz a prioridade embutida (`50 mo.mo.com.`), o `SRV` traz `prioridade peso porta alvo.` (`60 4 80 teste.com.`), o `HTTPS` traz `prioridade alvo parâmetros` (`3 teste.com. alpn=h2`), e todos os nomes de destino terminam com ponto. Isso é exatamente o resultado das regras de normalização descritas em [Zonas e Registros](/zones#normalizacao-por-que-o-valor-e-ajustado-antes-de-gravar).
- **O `SOA`** é apenas mais um RRset. Seu `content` tem os sete campos na ordem `mname rname serial refresh retry expire minimum` (`teste.com. teste.com. 2026060302 3600 600 604800 86400`). O `serial` — `2026060302` — está no formato de data (`AAAAMMDDnn`) porque a zona usa `soa_edit_api: "DEFAULT"`: foi o **próprio PowerDNS** que gerou esse número, não uma pessoa. É a garantia contra o problema do serial manual do BIND9 (ver [Por que SanchezDNS](/reason#bind9-poderoso-mas-manual-e-perigoso)).
- No topo, o campo **`"dnssec": true`** confirma que a zona está com DNSSEC ativo, e **`"kind": "Native"`** confirma o tipo de zona que o SanchezDNS cria.

> Os campos administrativos da zona (`masters`, `nsec3param`, `api_rectify`, etc.) foram omitidos do recorte por clareza. Uma zona criada pela interface do SanchezDNS ainda traria os RRsets de **DNSSEC** e o `NS`; esta zona de teste foi criada manualmente via `pdnsutil`, então mostra só os registros essenciais — o que a torna um exemplo limpo do conceito de RRset.

## Como o SanchezDNS usa a API — endpoint por endpoint

Todas as chamadas vão para `DNS_HOST` (a URL da API), autenticadas pelo cabeçalho **`X-API-Key`** com o valor de `DNS_API_KEY`, e usam o `DNS_SERVER_ID` (`localhost` por padrão) no caminho. O cliente HTTP tem timeout de 30s e 2 tentativas de retry.

### Criar uma zona — `POST /api/v1/servers/{id}/zones`

```go
zonePayload := pdnsCreateZoneRequest{
    Name:       "exemplo.com.",
    Kind:       "Native",
    SOAEditAPI: "DEFAULT",
}
```

O `Kind: "Native"` é uma decisão de projeto que vale explicar. O PowerDNS oferece três tipos de zona:

- **`Master`** — a zona é primária e notifica/transfere para secundários via protocolo DNS (AXFR).
- **`Slave`** — a zona é uma cópia recebida de um primário.
- **`Native`** — não há transferência DNS entre servidores; a **replicação, se existir, é feita pelo banco de dados** por baixo.

O SanchezDNS cria zonas como **`Native`** porque o modelo do projeto é de **uma instância PowerDNS operada pelo painel**. Se você precisar de réplicas, faz sentido replicar no nível do banco (o backend do PowerDNS) em vez de configurar AXFR — é mais simples e evita a complexidade de gerenciar relações master/slave pela interface. Isso é coerente com a decisão de remover o "painel de múltiplos servidores" das versões antigas.

### Ativar DNSSEC — `POST .../zones/{zona}/cryptokeys`

Logo após criar a zona, o backend ativa DNSSEC:

```go
dnssecPayload := pdnsCryptoKeyRequest{ Active: true, KeyType: "ksk" }
```

**O que é DNSSEC e por que ativar por padrão?** O DNS comum não tem autenticação: uma resposta pode ser forjada por um atacante no meio do caminho (cache poisoning). O **DNSSEC (DNS Security Extensions)** resolve isso assinando criptograficamente os registros da zona, de modo que o resolvedor consegue verificar que a resposta é autêntica e não foi adulterada.

O `KeyType: "ksk"` cria uma **KSK (Key Signing Key)** — a chave que assina as demais chaves da zona e cuja "impressão digital" (registro DS) é publicada na zona pai para estabelecer a cadeia de confiança. Ao ativar isso na criação, toda zona nasce **pronta para ser validada**, sem passo manual. Para que isso funcione, o backend do PowerDNS precisa suportar DNSSEC — no exemplo com SQLite, é o que a opção `gsqlite3-dnssec=yes` habilita (veja abaixo).

### Gravar o SOA inicial — `PATCH /api/v1/servers/{id}/zones/{zona}`

Todo `PATCH` numa zona carrega uma lista de RRsets. Na criação, o backend faz um `PATCH` com o RRset `SOA` (`REPLACE`) contendo a linha montada com os valores do formulário. O mesmo endpoint é reusado pelo `UpdateSoa` quando o administrador edita os parâmetros da zona depois.

### Adicionar/editar registros — `GET` + `PATCH` (o padrão de dois passos)

Adicionar um registro é a operação mais elaborada, e o motivo é o RRset + comentários. O fluxo em `InsertRecord` é:

1. **`GET` da zona** para descobrir se o RRset (nome+tipo) já existe e quais registros/comentários ele tem.
2. **`PATCH` nº 1** com o RRset completo (existentes + novo), `changetype REPLACE`.
3. **`GET` da zona de novo**, porque o PowerDNS pode **reordenar** ou **normalizar** os registros ao gravá-los.
4. **`PATCH` nº 2** reenviando o RRset já na ordem canônica do servidor, agora **com os comentários alinhados por posição** a cada registro.

**Por que dois `PATCH`?** Porque no PowerDNS os **comentários são posicionais**: o comentário `[0]` corresponde ao registro `[0]`, e assim por diante. Se o servidor reordenou os registros no passo 2, gravar os comentários "às cegas" alinharia o comentário no registro errado. Ao reler e só então gravar os comentários, o SanchezDNS garante que cada comentário fica no registro certo. É um detalhe sutil, mas é o que mantém a integridade da anotação de cada valor.

### Registro reverso (PTR) automático

Quando o registro é `A` ou `AAAA`, o backend chama `ensureReverseRecord` para manter o **PTR** (mapeamento reverso IP → nome) em sincronia. Essa lógica é detalhada na página de [Zonas e Registros](/zones), mas do ponto de vista da API ela: lista as zonas (`GET .../zones`), encontra por **maior sufixo em comum** a melhor zona reversa (`in-addr.arpa` para IPv4, `ip6.arpa` para IPv6), verifica se o PTR já existe e, se não, cria com um `PATCH`.

### Ler estatísticas — `GET .../statistics` e varredura de zonas

A página de estatísticas usa `GET /api/v1/servers/{id}/statistics` para métricas do servidor (uptime, `udp-queries`, `tcp-queries`) e depois percorre `GET .../zones` + `GET .../zones/{id}` de cada zona para **contar os registros somando o tamanho de cada RRset**. Detalhes em [Estatísticas](/statistics).

## Configurando o servidor PowerDNS

Toda a integração acima só funciona se o PowerDNS estiver configurado com a **API e o webserver habilitados**. Abaixo está uma configuração real de servidor primário com backend SQLite e DNSSEC, comentada campo a campo (`/etc/pdns/pdns.conf`):

```ini
# --- Papel do servidor ---
primary=yes
# allow-axfr-ips=<ip do secundário>     # libera transferência de zona para secundários
# also-notify=<ip do secundário>        # notifica secundários sobre mudanças

# --- Onde o PowerDNS escuta consultas DNS ---
local-address=127.0.0.1,<ip do host>
local-port=53

# --- Backend de armazenamento das zonas ---
launch=gsqlite3
gsqlite3-database=/var/lib/powerdns/pdns.sqlite3
gsqlite3-dnssec=yes

# --- Padrões e privacidade ---
default-ttl=3600
version-string=anonymous

# --- API HTTP usada pelo SanchezDNS ---
api=yes
api-key=<sua-chave-secreta>
webserver=yes
webserver-address=0.0.0.0
webserver-allow-from=0.0.0.0/0
webserver-port=8081
server-id=localhost
```

Explicação dos campos que mais importam:

- **`primary=yes`** — marca o servidor como primário (autoridade principal). As linhas comentadas `allow-axfr-ips` e `also-notify` só são necessárias se você tiver servidores **secundários** recebendo as zonas por AXFR. Como o SanchezDNS cria zonas `Native`, elas não dependem disso; os campos ficam disponíveis caso você opere secundários no futuro.
- **`local-address` / `local-port`** — os IPs e a porta (**53**, a porta padrão do DNS) onde o PowerDNS responde às consultas de resolução propriamente ditas. **Não confundir** com a porta da API.
- **`launch=gsqlite3`** — escolhe o **backend** de armazenamento. `gsqlite3` guarda todas as zonas e registros em um arquivo **SQLite**. É perfeito para um TCC ou ambiente pequeno: zero infraestrutura de banco, tudo em um arquivo. Para produção maior, trocaria-se por `gmysql` ou `gpgsql`.
- **`gsqlite3-database`** — o caminho do arquivo SQLite.
- **`gsqlite3-dnssec=yes`** — habilita o suporte a **DNSSEC** nesse backend. É **obrigatório** para que a ativação de DNSSEC feita pelo SanchezDNS (o `POST .../cryptokeys`) funcione. Sem isso, a criação de zona falharia no passo de DNSSEC.
- **`default-ttl=3600`** — TTL padrão (1h) quando um registro não especifica um.
- **`version-string=anonymous`** — esconde a versão do PowerDNS em respostas de diagnóstico, uma boa prática de segurança (não entregar de graça a versão para quem faz reconhecimento).
- **`api=yes` + `api-key`** — habilita a API REST e define a **chave secreta**. Esse valor precisa ser **idêntico** ao `DNS_API_KEY` do SanchezDNS; é ele que o backend envia no cabeçalho `X-API-Key`.
- **`webserver=yes` + `webserver-address` + `webserver-port=8081`** — sobem o servidor HTTP que hospeda a API. A porta **8081** é o que vai em `DNS_HOST` (ex.: `http://powerdns:8081`).
- **`webserver-allow-from`** — quais IPs podem falar com a API. `0.0.0.0/0` libera todos; **em produção, restrinja** isso à rede onde o SanchezDNS roda, porque quem alcança a API com a chave controla todo o seu DNS.
- **`server-id=localhost`** — o identificador do servidor que aparece no caminho da API (`/servers/localhost/...`). É exatamente o valor que vai em `DNS_SERVER_ID`.

### O elo entre a config do PowerDNS e as variáveis do SanchezDNS

| No `pdns.conf`                     | Na `.env` do SanchezDNS            | Precisa bater?            |
| ---------------------------------- | ---------------------------------- | ------------------------- |
| `webserver-address:webserver-port` | `DNS_HOST` (ex.: `http://ip:8081`) | sim — é o endereço da API |
| `api-key`                          | `DNS_API_KEY`                      | sim — mesma chave secreta |
| `server-id`                        | `DNS_SERVER_ID`                    | sim — mesmo identificador |

Se qualquer um dos três não bater, o SanchezDNS não consegue conversar com o PowerDNS (erro de conexão ou `401`/`404` nas rotas de zona).

## Resumo mental

- O **PowerDNS** é a autoridade do DNS; o **SanchezDNS** é o painel que fala com a API dele.
- A API trabalha com **RRsets** (nome + tipo + lista de registros), não com registros isolados — por isso o padrão é sempre "ler o conjunto, alterar, reenviar com `REPLACE`".
- Zonas nascem **`Native`** (uma instância operada pelo painel) e **com DNSSEC** (KSK ativa), desde que o backend suporte (`gsqlite3-dnssec=yes`).
- O **serial do SOA** é gerenciado pelo próprio PowerDNS (`soa_edit_api=DEFAULT`).
- Três valores precisam estar em sincronia entre o servidor e o painel: **host da API, chave e server-id**.
