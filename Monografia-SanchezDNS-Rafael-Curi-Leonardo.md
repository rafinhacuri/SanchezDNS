UNIVERSIDADE VEIGA DE ALMEIDA – UVA
BACHARELADO EM CIÊNCIA DA COMPUTAÇÃO

SANCHEZDNS: PLATAFORMA WEB PARA OPERAÇÃO DE DNS AUTORITATIVO COM CONTROLE DE ACESSO POR ZONA, AUTOMAÇÃO DE REGISTROS E TRILHA DE AUDITORIA

RAFAEL CURI LEONARDO

RIO DE JANEIRO
2026
UNIVERSIDADE VEIGA DE ALMEIDA - UVA
RAFAEL CURI LEONARDO

Monografia apresentada ao curso de Ciência da Computação da Universidade Veiga de Almeida, como requisito parcial para obtenção do título de Bacharel em Ciência da Computação.
Orientador: Prof. Carlos Alvaro de Macedo Soares Quintella

SANCHEZDNS: PLATAFORMA WEB PARA OPERAÇÃO DE DNS AUTORITATIVO COM CONTROLE DE ACESSO POR ZONA, AUTOMAÇÃO DE REGISTROS E TRILHA DE AUDITORIA

RIO DE JANEIRO
2026
UNIVERSIDADE VEIGA DE ALMEIDA - UVA
BACHARELADO EM CIÊNCIA DA COMPUTAÇÃO
RAFAEL CURI LEONARDO

SANCHEZDNS: PLATAFORMA WEB PARA OPERAÇÃO DE DNS AUTORITATIVO COM CONTROLE DE ACESSO POR ZONA, AUTOMAÇÃO DE REGISTROS E TRILHA DE AUDITORIA

Monografia apresentada como requisito parcial à conclusão do curso em Bacharel em Ciência da Computação.

APROVADA EM: \_**\_ / \_\_** / **\_\_\_\_**
CONCEITO: **********\_\_\_\_**********

BANCA EXAMINADORA:

---

PROF. CARLOS ALVARO DE MACEDO SOARES QUINTELLA
ORIENTADOR

---

PROF. DSc ou MSc NOME COMPLETO DO PROFESSOR DA BANCA

---

PROF. DSc ou MSc NOME COMPLETO DO PROFESSOR DA BANCA

Coordenação de Ciência da Computação
Rio de Janeiro

RESUMO
O Sistema de Nomes de Domínio (Domain Name System – DNS) é uma infraestrutura crítica para o funcionamento da Internet. Sua especificação original priorizou a escalabilidade hierárquica e a redundância distribuída, delegando a integridade das zonas autoritativas à manutenção manual por operadores humanos. Estudos de medição em larga escala evidenciam que erros de configuração nessas zonas são recorrentes e comprometem a robustez, a latência e a disponibilidade dos serviços dependentes. Este trabalho apresenta o desenvolvimento e a avaliação do SanchezDNS, uma plataforma web que atua como camada de operação sobre a API REST do PowerDNS Authoritative Server, mantido como fonte única da verdade dos dados de zona. O artefato busca reduzir falhas humanas por meio de três mecanismos: (i) automação de registros deriváveis, como a criação e remoção de registros reversos (PTR) a partir de registros A e AAAA, a normalização de valores por tipo de registro, a delegação do incremento do serial SOA ao servidor autoritativo e a ativação automática de DNSSEC; (ii) controle de acesso combinando papéis globais e permissões de leitura e escrita atribuídas por zona; e (iii) trilha de auditoria centralizada das operações sobre zonas, registros e usuários. A pesquisa segue a metodologia Design Science Research (DSR), estruturada conforme o modelo de processo de Peffers et al. (2007), com avaliação conduzida segundo o Framework for Evaluation in Design Science (FEDS). A avaliação compara a operação pela plataforma com a edição manual quanto à taxa de erros de configuração e ao tempo de execução de alterações, e mede a usabilidade percebida por meio da System Usability Scale (SUS).

Palavras-chave: DNS autoritativo; erro de configuração; PowerDNS; controle de acesso por zona; Design Science Research.

SUMÁRIO
// todo: ajeitar sumário no final para cada página
1 INTRODUÇÃO .................................................................................. 1
1.1 Contextualização e Problematização do Domínio ......................................... 1
1.2 Defasagens do Mercado e Brechas das Soluções Existentes ............................... 4
1.3 Problema de Pesquisa e Questões Norteadoras ........................................... 6
1.4 Objetivos do Trabalho ................................................................. 8
1.4.1 Objetivo Geral .................................................................. 8
1.4.2 Objetivos Específicos ........................................................... 8
1.5 Delimitação do Escopo ................................................................. 9
1.6 Estrutura da Monografia .............................................................. 10

1 INTRODUÇÃO
A introdução deste Trabalho de Conclusão de Curso apresenta o recorte contextual, as justificativas acadêmicas e a formulação metodológica que norteiam a concepção da plataforma SanchezDNS. O objetivo central consiste em investigar e prescrever diretrizes computacionais capazes de mitigar a incidência de falhas operacionais em servidores de nomes autoritativos por meio de abstração de interface, controle granular de permissões e automação de rotinas de validação e sincronização de registros.
1.1 CONTEXTUALIZAÇÃO E PROBLEMATIZAÇÃO DO DOMÍNIO
O Sistema de Nomes de Domínio (Domain Name System – DNS) é uma infraestrutura hierárquica e distribuída essencial para a tradução entre identificadores mnemônicos e endereços IP na Internet. A especificação original do protocolo (MOCKAPETRIS, 1987a, 1987b) priorizou a escalabilidade hierárquica e a resiliência de consultas por meio de redundância distribuída, confiando a integridade das zonas autoritativas à edição e manutenção humana dos arquivos de configuração.
Estudos empíricos de medição de redes apontam que inconsistências operacionais e erros de configuração na administração de zonas autoritativas são amplamente difundidos na infraestrutura global de TI. Bojović e Gajin (2017), ao analisarem cerca de 80 mil domínios distribuídos entre redes acadêmicas internacionais (NRENs), o ccTLD sérvio (.rs) e os domínios mais populares da Internet, identificaram taxas relevantes de falhas de configuração, sobretudo nos segmentos acadêmico e nacional. Entre os problemas diagnosticados destacam-se a presença de servidores não autoritativos declarados na zona pai (lame delegation), ocorrendo em quase 5% dos domínios acadêmicos avaliados, e a ausência completa de registros de ponteiro reverso (PTR) para os servidores de nomes em cerca de 40% a 45% dos domínios de cada grupo analisados.
A permanência dessas falhas de configuração gera degradações severas de infraestrutura. Conforme evidenciado no estudo clássico de Pappas et al. (2004), inconsistências operacionais comprometem a robustez do ecossistema: os autores constataram lame delegations em aproximadamente 15% das zonas analisadas e dependências circulares em 2%, demonstrando que essas falhas causam a degradação da redundância planejada e elevam a latência de resolução de consultas em até uma ordem de grandeza.
Dentre os pontos críticos de manutenção manual sobressai o gerenciamento do número serial (SOA Serial). Por coordenar a sincronização de dados entre instâncias primárias e nós replicadores via transferências de zona ou mecanismos de replicação de retaguarda, lapsos no incremento manual ou violações das regras de aritmética de números de série (ELZ; BUSH, 1996) provocam dessincronização de dados (zone drift), fazendo com que servidores secundários continuem respondendo consultas com versões desatualizadas da zona.
Adicionalmente, a ausência de mecanismos preventivos de validação antes da publicação de alterações potencializa falhas catastróficas. A gravidade da replicação instantânea de erros ficou evidente no incidente operacional do ccTLD da Suécia (.se), ocorrido em 13 de outubro de 2009 (PINGDOM, 2009). Um script automatizado de manutenção omitiu o ponto final nos registros de delegação dos servidores de nomes; uma vez submetida a configuração, a zona sintaticamente corrompida foi imediatamente propagada para toda a rede secundária, invalidando a resolução de pouco mais de 900 mil domínios sob o sufixo .se por cerca de uma hora e gerando impactos prolongados pelo armazenamento em cache dos resolvedores (PINGDOM, 2009). O episódio evidencia a importância de mecanismos que validem e normalizem os dados estruturalmente antes de sua publicação definitiva no servidor autoritativo.
As falhas de configuração não se restringem, contudo, à sintaxe dos arquivos de zona. Uma segunda classe decorre da divergência entre a zona pai e a zona filha. Sommese et al. (2020), ao compararem os conjuntos de registros NS publicados por pais e filhos nos domínios de segundo nível de .com, .net e .org, identificaram inconsistência em mais de treze milhões de domínios, entre os quais ao menos cinquenta mil integrantes da lista dos um milhão mais acessados. Como o resolvedor pode legitimamente utilizar qualquer um dos dois conjuntos, a divergência não produz falha determinística, mas indisponibilidade intermitente, cujo diagnóstico é particularmente custoso justamente porque a consulta de verificação do operador frequentemente obtém resposta correta.
Uma terceira classe relaciona-se ao ciclo de vida dos registros, e não à sua criação. Registros que permanecem publicados após o desprovisionamento do recurso para o qual apontam — denominados registros pendentes (dangling records) — convertem um resíduo operacional em vetor de ataque. Liu, Hao e Wang (2016) identificaram 467 registros pendentes exploráveis em 277 domínios entre os dez mil mais acessados e em 52 zonas do domínio .edu, demonstrando que um terceiro pode assumir o controle do subdomínio correspondente e até obter certificados TLS válidos em nome da organização titular. O problema é especialmente agudo nos registros reversos, cuja exclusão não acompanha automaticamente a remoção do registro A ou AAAA que lhes deu origem, uma vez que residem em zona distinta.
Uma quarta classe concentra-se na operação das extensões de segurança. Chung et al. (2017), em série histórica de vinte e um meses sobre a totalidade dos domínios assinados em .com, .net e .org, verificaram que mais de 30% das zonas assinadas não possuíam o registro DS correspondente publicado na zona pai — o que rompe a cadeia de confiança e torna a assinatura inócua — e que cerca de 70% jamais efetuaram rotação de chaves. O achado indica que o obstáculo à adoção efetiva do DNSSEC não reside no mecanismo criptográfico, mas na condução correta da sequência de etapas operacionais que o cercam. De modo convergente, Korczyński et al. (2016) encontraram, em amostra aleatória de 2,9 milhões de domínios, 1.877 domínios (0,065%) cujos servidores autoritativos aceitavam atualizações dinâmicas não autenticadas, permitindo a inserção de registros arbitrários por qualquer agente da Internet; ainda que percentualmente reduzida, a falha é integral onde ocorre, pois anula por completo o controle de acesso à zona.
O impacto dessas falhas transcende o domínio diretamente afetado. Por operar como precondição de praticamente toda comunicação na Internet, o DNS converte erros localizados em indisponibilidade sistêmica. Em 22 de julho de 2021, uma atualização de configuração no serviço Edge DNS da Akamai tornou inacessíveis serviços de larga escala por pouco mais de uma hora, exigindo a reversão da alteração (AKAMAI, 2021). Em 4 de outubro do mesmo ano, um comando de manutenção que retirou de operação a rede de backbone da Meta provocou a retirada dos anúncios BGP dos prefixos em que residiam seus servidores autoritativos, deixando Facebook, Instagram e WhatsApp inacessíveis por mais de sete horas (META, 2021). Em ambos os episódios, a causa imediata não foi falha de equipamento nem ação adversária, mas uma alteração de configuração conduzida por operadores, cujos efeitos se propagaram antes que houvesse oportunidade de detecção — o mesmo padrão observado no incidente do ccTLD sueco.
Esse conjunto de evidências delimita o problema tratado nesta pesquisa. A superfície de erro do DNS autoritativo não decorre primordialmente da complexidade do protocolo, mas da distância entre o modelo mental do operador e a representação textual que ele precisa manipular, agravada por três propriedades do meio: a ausência de verificação no momento da edição, a irreversibilidade prática decorrente do armazenamento em cache e a inexistência de mecanismos que tornem determinadas classes de engano impossíveis, e não apenas detectáveis. A seção 2.2 aprofunda a caracterização dessas classes, de suas evidências empíricas e de seus impactos, sintetizando-as em uma taxonomia que orienta tanto o projeto quanto a avaliação do artefato.
1.2 DEFASAGENS DO MERCADO E BRECHAS DAS SOLUÇÕES EXISTENTES
As ferramentas atuais voltadas à administração de servidores autoritativos dividem-se predominantemente em quatro abordagens operacionais, cujas limitações fundamentam a oportunidade de concepção de um novo artefato:
Edição Manual de Arquivos de Zona via Interface de Linha de Comando (CLI): Softwares tradicionais como o BIND 9 demandam a manipulação direta de arquivos planos em formato de arquivo mestre (master file format). Essa abordagem apresenta alta propensão a falhas tipográficas, ausência de validação estrutural imediata no momento da digitação, necessidade de recargas manuais de serviço e dependência excessiva da especialização técnica individual do operador.
Interfaces de Gestão Unificadas do Tipo Painel de Controle (ex.: cPanel, Plesk): Desenvolvidas primordialmente para o provisionamento de hospedagem web compartilhada, essas ferramentas realizam o isolamento por conta de cliente proprietária, mas não oferecem subdivisão granular de privilégios operacionais (como distinção estrita entre operadores de leitura e escrita) dentro de uma mesma equipe corporativa gerenciando zonas compartilhadas.
Painéis Web Dedicados de Código Aberto (ex.: PowerDNS-Admin): Embora amplamente consolidados no ecossistema de software livre e providos de recursos de controle de acesso e registro de eventos, tais sistemas sustentam-se sobre arquiteturas monolíticas tradicionais com renderização pelo lado do servidor (server-side rendering). Sob a ótica operacional, rotinas como a criação automática de registros PTR exigem ativações adicionais de configuração, o provisionamento de extensões de segurança (como DNSSEC) permanece desacoplado da criação da zona e não se evidenciam rotinas assistidas de triagem e remoção de registros reversos órfãos.
Dashboards Proprietários de Provedores em Nuvem (Cloud DNS): Plataformas gerenciadas (como AWS Route 53, Cloudflare e Google Cloud DNS) oferecem usabilidade moderna e alta resiliência. Contudo, introduzem aprisionamento tecnológico (vendor lock-in), custos financeiros contínuos atrelados ao volume de consultas e inviabilidade de adoção em redes corporativas ou governamentais restritas (on-premises) que exigem custódia estrita sobre os servidores autoritativos.
Observa-se, portanto, que a lacuna não reside na inexistência isolada dos mecanismos pretendidos — vários deles encontram-se implementados em ferramentas maduras —, mas na inexistência de uma combinação específica: uma camada de operação desacoplada do motor de resolução, aplicável a infraestruturas locais já padronizadas sobre o PowerDNS Authoritative Server, que reúna, em um mesmo fluxo, a gestão do ciclo de vida completo dos registros reversos (incluindo a remoção encadeada e a triagem de reversos órfãos), a normalização sintática por tipo de registro, o provisionamento de DNSSEC no ato da criação da zona, a distinção entre privilégios de leitura e de escrita no nível de cada zona e uma trilha de auditoria na própria camada de aplicação. O Capítulo 2 examina cada uma dessas ferramentas em detalhe, quantifica a sobreposição funcional em quadro comparativo e evidencia, a partir de um incidente de segurança recentemente corrigido no Technitium DNS Server, que a articulação entre automação de registros deriváveis e controle de acesso por zona constitui problema de projeto não trivial.
1.3 PROBLEMA DE PESQUISA E QUESTÕES NORTEADORAS
O conjunto de evidências reunido nas seções anteriores permite enunciar o problema de pesquisa nos seguintes termos: a operação de zonas DNS autoritativas depende de tarefas manuais, repetitivas e mutuamente dependentes cuja correção não é verificada no instante da edição, de modo que deslizes de execução e enganos de planejamento são publicados e propagados aos resolvedores antes de qualquer oportunidade de detecção, e permanecem ativos mesmo após a correção da origem. O problema não é, nesse enunciado, a existência do erro humano — que a literatura trata como inevitável (REASON, 1990; NORMAN, 2013) —, mas a ausência de um meio de operação que o torne estruturalmente impossível nas classes em que isso é tecnicamente viável.
Desse enunciado deriva a questão de pesquisa principal:
QP. Em que medida uma camada web de operação sobre um servidor autoritativo, que automatize registros deriváveis, restrinja privilégios ao nível da zona e registre integralmente as alterações efetuadas, reduz a incidência de erros de configuração e o esforço de operação em comparação com a edição manual assistida por linha de comando?
A questão principal é insuficiente, isoladamente, para orientar o projeto e a avaliação, pois agrega asserções de naturezas distintas — algumas verificáveis por argumentação técnica, outras dependentes de evidência empírica com operadores humanos. Desdobra-se, por isso, em quatro questões secundárias, cada uma vinculada a objetivos específicos (seção 1.4.2) e a episódios determinados do desenho de avaliação (seção 2.6.2):
QP1 — Caracterização. Quais classes de erro de configuração em zonas autoritativas são recorrentes, quais são seus impactos documentados e quais delas são passíveis de eliminação por construção — isto é, de serem tornadas impossíveis pelo próprio meio de edição — em oposição àquelas que só admitem detecção posterior?
QP2 — Mecanismos. Que mecanismos de projeto — automação de dados deriváveis, normalização sintática por tipo de registro, delegação do incremento do serial ao motor autoritativo e provisionamento de DNSSEC no fluxo de criação da zona — são necessários e suficientes para eliminar as classes identificadas em QP1, e que novos riscos esses mecanismos, ao automatizarem decisões antes tomadas pelo operador, introduzem na operação?
QP3 — Autoridade e automação. Que modelo de controle de acesso concilia a granularidade por zona com a automação de registros que residem em zona distinta daquela que o operador edita — caso do registro PTR, que é derivado de um registro A ou AAAA na zona direta, mas gravado na zona reversa, frequentemente sob responsabilidade administrativa diversa?
QP4 — Efeito observado. A operação mediada pela plataforma produz redução observável na taxa de erros de configuração e no tempo de execução de tarefas típicas de administração de zonas, e a solução é percebida como utilizável por operadores em contexto real de produção?
A distinção entre os dois níveis de asserção é deliberada e condiciona o desenho metodológico. As questões QP1 a QP3 são respondidas por revisão sistemática da literatura, por argumentação de projeto e por verificação técnica conduzida contra oráculos independentes do artefato; QP4, por sua vez, exige evidência obtida com operadores humanos em tarefas reais, razão pela qual a avaliação descrita na seção 2.6 combina episódios artificiais de laboratório com um estudo de caso naturalístico em ambiente de produção. Nenhuma das questões é respondida pela mera demonstração de que a plataforma funciona: a demonstração de funcionamento é condição necessária, porém não suficiente, para sustentar a asserção de que os problemas enunciados foram efetivamente mitigados.
Fontes Externas da Verdade e Validadores Independentes (ex.: NetBox DNS, Zonemaster): a primeira abordagem desloca a autoridade sobre os dados para um sistema de inventário, derivando registros diretos e reversos do endereçamento cadastrado; a segunda submete a delegação já publicada a uma suíte formal de casos de teste. Ambas contribuem para a qualidade das zonas, mas nenhuma atua no instante da edição: a fonte externa introduz um componente adicional de sincronização e a possibilidade de divergência entre o inventário e aquilo que efetivamente responde às consultas, ao passo que o validador opera a posteriori, quando o erro já foi publicado e já se propagou aos resolvedores.
Fontes Externas da Verdade e Validadores Independentes (ex.: NetBox DNS, Zonemaster): a primeira abordagem desloca a autoridade sobre os dados para um sistema de inventário, derivando registros diretos e reversos do endereçamento cadastrado; a segunda submete a delegação já publicada a uma suíte formal de casos de teste. Ambas contribuem para a qualidade das zonas, mas nenhuma atua no instante da edição: a fonte externa introduz um componente adicional de sincronização e a possibilidade de divergência entre o inventário e aquilo que efetivamente responde às consultas, ao passo que o validador opera a posteriori, quando o erro já foi publicado e já se propagou aos resolvedores.
Configura-se, assim, uma lacuna prática e metodológica: a demanda por uma camada de abstração web moderna, orientada a infraestruturas locais baseadas no PowerDNS Authoritative Server, que reúna arquitetura desacoplada, validação de sintaxe e formato em tempo real, automação nativa da integridade de zonas reversas, provisionamento simplificado de DNSSEC e distinção refinada de privilégios de leitura e escrita por domínio. O Capítulo 2 amplia essa análise, contemplando também as ferramentas de DNS como código.
1.4 OBJETIVOS DO TRABALHO
1.4.1 Objetivo Geral
Projetar, implementar e avaliar um artefato de software — denominado SanchezDNS — fundamentado no paradigma Design Science Research (DSR), que atue como uma camada de mediação sobre o PowerDNS Authoritative Server, capaz de reduzir a ocorrência de inconsistências operacionais por meio de automação e sincronização de registros reversos, validação estrutural de entrada, controle de acesso híbrido com permissões por zona e registro de eventos em trilha de auditoria centralizada.
1.4.2 Objetivos Específicos
Para responder às questões formuladas, estabelecem-se os objetivos específicos a seguir, indicando-se, ao final de cada um, a questão de pesquisa a que se vincula:
Mapear os padrões de erros operacionais e inconsistências mais frequentes na edição de zonas autoritativas a partir da literatura técnica e empírica de medição. (QP1)
Formular diretrizes de interface voltadas à mitigação de falhas humanas na digitação e à automação de dados derivados, integrando a sincronização automática de registros PTR (IPv4 e IPv6) na criação/edição de registros A/AAAA, ferramentas para detecção e remoção de registros reversos órfãos e delegação transparente do número serial ao motor do PowerDNS (SOA-EDIT-API). (QP2)
Automatizar o provisionamento inicial de segurança criptográfica em zonas gerenciadas, efetuando a geração da chave de assinatura de chaves (Key Signing Key – KSK) e a configuração dos parâmetros de negação autenticada de existência (NSEC3PARAM) no fluxo de criação do domínio. (QP2)
Implementar o artefato sob arquitetura desacoplada, dividindo o frontend reativo do backend intermediador e mantendo a API REST do PowerDNS Authoritative Server como fonte exclusiva da verdade para os dados de zona. (QP2)
Estruturar um modelo híbrido de controle de acesso que conjugue o perfil administrativo global à concessão isolada de privilégios de leitura e de escrita no nível individual de cada zona. (QP3)
Registrar em trilha de auditoria centralizada as operações de criação, alteração e remoção de zonas, registros e usuários, com identificação do autor e do instante de execução de cada alteração. (QP4)
Avaliar a eficácia da plataforma segundo o framework FEDS, mensurando a redução na taxa de erros sintáticos, a variação no tempo de execução de alterações operacionais e o índice de usabilidade percebida mediante a aplicação da escala SUS (System Usability Scale). (QP4)
1.5 DELIMITAÇÃO DO ESCOPO
A explicitação dos limites do trabalho é condição para a avaliação justa de seus resultados. Integram o escopo desta pesquisa a operação de zonas autoritativas diretas e reversas mantidas em instâncias do PowerDNS Authoritative Server acessíveis por sua API REST; a automação de registros deriváveis dessas zonas; o controle de acesso de operadores no nível da zona; e o registro de auditoria das operações realizadas pela plataforma.
Permanecem expressamente fora do escopo:
a resolução recursiva e o comportamento de resolvedores, tratados apenas como fonte de restrições de projeto — notadamente o efeito do armazenamento em cache sobre a reversibilidade das alterações;
a resiliência da infraestrutura a ataques de negação de serviço, o dimensionamento de nuvens anycast e o desempenho do motor de resolução, atributos que dependem do servidor autoritativo subjacente e não da camada de operação proposta;
a interação com registrários e com a zona pai, incluindo a publicação do registro DS e a manutenção da consistência entre os conjuntos de NS pai e filho, operações conduzidas por interfaces de terceiros; a plataforma limita-se a exibir ao operador os dados necessários e a sinalizar divergências detectadas;
a substituição de ferramentas de infraestrutura como código, cujo público e cujo modelo operacional são distintos, conforme discutido na seção 2.5.8;
a implementação de um motor DNS próprio, uma vez que o artefato pressupõe e preserva o servidor autoritativo como fonte única da verdade.
No plano da avaliação, delimitam-se igualmente a população de participantes do experimento controlado, restrita a profissionais e estudantes com experiência declarada em administração de redes, e o ambiente do estudo de caso, circunscrito a uma única organização. As consequências dessas delimitações sobre a validade externa dos resultados são discutidas na seção 2.6.4.
1.6 ESTRUTURA DA MONOGRAFIA
Esta monografia está organizada nos seguintes capítulos:
Capítulo 2 – Fundamentação Teórica e Trabalhos Correlatos: Revisa a arquitetura do protocolo DNS, os estudos empíricos sobre inconsistências de configuração em zonas autoritativas e a fundamentação epistemológica em Design Science Research. Estabelece, ainda, a estratégia de avaliação do artefato, com os episódios, as métricas e os critérios de êxito definidos previamente à construção.
Capítulo 3 – Metodologia e Modelagem do Artefato: Detalha as fases do método DSRM empregadas na condução da pesquisa, o levantamento de requisitos do SanchezDNS, a modelagem de autenticação e de controle de acesso, as diretrizes de normalização sintática, as estratégias de automação do ciclo de vida de registros e a modelagem da trilha de auditoria.
Capítulo 4 – Desenvolvimento da Plataforma SanchezDNS: Apresenta os aspectos de engenharia e implementação do protótipo, evidenciando as decisões de arquitetura de software, a integração com o PowerDNS e os algoritmos de tratamento de zonas reversas e auditoria.
Capítulo 5 – Avaliação e Resultados: Documenta a execução experimental sob as diretrizes do framework FEDS, expondo os dados coletados relativos à mitigação de erros sintáticos, tempo de execução de tarefas e os escores obtidos na escala SUS. Os dados são analisados à luz dos critérios de êxito fixados no Quadro 3.
Capítulo 6 – Conclusões: Sintetiza as contribuições práticas e teóricas alcançadas, discute as limitações atuais da ferramenta e do artefato, e propõe diretrizes para investigações futuras.

