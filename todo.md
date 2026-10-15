Inconsistência crítica: a seção 5.3.3 não foi propagada
O Quadro 15 ainda registra a QP3 como "verificado por inspeção; testes contra a API não executados" e a rastreabilidade como "não verificado". A seção 5.3.3, porém, relata dezesseis casos executados, com A8 reprovado, e mede a cobertura da trilha no Quadro 17. O mesmo descompasso aparece em outros cinco lugares:
6.1.3: afirma que os testes de autorização não foram executados.
6.2.1: diz que a fragilidade da sobrescrita na alteração "não foi exercitada pelos episódios", quando A8 a confirmou.
6.2.5: lista os testes de autorização e a triagem de órfãos como substituídos por inspeção, mas A9 a A11 os exercitaram.
6.3.1: propõe como trabalho futuro o que já foi feito.
5.1: diz que a remoção encadeada e a separação de permissões foram verificadas apenas por inspeção.
A seção 5.6 também precisa ser reescrita, porque afirma que as falhas se concentram em P4, P6 e N13 e omite A8 e A9.
Mais sério é o quarto princípio de projeto da seção 6.1.2, que apresenta como contribuição teórica a contenção da escrita derivada pela condição de "não sobrescrever PTR anterior". O caso A8 refuta esse princípio tal como foi implementado. Pela lógica da DSR, um princípio de design emergente refutado pela própria avaliação deve ser reformulado, por exemplo como princípio que exige a mesma guarda em toda operação que grave o registro derivado. Não pode permanecer como contribuição.
Há ainda problemas de forma decorrentes da inserção tardia. Os Quadros 16 e 17 não constam da Lista de Quadros e aparecem antes dos Quadros 14 e 15, quebrando a numeração sequencial.

