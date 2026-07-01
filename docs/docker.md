# 🐳 A imagem Docker

Esta página explica **como a imagem do SanchezDNS é construída e por que ela é montada dessa maneira**. O `Dockerfile` do projeto usa uma técnica chamada **multi-stage build** (construção em múltiplos estágios) e roda **três processos dentro de um único container**. Ambas as escolhas têm motivos concretos, detalhados abaixo.

## O problema que a imagem resolve

O SanchezDNS é composto por duas linguagens diferentes:

- um **frontend** em Nuxt/Vue, que precisa do Bun e do Node ecosystem para ser compilado;
- um **backend** em Go, que precisa do compilador Go.

Se colocássemos tudo isso na imagem final, ela carregaria dois toolchains inteiros (centenas de MB) que só servem para _compilar_, não para _rodar_. O multi-stage build resolve isso: **compila em estágios "descartáveis" e copia apenas os artefatos prontos para o estágio final.**

## Os três estágios de build

### Estágio 1 — `nuxt-builder` (compila o frontend)

```dockerfile
FROM oven/bun:1-debian AS nuxt-builder
...
COPY ./package.json ./bun.lock ./.npmrc ./nuxt.config.ts ./tsconfig.json ./
RUN bun install --ci
COPY ./app ./app
COPY ./public ./public
RUN bun run app:build
```

Aqui o Bun instala as dependências e roda `nuxt build`, gerando a pasta `.output` (o servidor SSR pronto para produção).

**Detalhe de otimização:** os arquivos de manifesto (`package.json`, `bun.lock`) são copiados **antes** do código-fonte. Isso aproveita o **cache de camadas** do Docker: enquanto as dependências não mudam, o `bun install` não é refeito a cada build, mesmo que você altere um componente Vue. O `--mount=type=cache` no `bun install` reforça isso guardando o cache de download entre builds.

### Estágio 2 — `go-builder` (compila o backend)

```dockerfile
FROM golang:1.26.4-bookworm AS go-builder
...
COPY ./api/go.mod ./api/go.sum ./
RUN go mod download
COPY ./api ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /server/api .
```

O backend Go é compilado com três flags que importam:

- **`CGO_ENABLED=0`** — gera um binário **estático**, sem dependência de bibliotecas C do sistema. É o que permite rodá-lo em qualquer imagem base sem instalar nada.
- **`-trimpath`** — remove os caminhos absolutos de compilação do binário (mais limpo e reproduzível).
- **`-ldflags "-s -w"`** — remove a tabela de símbolos e informações de debug, deixando o binário **menor**.

Assim como no estágio do Nuxt, o `go.mod`/`go.sum` são baixados antes do código, para aproveitar cache.

### Estágio 3 — `runner` (a imagem final que roda)

```dockerfile
FROM oven/bun:1-debian AS runner
RUN apt-get install -y --no-install-recommends ca-certificates curl tini bash
COPY --from=caddy:2-alpine /usr/bin/caddy /usr/bin/caddy
COPY --from=nuxt-builder /main/.output ./.output
COPY --from=go-builder /server/api ./api
```

A imagem final parte do `oven/bun` (necessário porque o servidor SSR do Nuxt roda com Bun) e **copia apenas os artefatos**:

- o binário do **Caddy**, direto da imagem oficial `caddy:2-alpine` (não precisamos instalar o Caddy, só pegar o executável);
- a pasta `.output` do Nuxt, vinda do estágio 1;
- o binário `api` do Go, vindo do estágio 2.

Nenhum compilador Go, nenhum `node_modules` de build, nenhum código-fonte entram na imagem final. Só o que é preciso para **executar**.

## Por que três processos em um container?

Esta é uma decisão que merece explicação, porque a "regra clássica" do Docker é _um processo por container_. O SanchezDNS roda **três** — API Go, servidor Nuxt e Caddy — através de um `entrypoint.sh`:

```bash
/main/api &                                    # backend Go em :8080
bun /main/.output/server/index.mjs &           # Nuxt SSR em :3000
caddy run --config /etc/caddy/Caddyfile &      # proxy em :80
wait -n
```

**Por que juntar tudo?** Porque frontend e backend são **duas metades da mesma aplicação** e sempre sobem juntos, na mesma versão. Distribuir uma imagem única simplifica enormemente o uso: quem quer rodar o SanchezDNS faz `docker compose up` e recebe o painel inteiro funcionando, sem precisar orquestrar dois ou três containers da própria aplicação e configurar a rede entre eles.