2 FUNDAMENTAÇÃO TEÓRICA E TRABALHOS CORRELATOS
Este capítulo reúne a base conceitual necessária à fundamentação e à avaliação do SanchezDNS. A seção 2.1 descreve a arquitetura do Sistema de Nomes de Domínio (DNS), enfatizando os componentes manipulados pelo artefato. A seção 2.2 discute o erro humano na operação de infraestrutura, examina os estudos empíricos sobre falhas de configuração em zonas autoritativas e consolida os achados em uma taxonomia de classes de erro, evidências e impactos. A seção 2.3 aborda os princípios de controle de acesso e auditoria. A seção 2.4 expõe o paradigma Design Science Research (DSR), que orienta metodologicamente esta pesquisa. A seção 2.5 analisa os trabalhos correlatos e delimita a contribuição da plataforma proposta. A seção 2.6 estabelece a estratégia de avaliação do artefato, especificando episódios, métricas, critérios de êxito e ameaças à validade. Por fim, a seção 2.7 sintetiza o capítulo.

2.1 ARQUITETURA DO SISTEMA DE NOMES DE DOMÍNIO
2.1.1 Espaço de nomes, zonas e delegação
O DNS organiza os nomes sob uma árvore hierárquica invertida cuja raiz é representada por um rótulo vazio. O nome de domínio plenamente qualificado (Fully Qualified Domain Name – FQDN) é definido pela sequência de rótulos de um nó até a raiz, separados por pontos (MOCKAPETRIS, 1987a). Na representação textual, o ponto final explicita a raiz; sua omissão faz com que o nome seja interpretado como relativo a uma origem, detalhe sintático que constitui fonte frequente de erros operacionais (BARR, 1996).
A administração do espaço de nomes é distribuída por meio de zonas — porções contíguas da árvore sob uma mesma autoridade administrativa. O repasse de autoridade sobre uma subárvore ocorre por delegação, materializada pela inserção de registros NS na zona pai apontando para os servidores da zona filha. Quando esses servidores pertencem ao próprio domínio delegado, a zona pai deve conter registros de endereço auxiliares (glue records) para evitar dependências circulares na resolução (MOCKAPETRIS, 1987a).
A infraestrutura do protocolo baseia-se em dois papéis funcionais fundamentais:
Servidores autoritativos: armazenam as definições originais das zonas e respondem às consultas com autoridade sobre os dados.
Resolvedores recursivos: consultam a hierarquia a partir da raiz em nome das aplicações clientes e retêm as respostas em cache pelo intervalo estipulado no Time to Live (TTL) de cada registro (MOCKAPETRIS, 1987a).
Devido ao cache nos resolvedores, qualquer erro publicado em zona autoritativa propaga-se e permanece ativo até a expiração do TTL, mesmo após ter sido corrigido no servidor de origem, a exemplo do incidente com o ccTLD sueco (.se) documentado por Pingdom (2009).

