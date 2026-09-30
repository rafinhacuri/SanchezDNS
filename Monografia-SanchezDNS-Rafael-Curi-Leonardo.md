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
CONCEITO: \***\*\*\*\*\***\_\_\_\_\***\*\*\*\*\***

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
O Sistema de Nomes de Domínio (Domain Name System – DNS) é uma infraestrutura crítica para o funcionamento da Internet. Sua especificação original priorizou a escalabilidade hierárquica e a redundância distribuída, delegando a integridade das zonas autoritativas à manutenção manual por operadores humanos. Estudos de medição em larga escala evidenciam que erros de configuração nessas zonas são recorrentes e comprometem a robustez, a latência e a disponibilidade dos serviços dependentes. Este trabalho apresenta o desenvolvimento e a avaliação do SanchezDNS, uma plataforma web que atua como camada de operação sobre a API REST do PowerDNS Authoritative Server, mantido como fonte única da verdade dos dados de zona. O artefato busca reduzir falhas humanas por meio de três mecanismos: (i) automação de registros deriváveis, como a criação e remoção de registros reversos (PTR) a partir de registros A e AAAA, a normalização de valores por tipo de registro, a delegação do incremento do serial SOA ao servidor autoritativo e a ativação automática de DNSSEC; (ii) controle de acesso combinando papéis globais e permissões de leitura e escrita atribuídas por zona; e (iii) trilha de auditoria centralizada das operações sobre zonas, registros e usuários. A pesquisa segue a metodologia Design Science Research (DSR), estruturada conforme o modelo de processo de Peffers et al. (2007), com avaliação conduzida segundo o Framework for Evaluation in Design Science (FEDS). A avaliação verifica, por casos de teste negativos e com o Zonemaster como oráculo independente, se os mecanismos impedem as classes de erro que se propõem a eliminar, compara as oportunidades de erro da operação pela plataforma com as da edição manual e demonstra o artefato implantado em ambiente real.

Palavras-chave: DNS autoritativo; erro de configuração; PowerDNS; controle de acesso por zona; Design Science Research.
ABSTRACT
The Domain Name System (DNS) is a critical infrastructure for the operation of the Internet. Its original specification prioritized hierarchical scalability and distributed redundancy, leaving the integrity of authoritative zones to manual maintenance by human operators. Large-scale measurement studies show that configuration errors in these zones are recurrent and compromise the robustness, latency and availability of dependent services. This work presents the development and evaluation of SanchezDNS, a web platform that acts as an operation layer over the REST API of the PowerDNS Authoritative Server, which is kept as the single source of truth for zone data. The artifact seeks to reduce human error through three mechanisms: (i) automation of derivable records, such as the creation and removal of reverse (PTR) records from A and AAAA records, normalization of values by record type, delegation of SOA serial increments to the authoritative server and automatic activation of DNSSEC; (ii) access control combining global roles with read and write permissions assigned per zone; and (iii) a centralized audit trail of operations on zones, records and users. The research follows the Design Science Research (DSR) methodology, structured according to the process model of Peffers et al. (2007), with evaluation conducted under the Framework for Evaluation in Design Science (FEDS). The evaluation checks, through negative test cases with Zonemaster as an independent oracle, whether the mechanisms prevent the error classes they target, compares the error opportunities of operating through the platform with those of manual editing, and demonstrates the artifact deployed in a real environment.

