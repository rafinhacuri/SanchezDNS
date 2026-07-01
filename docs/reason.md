# 💡 Por que SanchezDNS

Esta página apresenta a **motivação** do projeto: qual problema ele resolve, para quem, e quais princípios guiaram as decisões técnicas. É o ponto de partida conceitual antes das páginas de arquitetura e implementação.

## O problema

Administrar um servidor DNS autoritativo "na mão" é trabalhoso e propenso a erros. Nas ferramentas tradicionais, quem opera precisa:

- editar arquivos de zona ou tabelas de banco diretamente;
- lembrar de formatos exigentes (o ponto final do FQDN, as aspas do TXT, a ordem dos campos do SOA);
- manter os registros **reversos (PTR)** manualmente, em zonas separadas e num formato invertido pouco intuitivo;
- gerenciar DNSSEC com comandos avulsos;
- e, o mais delicado em uma equipe, **controlar quem pode mexer em quê** — normalmente inexistente ou tudo-ou-nada.

Não há, nesse cenário, uma trilha de **auditoria** clara de quem alterou o quê e quando. Um erro de digitação em uma zona pode derrubar um serviço inteiro, e descobrir a causa depois é difícil.

## O estado das ferramentas de mercado

Ao pesquisar o que já existe hoje para operar DNS autoritativo, o cenário reforça a necessidade do projeto: as soluções open-source disponíveis são **datadas, limitadas ou frágeis** para a operação do dia a dia.

### BIND9 — poderoso, mas manual e perigoso

O **BIND9** é o servidor DNS mais tradicional do mundo, mas sua operação continua fundamentalmente **manual**: zonas em arquivos de texto editados à mão e recarregados no servidor. Não há uma camada de administração web moderna de primeira classe, controle de acesso por zona ou auditoria integrada — tudo isso vira responsabilidade de scripts e disciplina da equipe.

O exemplo mais crítico dessa fragilidade é o **serial da zona**. No fluxo do BIND9, **quem edita o arquivo de zona é responsável por incrementar o serial do SOA manualmente**. E isso é uma armadilha silenciosa:

- se você **esquece de incrementar** o serial depois de mudar um registro, os servidores secundários **não percebem que a zona mudou** — eles comparam o serial e concluem que nada foi alterado. A mudança **nunca propaga**, e você acha que publicou algo que na prática não saiu do primário;
- se você **erra o serial para menos** (por exemplo, digita um número menor que o atual, ou "reseta" a numeração), os secundários passam a considerar sua cópia como **mais antiga** e param de aceitar atualizações. Na prática, **a zona congela ou a replicação quebra** — e, dependendo do formato usado, corrigir exige intervenção manual delicada em todos os servidores.

Ou seja: **um único erro de digitação em um número pode derrubar ou travar a propagação de todo o seu DNS**, sem nenhuma mensagem de erro óbvia no momento em que o erro é cometido. É o tipo de problema que só aparece quando já é tarde.

