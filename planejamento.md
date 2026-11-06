# Planejamento da revisão da monografia e do artefato

Checklist das mudanças pedidas pelo orientador (revisões de 15/10/2026 e de 19/10/2026, esta com a estrutura de 8 capítulos), conferidas contra o texto do Google Doc e o código em 15/10/2026, mais os achados da conferência.

Legenda: **[Orient.]** pedido do orientador · **[Achado]** encontrado na conferência do texto ou do código · ~~riscado~~ já resolvido no texto atual.

Decisão tomada: **avaliação sem participantes** por enquanto. Se aparecerem participantes, o estudo entra depois (ver "Adiado").

---

## Etapa 1 — Código 3.0.3: o que muda os resultados

- [x] **[Orient.]** A8: `records/update-reverso.go` grava o PTR novo com REPLACE sem a guarda de PTR existente que `InsertReverso` tem. Aplicar a mesma guarda na alteração de endereço.
- [x] **[Orient.]** A9: quando o nome removido é o alvo do PTR e o endereço continua em uso por outro nome, reapontar o PTR para o nome remanescente (o 2º princípio precisa incluir o reapontamento).
- [x] **[Achado]** Remoção encadeada apaga o PTR mesmo quando ele aponta para outro nome (PTR de outra equipe): só remover se o PTR apontar para o nome removido. Mesma classe do A8, no caminho da remoção.
- [x] **[Achado]** Guarda de referência trata falha de leitura de zona como "sem referência" e remove o PTR. Tornar conservadora: na falha, não remover.
- [x] **[Orient.]** Registrar na trilha cada escrita derivada de PTR (criação, reapontamento, remoção), com o operador que a disparou, a zona reversa, o endereço e o alvo.
- [x] **[Orient.]** Campo `zone` do `edit_record` guarda o nome do registro. **[Achado]** `insert_record` e `delete_record` também gravam o FQDN no lugar da zona. Corrigir os três.
- [x] **[Orient.]** RF13 exige "toda operação de escrita sobre ... usuários". **[Achado]** Não são auditados: solicitação de cadastro, aprovação/recusa, remoção de cadastro, edição de cadastro e troca de nível (promoção a admin, `update.Level`).
- [x] **[Orient.]** A12: o 422 do PowerDNS vira 500 genérico. Repassar como erro do cliente (422) com a mensagem do servidor.
- [x] Testes automatizados das guardas (A8, A9, remoção de PTR alheio) com PowerDNS simulado; `go test`, lint e build.
- [x] Release 3.0.3.

## Etapa 2 — Código: segurança e on-premises

- [x] **[Orient.]** Cookie de sessão sem `SameSite` (nem proteção CSRF). Definir `SameSite=Lax` ou `Strict`. Rotas POST (`/user`, `/logout`) aceitam corpo JSON sem checar o Content-Type.
- [x] **[Orient.]** Enumeração: o login responde "cadastro em análise" antes de conferir a senha. A solicitação guarda o hash da senha, então dá para conferir antes de revelar.
- [x] **[Orient.]** Geolocalização por serviço externo: tornar opcional por variável de ambiente (desligada por padrão). **[Achado]** Hoje usa `http://ip-api.com` sem TLS e o User-Agent `cbpf-sso/1.0`.
- [x] **[Orient.]** RF03 depende do DoH da Google e da Cloudflare: tornar os resolvedores configuráveis (rede isolada).
- [x] **[Achado]** Listagem de zonas exclui duas zonas por nome fixo (resíduo de teste, §4.3.3). Remover.
- [x] **[Orient.]** Sessão inválida responde 400 em vez de 401 (RFC 9110). Corrigido para 401 (o frontend não dependia do 400; script do E2 aceita os dois).
- [ ] **[Achado]** `PUT /file` e `GET /file/:id` fora do grupo autenticado, o que contradiz o RNF04 ("não existe caminho de escrita não autenticado"). O `PUT /file` é usado no cadastro, antes de existir sessão, então não dá para autenticar. Ajustar o RNF04 no texto (Etapa 9) para "nenhum caminho de escrita sobre dados de zona sem autenticação".

## Etapa 3 — Reexecutar a avaliação

