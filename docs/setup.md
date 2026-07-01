# ⚙️ Instalação

Este guia descreve como subir o SanchezDNS, tanto em Docker (recomendado) quanto em desenvolvimento local. Para o detalhamento das variáveis, veja [Configuração](/configuration); para preparar o servidor DNS, veja [PowerDNS a fundo](/powerdns).

## Pré-requisitos

- **Docker** e **Docker Compose**.
- Uma instância **PowerDNS Authoritative** com a **API habilitada** (fora do compose — ver nota abaixo).
- Os serviços de apoio, que o próprio compose provê: **MongoDB**, **Redis** e um **storage S3 compatível**.

> **Nota:** o `docker-compose.yaml` do projeto sobe a aplicação e seus serviços de apoio (Mongo, Redis, S3), mas **não** o PowerDNS — ele é tratado como infraestrutura externa, apontada por `DNS_HOST`. Você pode rodar o PowerDNS em outra máquina, outro container ou usar um já existente.

## Passo 1 — Criar o `.env`

```bash
cp .env.example .env
```

Edite o `.env` com os valores do seu ambiente. Os campos mais importantes de acertar de cara são os do **PowerDNS** (`DNS_HOST`, `DNS_API_KEY`, `DNS_SERVER_ID`), porque é a integração que dá sentido ao painel. Veja [Configuração](/configuration) para a explicação de cada variável.

## Passo 2 — Obter o `docker-compose.yaml`

Se ainda não tiver o arquivo, baixe-o do repositório:

```bash
curl -L -o docker-compose.yaml https://raw.githubusercontent.com/rafinhacuri/sanchezdns/main/docker-compose.yaml
# ou
wget -O docker-compose.yaml https://raw.githubusercontent.com/rafinhacuri/sanchezdns/main/docker-compose.yaml
```

## Passo 3 — Subir o stack

```bash
docker compose pull
docker compose up -d --force-recreate
```

O stack cria:

- `sanchezdns` — a aplicação (frontend + backend + proxy Caddy na mesma imagem);
- `mongo` na porta `27017` — banco de dados da aplicação;
- `redis` na porta `6379` — cache de sessão e nível;
- `s3` na porta `9000` (console em `9001`) — storage de fotos;
- `s3-create-bucket` — job efêmero que cria o bucket e encerra.

A composição interna dessa imagem (multi-stage build, os três processos, o healthcheck) está explicada em [A imagem Docker](/docker).

## Passo 4 — Abrir a interface

Acesse **`http://localhost:3000`** no navegador. O **primeiro cadastro aprovado vira administrador**; os seguintes entram como solicitação até aprovação manual (ver [Usuários](/users)).

## Desenvolvimento local

Para trabalhar no código sem Docker, o repositório tem scripts (via `bun`):

- `bun app:dev` — sobe o frontend Nuxt.
- `bun api:dev` — sobe a API Go (usando `air` para hot-reload).
- `bun dev` — sobe os dois juntos, lado a lado.

Em desenvolvimento, o Nuxt faz proxy de `/go` para o backend Go em `http://localhost:8080` (via `devProxy` no `nuxt.config.ts`), reproduzindo o papel que o Caddy exerce em produção. Assim o código do frontend é idêntico nos dois modos.

Você ainda precisará de MongoDB, Redis, um storage S3 e um PowerDNS acessíveis — o mais simples é subir os serviços de apoio pelo compose e rodar só a aplicação localmente.

## Validação rápida

Confira se os containers estão de pé:

```bash
docker compose ps
```

E teste a saúde da aplicação (a rota passa pelo Caddy e chega ao backend Go):

```bash
curl http://localhost:3000/go/healthcheck
```

Se o healthcheck responder, a cadeia Caddy → API Go está funcionando. Se as telas de zona/estatística derem erro, o mais provável é divergência nas variáveis do PowerDNS — reveja [Configuração](/configuration#powerdns).
