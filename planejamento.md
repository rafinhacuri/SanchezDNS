# Planejamento da revisão da monografia e do artefato

Checklist das mudanças pedidas pelo orientador (revisão de 15/10/2026), conferidas contra o texto do Google Doc e o código em 15/10/2026, mais os achados da conferência.

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
- [ ] Release 3.0.3.

## Etapa 2 — Código: segurança e on-premises

- [x] **[Orient.]** Cookie de sessão sem `SameSite` (nem proteção CSRF). Definir `SameSite=Lax` ou `Strict`. Rotas POST (`/user`, `/logout`) aceitam corpo JSON sem checar o Content-Type.
- [x] **[Orient.]** Enumeração: o login responde "cadastro em análise" antes de conferir a senha. A solicitação guarda o hash da senha, então dá para conferir antes de revelar.
- [x] **[Orient.]** Geolocalização por serviço externo: tornar opcional por variável de ambiente (desligada por padrão). **[Achado]** Hoje usa `http://ip-api.com` sem TLS e o User-Agent `cbpf-sso/1.0`.
- [x] **[Orient.]** RF03 depende do DoH da Google e da Cloudflare: tornar os resolvedores configuráveis (rede isolada).
- [x] **[Achado]** Listagem de zonas exclui duas zonas por nome fixo (resíduo de teste, §4.3.3). Remover.
- [x] **[Orient.]** Sessão inválida responde 400 em vez de 401 (RFC 9110). Corrigido para 401 (o frontend não dependia do 400; script do E2 aceita os dois).
- [ ] **[Achado]** `PUT /file` e `GET /file/:id` fora do grupo autenticado, o que contradiz o RNF04 ("não existe caminho de escrita não autenticado"). O `PUT /file` é usado no cadastro, antes de existir sessão, então não dá para autenticar. Ajustar o RNF04 no texto (Etapa 8) para "nenhum caminho de escrita sobre dados de zona sem autenticação".

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
- [ ] **[Achado]** `records/nome-completo.go` compara o nome com a zona como texto, sem fronteira de rótulo (mesma classe do P4): com a zona `curi.dev.br`, o nome `xcuri.dev.br` vira FQDN fora da zona. Corrigir numa 3.0.4 e criar caso de teste.
- [ ] **[Achado]** O E2 (3.0.2 e 3.0.3) rodou no ambiente de desenvolvimento desta máquina (`bun dev`, `mizle.cbpf.dev.br:3000`) ligado ao PowerDNS do E4, e não na imagem publicada. A §5.3.3 diz "sobre a instância do episódio E4": corrigir no texto (Etapa 9).
- [ ] **[Achado]** API inconsistente no MX: `PUT /records` recebe a prioridade em campo próprio; `DELETE /records` exige a prioridade dentro do valor (`"10 mail..."`). Uniformizar numa 3.0.4.
- [ ] **[Achado]** `DELETE /records` responde 200 quando nenhum registro casa com o valor enviado (nada é removido). Devolver 404 numa 3.0.4.
- [ ] **[Achado]** No `pdnsutil` 5.0.1 o nome do registro é absoluto (nome relativo é recusado como fora da zona) e o conteúdo passa pelo mesmo parser da API. A "edição manual" da linha de base não reproduz o erro de ponto final do .se: o argumento da §1.1 e o E3 precisam refletir isso.
- [ ] **[Achado]** O campo "Email do responsável" vai para o RNAME do SOA sem converter o `@` (`api/zonas/update-soa.go`). Verificar se um e-mail comum produz RNAME malformado; se sim, converter numa 3.0.4.

## Etapa 4 — Nova estrutura da monografia (esquema de Gregor e Hevner, 2013)

- [ ] **[Orient.]** Reorganizar em: 1 Introdução · 2 Base de conhecimento (Ω descritivo / Λ prescritivo) · 3 Método · 4 Artefato · 5 Demonstração e avaliação · 6 Discussão · 7 Conclusões.
- [ ] **[Orient.]** Mover a §2.4 (DSR), a §2.4.3 (FEDS) e a §2.6 (estratégia de avaliação) para o Cap. 3 (Método), junto com a §3.1 e a §3.2.
- [ ] **[Orient.]** Fundir os Caps. 3 e 4 atuais, eliminando repetições (§3.6.1 × §4.5.1, §3.6.2 × §4.5.1, §3.4.1 × §4.3.1, §3.7 × §4.6).
- [ ] **[Orient.]** Criar o capítulo de Discussão, para onde migra boa parte das §6.1 e §6.2.
- [ ] Mapa "seção atual → seção nova" para mover os blocos no Docs.