2.1.2 Registros de recurso
Os dados de uma zona são estruturados como registros de recurso (Resource Records – RR), compostos por nome, tipo, classe, TTL e dados específicos do tipo (RDATA) (MOCKAPETRIS, 1987b). O SanchezDNS oferece suporte a múltiplos tipos de registro:
A e AAAA: mapeiam nomes a endereços IPv4 e IPv6, respectivamente (MOCKAPETRIS, 1987b; THOMSON et al., 2003).
CNAME: estabelece um apelido canônico; por definição normativa, não pode coexistir com outros registros no mesmo nome (MOCKAPETRIS, 1987b; BARR, 1996).
MX: indica servidores de correio eletrônico, exigindo prioridade e FQDN terminado em ponto (MOCKAPETRIS, 1987b).
NS: especifica servidores de nomes autoritativos para a zona (MOCKAPETRIS, 1987b).
PTR: viabiliza a resolução reversa mapeando endereços IP a nomes (MOCKAPETRIS, 1987b).
TXT: armazena cadeias de texto genéricas, amplamente utilizado em mecanismos de autenticação como SPF e DKIM (MOCKAPETRIS, 1987b).
SRV: localiza serviços indicando prioridade, peso, porta e alvo canônico (GULBRANDSEN; VIXIE; ESIBOV, 2000).
CAA: restringe quais Autoridades Certificadoras podem emitir certificados TLS para o domínio (HALLAM-BAKER; STRADLING; HOFFMAN-ANDREWS, 2019).
HTTPS: publica parâmetros de conexão e negociação de protocolos para clientes web (SCHWARTZ; BISHOP; NYGREN, 2023).
TLSA: vincula certificados ou chaves públicas a serviços usando DANE (HOFFMAN; SCHLYTER, 2012).
ALIAS: pseudorregistro emulado pelo PowerDNS que provê funcionalidade análoga ao CNAME diretamente no ápice da zona, contornando a restrição de coexistência (POWERDNS, 2026).
A diversidade de formatos exigidos no RDATA eleva substancialmente a probabilidade de falhas na edição textual direta de arquivos de zona.
2.1.3 O registro SOA e a replicação de zonas
Toda zona autoritativa possui exatamente um registro Start of Authority (SOA) em seu ápice (MOCKAPETRIS, 1987b). Seus campos incluem o servidor primário, o e-mail do administrador e cinco parâmetros operacionais numéricos: SERIAL, REFRESH, RETRY, EXPIRE e MINIMUM.
O SERIAL opera como contador de versão de 32 bits. Na replicação primário-secundário tradicional, o servidor secundário avalia periodicamente o serial do primário e, ao detectar incremento, solicita a transferência de zona via AXFR ou IXFR. Essa comparação segue a aritmética de números de série em espaço modular circular (ELZ; BUSH, 1996).
A sincronia do sistema depende estritamente da atualização desse campo a cada alteração. Conforme adverte a RFC 1912: "Don't forget to change the serial number when you change data! If you don't, your secondaries will not transfer the new zone information" (BARR, 1996). A omissão na atualização do serial gera divergência de conteúdo entre servidores (zone drift), enquanto reduções indevidas impedem transferências subsequentes, exigindo intervenções manuais para recuperação da consistência.
2.1.4 Resolução reversa
A resolução reversa mapeia endereços IP a nomes canônicos utilizando registros PTR sob domínios especiais:
IPv4: os octetos numéricos são invertidos sob a raiz in-addr.arpa. O endereço 192.0.2.10, por exemplo, é consultado como 10.2.0.192.in-addr.arpa (MOCKAPETRIS, 1987b).
IPv6: cada um dos 32 nibbles hexadecimais é decomposto em rótulos isolados e dispostos em ordem inversa sob ip6.arpa (THOMSON et al., 2003).
Essa sintaxe fragmentada em 32 níveis para IPv6 torna o preenchimento manual propenso a erros de digitação. Além disso, zonas diretas e reversas são desacopladas administrativamente; nada no protocolo força sua correspondência. Contudo, essa paridade é cobrada por serviços de rede essenciais (como a mitigação de spam em servidores SMTP), motivando a prescrição expressa da RFC 1912 de assegurar a concordância estrita entre apontamentos A/AAAA e seus respectivos PTRs (BARR, 1996).
2.1.5 Extensões de segurança (DNSSEC)
As extensões de segurança do DNS (DNSSEC) incorporam integridade e autenticidade de origem às respostas por meio de assinaturas digitais, sem prover confidencialidade (ARENDS et al., 2005). O protocolo emprega os seguintes registros:
RRSIG: armazena a assinatura criptográfica de um determinado conjunto de registros (RRset).
DNSKEY: publica as chaves públicas da zona, convencionalmente divididas em chave de assinatura de chaves (Key Signing Key – KSK) e chave de assinatura de zona (Zone Signing Key – ZSK) (KOLKMAN; MEKKING; GIEBEN, 2012).
DS: publicado na zona pai, guarda o resumo criptográfico da KSK da zona filha, constituindo o elo da cadeia de confiança hierárquica.
Para provar a inexistência de nós sem permitir a enumeração de zonas, o NSEC3 utiliza resumos criptográficos dos nomes (LAURIE et al., 2008). Seus parâmetros operacionais residem no registro NSEC3PARAM. Conforme as diretrizes da RFC 9276, recomenda-se mitigar o impacto computacional nos resolvedores fixando a contagem de iterações em zero e dispensando o uso de salt, adotando a configuração expressa NSEC3PARAM 1 0 0 - (HARDAKER; DUKHOVNI, 2022). O provisionamento manual de cada uma dessas etapas adiciona múltiplos pontos de falha à operação.
2.1.6 O PowerDNS Authoritative Server
O PowerDNS Authoritative Server é um servidor autoritativo modular de código aberto que desacopla o motor de processamento DNS do armazenamento físico, utilizando backends como bancos de dados relacionais. A aplicação expõe uma API REST para gestão programática de zonas, registros e chaves criptográficas (POWERDNS, 2026), atuando como base de integração do SanchezDNS.
Destacam-se dois mecanismos centrais do PowerDNS para esta pesquisa:
Metadado SOA-EDIT-API: quando definido com o parâmetro DEFAULT, delega ao servidor a atualização automática do serial SOA a cada alteração via API REST, gerando valores no padrão AAAAMMDD01 ou incrementando o serial existente em uma unidade, eliminando o erro de zone drift por esquecimento do operador.
Modos de operação de zona: no modo Native, a replicação dos dados é delegada diretamente à camada de armazenamento do backend (como replicação de banco de dados), dispensando transferências clássicas de zona; nos modos Primary e Secondary, o PowerDNS utiliza notificações NOTIFY e transferências AXFR/IXFR (POWERDNS, 2026).
2.2 ERROS DE CONFIGURAÇÃO EM ZONAS AUTORITATIVAS
2.2.1 O erro humano na operação de infraestrutura
Reason (1990) categoriza as falhas humanas em dois grupos: falhas de execução, nas quais a intenção está correta mas a ação física diverge (deslizes e lapsos), e enganos (mistakes), quando a falha reside no próprio planejamento ou na compreensão da situação. Norman (2013) estende esse modelo à interação humano-computador, sustentando que deslizes e lapsos ocorrem com frequência mesmo entre operadores experientes durante rotinas repetitivas. A prevenção eficaz de tais incidentes demanda interfaces com restrições semânticas, funções de bloqueio (forcing functions) e automação de passos suscetíveis ao esquecimento.
No DNS, esquecer o ponto final de um FQDN ou não atualizar o serial SOA configuram deslizes típicos de execução. Por outro lado, delegar uma zona a um servidor sem prévia configuração constitui um engano de planejamento. Oppenheimer, Ganapathi e Patterson (2003) comprovaram empiricamente essa vulnerabilidade em serviços corporativos de Internet, revelando que o erro de operador foi o fator primário de indisponibilidade em dois dos três serviços avaliados, com erros de configuração correspondendo à principal subcategoria.
A RFC 1912 (BARR, 1996) compila sistematicamente essas falhas operacionais no DNS: inconsistências A/PTR, omissão de incremento no serial SOA, ausência do ponto terminal em FQDNs, coexistência indevida com CNAME e delegações defeituosas (lame delegations).
2.2.2 Estudos empíricos de medição
Medições na Internet aberta comprovam que essas falhas mantêm prevalência persistente:
Pappas et al. (2004): identificaram delegações defeituosas em 15% das zonas avaliadas e dependências cíclicas em cerca de 2%. Mostraram que configurações incorretas degradam a disponibilidade e aumentam o tempo de resolução em até dez vezes, refutando a premissa do DNS clássico de que falhas de zonas operam de forma isolada e independente.
Bojović e Gajin (2017): analisaram 79.933 domínios (redes acadêmicas, ccTLD .rs e top 1.000 da Internet). Constataram que quase 5% dos domínios acadêmicos apontavam para servidores não autoritativos na zona pai e que entre 39,7% e mais de 45% dos domínios não possuíam registro PTR para seus servidores de nomes.
Akiwate et al. (2020): em estudo longitudinal de nove anos sobre 499 milhões de domínios, verificaram delegações defeituosas em aproximadamente 14% das amostras ativas. O trabalho revelou que tais problemas permanecem mascarados por longos períodos devido à redundância parcial e expõem centenas de milhares de domínios a sequestro malicioso (hijacking), indicando que a prevalência do problema permanece em patamar semelhante ao observado por Pappas et al. (2004).
Sommese et al. (2020): compararam os conjuntos de registros NS declarados pela zona pai e pela zona filha na totalidade dos domínios de segundo nível de .com, .net e .org, identificando inconsistência em mais de treze milhões de domínios, dos quais ao menos cinquenta mil figuram entre o milhão mais acessado. Por meio de experimentos com sondas RIPE Atlas, demonstraram que o comportamento dos resolvedores diante da divergência não é uniforme, o que converte a inconsistência em indisponibilidade intermitente e de diagnóstico difícil.
Liu, Hao e Wang (2016): caracterizaram os registros pendentes (dangling records) — apontamentos que sobrevivem ao desprovisionamento do recurso referenciado — e localizaram 467 ocorrências exploráveis em 277 domínios entre os dez mil mais acessados e em 52 zonas do domínio .edu. Descreveram três vetores pelos quais um adversário assume o controle do subdomínio afetado, incluindo a obtenção de certificados TLS legítimos em nome da organização titular, o que demonstra que a omissão na remoção de um registro é falha de segurança, e não apenas de higiene operacional.
Chung et al. (2017): acompanharam, por vinte e um meses, instantâneos diários dos registros DNSSEC de todos os domínios assinados em .com, .net e .org. Constataram que mais de 30% das zonas assinadas não possuíam o registro DS correspondente na zona pai, situação em que a assinatura não produz efeito de validação, que aproximadamente 70% das zonas nunca efetuaram rotação de chaves e que 91,7% das chaves de assinatura de zona empregavam parâmetros criptográficos considerados fracos. Os autores atribuem o quadro predominantemente a deficiências no fluxo operacional oferecido por registrários e por ferramentas de gestão, e não a limitações do protocolo.
Korczyński et al. (2016): mediram a prevalência de servidores autoritativos que aceitam atualizações dinâmicas sem autenticação, identificando 1.877 domínios vulneráveis (0,065%) em amostra aleatória de 2,9 milhões e 587 domínios (0,062%) entre o milhão mais acessado. A falha, que os autores denominam envenenamento de zona (zone poisoning), permite a qualquer agente da Internet inserir, alterar ou remover registros na zona afetada, anulando integralmente o controle de acesso pretendido pelo operador.
Tomados em conjunto, os estudos revelam três regularidades relevantes para o projeto do artefato. A primeira é a estabilidade temporal da prevalência: as taxas de delegação defeituosa observadas por Pappas et al. (2004) e por Akiwate et al. (2020) são equivalentes, ainda que separadas por dezesseis anos e por uma ordem de grandeza na população estudada, o que sugere que a evolução das ferramentas de administração não alterou materialmente a incidência do problema. A segunda é o caráter silencioso das falhas: em praticamente todos os casos examinados, a configuração incorreta continua produzindo respostas aparentemente válidas, de modo que a detecção depende de verificação ativa e deliberada, e não da observação do funcionamento cotidiano. A terceira é a concentração das falhas em operações de manutenção do ciclo de vida — alteração, migração e remoção de registros — e não na criação inicial da zona, o que desloca o foco do projeto da tarefa de provisionamento para a de atualização consistente de dados interdependentes.
2.2.3 Impactos operacionais e de segurança
A caracterização dos impactos exige distinguir três mecanismos pelos quais uma falha de configuração se converte em dano. O primeiro é a amplificação temporal produzida pelo armazenamento em cache. Como a resposta obtida por um resolvedor permanece válida pelo intervalo declarado no TTL, um erro publicado continua sendo servido a clientes mesmo após a correção no servidor de origem; o intervalo de exposição não é, portanto, o tempo de permanência do erro na zona, mas a soma desse tempo com o TTL do registro afetado. No incidente do ccTLD sueco, a correção da zona ocorreu em cerca de uma hora, mas os efeitos persistiram por período significativamente superior em virtude desse mecanismo (PINGDOM, 2009).
O segundo mecanismo é a degradação da redundância planejada. Pappas et al. (2004) demonstraram que delegações defeituosas não produzem indisponibilidade imediata — o domínio continua respondendo por meio dos servidores restantes —, mas consomem silenciosamente a margem de tolerância a falhas que a operação supunha possuir, elevando o tempo de resolução em até uma ordem de grandeza quando o resolvedor consulta um servidor não autoritativo antes de obter resposta válida. Trata-se, nesse sentido, de uma falha latente: seu custo só se materializa quando ocorre um segundo evento, momento em que a redundância já não existe.
O terceiro mecanismo é a conversão da falha operacional em vulnerabilidade de segurança. Akiwate et al. (2020) demonstraram que delegações defeituosas expõem centenas de milhares de domínios ao sequestro por terceiros capazes de registrar o nome do servidor de nomes inexistente para o qual a delegação aponta; Liu, Hao e Wang (2016) documentaram trajetória análoga a partir de registros pendentes; e Korczyński et al. (2016) evidenciaram a inserção arbitrária de registros em zonas cujos servidores aceitam atualizações não autenticadas. Em todos os casos, o adversário não explora uma falha do protocolo, mas um estado inconsistente deixado por uma operação incompleta.
Os incidentes documentados em larga escala confirmam a materialidade desses mecanismos. Além do episódio do ccTLD sueco, registram-se a interrupção do serviço Edge DNS da Akamai em 22 de julho de 2021, decorrente de atualização de configuração posteriormente revertida (AKAMAI, 2021), e a indisponibilidade dos serviços da Meta em 4 de outubro do mesmo ano, por mais de sete horas, quando um comando de manutenção provocou a retirada dos anúncios de rota dos prefixos que abrigavam seus servidores autoritativos (META, 2021). O relato oficial deste último episódio é particularmente instrutivo para o problema aqui tratado: a organização dispunha de um mecanismo automatizado de auditoria destinado precisamente a impedir comandos dessa natureza, e o incidente ocorreu porque uma falha nesse mecanismo permitiu que o comando fosse aceito (META, 2021). O caso ilustra que a mera existência de validação é insuficiente quando a validação não é parte indissociável do caminho de execução da alteração.
2.2.4 Síntese: taxonomia de erros, evidências e mecanismos de mitigação
A revisão precedente permite organizar os achados em uma taxonomia operacional, apresentada no Quadro 1. A organização proposta não é meramente descritiva: cada classe de erro é associada à evidência empírica que atesta sua prevalência, ao impacto documentado, ao mecanismo de mitigação previsto no artefato e à questão de pesquisa correspondente. O quadro cumpre, assim, três funções ao longo do trabalho. Primeiro, responde à questão QP1, ao distinguir as classes elimináveis por construção daquelas que apenas admitem detecção posterior — distinção determinante para o escopo do artefato, uma vez que classes dependentes de ação junto à zona pai ou ao registrário não podem ser resolvidas por uma camada de operação sobre o servidor autoritativo. Segundo, orienta a derivação dos requisitos apresentados no Capítulo 3, assegurando que cada mecanismo implementado corresponda a um problema documentado, e não a uma suposição do projetista. Terceiro, fornece a taxonomia de classificação de erros empregada na avaliação descrita na seção 2.6, o que permite que a redução de erros seja medida por categoria, e não apenas em agregado.
Quadro 1 – Taxonomia de erros de configuração em zonas autoritativas, evidências, impactos e mecanismos de mitigação
Classe de erro
Evidência empírica
Impacto principal
Tratamento no SanchezDNS
QP
Divergência entre A/AAAA e PTR
BARR (1996); BOJOVIĆ e GAJIN (2017): 39,7% a 45,3% dos domínios sem PTR para seus servidores de nomes
Rejeição de mensagens por servidores SMTP; diagnóstico de rede prejudicado
Criação e atualização derivadas do registro A/AAAA; remoção com guarda de referência
QP2
QP3
Serial SOA não incrementado
BARR (1996); ELZ e BUSH (1996)
Dessincronização (zone drift): secundários servem versão desatualizada
Delegação do incremento ao motor autoritativo pelo metadado SOA-EDIT-API
QP2
Erro sintático no RDATA (ponto final ausente, prioridade de MX, FQDN relativo)
BARR (1996); PINGDOM (2009): cerca de 900 mil domínios .se afetados
Corrupção da zona publicada, propagada e retida em cache
Normalização e validação por tipo de registro antes da submissão à API
QP2
Delegação defeituosa (lame delegation)
PAPPAS et al. (2004): 15%; AKIWATE et al. (2020): cerca de 14%
Perda de redundância; latência até dez vezes maior; risco de sequestro
Não tratada pelo artefato (ver seção 1.5)
QP1
Inconsistência entre NS do pai e do filho
SOMMESE et al. (2020): mais de 13 milhões de domínios
Indisponibilidade intermitente dependente do resolvedor
Fora do escopo de escrita (depende do registrário); detecção e sinalização
QP1
Registros pendentes e reversos órfãos
LIU, HAO e WANG (2016): 467 registros exploráveis em 277 domínios populares
Sequestro de subdomínio; emissão indevida de certificados TLS
Remoção encadeada do reverso e rotina de triagem de registros órfãos
QP2
DNSSEC incompleto (DS ausente; ausência de rotação)
CHUNG et al. (2017): mais de 30% sem DS; cerca de 70% sem rotação de chaves
Assinatura sem efeito de validação; falhas de resolução em rollover malconduzido
Geração de KSK e de NSEC3PARAM conforme a RFC 9276 no fluxo de criação; exibição do DS para publicação no pai
QP2
Privilégio excessivo do operador
SALTZER e SCHROEDER (1975); OPPENHEIMER, GANAPATHI e PATTERSON (2003)
O raio de impacto de um erro estende-se a zonas alheias à responsabilidade do operador
Papéis globais combinados a permissões de leitura e de escrita atribuídas por zona
QP3
Atualização dinâmica não autenticada
KORCZYŃSKI et al. (2016): 0,065% dos domínios de amostra aleatória
Inserção, alteração ou remoção arbitrária de registros por terceiros
API autenticada como único caminho de escrita; ausência de canal de atualização anônima
QP3
Ausência de rastreabilidade das alterações
KENT e SOUPPAYA (2006)
Impossibilidade de análise forense e de atribuição de responsabilidade
Trilha de auditoria centralizada com ação, instante, autor e recurso afetado
QP4

