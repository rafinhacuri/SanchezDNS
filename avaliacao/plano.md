# Plano da segunda rodada de avaliação (versão 3.0.3)

Este plano foi escrito e commitado **antes** da execução da segunda rodada. O hash e a data do commit que o introduz servem de evidência de que os critérios e as regras abaixo foram fixados antes dos resultados. A primeira rodada (versão 3.0.2, 30/09 e 15/10/2026) seguiu os critérios do Quadro 3 da monografia, registrados no Google Docs, e seus resultados ficam preservados em `resultado_e2_v3.0.2.md` e `e2_execucao_v3.0.2.log`.

## 1. Objeto

- Artefato: SanchezDNS 3.0.3 (tag `v3.0.3`), implantado pela imagem publicada na mesma instância do E4.
- Servidor autoritativo: o mesmo PowerDNS Authoritative Server 5.0.1 do E4.
- Zonas descartáveis: `aval-sanchezdns.example.` (RFC 2606) e `2.0.192.in-addr.arpa.` (RFC 5737), criadas e removidas pelo script.

## 2. Episódios desta rodada

| Id | Episódio | Script | Executor |
| --- | --- | --- | --- |
| R1 | E2: autorização e trilha de auditoria contra a API | `e2_autorizacao.sh` | autor, na máquina local |
| R2 | E2: condição de referência pela API do PowerDNS, sem a plataforma | `e2_referencia_api.sh` | autor, na máquina local |
| R3 | E2: condição de referência pelo `pdnsutil`, sem a plataforma | `e2_referencia_pdnsutil.sh` | autor, no servidor do PowerDNS |
| R4 | E2: zona publicada com tipos estruturados, verificada por oráculo externo | `e2_zona_ampliada.sh` + Zonemaster | autor |
| R5 | E3: recontagem das seis tarefas pela regra da seção 5 | manual | autor |

## 3. Casos e critérios de êxito (R1)

Os casos A0 a A12 repetem os da primeira rodada. Os critérios são os do Quadro 3, com duas mudanças decididas antes desta execução e justificadas pela primeira rodada:

- **A9**: na primeira rodada o caso era aprovado se o PTR permanecesse. O critério passa a exigir que o PTR seja **reapontado** para o nome que ainda usa o endereço, porque um PTR que aponta para um nome removido é a própria divergência A/PTR (Quadro 1).
- **Rastreabilidade**: o critério de 100% passa a valer também para as escritas derivadas de PTR e para o ciclo de vida das contas (RF13: "toda operação de escrita sobre zonas, registros e usuários").

Casos novos:

| Caso | Situação | Esperado | Classe (Quadro 1) | QP |
| --- | --- | --- | --- | --- |
| A13 | opA cria A no IP de um PTR de outra equipe e depois apaga esse A | PTR de outra equipe intacto após a remoção | Privilégio excessivo pela escrita derivada | QP3 |
| A14 | opA cria CNAME com destino sem ponto final e relê a zona | a leitura pela API e o dig mostram o destino com ponto final | Erro sintático no RDATA (3º princípio: exibir o publicado) | QP2 |
| L1 | admin promove e rebaixa o nível de opC | duas entradas `update_level` na trilha | Ausência de rastreabilidade | QP4 |

Critérios desta rodada:

| Construto | Critério de êxito |
| --- | --- |
| Contenção da escrita derivada (QP3) | Zero sobrescritas e zero remoções de PTR de outra equipe (A4, A8, A13); zero escritas em faixa sem zona reversa (A5); operador sem permissão não lista nem edita a zona reversa (A6, A7) |
| Riscos introduzidos pela automação (QP2) | A9 reapontado; A10 removido; A11 listado como órfão |
| Exibição do estado publicado (QP2) | A14: valor lido pela API igual ao publicado e igual à forma canônica |
| Rastreabilidade | 100% das escritas diretas, 100% das escritas derivadas de PTR e 100% das operações sobre contas com autor, ação e instante; o campo zone igual à zona afetada nas operações sobre registros |
| Tentativas recusadas | Nenhuma entrada na trilha (resultado descritivo, sem critério de êxito) |

Um caso só é aprovado quando a API e o `dig` concordam. Divergência entre os dois, falta de resposta ou erro 5xx tornam o caso inconclusivo.

## 4. Condição de referência (R2 e R3)

Objetivo: separar o efeito do artefato do efeito do servidor subjacente. As mesmas entradas dos casos N1 a N13 do Quadro 12 são enviadas **sem a plataforma**: pela API REST do PowerDNS (R2) e pelo `pdnsutil` (R3), em uma zona descartável `ref-sanchezdns.example.`.

Cada entrada é classificada como:

- **Recusada pelo servidor**: a API responde 4xx ou o `pdnsutil` sai com erro, e nada é gravado.
- **Aceita e reescrita**: gravada, mas com forma diferente da enviada (por exemplo, ponto final acrescentado).
- **Aceita como enviada**: gravada sem alteração.

Não há critério de êxito: o resultado é descritivo. Uma classe só é atribuída ao artefato quando a entrada é aceita (como enviada ou reescrita para um valor errado) em pelo menos uma das linhas de base e é impedida ou corrigida pela plataforma. Uma entrada recusada pelas duas linhas de base é atribuída ao servidor, e não ao artefato.

## 5. Regra de contagem do E3 (R5)

Fixada antes da recontagem, para as duas formas de operação:

- **Passo**: cada comando executado na linha de comando; na plataforma, cada envio de formulário ou confirmação de ação. Navegar entre telas não conta.
- **Entrada**: cada valor que o operador fornece, seja argumento de comando ou campo de formulário, inclusive os que têm valor padrão alterado. Campos com valor padrão mantidos não contam.
- **Valor formatado**: entrada cuja forma correta depende de regra sintática do DNS: nome de domínio qualificado, endereço IP, nome reverso, RDATA com ordem fixa, aspas ou prioridade, e endereço de e-mail no formato RNAME.
- **Classe ainda possível**: classe do Quadro 1 que a forma de operação não impede, considerando o comportamento observado em R2 e R3 para a linha de comando.

Passos, entradas e valores formatados são relatados lado a lado. A ameaça de que um formulário com vários campos conte como um passo é tratada pela coluna de entradas.

## 6. Oráculos (R4)

- O Zonemaster cobre o serial, o DNSSEC, a delegação e a sintaxe de nomes de NS, SOA e MX. Ele **não** tem casos de teste para SRV, CAA, HTTPS ou PTR.
- Para esses tipos, os oráculos são o `dig` contra o servidor autoritativo e, se disponível, um validador de zona independente do PowerDNS e do artefato (`named-checkzone`, do BIND, ou `kzonecheck`, do Knot) aplicado à zona exportada.
- A zona de R4 é a `teste.curi.dev.br.`, já delegada, acrescida de MX, SRV, CAA, HTTPS e CNAME criados pela plataforma.

## 7. Registro dos resultados

Cada script grava os resultados com a versão no nome do arquivo (`resultado_e2_v3.0.3.md`, `resultado_referencia_api.md` etc.), sem sobrescrever os da primeira rodada.
