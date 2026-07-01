# 🧭 Zonas e Registros

Esta página aprofunda **como o SanchezDNS cria e manipula zonas e registros DNS**, o que cada tipo de registro significa e **por que o backend trata cada um de um jeito específico**. Ela pressupõe os conceitos de RRset, SOA e DNSSEC apresentados em [PowerDNS a fundo](/powerdns).

## Os três grupos de zona

A tela de zonas separa o que existe em três grupos, porque eles têm propósitos diferentes:

- **Zonas normais** — mapeiam **nome → endereço** (o uso comum: `exemplo.com`, `www.exemplo.com`).
- **Zonas reversas IPv4** (`in-addr.arpa`) — mapeiam **IP → nome** para endereços IPv4.
- **Zonas reversas IPv6** (`ip6.arpa`) — o mesmo, para IPv6.

Cada entrada mostra nome, serial, nível de acesso do usuário e status de DNSSEC. O **serial** vem do SOA e permite ver rapidamente se a zona foi alterada recentemente (lembrando que o PowerDNS o gerencia automaticamente).

**Por que separar reversas em uma categoria própria?** Porque uma zona reversa não é editada "à mão" no dia a dia — na maior parte das vezes seus registros PTR são **gerados automaticamente** pelo sistema quando você cria um `A` ou `AAAA` (ver a seção de reverse automático). Agrupá-las deixa claro que são zonas de infraestrutura, não de conteúdo.

## Criação de uma zona

O formulário pede o domínio e os campos do SOA (`Start of Authority`, email do responsável, `refresh`, `retry`, `expire`, `negativeCacheTtl`). Ao submeter, o backend executa **três chamadas à API do PowerDNS em sequência** ([criar-zona.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/zonas/criar-zona.go)):

1. **`POST .../zones`** cria a zona como `Native` com `soa_edit_api: DEFAULT`.
2. **`POST .../zones/{zona}/cryptokeys`** ativa DNSSEC com uma **KSK** ativa.
3. **`PATCH .../zones/{zona}`** grava o RRset `SOA` inicial com os valores do formulário.

Se qualquer passo falhar, o backend retorna o erro exato vindo do PowerDNS (status + corpo), o que facilita diagnosticar problemas de configuração. Ao final, grava um log `create_zone` de forma assíncrona.

