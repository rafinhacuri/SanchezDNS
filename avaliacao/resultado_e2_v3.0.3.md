# Resultado do E2 — autorização (QP3) e trilha de auditoria (QP4)

Execução: 2026-10-15 15:52 -0300 · versão 3.0.3 · rodada `20261015155205` · zonas descartáveis `aval-sanchezdns.example.` e `2.0.192.in-addr.arpa.` (RFC 2606 e RFC 5737) · verificação externa por `dig` contra o servidor autoritativo.

## Casos de teste

| Caso | Esperado | Obtido | Resultado |
| --- | --- | --- | --- |
| A0 — escrita sem sessão (PUT /records) | erro de sessão; nada gravado | 401 (sessão: Missing session ID); A=(vazio) | APROVADO |
| A1 — opC (sem vínculo) lê a zona direta | 403 | 403 | APROVADO |
| A2 — opB (leitura) cria A na zona direta | 403; nada gravado (dig) | 403; A=(vazio); PTR=(vazio) | APROVADO |
| A3 — opA cria A 192.0.2.10 sem PTR prévio | 200; PTR criado → a3-20261015155205.aval-sanchezdns.example. | 200; A=192.0.2.10; PTR=a3-20261015155205.aval-sanchezdns.example. | APROVADO |
| A4 — opA cria A no IP do PTR de outra equipe | 200; PTR intacto → servidor.outra-equipe.example. | 200; PTR=servidor.outra-equipe.example. | APROVADO |
| A5 — opA cria A em 198.51.100.7 (faixa sem zona reversa) | 200; nenhuma zona e nenhum PTR | 200; zonas novas: nenhuma; PTR=(vazio) (REFUSED) | APROVADO |
| A6a — opA lista a zona reversa (GET /records) | 403 | 403 | APROVADO |
| A6b — opA cria PTR na zona reversa (PUT /records) | 403; nada gravado | 403; PTR=(vazio) | APROVADO |
| A6c — opA apaga o PTR de outra equipe (DELETE /records) | 403; PTR intacto | 403; PTR=servidor.outra-equipe.example. | APROVADO |
| A7 — opA consulta órfãos da zona reversa | 403 | 403 | APROVADO |
| A7b — opA remove o PTR como órfão (DELETE /reverses) | 403; PTR intacto | 403; PTR=servidor.outra-equipe.example. | APROVADO |
| A8 — opA muda o IP de um host para 192.0.2.50 (PATCH /records) | 200; PTR de outra equipe mantido (QP3) | 200; PTR 192.0.2.50 mantido → servidor.outra-equipe.example.; PTR 192.0.2.20 removido | APROVADO |
| A9 — opA apaga A cujo IP (192.0.2.10) segue em uso por outro nome | 200; PTR reapontado → a9-20261015155205.aval-sanchezdns.example. | 200; A de a3=(vazio); PTR reapontado → a9-20261015155205.aval-sanchezdns.example. | APROVADO |
| A10 — opA apaga o último A de 192.0.2.10 | 200; PTR removido | 200; PTR=(vazio) | APROVADO |
| A11 — admin cria PTR sem A e consulta órfãos | 192.0.2.99 listado como órfão | 200; GET orphans 200; órfãos: 192.0.2.99 | APROVADO |
| A12 — opA grava PTR sob in-addr.arpa via FQDN na zona direta | recusado (≠200); nada gravado | 422; PTR=(vazio) | APROVADO |
| A13 — opA apaga o único A no IP de um PTR de outra equipe | 200; PTR intacto → outro-servidor.outra-equipe.example. | 200; PTR=outro-servidor.outra-equipe.example. | APROVADO |
| A14 — opA cria CNAME sem ponto final e relê a zona | 200; exibido = publicado = destino.example. | 200; exibido=destino.example.; publicado=destino.example. | APROVADO |

Total: 18 aprovados, 0 reprovados, 0 inconclusivos.

## Classes do Quadro 1 por caso

| Caso | Classe de erro | QP |
| --- | --- | --- |
| A0 | Atualização dinâmica não autenticada | QP3 |
| A1, A2, A6a–A6c, A7, A7b, A12 | Privilégio excessivo do operador | QP3 |
| A13 | Privilégio excessivo pela escrita derivada na remoção | QP3 |
| A14 | Erro sintático no RDATA (exibição do valor publicado) | QP2 |
| A3 | Divergência entre A/AAAA e PTR (dentro das condições da §3.4.5) | QP2, QP3 |
| A4, A5, A8 | Privilégio excessivo pela escrita derivada (§3.4.5, §6.2.1) | QP3 |
| A9 | Divergência entre A/AAAA e PTR (guarda de referência) | QP2 |
| A10, A11 | Registros pendentes e reversos órfãos | QP2 |
| Trilha, L1 | Ausência de rastreabilidade das alterações | QP4 |

## Trilha de auditoria