- [x] **[Orient.]** Commitar `avaliacao/` (hoje fora do git) com `avaliacao/plano.md` datado, contendo critérios e regras **antes** da nova rodada. Isso serve de evidência de que os critérios foram fixados antes da execução. Tirar o default `mizle.cbpf.dev.br` do script (endereço por variável de ambiente).
- [x] **[Orient.]** Reexecutar o E2 de autorização e trilha na 3.0.3 (A8, A9, escritas derivadas, contas).
- [x] **[Achado]** Novo caso: criar A sobre o PTR de outra equipe e depois apagá-lo (o PTR alheio deve ficar).
- [x] **[Orient.]** Condição de referência do E2 (prometida na §2.6.2, nunca executada): mandar as entradas N1–N13 direto à API do PowerDNS e ao `pdnsutil`, sem a plataforma. Hipótese a testar: a API recusa conteúdo não canônico e nome sem ponto; o pdnsutil trata conteúdo como absoluto, mas nome do registro como relativo (duplicação possível); TTL 30 passa nos dois.
- [x] **[Orient.]** Oráculo independente na normalização: incluir MX, SRV, CAA, HTTPS e PTR na zona de teste e rodar o Zonemaster de novo (hoje só A, NS e TXT).
- [x] **[Orient.]** Teste do 3º princípio (exibir o publicado): gravar pela API um valor sem ponto e reler por `GET /records`.
- [x] **[Orient.]** E3: definir a regra de contagem de passos antes de recontar. Reposicionar a comparação contra a API e o `pdnsutil`.
- [x] **[Achado]** Quadro 16, tarefa 1: diz "1 valor (nome da zona)", mas o formulário exige domínio, servidor primário e e-mail do SOA (`api/zonas/types.go`).
- [x] Rodar R1 a R5 do `avaliacao/plano.md`. Feitos: R1 (18/18; trilha 26/26; derivadas 5/5), R2 (API recusou 11 de 13), R3 (pdnsutil corrige ponto final e recusa RDATA malformado; TTL 30 passa nos dois), R5 (`resultado_e3_v3.0.3.md`) e Zonemaster; R4 com os seis tipos publicados na forma canônica (terceira execução) (scripts prontos; precisam das credenciais, por isso rodam na sua máquina e no servidor).
- [x] **[Achado]** `records/nome-completo.go` compara o nome com a zona como texto, sem fronteira de rótulo (mesma classe do P4): com a zona `curi.dev.br`, o nome `xcuri.dev.br` vira FQDN fora da zona. Corrigido para a 3.0.5, com teste.
- [ ] **[Achado]** O E2 (3.0.2 e 3.0.3) rodou no ambiente de desenvolvimento desta máquina (`bun dev`, `mizle.cbpf.dev.br:3000`) ligado ao PowerDNS do E4, e não na imagem publicada. A §5.3.3 diz "sobre a instância do episódio E4": corrigir no texto (Etapa 10).
- [x] **[Achado]** API inconsistente no MX: `PUT /records` recebe a prioridade em campo próprio; `DELETE /records` exige a prioridade dentro do valor (`"10 mail..."`). Corrigido para a 3.0.5: inclusão, edição e remoção aceitam os dois formatos.
- [x] **[Achado]** `DELETE /records` responde 200 quando nenhum registro casa com o valor enviado (nada é removido). Corrigido para a 3.0.5: remoção e edição sem registro correspondente respondem 404.
- [ ] **[Achado]** No `pdnsutil` 5.0.1 o nome do registro é absoluto (nome relativo é recusado como fora da zona) e o conteúdo passa pelo mesmo parser da API. A "edição manual" da linha de base não reproduz o erro de ponto final do .se: o argumento da §1.1 e o E3 precisam refletir isso.
- [x] **[Achado]** O campo "Email do responsável" vai para o RNAME do SOA sem converter o `@` (`api/zonas/update-soa.go`). Corrigido para a 3.0.5: o e-mail é convertido para RNAME (pontos da parte local escapados).
- [x] Release 3.0.5 com as quatro correções acima (commit 33ccb8a).

## Etapa 4 — Nova estrutura da monografia (revisão do orientador de 19/10/2026)

A estrutura proposta em 15/10 (7 capítulos) foi substituída pela de 19/10, com 8 capítulos e trabalhos correlatos separados da fundamentação. As Etapas 5 a 12 seguem esta estrutura.

- [x] **[Orient.]** Reorganizar em 8 capítulos, no esquema de Gregor e Hevner (2013):
  1. Introdução (enxuta)
  2. Fundamentação teórica
  3. Trabalhos correlatos
  4. Método
  5. Projeto e desenvolvimento do artefato (fusão dos atuais Caps. 3 e 4)
  6. Demonstração e avaliação
  7. Discussão (novo)
  8. Conclusões
- [x] **[Orient.]** Aplicar o mapa de migração (a reorganização é, em grande parte, mover texto):

  | Conteúdo atual | Destino |
  | --- | --- |
  | 1.1 (parágrafos empíricos) | 2.3 (já duplicados lá) |
  | 2.2.1 | 2.2 |
  | 2.2.2 a 2.2.4 | 2.3 |
  | 2.3 | 2.4 |
  | 2.4.1 a 2.4.3 | 4.2 e 4.3 |
  | 2.5 | Capítulo 3 |
  | 2.5.3 (segundo motivo, Technitium) | 3.2 resumido; análise em 7.3 |
  | 2.6 | 4.3 a 4.6 (E4 em 6.4) |
  | 3.1 e 3.2 | 4.1 e 4.2 |
  | 3.3 | 5.1 |
  | 3.4 a 3.7 com 4.1 a 4.7 | Capítulo 5, fundidos |
  | 5.1 a 5.6 | Capítulo 6 |
  | 6.1.1 e 6.1.2 | 5.8 (enunciado), 7.2 (reavaliação), 8.1 |
  | 6.1.3 | 7.1 |
  | 6.2 | 7.6 |
  | 6.3 | 8.2 |