## Etapa 5 — Cap. 1, Introdução

- [ ] **[Orient.]** Problema independente do artefato. Hoje diz "a falta de um meio de operação que torne esse erro impossível", o que embute a solução e torna a avaliação tautológica.
- [ ] **[Orient.]** Contexto e interessados explícitos: organizações com PowerDNS próprio e equipes distintas para zonas diretas e reversas.
- [ ] **[Orient.]** Enquadramento técnico: servidores validam a sintaxe de cada registro isoladamente, mas não a consistência entre registros interdependentes nem o alcance da autorização. O erro humano fica como motivação.
- [ ] **[Orient.]** Falar em grau ("em que medida reduz"), não em "impossível".
- [ ] **[Orient.]** QP2: retirar "necessários" (a necessidade não é testada; exigiria ablação).
- [ ] **[Orient.]** Rastreabilidade está sob a QP4, que não fala de auditoria. Realinhar a QP, o objetivo específico e o Quadro 3.
- [ ] **[Orient.]** Objetivos de pesquisa (caracterizar, projetar, avaliar), não funcionalidades. Tirar "implementar sob arquitetura desacoplada" dos objetivos.
- [ ] **[Orient.]** Natureza da contribuição pretendida: melhoria (Gregor e Hevner, 2013), com instanciação e princípios de projeto; declarar também a taxonomia e o protocolo de avaliação.
- [ ] **[Orient.]** §1.1, "Era um erro sintático, do tipo que uma validação [...] teria barrado": refinar conforme o resultado da condição de referência.
- [ ] **[Achado]** §1.3 "outras dependem de evidência empírica com operadores humanos" e "revisão sistemática da literatura" (não houve protocolo sistemático): restos do desenho antigo.
- [ ] **[Achado]** QP principal "registre integralmente" e objetivo geral "reduza a ocorrência": verbos mais fortes do que a avaliação sustenta.
- [ ] **[Orient.]** §1.2: o argumento on-premises precisa ficar coerente com a geolocalização externa e o DoH (ver Etapa 2).

## Etapa 6 — Cap. 2, Base de conhecimento

- [ ] **[Orient.]** Parte descritiva (Ω): DNS, estudos de medição, erro humano como kernel theory, controle de acesso; fecho no Quadro 1.
- [ ] **[Orient.]** Parte prescritiva (Λ): trabalhos correlatos; fecho no Quadro 2; a lacuna surge do confronto entre os dois.
- [ ] **[Orient.]** Classe "DNSSEC incompleto" inclui a rotação de chaves (Chung et al.), mas nenhum mecanismo trata rotação: cortar a classe ou declarar a rotação fora do escopo.
- [ ] **[Orient.]** Poweradmin aparece no Quadro 2 sem subseção própria na 2.5: criar.
- [ ] **[Achado]** §2.2.4: "para que a redução de erros possa ser medida por categoria e no total" não é mais medido.
- [ ] **[Achado]** §2.5.6: "o SanchezDNS precisa emular a atomicidade do Knot" contradiz o RNF03 (melhor esforço).
- [ ] **[Achado]** "As documentações do Quadro 2 foram consultadas em setembro de 2026": manter a data, mas ela é usada no ponto de entrada do DSRM (Etapa 7).

## Etapa 7 — Cap. 3, Método