Achados de A9 e da auditoria que o texto subestima
O caso A9 revela uma classe de erro nova, criada pela própria automação. Quando dois nomes compartilham um endereço e um deles é removido, o PTR sobrevive apontando para o nome removido, e nenhuma das duas triagens o detecta. Isso atinge a QP2 diretamente, na parte "que novos riscos esses mecanismos introduzem", e qualifica o segundo princípio de projeto: a guarda de referência impede a remoção, mas não reaponta o PTR. O Quadro 2 continua atribuindo "Sim" ao SanchezDNS em "remoção encadeada e triagem de órfãos", o que agora é uma autoavaliação mais generosa do que a que a evidência sustenta.
O mesmo vale para a coluna "trilha de auditoria". O Quadro 17 mostra 0 de 5 escritas derivadas registradas, embora o RF13 exija registrar "toda operação de escrita" e o Quadro 9 declare os códigos insert_reverses e delete_reverse. Some-se a isso o campo zone gravando o nome do registro no edit_record. O RF13 está, portanto, parcialmente descumprido, e não apenas "fora do escopo". O fato de A8 sobrescrever o PTR de outra equipe sem deixar rastro atribuível é a combinação mais grave do trabalho e merece destaque na seção 6.2, não uma ressalva no fim de 5.3.3.
Validade de construto e o problema da atribuição causal
O ponto metodológico mais substantivo é de atribuição. A seção 2.6.2 promete uma condição de referência: as mesmas classes de erro introduzidas em zonas editadas diretamente no PowerDNS, sem a plataforma, e submetidas ao mesmo oráculo. Essa condição nunca é relatada no Capítulo 5, e sem ela não se separa o efeito do artefato do efeito do servidor subjacente.
A lacuna importa porque o próprio caso A12 mostra o PowerDNS rejeitando nome fora da zona com 422. Até onde sei, a API REST do PowerDNS também recusa conteúdo não canônico, como um CNAME sem ponto final, e o pdnsutil faz parse do conteúdo por tipo, o que barraria endereços malformados. Recomendo verificar na documentação do PowerDNS e com testes. Se isso se confirmar, parte da classe "erro sintático no RDATA" já está bloqueada na linha de base, e o E3 superestima as classes "ainda possíveis" na linha de comando. A comparação correta é entre a plataforma e a API ou o pdnsutil, não entre a plataforma e a edição livre de arquivos no estilo BIND, que motiva a introdução.
Três fragilidades relacionadas reforçam esse ponto:
Oráculo do E2: os casos N1 a N13 têm como oráculo o valor esperado escrito pelo próprio autor. O Zonemaster, apresentado como oráculo independente, foi aplicado a uma zona com apenas A, NS e TXT, sem exercitar MX, SRV, CAA, HTTPS ou PTR. A independência do oráculo, argumento central da seção 2.5.7, vale portanto para serial e DNSSEC, mas não para a normalização.
Terminologia: vários casos "negativos" são, a rigor, testes de canonicalização positiva.
Métrica do E3: a regra de contagem de passos não foi definida a priori, e um formulário com seis campos conta como um passo enquanto cada comando conta como um. Isso é uma ameaça à validade de construto, no sentido de Wohlin et al. (2012), que deveria constar da seção 2.6.4.
Alinhamento entre questões, objetivos e critérios
A QP2 pergunta por mecanismos "necessários e suficientes". A suficiência é parcialmente testada, mas a necessidade não, pois isso exigiria algo como ablação ou comparação com subconjuntos. Recomendo retirar "necessários" ou justificar o termo analiticamente.
A rastreabilidade aparece sob a QP4, que trata de oportunidades de erro e funcionamento real. O objetivo específico sobre auditoria também está vinculado à QP4, sem que a pergunta mencione auditoria. Isso enfraquece a coerência da matriz do Quadro 3.
A classe "DNSSEC incompleto" inclui a ausência de rotação de chaves, citando os cerca de 70% medidos por Chung et al. (2017), mas é declarada eliminável "no ato da criação", e nenhum mecanismo trata rotação. O escopo da classe deveria ser cortado, ou a rotação declarada fora do escopo.
A afirmação de que os critérios foram "fixados antes da execução" não tem evidência verificável. Um hash de commit ou uma data do plano de avaliação no repositório resolveria isso, o que é especialmente relevante dado que o próprio texto admite, na seção 5.2, que os retornos do E1 não estão registrados no histórico.
Coerência entre motivação e projeto
A seção 1.2 justifica o artefato por redes corporativas ou governamentais restritas que exigem custódia local. O projeto, porém, envia o IP de cada operador a um serviço externo de geolocalização (4.3.1) e depende dos resolvedores DoH da Google e da Cloudflare para o RF03 (4.2.3). Em rede isolada, ambos falham, e o primeiro tem implicação de proteção de dados pessoais que o texto não discute. A seção 2.6.4 promete pseudonimizar domínios e endereços do E4, mas as seções 5.3.2 e 5.5 e as Figuras 4 a 7 expõem curi.dev.br, 152.84.126.7 e o AS 2715. Ou se pseudonimiza, ou se remove a promessa.
Pontos técnicos que a banca pode explorar
O Quadro 6 afirma impedir a enumeração de usuários. Contudo, segundo a seção 4.3.1, o login verifica o cadastro pendente antes de conferir a senha, de modo que qualquer pessoa descobre se um e-mail tem solicitação pendente.
O modelo de sessão por cookie com rotas de escrita (PUT, PATCH, DELETE) não menciona o atributo SameSite nem proteção contra CSRF. Ver o OWASP CSRF Prevention Cheat Sheet e o RFC 6265.
A guarda de referência e a criação do PTR não são atômicas, e a condição de corrida entre operações concorrentes, do tipo TOCTOU, não é discutida.
O caso A12 devolve 500 a um erro de entrada do cliente. O texto já admite o 400 em lugar de 401 (RFC 9110), mas o 500 merece o mesmo registro.
Aspectos formais
O resumo não apresenta resultados nem conclusões, o que a ABNT NBR 6028:2021 exige. Agora há números concretos para isso: 18/19 casos na 3.0.2, 15/16 casos de autorização, 0/5 escritas derivadas auditadas. A Figura 1 aparece vazia no texto extraído, e convém conferir as Figuras 2 e 3. O Poweradmin entra no Quadro 2 sem subseção própria na 2.5. Por fim, a seção 1.1 mantém a frase "Era um erro sintático, do tipo que uma validação [...] teria barrado". Se a API do PowerDNS já barra esse erro, o argumento precisa ser refinado justamente pelo ponto da atribuição causal discutido acima.
Prioridade de revisão
Em ordem de impacto na defesa, recomendo:
Propagar os resultados da 5.3.3 para os Quadros 15 e 2, para as seções 5.1, 5.6 e 6.1.2 a 6.3.1, e para o resumo.
Reformular o quarto princípio de projeto à luz de A8 e qualificar o segundo à luz de A9.
Executar, ou retirar explicitamente, a condição de referência do E2 e reposicionar o E3 contra a API ou o pdnsutil.
Resolver a pseudonimização e a contradição com o cenário on-premises.