- [ ] **[Orient.]** Se o prazo apertar, a prioridade é: (1) mover a §2.6 para o capítulo de Método; (2) criar o capítulo de Discussão a partir das §6.1 e §6.2; (3) fundir os Caps. 3 e 4. O resto é refinamento. Andamento: (1) feito na reestruturação; (2) capítulo criado com o texto movido, falta escrever (Etapa 11); (3) textos lado a lado no Cap. 5, falta fundir (Etapa 9).
- [ ] **[Orient.]** Enxugar o texto onde puder: cortar repetições entre capítulos e parágrafos que só antecipam o que vem depois. Andamento: Cap. 1 (cerca de 40% menor) e Cap. 2 feitos; faltam os Caps. 3 a 8.
- [x] ~~**[Orient.]** Quadros 16 e 17 antes do 14 e do 15 e fora da Lista de Quadros~~ (renumeração já feita; faltam as páginas, ver Fechamento).
Regra para todas as etapas (não é tarefa): todo texto novo ou reescrito passa pelo humanizer antes de ir para o Doc, e o Rafael é avisado antes de qualquer mudança no Doc.
- [x] Reestruturação aplicada no Doc em 26/10/2026, com texto movido sem reescrita. Cópia de segurança: "BACKUP 2026-10-26 antes da reestruturação". Conferido: 17/17 quadros, 10/10 figuras, todos os parágrafos de texto preservados (só uma linha em branco removida), 127 títulos na ordem do mapa, 11 comentários ancorados.
- [x] **Rafael:** mover à mão os 4 blocos marcados em amarelo e apagar as marcas depois: Quadro 4 (fim do Cap. 4 → §4.2), Figura 3 (depois da §5.7.2 → §5.2.1), Quadro 10 (→ §5.2.2) e Quadro 11 (→ §5.3.1). Cada bloco é legenda + quadro/figura + linha Fonte. Feito e conferido em 26/10/2026, sem marcas restantes.
- [ ] Notas "[A ESCREVER/REESCREVER NA ETAPA N]" marcam as seções novas vazias (1.6, 3.1, 3.2.3, 5.8, 7.3, 7.4, 7.5, 8.1); somem quando a etapa for feita. Andamento: 1.6 escrita (Etapa 5); faltam as outras sete.
- [ ] O texto movido ainda cita a numeração antiga ("seção 3.6.3", "Capítulo 3", "seção 2.6.1" etc.) e a §5.4 tem as duas descrições de sessão (antigos 3.4 e 4.3) lado a lado: corrigir referências e fundir nas Etapas 5 a 12. Andamento: Caps. 1 e 2 corrigidos; faltam os Caps. 3 a 8.
- [ ] Numeração dos quadros mudou de ordem (o Quadro 4 passa a vir antes do Quadro 3): renumerar no Fechamento, junto com Sumário e listas, que são texto digitado e estão desatualizados.

## Etapa 5 — Cap. 1, Introdução (enxuta) — aplicada no Doc em 26/10/2026

Estrutura: 1.1 Contexto e problema · 1.2 Lacuna das soluções existentes · 1.3 Problema de pesquisa e questões · 1.4 Objetivos · 1.5 Escopo e delimitação · 1.6 Contribuições esperadas (nova) · 1.7 Estrutura do documento.

- [x] **[Orient.]** §1.1 reduzida em cerca de dois terços, com três ou quatro evidências-âncora (Oppenheimer, Bojović e Gajin, Sommese e um incidente); os parágrafos empíricos restantes migram para a §2.3, onde já estão duplicados.
- [x] **[Orient.]** §1.2 em versão curta, sem antecipar o Quadro 2.
- [x] **[Orient.]** Problema independente do artefato. Hoje diz "a falta de um meio de operação que torne esse erro impossível", o que embute a solução e torna a avaliação tautológica.
- [x] **[Orient.]** Contexto e interessados explícitos: organizações com PowerDNS próprio e equipes distintas para zonas diretas e reversas.
- [x] **[Orient.]** Enquadramento técnico: servidores validam a sintaxe de cada registro isoladamente, mas não a consistência entre registros interdependentes nem o alcance da autorização. O erro humano fica como motivação.
- [x] **[Orient.]** Falar em grau ("em que medida reduz"), não em "impossível".
- [x] **[Orient.]** QP2: retirar "necessários" (a necessidade não é testada; exigiria ablação).
- [ ] **[Orient.]** Rastreabilidade está sob a QP4, que não fala de auditoria. Realinhar a QP, o objetivo específico e o Quadro 3. QP3 e objetivo específico já realinhados no Cap. 1 (26/10); falta o Quadro 3, na §4.5 (Etapa 8).
- [x] **[Orient.]** Objetivos de pesquisa (caracterizar, projetar, avaliar), não funcionalidades. Tirar "implementar sob arquitetura desacoplada" dos objetivos.
- [x] **[Orient.]** Nova §1.6, Contribuições esperadas: melhoria (Gregor e Hevner, 2013), com instanciação e princípios de projeto; também a taxonomia e o protocolo de avaliação. Recebe as alegações de contribuição que hoje estão na §2.5.9.
- [x] **[Orient.]** §1.1, "Era um erro sintático, do tipo que uma validação [...] teria barrado": refinar conforme a condição de referência (o PowerDNS já recusa ou corrige esses erros).
- [x] **[Achado]** §1.3 "outras dependem de evidência empírica com operadores humanos" e "revisão sistemática da literatura" (não houve protocolo sistemático): restos do desenho antigo.
- [x] **[Achado]** QP principal "registre integralmente" e objetivo geral "reduza a ocorrência": verbos mais fortes do que a avaliação sustenta.
- [x] **[Orient.]** §1.2: o argumento on-premises precisa ficar coerente com a geolocalização externa e o DoH (agora opcionais, Etapa 2).

## Etapa 6 — Cap. 2, Fundamentação teórica — aplicada no Doc em 26/10/2026

Estrutura: 2.1 Arquitetura do DNS (atual 2.1) · 2.2 Erro humano e projeto de interfaces (atual 2.2.1) · 2.3 Erros de configuração: evidência empírica, impactos e taxonomia (atuais 2.2.2 a 2.2.4) · 2.4 Controle de acesso e auditoria (atual 2.3).