Keywords: authoritative DNS; misconfiguration; PowerDNS; per-zone access control; Design Science Research.
LISTA DE ILUSTRAÇÕES
Figura 1 – Hierarquia do espaço de nomes do DNS: raiz, domínios de topo, zonas e delegação ...................................................................................................................................................... 26
Figura 2 – Sincronização de registros PTR: criação otimista e remoção conservadora .............. 69
Figura 3 – Arquitetura da plataforma SanchezDNS e fluxo de uma requisição .......................... 74
LISTA DE QUADROS
Quadro 1 – Taxonomia de erros de configuração em zonas autoritativas, evidências, impactos e mecanismos de mitigação ............................................................................................................ 34
Quadro 2 – Comparação entre ferramentas de gestão de zonas autoritativas .............................. 44
Quadro 3 – Matriz de avaliação: questões, construtos, métricas, instrumentos e critérios de êxito ...................................................................................................................................................... 49
Quadro 4 – Atividades do DSRM e sua realização nesta pesquisa .............................................. 52
Quadro 5 – Requisitos funcionais e rastreabilidade com a taxonomia de erros .......................... 54
Quadro 6 – Parâmetros do modelo de sessão e autenticação ....................................................... 58
Quadro 7 – Matriz de autorização e controle de rotas por nível de acesso .................................. 61
Quadro 8 – Regras de normalização sintática aplicadas no backend ........................................... 64
Quadro 9 – Operações e códigos de ação registrados na trilha de auditoria ............................... 71
Quadro 10 – Pilha tecnológica da plataforma .............................................................................. 75
Quadro 11 – Rotas do PowerDNS consumidas pela plataforma ................................................. 77
LISTA DE ABREVIATURAS E SIGLAS
ACL Access Control List
AFNIC Association Française pour le Nommage Internet en Coopération
API Application Programming Interface
AWS Amazon Web Services
AXFR Authoritative Zone Transfer
BGP Border Gateway Protocol
BIND Berkeley Internet Name Domain
CAA Certification Authority Authorization
ccTLD Country Code Top-Level Domain
CLI Command-Line Interface
CNAME Canonical Name
CSS Cascading Style Sheets
DANE DNS-based Authentication of Named Entities
DKIM DomainKeys Identified Mail
DNS Domain Name System
DNSKEY DNS Public Key
DNSSEC Domain Name System Security Extensions
DS Delegation Signer
DSR Design Science Research
DSRM Design Science Research Methodology
ECDSA Elliptic Curve Digital Signature Algorithm
FEDS Framework for Evaluation in Design Science
FQDN Fully Qualified Domain Name
HTTP Hypertext Transfer Protocol
HTTPS Hypertext Transfer Protocol Secure
IaC Infrastructure as Code
IAM Identity and Access Management
IEEE Institute of Electrical and Electronics Engineers
IP Internet Protocol
IPAM IP Address Management
IPv4 Internet Protocol version 4
IPv6 Internet Protocol version 6
ISC Internet Systems Consortium
IXFR Incremental Zone Transfer
JSON JavaScript Object Notation
JWT JSON Web Token
KSK Key Signing Key
LDAP Lightweight Directory Access Protocol
MX Mail Exchanger
NIST National Institute of Standards and Technology
NREN National Research and Education Network
NS Name Server
NSEC3 Next Secure, versão 3
NSEC3PARAM NSEC3 Parameters
PTR Pointer Record
QP Questão de Pesquisa
RBAC Role-Based Access Control
RDATA Resource Data
REST Representational State Transfer
RF Requisito Funcional
RFC Request for Comments
RIPE Réseaux IP Européens
RNF Requisito Não Funcional
RR Resource Record
RRSIG Resource Record Signature
RSA Rivest-Shamir-Adleman
S3 Simple Storage Service
SID Session Identifier
SMTP Simple Mail Transfer Protocol
SO Sistema Operacional
SOA Start of Authority
SPF Sender Policy Framework
SRV Service Record
SUS System Usability Scale
TI Tecnologia da Informação
TLS Transport Layer Security
TLSA TLS Authentication Record
TTL Time to Live
TXT Text Record
UI User Interface
ZSK Zone Signing Key
SUMÁRIO
1 INTRODUÇÃO .............................................................................................................. 15
1.1 CONTEXTUALIZAÇÃO E PROBLEMATIZAÇÃO DO DOMÍNIO ............................ 15
1.2 DEFASAGENS DO MERCADO E BRECHAS DAS SOLUÇÕES EXISTENTES ...... 18
1.3 PROBLEMA DE PESQUISA E QUESTÕES NORTEADORAS .................................. 19
1.4 OBJETIVOS DO TRABALHO ....................................................................................... 21
1.4.1 Objetivo Geral .................................................................................................................. 21
1.4.2 Objetivos Específicos ....................................................................................................... 21
1.5 DELIMITAÇÃO DO ESCOPO ....................................................................................... 22
1.6 ESTRUTURA DA MONOGRAFIA ............................................................................... 23
2 FUNDAMENTAÇÃO TEÓRICA E TRABALHOS CORRELATOS ....................... 24
2.1 ARQUITETURA DO SISTEMA DE NOMES DE DOMÍNIO ...................................... 25
2.1.1 Espaço de nomes, zonas e delegação ............................................................................... 25
2.1.2 Registros de recurso ......................................................................................................... 26
2.1.3 O registro SOA e a replicação de zonas ........................................................................... 27
2.1.4 Resolução reversa ............................................................................................................ 28
2.1.5 Extensões de segurança (DNSSEC) ................................................................................ 28
2.1.6 O PowerDNS Authoritative Server .................................................................................. 29
2.2 ERROS DE CONFIGURAÇÃO EM ZONAS AUTORITATIVAS ................................. 30
2.2.1 O erro humano na operação de infraestrutura .................................................................. 30
2.2.2 Estudos empíricos de medição ......................................................................................... 30
2.2.3 Impactos operacionais e de segurança ............................................................................. 32
2.2.4 Síntese: taxonomia de erros, evidências e mecanismos de mitigação ............................. 33
2.3 CONTROLE DE ACESSO E AUDITORIA ................................................................... 35
2.4 DESIGN SCIENCE RESEARCH ................................................................................... 36
2.4.1 Fundamentos metodológicos ........................................................................................... 36
2.4.2 Processo DSRM ............................................................................................................... 37
2.4.3 Avaliação: Framework FEDS e Escala SUS .................................................................... 37
2.5 TRABALHOS CORRELATOS ....................................................................................... 38
2.5.1 BIND 9 e manipulação direta de arquivos ....................................................................... 38
2.5.2 Interfaces web de software livre (PowerDNS-Admin) .................................................... 38
2.5.3 Technitium DNS Server ................................................................................................... 39
2.5.4 Microsoft DNS Server e o console DNS Manager .......................................................... 40
2.5.5 NetBox DNS e a abordagem de fonte externa da verdade .............................................. 41
2.5.6 Knot DNS e as interfaces de controle programáticas ...................................................... 41
2.5.7 Zonemaster e a validação independente de delegações ................................................... 42
2.5.8 DNS como código (IaC) e plataformas em nuvem .......................................................... 43
2.5.9 Quadro comparativo e delimitação da contribuição ........................................................ 43
2.6 ESTRATÉGIA DE AVALIAÇÃO DO ARTEFATO ........................................................ 45
2.6.1 Escolha da estratégia de avaliação ................................................................................... 46
2.6.2 Episódios de avaliação ..................................................................................................... 46
2.6.3 Implantação em ambiente real ......................................................................................... 48
2.6.4 Matriz de avaliação, ameaças à validade e limitações ..................................................... 49
2.7 CONSIDERAÇÕES DO CAPÍTULO ............................................................................. 51
3 METODOLOGIA E MODELAGEM DO ARTEFATO ............................................. 51
3.1 CLASSIFICAÇÃO DA PESQUISA ................................................................................ 52
3.2 APLICAÇÃO DO MODELO DSRM .............................................................................. 52
3.3 REQUISITOS .................................................................................................................. 54
3.4 AUTENTICAÇÃO E CONTROLE DE ACESSO .......................................................... 57
3.4.1 Autenticação e ciclo de vida da sessão ............................................................................ 57
3.4.2 Cadastro, aprovação e inicialização ................................................................................. 59
3.4.3 Estrutura das permissões .................................................................................................. 59
3.4.4 Aplicação nas rotas .......................................................................................................... 60
3.4.5 Escrita derivada em zona de outro escopo ....................................................................... 63
3.5 NORMALIZAÇÃO SINTÁTICA .................................................................................... 63
3.5.1 Estratégia ......................................................................................................................... 63
3.5.2 Regras por tipo ................................................................................................................. 64
3.5.3 Devolução do valor normalizado ..................................................................................... 66
3.6 SINCRONIZAÇÃO DE ZONAS REVERSAS ............................................................... 66
3.6.1 Derivação do nome e da zona .......................................................................................... 66
3.6.2 Criação e alteração ........................................................................................................... 67
3.6.3 Guarda de referência ........................................................................................................ 67
3.6.4 Semântica de melhor esforço ........................................................................................... 70
3.7 AUDITORIA ................................................................................................................... 70
3.8 SÍNTESE ......................................................................................................................... 72
4 DESENVOLVIMENTO DA PLATAFORMA SANCHEZDNS ................................ 73
4.1 ARQUITETURA GERAL E PILHA TECNOLÓGICA ................................................. 73
4.1.1 Arquitetura desacoplada .................................................................................................. 73
4.1.2 Escolhas de linguagem e de biblioteca ............................................................................ 74
4.1.3 Persistência híbrida .......................................................................................................... 76
4.2 INTEGRAÇÃO COM O POWERDNS AUTHORITATIVE SERVER .......................... 76
4.2.1 Cliente e autenticação da interface de programação ........................................................ 76
4.2.2 Criação de zona e delegação do serial ............................................................................. 78
4.2.3 Assinatura e verificação da cadeia de confiança .............................................................. 78
4.3 AUTENTICAÇÃO E CONTROLE DE ACESSO .......................................................... 79
4.3.1 Ciclo de vida da sessão e derivação da senha .................................................................. 79
4.3.2 Middlewares de sessão e de administração ...................................................................... 79
4.3.3 Resolução da permissão por zona .................................................................................... 80
4.4 NORMALIZAÇÃO E INTERFACE PREVENTIVA ...................................................... 80
4.4.1 Decomposição dos tipos compostos na interface ............................................................ 80
4.4.2 Normalização no backend ................................................................................................ 81
4.4.3 Segundo retorno do ciclo: exibir o valor publicado ......................................................... 81
4.5 SINCRONIZAÇÃO DO CICLO DE VIDA DOS REGISTROS REVERSOS ............... 82
4.5.1 Derivação do nome e escolha da zona ............................................................................. 82
4.5.2 Primeiro retorno do ciclo: guarda de referência na remoção ........................................... 82
4.5.3 Melhor esforço, tempo limite e operação em lote ........................................................... 83
4.6 TRILHA DE AUDITORIA ............................................................................................. 83
4.6.1 Gravação assíncrona ........................................................................................................ 83
4.6.2 Consulta e filtragem ......................................................................................................... 84
4.7 EMPACOTAMENTO E IMPLANTAÇÃO ..................................................................... 84
REFERÊNCIAS ......................................................................................................................... 90
1 INTRODUÇÃO
Esta introdução apresenta o contexto, as justificativas e a formulação metodológica do SanchezDNS. O objetivo é investigar e propor diretrizes para reduzir falhas operacionais em servidores de nomes autoritativos, usando uma interface que abstrai a sintaxe das zonas, permissões atribuídas por zona e automação da validação e da sincronização de registros.
1.1 CONTEXTUALIZAÇÃO E PROBLEMATIZAÇÃO DO DOMÍNIO
O Sistema de Nomes de Domínio (Domain Name System – DNS) é uma infraestrutura hierárquica e distribuída, essencial para traduzir identificadores mnemônicos em endereços IP na Internet. A especificação original do protocolo (MOCKAPETRIS, 1987a, 1987b) priorizou a escalabilidade hierárquica e a resiliência das consultas por meio de redundância distribuída, e deixou a integridade das zonas autoritativas a cargo da edição e da manutenção humana dos arquivos de configuração.
Estudos empíricos de medição de redes mostram que inconsistências operacionais e erros de configuração na administração de zonas autoritativas são bastante difundidos. Bojović e Gajin (2017) analisaram cerca de 80 mil domínios distribuídos entre redes acadêmicas internacionais (NRENs), o ccTLD sérvio (.rs) e os domínios mais populares da Internet, e encontraram taxas relevantes de falhas de configuração, principalmente nos segmentos acadêmico e nacional. Entre os problemas diagnosticados estão a presença de servidores não autoritativos declarados na zona pai (lame delegation), em quase 5% dos domínios acadêmicos avaliados, e a ausência completa de registros de ponteiro reverso (PTR) para os servidores de nomes em cerca de 40% a 45% dos domínios de cada grupo analisado.
Pappas et al. (2004) já haviam encontrado lame delegations em aproximadamente 15% das zonas analisadas e dependências circulares em 2%. Segundo os autores, essas inconsistências reduzem a redundância planejada e elevam a latência de resolução das consultas em até uma ordem de grandeza.
Outro ponto sensível da manutenção manual é o número serial (SOA Serial), que coordena a sincronização entre o servidor primário e os secundários. Quando o operador esquece de incrementá-lo, ou o altera violando as regras de aritmética de números de série (ELZ; BUSH, 1996), ocorre dessincronização (zone drift) e os secundários continuam respondendo com uma versão desatualizada da zona.
A falta de validação antes da publicação também pode ter efeitos graves. Em 13 de outubro de 2009, um script automatizado de manutenção do ccTLD da Suécia (.se) omitiu o ponto final nos registros de delegação dos servidores de nomes. A zona corrompida foi propagada imediatamente para toda a rede secundária, e pouco mais de 900 mil domínios sob .se ficaram sem resolução por cerca de uma hora, com efeitos prolongados pelo cache dos resolvedores (PINGDOM, 2009). Era um erro sintático, do tipo que uma validação feita antes da publicação teria barrado.
Uma segunda classe de falha, além dos erros de sintaxe, decorre da divergência entre a zona pai e a zona filha. Sommese et al. (2020) compararam os conjuntos de registros NS publicados por pais e filhos nos domínios de segundo nível de .com, .net e .org e encontraram inconsistência em mais de treze milhões de domínios, entre eles ao menos cinquenta mil da lista dos um milhão mais acessados. Como o resolvedor pode usar legitimamente qualquer um dos dois conjuntos, a divergência provoca indisponibilidade intermitente. O diagnóstico é difícil porque a consulta de verificação feita pelo operador costuma receber resposta correta.
Uma terceira classe aparece no fim do ciclo de vida dos registros. Registros que continuam publicados depois do desprovisionamento do recurso para o qual apontam, os chamados registros pendentes (dangling records), podem ser explorados por terceiros. Liu, Hao e Wang (2016) identificaram 467 registros pendentes exploráveis em 277 domínios entre os dez mil mais acessados e em 52 zonas do domínio .edu, e demonstraram que um terceiro pode assumir o controle do subdomínio correspondente e até obter certificados TLS válidos em nome da organização titular. O risco é maior nos registros reversos: eles ficam em outra zona, e sua exclusão não acompanha automaticamente a remoção do registro A ou AAAA que lhes deu origem.
Uma quarta classe envolve a operação das extensões de segurança. Chung et al. (2017), em uma série histórica de vinte e um meses sobre todos os domínios assinados em .com, .net e .org, verificaram que mais de 30% das zonas assinadas não tinham o registro DS correspondente publicado na zona pai, o que rompe a cadeia de confiança e deixa a assinatura sem efeito, e que cerca de 70% nunca haviam feito rotação de chaves. Isso sugere que a dificuldade de adotar o DNSSEC está mais nas etapas operacionais do que no mecanismo criptográfico. Korczyński et al. (2016), por sua vez, encontraram em uma amostra aleatória de 2,9 milhões de domínios 1.877 domínios (0,065%) cujos servidores autoritativos aceitavam atualizações dinâmicas não autenticadas, o que permite a qualquer agente da Internet inserir registros arbitrários. O percentual é pequeno, mas nesses domínios o controle de acesso à zona simplesmente deixa de existir.
Como quase toda comunicação na Internet depende do DNS, um erro localizado pode derrubar serviços muito além do domínio afetado. Em 22 de julho de 2021, uma atualização de configuração no serviço Edge DNS da Akamai deixou serviços de larga escala inacessíveis por pouco mais de uma hora, até a reversão da alteração (AKAMAI, 2021). Em 4 de outubro do mesmo ano, um comando de manutenção tirou de operação a rede de backbone da Meta e provocou a retirada dos anúncios BGP dos prefixos em que estavam seus servidores autoritativos, deixando Facebook, Instagram e WhatsApp fora do ar por mais de sete horas (META, 2021). Nos dois casos a causa imediata foi uma alteração de configuração feita por operadores, sem falha de equipamento ou ação adversária, e os efeitos se propagaram antes que alguém pudesse detectá-los, como no incidente do ccTLD sueco.
Essas evidências sugerem que boa parte dos erros no DNS autoritativo vem da distância entre o modelo mental do operador e a representação textual que ele precisa editar, mais do que da complexidade do protocolo. Essa distância pesa mais porque nada verifica a alteração no momento da edição, porque o cache torna a correção lenta e porque as ferramentas, no máximo, detectam certos enganos depois de publicados, quando poderiam impedi-los. A seção 2.2 detalha essas classes de erro, suas evidências e seus impactos, e as organiza em uma taxonomia usada no projeto e na avaliação do artefato.
1.2 DEFASAGENS DO MERCADO E BRECHAS DAS SOLUÇÕES EXISTENTES
As ferramentas atuais de administração de servidores autoritativos podem ser agrupadas em cinco abordagens operacionais:
Edição manual de arquivos de zona via linha de comando (CLI): softwares tradicionais como o BIND 9 exigem a manipulação direta de arquivos planos no formato de arquivo mestre (master file format). Essa abordagem é muito propensa a erros de digitação, não oferece validação estrutural no momento da edição, requer recargas manuais do serviço e depende excessivamente da experiência individual do operador.
Painéis de controle de hospedagem (ex.: cPanel, Plesk): desenvolvidos principalmente para o provisionamento de hospedagem web compartilhada, isolam os dados por conta de cliente, mas não permitem subdividir os privilégios operacionais de forma granular (por exemplo, separar operadores de leitura e de escrita) dentro de uma mesma equipe que gerencia zonas compartilhadas.
Painéis web dedicados de código aberto (ex.: PowerDNS-Admin): embora consolidados no ecossistema de software livre e dotados de controle de acesso e registro de eventos, esses sistemas se apoiam em arquiteturas monolíticas tradicionais com renderização no servidor (server-side rendering). Na prática, a criação automática de registros PTR depende de configuração adicional, o provisionamento de DNSSEC fica separado da criação da zona e não foram identificadas rotinas assistidas para localizar e remover registros reversos órfãos.
Painéis proprietários de provedores em nuvem (Cloud DNS): plataformas gerenciadas como AWS Route 53, Cloudflare e Google Cloud DNS oferecem usabilidade moderna e alta resiliência. Por outro lado, trazem aprisionamento tecnológico (vendor lock-in), custos contínuos atrelados ao volume de consultas e não podem ser adotadas em redes corporativas ou governamentais restritas (on-premises) que exigem custódia estrita sobre os servidores autoritativos.
Fontes externas da verdade e validadores independentes (ex.: NetBox DNS, Zonemaster): a primeira abordagem transfere a autoridade sobre os dados para um sistema de inventário, que deriva registros diretos e reversos do endereçamento cadastrado; a segunda submete a delegação já publicada a uma suíte formal de casos de teste. As duas contribuem para a qualidade das zonas, mas nenhuma atua no momento da edição. A fonte externa acrescenta um componente de sincronização e abre espaço para divergências entre o inventário e o que de fato responde às consultas; o validador, por sua vez, opera a posteriori, quando o erro já foi publicado e propagado aos resolvedores.
Vários dos mecanismos pretendidos já existem em ferramentas maduras. O que falta é uma camada de operação desacoplada do motor de resolução, aplicável a infraestruturas locais já padronizadas sobre o PowerDNS Authoritative Server, que reúna em um mesmo fluxo a gestão do ciclo de vida completo dos registros reversos (incluindo a remoção encadeada e a triagem de reversos órfãos), a normalização sintática por tipo de registro, o provisionamento de DNSSEC na criação da zona, a distinção entre privilégios de leitura e de escrita em cada zona e uma trilha de auditoria na própria camada de aplicação. O Capítulo 2 examina essas ferramentas em detalhe, incluindo as de DNS como código, compara suas funcionalidades em um quadro e discute um incidente de segurança corrigido recentemente no Technitium DNS Server que mostra como é difícil combinar a automação de registros deriváveis com o controle de acesso por zona.
1.3 PROBLEMA DE PESQUISA E QUESTÕES NORTEADORAS
Com base nas seções anteriores, o problema de pesquisa pode ser enunciado assim: a operação de zonas DNS autoritativas depende de tarefas manuais, repetitivas e mutuamente dependentes cuja correção não é verificada no instante da edição, de modo que deslizes de execução e enganos de planejamento são publicados e propagados aos resolvedores antes de qualquer oportunidade de detecção, e permanecem ativos mesmo após a correção da origem. A literatura considera o erro humano inevitável (REASON, 1990; NORMAN, 2013). O problema, portanto, é a falta de um meio de operação que torne esse erro impossível nas classes em que isso é tecnicamente viável.
Desse enunciado deriva a questão de pesquisa principal:
QP. Em que medida uma camada web de operação sobre um servidor autoritativo, que automatize registros deriváveis, restrinja privilégios ao nível da zona e registre integralmente as alterações efetuadas, reduz as oportunidades de erro de configuração e o esforço de operação em comparação com a edição manual assistida por linha de comando?
A questão principal reúne asserções de naturezas distintas: algumas podem ser verificadas por argumentação técnica, outras dependem de evidência empírica com operadores humanos. Por isso, ela se desdobra em quatro questões secundárias, cada uma vinculada a objetivos específicos (seção 1.4.2) e a episódios do desenho de avaliação (seção 2.6.2):
QP1. Caracterização. Quais classes de erro de configuração em zonas autoritativas são recorrentes, quais são seus impactos documentados e quais delas podem ser eliminadas por construção (isto é, tornadas impossíveis pelo próprio meio de edição), em oposição às que só admitem detecção posterior?
QP2. Mecanismos. Que mecanismos de projeto (automação de dados deriváveis, normalização sintática por tipo de registro, delegação do incremento do serial ao motor autoritativo e provisionamento de DNSSEC no fluxo de criação da zona) são necessários e suficientes para eliminar as classes identificadas em QP1? E que novos riscos esses mecanismos introduzem na operação ao automatizar decisões antes tomadas pelo operador?
QP3. Autoridade e automação. Que modelo de controle de acesso concilia a granularidade por zona com a automação de registros que ficam em zona distinta daquela que o operador edita? É o caso do registro PTR, derivado de um registro A ou AAAA na zona direta, mas gravado na zona reversa, muitas vezes sob responsabilidade administrativa de outra equipe.
QP4. Efeito observado. Nas tarefas típicas de administração de zonas, a operação mediada pela plataforma retira do operador as oportunidades de erro e os passos manuais que a edição manual exige, e o artefato funciona como especificado quando implantado sobre um servidor autoritativo em operação?
Essa divisão define o desenho metodológico. As questões QP1 a QP3 são respondidas por revisão sistemática da literatura, argumentação de projeto e verificação técnica feita com oráculos independentes do artefato. A QP4 exige evidência obtida com o artefato em uso, e por isso a avaliação descrita na seção 2.6 combina episódios de laboratório com uma implantação em ambiente real. Ter a plataforma funcionando é pré-requisito para responder a qualquer uma das questões, mas isso sozinho não mostra que os problemas enunciados foram mitigados.
1.4 OBJETIVOS DO TRABALHO
1.4.1 Objetivo Geral
Projetar, implementar e avaliar um artefato de software, denominado SanchezDNS, fundamentado no paradigma Design Science Research (DSR), que atue como camada de mediação sobre o PowerDNS Authoritative Server e reduza a ocorrência de inconsistências operacionais por meio de automação e sincronização de registros reversos, validação estrutural de entrada, controle de acesso híbrido com permissões por zona e registro de eventos em trilha de auditoria centralizada.
1.4.2 Objetivos Específicos
Para responder às questões formuladas, estabelecem-se os objetivos específicos a seguir, indicando-se, ao final de cada um, a questão de pesquisa a que se vincula:
Mapear os padrões de erros operacionais e inconsistências mais frequentes na edição de zonas autoritativas a partir da literatura técnica e empírica de medição. (QP1)
Formular diretrizes de interface voltadas à mitigação de falhas humanas na digitação e à automação de dados derivados, integrando a sincronização automática de registros PTR (IPv4 e IPv6) na criação/edição de registros A/AAAA, ferramentas para detecção e remoção de registros reversos órfãos e delegação transparente do número serial ao motor do PowerDNS (SOA-EDIT-API). (QP2)
Automatizar o provisionamento inicial de segurança criptográfica em zonas gerenciadas, efetuando a geração da chave de assinatura de chaves (Key Signing Key – KSK) e a configuração dos parâmetros de negação autenticada de existência (NSEC3PARAM) no fluxo de criação do domínio. (QP2)
Implementar o artefato sob arquitetura desacoplada, dividindo o frontend reativo do backend intermediador e mantendo a API REST do PowerDNS Authoritative Server como fonte exclusiva da verdade para os dados de zona. (QP2)
Estruturar um modelo híbrido de controle de acesso que conjugue o perfil administrativo global à concessão isolada de privilégios de leitura e de escrita no nível individual de cada zona. (QP3)
Registrar em trilha de auditoria centralizada as operações de criação, alteração e remoção de zonas, registros e usuários, com identificação do autor e do instante de execução de cada alteração. (QP4)
Avaliar a plataforma segundo o framework FEDS, verificando por casos de teste negativos e oráculo independente se os mecanismos impedem as classes de erro elimináveis, comparando as oportunidades de erro da operação pela plataforma com as da edição manual e demonstrando o artefato implantado em ambiente real. (QP2, QP3, QP4)
1.5 DELIMITAÇÃO DO ESCOPO
Integram o escopo desta pesquisa a operação de zonas autoritativas diretas e reversas mantidas em instâncias do PowerDNS Authoritative Server acessíveis por sua API REST; a automação de registros deriváveis dessas zonas; o controle de acesso de operadores no nível da zona; e o registro de auditoria das operações realizadas pela plataforma.
Permanecem expressamente fora do escopo:
a resolução recursiva e o comportamento de resolvedores, tratados apenas como fonte de restrições de projeto, em especial o efeito do armazenamento em cache sobre a reversibilidade das alterações;
a resiliência da infraestrutura a ataques de negação de serviço, o dimensionamento de nuvens anycast e o desempenho do motor de resolução, atributos que dependem do servidor autoritativo subjacente e não da camada de operação proposta;
a interação com registrários e com a zona pai, incluindo a publicação do registro DS e a manutenção da consistência entre os conjuntos de NS pai e filho, operações conduzidas por interfaces de terceiros; a plataforma limita-se a exibir ao operador os dados necessários e a sinalizar divergências detectadas;
a substituição de ferramentas de infraestrutura como código, cujo público e cujo modelo operacional são distintos, conforme discutido na seção 2.5.8;
a implementação de um motor DNS próprio, uma vez que o artefato pressupõe e preserva o servidor autoritativo como fonte única da verdade.
No plano da avaliação, fica fora do escopo a medição com operadores: os episódios ocorrem em laboratório e em uma implantação real. As consequências dessa delimitação sobre a validade externa dos resultados são discutidas na seção 2.6.4.
1.6 ESTRUTURA DA MONOGRAFIA
Esta monografia está organizada nos seguintes capítulos:
Capítulo 2. Fundamentação Teórica e Trabalhos Correlatos: Revisa a arquitetura do protocolo DNS, os estudos empíricos sobre inconsistências de configuração em zonas autoritativas e a fundamentação epistemológica em Design Science Research. Estabelece, ainda, a estratégia de avaliação do artefato, com os episódios, as métricas e os critérios de êxito definidos previamente à construção.
Capítulo 3. Metodologia e Modelagem do Artefato: Detalha as fases do método DSRM empregadas na condução da pesquisa, o levantamento de requisitos do SanchezDNS, a modelagem de autenticação e de controle de acesso, as diretrizes de normalização sintática, as estratégias de automação do ciclo de vida de registros e a modelagem da trilha de auditoria.
Capítulo 4. Desenvolvimento da Plataforma SanchezDNS: Apresenta os aspectos de engenharia e implementação do protótipo, evidenciando as decisões de arquitetura de software, a integração com o PowerDNS e os algoritmos de tratamento de zonas reversas e auditoria.
Capítulo 5. Avaliação e Resultados: Documenta a execução dos episódios de avaliação sob as diretrizes do framework FEDS: os casos de teste negativos verificados pelo Zonemaster, a comparação com a operação manual e a implantação em ambiente real. Os dados são analisados à luz dos critérios de êxito fixados no Quadro 3.
Capítulo 6. Conclusões: Sintetiza as contribuições práticas e teóricas alcançadas, discute as limitações atuais da ferramenta e do artefato, e propõe diretrizes para investigações futuras.

2 FUNDAMENTAÇÃO TEÓRICA E TRABALHOS CORRELATOS
Este capítulo reúne a base conceitual necessária à fundamentação e à avaliação do SanchezDNS. A seção 2.1 descreve a arquitetura do Sistema de Nomes de Domínio (DNS), enfatizando os componentes manipulados pelo artefato. A seção 2.2 discute o erro humano na operação de infraestrutura, examina os estudos empíricos sobre falhas de configuração em zonas autoritativas e consolida os achados em uma taxonomia de classes de erro, evidências e impactos. A seção 2.3 aborda os princípios de controle de acesso e auditoria. A seção 2.4 expõe o paradigma Design Science Research (DSR), que orienta metodologicamente esta pesquisa. A seção 2.5 analisa os trabalhos correlatos e delimita a contribuição da plataforma proposta. A seção 2.6 estabelece a estratégia de avaliação do artefato, especificando episódios, métricas, critérios de êxito e ameaças à validade. Por fim, a seção 2.7 sintetiza o capítulo.