Fonte: elaborado, a partir das referências indicadas na coluna de evidências.
Das dez classes relacionadas, seis são consideradas elimináveis por construção no âmbito de uma camada de operação sobre o servidor autoritativo — divergência A/PTR, serial não incrementado, erro sintático no RDATA, registros reversos órfãos, DNSSEC incompleto no ato da criação e privilégio excessivo —, no sentido de que a interface pode tornar impossível a produção do estado incorreto, e não apenas detectá-lo. Duas classes — delegação defeituosa e inconsistência entre pai e filho — dependem de ação em sistemas de terceiros e admitem, no máximo, detecção e sinalização. Uma classe — atualização dinâmica não autenticada — é eliminada pela própria arquitetura, ao inexistir canal de escrita não autenticado. A classe restante, ausência de rastreabilidade, não é um erro de configuração em sentido estrito, mas uma condição que impede o diagnóstico das demais, e por isso integra a taxonomia. Essa classificação constitui a base sobre a qual se define, na seção 2.6, o que significará afirmar que o artefato atendeu aos objetivos propostos.
2.3 CONTROLE DE ACESSO E AUDITORIA
A integridade operacional de sistemas multiusuário apoia-se no princípio do menor privilégio (Least Privilege), formalizado por Saltzer e Schroeder (1975). O postulado estabelece que todo sujeito deve operar munido apenas do conjunto estritamente necessário de permissões para executar sua função, mitigando o raio de impacto de acidentes e condutas inadequadas. No contexto de zonas DNS, isso implica limitar a atuação dos operadores exclusivamente aos domínios pelos quais respondem formalmente.
O modelo de Controle de Acesso Baseado em Papéis (Role-Based Access Control – RBAC), formalizado por Sandhu et al. (1996), reduz a sobrecarga de gestão desacoplando permissões de identidades individuais e atrelando-as a perfis funcionais (roles). Em cenários de gerenciamento de recursos pontuais, adotam-se abordagens híbridas nas quais papéis gerais coexistem com concessões atreladas a entidades específicas — estrutura adotada no SanchezDNS, permitindo atribuir direitos segregados de leitura e de escrita no escopo individual de cada zona.
A auditoria atua de forma complementar ao controle de acesso. Conforme o guia NIST SP 800-92 (KENT; SOUPPAYA, 2006), o registro contínuo de eventos de segurança viabiliza a análise forense de incidentes, o diagnóstico operacional e a conformidade regulatória. Para tanto, cada entrada de log deve registrar minimamente a ação executada, o carimbo de data e hora (timestamp), o identificador do agente e o recurso modificado.
2.4 DESIGN SCIENCE RESEARCH
2.4.1 Fundamentos metodológicos
A pesquisa fundamenta-se nas ciências do artificial estabelecidas por Simon (1996), cujo foco não reside em descrever fenômenos naturais pré-existentes, mas em prescrever soluções para alcançar metas predeterminadas. March e Smith (1995) definem quatro saídas materiais decorrentes da pesquisa em tecnologia da informação — construtos, modelos, métodos e instanciações —, estruturadas em dois ciclos principais: construir e avaliar.
Hevner et al. (2004) estruturaram o paradigma em Sistemas de Informação sob sete diretrizes norteadoras: artefato como objeto de estudo, relevância ao negócio, avaliação rigorosa, contribuições claras, rigor nos métodos, design como processo de busca heurística e comunicação ampla dos resultados.
Quanto à natureza da contribuição científica, Gregor e Hevner (2013) estabelecem uma matriz bidimensional correlacionando a maturidade do problema com a da solução (Invenção, Melhoria, Exaptação e Design Rotineiro) e três níveis de abstração dos artefatos: instanciações situadas (nível 1), teorias de design emergentes (nível 2) e teorias de design consolidadas (nível 3). Como os erros de operação em DNS autoritativo constituem uma dor conhecida na literatura e a plataforma propõe uma arquitetura integrada de mitigação para enfrentá-los, o projeto classifica-se formalmente no quadrante de Melhoria (Improvement), com contribuições situadas nos níveis 1 (o software como instanciação funcional) e 2 (princípios de design preventivo).
2.4.2 Processo DSRM
O ciclo de desenvolvimento do artefato segue a Design Science Research Methodology (DSRM), formulada por Peffers et al. (2007). O modelo estrutura o avanço da investigação em seis etapas sequenciais e iterativas:
Identificação do problema e motivação: delimitação do impacto dos erros manuais de DNS na disponibilidade da Internet.
Definição dos objetivos da solução: estabelecimento de requisitos como automação de PTR, validação estrutural e isolamento de permissões.
Projeto e desenvolvimento: concepção da arquitetura técnica e implementação dos módulos da plataforma.
Demonstração: emprego funcional do artefato em cenários controlados de administração de zonas.
Avaliação: mensuração da conformidade técnica e aplicação de métricas de usabilidade com usuários.
Comunicação: formalização textual dos achados e publicação acadêmica.
2.4.3 Avaliação: Framework FEDS e Escala SUS
A verificação do artefato segue o Framework for Evaluation in Design Science (FEDS), proposto por Venable, Pries-Heje e Baskerville (2016). O arcabouço orienta a concepção dos episódios de teste combinando duas variáveis:
Propósito funcional: formativo (orientado a realimentar o design e refinar funcionalidades) ou somativo (voltado a aferir o cumprimento global dos objetivos propostos).
Ambiente de teste: artificial (ensaios sintéticos em bancada de laboratório) ou naturalístico (avaliação de campo com operadores em contexto operacional autêntico).
Para quantificar a percepção de usabilidade humana na interface do artefato, adota-se a System Usability Scale (SUS) (BROOKE, 1996). Composta por dez afirmações avaliadas em escala Likert de 5 pontos, a SUS produz uma pontuação consolidada de 0 a 100. As métricas interpretativas estabelecidas por Bangor, Kortum e Miller (2008) servem como régua de corte para situar o nível de usabilidade da ferramenta frente aos padrões da indústria.
2.5 TRABALHOS CORRELATOS
Esta seção traça um paralelo comparativo entre as estratégias predominantes para gestão de zonas autoritativas, fundamentando-se nas documentações oficiais e nos históricos de versão consultados em setembro de 2026. A ordem de exposição parte das abordagens de menor mediação — a edição direta de arquivos — em direção às de maior automação, encerrando-se na validação independente e em um quadro comparativo que isola a lacuna endereçada pelo artefato.
2.5.1 BIND 9 e manipulação direta de arquivos
O BIND 9, desenvolvido pelo Internet Systems Consortium (ISC), representa a implementação referencial do DNS na Internet. Em configurações convencionais, seus registros residem em arquivos de texto sob o formato padrão master file (MOCKAPETRIS, 1987b). A edição é manual, sem interface gráfica de usuário nativa. O controle de acesso e os registros de alteração restringem-se ao sistema de arquivos hospedeiro ou a repositórios externos.A responsabilidade de manter registros A e PTR alinhados, bem como incrementar o serial SOA a cada edição, recai exclusivamente sobre a disciplina do operador.
2.5.2 Interfaces web de software livre (PowerDNS-Admin)
O ecossistema em torno do PowerDNS conta com ferramentas de código aberto consolidadas como o PowerDNS-Admin, que, desenvolvido em Python/Flask, provê painel web com suporte a papéis de usuário, permissões por zona e autenticação centralizada (LDAP, OAuth). Todavia, a sincronização de registros PTR atua como opção secundária desligada por padrão (auto_ptr), a geração de DNSSEC é tratada fora do fluxo elementar de abertura da zona e não se identificou na documentação rotina de detecção de reversos órfãos (POWERDNS-ADMIN, 2026).
2.5.3 Technitium DNS Server
O Technitium DNS Server é um servidor de código aberto que acumula as funções autoritativa e recursiva e integra, no mesmo produto, o motor de resolução e o painel web de administração. Seu histórico de versões documenta a incorporação progressiva de recursos que se sobrepõem de forma substancial aos pretendidos por esta pesquisa: a versão 5.4, de 18 de outubro de 2020, introduziu a opção de criação da zona reversa no momento da inclusão de registros A ou AAAA; a versão 8.0, de 26 de março de 2022, acrescentou a assinatura DNSSEC de zonas primárias com algoritmos RSA e ECDSA diretamente pela interface; e a versão 9.0, de 24 de setembro de 2022, implementou controle de acesso baseado em papéis com permissões atribuíveis no nível da zona a usuários e a grupos, além de credenciais de API não expiráveis para automação (TECHNITIUM, 2026). Trata-se, portanto, do trabalho correlato funcionalmente mais próximo do artefato aqui proposto, e sua omissão comprometeria a honestidade da análise comparativa.
Duas ordens de consideração, contudo, sustentam a pertinência do SanchezDNS. A primeira é arquitetural. O Technitium é um servidor completo, e sua camada administrativa é indissociável do motor de resolução: organizações que já operam sobre o PowerDNS Authoritative Server — escolha frequentemente motivada pelo desacoplamento entre o motor e o backend de armazenamento, que permite replicação pela própria camada de banco de dados — não podem adotá-lo sem substituir a infraestrutura autoritativa, decisão cujo custo excede em muito o da adoção de uma camada de operação. A proposta desta pesquisa situa-se, por isso, em ponto distinto do espaço de projeto: não o de um servidor com painel, mas o de um painel para um servidor já implantado.
A segunda ordem de consideração é substantiva e reforça a relevância da questão QP3. A versão 15.5 do Technitium, de 19 de setembro de 2026, corrigiu falha de autorização, reportada por Tao Pan, que permitia utilizar a opção ptr das chamadas de inclusão, atualização e exclusão de registros para criar, sobrescrever ou remover registros PTR em zonas reversas sobre as quais o usuário autenticado não detinha permissão de escrita (TECHNITIUM, 2026). O episódio é significativo por três razões: ocorreu quatro anos após a introdução das permissões por zona, em produto maduro e amplamente implantado; resultou precisamente da interação entre os dois mecanismos que o presente trabalho pretende combinar; e evidencia que a automação de registros deriváveis atravessa, por natureza, fronteiras de autorização, uma vez que o registro derivado reside em zona distinta daquela sobre a qual o operador age. Longe de enfraquecer a proposta, o incidente confirma que a articulação entre automação e controle de acesso granular constitui problema de projeto não trivial, que demanda tratamento explícito — e não meramente a justaposição de duas funcionalidades independentes. O modelo adotado para resolvê-lo é descrito no Capítulo 3 e verificado pelos testes de autorização negativos especificados na seção 2.6.2.
2.5.4 Microsoft DNS Server e o console DNS Manager
O servidor DNS do Windows Server, administrado pelo console DNS Manager, constitui a solução predominante em redes corporativas estruturadas sobre o Active Directory. Oferece, há mais de duas décadas, a opção de atualização do registro de ponteiro associado (Update associated pointer (PTR) record) na criação e na edição de registros de host, automatizando a manutenção do reverso correspondente, e permite delegar permissões no nível da zona por meio das listas de controle de acesso do diretório, acessíveis pela guia de segurança das propriedades da zona (MICROSOFT, 2026). Dispõe ainda do mecanismo de envelhecimento e remoção automática (aging and scavenging), que elimina registros dinamicamente registrados cujo carimbo temporal não é renovado — aproximação parcial ao problema dos registros órfãos, limitada, porém, aos registros de origem dinâmica, sem alcançar os criados manualmente pelo operador.
Três restrições delimitam sua aplicabilidade ao problema tratado. A delegação granular de permissões exige zonas integradas ao Active Directory Domain Services, o que vincula a solução ao ecossistema Windows e ao modelo de diretório corporativo, e a torna inaplicável a infraestruturas baseadas em servidores autoritativos de código aberto. A automação do registro PTR pela interface gráfica está documentada como inoperante quando a zona reversa é subdividida em sub-redes menores que um octeto (MICROSOFT, 2026), cenário comum em organizações que recebem blocos de endereçamento inferiores a /24 e que constitui, justamente, o caso em que o preenchimento manual é mais suscetível a erro. Por fim, a trilha de auditoria não é recurso do console de DNS, mas dos subsistemas de auditoria do sistema operacional e do diretório, o que dispersa a evidência entre componentes distintos e dificulta a reconstituição do histórico de uma zona específica.
2.5.5 NetBox DNS e a abordagem de fonte externa da verdade
O plugin NetBox DNS estende o NetBox — sistema de gestão de infraestrutura e de endereçamento — para modelar servidores de nomes, zonas, visões e registros. Entre seus recursos figuram a criação e a atualização automáticas de registros PTR a partir de registros de endereço, a geração automática do serial da zona e dos registros SOA e NS, e a validação de nomes e valores conforme as RFCs pertinentes; o mecanismo IPAM DNSsync permite derivar registros diretos e reversos diretamente dos objetos de endereço IP associados a prefixos vinculados a visões de DNS (NETBOX DNS, 2026). O plugin herda ainda, do NetBox, um registro de alterações por objeto e um modelo de permissões baseado em objetos e em conjuntos de restrições.
A abordagem é conceitualmente distinta da adotada nesta pesquisa. No modelo do NetBox DNS, o inventário torna-se a fonte da verdade e o servidor autoritativo passa a ser um destino de sincronização, o que exige um componente adicional de publicação e admite, por construção, a divergência entre o que está cadastrado e o que efetivamente responde às consultas — divergência que constitui, ela própria, uma nova classe de erro operacional, ausente das arquiteturas em que o servidor é a autoridade. O SanchezDNS adota orientação inversa: mantém a API do PowerDNS como fonte única da verdade e opera sobre o estado real do servidor, dispensando reconciliação. Ambas as escolhas são legítimas e correspondem a contextos organizacionais diferentes; explicitá-las evita que a comparação entre as ferramentas seja conduzida como se disputassem o mesmo espaço de projeto.
2.5.6 Knot DNS e as interfaces de controle programáticas
O Knot DNS, desenvolvido pelo CZ.NIC e amplamente empregado na operação de domínios de topo, é um servidor autoritativo cuja administração se dá pelo utilitário knotc sobre um soquete de controle local, complementado por vinculações de biblioteca que permitem automação programática. O modelo de edição de zonas é transacional: as alterações são iniciadas com zone-begin, aplicadas com zone-set e zone-unset, inspecionadas antes da efetivação com zone-diff e confirmadas com zone-commit ou descartadas com zone-abort; enquanto a transação permanece aberta, todos os eventos que modificam a zona — inclusive a assinatura DNSSEC automática — ficam bloqueados (CZ.NIC, 2026). O servidor automatiza igualmente a política de serial e a gestão do ciclo de vida das chaves criptográficas, incluindo as trocas programadas.
Duas características o distanciam do problema tratado nesta pesquisa. A interface é de linha de comando e de biblioteca, sem painel web, permanecendo inacessível a operadores sem familiaridade com o ambiente de terminal — precisamente o perfil para o qual a literatura de erro humano recomenda restrições semânticas na interface (NORMAN, 2013). E o controle de acesso é exercido pelas permissões do soquete e do sistema operacional, inexistindo noção de usuário de aplicação, de papel ou de permissão por zona; a segregação de responsabilidades, quando necessária, é implementada fora do servidor. A semântica transacional, por outro lado, constitui referência de projeto diretamente aplicável ao artefato: é a garantia de atomicidade oferecida pelo Knot que o SanchezDNS precisa emular ao propagar uma alteração que afeta simultaneamente a zona direta e a reversa, situação em que a efetivação parcial produziria exatamente a divergência que o mecanismo pretende evitar. Essa exigência é incorporada aos requisitos apresentados no Capítulo 3.
2.5.7 Zonemaster e a validação independente de delegações
O Zonemaster é uma ferramenta de validação da qualidade de delegações DNS desenvolvida conjuntamente pela Internetstiftelsen, operadora do ccTLD sueco, e pela AFNIC, operadora do ccTLD francês, em substituição às ferramentas DNSCheck e Zonecheck anteriormente mantidas por cada uma dessas organizações. Sua característica mais relevante para esta pesquisa é a existência de um plano de testes formalmente especificado segundo a estrutura da norma IEEE 829-2008, composto por setenta e quatro casos de teste distribuídos em nove categorias — Address, Basic, Connectivity, Consistency, Delegation, DNSSEC, Nameserver, Syntax e Zone —, cada um com critérios de aprovação, mensagens e níveis de severidade definidos, e executável por interface web, por linha de comando ou por serviço de retaguarda com saída em JSON (ZONEMASTER, 2026).
O Zonemaster não é, contudo, uma ferramenta de operação: atua sobre a delegação já publicada, não previne o erro, não gerencia permissões e não registra autoria das alterações. Representa, na taxonomia de abordagens desta seção, a estratégia de verificação independente e posterior — complementar, e não alternativa, à prevenção no momento da edição. Nesta pesquisa, o Zonemaster ocupa dupla condição: é trabalho correlato, por materializar essa estratégia, e é instrumento de medição na avaliação do artefato, conforme especificado na seção 2.6.2. Sua adequação a esse segundo papel decorre de três propriedades: a especificação formal dos casos de teste, que torna a medição reprodutível e auditável por terceiros; a saída estruturada, que permite a coleta automatizada; e, sobretudo, sua independência em relação ao objeto avaliado, que evita a circularidade de aferir o artefato pelos mesmos critérios que ele próprio implementa.
2.5.8 DNS como código (IaC) e plataformas em nuvem
DNS como código: ferramentas como octoDNS e DNSControl empregam o paradigma de Infraestrutura como Código (IaC). As zonas são modeladas como arquivos de configuração versionados (Git), incorporando validações pré-publicação (dry-run/preview) e auditoria herdada do histórico de commits (OCTODNS, 2026; DNSCONTROL, 2026). Essa abordagem, contudo, carece de interface gráfica para operadores não técnicos e exige políticas de controle de acesso baseadas em repositórios, e não na granularidade fina por zona.
Provedores em nuvem: serviços corporativos como Cloudflare DNS, AWS Route 53 e Google Cloud DNS oferecem alta disponibilidade, suporte nativo a DNSSEC e interfaces amigáveis com controle de acesso integrado. No entanto, envolvem custos recorrentes, risco de dependência de fornecedor (vendor lock-in) e incompatibilidade com ambientes corporativos que demandam custódia estrita de dados em infraestruturas locais (on-premises).
2.5.9 Quadro comparativo e delimitação da contribuição
O Quadro 2 sintetiza a análise precedente sobre cinco dimensões derivadas diretamente das classes de erro elimináveis por construção identificadas no Quadro 1. A classificação utiliza três valores: "Sim", quando o recurso integra o fluxo principal da ferramenta e opera sem configuração adicional; "Parcial", quando existe mas depende de ativação explícita, cobre apenas parte dos casos ou é oferecido por componente externo à ferramenta; e "Não", quando inexiste.
Quadro 2 – Comparação entre ferramentas de gestão de zonas autoritativas
Ferramenta ou abordagem
PTR automático na criação
Remoção encadeada e triagem de órfãos
DNSSEC no fluxo de criação da zona
Permissão por zona (leitura/escrita)
Trilha de auditoria na aplicação
BIND 9 (arquivos de zona)
Não
Não
Parcial (política de assinatura, sem interface)
Parcial (ACL de atualização dinâmica)
Não (sistema de arquivos)
PowerDNS-Admin
Parcial (opcional, desativado por padrão)
Não
Não (fora do fluxo de criação)
Parcial (acesso à zona, sem separar leitura de escrita)
Sim
Poweradmin
Parcial
Não
Não
Parcial
Sim (estado anterior e posterior)
Technitium DNS Server
Sim (desde a v5.4, 2020)
Não
Parcial (assinatura pela interface, v8.0, 2022)
Sim (desde a v9.0, 2022)
Parcial
Microsoft DNS Manager
Sim (exceto zonas reversas classless)
Parcial (scavenging de registros dinâmicos)
Sim (assistente de assinatura)
Sim (ACL do Active Directory)
Parcial (auditoria do SO e do diretório)
NetBox DNS
Sim (derivado do IPAM)
Sim (ciclo de vida derivado)
Não
Parcial (permissões por objeto)
Sim (changelog do NetBox)
Knot DNS
Não
Não
Sim (assinatura automática)
Não
Não
octoDNS e DNSControl
Parcial (conforme provedor)
Sim (estado declarativo)
Parcial
Não (por repositório)
Sim (histórico de commits)
Provedores em nuvem (Route 53, Cloudflare, Google)
Não (reverso raramente delegado)
Não
Sim
Sim (políticas de IAM)
Sim
Zonemaster
Não se aplica (validador)
Detecta, não corrige
Verifica, não provisiona
Não se aplica
Não se aplica
SanchezDNS (proposto)
Sim
Sim
Sim
Sim (leitura e escrita separadas)
Sim

