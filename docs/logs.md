# 🧾 Logs

## O que aparece na página

A página de logs mostra entradas com:

- zona
- usuário
- ação
- detalhes
- data de criação

## Eventos registrados

Os logs são gerados para operações como:

- criação e remoção de zonas
- criação, edição e remoção de registros
- inclusão, edição e remoção de usuários por zona
- atualização de SOA
- aprovação ou rejeição de solicitações

## Comportamento

- a listagem tem busca textual;
- a página trabalha com paginação;
- os eventos são gravados no MongoDB;
- a ordenação padrão mostra os eventos mais recentes primeiro.

## Acesso

Somente `admin` acessa essa página.

## Objetivo

O log existe para auditoria operacional. Ele permite rastrear quem fez o quê e em qual zona, sem depender de histórico manual na interface.