- [ ] **[Orient.]** Classificação da pesquisa (atual §3.1).
- [ ] **[Orient.]** Ponto de entrada do DSRM: a §3.2 diz "centrada no problema", mas o código começou em out/2025 e a literatura e as fontes do Quadro 2 são de set/2026. Declarar entrada centrada em projeto e desenvolvimento (Peffers et al., 2007).
- [ ] **[Orient.]** Estratégia FEDS (Technical Risk and Efficacy) justificada pelo risco dominante; declarar que se avalia a eficácia do meio, não a efetividade junto a operadores.
- [ ] **[Orient.]** Proposição de projeto em forma CIMO (contexto, intervenção, mecanismo, resultado), indicando qual parte cada episódio testa.
- [ ] **[Orient.]** Episódios, Quadro 3 com critérios de êxito e ameaças à validade (atual §2.6).
- [ ] **[Orient.]** Evidência de que os critérios foram fixados antes da execução: commit datado do plano (Etapa 3). **[Achado]** Remover "definidos previamente à construção" (§1.6, §2.7), que é falso.
- [ ] **[Orient.]** Ameaça à validade de construto: regra de contagem de passos do E3 (formulário de seis campos = 1 passo; cada comando = 1 passo) (Wohlin et al., 2012).
- [ ] **[Orient.]** Ameaça: oráculo dos casos N é o valor esperado escrito pelo autor; o Zonemaster só é independente para serial e DNSSEC (até a Etapa 3 ampliar a zona).
- [ ] **[Orient.]** Terminologia: vários casos "negativos" são testes de canonicalização positiva. Separar canonicalização (N1–N10) de rejeição (N11–N13).
- [ ] **[Orient.]** Condição de referência: descrever como executada (Etapa 3) ou retirar a promessa.
- [ ] **[Orient.]** Pseudonimização: a §2.6.4 promete, mas o texto e as Figuras 4–7 expõem curi.dev.br, 152.84.126.7 e o AS 2715. Decidir entre pseudonimizar ou retirar a promessa. **[Achado]** O servidor fica em infraestrutura do CBPF.
- [ ] **[Achado]** Remover restos do desenho antigo: título "Framework FEDS e Escala SUS" (§2.4.3) e Quadro 4 "Demonstração ... com conferência por oráculo externo".

## Etapa 8 — Cap. 4, Artefato (fusão dos Caps. 3 e 4)

- [ ] **[Orient.]** Requisitos rastreáveis à taxonomia, arquitetura, mecanismos e implementação numa sequência única.
- [ ] **[Orient.]** Enunciar os quatro princípios de projeto aqui (não só na conclusão), cada um ligado ao mecanismo que o realiza e ao caso de teste que pode refutá-lo (Gregor e Jones, 2007).
- [ ] **[Orient.]** Manter aqui os retornos do ciclo construir-avaliar (guarda de referência; exibição do valor publicado).
- [ ] **[Orient.]** Quadro 6 afirma impedir a enumeração de usuários: ajustar à Etapa 2.
- [ ] **[Orient.]** Mencionar SameSite/CSRF no modelo de sessão.
- [ ] **[Achado]** Quadro 7 omite `DELETE /reverses`, que o caso A7b usa.
- [ ] **[Achado]** RNF04: restringir a "nenhum caminho de escrita sobre dados de zona sem autenticação" (`PUT /file` e `POST /cadastro` são escritas sem sessão, fora das zonas).
- [ ] **[Achado]** RF10 diz "criação atômica em lote", mas a §3.6.4 diz não transacional.
- [ ] **[Achado]** §3.7 diz "cada operação de escrita grava um documento" e que o ciclo de contas fica fora da trilha: alinhar com o RF13 e com a 3.0.3.
- [ ] Descrever as mudanças da 3.0.3 (guardas, reapontamento, auditoria).

## Etapa 9 — Cap. 5, Demonstração e avaliação

- [ ] **[Orient.]** Separar demonstração (zona publicada, E4) de avaliação (casos N e A, E3, condição de referência).
- [ ] **[Orient.]** Relatar os episódios na ordem do método, cada um contra o seu critério, sem interpretação; a interpretação vai para a Discussão.
- [ ] **[Orient.]** Reportar as reprovações (N13, P4, P6, A8, A9) e as lacunas de execução.
- [ ] Atualizar com os resultados da Etapa 3 (3.0.3, referência, Zonemaster ampliado, E3 recontado).
- [ ] **[Orient.]** A12: registrar o 500 devolvido a um erro do cliente (ou a correção da 3.0.3).
- [ ] **[Orient.]** RF13 parcialmente descumprido (0/5 escritas derivadas; contas fora da trilha), não "fora do escopo".
- [ ] **[Achado]** Quadro 17, linha QP1: "casos de teste executados para parte delas" contradiz a §5.1 ("todas").
- [ ] ~~**[Orient.]** Quadro de resultados com QP3 "não executado" e rastreabilidade "não verificado"~~ (já atualizado no Quadro 17).
- [ ] ~~**[Orient.]** §5.1 dizendo "apenas por inspeção"~~ (já atualizado).
- [ ] **[Orient.]** §5.6: incluir o A9 no parágrafo (o Quadro 17 já cita; o texto não).

