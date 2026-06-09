<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=venom&height=300&color=2FC851&text=SanchezDNS&section=header&fontColor=ffffff"/>
</p>

## SanchezDNS

SanchezDNS é uma aplicação web para administrar um ambiente PowerDNS Authoritative com foco em operação diária, acesso por zona e automações que reduzem o trabalho manual.

### O que ele faz hoje

- Autenticação por sessão em cookie.
- Cadastro inicial com promoção automática para `admin`.
- Fluxo de solicitação de cadastro com aprovação manual pelo administrador.
- Gestão de zonas normais, reversas IPv4 e reversas IPv6.
- Gestão de registros DNS com atualização automática de PTR reverso para `A` e `AAAA`.
- Controle de acesso por zona com permissões de leitura e escrita.
- Logs de auditoria para operações de zonas, registros, usuários e cadastros.
- Upload de fotos de perfil em storage S3 compatível.

### Stack atual

- Backend em Go com Gin.
- Frontend em Nuxt 4 com Nuxt UI 4.
- MongoDB para cadastros, sessões, permissões, solicitações e logs.
- Redis para cache de sessão e nível de acesso.
- PowerDNS Authoritative API como fonte de verdade das zonas e registros.
- S3 compatível para arquivos de perfil.

### Funcionalidades principais

- Zonas `normal`, `reverse` e `reverse-ipv6`.
- Criação e edição de SOA.
- DNSSEC ativado nas zonas criadas pelo sistema.
- Registros `A`, `AAAA`, `ALIAS`, `CAA`, `CNAME`, `HTTPS`, `MX`, `NS`, `PTR`, `TXT`, `SRV` e `TLSA`.
- Painel de estatísticas com zonas, registros, uptime e tráfego UDP/TCP.
- Página de logs com busca e paginação.
- Página administrativa para solicitações e cadastros.

### Executando com Docker

**Step 1 — Create the `.env` file**

Copy the example environment file and configure your variables:

```bash
cp .env.example .env
```

Edit the `.env` file with your PowerDNS configuration and other required variables.

**Step 2 — Get the docker-compose.yaml**

Choose one of the options below to download the configuration file:

🔽 Using curl

```bash
curl -L -o docker-compose.yaml https://raw.githubusercontent.com/rafinhacuri/sanchezdns/main/docker-compose.yaml
```

🔽 Using wget

```bash
wget -O docker-compose.yaml https://raw.githubusercontent.com/rafinhacuri/sanchezdns/main/docker-compose.yaml
```

Alternatively, copy it directly from the repository.

**Step 3 — Start the services**

```bash
docker compose pull
docker compose up -d --force-recreate
```

**Step 4 — Open the interface**

Navigate to `http://localhost:3000` in your browser.

### Variáveis de ambiente mais importantes

- `NUXT_PUBLIC_PRODUCTION`
- `NUXT_PUBLIC_SITE_URL`
- `MONGO_URL`
- `REDIS_URL`
- `DNS_HOST`
- `DNS_API_KEY`
- `DNS_SERVER_ID`
- `FS_USERNAME`
- `FS_PASSWORD`
- `FS_BUCKET`
- `FS_ENDPOINT`

### Observações de operação

- O sistema trabalha com uma única instância PowerDNS definida por `DNS_HOST` e `DNS_SERVER_ID`.
- Quando uma zona é criada, o backend ativa a zona e gera a chave KSK de DNSSEC.
- Quando um registro `A` ou `AAAA` muda, o sistema tenta manter o PTR reverso correspondente em sincronia.
- O primeiro cadastro aprovado vira administrador; os próximos entram como solicitações até aprovação.

### Contribuição

Issues e pull requests são bem-vindos.

### Licença

Projeto licenciado sob a [MIT License](LICENSE).