Para fazer isso corretamente, o entrypoint cuida de dois problemas que "vários processos em um container" normalmente causam:

1. **Encerramento coordenado.** Um `trap` captura os sinais `TERM`/`INT` e o `wait -n` faz o script acordar assim que **qualquer** um dos três processos morrer. Se um cair, o container inteiro é derrubado (`kill -TERM 0`) e o Docker o reinicia limpo. Não existe estado "meio vivo".
2. **Reaping de processos zumbi.** O `ENTRYPOINT` usa o **`tini`** como PID 1. Ele adota e "colhe" processos filhos que terminam, evitando o acúmulo de zumbis — problema conhecido quando se roda múltiplos processos sob um shell.

### O Caddyfile embutido

```text
:80
handle_path /go/* {
    reverse_proxy 127.0.0.1:8080
}
handle {
    reverse_proxy 127.0.0.1:3000
}
```

O Caddy é o que dá a **origem única** descrita na página de [Arquitetura](/architecture): tudo que começa com `/go/` vai para a API Go (com o prefixo removido pelo `handle_path`), e o resto vai para o Nuxt. A diretiva `trusted_proxies` garante que o backend confie nos cabeçalhos `X-Forwarded-For` das redes internas do Docker, para obter o IP real do cliente.

## Healthcheck

```dockerfile
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
    CMD curl -fsS http://127.0.0.1:80/go/healthcheck || exit 1
```

A cada 30 segundos o Docker chama `/go/healthcheck`, que passa pelo Caddy e chega ao backend Go. Isso testa a **cadeia inteira de uma vez**: se o Caddy caiu, ou a API Go não responde, o healthcheck falha e o orquestrador sabe que o container está insalubre. O `--start-period=20s` dá um tempo de graça no boot, antes de começar a contar falhas.

## O `docker-compose.yaml`

O compose descreve o **stack completo** de execução — a aplicação mais os serviços de apoio de que ela depende:

```yaml
services:
  sanchezdns: # a imagem descrita acima
    depends_on: [mongo, redis, s3, s3-create-bucket]
    env_file: [.env]
  mongo: # banco de dados da aplicação
  redis: # cache de sessão e nível
  s3: # storage compatível S3 (RustFS) para fotos
  s3-create-bucket: # job que cria o bucket na primeira subida
```

Alguns pontos que valem justificativa:

- **`depends_on`** garante a ordem de inicialização: o SanchezDNS só sobe depois de Mongo, Redis e S3 estarem de pé.
- **`s3-create-bucket`** é um container efêmero (imagem `minio/mc`) que roda uma vez, cria o bucket configurado em `FS_BUCKET` e libera leitura pública das fotos, depois encerra (`exit 0`). É a forma de ter um bucket pronto sem passo manual.
- **`volumes`** (`mongo_data`, `redis_data`, `s3_data`) mantêm os dados persistentes **fora** dos containers. É o que permite recriar qualquer serviço sem perder dados — coerente com a ideia de containers descartáveis.
- **`restart: unless-stopped`** faz os serviços voltarem sozinhos após uma falha ou reboot da máquina.
- **`AWS_REQUEST_CHECKSUM_CALCULATION=WHEN_REQUIRED`** ajusta o comportamento do SDK da AWS para ser compatível com storages S3 alternativos (como o RustFS), que não implementam todos os checksums da AWS.

> **Nota:** o `docker-compose.yaml` do repositório **não inclui um serviço PowerDNS**. Isso é intencional — o PowerDNS costuma ser uma infraestrutura existente e externa, apontada pela variável `DNS_HOST`. A página [PowerDNS a fundo](/powerdns) mostra como configurá-lo.

## Resumo das decisões

| Decisão                               | Motivo                                                                         |
| ------------------------------------- | ------------------------------------------------------------------------------ |
| Multi-stage build                     | Imagem final sem toolchains de compilação — menor e mais segura                |
| Copiar Caddy de outra imagem          | Ter o proxy sem instalá-lo manualmente                                         |
| `CGO_ENABLED=0` + `-ldflags "-s -w"`  | Binário Go estático e enxuto                                                   |
| Três processos + `tini`               | Distribuir o painel como uma imagem única, com encerramento e reaping corretos |
| Healthcheck via `/go/healthcheck`     | Testar Caddy + API numa só verificação                                         |
| Estado só em Mongo/Redis/S3 (volumes) | Container descartável, atualização sem perda de dados                          |