- [x] **[Orient.]** §2.2 ampliada, com Reason (1990) e Norman (2013) apresentados como kernel theories.
- [x] **[Orient.]** §2.3 recebe os parágrafos empíricos tirados da §1.1, sem duplicar o que já está lá; fecho no Quadro 1.
- [x] **[Orient.]** Classe "DNSSEC incompleto" inclui a rotação de chaves (Chung et al.), mas nenhum mecanismo trata rotação: cortar a classe ou declarar a rotação fora do escopo.
- [x] **[Achado]** §2.2.4: "para que a redução de erros possa ser medida por categoria e no total" não é mais medido.
- [x] Abertura e §2.5 reescritas (ainda citavam DSR, correlatos e avaliação); remissões a "Capítulo 3" e "seção 2.6" atualizadas para Cap. 5 e §4.5; Quadro 1 com rastreabilidade na QP3.
- [x] Cap. 1: a QP1 cita os estudos de medição (comentário "literatura: quais?") e a abertura não descreve mais o capítulo.
- [ ] Para a Etapa 11 (Discussão): a §2.3.3 diz que a camada elimina o erro sintático no RDATA, mas a condição de referência mostrou que o PowerDNS já recusa ou corrige a maior parte desses erros.
- [ ] Para o Fechamento: siglas com travessão no Cap. 2 ("(*Fully Qualified Domain Name* – FQDN)") e com "do inglês" no Cap. 1; uniformizar.

## Etapa 7 — Cap. 3, Trabalhos correlatos — aplicada no Doc em 26/10/2026

Estrutura: 3.1 Critérios de comparação (nova) · 3.2 Ferramentas por abordagem (atuais 2.5.1 a 2.5.8) · 3.3 Quadro comparativo e lacuna (atual 2.5.9).

- [x] **[Orient.]** Nova §3.1 com os critérios de comparação (as colunas do Quadro 2 e a escala Sim/Parcial/Não).
- [x] **[Orient.]** Poweradmin aparece no Quadro 2 sem subseção própria: criar.
- [x] **[Orient.]** §2.5.3 (Technitium) antecipa a avaliação: manter o incidente resumido na §3.2 e levar a análise para a §7.3; tirar a remissão aos testes do Cap. 5.
- [x] **[Orient.]** §2.5.9 sem as alegações de contribuição ("o trabalho se situa no quadrante de melhoria"), que vão para a §1.6 e para a Discussão.
- [x] **[Achado]** §2.5.6: "o SanchezDNS precisa emular a atomicidade do Knot" contradiz o RNF03 (melhor esforço).
- [x] **[Achado]** "As documentações do Quadro 2 foram consultadas em setembro de 2026": manter a data, que sustenta o ponto de entrada do DSRM (Etapa 8).
- [x] Poweradmin conferido na documentação em 26/10/2026: a v4.5.0 (20/10/2026) mudou a linha do Quadro 2 (remoção e DNSSEC de "Não" para "Parcial"; auditoria de "Sim" para "Parcial", porque o registro em banco vem desligado). A §1.2 passou a citar a verificação de referência antes de remover o reverso.
- [x] Remissões "seção 2.6.2" (para §4.4) e "seção 2.5.5" (para §3.2.6) corrigidas; análise do Technitium e frases de contribuição guardadas como texto de origem nas §7.3 e §7.4.
- [x] **[Orient.]** Apêndice A com a fonte (link da documentação) de cada célula do Quadro 2, como o orientador sugeriu ("se fez a análise, tem que constar; pode até ir para o apêndice"). Aplicado em 26/10/2026, com 10 células do Quadro 2 corrigidas, a frase da §3.3 sobre remoção do reverso ajustada e 3 referências novas (ISC, AWS, Google Cloud). Pendente: auditoria do Technitium ficou "Não" por falta de fonte; Cloudflare não conferido célula a célula.

## Etapa 8 — Cap. 4, Método — aplicada no Doc em 04/11/2026

Estrutura: 4.1 Classificação da pesquisa (atual 3.1) · 4.2 DSR e DSRM, com ponto de entrada (atuais 2.4.1, 2.4.2 e 3.2) · 4.3 Estratégia de avaliação FEDS (atuais 2.4.3 e 2.6.1) · 4.4 Episódios de avaliação (atual 2.6.2) · 4.5 Matriz de avaliação e critérios de êxito (Quadro 3) · 4.6 Ameaças à validade (atual 2.6.4).