2.1 ARQUITETURA DO SISTEMA DE NOMES DE DOMÍNIO
2.1.1 Espaço de nomes, zonas e delegação
O DNS organiza os nomes sob uma árvore hierárquica invertida cuja raiz é representada por um rótulo vazio. O nome de domínio plenamente qualificado (Fully Qualified Domain Name – FQDN) é definido pela sequência de rótulos de um nó até a raiz, separados por pontos (MOCKAPETRIS, 1987a). Na representação textual, o ponto final explicita a raiz; sua omissão faz com que o nome seja interpretado como relativo a uma origem, detalhe sintático que constitui fonte frequente de erros operacionais (BARR, 1996).
A administração do espaço de nomes é distribuída por meio de zonas, que são porções contíguas da árvore sob uma mesma autoridade administrativa. A autoridade sobre uma subárvore é repassada por delegação, materializada pela inserção, na zona pai, de registros NS que apontam para os servidores da zona filha. Quando esses servidores pertencem ao próprio domínio delegado, a zona pai precisa conter registros de endereço auxiliares (glue records) para evitar dependências circulares na resolução (MOCKAPETRIS, 1987a). A Figura 1 representa essa hierarquia, com a raiz, os domínios de topo e a delegação de uma zona filha.
A infraestrutura do protocolo baseia-se em dois papéis funcionais fundamentais:
Servidores autoritativos: armazenam as definições originais das zonas e respondem às consultas com autoridade sobre os dados.
Resolvedores recursivos: consultam a hierarquia a partir da raiz em nome das aplicações clientes e retêm as respostas em cache pelo intervalo estipulado no Time to Live (TTL) de cada registro (MOCKAPETRIS, 1987a).
Devido ao cache nos resolvedores, qualquer erro publicado em zona autoritativa propaga-se e permanece ativo até a expiração do TTL, mesmo após ter sido corrigido no servidor de origem, a exemplo do incidente com o ccTLD sueco (.se) documentado por Pingdom (2009).
Figura 1 – Hierarquia do espaço de nomes do DNS: raiz, domínios de topo, zonas e delegação

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
Essa sintaxe fragmentada em 32 níveis torna o preenchimento manual de reversos IPv6 propenso a erros de digitação. Além disso, zonas diretas e reversas são administrativamente independentes, e nada no protocolo obriga que correspondam entre si. Essa paridade, porém, é cobrada por serviços de rede essenciais (como a mitigação de spam em servidores SMTP), o que motivou a recomendação expressa da RFC 1912 de manter a concordância entre apontamentos A/AAAA e seus respectivos PTRs (BARR, 1996).
2.1.5 Extensões de segurança (DNSSEC)
As extensões de segurança do DNS (DNSSEC) incorporam integridade e autenticidade de origem às respostas por meio de assinaturas digitais, sem prover confidencialidade (ARENDS et al., 2005). O protocolo emprega os seguintes registros:
RRSIG: armazena a assinatura criptográfica de um determinado conjunto de registros (RRset).
DNSKEY: publica as chaves públicas da zona, convencionalmente divididas em chave de assinatura de chaves (Key Signing Key – KSK) e chave de assinatura de zona (Zone Signing Key – ZSK) (KOLKMAN; MEKKING; GIEBEN, 2012).
DS: publicado na zona pai, guarda o resumo criptográfico da KSK da zona filha, constituindo o elo da cadeia de confiança hierárquica.
Para provar a inexistência de nós sem permitir a enumeração de zonas, o NSEC3 utiliza resumos criptográficos dos nomes (LAURIE et al., 2008). Seus parâmetros operacionais residem no registro NSEC3PARAM. Conforme as diretrizes da RFC 9276, recomenda-se mitigar o impacto computacional nos resolvedores fixando a contagem de iterações em zero e dispensando o uso de salt, adotando a configuração expressa NSEC3PARAM 1 0 0 - (HARDAKER; DUKHOVNI, 2022). O provisionamento manual de cada uma dessas etapas adiciona múltiplos pontos de falha à operação.
2.1.6 O PowerDNS Authoritative Server
O PowerDNS Authoritative Server é um servidor autoritativo modular de código aberto que separa o motor de processamento DNS do armazenamento, usando backends como bancos de dados relacionais. Ele expõe uma API REST para gestão programática de zonas, registros e chaves criptográficas (POWERDNS, 2026), e é sobre essa API que o SanchezDNS foi construído.
Destacam-se dois mecanismos centrais do PowerDNS para esta pesquisa:
Metadado SOA-EDIT-API: quando definido com o parâmetro DEFAULT, delega ao servidor a atualização automática do serial SOA a cada alteração via API REST, gerando valores no padrão AAAAMMDD01 ou incrementando o serial existente em uma unidade, eliminando o erro de zone drift por esquecimento do operador.
Modos de operação de zona: no modo Native, a replicação dos dados é delegada diretamente à camada de armazenamento do backend (como replicação de banco de dados), dispensando transferências clássicas de zona; nos modos Primary e Secondary, o PowerDNS utiliza notificações NOTIFY e transferências AXFR/IXFR (POWERDNS, 2026).
2.2 ERROS DE CONFIGURAÇÃO EM ZONAS AUTORITATIVAS
2.2.1 O erro humano na operação de infraestrutura
Reason (1990) divide as falhas humanas em dois grupos: falhas de execução, nas quais a intenção está correta mas a ação diverge (deslizes e lapsos), e enganos (mistakes), quando a falha está no próprio planejamento ou na compreensão da situação. Norman (2013) leva esse modelo para a interação humano-computador e observa que deslizes e lapsos são frequentes mesmo entre operadores experientes em rotinas repetitivas. Para prevenir esses erros, o autor recomenda interfaces com restrições semânticas, funções de bloqueio (forcing functions) e automação dos passos que costumam ser esquecidos.
No DNS, esquecer o ponto final de um FQDN ou não atualizar o serial SOA são deslizes típicos de execução. Já delegar uma zona a um servidor que ainda não foi configurado é um engano de planejamento. Oppenheimer, Ganapathi e Patterson (2003) observaram em serviços corporativos de Internet que o erro de operador foi a principal causa de indisponibilidade em dois dos três serviços avaliados, e que os erros de configuração formavam a maior subcategoria.
A RFC 1912 (BARR, 1996) lista essas falhas operacionais no DNS: inconsistências A/PTR, omissão de incremento no serial SOA, ausência do ponto terminal em FQDNs, coexistência indevida com CNAME e delegações defeituosas (lame delegations).
2.2.2 Estudos empíricos de medição
Medições na Internet aberta mostram que essas falhas continuam prevalentes:
Pappas et al. (2004): identificaram delegações defeituosas em 15% das zonas avaliadas e dependências cíclicas em cerca de 2%. Mostraram que configurações incorretas degradam a disponibilidade e aumentam o tempo de resolução em até dez vezes, refutando a premissa do DNS clássico de que falhas de zonas operam de forma isolada e independente.
Bojović e Gajin (2017): analisaram 79.933 domínios (redes acadêmicas, ccTLD .rs e top 1.000 da Internet). Constataram que quase 5% dos domínios acadêmicos apontavam para servidores não autoritativos na zona pai e que entre 39,7% e mais de 45% dos domínios não possuíam registro PTR para seus servidores de nomes.
Akiwate et al. (2020): em estudo longitudinal de nove anos sobre 499 milhões de domínios, verificaram delegações defeituosas em aproximadamente 14% das amostras ativas. Esses problemas ficam mascarados por longos períodos graças à redundância parcial e expõem centenas de milhares de domínios a sequestro (hijacking). A prevalência se mantém em patamar semelhante ao observado por Pappas et al. (2004).
Sommese et al. (2020): compararam os conjuntos de registros NS declarados pela zona pai e pela zona filha na totalidade dos domínios de segundo nível de .com, .net e .org, identificando inconsistência em mais de treze milhões de domínios, dos quais ao menos cinquenta mil figuram entre o milhão mais acessado. Por meio de experimentos com sondas RIPE Atlas, demonstraram que o comportamento dos resolvedores diante da divergência não é uniforme, o que converte a inconsistência em indisponibilidade intermitente e de diagnóstico difícil.
Liu, Hao e Wang (2016): caracterizaram os registros pendentes (dangling records), isto é, apontamentos que sobrevivem ao desprovisionamento do recurso referenciado, e localizaram 467 ocorrências exploráveis em 277 domínios entre os dez mil mais acessados e em 52 zonas do domínio .edu. Descreveram três vetores pelos quais um adversário assume o controle do subdomínio afetado, incluindo a obtenção de certificados TLS legítimos em nome da organização titular. Deixar de remover um registro, portanto, cria um risco de segurança.
Chung et al. (2017): acompanharam, por vinte e um meses, instantâneos diários dos registros DNSSEC de todos os domínios assinados em .com, .net e .org. Constataram que mais de 30% das zonas assinadas não possuíam o registro DS correspondente na zona pai, situação em que a assinatura não produz efeito de validação, que aproximadamente 70% das zonas nunca efetuaram rotação de chaves e que 91,7% das chaves de assinatura de zona empregavam parâmetros criptográficos considerados fracos. Os autores atribuem o quadro principalmente a deficiências no fluxo operacional oferecido por registrários e ferramentas de gestão, mais do que a limitações do protocolo.
Korczyński et al. (2016): mediram a prevalência de servidores autoritativos que aceitam atualizações dinâmicas sem autenticação, identificando 1.877 domínios vulneráveis (0,065%) em amostra aleatória de 2,9 milhões e 587 domínios (0,062%) entre o milhão mais acessado. A falha, que os autores denominam envenenamento de zona (zone poisoning), permite a qualquer agente da Internet inserir, alterar ou remover registros na zona afetada, anulando integralmente o controle de acesso pretendido pelo operador.
Lidos em conjunto, esses estudos mostram, antes de tudo, que a prevalência é estável: as taxas de delegação defeituosa de Pappas et al. (2004) e de Akiwate et al. (2020) são equivalentes, apesar dos dezesseis anos que as separam e da diferença de uma ordem de grandeza na população estudada. Isso sugere que a evolução das ferramentas de administração pouco mudou a incidência do problema. As falhas também são silenciosas. Em quase todos os casos, a configuração incorreta continua produzindo respostas aparentemente válidas, e só uma verificação ativa revela o problema. Por fim, elas se concentram na manutenção (alteração, migração e remoção de registros) mais do que na criação inicial da zona, e por isso o projeto dá mais atenção à atualização consistente de dados interdependentes do que ao provisionamento.
2.2.3 Impactos operacionais e de segurança
Uma falha de configuração pode causar dano de maneiras diferentes. A mais direta passa pelo cache: como a resposta obtida por um resolvedor vale pelo intervalo declarado no TTL, um erro publicado continua sendo servido aos clientes depois de corrigido no servidor de origem. A exposição dura, portanto, o tempo de permanência do erro na zona somado ao TTL do registro afetado. No incidente do ccTLD sueco, a zona foi corrigida em cerca de uma hora, mas os efeitos duraram bem mais (PINGDOM, 2009).
Delegações defeituosas agem de outro modo e degradam a redundância. Pappas et al. (2004) mostraram que elas não causam indisponibilidade imediata, porque o domínio continua respondendo pelos servidores restantes, mas consomem sem aviso a margem de tolerância a falhas que a operação supunha ter. Também elevam o tempo de resolução em até uma ordem de grandeza quando o resolvedor consulta um servidor não autoritativo antes de obter uma resposta válida. O custo só aparece quando ocorre um segundo evento e a redundância já não existe.
Há ainda falhas operacionais que se tornam vulnerabilidades. Akiwate et al. (2020) demonstraram que delegações defeituosas expõem centenas de milhares de domínios ao sequestro por terceiros capazes de registrar o nome do servidor de nomes inexistente para o qual a delegação aponta; Liu, Hao e Wang (2016) documentaram trajetória análoga a partir de registros pendentes; e Korczyński et al. (2016) mostraram a inserção arbitrária de registros em zonas cujos servidores aceitam atualizações não autenticadas. Nesses casos o protocolo funciona como especificado, e o adversário explora um estado inconsistente deixado por uma operação incompleta.
Incidentes de grande escala mostram esses mecanismos na prática. Além do episódio do ccTLD sueco, houve a interrupção do serviço Edge DNS da Akamai em 22 de julho de 2021, causada por uma atualização de configuração depois revertida (AKAMAI, 2021), e a indisponibilidade dos serviços da Meta em 4 de outubro do mesmo ano, por mais de sete horas, quando um comando de manutenção provocou a retirada dos anúncios de rota dos prefixos que abrigavam seus servidores autoritativos (META, 2021). O relato oficial deste último caso interessa diretamente a esta pesquisa: a organização tinha um mecanismo automatizado de auditoria para impedir comandos desse tipo, e o incidente ocorreu porque uma falha nesse mecanismo deixou o comando passar (META, 2021). Ter um mecanismo de validação, portanto, não basta se a alteração puder seguir adiante quando ele falha.
2.2.4 Síntese: taxonomia de erros, evidências e mecanismos de mitigação
A revisão anterior permite organizar os achados em uma taxonomia operacional, apresentada no Quadro 1, que associa cada classe de erro à evidência empírica de sua prevalência, ao impacto documentado, ao mecanismo de mitigação previsto no artefato e à questão de pesquisa correspondente. O quadro responde à questão QP1 ao separar as classes elimináveis por construção daquelas que só admitem detecção posterior. Essa separação define o escopo do artefato, porque classes que dependem de ação junto à zona pai ou ao registrário não podem ser resolvidas por uma camada de operação sobre o servidor autoritativo. O quadro também é a origem dos requisitos do Capítulo 3, o que garante que cada mecanismo implementado corresponda a um problema documentado, e fornece a classificação de erros usada na avaliação da seção 2.6, para que a redução de erros possa ser medida por categoria e no total.
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

