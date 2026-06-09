# 📊 Estatísticas

## Visão geral

A página de estatísticas mostra uma visão operacional da instância PowerDNS configurada no ambiente.

## Métricas exibidas

- quantidade de zonas
- quantidade total de registros
- uptime formatado
- status da instância
- `UDP Queries`
- `TCP Queries`
- `Server ID`
- momento de início calculado a partir do uptime

## Atualização

- os dados são buscados no backend;
- a interface atualiza automaticamente a cada 60 segundos.

## Fonte dos dados

O backend consulta a API de estatísticas do PowerDNS e também percorre as zonas para contar os registros.

## Limite atual

Essa tela reflete uma única instância PowerDNS por ambiente. Não há agregação entre vários servidores.