- [x] **[Orient.]** Ponto de entrada do DSRM: a §3.2 diz "centrada no problema", mas o código começou em out/2025 e a literatura e as fontes do Quadro 2 são de set/2026. Declarar entrada centrada em projeto e desenvolvimento (Peffers et al., 2007).
- [x] **[Orient.]** Estratégia FEDS (Technical Risk and Efficacy) justificada pelo risco dominante; declarar que se avalia a eficácia do meio, não a efetividade junto a operadores.
- [x] **[Orient.]** Proposição de projeto em forma CIMO (contexto, intervenção, mecanismo, resultado), indicando qual parte cada episódio testa.
- [x] **[Orient.]** A §2.6.3 (implantação real) é detalhe do E4: vai para a §6.4 ou fica como subseção da §4.4.
- [x] **[Orient.]** A §3.2 remete a "5.3 e 5.5" como demonstração sem seção própria: apontar para a demonstração do Cap. 6.
- [x] **[Orient.]** Evidência de que os critérios foram fixados antes da execução: commit datado do plano (`avaliacao/plano.md`, 48d40b4). **[Achado]** Remover "definidos previamente à construção" (§1.6, §2.7), que é falso.
- [x] **[Orient.]** Ameaça à validade de construto: regra de contagem de passos do E3 (formulário de seis campos = 1 passo; cada comando = 1 passo) (Wohlin et al., 2012), tratada pela coluna de entradas do R5.
- [x] **[Orient.]** Ameaça: oráculo dos casos N é o valor esperado escrito pelo autor; o Zonemaster só é independente para serial e DNSSEC (e sintaxe do MX, depois da zona ampliada).
- [x] **[Orient.]** Terminologia: vários casos "negativos" são testes de canonicalização positiva. Separar canonicalização (N1–N10) de rejeição (N11–N13).
- [x] **[Orient.]** Condição de referência: descrever como executada (R2 e R3).
- [x] **[Orient.]** Pseudonimização: a §2.6.4 promete, mas o texto e as Figuras 4–7 expõem curi.dev.br, 152.84.126.7 e o AS 2715. Decidir entre pseudonimizar ou retirar a promessa. **[Achado]** O servidor fica em infraestrutura do CBPF.
- [x] **[Achado]** Remover restos do desenho antigo: título "Framework FEDS e Escala SUS" (§2.4.3) e Quadro 4 "Demonstração ... com conferência por oráculo externo".
- [x] Decisões desta etapa: a §2.6.3 ficou como §4.4.1; a pseudonimização foi retirada (domínio do autor, servidor do CBPF usado com autorização); a lista das seis atividades do DSRM saiu por repetir o Quadro 4; a §4.1 não repete mais March e Smith, Hevner e Gregor e Hevner. Critérios do Quadro 3 mantidos; só a rastreabilidade passou a QP3 e as remissões foram atualizadas. Referência nova: Denyer, Tranfield e van Aken (2008).
- [ ] Para a Etapa 10: o título do Quadro 12 ainda diz "Casos de teste negativos"; trocar para canonicalização e rejeição, como no E2 da §4.4.
- [x] Para a Etapa 12: o Resumo ainda fala em "casos de teste negativos".

## Etapa 9 — Cap. 5, Projeto e desenvolvimento do artefato (fusão dos Caps. 3 e 4) — aplicada no Doc em 04/11/2026

Estrutura: 5.1 Requisitos e rastreabilidade (3.3, Quadro 5) · 5.2 Arquitetura e pilha tecnológica (4.1) · 5.3 Integração com o PowerDNS (4.2) · 5.4 Autenticação e controle de acesso (3.4 + 4.3) · 5.5 Normalização e interface preventiva (3.5 + 4.4) · 5.6 Sincronização de reversos (3.6 + 4.5) · 5.7 Auditoria (3.7 + 4.6) · 5.8 Princípios de projeto instanciados (de 6.1.2) · 5.9 Empacotamento e implantação (4.7).

- [x] **[Orient.]** Cada subseção segue a ordem: decisão de projeto, justificativa, implementação, limitação conhecida. Cada mecanismo descrito uma vez só (hoje duplicados: 3.4.1 × 4.3.1, 3.5 × 4.4, 3.6 × 4.5, 3.7 × 4.6).
- [x] **[Orient.]** Cuidar para não perder a rastreabilidade com os requisitos na fusão.
- [x] **[Orient.]** Retornos do ciclo construir-avaliar (guarda de referência; exibição do valor publicado) registrados dentro das §5.5 e §5.6.
- [x] **[Orient.]** §5.8: os quatro princípios de projeto, cada um ligado ao mecanismo que o realiza e ao caso de teste que pode refutá-lo (Gregor e Jones, 2007).
- [x] **[Orient.]** Quadro 6 afirma impedir a enumeração de usuários: ajustar à Etapa 2.
- [x] **[Orient.]** Mencionar SameSite/CSRF no modelo de sessão.
- [x] **[Achado]** Quadro 7 omite `DELETE /reverses`, que o caso A7b usa.
- [x] **[Achado]** RNF04: restringir a "nenhum caminho de escrita sobre dados de zona sem autenticação" (`PUT /file` e `POST /cadastro` são escritas sem sessão, fora das zonas).
- [x] **[Achado]** RF10 diz "criação atômica em lote", mas a §3.6.4 diz não transacional.
- [x] **[Achado]** §3.7 diz "cada operação de escrita grava um documento" e que o ciclo de contas fica fora da trilha: alinhar com o RF13 e com a 3.0.3.
- [x] Descrever as mudanças da 3.0.3 e da 3.0.5 (guardas, reapontamento, auditoria, nome por rótulo, MX, 404, RNAME).
- [x] Fusão: saíram as subseções 5.4.6–5.4.8, 5.5.4–5.5.6, 5.6.5–5.6.7 e 5.7.1 (o que só elas tinham foi para as que ficaram); a 5.7.2 virou 5.7.1. Quadros e figuras no lugar (17/17, 10/10). Código conferido: 401, SameSite=Lax, "em análise" só com a senha, GEOIP_URL e DOH_RESOLVERS opcionais, DELETE /reverses, 21 códigos de auditoria, reapontamento do PTR, 404, 422, RNAME.
- [ ] Para a Etapa 13 (enxugar): o Quadro 6 ainda registra que o pedido de cadastro responde 409 a e-mail já registrado; se quiser fechar a enumeração, é mudança de código, não de texto.
- [ ] Para as Etapas 10 e 11: os Caps. 6 e 7 ainda citam subseções que saíram (5.4.6, 5.5.6, 5.6.6 etc.) e a numeração antiga dos Caps. 3 e 4; corrigir nas etapas deles.