Fonte: elaborado pelo autor a partir das referências indicadas na coluna de evidências.
Das dez classes relacionadas, seis podem ser eliminadas por construção em uma camada de operação sobre o servidor autoritativo: divergência A/PTR, serial não incrementado, erro sintático no RDATA, registros reversos órfãos, DNSSEC incompleto no ato da criação e privilégio excessivo. Nelas, a interface consegue impedir que o estado incorreto seja produzido. Duas classes, delegação defeituosa e inconsistência entre pai e filho, dependem de ação em sistemas de terceiros e admitem no máximo detecção e sinalização. A atualização dinâmica não autenticada é eliminada pela própria arquitetura, já que não existe canal de escrita sem autenticação. A classe restante, ausência de rastreabilidade, não é propriamente um erro de configuração; ela entra na taxonomia porque impede o diagnóstico das demais. A seção 2.6 usa essa classificação para definir os critérios de êxito do artefato.
2.3 CONTROLE DE ACESSO E AUDITORIA
Em sistemas multiusuário, a integridade operacional depende do princípio do menor privilégio (Least Privilege), formalizado por Saltzer e Schroeder (1975): cada sujeito deve ter apenas as permissões necessárias para sua função, o que limita o alcance de acidentes e de condutas inadequadas. Em zonas DNS, isso significa restringir cada operador aos domínios pelos quais ele responde formalmente.
O Controle de Acesso Baseado em Papéis (Role-Based Access Control, RBAC), formalizado por Sandhu et al. (1996), reduz o trabalho de gestão ao associar as permissões a perfis funcionais (roles) em vez de identidades individuais. Quando é preciso controlar recursos específicos, usam-se abordagens híbridas, em que papéis gerais convivem com concessões ligadas a entidades específicas. O SanchezDNS usa essa estrutura para atribuir direitos separados de leitura e de escrita em cada zona.
A auditoria complementa o controle de acesso. Segundo o guia NIST SP 800-92 (KENT; SOUPPAYA, 2006), o registro contínuo de eventos de segurança permite a análise forense de incidentes, o diagnóstico operacional e a conformidade regulatória. Para isso, cada entrada de log deve registrar no mínimo a ação executada, a data e hora (timestamp), o identificador do agente e o recurso modificado.
2.4 DESIGN SCIENCE RESEARCH
2.4.1 Fundamentos metodológicos
A pesquisa se apoia nas ciências do artificial de Simon (1996), voltadas a prescrever soluções para metas definidas, e não a descrever fenômenos naturais. March e Smith (1995) definem quatro tipos de produto da pesquisa em tecnologia da informação (construtos, modelos, métodos e instanciações), obtidos em dois ciclos principais: construir e avaliar.
Hevner et al. (2004) estruturaram o paradigma em Sistemas de Informação sob sete diretrizes norteadoras: artefato como objeto de estudo, relevância ao negócio, avaliação rigorosa, contribuições claras, rigor nos métodos, design como processo de busca heurística e comunicação ampla dos resultados.
Quanto à natureza da contribuição, Gregor e Hevner (2013) propõem uma matriz que relaciona a maturidade do problema com a da solução (Invenção, Melhoria, Exaptação e Design Rotineiro) e três níveis de abstração dos artefatos: instanciações situadas (nível 1), teorias de design emergentes (nível 2) e teorias de design consolidadas (nível 3). Como os erros de operação em DNS autoritativo já são conhecidos na literatura e a plataforma propõe uma arquitetura integrada para mitigá-los, o projeto se enquadra no quadrante de Melhoria (Improvement), com contribuições nos níveis 1 (o software como instanciação funcional) e 2 (princípios de design preventivo).
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
Para medir a usabilidade percebida, a literatura emprega a System Usability Scale (SUS) (BROOKE, 1996). Com dez afirmações avaliadas em escala Likert de 5 pontos, a SUS produz uma pontuação de 0 a 100, interpretada segundo as faixas de Bangor, Kortum e Miller (2008). Neste trabalho a SUS não é aplicada, porque depende de um grupo de operadores (seção 2.6.1); ela fica indicada como instrumento da avaliação com usuários proposta como trabalho futuro.
2.5 TRABALHOS CORRELATOS
Esta seção compara as principais estratégias de gestão de zonas autoritativas, com base nas documentações oficiais e nos históricos de versão consultados em setembro de 2026. A ordem vai da edição direta de arquivos às ferramentas com mais automação, passa pela validação independente e termina em um quadro comparativo que mostra a lacuna tratada pelo artefato.
2.5.1 BIND 9 e manipulação direta de arquivos
O BIND 9, desenvolvido pelo Internet Systems Consortium (ISC), é a implementação de referência do DNS na Internet. Em configurações convencionais, seus registros ficam em arquivos de texto no formato padrão master file (MOCKAPETRIS, 1987b). A edição é manual, sem interface gráfica nativa, e o controle de acesso e o histórico de alterações dependem do sistema de arquivos hospedeiro ou de repositórios externos. Manter os registros A e PTR alinhados e incrementar o serial SOA a cada edição fica inteiramente por conta da disciplina do operador.
2.5.2 Interfaces web de software livre (PowerDNS-Admin)
O ecossistema do PowerDNS conta com ferramentas de código aberto consolidadas, como o PowerDNS-Admin, desenvolvido em Python/Flask, que oferece painel web com suporte a papéis de usuário, permissões por zona e autenticação centralizada (LDAP, OAuth). No entanto, a sincronização de registros PTR é uma opção secundária, desligada por padrão (auto_ptr), a geração de DNSSEC fica fora do fluxo básico de criação da zona e não foi encontrada na documentação nenhuma rotina de detecção de reversos órfãos (POWERDNS-ADMIN, 2026).
2.5.3 Technitium DNS Server
O Technitium DNS Server é um servidor de código aberto que acumula as funções autoritativa e recursiva e reúne, no mesmo produto, o motor de resolução e o painel web de administração. Seu histórico de versões mostra a incorporação gradual de recursos muito parecidos com os pretendidos nesta pesquisa: a versão 5.4, de 18 de outubro de 2020, introduziu a opção de criação da zona reversa no momento da inclusão de registros A ou AAAA; a versão 8.0, de 26 de março de 2022, acrescentou a assinatura DNSSEC de zonas primárias com algoritmos RSA e ECDSA diretamente pela interface; e a versão 9.0, de 24 de setembro de 2022, implementou controle de acesso baseado em papéis com permissões atribuíveis no nível da zona a usuários e a grupos, além de credenciais de API não expiráveis para automação (TECHNITIUM, 2026). É o trabalho correlato funcionalmente mais próximo do artefato proposto.
Mesmo assim, o SanchezDNS se justifica. O primeiro motivo é de arquitetura: o Technitium é um servidor completo, e sua camada administrativa não se separa do motor de resolução. Organizações que já operam sobre o PowerDNS Authoritative Server (escolha muitas vezes motivada pelo desacoplamento entre o motor e o backend de armazenamento, que permite replicação pela própria camada de banco de dados) teriam de substituir a infraestrutura autoritativa para adotá-lo, o que custa muito mais do que adotar uma camada de operação. A proposta desta pesquisa é, por isso, um painel para um servidor já implantado.
O segundo motivo reforça a relevância da questão QP3. A versão 15.5 do Technitium, de 19 de setembro de 2026, corrigiu uma falha de autorização, reportada por Tao Pan, que permitia usar a opção ptr das chamadas de inclusão, atualização e exclusão de registros para criar, sobrescrever ou remover registros PTR em zonas reversas nas quais o usuário autenticado não tinha permissão de escrita (TECHNITIUM, 2026). A falha apareceu quatro anos depois da introdução das permissões por zona, em um produto maduro e amplamente implantado, e resultou da interação entre os dois mecanismos que este trabalho pretende combinar. Isso acontece porque o registro derivado fica em zona diferente daquela em que o operador atua, e a automação acaba atravessando fronteiras de autorização. Combinar automação e controle de acesso granular exige, portanto, um tratamento de projeto próprio. O modelo adotado para isso é descrito no Capítulo 3 e verificado pelos testes negativos de autorização especificados na seção 2.6.2.
2.5.4 Microsoft DNS Server e o console DNS Manager
O servidor DNS do Windows Server, administrado pelo console DNS Manager, é a solução predominante em redes corporativas baseadas em Active Directory. Há mais de duas décadas ele oferece a opção de atualizar o registro de ponteiro associado (Update associated pointer (PTR) record) na criação e na edição de registros de host, o que automatiza a manutenção do reverso correspondente, e permite delegar permissões no nível da zona pelas listas de controle de acesso do diretório, na guia de segurança das propriedades da zona (MICROSOFT, 2026). Tem ainda o mecanismo de envelhecimento e remoção automática (aging and scavenging), que elimina registros criados dinamicamente cujo carimbo de tempo não é renovado. É uma solução parcial para os registros órfãos, porque vale só para registros de origem dinâmica e não alcança os criados manualmente pelo operador.
Algumas restrições limitam seu uso no problema tratado aqui. A delegação granular de permissões exige zonas integradas ao Active Directory Domain Services, o que prende a solução ao ecossistema Windows e ao modelo de diretório corporativo e a torna inaplicável a infraestruturas baseadas em servidores autoritativos de código aberto. A automação do registro PTR pela interface gráfica está documentada como inoperante quando a zona reversa é subdividida em sub-redes menores que um octeto (MICROSOFT, 2026), cenário comum em organizações que recebem blocos de endereçamento inferiores a /24 e em que o preenchimento manual é mais sujeito a erro. Além disso, a trilha de auditoria fica nos subsistemas de auditoria do sistema operacional e do diretório, fora do console de DNS, o que espalha a evidência entre componentes distintos e dificulta reconstituir o histórico de uma zona específica.
2.5.5 NetBox DNS e a abordagem de fonte externa da verdade
O plugin NetBox DNS estende o NetBox, sistema de gestão de infraestrutura e de endereçamento, para modelar servidores de nomes, zonas, visões e registros. Entre seus recursos figuram a criação e a atualização automáticas de registros PTR a partir de registros de endereço, a geração automática do serial da zona e dos registros SOA e NS, e a validação de nomes e valores conforme as RFCs pertinentes; o mecanismo IPAM DNSsync permite derivar registros diretos e reversos diretamente dos objetos de endereço IP associados a prefixos vinculados a visões de DNS (NETBOX DNS, 2026). O plugin herda ainda, do NetBox, um registro de alterações por objeto e um modelo de permissões baseado em objetos e em conjuntos de restrições.
A abordagem é conceitualmente distinta da adotada nesta pesquisa. No NetBox DNS, o inventário passa a ser a fonte da verdade e o servidor autoritativo vira um destino de sincronização, o que exige um componente adicional de publicação e admite, por construção, divergência entre o que está cadastrado e o que efetivamente responde às consultas. Essa divergência é uma classe de erro operacional que não existe quando o servidor é a autoridade. O SanchezDNS faz o contrário: mantém a API do PowerDNS como fonte única da verdade e opera sobre o estado real do servidor, sem precisar de reconciliação. As duas escolhas são legítimas e atendem a contextos organizacionais diferentes, de modo que as ferramentas não disputam o mesmo espaço de projeto.
2.5.6 Knot DNS e as interfaces de controle programáticas
O Knot DNS, desenvolvido pelo CZ.NIC e amplamente empregado na operação de domínios de topo, é um servidor autoritativo cuja administração se dá pelo utilitário knotc sobre um soquete de controle local, complementado por vinculações de biblioteca que permitem automação programática. O modelo de edição de zonas é transacional: as alterações são iniciadas com zone-begin, aplicadas com zone-set e zone-unset, inspecionadas antes da efetivação com zone-diff e confirmadas com zone-commit ou descartadas com zone-abort; enquanto a transação permanece aberta, todos os eventos que modificam a zona, inclusive a assinatura DNSSEC automática, ficam bloqueados (CZ.NIC, 2026). O servidor automatiza igualmente a política de serial e a gestão do ciclo de vida das chaves criptográficas, incluindo as trocas programadas.
Duas características afastam o Knot do problema tratado nesta pesquisa. A interface é de linha de comando e de biblioteca, sem painel web, e por isso é pouco acessível a operadores sem familiaridade com o terminal, o perfil para o qual a literatura de erro humano recomenda restrições semânticas na interface (NORMAN, 2013). Além disso, o controle de acesso se resume às permissões do soquete e do sistema operacional, sem noção de usuário de aplicação, de papel ou de permissão por zona; quando é preciso segregar responsabilidades, isso é feito fora do servidor. A semântica transacional, por outro lado, é uma referência de projeto direta para o artefato. O SanchezDNS precisa emular a atomicidade do Knot ao propagar uma alteração que afeta ao mesmo tempo a zona direta e a reversa, já que uma efetivação parcial produziria a divergência que o mecanismo pretende evitar. Essa exigência entra nos requisitos apresentados no Capítulo 3.
2.5.7 Zonemaster e a validação independente de delegações
O Zonemaster é uma ferramenta de validação da qualidade de delegações DNS desenvolvida conjuntamente pela Internetstiftelsen, operadora do ccTLD sueco, e pela AFNIC, operadora do ccTLD francês, em substituição às ferramentas DNSCheck e Zonecheck anteriormente mantidas por cada uma dessas organizações. Sua característica mais relevante para esta pesquisa é a existência de um plano de testes formalmente especificado segundo a estrutura da norma IEEE 829-2008, composto por setenta e quatro casos de teste distribuídos em nove categorias (Address, Basic, Connectivity, Consistency, Delegation, DNSSEC, Nameserver, Syntax e Zone), cada um com critérios de aprovação, mensagens e níveis de severidade definidos, e executável por interface web, por linha de comando ou por serviço de retaguarda com saída em JSON (ZONEMASTER, 2026).
O Zonemaster, contudo, não é uma ferramenta de operação: atua sobre a delegação já publicada, não previne o erro, não gerencia permissões e não registra autoria das alterações. Ele representa a estratégia de verificação independente e posterior, que complementa a prevenção no momento da edição sem substituí-la. Nesta pesquisa, o Zonemaster aparece como trabalho correlato e também como instrumento de medição na avaliação do artefato, conforme a seção 2.6.2. Serve a esse segundo papel porque seus casos de teste são especificados formalmente, o que torna a medição reprodutível e auditável por terceiros, porque sua saída estruturada permite coleta automatizada e, principalmente, porque é independente do objeto avaliado, o que evita aferir o artefato pelos mesmos critérios que ele próprio implementa.
2.5.8 DNS como código (IaC) e plataformas em nuvem
DNS como código: ferramentas como octoDNS e DNSControl empregam o paradigma de Infraestrutura como Código (IaC). As zonas são modeladas como arquivos de configuração versionados (Git), incorporando validações pré-publicação (dry-run/preview) e auditoria herdada do histórico de commits (OCTODNS, 2026; DNSCONTROL, 2026). Essa abordagem, contudo, carece de interface gráfica para operadores não técnicos e o controle de acesso é definido por repositório, sem granularidade por zona.
Provedores em nuvem: serviços corporativos como Cloudflare DNS, AWS Route 53 e Google Cloud DNS oferecem alta disponibilidade, suporte nativo a DNSSEC e interfaces amigáveis com controle de acesso integrado. No entanto, envolvem custos recorrentes, risco de dependência de fornecedor (vendor lock-in) e incompatibilidade com ambientes corporativos que demandam custódia estrita de dados em infraestruturas locais (on-premises).
2.5.9 Quadro comparativo e delimitação da contribuição
O Quadro 2 resume a análise em cinco dimensões, tiradas das classes de erro elimináveis por construção do Quadro 1. A classificação usa três valores: "Sim", quando o recurso faz parte do fluxo principal da ferramenta e funciona sem configuração adicional; "Parcial", quando existe mas depende de ativação explícita, cobre apenas parte dos casos ou é oferecido por componente externo à ferramenta; e "Não", quando não existe.
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
Nenhum dos recursos do SanchezDNS é inédito isoladamente. A criação automática de registros reversos existe no Technitium desde 2020 e no produto da Microsoft há mais de duas décadas; o controle de acesso por zona está presente no Technitium, no Microsoft DNS e, de forma parcial, no PowerDNS-Admin; a trilha de auditoria é oferecida pelo Poweradmin com registro do estado anterior e posterior dos registros (POWERADMIN, 2026); a validação prévia à publicação é central nas ferramentas de infraestrutura como código; e a verificação formal da qualidade da delegação é o objeto do Zonemaster. Não há, portanto, ineditismo funcional.
A contribuição do trabalho está em outro lugar. Primeiro, na combinação, em uma única camada desacoplada e aplicável a infraestruturas locais já padronizadas sobre o PowerDNS, de mecanismos hoje espalhados por ferramentas que não podem ser adotadas ao mesmo tempo, porque competem pelo papel de servidor autoritativo ou de fonte da verdade. Depois, no tratamento do ciclo de vida completo do registro reverso sobre o próprio servidor autoritativo. Entre as ferramentas examinadas, só o NetBox DNS e as de DNS como código removem o reverso junto com o registro de origem, e as duas o fazem a partir de uma fonte da verdade externa ao servidor, o que admite a divergência entre inventário e servidor discutida na seção 2.5.5. Entre as que operam diretamente sobre o servidor autoritativo, nenhuma condiciona a remoção do reverso à verificação de que o endereço deixou de ser referenciado por outro nome, nem oferece triagem dos reversos órfãos já existentes. Essa lacuna aparece nas colunas correspondentes do Quadro 2, e os achados de Liu, Hao e Wang (2016) mostram por que ela importa. Por fim, há a avaliação desses mecanismos, com taxonomia de erros definida antes dos testes, oráculo de medição independente e implantação em ambiente real. A documentação das ferramentas comparadas descreve funcionalidades, mas raramente verifica se elas impedem os erros que se propõem a evitar. Nos termos de Gregor e Hevner (2013), o trabalho se situa no quadrante de melhoria (improvement), com contribuições nos níveis 1 e 2 de abstração.
2.6 ESTRATÉGIA DE AVALIAÇÃO DO ARTEFATO
A seção 2.4.3 apresentou o FEDS. Esta seção define como ele é aplicado para sustentar ou refutar a afirmação de que o SanchezDNS mitiga os problemas caracterizados na seção 2.2. Os episódios, as métricas e os critérios de êxito foram fixados antes da execução de cada episódio, o que impede que os critérios sejam ajustados aos resultados obtidos.
2.6.1 Escolha da estratégia de avaliação
Venable, Pries-Heje e Baskerville (2016) propõem quatro estratégias de avaliação, cuja seleção decorre da natureza do risco dominante no projeto e do custo relativo dos episódios: Quick & Simple, apropriada quando riscos técnicos e sociais são baixos; Technical Risk and Efficacy, quando o risco predominante é a viabilidade técnica; Human Risk and Effectiveness, quando o risco predominante é o de que o artefato, ainda que tecnicamente correto, não produza o efeito pretendido em uso real; e Purely Technical Artifact, quando não há usuários humanos envolvidos.
O risco que mais pesa neste trabalho é o de os mecanismos não se comportarem como especificado. A automação do reverso opera sobre duas zonas, a normalização decide pelo operador e a escrita derivada atravessa a fronteira de permissão por zona; um defeito em qualquer um desses pontos produz justamente a classe de erro que o mecanismo deveria eliminar. Adota-se, por isso, a estratégia Technical Risk and Efficacy, com episódios artificiais e um episódio final em ambiente real.
A estratégia Human Risk and Effectiveness exigiria um grupo de operadores e uma organização disposta a submeter sua infraestrutura de produção a um período de observação, recursos que não estão ao alcance de um trabalho de conclusão de curso. A consequência é que a avaliação mostra se os mecanismos impedem os erros quando exercitados, mas não mede se operadores reais passam a cometer menos erros ao usá-los. Essa distinção delimita o alcance das conclusões sobre a QP4 e é retomada na seção 2.6.4.
2.6.2 Episódios de avaliação
O desenho compreende cinco episódios. A numeração é empregada de forma consistente no Quadro 3 e nos capítulos subsequentes.
E0. Caracterização documental (formativo, artificial). Consiste na construção da taxonomia apresentada no Quadro 1 e na matriz de rastreabilidade do Quadro 5, que associa cada requisito do artefato a uma classe de erro documentada na literatura. O episódio responde à QP1 e estabelece o critério de suficiência do escopo: um requisito sem classe de erro associada é candidato à remoção, e uma classe eliminável sem requisito associado indica lacuna de projeto.
E1. Avaliação durante a construção (formativo, artificial). Corresponde aos ciclos de construir e avaliar percorridos ao longo do desenvolvimento, de outubro de 2025 a setembro de 2026, registrados no histórico de versões do repositório. A cada ciclo, as funcionalidades novas foram exercitadas sobre uma instância de laboratório do PowerDNS, com conferência do estado resultante das zonas. O episódio não produz medida quantitativa; seu produto são os retornos da avaliação ao projeto, que alteraram requisitos e estão documentados no Capítulo 4.
E2. Conformidade técnica (somativo, artificial). Executado sobre instância de laboratório composta pelo artefato, por uma instância do PowerDNS e por zonas de teste. Para cada classe de erro eliminável do Quadro 1, constrói-se ao menos um caso de teste negativo que tenta reproduzir o erro pela interface e pela API do artefato; o caso é aprovado quando o erro é impedido ou corrigido automaticamente, e reprovado quando é aceito e publicado. Incluem-se testes de autorização negativos no caminho de automação do registro reverso, que endereçam a QP3 e a classe de falha do incidente descrito na seção 2.5.3. O estado resultante de cada zona é verificado por dois oráculos independentes do artefato: a consulta direta ao servidor autoritativo, comparada ao estado esperado, e a execução do Zonemaster, com registro das reprovações por categoria de caso de teste. Como referência, as mesmas classes de erro são introduzidas de propósito em zonas editadas diretamente no PowerDNS, sem a plataforma, e submetidas ao mesmo oráculo. O uso de oráculo externo é intencional: avaliar o artefato apenas pelas validações que ele mesmo implementa seria um raciocínio circular, incapaz de detectar justamente os erros que o projetista não previu.
E3. Comparação com a operação manual (somativo, artificial). Seis tarefas típicas de administração de zonas são executadas de duas formas, pela linha de comando do PowerDNS e pela plataforma: criação de zona com ativação de DNSSEC; inclusão de registro de host com o reverso correspondente; alteração do endereço de um host já publicado; remoção de um host e de seu reverso; inclusão de registros MX e TXT com sintaxe específica; e correção de um conjunto de registros NS. Para cada tarefa, registram-se o número de passos que o operador precisa executar, o número de valores que precisa digitar em formato específico e as classes de erro do Quadro 1 que continuam possíveis. A medida é analítica: indica onde a plataforma retira do operador oportunidades de erro, mas não a frequência com que operadores reais errariam em cada caso.
E4. Implantação em ambiente real (somativo, naturalístico). Detalhado na seção 2.6.3.
2.6.3 Implantação em ambiente real
O episódio E4 verifica se o artefato funciona fora do laboratório, sobre um servidor autoritativo em operação. A plataforma é implantada a partir da imagem de contêiner publicada e da composição de serviços descrita na seção 4.7, seguindo a documentação do projeto, escrita com VitePress no diretório docs do repositório, e conectada a uma instância do PowerDNS Authoritative Server em uso.
O episódio tem dois objetivos. O primeiro é mostrar que o artefato pode ser posto em operação apenas com a documentação e a imagem publicada, sem compilação a partir do código-fonte, condição para que a contribuição seja utilizável além deste trabalho. O segundo é observar os mecanismos sobre zonas que não foram preparadas para o teste, nas quais aparecem situações que o laboratório não reproduz, como registros criados fora da plataforma e reversos que já estavam órfãos antes da implantação.
São coletados:
os passos necessários para implantar a plataforma a partir da documentação;
o resultado do Zonemaster sobre as zonas sob gestão, no início e ao fim do período de operação;
o número de registros A e AAAA sem PTR correspondente e de reversos órfãos apontados pelas rotinas de triagem, nos mesmos dois instantes;
a proporção das alterações feitas pela plataforma que constam da trilha de auditoria com autor, instante e recurso.
O episódio não envolve uma equipe de operadores. Ele mostra que o artefato funciona em ambiente real e que os mecanismos se comportam como especificado sobre dados que não foram preparados para o teste. Não constitui, porém, um estudo de caso no sentido de Yin (2018), que exigiria uma organização, uma equipe de operadores e um período de observação comparado a uma linha de base.
2.6.4 Matriz de avaliação, ameaças à validade e limitações
O Quadro 3 relaciona questões de pesquisa, construtos, métricas, instrumentos, episódios e critérios de êxito, de modo que cada questão formulada na seção 1.3 tenha um critério verificável e uma fonte de dados definida previamente.
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
Casos de teste negativos; consulta direta ao servidor; Zonemaster
E2
100% das tentativas nas classes elimináveis impedidas; nenhuma reprovação nova no Zonemaster após operações pelo artefato
QP2
Riscos introduzidos pela automação
Número de efeitos colaterais não intencionais (registros criados ou removidos fora da intenção declarada)
Diferença do estado da zona antes e depois de cada operação
E1, E2
Zero efeitos colaterais não declarados; cada risco identificado registrado com sua mitigação
QP3
Contenção da escrita derivada
Escritas em zona reversa sem permissão do operador fora das condições da seção 3.4.5
Testes de autorização negativos no caminho de automação do PTR
E2
Zero sobrescritas de PTR existente; zero escritas em faixa sem zona reversa provisionada; operador sem permissão não lista nem edita a zona reversa
QP4
Redução das oportunidades de erro
Passos, valores em formato específico e classes do Quadro 1 ainda possíveis, por tarefa
Execução das seis tarefas pela linha de comando e pela plataforma
E3
Nenhuma tarefa com mais passos ou valores formatados na plataforma; classes elimináveis impossíveis na plataforma
QP4
Funcionamento em ambiente real
Implantação a partir da documentação; reprovações do Zonemaster e reversos ausentes ou órfãos no início e no fim da operação
Documentação do projeto; Zonemaster; rotinas de triagem
E4
Implantação concluída só com a documentação; nenhuma reprovação nova atribuível ao artefato; reversos ausentes e órfãos não aumentam
QP4
Rastreabilidade
Proporção das alterações feitas pela plataforma com autor, instante e recurso registrados
Trilha de auditoria confrontada com o estado do servidor
E2, E4
100% das alterações efetuadas pela plataforma