As fontes são posteriores ao início da construção. As documentações do Quadro 2 foram consultadas em setembro de 2026, ao fim da construção. obs: nao entendi esse aq

Episódios, métricas, critérios de êxito e ameaças à validade são decisões de método, não fundamentos teóricos

esse problema estah adeqaudo ao metodo?

Pensou por 22s
Na natureza, sim. Na formulação, não inteiramente. O problema tem a forma de um problema de projeto e se encaixa bem na DSR, mas a maneira como foi redigido gera atritos com três exigências do método: a independência entre problema e artefato, a explicitação do contexto e a coerência entre o fenômeno declarado e a estratégia de avaliação.

O que está adequado
O problema trata de como melhorar uma prática, não de explicar um fenômeno, o que o situa nas ciências do artificial (Simon, 1996). A entrada centrada no problema do DSRM (Peffers et al., 2007) é coerente com uma motivação extraída da literatura de medição. O enquadramento de Gregor e Hevner (2013) também se sustenta: o problema é conhecido e a solução proposta é nova, o que caracteriza "melhoria". O ciclo de rigor de Hevner (2007) está presente, porque a taxonomia do Quadro 1 ancora o problema na base de conhecimento e as teorias de Reason (1990) e Norman (2013) funcionam como kernel theories que justificam por que funções de bloqueio deveriam reduzir erros.
Onde a formulação conflita com o método

1. O problema embute a solução. Defini-lo como "a falta de um meio de operação que torne esse erro impossível" torna a avaliação quase tautológica, porque qualquer instanciação "resolve" a falta de si mesma. Na DSR, o problema deve descrever uma classe de situações e um objetivo dos interessados, de forma independente do artefato, para que a avaliação possa refutar a solução. O modelo de Wieringa (2014) explicita essa separação: melhorar um contexto-problema por meio de um artefato que satisfaça certos requisitos, a fim de atingir os objetivos de determinados stakeholders.
2. Faltam o contexto e os interessados. O ciclo de relevância de Hevner (2007) exige um ambiente definido, com pessoas, organizações e sistemas. O problema não diz quem opera, em que tipo de organização, com que divisão de responsabilidades. Esse contexto aparece só disperso no texto: organizações com PowerDNS próprio e infraestrutura local, com equipes distintas para zonas diretas e reversas. Justamente esse contexto justifica a QP3 e a escrita derivada entre zonas, e por isso deveria estar no próprio problema.
3. O fenômeno declarado é humano, mas a avaliação é técnica. Este é o desalinhamento mais sério. O problema é formulado em termos de erro humano: deslizes, enganos e a inevitabilidade defendida por Reason. No FEDS (Venable, Pries-Heje e Baskerville, 2016), a estratégia de avaliação decorre do risco dominante. Um problema definido pelo comportamento de operadores sugere risco humano e, portanto, a estratégia Human Risk and Effectiveness. O trabalho escolhe Technical Risk and Efficacy, com argumento razoável na seção 2.6.1, mas o resultado é um problema sobre pessoas avaliado sem pessoas. Há duas saídas coerentes:
   Reformular o problema em termos técnicos. Por exemplo: os servidores validam a sintaxe de cada registro isoladamente, mas não a consistência entre registros interdependentes nem o alcance da autorização. Nesse caso, a avaliação atual responde ao problema por completo.
   Manter o enquadramento humano. Nesse caso, o trabalho deve declarar desde a seção 1.3 que avalia apenas a eficácia do meio, isto é, a redução de oportunidades de erro, e não a efetividade junto aos operadores.
   A reformulação que sugeri anteriormente segue a primeira via e preserva o erro humano apenas como motivação.
