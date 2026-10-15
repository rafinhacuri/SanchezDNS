# Resultado do E3 — recontagem das seis tarefas (R5)

Recontagem feita em 15/10/2026 pela regra da seção 5 de `avaliacao/plano.md`, fixada antes da execução. Linha de comando: `pdnsutil` 5.0.1, com o comportamento observado em `resultado_referencia_pdnsutil.md` (R3). Plataforma: SanchezDNS 3.0.3, pelos formulários de `app/components/Zones.vue` e `app/components/Record.vue`.

Regras aplicadas:

- **Passo**: comando executado; na plataforma, envio de formulário ou confirmação de ação.
- **Entrada**: valor fornecido pelo operador (argumento ou campo). Campos com valor padrão mantido não contam. O nome da zona repetido em vários comandos conta uma vez por comando.
- **Valor formatado**: entrada cuja forma depende de regra sintática do DNS (nome de domínio, endereço IP, nome reverso, RDATA estruturado, RNAME). Entre parênteses, quantos são distintos.
- **Classe ainda possível**: classe do Quadro 1 que a forma de operação não impede. Na linha de comando, segundo o R3, o ponto final ausente é corrigido e o RDATA malformado é recusado; o TTL abaixo de 60 segundos é aceito.

## Quadro de contagem

| Tarefa | Linha de comando: passos | Entradas | Formatados | Classes ainda possíveis | Plataforma: passos | Entradas | Formatados | Classes ainda possíveis |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1. Criação de zona com DNSSEC | 5 (`create-zone Z ns`, `set-meta Z SOA-EDIT-API DEFAULT`, `secure-zone Z`, `set-nsec3 Z '1 0 0 -'`, `rectify-zone Z`) | 9 | 7 (3) | DNSSEC incompleto e serial sem delegação, se um comando for omitido | 2 (formulário da zona; inclusão do NS, que a criação não grava) | 5 (domínio, servidor primário, e-mail do responsável; tipo e valor do NS) | 4 (4) | Nenhuma eliminável; RNAME malformado se o e-mail for digitado com `@` (ver nota 2); DS no pai continua manual |
| 2. Inclusão de host com reverso | 2 (`add-record` do A e do PTR) | 8 | 6 (5) | Divergência A/PTR (PTR omitido ou com nome reverso errado); TTL < 60 | 1 | 3 (nome, tipo, endereço) | 1 (1) | Nenhuma eliminável, desde que a zona reversa exista |
| 3. Alteração do endereço de um host | 3 (`replace-rrset` do A, `delete-rrset` do PTR antigo, `add-record` do PTR novo) + busca por outros nomes no endereço antigo (≥ 2: `list-all-zones` e `list-zone` por zona) | 11 | 8 (6) | Divergência A/PTR; reverso órfão; remoção de PTR ainda referenciado; sobrescrita de PTR de outra equipe | 1 | 1 (endereço novo) | 1 (1) | Nenhuma eliminável |
| 4. Remoção de host e reverso | 2 (`delete-rrset` do A e do PTR) + busca (≥ 2) | 6 | 4 (4) | Reverso órfão; remoção de PTR ainda referenciado ou de outra equipe | 1 | 1 (nome digitado para confirmar a remoção) | 1 (1) | Nenhuma eliminável |
| 5. Inclusão de MX e TXT | 2 (`add-record` do MX e do TXT) | 8 | 6 (4) | TTL < 60 (o ponto final ausente é corrigido e o TXT sem aspas é recusado pelo `pdnsutil`) | 2 | 7 (MX: nome, tipo, destino, prioridade; TXT: nome, tipo, texto) | 1 (1) | Nenhuma eliminável |
| 6. Correção do conjunto de NS (dois servidores) | 1 (`replace-rrset Z Z NS ns1 ns2`) + ajuste no registrário | 5 | 4 (3) | Inconsistência entre pai e filho (fora do escopo) | 2 (uma ação por servidor) + ajuste no registrário | 2 | 2 (2) | Inconsistência entre pai e filho (fora do escopo) |

## Síntese

- **Passos**: menos na plataforma nas tarefas 1 a 4, iguais na tarefa 5 e mais na tarefa 6. O critério do Quadro 3 ("nenhuma tarefa com mais passos") não é atendido na tarefa 6.
- **Entradas e valores formatados**: menos na plataforma nas seis tarefas. A coluna de entradas responde à ameaça de que um formulário com vários campos conte como um passo só: na tarefa 5, os passos empatam (2 × 2), mas a plataforma pede 1 valor formatado contra 6.
- **Classes ainda possíveis**: contra o `pdnsutil` 5.0.1, o erro sintático no RDATA deixa de ser uma classe que a linha de comando permite publicar: o R3 mostrou o servidor corrigindo o ponto final e recusando o RDATA malformado. As classes que só a plataforma elimina são as que envolvem **mais de um registro ou mais de uma zona** (divergência A/PTR, reverso órfão, remoção ou sobrescrita de PTR alheio), as que dependem de **omitir um comando** (DNSSEC incompleto, serial sem delegação) e o **TTL abaixo de 60 segundos**.

## Notas

1. A comparação anterior (Quadro 16 da monografia, versão 3.0.1) contava "1 valor" na tarefa 1 e não incluía a inclusão do NS, que a criação de zona pela plataforma não grava. Esta recontagem corrige os dois pontos.
2. O campo "Email do responsável" é gravado no RNAME do SOA sem conversão do `@` (`api/zonas/update-soa.go`). A zona de teste foi criada com o RNAME já no formato correto (`rafael.curi.dev.br`), e por isso o Zonemaster não apontou o problema. A hipótese de que um e-mail digitado com `@` produz um RNAME malformado ainda precisa ser verificada.
3. Na linha de comando, o `pdnsutil` 5.0.1 trata o nome do registro como absoluto e recusa nome relativo (R3, caso N8). Por isso a duplicação do domínio, possível em arquivos de zona no estilo BIND, não aparece nesta linha de base.