Fonte: elaborado pelo autor (2026).
As ameaças à validade do desenho são classificadas a seguir segundo as categorias de Wohlin et al. (2012).
Validade interna. A avaliação do artefato por quem o desenvolveu configura risco de viés do experimentador: os casos de teste podem refletir, sem intenção, apenas os erros que o projetista já previu. As contramedidas são a derivação dos casos de teste a partir da taxonomia do Quadro 1, fixada antes da execução, e a verificação do estado das zonas por oráculo externo ao artefato.
Validade de construto. O construto "erro de configuração" é operacionalizado pela taxonomia do Quadro 1 e pelos casos de teste do Zonemaster, recorte que não captura erros semânticos: para esses instrumentos, um registro sintaticamente válido que aponta para o endereço errado é indistinguível de um registro correto. Da mesma forma, o E3 mede oportunidades de erro, e não erros cometidos: uma tarefa com menos passos oferece menos ocasiões de deslize, mas a avaliação não observa quantos deslizes operadores reais cometeriam.
Validade externa. Todos os episódios usam um único servidor autoritativo, o PowerDNS, e zonas de porte moderado, e nenhum deles envolve uma equipe de operadores. Os resultados não se estendem a operadores de domínios de topo, a infraestruturas com dezenas de milhares de zonas, a servidores autoritativos diversos do PowerDNS nem ao comportamento de equipes de operadores.
Validade de conclusão. Os episódios não envolvem amostras nem testes estatísticos. As conclusões são técnicas e analíticas: o artefato impede, ou não, cada classe de erro exercitada, e retira, ou não, passos e valores que a operação manual exige.
A principal limitação do desenho é a ausência de avaliação com operadores. Um experimento controlado com participantes e um estudo de caso em organização real, que mediriam a taxa de erro e a usabilidade percebida pela System Usability Scale, ficam como trabalho futuro e são retomados no Capítulo 6. Como nenhum episódio envolve participantes, não há coleta de dados pessoais; nomes de domínio e endereços das zonas do E4 são pseudonimizados no relato.
2.7 CONSIDERAÇÕES DO CAPÍTULO
Este capítulo apresentou a fundamentação conceitual do SanchezDNS. Revisou o modelo operacional do DNS e os componentes manipulados pelo artefato; mostrou, com estudos empíricos que cobrem duas décadas de medições, que os erros de configuração em zonas autoritativas continuam frequentes, e os organizou em uma taxonomia de dez classes associadas a evidências, impactos e mecanismos de mitigação; apresentou os fundamentos de controle de acesso e auditoria e o método Design Science Research; e examinou os trabalhos correlatos, concluindo, a partir do quadro comparativo, que a contribuição do artefato está na combinação de mecanismos hoje dispersos, no tratamento do ciclo de vida completo dos registros reversos e na avaliação empírica de seus efeitos. Também definiu a estratégia de avaliação, com episódios, métricas, critérios de êxito fixados antes da execução e as ameaças à validade. O Capítulo 3 detalha a metodologia da pesquisa e o modelo de engenharia do artefato, derivando os requisitos das classes de erro do Quadro 1.
3 METODOLOGIA E MODELAGEM DO ARTEFATO
Este capítulo apresenta a metodologia aplicada ao desenvolvimento do SanchezDNS e a modelagem técnica de seus componentes. A condução do projeto adota o método Design Science Research (DSRM) de Peffers et al. (2007), derivando os requisitos das falhas operacionais levantadas no Capítulo 2. Em seguida, especificam-se os cinco modelos que orientam a implementação: autenticação e sessão, controle de acesso, normalização sintática de registros, sincronização de zonas reversas e auditoria de eventos.
3.1 CLASSIFICAÇÃO DA PESQUISA
A pesquisa é aplicada quanto à natureza, descritiva e prescritiva quanto aos objetivos, e combina uma abordagem qualitativa (revisão da literatura e análise de soluções existentes) com uma abordagem quantitativa na contagem dos casos de teste aprovados, das reprovações do Zonemaster e dos passos exigidos por tarefa.
O delineamento é o da Design Science Research. March e Smith (1995) distinguem as ciências naturais, que explicam fenômenos existentes, das ciências do artificial, que constroem e avaliam artefatos. Hevner et al. (2004) sistematizam essa distinção em sete diretrizes, das quais três condicionam o desenho deste trabalho: produzir um artefato viável, responder a um problema documentado e demonstrar utilidade por método de avaliação rigoroso.
As falhas de configuração em DNS são documentadas pela literatura de medição há mais de duas décadas; o que falta é uma ferramenta que as impeça no momento da edição. Trata-se, portanto, de um problema de projeto. Na classificação de Gregor e Hevner (2013), a contribuição fica no quadrante de melhoria: solução nova para problema conhecido. Os mecanismos isolados já existem em ferramentas maduras, e a contribuição está em integrá-los em um único fluxo e medir o efeito dessa integração sobre a taxa de erro.
3.2 APLICAÇÃO DO MODELO DSRM
O modelo de Peffers et al. (2007) prevê seis atividades e quatro pontos possíveis de entrada. Adotou-se a entrada centrada no problema, porque a motivação partiu de falhas observadas na operação de zonas. O Quadro 4 mapeia cada atividade à seção correspondente.
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
   Modelagem dos cinco componentes técnicos; implementação dos serviços com Go, Nuxt e banco de dados
   3.4 a 3.7; Capítulo 4
4. Demonstração
   Execução de tarefas de administração em bancada de laboratório, com conferência por oráculo externo independente
   5.1
5. Avaliação
   Episódios formativos e somativos sob o framework FEDS; casos de teste negativos com oráculo externo, comparação com a operação manual e implantação em ambiente real
   2.6; Capítulo 5
6. Comunicação
   Formalização da monografia; disponibilização do código-fonte e documentação técnica sob licença livre
   Capítulo 6

Fonte: elaborado com base em Peffers et al. (2007).
O ciclo foi percorrido mais de uma vez. Dois retornos da avaliação ao projeto alteraram requisitos e estão documentados no Capítulo 4: a introdução da guarda de referência na remoção de registros PTR (seção 3.6.3) e a exibição do valor normalizado no lugar do valor digitado (seção 3.5.3). Os dois tiveram a mesma origem, já antecipada na questão QP2: ao automatizar uma decisão que antes era do operador, a plataforma pode apenas mudar o erro de lugar.
3.3 REQUISITOS
Os requisitos foram derivados diretamente da taxonomia do Quadro 1, sem entrevistas ou questionários, de modo que cada um responda a uma classe de erro com evidência empírica na literatura. O procedimento teve três passos:
Selecionaram-se as seis classes classificadas como elimináveis por construção, mais a classe de ausência de rastreabilidade;
Para cada classe, especificou-se o mecanismo que torna o estado incorreto inalcançável pela interface;
As classes dependentes de terceiros (delegação inconsistente e divergência de NS entre pai e filho) foram convertidas em requisitos de detecção, sem correção automática, conforme a delimitação da seção 1.5.
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
RNF01: o PowerDNS é a fonte única da verdade dos dados de zona. A aplicação não mantém cópia local dos registros; toda leitura é uma chamada à API do servidor autoritativo.
RNF02: toda autorização é decidida no backend, por requisição. O frontend apenas reflete o resultado.
RNF03: a falha de uma escrita derivada não aborta a operação principal.
RNF04: não existe caminho de escrita não autenticado. Todas as rotas de alteração exigem sessão válida, o que elimina por arquitetura a classe de atualização dinâmica não autenticada medida por Korczyński et al. (2016).
O RNF03 é uma escolha entre dois modos de falha. A seção 3.6.4 examina suas consequências para a integridade da solução.
3.4 AUTENTICAÇÃO E CONTROLE DE ACESSO
3.4.1 Autenticação e ciclo de vida da sessão
A autenticação usa uma sessão opaca mantida no servidor, em vez de um token autocontido. A escolha decorre do RNF02: como toda autorização é decidida no backend a cada requisição, o identificador entregue ao cliente não precisa carregar nenhuma informação; basta que seja impossível de adivinhar e possível de revogar. Um token autocontido do tipo JWT inverteria essas duas propriedades, pois embutiria o nível de acesso no próprio artefato entregue ao cliente e continuaria válido até expirar, mesmo depois da revogação do privilégio que declara.
Na autenticação bem-sucedida, o backend gera um identificador de sessão (SID) de 32 bytes obtidos do gerador criptográfico do sistema operacional, codificados em base64 sem preenchimento. Isso corresponde a 256 bits de entropia, o que torna inviável adivinhar o identificador por força bruta. O SID é persistido na coleção sessions do MongoDB, que mantém índice único sobre o campo sid, e devolvido ao cliente em cookie.
O documento de sessão registra, além do identificador e do e-mail do operador, os metadados de origem da conexão: endereço IP, sistema operacional e navegador derivados do cabeçalho User-Agent, e localização geográfica aproximada resolvida a partir do IP. Esses campos não participam da decisão de autorização; existem para dar contexto forense à trilha de auditoria da seção 3.7, que registra o autor da operação mas não as condições em que ela foi realizada.
O cookie é declarado com o atributo HttpOnly, o que impede sua leitura por código JavaScript na página e neutraliza o roubo de sessão por injeção de script. O atributo Secure, que restringe o envio ao canal TLS, é condicionado ao ambiente: ativo em produção e inativo em desenvolvimento, onde a aplicação roda sobre HTTP simples.
A validação ocorre no middleware de sessão, em duas camadas. Consulta-se primeiro o cache em Redis; havendo acerto, verifica-se apenas o campo de atividade. Na ausência do registro em cache, recorre-se ao MongoDB com filtro composto por SID e sessão ativa. Resolvido o e-mail, o nível global é obtido e injetado no contexto da requisição. O cache existe para evitar uma consulta ao banco por requisição; a fonte da verdade continua sendo o MongoDB, e a ausência no Redis nunca é interpretada como sessão inválida.
A sessão expira por inatividade, e não após um tempo absoluto. Cada requisição autenticada atualiza o campo de último acesso da sessão. Uma rotina agendada executa a cada dez minutos, seleciona as sessões ativas cujo último acesso seja anterior a sete dias, marca-as como inativas com o motivo registrado e remove as entradas correspondentes do cache. Com isso, a revogação tem efeito imediato: basta marcar a sessão como inativa para que a próxima requisição seja recusada, sem esperar a expiração de um token.
Há uma assimetria deliberada entre a validade do cookie e a da sessão. O cookie é emitido com prazo longo, para que o operador não precise reautenticar-se a cada visita; a sessão no servidor expira em sete dias de inatividade. A decisão cabe sempre ao servidor: se a sessão tiver sido desativada, a requisição é recusada mesmo que o cookie ainda esteja dentro do prazo, e o cookie é descartado na própria resposta. O cookie apenas transporta a referência à sessão; a credencial propriamente dita fica no servidor.
As senhas nunca são armazenadas em texto claro. O cadastro grava o resultado de bcrypt com o fator de custo padrão da biblioteca, prefixado pelo identificador do esquema, o que permite reconhecer o algoritmo empregado e migrar para outro sem ambiguidade. A verificação usa comparação em tempo constante da própria biblioteca, e hashes marcados com o esquema legado do crypt tradicional são recusados de forma incondicional, sem tentativa de verificação.
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
O primeiro cadastro realizado em uma instância vazia é gravado diretamente na coleção de cadastros com nível administrativo. Essa regra resolve a inicialização do sistema sem depender de uma credencial padrão embutida no código ou em variável de ambiente (o que seria, por si só, uma vulnerabilidade) e sem exigir intervenção manual no banco.
Os demais cadastros são gravados na coleção de solicitações, com estado pendente. A tentativa de autenticação de um solicitante ainda não apreciado é recusada com mensagem específica informando que o cadastro está em análise. Essa mensagem é propositalmente diferente da de credencial inválida, já que não revela nada que o próprio solicitante não saiba.
A aprovação por um administrador copia os dados da solicitação para a coleção de cadastros com nível de membro e registra o desfecho na solicitação, preservando o histórico da decisão. Um membro recém-aprovado não tem acesso a zona alguma: o vínculo por zona é concedido separadamente, conforme a seção 3.4.3.
A validação de entrada no cadastro exige endereço de correio eletrônico sintaticamente válido e senha de no mínimo oito caracteres, e normaliza o endereço em caixa baixa antes da gravação, de modo que a comparação posterior nos vetores de permissão seja determinística.
3.4.3 Estrutura das permissões
O SanchezDNS combina papéis globais com permissões por zona, em um modelo RBAC híbrido derivado de Sandhu et al. (1996). O objetivo é isolar o raio de impacto operacional: um operador responsável por um domínio não deve enxergar nem alterar zonas sob custódia de outra equipe, conforme o princípio de menor privilégio de Saltzer e Schroeder (1975).
O RBAC puro não atende bem a esse caso. Um papel global de "operador de zona" concederia acesso a todas as zonas; um papel por zona multiplicaria papéis sem ganho de abstração. A autorização opera, por isso, em dois níveis:
Nível global: Definido pelo campo level do documento de cadastro do operador, com dois valores possíveis. O perfil admin detém privilégio irrestrito sobre criação e remoção de zonas, gestão de usuários e consulta à auditoria. O perfil member não concede, por si, acesso a zona alguma. O primeiro cadastro criado na instância recebe admin automaticamente, o que dispensa credencial embutida no código ou em variável de ambiente.
Nível de zona: A coleção users do MongoDB mantém um documento por zona, contendo dois vetores de endereços de e-mail:
JSON
{
"zona": "exemplo.com.br.",
"leitura": ["operador1@org.br", "operador2@org.br"],
"escrita": ["admin-zona@org.br"]
}