Fonte: elaborado a partir das documentações oficiais e dos históricos de versão consultados em setembro de 2026.
A leitura do quadro impõe uma delimitação honesta da contribuição. Nenhum dos recursos do SanchezDNS é inédito isoladamente: a criação automática de registros reversos existe no Technitium desde 2020 e no produto da Microsoft há mais de duas décadas; o controle de acesso por zona está presente no Technitium, no Microsoft DNS e, de forma parcial, no PowerDNS-Admin; a trilha de auditoria é oferecida pelo Poweradmin com registro do estado anterior e posterior dos registros (POWERADMIN, 2026); a validação prévia à publicação é central nas ferramentas de infraestrutura como código; e a verificação formal da qualidade da delegação é o objeto do Zonemaster. A afirmação de ineditismo funcional seria, portanto, insustentável.
A contribuição reside em três aspectos precisos. O primeiro é a combinação, em uma única camada desacoplada e aplicável a infraestruturas locais já padronizadas sobre o PowerDNS, de mecanismos hoje dispersos entre ferramentas que não podem ser adotadas simultaneamente, uma vez que competem pelo papel de servidor autoritativo ou de fonte da verdade. O segundo é o tratamento do ciclo de vida completo do registro reverso: nenhuma das ferramentas examinadas que cria o PTR automaticamente também remove o reverso quando o registro de origem é excluído nem oferece rotina de triagem de reversos órfãos — lacuna que o Quadro 2 evidencia na coluna correspondente e cuja relevância é atestada pelos achados de Liu, Hao e Wang (2016). O terceiro, de natureza propriamente científica, é a avaliação empírica desses mecanismos segundo um desenho explicitado previamente à construção do artefato, com taxonomia de erros definida a priori, oráculo de medição independente e estudo de caso em ambiente real — elemento raramente presente na documentação das ferramentas comparadas, que reportam funcionalidades, e não efeitos medidos sobre a incidência de erros. Nos termos de Gregor e Hevner (2013), o trabalho situa-se no quadrante de melhoria (improvement), com contribuições nos níveis 1 e 2 de abstração.
2.6 ESTRATÉGIA DE AVALIAÇÃO DO ARTEFATO
A seção 2.4.3 apresentou o arcabouço FEDS e a escala SUS como instrumentos metodológicos. Esta seção estabelece como esses instrumentos serão efetivamente aplicados para sustentar — ou refutar — a asserção de que o SanchezDNS mitiga os problemas caracterizados na seção 2.2. A explicitação do desenho de avaliação antes da construção do artefato cumpre duas funções de rigor: impede que os critérios de êxito sejam formulados a posteriori, em função dos resultados obtidos, prática que compromete a validade das conclusões; e permite que a instrumentação necessária à coleta de dados seja incorporada ao artefato desde o início do desenvolvimento, e não improvisada ao término.
2.6.1 Escolha da estratégia de avaliação
Venable, Pries-Heje e Baskerville (2016) propõem quatro estratégias de avaliação, cuja seleção decorre da natureza do risco dominante no projeto e do custo relativo dos episódios: Quick and Dirty, apropriada quando riscos técnicos e sociais são baixos; Technical Risk and Efficacy, quando o risco predominante é a viabilidade técnica; Human Risk and Effectiveness, quando o risco predominante é o de que o artefato, ainda que tecnicamente correto, não produza o efeito pretendido em uso real; e Purely Technical Artifact, quando não há usuários humanos envolvidos.
No caso em exame, o risco técnico é reduzido: o artefato consiste em uma camada de mediação sobre uma API REST documentada e estável, sem componentes algorítmicos de resultado incerto. O risco dominante é humano e organizacional — a possibilidade de que os mecanismos propostos, embora funcionem conforme especificado, não reduzam a incidência de erros porque os operadores os contornam, os desativam, passam a confiar excessivamente na automação ou porque o erro se desloca para etapas não cobertas pela plataforma. Adota-se, por conseguinte, a estratégia Human Risk and Effectiveness, caracterizada por episódios formativos e artificiais nas etapas iniciais, com progressão relativamente rápida para episódios naturalísticos e somativos com usuários reais, de modo a submeter precocemente as hipóteses de projeto ao confronto com o uso efetivo.
2.6.2 Episódios de avaliação
O desenho compreende quatro episódios, precedidos pela caracterização documental que responde a QP1. A numeração é empregada de forma consistente no Quadro 3 e nos capítulos subsequentes.
E0 — Caracterização documental (formativo, artificial). Consiste na construção da taxonomia apresentada no Quadro 1 e na matriz de rastreabilidade que associa cada requisito do artefato a uma classe de erro documentada na literatura. O episódio responde a QP1 e estabelece o critério de suficiência do escopo: um requisito sem classe de erro associada é candidato à remoção, e uma classe eliminável sem requisito associado indica lacuna de projeto.
E1 — Conformidade técnica (formativo, artificial). Executado sobre instância de laboratório composta pelo artefato, por uma instância do PowerDNS e por um ambiente de delegação de teste. Para cada classe de erro eliminável do Quadro 1, constrói-se ao menos um caso de teste negativo que tenta reproduzir o erro pela interface e pela API do artefato; o caso é aprovado quando o erro é impedido ou corrigido automaticamente, e reprovado quando é aceito e publicado. Incluem-se testes de autorização negativos que verificam especificamente o caminho de automação do registro reverso, endereçando a questão QP3 e a classe de falha documentada no incidente descrito na seção 2.5.3. O estado resultante de cada zona é verificado por dois oráculos independentes do artefato: consulta direta ao servidor autoritativo, comparada ao estado esperado, e execução do Zonemaster sobre a zona publicada, registrando-se as reprovações por categoria de caso de teste. O emprego de oráculo externo é deliberado: avaliar o artefato exclusivamente pelas validações que ele próprio implementa constituiria raciocínio circular, incapaz de detectar precisamente os erros que o projetista não antecipou.
E2 — Experimento controlado com operadores (somativo, artificial). Desenho intrassujeitos, com contrabalanceamento da ordem dos tratamentos para neutralizar o efeito de aprendizagem. Cada participante executa dois conjuntos equivalentes de tarefas — criação de zona com ativação de DNSSEC; inclusão de registro de host com o reverso correspondente; alteração do endereço de um host já publicado; remoção de um host e de seu reverso; inclusão de registros MX e TXT com sintaxe específica; e correção de um conjunto de registros NS —, um conjunto sob edição manual assistida por linha de comando e outro sob o SanchezDNS. As variáveis dependentes são a taxa de tarefas concluídas com a zona em estado correto, o número de erros por tarefa classificados segundo o Quadro 1, o tempo até a conclusão e o número de tentativas. A análise emprega testes não paramétricos para amostras pareadas, com nível de significância de 0,05 e relato do tamanho de efeito. Sauro e Lewis (2016) indicam que amostras de dez a doze participantes em desenho pareado são suficientes para detectar efeitos de grande magnitude em métricas de desempenho; adota-se doze participantes como meta, reconhecendo-se a limitação de poder estatístico para efeitos pequenos, explicitada na seção 2.6.4.
E3 — Estudo de caso em ambiente de produção (somativo, naturalístico). Constitui o episódio decisivo para a questão QP4 e é detalhado na seção 2.6.3.
E4 — Usabilidade percebida (somativo). Aplicação da System Usability Scale ao término de cada tratamento em E2 e ao encerramento do período de observação em E3, com interpretação segundo as faixas estabelecidas por Bangor, Kortum e Miller (2008) e decomposição nas subescalas de usabilidade e de facilidade de aprendizado. O escore é tratado como medida de percepção, e não de desempenho, e sua leitura é condicionada aos resultados objetivos dos demais episódios.
2.6.3 O estudo de caso em ambiente de produção
O episódio E3 responde à exigência metodológica de demonstrar o efeito do artefato em condições reais de operação, e não apenas em ambiente controlado. Adota-se a estratégia de estudo de caso único e incorporado (single embedded case), nos termos propostos por Yin (2018), conduzido em organização que mantenha zonas autoritativas em produção sobre o PowerDNS Authoritative Server. A unidade de análise é a operação de alterações em zonas ao longo do período de observação; as unidades incorporadas são as zonas individuais e os operadores que sobre elas atuam. A escolha pelo estudo de caso, e não por um quase-experimento de campo, decorre da impossibilidade prática de manter dois grupos de operadores atuando simultaneamente sobre a mesma infraestrutura sob procedimentos distintos.
O desenho compreende três fases:
Fase 1 — Linha de base (quatro semanas). Registro das alterações conduzidas pelo procedimento vigente na organização, com coleta dos indicadores definidos adiante, sem qualquer intervenção. Realiza-se, ao início da fase, uma varredura completa das zonas sob gestão com o Zonemaster e um inventário de registros reversos órfãos, estabelecendo o estado inicial.
Fase 2 — Intervenção. Implantação do artefato, migração das credenciais e das permissões, e treinamento dos operadores, com registro do tempo despendido nessas atividades, que integra o custo de adoção.
Fase 3 — Observação (oito semanas). Operação das zonas exclusivamente pela plataforma, com coleta dos mesmos indicadores da Fase 1 e varreduras semanais com o Zonemaster.
Em atendimento ao princípio da triangulação de fontes de evidência (YIN, 2018), os dados provêm de quatro origens independentes: a trilha de auditoria do artefato, que fornece autoria, instante e conteúdo de cada alteração; os registros do sistema de chamados da organização, que permitem identificar retrabalho e incidentes atribuíveis a configuração; as varreduras periódicas com o Zonemaster, que medem a qualidade objetiva das zonas publicadas independentemente do que a plataforma reporta; e entrevistas semiestruturadas de encerramento com os operadores, que captam contornos, dificuldades e percepções não observáveis nos registros. A convergência entre fontes é condição para a aceitação de cada achado; divergências entre elas são reportadas, e não suprimidas.
Os indicadores coletados em ambas as fases são:
número de alterações publicadas e número de alterações que exigiram correção subsequente em até setenta e duas horas, adotado como indicador operacional de erro;
número de incidentes registrados no sistema de chamados cuja causa raiz seja atribuída a erro de configuração de DNS;
número de casos de teste do Zonemaster reprovados por categoria, ao início da Fase 1 e ao término da Fase 3;
proporção de registros A e AAAA sem o registro PTR correspondente e número de registros reversos órfãos, medidos nos mesmos dois instantes;
tempo mediano decorrido entre a solicitação da alteração e sua publicação efetiva;
proporção de alterações rastreáveis a um autor identificado.
Reconhece-se que o estudo de caso único não autoriza generalização estatística para uma população de organizações. A inferência pretendida é de natureza analítica (YIN, 2018): o caso é confrontado com a proposição teórica de que a eliminação por construção de determinadas classes de erro reduz sua incidência em operação real, e o resultado corrobora ou refuta essa proposição, sem pretender estimar a magnitude do efeito em outros contextos. Prevê-se, ademais, um plano de contingência: caso não se obtenha acesso a uma organização disposta a submeter sua infraestrutura de produção ao estudo dentro do cronograma da pesquisa, a Fase 3 será conduzida no ambiente de homologação de uma organização parceira, com zonas reais porém sem tráfego de usuários finais, registrando-se explicitamente a consequente redução de validade externa e a impossibilidade de medir incidentes percebidos por terceiros.
2.6.4 Matriz de avaliação, ameaças à validade e limitações
O Quadro 3 consolida a correspondência entre questões de pesquisa, construtos, métricas, instrumentos, episódios e critérios de êxito. Sua função é impedir que a avaliação se dissolva em relato de funcionalidades: cada questão formulada na seção 1.3 possui, nesse arranjo, um critério verificável e uma fonte de dados definida previamente.
Quadro 3 – Matriz de avaliação: questões, construtos, métricas, instrumentos e critérios de êxito
QP
Construto
Métrica
Instrumento ou fonte
Ep.
Critério de êxito
QP1
Cobertura da caracterização
Proporção das classes do Quadro 1 classificadas quanto à eliminabilidade e associadas a um mecanismo ou a uma justificativa de exclusão
Revisão da literatura; matriz de rastreabilidade requisito-problema
E0
Todas as dez classes classificadas; toda classe eliminável possui mecanismo e caso de teste correspondentes
QP2
Eficácia técnica dos mecanismos
Percentual de tentativas de produzir cada erro que são impedidas ou corrigidas automaticamente pelo artefato
Suíte de casos de teste negativos; consulta direta ao servidor; Zonemaster
E1
100% das tentativas nas classes elimináveis impedidas; nenhuma reprovação nova no Zonemaster após operações pelo artefato
QP2
Riscos introduzidos pela automação
Número de efeitos colaterais não intencionais observados (registros criados ou removidos fora da intenção declarada)
Diferencial do estado da zona antes e depois de cada operação
E1
Zero efeitos colaterais não declarados na interface antes da confirmação
QP3
Isolamento de privilégio
Número de operações que produzem escrita em zona sobre a qual o autor não detém permissão
Testes de autorização negativos, inclusive no caminho de automação do PTR
E1
Zero. A automação do reverso deve falhar de modo explícito, e não silencioso, quando faltar permissão na zona reversa
QP4
Eficácia operacional
Taxa de tarefas concluídas com a zona em estado correto; número de erros por tarefa, classificados pelo Quadro 1
Experimento controlado intrassujeitos; oráculo de verificação de E1
E2
E3
Redução da taxa de erro em relação ao tratamento manual, com significância estatística (p < 0,05) e tamanho de efeito reportado
QP4
Eficiência
Tempo mediano até a conclusão correta de cada tarefa; número de tentativas
Cronometragem no experimento; carimbos temporais da trilha de auditoria
E2
E3
Tempo não superior ao do tratamento manual; qualquer aumento deve ser justificado pela redução de erros
QP4
Usabilidade percebida
Escore SUS (0 a 100) e escala adjetiva correspondente
Questionário SUS aplicado ao final de cada tratamento
E2
E4
Escore igual ou superior a 68, situando-se acima da média das aplicações relatadas na literatura
QP4
Rastreabilidade
Proporção das alterações publicadas no período de observação que possuem autor, instante e recurso registrados
Trilha de auditoria confrontada com o histórico de alterações do servidor
E3
100% das alterações efetuadas pela plataforma; divergências em relação ao servidor investigadas e reportadas