**26 de 26 alterações registradas** com autor, ação, zona e instante (zonas, registros e permissões: 18 de 18; contas: 8 de 8).
Alterações derivadas (PTR gravado ou removido pela automação): 5 de 5 registradas.
Tentativas recusadas presentes na trilha: 0 de 5.

| Alteração | Autor | Ação | Zona | Instante | Registrada |
| --- | --- | --- | --- | --- | --- |
| admin cria a zona aval-sanchezdns.example. | sim | sim (create_zone) | sim (campo zone) | sim | sim |
| admin cria a zona 2.0.192.in-addr.arpa. | sim | sim (create_zone) | sim (campo zone) | sim | sim |
| opA envia solicitação de cadastro | sim | sim (request_account) | n/a | sim | sim |
| admin aprova a solicitação de opA | sim | sim (review_account) | n/a | sim | sim |
| admin concede escrita em aval-sanchezdns.example. a opA | sim | sim (insert_user) | sim (campo zone) | sim | sim |
| opB envia solicitação de cadastro | sim | sim (request_account) | n/a | sim | sim |
| admin aprova a solicitação de opB | sim | sim (review_account) | n/a | sim | sim |
| admin concede leitura em aval-sanchezdns.example. a opB | sim | sim (insert_user) | sim (campo zone) | sim | sim |
| opC envia solicitação de cadastro | sim | sim (request_account) | n/a | sim | sim |
| admin aprova a solicitação de opC | sim | sim (review_account) | n/a | sim | sim |
| L1 admin altera o nível de opC para admin | sim | sim (update_level) | n/a | sim | sim |
| L1 admin altera o nível de opC para member | sim | sim (update_level) | n/a | sim | sim |
| admin semeia PTR de outra equipe em 192.0.2.50 | sim | sim (insert_record) | sim (campo zone) | sim | sim |
| A3 opA cria A a3-20261015155205.aval-sanchezdns.example. → 192.0.2.10 | sim | sim (insert_record) | sim (campo zone) | sim | sim |
| A4 opA cria A a4-20261015155205.aval-sanchezdns.example. → 192.0.2.50 | sim | sim (insert_record) | sim (campo zone) | sim | sim |
| A5 opA cria A a5-20261015155205.aval-sanchezdns.example. → 198.51.100.7 | sim | sim (insert_record) | sim (campo zone) | sim | sim |
| A8 (preparo) opA cria A a8-20261015155205.aval-sanchezdns.example. → 192.0.2.20 | sim | sim (insert_record) | sim (campo zone) | sim | sim |
| A8 opA altera a8-20261015155205.aval-sanchezdns.example. de 192.0.2.20 para 192.0.2.50 | sim | sim (edit_record) | sim (campo zone) | sim | sim |
| A9 (preparo) opA cria A a9-20261015155205.aval-sanchezdns.example. → 192.0.2.10 | sim | sim (insert_record) | sim (campo zone) | sim | sim |
| A9 opA apaga A a3-20261015155205.aval-sanchezdns.example. (192.0.2.10 segue em a9-20261015155205.aval-sanchezdns.example.) | sim | sim (delete_record) | sim (campo zone) | sim | sim |
| A10 opA apaga o último A de 192.0.2.10 (a9-20261015155205.aval-sanchezdns.example.) | sim | sim (delete_record) | sim (campo zone) | sim | sim |
| A11 admin cria PTR órfão em 192.0.2.99 | sim | sim (insert_record) | sim (campo zone) | sim | sim |
| A13 (preparo) admin semeia PTR de outra equipe em 192.0.2.70 | sim | sim (insert_record) | sim (campo zone) | sim | sim |
| A13 (preparo) opA cria A a13-20261015155205.aval-sanchezdns.example. → 192.0.2.70 | sim | sim (insert_record) | sim (campo zone) | sim | sim |
| A13 opA apaga A a13-20261015155205.aval-sanchezdns.example. (único nome de 192.0.2.70) | sim | sim (delete_record) | sim (campo zone) | sim | sim |
| A14 opA cria CNAME a14-20261015155205.aval-sanchezdns.example. → destino.example (sem ponto) | sim | sim (insert_record) | sim (campo zone) | sim | sim |

| Alteração derivada | Registrada |
| --- | --- |
| A3: PTR 192.0.2.10 criado pela automação | sim |
| A8 (preparo): PTR 192.0.2.20 criado pela automação | sim |
| A8: PTR 192.0.2.20 removido pela automação | sim |
| A9: PTR 192.0.2.10 reapontado pela automação | sim |
| A10: PTR 192.0.2.10 removido pela automação | sim |

| Tentativa recusada | Entradas na trilha |
| --- | --- |
| A2 opB PUT /records | 0 |
| A6b opA PUT /records na zona reversa | 0 |
| A6c opA DELETE /records na zona reversa | 0 |
| A7b opA DELETE /reverses | 0 |
| A12 opA PUT /records com FQDN fora da zona | 0 |