A resolução da permissão segue precedência fixa: verifica-se se level == "admin"; na ausência de privilégio global, verifica-se a pertinência ao vetor escrita; na ausência desta, a pertinência ao vetor leitura; esgotadas as verificações, o acesso é negado com status HTTP 403 (Forbidden). A comparação de endereços de e-mail é normalizada em caixa baixa (case-insensitive). Como a escrita pressupõe e engloba a visualização, a concessão de escrita herda implicitamente o privilégio de leitura.
3.4.4 Aplicação nas rotas
A autorização é implementada em dois middlewares e uma função de resolução por zona. O middleware de sessão, descrito na seção 3.4.1, injeta email e level no contexto da requisição e recusa com HTTP 400 a requisição sem sessão válida, descartando o cookie na resposta.
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
O operador não recebe acesso à zona reversa: não pode listá-la, editá-la nem remover dela registros arbitrários. O que ele obtém é apenas o efeito determinístico de uma escrita que já estava autorizado a fazer na zona direta: o nome do PTR é função do endereço, e seu conteúdo, do nome do registro de origem.
É a primeira condição que delimita essa concessão. Como a plataforma não cria zonas reversas por iniciativa própria, a equipe responsável pelo espaço reverso controla quais faixas ficam sujeitas à derivação: basta não provisionar a zona. Onde ela não existe, nenhuma escrita ocorre e a operação direta se completa normalmente.
A solução é um afrouxamento deliberado do menor privilégio, cujos desdobramentos e riscos operacionais associados são discutidos criticamente entre as limitações do artefato no Capítulo 6.
3.5 NORMALIZAÇÃO SINTÁTICA
3.5.1 Estratégia
Para impedir que RDATA malformado chegue à API do PowerDNS, o backend aplica rotinas de normalização antes da submissão. Sempre que possível, a plataforma converte o valor para a forma canônica do tipo, em vez de rejeitar a entrada, e o operador não precisa conhecer a sintaxe exigida. É uma aplicação do conceito de forcing function descrito por Norman (2013), em que o estado incorreto nem chega a ser representável.
A validação é reservada aos casos sem conversão determinística. O esquema de validação do formulário, avaliado no frontend com a biblioteca Valibot e reavaliado no backend, exige os campos obrigatórios de cada tipo (prioridade, peso, porta e alvo para SRV; prioridade e nome de destino para HTTPS) e recusa TTL inferior a 60 segundos. A ausência de um campo obrigatório não admite suprimento automático e produz HTTP 400.
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
Esta seção especifica a automação que atende aos requisitos RF06 a RF10. É o componente mais complexo do artefato: opera sobre duas zonas, depende de uma correspondência que o protocolo não impõe mas que serviços de rede exigem (BARR, 1996), e é o único em que a remoção pode criar um erro novo.
3.6.1 Derivação do nome e da zona
A gravação do PTR exige duas informações derivadas do endereço.
Nome reverso. Para IPv4, os quatro octetos são invertidos sob in-addr.arpa: 192.0.2.10 produz 10.2.0.192.in-addr.arpa.. Para IPv6, o endereço é primeiro expandido para a forma plena de 32 dígitos hexadecimais, e cada nibble vira um rótulo isolado em ordem inversa sob ip6.arpa.. A expansão prévia é indispensável: a notação abreviada, que suprime sequências de zeros, não admite conversão direta em rótulos.
Zona de destino. O backend lista as zonas do servidor e seleciona, entre as que são sufixo do nome reverso, a de maior comprimento. A regra do sufixo mais longo acomoda a delegação de faixas menores que o octeto (EIDNES; DE GROOT; VIXIE, 1998): coexistindo 2.0.192.in-addr.arpa. e 128/25.2.0.192.in-addr.arpa., o registro vai para a mais específica, que é a que responde pelo endereço. Escolher a mais ampla gravaria o registro em zona sem autoridade sobre o nome: o registro existiria no banco, mas não responderia às consultas.
Se nenhuma zona reversa casar com o endereço, a função retorna sem erro e nenhuma escrita é feita.
3.6.2 Criação e alteração
Na criação de um registro A ou AAAA, o backend deriva o nome reverso, localiza a zona e verifica a existência de RRset PTR anterior. Havendo registro anterior, a rotina não o substitui: a escrita derivada não sobrepõe uma decisão tomada pelo responsável da zona reversa. O estado resultante, um registro direto sem PTR, pode ser detectado pelo endpoint GET /reverses, que implementa o RF10.
A alteração de endereço é executada em duas operações, em vez de uma substituição única:
Grava o PTR do novo endereço;
Se o endereço anterior for diferente do novo, avalia a remoção do PTR anterior.
A ordem é importante. Como a criação vem antes da remoção, uma falha na segunda etapa deixa o sistema com dois reversos, um deles obsoleto, o que é redundante, mas resolvível. A ordem inversa deixaria o endereço sem reverso, que é o estado que faz servidores SMTP rejeitarem mensagens.
3.6.3 Guarda de referência
A remoção é a operação que exige mais cautela, e foi a que motivou o primeiro retorno de iteração da seção 3.2. A formulação direta do RF08 (apagar o PTR ao apagar o registro A) falha sempre que mais de um nome aponta para o mesmo IP, caso corriqueiro em servidores com vários serviços. A remoção de um dos nomes apagaria o reverso ainda referenciado pelos demais, transformando uma limpeza na própria divergência A/PTR que o requisito combate.
O backend aplica, por isso, uma guarda de referência antes de qualquer remoção de PTR: percorre todas as zonas diretas do servidor, ignorando as terminadas em in-addr.arpa e ip6.arpa, e verifica se algum RRset A ou AAAA ainda contém o endereço. Havendo referência, a remoção é suprimida. A mesma guarda roda na etapa de remoção da alteração descrita acima.
Daí decorre uma assimetria entre criação e remoção, adotada como regra de projeto. A criação é otimista e ocorre a menos que exista uma decisão anterior em contrário; a remoção é conservadora e só ocorre se nenhuma referência restar. A razão é a diferença de custo entre os dois erros. Um reverso a mais é detectado pela triagem de órfãos, não quebra a resolução e pode ser corrigido a qualquer momento. Um reverso removido indevidamente quebra serviços na hora e só é detectado com verificação ativa. A automação foi, por isso, enviesada na direção do erro mais barato. A Figura 2 resume os dois fluxos.
Figura 2 – Sincronização de registros PTR: criação otimista e remoção conservadora

Fonte: elaborado pelo autor (2026).
3.6.4 Semântica de melhor esforço
As rotinas de reverso operam sob o RNF03: a falha da escrita derivada não aborta a operação principal. Cada uma roda com contexto de tempo limite próprio de seis segundos, para que a lentidão do servidor autoritativo não trave a requisição do operador.
Essa decisão introduz um risco que faz parte da resposta à segunda parte da QP2. O operador recebe HTTP 200 e a mensagem de sucesso, mas a escrita derivada pode não ter ocorrido. O resultado é a divergência A/PTR que o RF06 pretende eliminar, agravada pelo fato de o operador ter motivos para acreditar que ela não existe. Nesse caso, a automação apenas transfere a classe de erro do esquecimento do operador para uma falha silenciosa do sistema, que é ainda mais difícil de detectar.
O projeto inclui três mitigações, embora nenhuma resolva o problema por completo:
O endpoint GET /reverses (RF10) lista os registros A/AAAA sem PTR e permite criação em lote por PUT /reverses, o que converte a falha silenciosa em achado verificável, desde que o operador execute a verificação;
A falha da rotina é gravada no log da aplicação, preservando rastreabilidade do evento ainda que ele não seja comunicado na hora;
A ordem de operações da seção 3.6.2 faz com que uma falha parcial resulte em redundância, e não em ausência de reverso.
Com essa escolha, preserva-se a autonomia sobre o recurso principal ao custo de uma garantia mais fraca sobre o derivado. Essa limitação arquitetural é retomada no Capítulo 6.
A criação em lote do RF10 é sequencial e não transacional: cada PTR selecionado é submetido em requisição própria à API do PowerDNS e a primeira falha interrompe o processamento, preservando os registros já gravados e deixando os restantes por criar. O estado parcial é convergente, já que basta reexecutar a rotina, que recalcula os ausentes a partir do estado atual das zonas. Ainda assim, a operação não oferece garantia de atomicidade sobre o conjunto selecionado.
3.7 AUDITORIA
A auditoria atende aos requisitos RF13 e RF14. Kent e Souppaya (2006) identificam três finalidades para o registro de eventos (reconstrução de incidentes, atribuição de responsabilidade e detecção de uso indevido), cada uma com exigências próprias de granularidade e retenção.
Cada operação de escrita grava um documento na coleção logs do MongoDB:
JSON
{
"zone": "exemplo.com.br.",
"username": "operador@org.br",
"action": "edit_record",
"details": "Alterado registro A de www.exemplo.com.br.",
"createdAt": "2026-03-14T10:22:47-03:00"
}

O campo action usa códigos fixos em vez de texto livre, o que permite agregar os registros. O endpoint GET /logs, restrito ao administrador, oferece paginação e filtro por expressão regular sobre username, action, zone e details, com a entrada escapada antes da composição da consulta.
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
Três decisões de modelagem merecem destaque:
A trilha fica no MongoDB, e não no PowerDNS. Isso decorre do RNF01. O PowerDNS é a fonte da verdade sobre o estado das zonas, mas não registra quem alterou o quê, pois sua API autentica a aplicação por chave, e não o operador. Como a identidade do operador só existe na camada de operação, é nela que precisa ser persistida.
O registro é descritivo, não diferencial. Grava-se uma descrição do efeito, sem o estado anterior e posterior do RRset. Isso basta para atribuir responsabilidade e reconstruir a sequência de eventos, mas não para restaurar um estado anterior. A limitação também restringe o que a avaliação poderá afirmar sobre a QP4: a trilha informa quem alterou o quê e quando, mas não permite desfazer a alteração.
A gravação é assíncrona e desacoplada da operação. O controlador dispara a inserção em goroutine após responder ao cliente, e a falha de gravação não aborta a operação nem é comunicada ao operador. A decisão prioriza a disponibilidade da operação sobre a completude da trilha, coerente com o RNF03, mas se afasta da recomendação de Kent e Souppaya (2006) quanto à confiabilidade do registro para fins de auditoria. Essa restrição é detalhada no Capítulo 6, e a gravação síncrona com confirmação de escrita figura como proposta de trabalho futuro.
3.8 SÍNTESE
O capítulo especificou a metodologia e os cinco modelos que orientam a implementação. Os requisitos vieram da taxonomia do Quadro 1, o que mantém a rastreabilidade entre cada mecanismo e a classe de erro que o justifica.
Os modelos resultantes são: sessão opaca no servidor, revogável e expirada por inatividade, com senha derivada por bcrypt; RBAC híbrido com papel global e permissão de leitura ou escrita por par (usuário, zona), decidido no backend por requisição; normalização do RDATA como estratégia primária, com validação restrita aos casos sem conversão determinística; sincronização de reversos assimétrica entre criação otimista e remoção guardada por verificação de referência; e auditoria descritiva em coleção própria.
Três decisões foram registradas com suas limitações, por se afastarem de princípios que o próprio trabalho invoca: a escrita derivada em zona fora do escopo do operador, o melhor esforço das rotinas de reverso e o desacoplamento entre operação e registro de auditoria. As implicações dessas escolhas de projeto são discutidas no balanço crítico do Capítulo 6.
O Capítulo 4 documenta a construção do artefato a partir desses modelos.

4 DESENVOLVIMENTO DA PLATAFORMA SANCHEZDNS
Este capítulo documenta a construção do SanchezDNS a partir dos modelos do Capítulo 3. O produto é uma instanciação, o tipo de artefato que Gregor e Hevner (2013) situam no primeiro nível de contribuição em Design Science Research: um sistema em funcionamento, do qual se extraem depois os princípios de projeto generalizáveis. O texto descreve as decisões de engenharia que materializam cada modelo e registra os dois pontos em que a construção obrigou a revisar um requisito.
4.1 ARQUITETURA GERAL E PILHA TECNOLÓGICA
4.1.1 Arquitetura desacoplada
A plataforma tem três camadas com responsabilidades separadas. A camada de apresentação é uma aplicação Nuxt que roda no navegador e em um servidor próprio, executado pelo Bun. O intermediário é um serviço HTTP escrito em Go, que autentica o operador, decide a autorização e traduz cada pedido em chamadas à interface de programação do servidor autoritativo. A camada de dados é o PowerDNS Authoritative Server, único detentor dos dados de zona.
Um proxy reverso Caddy recebe as requisições na porta 80 e separa os dois destinos: o prefixo /go vai para o serviço Go na porta 8080 e todo o restante vai para o servidor Nuxt na porta 3000. O navegador conversa com uma única origem, o que dispensa a configuração de compartilhamento de recursos entre origens e permite emitir o cookie de sessão para o mesmo domínio que serve a interface.
A separação atende ao RNF01 e ao RNF02. O RNF01 aparece na ausência de qualquer cópia local dos registros: o serviço Go não guarda zonas nem registros, e cada tela lê os dados do PowerDNS no momento em que é montada. O RNF02 aparece na assimetria entre as duas camadas. O frontend tem um middleware de rota que esconde páginas conforme o perfil do usuário, mas essa verificação só orienta a navegação. Quem autoriza é o serviço Go, a cada requisição, e uma chamada feita fora da interface passa pela mesma checagem. A Figura 3 mostra as três camadas e o caminho de uma requisição.
Figura 3 – Arquitetura da plataforma SanchezDNS e fluxo de uma requisição