A validade das conclusões depende do reconhecimento explícito das ameaças que incidem sobre o desenho proposto, classificadas a seguir segundo as categorias empregadas por Wohlin et al. (2012).
Validade interna. No experimento E2, o efeito de aprendizagem e a ordem de apresentação dos tratamentos podem ser confundidos com o efeito da plataforma, ameaça mitigada pelo contrabalanceamento e pelo uso de conjuntos de tarefas equivalentes, porém não idênticos. No estudo de caso E3, a comparação entre fases sucessivas é vulnerável à maturação dos operadores, à variação sazonal da carga de trabalho e ao efeito Hawthorne, decorrente da consciência de estarem sendo observados; mitiga-se parcialmente pela extensão do período de observação e pelo uso de indicadores objetivos coletados automaticamente. Registra-se, ainda, que o pesquisador é também o desenvolvedor do artefato, configurando risco de viés do experimentador; adotam-se como contramedidas o roteiro padronizado de condução, a coleta automatizada das medidas, a verificação por oráculo externo e a classificação dos erros segundo taxonomia definida antes da coleta.
Validade de construto. O construto "erro de configuração" é operacionalizado pela taxonomia do Quadro 1 e pelos casos de teste do Zonemaster, recorte que não captura erros semânticos — um registro sintaticamente válido que aponta para o endereço errado é indistinguível, para esses instrumentos, de um registro correto. Essa limitação é intrínseca à abordagem e deve ser explicitada no relato dos resultados. Do mesmo modo, a escala SUS mede usabilidade percebida, e não desempenho, de modo que escores elevados não constituem evidência de redução de erros, nem a substituem.
Validade externa. A amostra de participantes de E2 é de conveniência, e o estudo de caso E3 restringe-se a uma organização, a um servidor autoritativo específico e a zonas de porte moderado. Os resultados não se estendem a operadores de domínios de topo, a infraestruturas com dezenas de milhares de zonas ou a servidores autoritativos diversos do PowerDNS.
Validade de conclusão. O tamanho amostral previsto para E2 confere poder estatístico adequado apenas à detecção de efeitos de grande magnitude; efeitos pequenos, ainda que reais, provavelmente não alcançarão significância. Adotam-se testes não paramétricos, em razão da distribuição esperada das medidas de tempo, e reportam-se tamanhos de efeito e intervalos de confiança em vez de valores-p isolados. Comparações múltiplas entre as seis tarefas são tratadas como exploratórias, com a devida ressalva no relato.
Ameaça específica de comparabilidade dos tratamentos. O contraste entre o SanchezDNS e a edição manual por linha de comando compara, simultaneamente, dois fatores: os mecanismos de prevenção propostos e a modalidade de interface, gráfica em um caso e textual no outro. Parte do efeito eventualmente observado pode, portanto, ser atribuível à modalidade, e não aos mecanismos. Trata-se da principal fragilidade do desenho. Prevê-se, para isolá-la, a inclusão de um terceiro tratamento sempre que a disponibilidade dos participantes permitir: um painel web sem os mecanismos propostos — o PowerDNS-Admin com a automação de PTR desativada, em sua configuração padrão —, de modo que a diferença entre esse tratamento e o artefato isole o efeito dos mecanismos, e a diferença entre ele e a linha de comando isole o efeito da modalidade. Caso o terceiro tratamento não se viabilize, a limitação é declarada e as conclusões são correspondentemente restringidas.
Quanto aos aspectos éticos, a participação em E2 e em E3 é voluntária e precedida de termo de consentimento livre e esclarecido, com esclarecimento de que o objeto de avaliação é a ferramenta, e não o desempenho individual do participante. Os dados são anonimizados na coleta, e nenhum registro que identifique operadores é retido após a consolidação dos resultados. No estudo de caso, nomes de domínio, endereços e demais dados de produção da organização são pseudonimizados no relato, preservando-se apenas as propriedades estruturais relevantes à análise. A submissão do protocolo ao comitê de ética em pesquisa segue as exigências aplicáveis da instituição.
2.7 CONSIDERAÇÕES DO CAPÍTULO
Este capítulo delimitou a fundamentação conceitual que sustenta o SanchezDNS. Revisou o modelo operacional do DNS e os componentes manipulados pelo artefato; demonstrou, a partir de estudos empíricos que abrangem duas décadas de medições, a prevalência persistente dos erros de configuração em zonas autoritativas, consolidando-os em uma taxonomia de dez classes associadas a evidências, impactos e mecanismos de mitigação; apresentou os pilares de controle de acesso e auditoria e o arcabouço metodológico do Design Science Research; e examinou os trabalhos correlatos, evidenciando, por meio de quadro comparativo, que a contribuição do artefato reside na combinação de mecanismos hoje dispersos, no tratamento do ciclo de vida completo dos registros reversos e na avaliação empírica de seus efeitos. Por fim, estabeleceu a estratégia de avaliação, com episódios, métricas, critérios de êxito definidos previamente à construção e reconhecimento explícito das ameaças à validade. O Capítulo 3 detalha a metodologia de condução da pesquisa e o modelo de engenharia do artefato, derivando os requisitos diretamente das classes de erro consolidadas no Quadro 1.
3 METODOLOGIA E MODELAGEM DO ARTEFATO
Este capítulo apresenta a metodologia aplicada ao desenvolvimento do SanchezDNS e a modelagem técnica de seus componentes. A condução do projeto adota o método Design Science Research (DSRM) de Peffers et al. (2007), derivando os requisitos das falhas operacionais levantadas no Capítulo 2. Em seguida, especificam-se os cinco modelos que orientam a implementação: controle de acesso, normalização sintática de registros, sincronização de zonas reversas e auditoria de eventos.
3.1 CLASSIFICAÇÃO DA PESQUISA
A pesquisa é aplicada quanto à natureza, descritiva e prescritiva quanto aos objetivos, e combina abordagem qualitativa — revisão da literatura e análise de soluções existentes — com abordagem quantitativa na mensuração da taxa de erros, do tempo de execução de tarefas e do escore de usabilidade.
O delineamento é o da Design Science Research. March e Smith (1995) distinguem as ciências naturais, que explicam fenômenos existentes, das ciências do artificial, que constroem e avaliam artefatos. Hevner et al. (2004) sistematizam essa distinção em sete diretrizes, das quais três condicionam o desenho deste trabalho: produzir um artefato viável, responder a um problema documentado e demonstrar utilidade por método de avaliação rigoroso.
O problema tratado no Capítulo 2 não é a falta de conhecimento sobre as falhas de configuração em DNS — a literatura de medição as documenta há mais de duas décadas —, mas a ausência de uma ferramenta que as impeça no momento da edição. É, portanto, um problema de projeto. Na classificação de Gregor e Hevner (2013), a contribuição situa-se no quadrante de melhoria: solução nova para problema conhecido. Os mecanismos isolados já existem em ferramentas maduras; a contribuição está na sua integração em um único fluxo e na medição do efeito dessa integração sobre a taxa de erro.
3.2 APLICAÇÃO DO MODELO DSRM
O modelo de Peffers et al. (2007) prevê seis atividades e quatro pontos possíveis de entrada. Adotou-se a entrada centrada no problema, já que a motivação partiu de falhas observadas na operação de zonas, e não de uma tecnologia à procura de aplicação. O Quadro 4 mapeia cada atividade à seção correspondente.
Quadro 4 – Atividades do DSRM e sua realização nesta pesquisa.
Atividade (PEFFERS et al., 2007)
Realização nesta pesquisa
Seção

1. Identificação do problema e motivação
   Revisão da literatura de medição; consolidação da taxonomia de dez classes de erro com evidência empírica e impacto documentado
   2.2; Quadro 1
2. Definição dos objetivos da solução
   Questões QP1 a QP4; formulação de requisitos a partir das classes elimináveis por construção
   1.3; 1.4; 3.3
3. Projeto e desenvolvimento
   Modelagem dos quatro componentes técnicos; implementação dos serviços com Go, Nuxt e banco de dados
   3.4 a 3.7; Capítulo 4
4. Demonstração
   Execução de tarefas de administração em bancada de laboratório, com conferência por oráculo externo independente
   5.1
5. Avaliação
   Ensaios formativos e somativos sob o framework FEDS; medição de erros, tempos de execução e escore SUS
   2.6; 5.2
6. Comunicação
   Formalização da monografia; disponibilização do código-fonte e documentação técnica sob licença livre
   Capítulo 6

Fonte: elaborado com base em Peffers et al. (2007).
O ciclo foi percorrido mais de uma vez. Dois retornos da avaliação ao projeto produziram alteração de requisito e estão documentados no Capítulo 4: a introdução da guarda de referência na remoção de registros PTR (seção 3.6.3) e a exibição do valor normalizado em lugar do valor digitado (seção 3.5.3). Ambos decorrem do mesmo efeito, antecipado pela questão QP2: a automação de uma decisão antes tomada pelo operador desloca o erro em vez de eliminá-lo.
3.3 REQUISITOS
Os requisitos não foram levantados por entrevista ou questionário. Foram derivados diretamente da taxonomia do Quadro 1, de modo que cada um responda a uma classe de erro com evidência empírica na literatura. O procedimento teve três passos:
Selecionaram-se as seis classes classificadas como elimináveis por construção, mais a classe de ausência de rastreabilidade;
Para cada classe, especificou-se o mecanismo que torna o estado incorreto inalcançável pela interface;
As classes dependentes de terceiros — delegação inconsistente e divergência de NS entre pai e filho — foram convertidas em requisitos de detecção, sem correção automática, conforme a delimitação da seção 1.5.
O Quadro 5 apresenta os requisitos funcionais resultantes, com a rastreabilidade para a classe de erro de origem e para a questão de pesquisa correspondente.
Quadro 5 – Requisitos funcionais e rastreabilidade com a taxonomia de erros.
Id
Requisito Funcional
Classe de erro (Quadro 1)
Questão de Pesquisa
RF01
Criar zonas com soa_edit_api = DEFAULT no payload de criação, sem expor o campo serial à edição manual
Serial SOA não incrementado
QP2
RF02
Gerar a KSK e aplicar NSEC3PARAM automaticamente no mesmo fluxo de criação da zona
DNSSEC incompleto
QP2
RF03
Exibir o resumo criptográfico (DS) da zona pronto para publicação no registrário e confrontá-lo com o DS publicado no pai
DNSSEC incompleto
QP2
RF04
Normalizar o RDATA conforme a sintaxe exigida pelo tipo de registro antes do envio à API do PowerDNS
Erro sintático no RDATA
QP2
RF05
Completar nomes relativos em nomes plenamente qualificados (FQDN) a partir da origem da zona
Erro sintático no RDATA
QP2
RF06
Criar o registro PTR correspondente ao gravar um registro direto A ou AAAA
Divergência entre A/AAAA e PTR
QP2, QP3
RF07
Reposicionar o registro PTR ao alterar o endereço IP de um registro A ou AAAA
Divergência entre A/AAAA e PTR
QP2, QP3
RF08
Remover o registro PTR ao excluir o registro A/AAAA de origem, desde que nenhum outro registro referencie o mesmo IP
Registros pendentes e reversos órfãos
QP2
RF09
Listar os registros PTR de uma zona reversa sem apontamento direto correspondente e permitir remoção assistida
Registros pendentes e reversos órfãos
QP2
RF10
Listar os registros A/AAAA sem PTR correspondente e permitir criação atômica em lote
Divergência entre A/AAAA e PTR
QP2
RF11
Atribuir permissão de leitura ou escrita por par (usuário, zona), conjugando com um papel administrativo global
Privilégio excessivo do operador
QP3
RF12
Condicionar o cadastro de novos usuários à aprovação explícita de um administrador
Privilégio excessivo do operador
QP3
RF13
Gravar em trilha centralizada toda operação de escrita sobre zonas, registros e usuários
Ausência de rastreabilidade
QP4
RF14
Permitir consulta à trilha de auditoria com paginação e filtragem por autor, ação, zona e conteúdo
Ausência de rastreabilidade
QP4

Fonte: elaborado pelo autor (2026).
Quatro requisitos não funcionais decorrem de decisões de arquitetura, detalhadas no Capítulo 4:
RNF01 — O PowerDNS é a fonte única da verdade dos dados de zona. A aplicação não mantém cópia local dos registros; toda leitura é uma chamada à API do servidor autoritativo.
RNF02 — Toda autorização é decidida no backend, por requisição. O frontend apenas reflete o resultado.
RNF03 — A falha de uma escrita derivada não aborta a operação principal.
RNF04 — Não existe caminho de escrita não autenticado. Todas as rotas de alteração exigem sessão válida, o que elimina por arquitetura a classe de atualização dinâmica não autenticada medida por Korczyński et al. (2016).
O RNF03 é uma escolha entre dois modos de falha, e não uma propriedade desejável em si. A seção 3.6.4 examina sua consequência na integridade global da solução.
3.4 AUTENTICAÇÃO E CONTROLE DE ACESSO
3.4.1 Autenticação e ciclo de vida da sessão
A autenticação é baseada em sessão opaca mantida no servidor, e não em token autocontido. A escolha decorre do RNF02: como toda autorização é decidida no backend a cada requisição, o identificador entregue ao cliente não precisa transportar informação alguma — precisa apenas ser impossível de adivinhar e possível de revogar. Um token autocontido do tipo JWT inverteria as duas propriedades, ao embutir o nível de acesso no próprio artefato entregue ao cliente e ao permanecer válido até expirar, mesmo após a revogação do privilégio que ele declara.
Na autenticação bem-sucedida, o backend gera um identificador de sessão (SID) de 32 bytes obtidos do gerador criptográfico do sistema operacional, codificados em base64 sem preenchimento. São 256 bits de entropia, o que torna a adivinhação inviável por força bruta. O SID é persistido na coleção sessions do MongoDB, que mantém índice único sobre o campo sid, e devolvido ao cliente em cookie.
O documento de sessão registra, além do identificador e do e-mail do operador, os metadados de origem da conexão: endereço IP, sistema operacional e navegador derivados do cabeçalho User-Agent, e localização geográfica aproximada resolvida a partir do IP. Esses campos não participam da decisão de autorização; existem para dar contexto forense à trilha de auditoria da seção 3.7, que registra o autor da operação mas não as condições em que ela foi realizada.
O cookie é declarado com o atributo HttpOnly, o que impede sua leitura por código JavaScript na página e neutraliza o roubo de sessão por injeção de script. O atributo Secure, que restringe o envio ao canal TLS, é condicionado ao ambiente: ativo em produção e inativo em desenvolvimento, onde a aplicação roda sobre HTTP simples.
A validação ocorre no middleware de sessão, em duas camadas. Consulta-se primeiro o cache em Redis; havendo acerto, verifica-se apenas o campo de atividade. Na ausência do registro em cache, recorre-se ao MongoDB com filtro composto por SID e sessão ativa. Resolvido o e-mail, o nível global é obtido e injetado no contexto da requisição. O cache existe para evitar uma consulta ao banco por requisição; a fonte da verdade continua sendo o MongoDB, e a ausência no Redis nunca é interpretada como sessão inválida.
A expiração é por inatividade, não por tempo absoluto. Cada requisição autenticada atualiza o campo de último acesso da sessão. Uma rotina agendada executa a cada dez minutos, seleciona as sessões ativas cujo último acesso seja anterior a sete dias, marca-as como inativas com o motivo registrado e remove as entradas correspondentes do cache. A revogação é, portanto, imediata em efeito: basta marcar a sessão como inativa para que a próxima requisição seja recusada, sem esperar a expiração de um token.
Há uma assimetria deliberada entre a validade do cookie e a da sessão. O cookie é emitido com prazo longo, para que o operador não precise reautenticar-se a cada visita; a sessão no servidor expira em sete dias de inatividade. Quem decide é sempre o servidor: um cookie ainda dentro do prazo cuja sessão tenha sido desativada produz recusa e o cookie é descartado na própria resposta. O cookie é portador de credencial, não a credencial.
As senhas nunca são armazenadas em texto claro. O cadastro grava o resultado de bcrypt com o fator de custo padrão da biblioteca, prefixado pelo identificador do esquema, o que permite reconhecer o algoritmo empregado e migrar para outro sem ambiguidade. A verificação usa comparação em tempo constante da própria biblioteca, e hashes marcados com o esquema legado de cripta tradicional são recusados de forma incondicional, sem tentativa de verificação.
O Quadro 6 consolida os parâmetros do modelo.
Quadro 6 – Parâmetros do modelo de sessão e autenticação.
Parâmetro
Valor adotado
Finalidade
Entropia do identificador de sessão
256 bits (32 bytes do gerador criptográfico)
Inviabilizar adivinhação do SID
Formato do identificador
Opaco, sem dados embutidos
Manter a decisão de autorização no servidor (RNF02)
Armazenamento
MongoDB com índice único sobre o SID; cache em Redis
Revogação imediata e leitura de baixo custo
Atributos do cookie
HttpOnly sempre; Secure em produção
Impedir leitura por script e trânsito fora de TLS
Expiração
Sete dias de inatividade, apurada a cada dez minutos
Encerrar sessões abandonadas sem penalizar o uso contínuo
Derivação da senha
bcrypt com fator de custo padrão e prefixo de esquema
Impedir recuperação da senha a partir do banco
Resposta a credencial inválida
Mensagem única para e-mail inexistente e senha incorreta
Impedir enumeração de usuários