## Etapa 10 — Cap. 6, Demonstração e avaliação — aplicada no Doc em 04/11/2026

Estrutura: 6.1 E0 e E1 (atuais 5.1 e 5.2) · 6.2 E2: conformidade técnica (5.3: normalização, Zonemaster, autorização, auditoria) · 6.3 E3: comparação com a operação manual (5.4) · 6.4 E4: implantação real (5.5 e 2.6.3) · 6.5 Síntese frente aos critérios (5.6, com o quadro de resultados atualizado).

- [x] **[Orient.]** Separar demonstração (zona publicada, E4) de avaliação (casos N e A, E3, condição de referência).
- [x] **[Orient.]** Resultados relatados sem interpretação, na ordem do método; as reprovações aparecem com o seu critério (N13, P4, P6, A8, A9) e as lacunas de execução também.
- [x] Atualizar com os resultados da Etapa 3 (3.0.3: R1 a R5) e indicar o que a 3.0.5 mudou.
- [x] **[Achado]** O E2 rodou no ambiente de desenvolvimento (`bun dev`) ligado ao PowerDNS do E4, não na imagem publicada: corrigir "sobre a instância do episódio E4".
- [x] **[Orient.]** A12: registrar o 500 devolvido a um erro do cliente na 3.0.2 e a correção para 422 na 3.0.3.
- [x] **[Orient.]** RF13 parcialmente descumprido na 3.0.2 (0/5 escritas derivadas; contas fora da trilha), não "fora do escopo"; na 3.0.3, 26/26 e 5/5.
- [x] **[Achado]** Quadro de síntese, linha QP1: "casos de teste executados para parte delas" contradiz a §5.1 ("todas").
- [x] ~~**[Orient.]** Quadro de resultados com QP3 "não executado" e rastreabilidade "não verificado"~~ (já atualizado no Quadro 17).
- [x] ~~**[Orient.]** §5.1 dizendo "apenas por inspeção"~~ (já atualizado).
- [x] **[Orient.]** §5.6: incluir o A9 no parágrafo (o Quadro 17 já cita; o texto não).
- [x] **[Achado]** E3 recontado (R5): tarefa 1 com 2 passos na plataforma (inclusão do NS) e linha de comando pelo `pdnsutil` 5.0.1, que corrige o ponto final e usa nome absoluto.
- [x] Feito também: §6.2.4 nova (condição de referência R2 e R3); zona ampliada (R4) na §6.2.2; Quadro 14 com A13 e A14 e resultados das duas versões; Quadro 15 e Quadro 16 com a 3.0.3; casos N reexecutados sobre a 3.0.5 em 04/11/2026 (18 aprovados, P6 ignorado como limitação); E2 em `bun dev` registrado como desvio.
- [ ] Para a Etapa 11: os Caps. 7 e 8 ainda remetem à "seção 5.3.3" (testes contra a API, agora §6.2.3) e à "seção 3.4.5" (agora §5.4.5), e o texto da Discussão ainda descreve o A8 e a trilha como na 3.0.2.
- [ ] Para o Fechamento: a Lista de Quadros ainda traz o título antigo do Quadro 12.

## Etapa 11 — Cap. 7, Discussão (novo) — aplicada no Doc em 04/11/2026

Estrutura: 7.1 Respostas às questões de pesquisa (6.1.3) · 7.2 Princípios de projeto reavaliados (6.1.2) · 7.3 Automação e autorização granular: A8, A9 e A12 à luz do Technitium (2.5.3 e 5.3.3) · 7.4 Reposicionamento no quadro comparativo · 7.5 Eficácia demonstrada versus efetividade não medida · 7.6 Limitações (6.2).