Fonte: elaborado pelo autor (2026).
4.1.2 Escolhas de linguagem e de biblioteca
O backend usa Go com o framework Gin. Três propriedades da linguagem pesaram na escolha. A tipagem estática permite declarar como estruturas os payloads trocados com o PowerDNS, de modo que um campo errado falhe na compilação. As goroutines dão concorrência sem biblioteca externa, recurso usado na gravação da trilha de auditoria descrita na seção 4.6.1. A compilação para binário único, sem dependência de bibliotecas do sistema, obtida com CGO_ENABLED=0, reduz o serviço a um arquivo copiado para o contêiner.
As chamadas HTTP de saída, tanto ao PowerDNS quanto aos resolvedores públicos, usam a biblioteca Resty, que oferece repetição automática de tentativas e desserialização direta em estrutura.
O frontend usa Nuxt 4 sobre Vue 3, com a biblioteca de componentes Nuxt UI e os estilos utilitários do Tailwind CSS. A validação de formulários é feita com Valibot, e os esquemas ficam em composables reaproveitados pelas telas. O módulo nuxt-security define os cabeçalhos de segurança da resposta, entre eles a política de segurança de conteúdo. O Quadro 10 reúne a pilha.
Quadro 10 – Pilha tecnológica da plataforma.
Camada
Tecnologia
Papel no artefato
Apresentação
Nuxt 4, Vue 3, Nuxt UI, Tailwind CSS
Telas, decomposição de campos por tipo de registro e validação no cliente
Validação de formulário
Valibot
Esquemas de campos obrigatórios por tipo, reavaliados no backend
Cabeçalhos de segurança
nuxt-security
Política de conteúdo e cabeçalhos de resposta
Intermediário
Go, Gin
Autenticação, autorização por zona, normalização e automações
Cliente HTTP
Resty
Chamadas ao PowerDNS e aos resolvedores públicos
Estado de aplicação
MongoDB
Sessões, cadastros, solicitações, permissões por zona e auditoria
Cache de sessão
Redis
Leitura de sessão sem consulta ao banco a cada requisição
Arquivos
Armazenamento compatível com S3
Fotos de perfil dos operadores
Dados de zona
PowerDNS Authoritative Server
Fonte única da verdade
Proxy reverso
Caddy
Roteamento entre o serviço Go e o servidor Nuxt

Fonte: elaborado pelo autor (2026).
4.1.3 Persistência híbrida
O MongoDB guarda o estado que pertence à camada de operação, distribuído em cinco coleções. A coleção sessions mantém as sessões, com índice único sobre o identificador; cadastros guarda os operadores aprovados, com o nível global e o resumo da senha; solicitacoes guarda os pedidos de cadastro pendentes e o desfecho de cada um; users guarda um documento por zona, com os vetores de leitura e de escrita; e logs guarda a trilha de auditoria.
O Redis funciona como cache de sessão com expiração de uma hora. A validação consulta o Redis primeiro e recorre ao MongoDB quando a chave não está lá. Isso evita uma leitura de banco por requisição sem transferir a autoridade para o cache: a ausência da chave no Redis é tratada como cache frio, nunca como sessão inválida.
As fotos de perfil ficam em um armazenamento de objetos compatível com o protocolo S3 e são servidas por um endpoint próprio, que devolve o binário e o tipo de conteúdo. Nenhum desses armazenamentos contém dados de zona, o que preserva o RNF01.
4.2 INTEGRAÇÃO COM O POWERDNS AUTHORITATIVE SERVER
4.2.1 Cliente e autenticação da interface de programação
Toda comunicação com o servidor autoritativo passa por um cliente construído em um ponto único da aplicação. Ele lê o endereço base da variável de ambiente DNS_HOST, envia a chave de acesso no cabeçalho X-API-Key, aceita resposta em JSON, tem tempo limite de trinta segundos e repete a chamada até duas vezes em caso de falha. O identificador do servidor vem da variável DNS_SERVER_ID e compõe o caminho das rotas, que seguem o padrão /api/v1/servers/{servidor}/zones.
Esse desenho tem uma consequência já tratada no Capítulo 3: a interface do PowerDNS autentica a aplicação por chave, e não o operador que iniciou a ação. Para o servidor autoritativo, todas as alterações vêm do mesmo cliente. A identidade de quem alterou existe apenas na camada de operação, e por isso a trilha de auditoria fica no MongoDB, conforme a decisão registrada na seção 3.7. O Quadro 11 relaciona as rotas consumidas.
Quadro 11 – Rotas do PowerDNS consumidas pela plataforma.
Rota
Método
Uso no artefato
/zones
GET
Listagem de zonas para a tela e para a busca da zona reversa mais específica
/zones
POST
Criação de zona com o metadado de derivação do serial
/zones/{zona}
GET
Leitura dos conjuntos de registros e do SOA
/zones/{zona}
PUT
Aplicação dos parâmetros de negação autenticada de existência
/zones/{zona}
PATCH
Inclusão, alteração e remoção de conjuntos de registros e gravação dos campos do SOA
/zones/{zona}
DELETE
Remoção de zona
/zones/{zona}/cryptokeys
POST
Geração da chave de assinatura de chaves
/zones/{zona}/cryptokeys
GET
Leitura das chaves e do resumo criptográfico para exibição
/zones/{zona}/rectify
PUT
Recálculo dos dados derivados da zona após a criação
/statistics
GET
Indicadores do servidor autoritativo exibidos no painel