Fonte: elaborado pelo autor (2026).
3.4.2 Cadastro, aprovação e inicialização
O RF12 condiciona o acesso à aprovação explícita de um administrador. O fluxo tem três caminhos.
O primeiro cadastro realizado em uma instância vazia é gravado diretamente na coleção de cadastros com nível administrativo. A regra resolve a inicialização do sistema sem credencial padrão embutida no código ou em variável de ambiente, que seria ela própria uma vulnerabilidade, e sem exigir intervenção manual no banco.
Os demais cadastros são gravados na coleção de solicitações, com estado pendente. A tentativa de autenticação de um solicitante ainda não apreciado é recusada com mensagem específica informando que o cadastro está em análise — distinção deliberada em relação à mensagem de credencial inválida, por não revelar nada que o próprio solicitante já não saiba.
A aprovação por um administrador copia os dados da solicitação para a coleção de cadastros com nível de membro e registra o desfecho na solicitação, preservando o histórico da decisão. Um membro recém-aprovado não tem acesso a zona alguma: o vínculo por zona é concedido separadamente, conforme a seção 3.4.3.
A validação de entrada no cadastro exige endereço de correio eletrônico sintaticamente válido e senha de no mínimo oito caracteres, e normaliza o endereço em caixa baixa antes da gravação, de modo que a comparação posterior nos vetores de permissão seja determinística.
3.4.3 Estrutura das permissões
O SanchezDNS combina papéis globais com permissões por zona, em um modelo RBAC híbrido derivado de Sandhu et al. (1996). O objetivo é isolar o raio de impacto operacional: um operador responsável por um domínio não deve enxergar nem alterar zonas sob custódia de outra equipe, conforme o princípio de menor privilégio de Saltzer e Schroeder (1975).
O RBAC puro não resolve o caso. Um papel global de "operador de zona" concederia acesso a todas as zonas; um papel por zona multiplicaria papéis sem ganho de abstração. A autorização opera, por isso, em dois níveis:
Nível global: Definido pelo campo level do documento de cadastro do operador, com dois valores possíveis. O perfil admin detém privilégio irrestrito sobre criação e remoção de zonas, gestão de usuários e consulta à auditoria. O perfil member não concede, por si, acesso a zona alguma. O primeiro cadastro criado na instância recebe admin automaticamente, o que dispensa credencial embutida no código ou em variável de ambiente.
Nível de zona: A coleção users do MongoDB mantém um documento por zona, contendo dois vetores de endereços de e-mail:
JSON
{
"zona": "exemplo.com.br.",
"leitura": ["operador1@org.br", "operador2@org.br"],
"escrita": ["admin-zona@org.br"]
}

A resolução da permissão segue precedência fixa: verifica-se se level == "admin"; na ausência de privilégio global, verifica-se a pertinência ao vetor escrita; na ausência desta, a pertinência ao vetor leitura; esgotadas as verificações, o acesso é sumariamente negado com status HTTP 403 (Forbidden). A comparação de endereços de e-mail é normalizada em caixa baixa (case-insensitive). Como a escrita pressupõe e engloba a visualização, a concessão de escrita herda implicitamente o privilégio de leitura.
3.4.4 Aplicação nas rotas
A autorização é implementada em dois middlewares e uma função de resolução por zona. O middleware de sessão, descrito na seção 3.4.1, injeta email e level no contexto da requisição e recusa com HTTP 400 a requisição sem sessão válida, descartando o cookie na resposta. O middleware administrativo protege o grupo de rotas administrativas, recusando com HTTP 403 qualquer requisição cujo contexto não tenha level == "admin".
O middleware administrativo protege o grupo de rotas administrativas, recusando com HTTP 403 qualquer requisição cujo contexto não tenha level == "admin".
As rotas de zona não são cobertas por middleware, porque a decisão depende do parâmetro zone da requisição. Cada controlador invoca a resolução de permissão antes de executar a operação, e retorna HTTP 403 quando o nível apurado é insuficiente. O Quadro 7 apresenta a matriz resultante.
Quadro 7 – Matriz de autorização e controle de rotas por nível de acesso.
Rota da API
Método HTTP
Administrador Global
Escrita na Zona
Leitura na Zona
Sem Vínculo
/zones
GET
Todas as zonas
Próprias
Próprias
Lista vazia
/records
GET
Sim
Sim
Sim
HTTP 403
/zone/dnssec
GET
Sim
Sim
Sim
HTTP 403
/records
PUT, PATCH, DELETE
Sim
Sim
HTTP 403
HTTP 403
/reverses
GET, PUT
Sim
Sim
HTTP 403
HTTP 403
/reverses/orphans
GET
Sim
Sim
HTTP 403
HTTP 403
/zone
PUT, DELETE
Sim
HTTP 403
HTTP 403
HTTP 403
/soa
PATCH
Sim
HTTP 403
HTTP 403
HTTP 403
/zone/nsec3
PATCH
Sim
HTTP 403
HTTP 403
HTTP 403
/user, /users
POST, PATCH, DELETE, GET
Sim
HTTP 403
HTTP 403
HTTP 403
/solicitacoes/status
PUT
Sim
HTTP 403
HTTP 403
HTTP 403
/logs
GET
Sim
HTTP 403
HTTP 403
HTTP 403

Fonte: elaborado pelo autor (2026).
Duas restrições do modelo merecem registro. A criação e a remoção de zonas são privativas do administrador porque produzem efeito fora da zona: a criação inaugura uma delegação que depende de ação no registrário, e a remoção apaga dados que a plataforma não é capaz de restaurar. A alteração do SOA e a aplicação de NSEC3 também são privativas do administrador, por alterarem parâmetros de replicação e de assinatura da zona inteira.
3.4.5 Escrita derivada em zona de outro escopo
A sincronização do registro reverso introduz um conflito de autorização. O operador edita a zona direta, onde tem permissão de escrita, mas o PTR correspondente é gravado sob in-addr.arpa ou ip6.arpa, que pode pertencer a outra equipe.
No SanchezDNS a sincronização é tratada como rotina derivada do sistema. O backend usa a chave de API da própria plataforma junto ao PowerDNS para gravar na zona reversa, sob duas condições:
A zona reversa correspondente ao endereço já existe no servidor autoritativo;
Não há RRset PTR anterior para o mesmo nome derivado.
O operador não recebe acesso à zona reversa: não pode listá-la, editá-la nem remover dela registros arbitrários. O que ele obtém é o efeito determinístico da escrita que já estava autorizado a fazer na zona direta — o nome do PTR é função do endereço, e seu conteúdo é função do nome do registro de origem.
A primeira condição é o que delimita a concessão. Como a plataforma não cria zonas reversas por iniciativa própria, a equipe responsável pelo espaço reverso controla quais faixas ficam sujeitas à derivação: basta não provisionar a zona. Onde ela não existe, nenhuma escrita ocorre e a operação direta se completa normalmente.
A solução é um afrouxamento deliberado do menor privilégio, cujos desdobramentos e riscos operacionais associados são discutidos criticamente entre as limitações do artefato no Capítulo 6.
3.5 NORMALIZAÇÃO SINTÁTICA
3.5.1 Estratégia
Para impedir que RDATA malformado chegue à API do PowerDNS, o backend aplica rotinas de normalização antes da submissão. Em vez de apenas rejeitar a entrada, a plataforma converte o valor para a forma canônica do tipo, o que dispensa o operador de conhecer a sintaxe exigida. É a aplicação do conceito de forcing function descrito por Norman (2013): o estado incorreto não chega a ser representável.
A validação é reservada aos casos sem conversão determinística. O esquema de validação do formulário, avaliado no frontend com a biblioteca Valibot e reavaliado no backend, exige os campos obrigatórios de cada tipo — prioridade, peso, porta e alvo para SRV; prioridade e nome de destino para HTTPS — e recusa TTL inferior a 60 segundos. A ausência de um campo obrigatório não admite suprimento automático e produz HTTP 400.
3.5.2 Regras por tipo
Quadro 8 – Regras de normalização sintática aplicadas no backend.
Tipo de Registro
Transformação e Sanitização Aplicada
Erro Prevenido
Todos (campo nome)
Nome sem ponto final recebe a origem da zona e o ponto terminal; nome já terminado em ponto é preservado; nome que já contém a zona recebe apenas o ponto
Duplicação do domínio base (host.dominio.com.br.dominio.com.br.)
A, AAAA
Conteúdo é validado e parseado como endereço IP; valor não interpretável aborta a operação
Endereço malformado gravado como cadeia textual inválida
CNAME, NS, ALIAS, PTR
Ponto final acrescentado ao conteúdo quando ausente
Destino interpretado como relativo à origem da zona (BARR, 1996)
MX
Ponto final inserido no host de destino; prioridade numérica capturada em campo isolado e prefixada ao conteúdo
Omissão da prioridade ou inversão na ordem dos campos
TXT
Aspas delimitadoras acrescentadas automaticamente nas extremidades quando ausentes
Cadeia recusada pelo servidor autoritativo
SRV
Prioridade, peso, porta e alvo capturados em campos isolados e serializados na ordem canônica; ponto final no alvo
Inversão ou omissão da ordem dos quatro parâmetros
CAA
Composição estruturada de flag e tag quando o valor não contém issue, issuewild ou iodef
Omissão da diretiva de emissão de certificados TLS
HTTPS
Prioridade e nome de destino em campos próprios; destino vazio serializado como .; SvcParams recebe alpn=h2 quando omitido
Omissão do destino canônico no modo de serviço
Todos (campo TTL)
Valores numéricos inferiores a 60 segundos são bloqueados na validação
TTL excessivamente volátil, elevando a carga sobre resolvedores

Fonte: elaborado com base em Barr (1996) e nas especificações normativas de cada RRset.
A diretriz comum ao quadro é a decomposição dos tipos de conteúdo estruturado em campos isolados do formulário. MX, SRV e HTTPS têm RDATA composto por elementos em ordem fixa; editá-lo como cadeia única transfere ao operador a responsabilidade de lembrar a ordem. Com a decomposição, o operador preenche campos rotulados e a serialização na ordem das RFCs é feita pelo backend.
3.5.3 Devolução do valor normalizado
A normalização introduz um risco: ao converter o valor, o sistema toma uma decisão que o operador supõe ter tomado. Se a conversão não for devolvida à tela, o operador fica com modelo mental divergente do que foi publicado. O caso mais sensível é o do ponto final, em que a diferença entre o digitado e o publicado é de um caractere.
A regra adotada, em resposta ao segundo retorno de iteração da seção 3.2, é que a tela nunca exibe o valor digitado. Após a operação, o frontend recarrega os registros a partir do endpoint GET /records, que os lê do PowerDNS. O valor exibido é o valor efetivamente publicado. A regra decorre do RNF01 e torna imediatamente visível qualquer divergência entre a intenção do operador e o estado da zona.
3.6 SINCRONIZAÇÃO DE ZONAS REVERSAS
Esta seção especifica a automação que atende aos requisitos RF06 a RF10. É o componente de maior complexidade do artefato: opera sobre duas zonas, depende de uma correspondência que o protocolo não impõe mas que serviços de rede exigem (BARR, 1996), e é o único em que a remoção pode criar um erro que não existia.
3.6.1 Derivação do nome e da zona
A gravação do PTR exige duas informações derivadas do endereço.
Nome reverso. Para IPv4, os quatro octetos são invertidos sob in-addr.arpa: 192.0.2.10 produz 10.2.0.192.in-addr.arpa.. Para IPv6, o endereço é primeiro expandido para a forma plena de 32 dígitos hexadecimais, e cada nibble vira um rótulo isolado em ordem inversa sob ip6.arpa.. A expansão prévia é indispensável: a notação abreviada, que suprime sequências de zeros, não admite conversão direta em rótulos.
Zona de destino. O backend lista as zonas do servidor e seleciona, entre as que são sufixo do nome reverso, a de maior comprimento. A regra do sufixo mais longo acomoda a delegação de faixas menores que o octeto (EIDNES; DE GROOT; VIXIE, 1998): coexistindo 2.0.192.in-addr.arpa. e 128/25.2.0.192.in-addr.arpa., o registro vai para a mais específica, que é a que responde pelo endereço. Escolher a mais ampla gravaria o registro em zona sem autoridade sobre o nome — ele existiria no banco e não responderia às consultas.
Se nenhuma zona reversa casar com o endereço, a função retorna sem erro e nenhuma escrita é feita.
3.6.2 Criação e alteração
Na criação de um registro A ou AAAA, o backend deriva o nome reverso, localiza a zona e verifica a existência de RRset PTR anterior. Havendo registro anterior, a rotina não o substitui: a escrita derivada não sobrepõe uma decisão tomada pelo responsável da zona reversa. O estado resultante — registro direto sem PTR — é detectável pelo endpoint GET /reverses, que implementa o RF10.
A alteração de endereço é executada como duas operações, e não como uma substituição:
Grava o PTR do novo endereço;
Se o endereço anterior for diferente do novo, avalia a remoção do PTR anterior.
A ordem importa. Com a criação antes da remoção, uma falha na segunda etapa deixa o sistema com dois reversos, um deles obsoleto — estado redundante, mas resolvível. A ordem inversa deixaria o endereço sem reverso, que é o estado que faz servidores SMTP rejeitarem mensagens.
3.6.3 Guarda de referência
A remoção é a operação que exige mais cautela, e foi a que motivou o primeiro retorno de iteração da seção 3.2. A formulação direta do RF08 — apagar o PTR ao apagar o registro A — quebra sempre que mais de um nome aponta para o mesmo IP, caso corriqueiro em servidores com vários serviços. A remoção de um dos nomes apagaria o reverso ainda referenciado pelos demais, transformando uma limpeza na própria divergência A/PTR que o requisito combate.
O backend aplica, por isso, uma guarda de referência antes de qualquer remoção de PTR: percorre todas as zonas diretas do servidor — ignorando as terminadas em in-addr.arpa e ip6.arpa — e verifica se algum RRset A ou AAAA ainda contém o endereço. Havendo referência, a remoção é suprimida. A mesma guarda roda na etapa de remoção da alteração descrita acima.
Daí decorre uma assimetria entre criação e remoção que vale explicitar como regra de projeto. A criação é otimista: executa salvo decisão anterior em contrário. A remoção é conservadora: executa apenas se nenhuma referência subsistir. A assimetria vem da diferença de custo entre os dois erros. Um reverso a mais é detectado pela triagem de órfãos, não quebra resolução e é corrigível a qualquer tempo. Um reverso removido indevidamente quebra serviços imediatamente e não é detectável sem verificação ativa. Diante de erros de custo assimétrico, a automação é enviesada na direção do erro mais barato.
(Local reservado para inserção da Figura 1: Fluxograma da Criação Otimista e Remoção Conservadora de PTR)
3.6.4 Semântica de melhor esforço
As rotinas de reverso operam sob o RNF03: a falha da escrita derivada não aborta a operação principal. Cada uma roda com contexto de tempo limite próprio de seis segundos, para que a lentidão do servidor autoritativo não trave a requisição do operador.
A decisão introduz um risco que integra a resposta à segunda parte da QP2. O operador recebe HTTP 200 e a mensagem de sucesso; a escrita derivada pode não ter ocorrido. O estado resultante é a divergência A/PTR que o RF06 pretende eliminar — agravada pelo fato de o operador ter motivo para supor que ela não existe. A automação, nesse caso, não elimina a classe de erro: transfere-a do esquecimento do operador para uma falha silenciosa do sistema, e piora a condição de detecção.
Três mitigações fazem parte do projeto, e nenhuma resolve o problema por completo:
O endpoint GET /reverses (RF10) lista os registros A/AAAA sem PTR e permite criação em lote por PUT /reverses, o que converte a falha silenciosa em achado verificável — desde que o operador execute a verificação;
A falha da rotina é gravada no log da aplicação, preservando rastreabilidade do evento ainda que ele não seja comunicado na hora;
A ordenação da seção 3.6.2 restringe o estado de falha parcial à redundância, não à ausência.
A escolha preserva a autonomia sobre o recurso próprio ao custo de uma garantia mais fraca sobre o derivado, limitação arquitetural retomada e aprofundada no Capítulo 6.
A criação em lote do RF10 é sequencial e não transacional: cada PTR selecionado é submetido em requisição própria à API do PowerDNS e a primeira falha interrompe o processamento, preservando os registros já gravados e deixando os restantes por criar. O estado parcial é convergente — basta reexecutar a rotina, que recalcula os ausentes a partir do estado corrente das zonas —, mas a operação não oferece garantia de atomicidade sobre o conjunto selecionado.
3.7 AUDITORIA
A auditoria atende aos requisitos RF13 e RF14. Kent e Souppaya (2006) identificam três finalidades para o registro de eventos — reconstrução de incidentes, atribuição de responsabilidade e detecção de uso indevido —, cada uma com exigências próprias de granularidade e retenção.
Cada operação de escrita grava um documento na coleção logs do MongoDB:
JSON
{
"zone": "exemplo.com.br.",
"username": "operador@org.br",
"action": "edit_record",
"details": "Alterado registro A de www.exemplo.com.br.",
"createdAt": "2026-03-14T10:22:47-03:00"
}