- [x] **[Orient.]** §7.2, destino dos princípios: 1º se sustenta (N13 o refina: regra no servidor); 2º parcial (P4 + A9: precisa do reapontamento); 3º se sustenta (A14); 4º refutado como implementado (A8) e reformulado como "a guarda vale para toda operação que grava o registro derivado".
- [x] **[Orient.]** §7.3: A8 e A9 à luz do incidente do Technitium 15.5; A8 sobrescrevendo o PTR de outra equipe sem rastro na trilha é a combinação mais grave e merece destaque; A9 é uma classe de erro nova criada pela automação (QP2, "novos riscos").
- [x] **[Achado]** §7.3: só o A8 reproduz a falha do Technitium; no A12 a proteção veio do PowerDNS, não da camada de autorização.
- [x] **[Achado]** "A violação tende a ocorrer nos caminhos menos óbvios": apresentar como hipótese (dois casos), não como proposição demonstrada.
- [x] **[Orient.]** §7.4: reposicionar o SanchezDNS no Quadro 2 pela evidência ("Parcial" na 3.0.2 em remoção encadeada e triagem, permissão por zona e trilha; dizer o que a 3.0.3 mudou).
- [x] **[Orient.]** §7.5: separar a eficácia demonstrada da efetividade, que continua hipótese (Venable, Pries-Heje e Baskerville, 2016).
- [x] **[Achado]** Condição de referência: no PowerDNS, o erro sintático de RDATA já é recusado (API) ou corrigido (`pdnsutil`); o que só a plataforma elimina é o que envolve mais de um registro ou zona, o que depende de lembrar um comando e o TTL abaixo de 60 s.
- [x] **[Orient.]** §7.6: condição de corrida (TOCTOU) entre a guarda de referência e a criação do PTR.
- [x] **[Orient.]** §7.6: geolocalização externa (dados pessoais, LGPD) e dependência de DoH em rede isolada, agora opcionais.
- [x] **[Orient.]** Ameaças à validade depois da execução.
- [x] **[Achado]** Protocolo "replicável em outras ferramentas": o script é específico da API do SanchezDNS; replicável é o protocolo, não o script.
- [x] **[Achado]** §6.2.5: "Parte do plano também ficou por executar" seguido de "foram executados depois" e "não verificados" (o Quadro 17 não tem mais esse status).
- [x] Feito também: §7.1 com as duas rodadas (rastreabilidade na QP3); §7.3, §7.4 e §7.5 escritas no lugar das notas e dos textos de origem; §7.6 com a condição de corrida, LGPD e DoH; referência nova BRASIL (2018), LGPD.
- [x] Para a Etapa 12: o Cap. 8 ainda cita "seção 5.3.3" para os testes contra a API (agora §6.2.3) e o "plano do Capítulo 2".

## Etapa 12 — Cap. 8, Conclusões e acabamento — aplicada no Doc em 04/11/2026

Estrutura: 8.1 Contribuições · 8.2 Trabalhos futuros (6.3) · 8.3 Considerações finais.

- [x] **[Orient.]** Conclusão curta: responder ponto a ponto ao que a introdução prometeu; nada novo; só o que a avaliação demonstrou e a discussão interpretou.
- [x] **[Orient.]** §8.1, contribuições em três grupos: descritiva (taxonomia), prescritiva (instanciação e princípios), evidência sobre automação entre zonas e protocolo de avaliação.
- [x] **[Orient.]** Resumo e abstract com resultados e conclusões (ABNT NBR 6028:2021), com os números da 3.0.3. Escrever por último.
- [x] **[Orient.]** Paginação: cada capítulo começa em página nova, e nenhum título fica sozinho no pé da página (no Google Docs, "Manter com o próximo" e "Quebra de página antes" nos estilos de título); tabelas e quadros sem quebra no meio quando couberem numa página.
- [x] Numeração de páginas, Sumário e listas: ver "Fechamento", no fim deste arquivo.
- [x] **[Orient.]** Conferir as Figuras 1, 2 e 3. **[Achado]** As imagens existem no Doc; a Figura 1 está sem a linha "Fonte".
- [x] **[Achado]** Referências novas: ~~Gregor e Jones (2007)~~ (incluída em 26/10/2026), ~~Denyer, Tranfield e van Aken (2008)~~ (incluída em 04/11/2026) para o CIMO e, se usar, Wieringa (2014); conferir que Gregor e Hevner (2013) e Peffers et al. (2007) seguem com DOI.
- [x] Revisão final de coerência entre os capítulos e numeração de seções, quadros e figuras.
- [x] Aplicado: abertura do Cap. 8 com os seis objetivos e a questão principal; §8.1 sem subseções, em três parágrafos (descritiva, prescritiva, evidência e protocolo); §8.2.1, §8.2.3 (agora "Melhorias no mecanismo de reverso") e §8.2.4 sem o que a 3.0.3 já resolveu; remissões corrigidas (Cap. 4, 6.2.2, 6.4, 3.2.7, 7.2, 7.6.x); Resumo e Abstract com os números da 3.0.3; Figuras 2 e 3 renumeradas na ordem do texto; Fonte da Figura 1; DOI de Gregor e Hevner (2013); Wieringa não entra (não é citado). Corrigida também a cor herdada das notas de rascunho em 17 trechos das §3.1, §3.2.3 e §5.8. Conferido: 17 quadros, 10 imagens, nenhum trecho com a cor das notas.
- [ ] Para o Fechamento: a Lista de Ilustrações ainda tem as Figuras 2 e 3 na ordem antiga.

---

## Adiado

- [ ] Avaliação com participantes: tarefas do E3 no pdnsutil, no SanchezDNS e no PowerDNS-Admin (mesmo PowerDNS); erros publicados (conferidos com dig), tempo e SUS. Antes, confirmar com a coordenação da UVA se precisa de aprovação no CEP (Resolução CNS 510/2016) e TCLE.
- [x] Comentários do orientador no Google Doc (12 comentários; conferidos em 05/11/2026): 2 já resolvidos e 8 atendidos pela reescrita; C12 ("erro impossível") e C4 ("literatura: quais?") fechados com quatro trocas (Quadro 4, §5.1, §5.5.1 e §4.2). O Rafael marca como feitos na interface; sugestão de resposta no C9 antes de resolver. Sumário conferido de novo no PDF (4.6 passou para a p. 49).

---

## Etapa 13 — Enxugar e retirar o Apêndice A (pedido do Rafael em 04/11/2026, antes do Fechamento) — aplicada no Doc em 05/11/2026

