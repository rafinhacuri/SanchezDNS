# 💡 Por que SanchezDNS

SanchezDNS existe para resolver o trabalho repetitivo de operação DNS em um fluxo único, com regras claras de acesso e menos intervenção manual.

## O foco do projeto

- Centralizar a operação de zonas e registros em uma UI só.
- Automatizar o que puder ser derivado, como PTR reverso e sessão do usuário.
- Separar leitura, escrita e administração sem complicar o modelo.
- Manter auditoria das mudanças que realmente importam.

## O que ele faz bem

- gera e mantém zonas no PowerDNS de forma direta;
- trata zonas normais e reversas sem exigir procedimentos paralelos;
- reduz erros de manutenção ao completar FQDNs e sincronizar reverse records;
- expõe métricas úteis para operação diária;
- permite aprovações e permissões sem espalhar a lógica em várias ferramentas.

## O que o refactor removeu da visão anterior

- não há painel de múltiplas conexões;
- não há “cluster” de servidores DNS dentro da interface;
- a documentação antiga de storage e credenciais distribuídas não se aplica mais ao estado atual.

## Resumo

O projeto ficou mais simples de operar e mais explícito nas regras: uma instância PowerDNS, um conjunto de permissões por zona e automações nos pontos certos.
