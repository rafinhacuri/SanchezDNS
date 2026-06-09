# ⚙️ Instalação

Este guia descreve a forma atual de subir o projeto, tanto em Docker quanto em desenvolvimento local.

## Pré-requisitos

- Docker e Docker Compose.
- Uma instância PowerDNS Authoritative com API habilitada.
- MongoDB.
- Redis.
- Um storage S3 compatível para fotos de perfil.

## Subindo com Docker

O caminho mais simples é usar o `docker-compose.yaml` da raiz do repositório.

```bash
docker compose pull
docker compose up -d --force-recreate
```

O stack cria:

- `sanchezdns` na porta `3000`
- `mongo` na porta `27017`
- `redis` na porta `6379`
- `s3` na porta `9000` e console na `9001`

Depois da subida, abra `http://localhost:3000`.

## Variáveis de ambiente

Crie um arquivo `.env` na raiz com os valores do seu ambiente:

```bash
NUXT_PUBLIC_PRODUCTION=true
NUXT_PUBLIC_SITE_URL=https://dns.exemplo.com
MONGO_URL=mongodb://mongo:27017
REDIS_URL=redis://redis:6379
DNS_HOST=http://powerdns:8081
DNS_API_KEY=chave-da-api
DNS_SERVER_ID=localhost
FS_USERNAME=access-key
FS_PASSWORD=secret-key
FS_BUCKET=sanchezdns
FS_ENDPOINT=http://s3:9000
```

## PowerDNS

No servidor PowerDNS, a API precisa estar habilitada. Um exemplo mínimo:

```ini
api=yes
api-key=chave-da-api
webserver=yes
webserver-address=0.0.0.0
webserver-port=8081
server-id=localhost
```

## Desenvolvimento local

O repositório também possui scripts para rodar frontend e backend separadamente:

- `bun app:dev` para o Nuxt.
- `bun api:dev` para a API Go.
- `bun dev` para subir os dois juntos.

## Validação rápida

Depois de subir o stack, confira se os containers estão ativos:

```bash
docker compose ps
```

Se quiser uma checagem da API, use:

```bash
curl http://localhost:3000/server/healthcheck
```

---