## Etapa 10 — Cap. 6, Discussão

- [ ] **[Orient.]** A8 e A9 à luz do incidente do Technitium 15.5: dificuldade de combinar automação e autorização granular.
- [ ] **[Orient.]** A8 sobrescrevendo o PTR de outra equipe sem rastro na trilha é a combinação mais grave: destacar, não deixar como ressalva.
- [ ] **[Orient.]** A9 é uma classe de erro nova criada pela automação: atinge a QP2 ("novos riscos").
- [ ] **[Orient.]** Destino dos princípios: 1º se sustenta (N13 o refina: regra no servidor); 2º parcial (P4 + A9: precisa do reapontamento); 3º se sustenta (com o teste da Etapa 3); 4º refutado como implementado (A8) e reformulado como "a guarda vale para toda operação que grava o registro derivado".
- [ ] **[Orient.]** Reposicionar o SanchezDNS no Quadro 2 pela evidência: "Parcial" em remoção encadeada e triagem (A9), permissão por zona (A8) e trilha (escritas derivadas), ou "Sim" com nota se a 3.0.3 for aprovada.
- [ ] **[Orient.]** Contribuições não reconhecidas: taxonomia do Quadro 1 (com ressalva: não é revisão sistemática); evidência sobre automação × autorização; protocolo de avaliação.
- [ ] **[Achado]** Corrigir a afirmação de que "A8 e A12 reproduzem a falha do Technitium": só o A8 reproduz; no A12 a proteção veio do PowerDNS, não da camada de autorização.
- [ ] **[Achado]** "A violação tende a ocorrer nos caminhos menos óbvios": apresentar como hipótese (dois casos), não como proposição demonstrada.
- [ ] **[Achado]** Protocolo "replicável em outras ferramentas": o script é específico da API do SanchezDNS; replicável é o protocolo, não o script.
- [ ] **[Orient.]** Separar a eficácia demonstrada da efetividade, que continua hipótese (Venable, Pries-Heje e Baskerville, 2016).
- [ ] **[Orient.]** Condição de corrida (TOCTOU) entre a guarda de referência e a criação do PTR.
- [ ] **[Orient.]** Geolocalização externa (dados pessoais, LGPD) e dependência de DoH em rede isolada.
- [ ] **[Orient.]** Ameaças à validade depois da execução.
- [ ] **[Achado]** §6.2.5: "Parte do plano também ficou por executar" seguido de "foram executados depois" e "não verificados" (o Quadro 17 não tem mais esse status).

## Etapa 11 — Cap. 7, Conclusões e acabamento

- [ ] **[Orient.]** Conclusão curta: responder ponto a ponto ao que a introdução prometeu; nada novo; só o que a avaliação demonstrou e a discussão interpretou.
- [ ] **[Orient.]** Resumo e abstract com resultados e conclusões (ABNT NBR 6028:2021): 18/19 na 3.0.2, 15/16 de autorização, 0/5 escritas derivadas auditadas (atualizar com a 3.0.3). Escrever por último.
- [ ] **[Orient.]** Lista de Quadros: Quadros 14 e 15 com página "??".
- [ ] **[Orient.]** Conferir as Figuras 1, 2 e 3. **[Achado]** As imagens existem no Doc; a Figura 1 está sem a linha "Fonte".
- [ ] **[Achado]** Referências novas: Gregor e Jones (2007), Denyer, Tranfield e van Aken (2008) para o CIMO e, se usar, Wieringa (2014).
- [ ] Revisão final de coerência entre os capítulos e numeração de seções, quadros e figuras.

---

## Adiado

- [ ] Avaliação com participantes: tarefas do E3 no pdnsutil, no SanchezDNS e no PowerDNS-Admin (mesmo PowerDNS); erros publicados (conferidos com dig), tempo e SUS. Antes, confirmar com a coordenação da UVA se precisa de aprovação no CEP (Resolução CNS 510/2016) e TCLE.
- [ ] Comentários do orientador no Google Doc (10 comentários de 15/10/2026): tratar depois.
