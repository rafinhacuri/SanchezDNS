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

O formulário pede o domínio e os campos do SOA (`Start of Authority`, email do responsável, `refresh`, `retry`, `expire`, `negativeCacheTtl`). Ao submeter, o controlador executa **três chamadas à API do PowerDNS em sequência**, cada uma em uma função própria do pacote `zonas`:

1. **`POST .../zones`** cria a zona como `Native` com `soa_edit_api: DEFAULT` ([insert.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/zonas/insert.go)).
2. **`POST .../zones/{zona}/cryptokeys`** ativa DNSSEC com uma **KSK** ativa ([insert-dnssec.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/zonas/insert-dnssec.go)).
3. **`PATCH .../zones/{zona}`** grava o RRset `SOA` inicial com os valores do formulário ([update-soa.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/zonas/update-soa.go)).

Quem orquestra os três passos e decide a mensagem de cada falha é o controlador ([insert/zona.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/controller/insert/zona.go)) — as funções do pacote `zonas` só executam a chamada e devolvem o erro cru.

Se qualquer passo falhar, o backend retorna o erro exato vindo do PowerDNS (status + corpo), o que facilita diagnosticar problemas de configuração. Ao final, grava um log `create_zone` de forma assíncrona.

O detalhe de **por que `Native`** e **por que DNSSEC por padrão** está justificado em [PowerDNS a fundo](/powerdns#como-o-sanchezdns-usa-a-api-endpoint-por-endpoint).

## Edição do SOA

Editar a zona é, na prática, **regravar o RRset SOA** ([update-soa.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/zonas/update-soa.go)): um único `PATCH` com `changetype REPLACE` sobre o registro `SOA`. O serial continua sendo gerenciado pelo PowerDNS. Só administradores acessam essa operação (rota `PATCH /soa` protegida por `AdminOnly`).

## Os tipos de registro suportados e por que cada um existe

O editor trabalha com estes tipos. Entender cada um ajuda a justificar por que o backend precisa tratá-los de forma diferente:

| Tipo    | Para que serve                                                                    |
| ------- | --------------------------------------------------------------------------------- |
| `A`     | aponta um nome para um endereço **IPv4**                                          |
| `AAAA`  | aponta um nome para um endereço **IPv6**                                          |
| `CNAME` | apelido: aponta um nome para **outro nome**                                       |
| `ALIAS` | como CNAME, mas pode ser usado no apex da zona (raiz do domínio)                  |
| `MX`    | define o **servidor de email** do domínio, com prioridade                         |
| `NS`    | delega uma subzona a outros servidores de nomes                                   |
| `TXT`   | texto livre — usado por SPF, DKIM, verificações de propriedade                    |
| `PTR`   | o **reverso**: mapeia IP → nome (vive nas zonas `.arpa`)                          |
| `SRV`   | localiza serviços (host + porta) para protocolos como SIP, XMPP                   |
| `CAA`   | diz **quais autoridades certificadoras** podem emitir certificados para o domínio |
| `HTTPS` | parâmetros de conexão HTTPS (ALPN, etc.) já na resolução do nome                  |
| `TLSA`  | associa um certificado/chave TLS ao nome (DANE)                                   |

## Normalização: por que o valor é ajustado antes de gravar

DNS tem convenções de formatação que, se o usuário tivesse que digitar perfeitamente, gerariam erros constantes. Por isso o backend **normaliza** o valor antes de montá-lo no RRset ([normalizar-valor.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/records/normalizar-valor.go)). Cada regra existe por um motivo concreto:

- **`TXT`** — o conteúdo é envolvido em aspas (`"..."`) se ainda não estiver. O formato de registro TXT exige aspas para delimitar a string; sem elas, o PowerDNS rejeitaria ou interpretaria errado.
- **`CNAME`, `NS`, `ALIAS`, `MX`, `PTR`** — recebem um **ponto final** se não tiverem. Esses tipos apontam para **nomes**, e um nome sem ponto final seria interpretado como relativo à zona. O ponto o torna absoluto (FQDN), evitando que `mail.exemplo.com` vire `mail.exemplo.com.exemplo.com`.
- **`MX`** — a **prioridade** é prefixada ao valor (`10 mail.exemplo.com.`). O formato do MX exige o número de prioridade antes do host; a interface coleta isso em um campo separado e o backend junta.
- **`CAA`** — se o usuário digitou só a autoridade, o backend completa para o formato `0 issue "letsencrypt.org"`. Assim quem não conhece a sintaxe do CAA ainda consegue criar um registro válido.
- **`SRV`** — montado como `prioridade peso porta alvo.`, cada parte vinda de um campo próprio da interface, com ponto final garantido no alvo.
- **`HTTPS`** — montado como `prioridade alvo parâmetros` (por exemplo `1 . alpn=h2`), com um padrão sensato (`alpn=h2`) quando os parâmetros não são informados.

Repare que `HTTPS` e `SRV` **não usam um único campo de valor** na interface — eles têm campos estruturados (`svcPriority`, `targetName`, `svcParams`, `weight`, `port`, `target`, `priority`), porque seu conteúdo é composto por várias partes. Por isso o `records.Registro` no backend tem campos opcionais dedicados a esses tipos.

## Como o nome do registro é completado (FQDN)

Ao inserir, o backend ajusta o **nome** do registro para um FQDN absoluto ([nome-completo.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/records/nome-completo.go)):

```go
func NomeCompleto(zona, name string) string {
    if strings.HasSuffix(name, ".") {
        return name                                  // já é absoluto
    }

    zone := strings.TrimSuffix(zona, ".")
    if strings.HasSuffix(name, zone) {
        return name + "."                            // já tinha a zona, só falta o ponto
    }

    return fmt.Sprintf("%s.%s.", name, zone)         // "www" → "www.exemplo.com."
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

Este é o comportamento mais característico do SanchezDNS. Sempre que um registro **`A`** ou **`AAAA`** é criado, editado ou removido, o backend mantém o **PTR reverso** correspondente em sincronia, sem intervenção manual. Cada operação tem sua função: [insert-reverso.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/records/insert-reverso.go), [update-reverso.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/records/update-reverso.go) e [delete-reverso.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/records/delete-reverso.go).

As três compartilham os mesmos blocos de construção — `ParseIP`, `NomeReverso`, `ZonaReversa`, `TemForward`, `InsertPTR` e `DeletePTR` — o que faz IPv4 e IPv6 seguirem exatamente o mesmo caminho, em vez de terem implementações separadas.

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

O backend lista todas as zonas (`GET .../zones`) e procura, entre as que terminam em `in-addr.arpa`/`ip6.arpa`, aquela cujo nome é o **maior sufixo** do nome reverso calculado ([zona-reversa.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/records/zona-reversa.go)):

```go
for _, zona := range zonas {
    nome := strings.TrimSuffix(zona.Name, ".")

    if !strings.HasSuffix(nome, sufixo) || !strings.HasSuffix(reverso, nome) {
        continue
    }

    if len(nome) > len(strings.TrimSuffix(melhor, ".")) {   // vence o sufixo mais longo
        melhor = zona.Name
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

## Conferência de reversos: ausentes e órfãos

A automação acima só age no momento em que o registro é criado, editado ou removido. Ela não cobre dois casos que aparecem na prática:

- registros `A`/`AAAA` que **já existiam** antes de a automação entrar em cena, ou criados enquanto a zona reversa ainda não existia;
- PTRs que **sobraram** de registros diretos que foram apagados por fora do painel.

Para isso a tela da zona tem dois botões de conferência sob demanda, que aparecem conforme o tipo de zona.

### Reversos ausentes (zonas normais)

O botão **Verificar reversos** varre a zona atual atrás de `A`/`AAAA` que **deveriam** ter PTR e não têm ([fetch-reversos-ausentes.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/records/fetch-reversos-ausentes.go)):

1. Lista os registros de IP da zona (`RegistrosIP`).
2. Para cada IP, calcula o nome reverso e procura a zona reversa mais específica que o cobre.
3. **Descarta** os IPs sem zona reversa cadastrada — não é erro, apenas não é um reverso que você gerencia.
4. Dos que sobraram, mantém só os que **não têm PTR** na zona reversa (`semPTR`).

O resultado abre em um modal com seleção múltipla, e você escolhe quais criar. O `PUT /reverses` recebe a lista de nomes reversos escolhidos, **recalcula os ausentes no servidor** e cria só os que estiverem de fato na interseção — assim uma lista desatualizada na tela nunca cria um PTR indevido.

### Reversos órfãos (zonas reversas)

Dentro de uma zona `in-addr.arpa`/`ip6.arpa`, o botão **Verificar órfãos** faz o caminho inverso ([fetch-reversos-orfaos.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/records/fetch-reversos-orfaos.go)):

1. Monta o conjunto de **todos os IPs** que aparecem em `A`/`AAAA` de todas as zonas normais (`IPsForward`).
2. Percorre os PTRs da zona reversa, converte cada nome reverso de volta para IP (`IPDoReverso`).
3. Marca como órfão todo PTR cujo IP **não aparece** em nenhum registro direto.

Cada órfão pode ser removido individualmente, com confirmação. As duas operações de escrita (`PUT` e `DELETE /reverses`) exigem **escrita** na zona e geram log (`insert_reverses`, `delete_reverse`).

**Por que sob demanda e não automático?** Porque as duas varreduras leem **todas** as zonas do servidor para montar o panorama — é uma operação cara demais para rodar a cada requisição. E remover PTR órfão é destrutivo: a decisão fica com o operador, não com uma rotina automática.

## Permissões na tela de zona

- O botão de gestão de usuários e a edição de SOA só aparecem quando a zona está em nível **`ADMINISTRADOR`** para o usuário.
- Para usuários com nível **`LEITURA`**, a tabela de registros fica **somente leitura**.
- Para criar/editar/remover registros, o usuário precisa de **`ESCRITA`** naquela zona.

O modelo de permissões por zona é detalhado em [Usuários](/users), e a verificação acontece no backend a cada operação — a interface apenas reflete o que o servidor já decide.