O SanchezDNS elimina essa classe de erro pela raiz: o **serial nunca é digitado por uma pessoa**. Como as zonas são criadas com `soa_edit_api: DEFAULT`, é o **próprio PowerDNS que reescreve o serial automaticamente** a cada alteração feita pela API, no formato correto e sempre crescente (ver [PowerDNS a fundo](/powerdns#registro-resource-record-e-o-soa)). Não há como esquecer de incrementar nem como regredir o número.

### PowerDNS-Admin — a referência que envelheceu

Do lado do PowerDNS, a ferramenta open-source mais conhecida é o **PowerDNS-Admin**. Ele resolveu o problema de dar uma interface web ao PowerDNS, mas é um projeto que **envelheceu**: interface antiquada, base tecnológica datada e manutenção irregular.

### A lacuna

Em resumo, há uma lacuna clara: **falta uma ferramenta de administração de DNS autoritativo que seja ao mesmo tempo moderna, segura e automatizada.** É essa lacuna que o SanchezDNS se propõe a preencher — uma interface atual (Nuxt 4), com controle de acesso por zona, auditoria e automações que removem justamente as classes de erro (como o serial) que tornam a operação tradicional perigosa.

## A proposta

O SanchezDNS é uma **plataforma web de operação** para PowerDNS Authoritative que transforma essas tarefas em um fluxo único, seguro e auditável. Ele parte de três decisões de fundo:

1. **Não reinventar o DNS.** Quem serve o DNS continua sendo o PowerDNS. O SanchezDNS é a camada de operação por cima dele, falando com sua API. Isso evita duplicar a fonte da verdade e aproveita um servidor DNS maduro e robusto.
2. **Automatizar o que é derivável.** Coisas que podem ser calculadas — o PTR reverso de um `A`, o serial do SOA, o FQDN de um nome curto, a formatação de cada tipo de registro — o sistema faz sozinho. Menos digitação, menos erro.
3. **Tornar o acesso explícito.** Cada zona tem permissões de **leitura** e **escrita** por usuário, e há um nível **administrativo** global. Toda mudança relevante vira **log de auditoria**.

## O que ele faz bem

- **Gera e mantém zonas no PowerDNS** de forma direta, já com DNSSEC ativado.
- **Trata zonas normais e reversas** sem exigir procedimentos paralelos — o reverso acompanha o direto automaticamente.
- **Reduz erros de manutenção** completando FQDNs, normalizando valores por tipo e sincronizando registros reversos.
- **Expõe métricas úteis** para a operação diária (zonas, registros, uptime, tráfego).
- **Concentra aprovações, permissões e auditoria** em um lugar só, sem espalhar a lógica por várias ferramentas.

## Público-alvo

O projeto foi pensado para **equipes pequenas e médias que operam seu próprio DNS autoritativo** — provedores, laboratórios, times de infraestrutura — e que precisam de controle de acesso e rastreabilidade sem montar uma solução complexa. Como Trabalho de Conclusão de Curso, ele também serve de estudo prático de integração entre uma aplicação web moderna e um serviço de infraestrutura crítico via API.

## Princípios de projeto

Estes princípios reaparecem em todas as decisões documentadas nas outras páginas:

- **Fonte da verdade única.** As zonas vivem no PowerDNS; as permissões e a auditoria, no MongoDB. O painel nunca guarda uma cópia divergente do DNS.
- **Serviços descartáveis, estado externo.** A aplicação roda em um container sem estado; tudo que precisa persistir vive em serviços dedicados (Mongo, Redis, S3, PowerDNS). Ver [Arquitetura](/architecture).
- **Segurança no backend.** Autenticação, permissão e auditoria são decididas pelo servidor; a interface só reflete o que ele autoriza. Ver [Autenticação](/authentication).
- **Automação de melhor esforço.** As automações (como o PTR reverso) nunca derrubam a operação principal: se falham, registram no log e seguem.

## O que o projeto deliberadamente **não** é

Deixar claro o escopo também é uma decisão de projeto:

- **Não é um resolvedor DNS.** Ele administra um servidor autoritativo, não resolve nomes para clientes.
- **Não gerencia um cluster de servidores pela interface.** O modelo é de **uma instância PowerDNS por ambiente** (definida por `DNS_HOST`/`DNS_SERVER_ID`). Não há painel de múltiplas conexões — a replicação, quando necessária, é resolvida na camada de infraestrutura/banco, não na aplicação.
- **Não substitui o PowerDNS.** Ele depende dele; é uma camada de operação, não um servidor DNS.

## Resumo

O SanchezDNS troca a operação manual e arriscada de um DNS autoritativo por um fluxo **web, automatizado, com controle de acesso por zona e auditoria** — mantendo o PowerDNS como o motor confiável por baixo e adicionando exatamente as camadas que faltam para operar com segurança em equipe.
