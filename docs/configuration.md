# ⚙️ Configuração

Esta página reúne **todas as variáveis de ambiente** do SanchezDNS, o que cada uma faz e por que existe. É a referência para montar o arquivo `.env`. Para a configuração do servidor PowerDNS em si, veja [PowerDNS a fundo](/powerdns); para o modelo de sessão, veja [Autenticação](/authentication).

## O modelo: uma instância PowerDNS por ambiente

O SanchezDNS opera **uma única instância PowerDNS Authoritative** por ambiente. Não existe tela de múltiplas conexões nem "cluster" de servidores dentro da interface. Toda a integração é definida por três variáveis: `DNS_HOST`, `DNS_API_KEY` e `DNS_SERVER_ID`.

**Por que um só servidor?** Porque isso mantém o modelo mental simples e a fonte da verdade única (ver [Por que SanchezDNS](/reason)). Se você precisa de redundância, ela é resolvida na camada de infraestrutura (réplicas do backend do PowerDNS), não multiplicando conexões no painel.

## O arquivo `.env`

Copie o exemplo e ajuste aos valores do seu ambiente:

```bash
cp .env.example .env
```

```bash
# Frontend / SSR
NUXT_PUBLIC_PRODUCTION="false"
NUXT_SITE_URL="http://localhost:3000"
NUXT_PUBLIC_SITE_URL="http://localhost:3000"

# Persistência e cache
MONGO_URL="mongodb://mongo:27017/dns"
REDIS_URL="redis://redis:6379"

# Integração PowerDNS
DNS_HOST="http://powerdns:8081"
DNS_API_KEY="sua-chave-de-api-aqui"
DNS_SERVER_ID="localhost"

# Storage S3 (fotos de perfil)
FS_USERNAME="minioadmin"
FS_PASSWORD="minioadmin"
FS_BUCKET="sanchez-dns"
FS_ENDPOINT="http://s3:9000"
```

## Variáveis, uma a uma

### Frontend e sessão

| Variável | Papel |
|---|---|
| `NUXT_PUBLIC_PRODUCTION` | Quando `true`, ativa comportamentos de produção — em especial, marca o cookie de sessão como **`Secure`** (só trafega por HTTPS). Deixe `false` em desenvolvimento local sem HTTPS. |
| `NUXT_PUBLIC_SITE_URL` | URL pública do site, exposta ao cliente. É a base para o frontend montar as chamadas à API (`.../go`). |
| `NUXT_SITE_URL` | URL do site usada pelo servidor Nuxt (SSR), para SEO e metadados. |

**Por que `PRODUCTION` controla o cookie?** Porque em desenvolvimento você acessa por `http://localhost` (sem TLS); se o cookie fosse `Secure`, o navegador não o enviaria e o login não funcionaria. Em produção, atrás de HTTPS, `Secure` é obrigatório para não expor o token. Ver [Autenticação](/authentication#o-cookie-de-sessao).

### Persistência e cache

| Variável | Papel |
|---|---|
| `MONGO_URL` | Conexão com o MongoDB — a **fonte da verdade** de cadastros, sessões, permissões, solicitações e logs. |
| `REDIS_URL` | Conexão com o Redis — **cache** de sessão e de nível de acesso. Ver [Autenticação](/authentication#o-papel-do-redis). |

### PowerDNS

| Variável | Papel | Deve corresponder a |
|---|---|---|
| `DNS_HOST` | URL da **API HTTP** do PowerDNS (não a porta 53 do DNS). | `webserver-address:webserver-port` do `pdns.conf` |
| `DNS_API_KEY` | Chave secreta enviada no cabeçalho `X-API-Key`. | `api-key` do `pdns.conf` |
| `DNS_SERVER_ID` | Identificador do servidor usado no caminho da API (`/servers/<id>/...`). | `server-id` do `pdns.conf` (padrão `localhost`) |

Se qualquer um dos três não bater com a configuração do PowerDNS, as rotas de zona, registro e estatística falham. A relação completa está em [PowerDNS a fundo](/powerdns#o-elo-entre-a-config-do-powerdns-e-as-variaveis-do-sanchezdns).

### Upload de arquivos (S3)

| Variável | Papel |
|---|---|
| `FS_USERNAME` | Access key do storage S3 compatível. |
| `FS_PASSWORD` | Secret key do storage S3. |
| `FS_BUCKET` | Nome do bucket onde as fotos de perfil ficam. É criado automaticamente pelo serviço `s3-create-bucket` do compose. |
| `FS_ENDPOINT` | Endpoint do storage (ex.: `http://s3:9000`). Permite usar RustFS, MinIO ou AWS S3 sem mudar código. |

**Por que S3 em vez de disco?** Para manter o container sem estado — ver [Arquitetura](/architecture#s3-compativel-rustfs-minio-arquivos).

## DNSSEC nas zonas criadas

Quando uma zona é criada pela interface, o backend a cria como `Native`, adiciona uma `cryptokey` ativa do tipo `ksk` e grava o SOA inicial. Em outras palavras, **as zonas nascem com DNSSEC ativado** — desde que o backend do PowerDNS suporte (por exemplo, `gsqlite3-dnssec=yes` no SQLite). Detalhes em [PowerDNS a fundo](/powerdns#ativar-dnssec-post-zones-zona-cryptokeys).

## Sessão e permissões (resumo)

- O login cria o cookie **`sanchezdns_session_id`** (`HttpOnly`; `Secure` em produção), validado no backend a cada requisição protegida.
- O nível do usuário é **`admin`** ou **`member`**, lido a cada request (não embutido em token).
- O acesso aos registros de cada zona é controlado pelas listas **`leitura`** e **`escrita`** na coleção `users`. `admin` ignora essas restrições.

O funcionamento detalhado está em [Autenticação](/authentication) e [Usuários](/users).

## Observações importantes

- A aplicação **não usa JWT**; a autenticação é por sessão do lado do servidor (justificado em [Autenticação](/authentication)).
- Zonas e registros **não** ficam no MongoDB — a fonte da verdade do DNS é o PowerDNS.
- O sistema registra logs das operações administrativas e de DNS para auditoria (ver [Logs](/logs)).
