# ⚙️ Configuração

Esta página resume o que o projeto espera do PowerDNS e quais variáveis controlam a aplicação.

## Modelo atual

SanchezDNS trabalha com **uma instância PowerDNS Authoritative por ambiente**.

- Não existe tela de múltiplas conexões.
- `DNS_HOST` define a URL da API do PowerDNS.
- `DNS_SERVER_ID` define qual servidor será consultado nas rotas de zona, registros e estatísticas.

## PowerDNS

Exemplo mínimo de configuração do PowerDNS:

```ini
api=yes
api-key=chave-da-api
webserver=yes
webserver-address=0.0.0.0
webserver-port=8081
server-id=localhost
```

## DNSSEC

Quando uma nova zona é criada pela interface, o backend:

- cria a zona como `Native`;
- adiciona uma `cryptokey` ativa do tipo `ksk`;
- grava o SOA inicial com os valores informados no formulário.

Em outras palavras, as zonas criadas por SanchezDNS já nascem com DNSSEC ativado no PowerDNS, desde que o backend do servidor suporte isso.

## Variáveis usadas pela aplicação

### Frontend e sessão

- `NUXT_PUBLIC_PRODUCTION`
- `NUXT_PUBLIC_SITE_URL`
- `NUXT_SITE_URL`

### Persistência e cache

- `MONGO_URL`
- `REDIS_URL`

### PowerDNS

- `DNS_HOST`
- `DNS_API_KEY`
- `DNS_SERVER_ID`

### Upload de arquivos

- `FS_USERNAME`
- `FS_PASSWORD`
- `FS_BUCKET`
- `FS_ENDPOINT`

## Sessão

O login cria o cookie `sanchezdns_session_id`.

- Em produção o cookie é marcado como seguro.
- A sessão é validada no backend a cada requisição protegida.
- O nível do usuário vem do cadastro e pode ser `admin` ou `member`.

## Permissões por zona

O acesso aos registros é controlado pela coleção `users` no MongoDB.

- `leitura` permite visualizar registros da zona.
- `escrita` permite criar, editar e remover registros.
- `admin` bypassa essas restrições e pode administrar zonas, logs, cadastros e solicitações.

## Observações importantes

- A aplicação não usa JWT para autenticação de usuário no fluxo atual.
- O armazenamento de fotos passa por um bucket S3 compatível.
- O sistema registra logs de operações administrativas e de DNS para auditoria.
