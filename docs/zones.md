# 🧭 Zonas

## Visão geral

A tela de zonas mostra três grupos:

- zonas normais
- zonas reversas IPv4 (`in-addr.arpa`)
- zonas reversas IPv6 (`ip6.arpa`)

Cada entrada traz nome, serial, nível de acesso e status de DNSSEC.

## Criação de zonas

Ao criar uma zona, o formulário pede:

- domínio da zona
- `Start of Authority`
- email do responsável
- `refresh`
- `retry`
- `expire`
- `negativeCacheTtl`

O backend sempre cria a zona como `Native` e já ativa DNSSEC no PowerDNS.

## SOA

O SOA é atualizado diretamente na zona do PowerDNS. Os campos usados na interface são:

- `startOfAuthority`
- `email`
- `refresh`
- `retry`
- `expire`
- `negativeCacheTtl`

## Registros suportados

O editor de registros trabalha com:

- `A`
- `AAAA`
- `ALIAS`
- `CAA`
- `CNAME`
- `HTTPS`
- `MX`
- `NS`
- `PTR`
- `TXT`
- `SRV`
- `TLSA`

## Regras de preenchimento

- Se o nome do registro estiver vazio, o sistema usa o apex da zona.
- Se o nome não terminar com a zona, o backend completa automaticamente o FQDN.
- `HTTPS` e `SRV` usam campos próprios na interface.
- `MX` usa prioridade no valor normalizado.

## Reverse automático

O comportamento mais importante do refactor é este:

- ao criar um registro `A`, o sistema tenta criar o PTR reverso correspondente;
- ao criar um registro `AAAA`, o sistema tenta criar o PTR reverso correspondente;
- ao editar um `A` ou `AAAA`, o PTR é atualizado;
- ao remover um `A` ou `AAAA`, o PTR é removido junto.

O backend procura a melhor zona reversa disponível e mantém o registro alinhado com o valor atual.

## Permissões

O botão de usuário e a edição de SOA só aparecem quando a zona tem nível `ADMINISTRADOR`.

Para usuários com nível `LEITURA`, a tabela fica em modo somente leitura.

## O que mudou em relação à doc antiga

- Não existe fluxo de múltiplos servidores.
- Não há tela de conexões.
- O foco atual é uma única instância PowerDNS com zonas, permissões e auditoria centralizadas.