O campo action usa códigos estáveis, e não texto livre, para permitir agregação. O endpoint GET /logs, restrito ao administrador, oferece paginação e filtro por expressão regular sobre username, action, zone e details, com a entrada escapada antes da composição da consulta.
O Quadro 9 relaciona as operações contempladas.
Quadro 9 – Operações e códigos de ação registrados na trilha de auditoria.
Domínio Operacional
Códigos de Ação Padronizados
Zonas
create_zone, delete_zone, update_soa, update_nsec3
Registros
insert_record, edit_record, delete_record
Registros Reversos
insert_reverses, delete_reverse
Usuários e Acessos
insert_user, update_user, delete_user

Fonte: elaborado pelo autor (2026).
Três decisões de modelagem:
A trilha fica no MongoDB, não no PowerDNS. É consequência do RNF01. O PowerDNS é a fonte da verdade do estado das zonas, mas não guarda quem alterou o quê: sua API autentica a aplicação por chave, não o operador. A identidade só existe na camada de operação e é lá que precisa ser persistida.
O registro é descritivo, não diferencial. Grava-se a descrição do efeito, não o estado anterior e posterior do RRset. A opção atende à atribuição de responsabilidade e à reconstrução da sequência de eventos, mas não à restituição de estado. A limitação delimita o que a avaliação poderá afirmar sobre a QP4: a trilha responde quem alterou o quê e quando, mas não permite desfazer a alteração.
A gravação é assíncrona e desacoplada da operação. O controlador dispara a inserção em goroutine após responder ao cliente, e a falha de gravação não aborta a operação nem é comunicada ao operador. A decisão prioriza a disponibilidade da operação sobre a completude da trilha, coerente com o RNF03, mas se afasta da recomendação de Kent e Souppaya (2006) quanto à confiabilidade do registro para fins de auditoria. Essa restrição é detalhada no Capítulo 6, e a gravação síncrona com confirmação de escrita figura como proposta de trabalho futuro.
3.8 SÍNTESE
O capítulo especificou a metodologia e os cinco modelos que orientam a implementação. Os requisitos vieram da taxonomia do Quadro 1, o que mantém rastreabilidade entre cada mecanismo e a classe de erro que o justifica.
Os modelos resultantes são: sessão opaca no servidor, revogável e expirada por inatividade, com senha derivada por bcrypt; RBAC híbrido com papel global e permissão de leitura ou escrita por par (usuário, zona), decidido no backend por requisição; normalização do RDATA como estratégia primária, com validação restrita aos casos sem conversão determinística; sincronização de reversos assimétrica entre criação otimista e remoção guardada por verificação de referência; e auditoria descritiva em coleção própria.
Três decisões foram registradas com suas limitações, por se afastarem de princípios que o próprio trabalho invoca: a escrita derivada em zona fora do escopo do operador, o melhor esforço das rotinas de reverso e o desacoplamento entre operação e registro de auditoria. As implicações dessas escolhas de projeto são discutidas no balanço crítico do Capítulo 6.
O Capítulo 4 documenta a construção do artefato a partir desses modelos.

REFERÊNCIAS
AKAMAI. Akamai summarizes service disruption (resolved). Akamai Blog, [s. l.], 22 jul. 2021. Disponível em: https://www.akamai.com/blog/news/akamai-summarizes-service-disruption-resolved. Acesso em: 14 set. 2026.

AKIWATE, G. et al. Unresolved issues: prevalence, persistence, and perils of lame delegations. In: ACM INTERNET MEASUREMENT CONFERENCE, 2020, [s. l.]. Proceedings [...]. New York: ACM, 2020. DOI: 10.1145/3419394.3423623.

ARENDS, R. et al. DNS security introduction and requirements. RFC 4033. [S. l.]: IETF, 2005. Disponível em: https://www.rfc-editor.org/rfc/rfc4033.

BANGOR, A.; KORTUM, P. T.; MILLER, J. T. An empirical evaluation of the System Usability Scale. International Journal of Human-Computer Interaction, v. 24, n. 6, p. 574-594, 2008.

BARR, D. Common DNS operational and configuration errors. RFC 1912. [S. l.]: IETF, 1996. Disponível em: https://www.rfc-editor.org/rfc/rfc1912.

BOJOVIĆ, P.; GAJIN, S. An approach to evaluation of common DNS misconfigurations. arXiv preprint arXiv:1711.05696, 2017. Disponível em: https://arxiv.org/abs/1711.05696.

BROOKE, J. SUS: a "quick and dirty" usability scale. In: JORDAN, P. W. et al. (ed.). Usability evaluation in industry. London: Taylor & Francis, 1996. p. 189-194.

CHUNG, T. et al. A longitudinal, end-to-end view of the DNSSEC ecosystem. In: USENIX SECURITY SYMPOSIUM, 26., 2017, Vancouver. Proceedings [...]. Berkeley: USENIX Association, 2017. p. 1307-1322.

CZ.NIC. Knot DNS documentation: operation. [S. l.], 2026. Disponível em: https://www.knot-dns.cz/docs/latest/html/operation.html.

DNSCONTROL. DNSControl. [S. l.], 2026. Disponível em: https://github.com/StackExchange/dnscontrol.

EIDNES, H.; DE GROOT, G.; VIXIE, P. Classless IN-ADDR.ARPA delegation. RFC 2317. [S. l.]: IETF, 1998. Disponível em: https://www.rfc-editor.org/rfc/rfc2317.
ELZ, R.; BUSH, R. Serial number arithmetic. RFC 1982. [S. l.]: IETF, 1996. Disponível em: https://www.rfc-editor.org/rfc/rfc1982.

GREGOR, S.; HEVNER, A. R. Positioning and presenting design science research for maximum impact. MIS Quarterly, v. 37, n. 2, p. 337-355, 2013.

GULBRANDSEN, A.; VIXIE, P.; ESIBOV, L. A DNS RR for specifying the location of services (DNS SRV). RFC 2782. [S. l.]: IETF, 2000. Disponível em: https://www.rfc-editor.org/rfc/rfc2782.

HALLAM-BAKER, P.; STRADLING, R.; HOFFMAN-ANDREWS, J. DNS Certification Authority Authorization (CAA) resource record. RFC 8659. [S. l.]: IETF, 2019. Disponível em: https://www.rfc-editor.org/rfc/rfc8659.

HARDAKER, W.; DUKHOVNI, V. Guidance for NSEC3 parameter settings. RFC 9276. [S. l.]: IETF, 2022. Disponível em: https://www.rfc-editor.org/rfc/rfc9276.

HEVNER, A. R. et al. Design science in information systems research. MIS Quarterly, v. 28, n. 1, p. 75-105, 2004.

HOFFMAN, P.; SCHLYTER, J. The DNS-Based Authentication of Named Entities (DANE) Transport Layer Security (TLS) protocol: TLSA. RFC 6698. [S. l.]: IETF, 2012. Disponível em: https://www.rfc-editor.org/rfc/rfc6698.

KENT, K.; SOUPPAYA, M. Guide to computer security log management. NIST Special Publication 800-92. Gaithersburg: National Institute of Standards and Technology, 2006.

KOLKMAN, O.; MEKKING, W.; GIEBEN, R. DNSSEC operational practices, version 2. RFC 6781. [S. l.]: IETF, 2012. Disponível em: https://www.rfc-editor.org/rfc/rfc6781.

KORCZYŃSKI, M. et al. Zone poisoning: the how and where of non-secure DNS dynamic updates. In: ACM INTERNET MEASUREMENT CONFERENCE, 2016, Santa Monica. Proceedings [...]. New York: ACM, 2016. p. 271-278. DOI: 10.1145/2987443.2987477.

LAURIE, B. et al. DNS Security (DNSSEC) hashed authenticated denial of existence. RFC 5155. [S. l.]: IETF, 2008. Disponível em: https://www.rfc-editor.org/rfc/rfc5155.

LIU, D.; HAO, S.; WANG, H. All your DNS records point to us: understanding the security threats of dangling DNS records. In: ACM SIGSAC CONFERENCE ON COMPUTER AND COMMUNICATIONS SECURITY, 2016, Vienna. Proceedings [...]. New York: ACM, 2016. p. 1414-1425. DOI: 10.1145/2976749.2978387.

MARCH, S. T.; SMITH, G. F. Design and natural science research on information technology. Decision Support Systems, v. 15, n. 4, p. 251-266, 1995.

META. More details about the October 4 outage. Engineering at Meta, [s. l.], 5 out. 2021. Disponível em: https://engineering.fb.com/2021/10/05/networking-traffic/outage-details/.

MICROSOFT. DNS Manager: manage resource records and zone security. Microsoft Learn, [s. l.], 2026. Disponível em: https://learn.microsoft.com/en-us/windows-server/networking/dns/.

MOCKAPETRIS, P. Domain names – concepts and facilities. RFC 1034. [S. l.]: IETF, 1987a. Disponível em: https://www.rfc-editor.org/rfc/rfc1034.

MOCKAPETRIS, P. Domain names – implementation and specification. RFC 1035. [S. l.]: IETF, 1987b. Disponível em: https://www.rfc-editor.org/rfc/rfc1035.

NETBOX DNS. NetBox DNS: a NetBox plugin for managing DNS data. [S. l.], 2026. Disponível em: https://github.com/peteeckel/netbox-plugin-dns.

NORMAN, D. A. The design of everyday things. Revised and expanded edition. New York: Basic Books, 2013.

OCTODNS. octoDNS. [S. l.], 2026. Disponível em: https://github.com/octodns/octodns

OPPENHEIMER, D.; GANAPATHI, A.; PATTERSON, D. A. Why do Internet services fail, and what can be done about it? In: USENIX SYMPOSIUM ON INTERNET TECHNOLOGIES AND SYSTEMS, 4., 2003, Seattle. Proceedings [...]. Berkeley: USENIX Association, 2003.

PAPPAS, V. et al. Impact of configuration errors on DNS robustness. ACM SIGCOMM Computer Communication Review, New York, v. 34, n. 4, p. 319-330, 2004. Disponível em: https://dl.acm.org/doi/10.1145/1030194.1015503.

PEFFERS, K. et al. A design science research methodology for information systems research. Journal of Management Information Systems, v. 24, n. 3, p. 45-77, 2007. DOI: 10.2753/MIS0742-1222240302.

PINGDOM. Sweden's Internet broken by DNS mistake. Pingdom Blog, 13 out. 2009. Disponível em: https://www.pingdom.com/blog/swedens-internet-broken-by-dns-mistake/.

POWERADMIN. Poweradmin. [S. l.], 2026. Disponível em: https://github.com/poweradmin/poweradmin.

POWERDNS. PowerDNS Authoritative Server documentation. [S. l.], 2026. Disponível em: https://doc.powerdns.com/authoritative/.

POWERDNS-ADMIN. PowerDNS-Admin. [S. l.], 2026. Disponível em: https://github.com/PowerDNS-Admin/PowerDNS-Admin.

REASON, J. Human error. Cambridge: Cambridge University Press, 1990.

SALTZER, J. H.; SCHROEDER, M. D. The protection of information in computer systems. Proceedings of the IEEE, v. 63, n. 9, p. 1278-1308, 1975.

SANDHU, R. S. et al. Role-based access control models. Computer, v. 29, n. 2, p. 38-47, 1996.

SAURO, J.; LEWIS, J. R. Quantifying the user experience: practical statistics for user research. 2. ed. Cambridge: Morgan Kaufmann, 2016.

SCHWARTZ, B.; BISHOP, M.; NYGREN, E. Service binding and parameter specification via the DNS (SVCB and HTTPS resource records). RFC 9460. [S. l.]: IETF, 2023. Disponível em: https://www.rfc-editor.org/rfc/rfc9460.

SIMON, H. A. The sciences of the artificial. 3. ed. Cambridge: MIT Press, 1996.

SOMMESE, R. et al. When parents and children disagree: diving into DNS delegation inconsistency. In: PASSIVE AND ACTIVE MEASUREMENT CONFERENCE, 21., 2020, Eugene. Proceedings [...]. Cham: Springer, 2020. p. 175-189. DOI: 10.1007/978-3-030-44081-7_11.

TECHNITIUM. Technitium DNS Server change log. [S. l.], 2026. Disponível em: https://github.com/TechnitiumSoftware/DnsServer/blob/master/CHANGELOG.md.

THOMSON, S. et al. DNS extensions to support IP version 6. RFC 3596. [S. l.]: IETF, 2003. Disponível em: https://www.rfc-editor.org/rfc/rfc3596.

VENABLE, J.; PRIES-HEJE, J.; BASKERVILLE, R. FEDS: a framework for evaluation in design science research. European Journal of Information Systems, v. 25, n. 1, p. 77-89, 2016. DOI: 10.1057/ejis.2014.36.

WOHLIN, C. et al. Experimentation in software engineering. Berlin: Springer, 2012.

YIN, R. K. Case study research and applications: design and methods. 6. ed. Thousand Oaks: Sage, 2018.

ZONEMASTER. Zonemaster master test plan. [S. l.], 2026. Disponível em: https://doc.zonemaster.net/latest/specifications/tests/MasterTestPlan.html.