O detalhe de **por que `Native`** e **por que DNSSEC por padrão** está justificado em [PowerDNS a fundo](/powerdns#como-o-sanchezdns-usa-a-api-endpoint-por-endpoint).

## Edição do SOA

Editar a zona é, na prática, **regravar o RRset SOA** ([update-soa.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/zonas/update-soa.go)): um único `PATCH` com `changetype REPLACE` sobre o registro `SOA`. O serial continua sendo gerenciado pelo PowerDNS. Só administradores acessam essa operação (rota `PATCH /soa` protegida por `AdminOnly`).

## Os tipos de registro suportados e por que cada um existe

O editor trabalha com estes tipos. Entender cada um ajuda a justificar por que o backend precisa tratá-los de forma diferente:

| Tipo | Para que serve |
|---|---|
| `A` | aponta um nome para um endereço **IPv4** |
| `AAAA` | aponta um nome para um endereço **IPv6** |
| `CNAME` | apelido: aponta um nome para **outro nome** |
| `ALIAS` | como CNAME, mas pode ser usado no apex da zona (raiz do domínio) |
| `MX` | define o **servidor de email** do domínio, com prioridade |
| `NS` | delega uma subzona a outros servidores de nomes |
| `TXT` | texto livre — usado por SPF, DKIM, verificações de propriedade |
| `PTR` | o **reverso**: mapeia IP → nome (vive nas zonas `.arpa`) |
| `SRV` | localiza serviços (host + porta) para protocolos como SIP, XMPP |
| `CAA` | diz **quais autoridades certificadoras** podem emitir certificados para o domínio |
| `HTTPS` | parâmetros de conexão HTTPS (ALPN, etc.) já na resolução do nome |
| `TLSA` | associa um certificado/chave TLS ao nome (DANE) |

## Normalização: por que o valor é ajustado antes de gravar

DNS tem convenções de formatação que, se o usuário tivesse que digitar perfeitamente, gerariam erros constantes. Por isso o backend **normaliza** o valor antes de montá-lo no RRset ([normalize-record-value.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/records/normalize-record-value.go)). Cada regra existe por um motivo concreto:

- **`TXT`** — o conteúdo é envolvido em aspas (`"..."`) se ainda não estiver. O formato de registro TXT exige aspas para delimitar a string; sem elas, o PowerDNS rejeitaria ou interpretaria errado.
- **`CNAME`, `NS`, `ALIAS`, `MX`, `PTR`** — recebem um **ponto final** se não tiverem. Esses tipos apontam para **nomes**, e um nome sem ponto final seria interpretado como relativo à zona. O ponto o torna absoluto (FQDN), evitando que `mail.exemplo.com` vire `mail.exemplo.com.exemplo.com`.
- **`MX`** — a **prioridade** é prefixada ao valor (`10 mail.exemplo.com.`). O formato do MX exige o número de prioridade antes do host; a interface coleta isso em um campo separado e o backend junta.
- **`CAA`** — se o usuário digitou só a autoridade, o backend completa para o formato `0 issue "letsencrypt.org"`. Assim quem não conhece a sintaxe do CAA ainda consegue criar um registro válido.
- **`SRV`** — montado como `prioridade peso porta alvo.`, cada parte vinda de um campo próprio da interface, com ponto final garantido no alvo.
- **`HTTPS`** — montado como `prioridade alvo parâmetros` (por exemplo `1 . alpn=h2`), com um padrão sensato (`alpn=h2`) quando os parâmetros não são informados.

Repare que `HTTPS` e `SRV` **não usam um único campo de valor** na interface — eles têm campos estruturados (`svcPriority`, `targetName`, `svcParams`, `weight`, `port`, `target`, `priority`), porque seu conteúdo é composto por várias partes. Por isso o `AddRecordRequest` no backend tem campos opcionais dedicados a esses tipos.

## Como o nome do registro é completado (FQDN)

Ao inserir, o backend ajusta o **nome** do registro para um FQDN absoluto:

```go
if !strings.HasSuffix(name, ".") {
    if !strings.HasSuffix(name, zone) {
        name = fmt.Sprintf("%s.%s.", name, zone)   // "www" → "www.exemplo.com."
    } else {
        name += "."                                 // já tinha a zona, só falta o ponto
    }
}
```

Ou seja:

- nome **vazio** → usa o **apex** da zona (a raiz do domínio);
- nome **curto** (`www`) → vira `www.exemplo.com.`;
- nome que **já termina com a zona** → só recebe o ponto final.

Isso libera o usuário de digitar o domínio inteiro toda vez, sem risco de duplicar o sufixo.

## O padrão de dois PATCHes na inserção

Como explicado em [PowerDNS a fundo](/powerdns#adicionar-editar-registros-get-patch-o-padrao-de-dois-passos), adicionar um registro envolve **ler o RRset atual, gravar o conjunto completo, reler e então gravar os comentários alinhados**. O motivo, em resumo:

1. **Preservar registros existentes** do mesmo nome+tipo (senão o `REPLACE` os apagaria).
2. **Alinhar comentários por posição** — o PowerDNS pode reordenar os registros ao gravar, e os comentários são posicionais, então o backend relê a ordem canônica antes de anexar cada comentário ao registro certo.

Se o segundo `PATCH` (dos comentários) falhar, a operação retorna erro; se o **reverso** falhar, apenas registra no log e segue — a criação do registro principal não é revertida por causa do reverso.

## Reverse automático — o PTR que se cuida sozinho

Este é o comportamento mais característico do SanchezDNS. Sempre que um registro **`A`** ou **`AAAA`** é criado, editado ou removido, o backend mantém o **PTR reverso** correspondente em sincronia, sem intervenção manual ([ensure-reverse-record.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/records/ensure-reverse-record.go)).

### Por que isso importa

O PTR (reverso) é o que responde "qual nome corresponde a este IP?". Ele é exigido por servidores de email, aparece em logs e diagnósticos, e é fácil de **esquecer** de manter atualizado — porque vive em uma zona separada (`in-addr.arpa`/`ip6.arpa`) e num formato pouco intuitivo. Automatizá-lo elimina uma das fontes de erro mais comuns na operação de DNS.

### Como o nome reverso é construído

**IPv4** — os octetos são invertidos e recebem o sufixo `in-addr.arpa`:

```
192.0.2.10  →  10.2.0.192.in-addr.arpa.
```

**IPv6** — o endereço é expandido para 32 dígitos hexadecimais, cada **nibble** (dígito) é invertido e separado por pontos, com sufixo `ip6.arpa`:

```
2001:db8::1  →  1.0.0.0. ... .8.b.d.0.1.0.0.2.ip6.arpa.
```

### Como o backend escolhe a zona reversa

O backend lista todas as zonas (`GET .../zones`) e procura, entre as que terminam em `in-addr.arpa`/`ip6.arpa`, aquela cujo nome é o **maior sufixo** do nome reverso calculado:

```go
if strings.HasSuffix(fullRevNoDot, zNameNoDot) {
    if len(zNameNoDot) > len(bestZone) {   // vence o sufixo mais longo/específico
        bestZone = z.Name
    }
}
```

**Por que o sufixo mais longo?** Porque pode haver zonas reversas de granularidades diferentes (por exemplo, uma `/16` e uma `/24` da mesma rede). A mais específica (nome mais longo) é a dona correta daquele IP — a mesma lógica de "correspondência mais específica" que o roteamento usa. Se **nenhuma** zona reversa cobrir o IP, o backend simplesmente não faz nada (não é erro: talvez você não gerencie o reverso daquela faixa).

Antes de criar, ele verifica se o PTR já existe naquela zona; se existir, não sobrescreve. Se não, faz um `PATCH` criando o PTR que aponta de volta para o FQDN do registro `A`/`AAAA`.

### O ciclo completo

- **criar** `A`/`AAAA` → cria o PTR correspondente;
- **editar** o valor → atualiza o PTR (remove o antigo, garante o novo);
- **remover** `A`/`AAAA` → remove o PTR junto.

Tudo isso roda com um **timeout curto (6s)** e em caráter de "melhor esforço": se a parte reversa falhar, ela é registrada no log, mas **não derruba** a operação principal sobre o registro direto.

## Permissões na tela de zona

- O botão de gestão de usuários e a edição de SOA só aparecem quando a zona está em nível **`ADMINISTRADOR`** para o usuário.
- Para usuários com nível **`LEITURA`**, a tabela de registros fica **somente leitura**.
- Para criar/editar/remover registros, o usuário precisa de **`ESCRITA`** naquela zona.

O modelo de permissões por zona é detalhado em [Usuários](/users), e a verificação acontece no backend a cada operação — a interface apenas reflete o que o servidor já decide.