4. Promete-se impossibilidade, mas avalia-se grau. O problema fala em tornar o erro "impossível", que é uma propriedade de garantia. A QP e a avaliação medem "em que medida reduz". A DSR aceita ambas as formas, mas exige que problema, questão e critério falem a mesma língua. Os casos A8 e A9 e a semântica de melhor esforço do RNF03 mostram que o artefato não oferece garantia, de modo que é o problema que precisa ser ajustado.
5. Não há proposição de projeto explícita. Na DSR, o problema deveria conduzir a uma proposição testável que ligue a intervenção ao resultado (Gregor e Jones, 2007). A estrutura CIMO de Denyer, Tranfield e van Aken (2008) ajuda: num contexto C, a intervenção I aciona o mecanismo M e produz o resultado O. Aplicada ao caso:
   Contexto (C): organizações que operam PowerDNS próprio, com zonas diretas e reversas sob equipes distintas.
   Intervenção (I): derivação automática, normalização e permissão por zona.
   Mecanismo (M): retirar do operador as decisões e digitações em que ocorrem deslizes, que funcionam como funções de bloqueio.
   Resultado (O): menos estados inconsistentes publicados.
   Essa forma explicita o que cada episódio testa. O E2 e o E3 testam o mecanismo, o E4 testa o contexto, e o resultado (O) permanece não testado sem operadores.
   Veredito
   O problema é adequado à DSR quanto à natureza e à forma de entrada, mas precisa de três ajustes para que a avaliação de fato o responda: formulá-lo independentemente do artefato, incluir o contexto organizacional e alinhar o fenômeno, seja técnico ou humano, à estratégia do FEDS adotada. Sem esses ajustes, a banca pode objetar que o trabalho responde a uma pergunta diferente da que formulou.

vc mudou a metodologia, mas nao reescreveu o que jah tinha feito apra a metodologia

1. Introdução
   A introdução deve conter o problema formulado independentemente do artefato, com contexto e interessados explícitos, a questão de pesquisa e a natureza da contribuição pretendida. Esta última é a melhoria, nos termos de Gregor e Hevner (2013), com instanciação e princípios de projeto. Aqui entra a reformulação discutida antes: servidores que validam a sintaxe isoladamente, mas não a consistência entre registros interdependentes nem o alcance da autorização, em organizações com equipes distintas para zonas diretas e reversas.
   Os objetivos do trabalho devem ser de pesquisa: caracterizar, projetar, avaliar. Não devem ser funcionalidades. Os objetivos atuais misturam os dois níveis: "implementar sob arquitetura desacoplada" é um requisito do artefato, não um objetivo de pesquisa.
2. Revisão da literatura (base de conhecimento)
   Gregor e Hevner distinguem dois tipos de conhecimento, e o capítulo deveria se organizar em torno dessa distinção. O conhecimento descritivo (Ω) cobre o que se sabe sobre o fenômeno: a arquitetura do DNS, os estudos de medição e a teoria do erro humano como kernel theory. O conhecimento prescritivo (Λ) cobre as soluções já existentes: os trabalhos correlatos e o Quadro 2.
   A taxonomia do Quadro 1 é o fecho natural da parte descritiva. O quadro comparativo é o fecho da parte prescritiva, e a lacuna emerge do confronto entre os dois. O FEDS, o DSRM e o desenho de avaliação não pertencem a este capítulo.
3. Método
   O capítulo de método reúne o que hoje está disperso entre as seções 2.4, 2.6, 3.1 e 3.2:
   a classificação da pesquisa;
   o DSRM com o ponto de entrada;
   a estratégia FEDS escolhida e sua justificativa pelo risco dominante;
   os episódios, a matriz do Quadro 3 com critérios de êxito e as ameaças à validade.
   É aqui que a proposição de projeto deve ser declarada, em forma CIMO (contexto, intervenção, mecanismo, resultado), indicando qual parte cada episódio testa. Se os critérios foram fixados antes da execução, o capítulo deve apontar a evidência disso, como um commit datado do plano de avaliação.