Fonte: elaborado pelo autor (2026).
4.2.2 Criação de zona e delegação do serial
O payload de criação declara três campos: o nome da zona, o modo de operação Native e o metadado soa_edit_api com o valor DEFAULT. O modo Native delega a replicação à camada de armazenamento do backend e dispensa as transferências clássicas de zona, conforme a seção 2.1.6. O metadado implementa o RF01: a partir dele, o próprio PowerDNS reescreve o serial a cada alteração recebida pela interface de programação, e o campo não aparece em nenhuma tela.
O controlador de criação executa as etapas em sequência, e qualquer falha interrompe a operação com resposta de erro e mensagem específica. Primeiro valida o tipo declarado (zona direta, reversa IPv4 ou reversa IPv6) e confere se o domínio informado corresponde a ele. Depois cria a zona, gera a chave de assinatura, aplica os parâmetros de negação autenticada de existência, grava os campos do SOA informados pelo administrador e chama a retificação da zona. A retificação recalcula os dados derivados que o PowerDNS mantém para responder às consultas e roda por último porque as etapas anteriores alteram o conteúdo da zona.
4.2.3 Assinatura e verificação da cadeia de confiança
A geração de chave envia à rota de chaves criptográficas da zona um payload com a chave ativa e do tipo ksk. Em seguida, uma atualização da zona grava o registro de parâmetros NSEC3 com o valor 1 0 0 -, que fixa zero iterações adicionais e dispensa o salt, conforme a recomendação da RFC 9276 discutida na seção 2.1.5. As duas chamadas ficam no mesmo fluxo da criação, o que atende ao RF02: não há etapa posterior que o operador possa esquecer.
Para o RF03, a plataforma também confere o registro DS publicado na zona pai, por meio da interface de consulta sobre HTTPS dos resolvedores públicos. A consulta vai primeiro ao resolvedor da Google e, se falhar, ao da Cloudflare, com pedido de dados de validação e tempo limite de cinco segundos por tentativa. O resultado é comparado com os resumos das chaves da própria zona, depois de normalizar caixa e espaços, e classificado em quatro estados: publicado e conferindo, ausente no pai, divergente e indisponível para consulta. A tela apresenta o estado e, no caso divergente, explica que resolvedores validantes deixam de resolver o domínio até o registro ser corrigido no registrário. A flag de dados autenticados da resposta também é registrada e indica se o resolvedor consultado validou a cadeia.
4.3 AUTENTICAÇÃO E CONTROLE DE ACESSO
4.3.1 Ciclo de vida da sessão e derivação da senha
A autenticação bem-sucedida gera um identificador de sessão com trinta e dois bytes lidos do gerador criptográfico do sistema operacional e codificados em base64 sem preenchimento. O documento gravado na coleção de sessões registra o identificador, o endereço de correio eletrônico, o estado de atividade, os instantes de criação e de último acesso, o endereço IP, o sistema operacional e o navegador extraídos do cabeçalho de agente de usuário, a localização aproximada e o motivo de encerramento. A localização vem de um serviço externo de geolocalização por IP com tempo limite de cinco segundos, e a falha dessa consulta devolve cadeia vazia sem interromper o login.
O cookie é emitido sempre com o atributo que impede sua leitura por script, e com o atributo que o restringe ao canal cifrado apenas no ambiente de produção. O controlador de login recusa a tentativa quando já existe cookie de sessão e, antes de conferir a senha, verifica se o solicitante tem cadastro pendente de aprovação; nesse caso, responde que o cadastro está em análise. Credencial inválida devolve a mesma mensagem para endereço inexistente e para senha errada, o que impede a enumeração de usuários.
As senhas são derivadas com bcrypt (PROVOS; MAZIÈRES, 1999) no fator de custo padrão da biblioteca e gravadas com um prefixo que identifica o esquema. A verificação usa a comparação da própria biblioteca, e resumos marcados com o esquema tradicional de crypt são recusados sem tentativa de conferência.
A expiração é por inatividade. Cada requisição autenticada atualiza o instante de último acesso no MongoDB e renova a chave no Redis. Uma rotina agendada roda a cada dez minutos, seleciona as sessões ativas cujo último acesso é anterior a sete dias, marca-as como inativas com o motivo registrado e apaga as chaves correspondentes do cache.
4.3.2 Middlewares de sessão e de administração
O middleware de sessão lê o cookie, resolve a sessão pelo Redis, com recurso ao MongoDB, obtém o nível global do cadastro e injeta no contexto da requisição o identificador de sessão, o endereço de correio eletrônico e o nível. Quando a sessão não resolve, ele invalida o cookie na própria resposta e recusa o pedido. A implementação atual usa o código 400 nesse caso, embora a especificação HTTP reserve o 401 para falha de autenticação (FIELDING; NOTTINGHAM; RESCHKE, 2022).
O middleware administrativo protege o grupo de rotas reservadas ao perfil global e recusa com 403 qualquer requisição cujo contexto não traga esse nível. Ele cobre a criação e a remoção de zonas, a alteração do SOA, a aplicação dos parâmetros de negação autenticada de existência, a gestão de operadores, a apreciação de solicitações de cadastro e a consulta à auditoria.
4.3.3 Resolução da permissão por zona
As rotas de zona não passam por middleware, porque a decisão depende do parâmetro de zona da requisição e não apenas do contexto da sessão. Cada controlador chama a função de resolução antes de executar a operação. A função devolve o nível administrativo quando o perfil global é de administrador. Nos demais casos, lê o documento da zona na coleção de permissões e verifica a pertinência ao vetor de escrita e depois ao de leitura, comparando endereços sem diferenciar caixa e ignorando espaços nas extremidades. A ausência do documento da zona é tratada como ausência de permissão, sem erro.
A listagem de zonas aplica a mesma resolução a cada item antes de montar a resposta, e um operador sem vínculo recebe uma lista vazia, sem os nomes das zonas que não pode acessar. A implementação atual também exclui da listagem duas zonas por nome fixo, resíduo do ambiente de testes que ainda não foi removido do controlador.
4.4 NORMALIZAÇÃO E INTERFACE PREVENTIVA
4.4.1 Decomposição dos tipos compostos na interface
Os tipos de registro com conteúdo estruturado são editados em campos separados, e não como cadeia única. O formulário oferece campo próprio para a prioridade do MX; para a prioridade, o peso, a porta e o alvo do SRV; e para a prioridade de serviço, o nome de destino e os parâmetros do HTTPS. A obrigatoriedade depende do tipo escolhido: o esquema de validação exige a prioridade, o peso, a porta e o alvo apenas quando o tipo é SRV, e a prioridade de serviço e o nome de destino apenas quando é HTTPS. O tempo de vida tem mínimo de sessenta segundos.
Cada tipo carrega na interface uma frase curta que descreve sua função, como a indicação de que o PTR faz a resolução do endereço para o nome ou de que o CAA define quais autoridades podem emitir certificados. A frase serve para reduzir o engano de planejamento descrito por Reason (1990), que a validação sintática não alcança: um registro pode estar sintaticamente correto e ser do tipo errado para a intenção do operador.
4.4.2 Normalização no backend
A validação no cliente serve à experiência de uso, e o controle fica no serviço Go. O mesmo conjunto de exigências é reavaliado ali, e a normalização acontece apenas no backend, imediatamente antes da submissão ao PowerDNS.
A rotina de normalização escolhe a transformação pelo tipo do registro. Para TXT, acrescenta aspas quando o valor não começa por elas. Para CNAME, NS, ALIAS e PTR, acrescenta o ponto final quando ausente. Para MX, acrescenta o ponto final ao destino e prefixa a prioridade lida do campo próprio. Para SRV e HTTPS, serializa os campos separados na ordem exigida pela especificação, substitui o destino vazio do HTTPS pelo ponto e atribui um conjunto padrão de parâmetros quando o operador não informa nenhum. Para CAA, compõe a marcação e a propriedade quando o valor informado não as contém.
A complementação do nome tem três casos. O nome já terminado em ponto é preservado. O nome que já contém a origem da zona recebe apenas o ponto final. O nome relativo recebe a origem e o ponto. A distinção entre o segundo e o terceiro caso evita a duplicação do domínio, erro em que o operador digita o nome completo em um campo que o sistema supõe relativo.
Para A e AAAA, o valor é interpretado como endereço antes de qualquer outra transformação, e um valor não interpretável interrompe a operação em vez de ser gravado como texto.
4.4.3 Segundo retorno do ciclo: exibir o valor publicado
O desenho inicial mantinha na tela, depois da submissão, o valor digitado pelo operador. A avaliação formativa da tela de registros mostrou o problema: como a normalização acontece no servidor, o operador terminava a operação vendo o que escreveu, e não o que foi publicado. No caso do ponto final a diferença é de um caractere, e o operador não tinha motivo para desconfiar dela. Nesse arranjo, a automação transferia para o sistema uma decisão que o operador supunha ter tomado, efeito que a questão QP2 antecipa.
A correção exigiu a transformação inversa da normalização. A plataforma tem um conjunto de funções que leem os conjuntos de registros publicados e decompõem o conteúdo de volta nos campos do formulário: a prioridade do MX é separada do host, os quatro campos do SRV são extraídos por posição e o conteúdo do HTTPS é dividido em prioridade, destino e parâmetros. Com essa transformação, a tela passou a recarregar os registros pela rota de leitura depois de cada escrita, e essa rota lê o PowerDNS. O valor apresentado é sempre o que está publicado, e uma divergência entre a intenção e o resultado aparece na própria tela em que a operação foi feita.
4.5 SINCRONIZAÇÃO DO CICLO DE VIDA DOS REGISTROS REVERSOS
4.5.1 Derivação do nome e escolha da zona
A derivação do nome reverso tem dois caminhos. Para IPv4, os quatro octetos são invertidos e concatenados sob a raiz reservada. Para IPv6, o endereço é primeiro expandido para a forma plena, com os oito grupos de dezesseis bits formatados em quatro dígitos hexadecimais cada; a cadeia de trinta e dois dígitos resultante é então percorrida do fim para o começo, um dígito por rótulo. A expansão prévia é indispensável porque a notação abreviada suprime zeros e não admite conversão direta em rótulos.
A escolha da zona de destino percorre as zonas existentes no servidor e seleciona, entre as que são sufixo do nome derivado, a de maior comprimento. A comparação usa o sufixo reservado da família do endereço, o que evita confundir uma zona IPv4 com uma IPv6. A regra do sufixo mais longo acomoda a delegação de faixas menores que o octeto, conforme a RFC 2317 discutida na seção 3.6.1. Quando nenhuma zona corresponde, a função retorna sem erro e nenhuma escrita acontece.
A criação segue essa sequência e termina com uma verificação: se já existe conjunto de registros PTR para o nome derivado, a rotina não o sobrescreve. Assim se preserva a decisão anterior de quem administra a zona reversa, e o estado resultante fica detectável pela rotina de reversos ausentes.
A alteração de endereço é executada em duas operações. A primeira grava o PTR do endereço novo. A segunda avalia a remoção do PTR antigo e só ocorre quando os dois endereços diferem. Com a criação antes da remoção, uma falha na segunda etapa deixa dois reversos, um deles obsoleto, em vez de deixar o endereço sem reverso.
4.5.2 Primeiro retorno do ciclo: guarda de referência na remoção
A implementação direta do RF08 removia o PTR sempre que o registro A ou AAAA de origem era removido. A verificação do comportamento mostrou que esse desenho produz um erro novo quando mais de um nome aponta para o mesmo endereço, situação corriqueira em servidor que hospeda vários serviços. A remoção de um dos nomes apagava o reverso ainda referenciado pelos demais, e o resultado era a divergência entre registros diretos e reversos que o próprio requisito pretende combater.
A resposta foi uma guarda de referência executada antes de qualquer remoção de PTR. A rotina percorre as zonas do servidor, descarta as que terminam nos sufixos reservados de resolução reversa, lê cada zona direta e compara com o endereço em questão os endereços de todos os conjuntos A ou AAAA. A comparação é feita sobre o endereço interpretado, para que formas diferentes de escrever o mesmo endereço IPv6 sejam reconhecidas como iguais. Havendo qualquer referência remanescente, a remoção é suprimida. A mesma guarda foi aplicada à etapa de remoção da operação de alteração.
O custo é uma leitura por zona direta a cada remoção de reverso. Em instalações com muitas zonas isso multiplica as chamadas ao servidor autoritativo, e a otimização por um índice invertido de endereços permanece em aberto. A escolha de pagar esse custo decorre da assimetria discutida na seção 3.6.3: um reverso a mais é detectável e corrigível, enquanto um reverso removido indevidamente quebra serviços de imediato.
4.5.3 Melhor esforço, tempo limite e operação em lote
Cada rotina de reverso roda com contexto de tempo limite próprio de seis segundos, derivado do contexto da requisição, para que a lentidão do servidor autoritativo não trave a operação do operador na zona direta, conforme o RNF03. A falha é registrada no log da aplicação e a operação principal se completa.
A rotina de reversos ausentes percorre os registros de endereço da zona direta, deriva o nome reverso de cada um, descarta as repetições pelo nome derivado, localiza a zona reversa correspondente e consulta quais desses nomes ainda não têm PTR. A criação em lote recebe a seleção do operador, filtra a lista pelos itens escolhidos e submete cada PTR em requisição própria. A primeira falha interrompe o processamento e devolve erro, preservando o que já foi gravado, de modo que a operação não é transacional sobre o conjunto selecionado. O estado parcial é convergente, porque a reexecução recalcula os ausentes a partir do estado corrente das zonas.
A rotina de órfãos faz o caminho inverso. Ela monta o conjunto de todos os endereços presentes em registros A ou AAAA das zonas diretas, percorre os conjuntos PTR da zona reversa indicada, converte cada nome reverso de volta em endereço e relaciona os que não constam desse conjunto, com o nome, o endereço e o alvo atual. A remoção assistida opera sobre essa lista.
4.6 TRILHA DE AUDITORIA
4.6.1 Gravação assíncrona
A gravação recebe quatro dados: a zona ou o recurso afetado, o endereço de correio eletrônico do autor, obtido do contexto da sessão, o código da ação e a descrição do efeito. O instante é gerado no momento da gravação, em formato de data e hora com fuso horário. Os códigos de ação são os doze do Quadro 9, estáveis para permitir agregação, e a descrição em texto livre serve à leitura humana.
O controlador dispara a inserção em goroutine e responde ao cliente sem esperar a gravação, cujo erro é descartado. A operação auditada se completa mesmo quando o registro falha, e o operador não é avisado. A decisão prioriza a disponibilidade da operação sobre a completude da trilha, coerente com o RNF03, e se afasta da recomendação de Kent e Souppaya (2006) quanto à confiabilidade do registro. A consequência para o alcance das conclusões consta na seção 6.2, e a gravação com confirmação de escrita figura entre os trabalhos futuros.
4.6.2 Consulta e filtragem
A rota de consulta é restrita ao perfil administrativo. Ela aceita os parâmetros de página e de limite, com valores padrão de um e de dez e correção dos valores inválidos, converte a página em deslocamento e ordena os resultados pelo instante, do mais recente para o mais antigo. A resposta traz a página pedida e o total de documentos que satisfazem o filtro, o que permite à interface montar a paginação.
O filtro recebe uma cadeia única e a aplica como expressão regular, sem diferenciar caixa, sobre quatro campos ao mesmo tempo: autor, código de ação, zona e descrição. A entrada do usuário passa por escape antes de compor a consulta, e os caracteres com significado especial em expressão regular são tratados como literais. Sem esse cuidado, uma cadeia de busca poderia alterar a semântica da consulta ou provocar avaliação custosa no banco.
4.7 EMPACOTAMENTO E IMPLANTAÇÃO
A imagem de execução é construída em três estágios. O primeiro instala as dependências do frontend com o Bun e gera a saída de produção do Nuxt. O segundo compila o serviço Go sem informações de depuração e sem vínculo a bibliotecas do sistema. O terceiro parte de uma imagem Debian com o Bun, copia a saída do Nuxt, o binário do Go e o executável do Caddy, e grava o arquivo de configuração do proxy e o script de entrada. Os dois primeiros estágios usam cache de dependências, o que reduz o tempo das reconstruções.
O contêiner executa três processos sob um init mínimo, responsável por encaminhar sinais e recolher processos órfãos: o serviço Go, o servidor Nuxt e o Caddy. O script de entrada instala um tratador de sinais e encerra o grupo inteiro quando qualquer um dos três termina, o que evita que o contêiner continue no ar com uma das camadas fora de operação. A verificação de saúde consulta a rota de checagem do serviço Go através do proxy a cada trinta segundos, com período inicial de tolerância.
A composição de serviços declara, além da aplicação, o MongoDB, o Redis, o armazenamento de objetos e um serviço auxiliar de inicialização que cria o compartimento de arquivos na primeira subida. A configuração vem de variáveis de ambiente, entre elas o endereço e a chave da interface do PowerDNS, o identificador do servidor autoritativo, as cadeias de conexão do MongoDB e do Redis, as credenciais do armazenamento de objetos e o endereço público da aplicação. O PowerDNS não integra a composição, porque o artefato pressupõe uma instância já em operação, conforme a delimitação da seção 1.5.
5 AVALIAÇÃO E RESULTADOS
Considerações Iniciais do Capítulo 5:
Apresentar a execução rigorosa do plano de avaliação FEDS desenhado na seção 2.6. Demonstrar a eficácia e a eficiência do artefato confrontando métricas quantitativas e qualitativas coletadas perante os critérios de êxito pré-fixados no Quadro 3.
5.1 DEMONSTRAÇÃO E CONFORMIDADE TÉCNICA (EPISÓDIOS E0 E E1)
5.1.1 Rastreabilidade e Cobertura da Caracterização (Episódio E0)
O que escrever: Análise da suficiência da matriz requisito-problema, comprovando que todas as classes elimináveis por construção do Quadro 1 foram cobertas pelo artefato.
Citação interna anterior: Previsto no episódio E0 da §2.6.2 e no critério de êxito de QP1 no Quadro 3 (§2.6.4).
5.1.2 Testes em Bancada de Laboratório e Casos de Teste Negativos (Episódio E1)
O que escrever: Resultados dos testes automatizados tentando injetar dados malformados na interface (FQDN relativo, sintaxe MX/TXT quebrada, estouro de serial). Relato da taxa de bloqueio (esperado: 100%).
Citação interna anterior: Vinculado ao episódio E1 da §2.6.2, atendendo ao marco de demonstração DSRM anunciado na §3.2 (Quadro 4, item 4).
5.1.3 Verificação com Oráculos Externos Independentes (Zonemaster e PowerDNS)
O que escrever: Execução de consultas diretas ao PowerDNS e varreduras com a suíte de casos de teste do Zonemaster sobre zonas configuradas pelo SanchezDNS, demonstrando ausência de novas reprovações.
Citação interna anterior: Previsto como metodologia de oráculo externo na §2.5.7 e §2.6.2 (E1).
5.1.4 Testes Negativos de Autorização e Isolamento de Escopo
O que escrever: Ensaios simulando a tentativa de escrita de registros PTR em zonas reversas não autorizadas ao operador, comprovando que o artefato não replica a falha do Technitium (v15.5).
Citação interna anterior: Mapeado na justificativa da QP3 na §2.5.3 e nos critérios do Quadro 3 (E1).
5.2 EXPERIMENTO CONTROLADO COM OPERADORES (EPISÓDIO E2)
5.2.1 Protocolo Experimental, Participantes e Amostra
O que escrever: Perfil dos 12 participantes (profissionais/estudantes de redes), desenho intrassujeitos contrabalanceado e roteiro das seis tarefas operacionais padronizadas.
Citação interna anterior: Especificado detalhadamente no episódio E2 da §2.6.2 e amparado na bibliografia de Sauro & Lewis (2016).
5.2.2 Análise Comparativa da Taxa de Erros Operacionais
O que escrever: Dados quantitativos da contagem de erros sintáticos e operacionais cometidos na edição manual (CLI/BIND) versus no SanchezDNS, categorizados pela taxonomia do Quadro 1. Aplicação de testes não paramétricos pareados (p < 0,05).
Citação interna anterior: Métrica definida para QP4 no Quadro 3 (§2.6.4).
5.2.3 Análise Comparativa do Tempo de Execução e Eficiência
O que escrever: Comparação dos tempos medianos de conclusão de cada tarefa, analisando se o SanchezDNS reduz ou mantém o tempo de operação enquanto zera os erros.
Citação interna anterior: Métrica de eficiência definida para QP4 no Quadro 3 (§2.6.4).
5.3 ESTUDO DE CASO EM AMBIENTE REAL DE PRODUÇÃO (EPISÓDIO E3)
5.3.1 Caracterização do Ambiente Organizacional e Zonas Avaliadas
O que escrever: Descrição da organização parceira, parque de servidores PowerDNS e volumetria das zonas e registros em produção (dados pseudonimizados).
Citação interna anterior: Metodologia de estudo de caso único incorporado de Yin (2018) definida na §2.6.3.
5.3.2 Fase 1: Coleta da Linha de Base (Baseline de 4 Semanas)
O que escrever: Indicadores operacionais do método anterior (chamados abertos por falhas de DNS, taxa de retrabalho em 72h, varredura inicial do Zonemaster e contagem de reversos órfãos).
Citação interna anterior: Procedimento detalhado na §2.6.3 (Fase 1).
5.3.3 Fase 2: Intervenção, Migração e Treinamento
O que escrever: Registro do processo de implantação da plataforma, migração de permissões e horas de treinamento dispensadas aos operadores.
Citação interna anterior: Detalhado na §2.6.3 (Fase 2).
5.3.4 Fase 3: Observação Operacional (8 Semanas) e Triangulação de Dados
O que escrever: Análise da convergência de evidências entre os logs da plataforma, chamados de TI, varreduras periódicas com Zonemaster e entrevistas com operadores.
Citação interna anterior: Critério de triangulação de evidências de Yin (2018) estabelecido na §2.6.3.
5.4 AVALIAÇÃO DE USABILIDADE PERCEBIDA (EPISÓDIO E4)
5.4.1 Aplicação da Escala SUS (System Usability Scale)
O que escrever: Coleta das respostas dos 10 itens do questionário SUS aplicados ao final de E2 e E3.
Citação interna anterior: Previsto no episódio E4 da §2.6.2 e referenciado em Brooke (1996).
5.4.2 Análise Psicométrica e Comparação com a Linha de Corte Industrial
O que escrever: Cálculo do escore global consolidado (0 a 100), análise das subescalas de usabilidade e facilidade de aprendizado e enquadramento nas faixas adjetivas de Bangor, Kortum & Miller (2008) (esperado: SUS ≥ 68).
Citação interna anterior: Critério de êxito estabelecido no Quadro 3 (§2.6.4).
6 CONCLUSÕES
Considerações Iniciais do Capítulo 6:
Sintetizar o substrato teórico e prático da pesquisa, demonstrando como o ciclo DSRM e os resultados do Capítulo 5 responderam às questões QP principal e QP1 a QP4, além de conduzir um balanço crítico honesto das limitações do artefato e propor os trabalhos futuros.
6.1 SÍNTESE DAS CONTRIBUIÇÕES E RESPOSTAS ÀS QUESTÕES DE PESQUISA
6.1.1 Contribuição Prática: A Instanciação Funcional (Nível 1 em DSR)
O que escrever: Resumo da entrega da plataforma SanchezDNS como artefato de software livre eficiente, seguro e desacoplado para ambientes com PowerDNS.
Citação interna anterior: Enquadramento de contribuição estabelecido na §2.4.1 (Gregor & Hevner, 2013).
6.1.2 Contribuição Teórica e Princípios de Design Emergentes (Nível 2 em DSR)
O que escrever: Formalização dos princípios de design preventivo para interfaces operacionais de infraestruturas críticas (validação em tempo de entrada, automação de dados deriváveis com checagem de referência e RBAC por escopo).
Citação interna anterior: Fundamentado em Gregor & Hevner (2013) na §2.4.1 e na síntese da §3.8.
6.1.3 Resposta Consolidada às Questões de Pesquisa (QP, QP1, QP2, QP3 e QP4)
O que escrever: Apresentar um parágrafo conclusivo e direto respondendo objetivamente a cada uma das questões formuladas na seção 1.3 com base nas evidências empíricas de E0 a E4.
6.2 ANÁLISE CRÍTICA DAS LIMITAÇÕES DO ARTEFATO
(Atenção: esta seção reúne exatamente todas as quatro concessões e trade-offs técnicos assumidos no Capítulo 3!)
6.2.1 Afrouxamento do Menor Privilégio na Escrita Derivada em Zona Externa
O que escrever: Discutir criticamente o fato de a automação do PTR gravar sob in-addr.arpa ou ip6.arpa usando credenciais de sistema, permitindo que a escrita atravesse fronteiras de escopo sem que o operador tenha permissão explícita na zona reversa.
Citação interna anterior: MUITO IMPORTANTE! Prometido explicitamente na §3.4.5 e na síntese da §3.8.
6.2.2 Semântica de Melhor Esforço e Riscos de Falha Silenciosa no Registro Derivado
O que escrever: Avaliar o risco do RNF03 (a operação direta retornar HTTP 200 de sucesso enquanto a gravação do reverso sofre timeout ou falha), gerando divergência assíncrona imperceptível se o operador não auditar a ferramenta de triagem.
Citação interna anterior: MUITO IMPORTANTE! Remetido na §3.6.4 e na §3.8.
6.2.3 Auditoria Descritiva Assíncrona e Ausência de Registro Diferencial (Rollback)
O que escrever: Explicar a limitação de a auditoria não armazenar o estado anterior/posterior (diff) do registro de recurso, permitindo atribuir responsabilidade histórica, mas impedindo rotinas automatizadas de desfazimento (undo/rollback), além do risco probabilístico do desacoplamento assíncrono em goroutines.
Citação interna anterior: MUITO IMPORTANTE! Remetido na §3.7 e na §3.8.
6.2.4 Ameaças Remanescentes à Validade da Pesquisa
O que escrever: Retomar as fronteiras de validade externa (estudo circunscrito a uma única organização parceira e ao ecossistema PowerDNS).
Citação interna anterior: Vinculado à análise de ameaças à validade da §2.6.4.
6.3 TRABALHOS FUTUROS
6.3.1 Auditoria Diferencial e Mecanismos Nativos de Rollback de Zonas
O que escrever: Proposta de evolução para persistir instantâneos antes/depois do RRset, viabilizando reversão de registros corrompidos com um clique.
Citação interna anterior: Sugerido no fechamento da §3.7.
6.3.2 Persistência Síncrona e Garantias Transacionais entre Zonas Diretas e Reversas
O que escrever: Modelagem de transações distribuídas (estilo protocolo two-phase commit ou saga) para garantir atomicidade real entre zonas diretas e reversas sem degradar o tempo de resposta.
Citação interna anterior: Referenciado na discussão sobre semântica transacional do Knot DNS (§2.5.6) e §3.7.
6.3.3 Extensão do Modelo para Múltiplos Backends Autoritativos (BIND 9 e Knot DNS)
O que escrever: Criação de adaptadores abstratos para permitir que o plano de controle do SanchezDNS atue sobre motores BIND 9 e Knot via soquetes de controle ou arquivos dinâmicos.

REFERÊNCIAS
AKAMAI. Akamai summarizes service disruption (resolved). Akamai Blog, [s. l.], 22 jul. 2021. Disponível em: https://www.akamai.com/blog/news/akamai-summarizes-service-disruption-resolved.

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

FIELDING, R.; NOTTINGHAM, M.; RESCHKE, J. HTTP semantics. RFC 9110. [S. l.]: IETF, 2022. Disponível em: https://www.rfc-editor.org/rfc/rfc9110.

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

PROVOS, N.; MAZIÈRES, D. A future-adaptable password scheme. In: USENIX ANNUAL TECHNICAL CONFERENCE, 1999, Monterey. Proceedings [...]. Berkeley: USENIX Association, 1999. p. 81-91.

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