- [x] Retirar o Apêndice A e as remissões a ele (§3.1 e, se houver, outras); manter na §3.1 a frase que diz de onde vêm as classificações do Quadro 2. Conferir se as referências ISC, AWS e Google Cloud ainda são citadas em outro lugar; se não forem, retirar também.
- [x] Passada de enxugamento no texto inteiro: cortar repetições entre capítulos, parágrafos que só anunciam o que vem depois e detalhes de implementação que não sustentam nenhuma QP. Mostrar o que sai antes de aplicar.
- [x] Aplicado (backup antes: "BACKUP 2026-11-05 antes da Etapa 13"): Apêndice A e RFC 9110 (sem citação) retirados; as células do Quadro 2 que só o apêndice sustentava ganharam uma frase nas §3.2.1, §3.2.2, §3.2.7 e §3.2.9, de modo que ISC, AWS e Google Cloud seguem citados. Cortes: aberturas dos Caps. 2, 3, 4 e 7 (Caps. 5 e 6 encurtadas), §2.5 e §5.10 inteiras, §5.7.1 absorvida na §5.7, §5.2.2 e §5.9 de três parágrafos para um, repetições (.se, Oppenheimer, cadastro em análise, cache de sessão, correções da 3.0.5, §7.6.1 e §7.6.2 contra §5.4.5 e §5.6.4). Corrigida a remissão da §6.2.2 (seção 5.5 → 6.4). Corpo e quadros passaram de ~33,1 mil para ~30,0 mil palavras. Conferido: 17 quadros, 10 imagens, cada texto novo uma vez.
- [ ] Ficou para decidir: retirar os Quadros 10 e 11 (exige renumerar os Quadros 12 a 17); o 409 do cadastro é mudança de código.

---

## Fechamento — paginação, Sumário e listas (por último, depois de todo o texto pronto) — aplicado no Doc em 05/11/2026

- [ ] **[Orient.]** Numeração de páginas (ABNT NBR 14724:2011) — falta só o Rafael inserir o campo de número (a API não insere): já há quebra de seção antes da Introdução; no cabeçalho de uma página da Introdução, Inserir > Números de página > Mais opções, "Esta seção", começar em 14, à direita.: contadas desde a folha de rosto e impressas a partir da primeira folha textual (Introdução), em algarismos arábicos no canto superior direito.
- [x] **[Orient.]** Sumário refeito com a estrutura de 8 capítulos e os títulos novos (o atual é texto digitado e ainda mostra a estrutura antiga, inclusive os títulos do Cap. 1 de antes de 26/10/2026), com as páginas corretas.
- [x] **[Orient.]** Lista de Ilustrações atualizada: Figuras 1, 2 e 3 nas posições novas (a Figura 3 foi para a §5.2.1), com as páginas corretas.
- [x] Alinhamento: deixar justificados os parágrafos de corpo que ainda estão à esquerda (Caps. 2 a 8, sobretudo texto movido e notas antigas).
- [x] **[Orient.]** Lista de Quadros atualizada: renumerar os quadros na ordem em que aparecem depois da reestruturação (o Quadro 4 vem antes do Quadro 3), corrigir as remissões no texto e as páginas (hoje Quadros 14 e 15 com "??").
- [x] Aplicado em 05/11/2026 (backup antes: "BACKUP 2026-11-05 antes do Fechamento"): papel A4 e margens ABNT (3/3/2/2 cm); capa, folha de rosto e folha de aprovação em páginas próprias e na ordem ABNT; quebra de página antes de cada capítulo e das Referências; "manter com o próximo" nos títulos e legendas; Quadros renumerados na ordem de aparição (4↔3, 10→6, 11→7, 6→8, 7→9, 8→10, 9→11) com as 23 remissões; 92 parágrafos justificados e legendas/fontes à esquerda; Lista de Siglas sem BGP, CLI, NREN e TI e com CBPF, CIMO, CSK, CSRF, DoH, LGPD, NSEC, PHP, SQL e URL; Sumário, Lista de Ilustrações e Lista de Quadros refeitos com as páginas do PDF exportado (contagem a partir da folha de rosto; Introdução na p. 14), conferidos numa segunda exportação (131 de 131 entradas).

## Revisão de 06/11/2026 (lista de pendências conferida item a item)

- [x] TTL fora da taxonomia: N13 reclassificado como política de TTL da plataforma (Quadro 12, §6.2.1); §6.2.4 explica o N9 e conclui que, no erro sintático, o ganho sobre o PowerDNS 5.0.1 é praticamente nulo; ajustados Quadro 16 (tarefas 2 e 5), §6.3, Quadro 17, §7.1 (QP2, QP4 e questão principal), abertura do Cap. 8, Resumo e Abstract.
- [x] §4.4: regra de atribuição alinhada ao plano ("grava um valor diferente do pretendido").
- [x] §1.3 ("À luz dessa premissa"), §4.6 (registro verificável só da terceira rodada; mudança do critério do A9), §5.6.2 (alteração segue a condição da criação), abertura da §6.2 (quatro partes).
- [x] Quebra de página antes da 5.6.3 retirada; Sumário e Lista de Quadros acertados pelo PDF (26 números mudaram).
- [x] Figura 3 substituída pela versão nova (figuras/figura3-sincronizacao-ptr.png, conforme a 3.0.5); figura, legenda e Fonte na mesma página (p. 75); quebra de página antes da 5.6.4, que ficava sozinha no pé da página; Sumário e listas reconferidos pelo PDF (50 números acertados).
- [ ] Opcionais não feitos: Quadro 4 (ordem da linha de rastreabilidade), Figura 5 em três páginas, enxugar §7.1 e §7.6.3.