4. Método
   O capítulo de método reúne o que hoje está disperso entre as seções 2.4, 2.6, 3.1 e 3.2:
   a classificação da pesquisa;
   o DSRM com o ponto de entrada;
   a estratégia FEDS escolhida e sua justificativa pelo risco dominante;
   os episódios, a matriz do Quadro 3 com critérios de êxito e as ameaças à validade.
   É aqui que a proposição de projeto deve ser declarada, em forma CIMO (contexto, intervenção, mecanismo, resultado), indicando qual parte cada episódio testa. Se os critérios foram fixados antes da execução, o capítulo deve apontar a evidência disso, como um commit datado do plano de avaliação.
5. Descrição do artefato (projeto e desenvolvimento)
   Este capítulo funde os atuais Capítulos 3 (modelagem) e 4 (implementação) em uma sequência única: objetivos da solução em forma de requisitos rastreáveis à taxonomia, arquitetura, mecanismos e implementação.
   O ponto distintivo em DSR é que o capítulo deve tornar explícitos os princípios de projeto que o artefato instancia. Eles não devem aparecer apenas na conclusão. Gregor e Jones (2007) tratam esses princípios como componentes de uma teoria de design: propósito e escopo, construtos, princípios de forma e função, mutabilidade do artefato e proposições testáveis. Os quatro princípios hoje na seção 6.1.2 deveriam ser enunciados aqui, cada um ligado ao mecanismo que o realiza e ao caso de teste que pode refutá-lo.
   Os retornos do ciclo construir-avaliar (guarda de referência e exibição do valor publicado) também ficam aqui, como evidência da natureza iterativa do processo.
6. Demonstração e avaliação
   Os episódios devem ser relatados na ordem do método, cada um confrontado com seu critério. Peffers et al. (2007) separam demonstração de avaliação, e a separação é útil neste caso. A demonstração mostra que o artefato resolve instâncias do problema; é o papel da zona publicada e da implantação no E4. A avaliação mede o quanto o artefato atinge os objetivos, que é o papel dos casos negativos, dos testes de autorização e da comparação do E3.
   As reprovações (N13, P4, P6, A8, A9) e as lacunas de execução devem ser reportadas aqui sem interpretação. A interpretação fica para o capítulo seguinte.
7. Discussão
   Este é o capítulo que o TCC atual não tem, e é o mais característico do esquema de Gregor e Hevner. Nele se interpretam os resultados à luz da base de conhecimento: o que os casos A8 e A9 dizem sobre a dificuldade de combinar automação e autorização granular, à luz do incidente do Technitium. Discute-se também o que sobrevive dos princípios de projeto após a avaliação. Aqui o quarto princípio seria reformulado, e o segundo, qualificado.
   A discussão deve ainda reposicionar o artefato em relação ao Quadro 2 com base na evidência, e não na autoavaliação. Deve separar com clareza o que foi demonstrado como eficácia do que permanece como hipótese de efetividade, retomando a distinção de Venable, Pries-Heje e Baskerville (2016). A análise das ameaças à validade pós-execução também pertence a este capítulo.
8. Conclusões
   A conclusão deve ser curta e responder ponto a ponto ao que a introdução prometeu: a resposta à questão, as contribuições nos dois níveis, as limitações e os trabalhos futuros. Nada de novo deve aparecer aqui. A conclusão só pode afirmar o que a avaliação demonstrou e a discussão interpretou.
   Mudanças em relação à versão atual
   No TCC do Rafael, as mudanças estruturais seriam três:
   Transferir a seção 2.6 e a 2.4.3 para um capítulo de método único.
   Fundir os Capítulos 3 e 4, eliminando a repetição que hoje existe entre modelagem e implementação, como a derivação do reverso descrita duas vezes nas seções 3.6.1 e 4.5.1.
   Criar o capítulo de discussão, para onde migram os princípios reavaliados e boa parte do que hoje está nas seções 6.1 e 6.2.

RAFAEL, seu trabalho tá mais ou menos pronto, mas precisa reescrever pro DSR, ele começa sem colocar o DSR no meio do projeto. O que talvez precise é incluir gente no processo de avaliacao. chamar gente pra fazer o teste do artefato para comparar com outros que vc citou.

Contribuições declaradas
Contribuição prática, de nível 1: a instanciação. O SanchezDNS é uma camada de operação sobre a API do PowerDNS, distribuída como imagem de contêiner, documentada e licenciada sob MIT. A seção 2.5.9 situa seu valor na combinação, num único fluxo e sem trocar o servidor autoritativo, de mecanismos hoje espalhados por ferramentas que competem pelo papel de servidor ou de fonte da verdade.
Essa contribuição se sustenta. O E4 demonstrou implantação apenas com a documentação e operação sobre um servidor real, e o Zonemaster e o pdnsutil confirmaram o serial e o DNSSEC na zona publicada. A ressalva é que o Quadro 2 atribui "Sim" ao artefato em todas as colunas, enquanto a avaliação mostra cumprimento parcial em três delas: remoção encadeada (A9), auditoria (escritas derivadas não registradas) e contenção da escrita (A8). A contribuição vale, mas a autoavaliação comparativa precisa ser rebaixada.
Contribuição teórica, de nível 2: os princípios de projeto da seção 6.1.2. Os quatro princípios têm destinos diferentes depois da avaliação:
Canonicalizar na edição e rejeitar apenas o que não admite conversão determinística: se sustenta. O caso N13 o refina: a regra precisa residir no servidor da aplicação, não no cliente.
Derivar automaticamente, com criação otimista e remoção guardada: se sustenta parcialmente. O caso P4 mostrou que a derivação precisa respeitar a estrutura rótulo a rótulo do espaço de nomes. O caso A9 mostrou que a guarda impede a remoção, mas não mantém a correspondência, porque o PTR sobrevivente aponta para um nome inexistente. O princípio precisa incluir o reapontamento.
Exibir o estado publicado, não o digitado: se sustenta, embora sem teste específico. A evidência vem da consistência entre a interface e o resolvedor público.
Conter a escrita derivada por condições controladas pelo dono do recurso: refutado como implementado pelo caso A8. Pode ser reformulado como "a guarda deve valer para toda operação que grave o registro derivado, não apenas para a criação", e assim reformulado se torna uma lição mais valiosa do que a versão original.
Os princípios são uma contribuição real, mas emergente no sentido estrito: derivam de uma única instanciação e foram parcialmente refutados e corrigidos pela própria avaliação.
Contribuições não reconhecidas pelo próprio trabalho
A taxonomia do Quadro 1. Ela organiza duas décadas de estudos de medição em dez classes e as separa pelo critério de eliminabilidade por construção numa camada de operação, em oposição às que só admitem detecção. Isso é conhecimento descritivo independente do artefato e reutilizável por quem projete ou avalie qualquer ferramenta de gestão de zonas. O texto a trata como etapa intermediária, mas ela merece ser declarada como contribuição.
A evidência empírica sobre a interação entre automação e autorização granular. Os casos A8 e A12 reproduzem, num artefato projetado especificamente para evitá-la, a mesma classe de falha corrigida no Technitium 15.5: a escrita derivada que atravessa a fronteira de permissão. Isso sustenta uma proposição mais forte do que o texto atual faz. A fronteira de autorização precisa ser verificada em todo caminho de escrita derivada (criação, alteração e remoção), e sua violação tende a ocorrer nos caminhos menos óbvios, como a alteração. Junto com o caso A9, que expõe um estado inconsistente invisível às duas triagens, isso é conhecimento prescritivo novo, obtido justamente por causa do rigor da avaliação.
O protocolo de avaliação, como contribuição metodológica menor. Ele reúne casos negativos derivados de uma taxonomia prévia, oráculos independentes (Zonemaster, pdnsutil, dig) e testes de autorização de ponta a ponta com zonas descartáveis em espaços reservados pelas RFC 2606 e 5737. Os testes estão versionados no repositório, o que permite replicar a avaliação em outras ferramentas.
