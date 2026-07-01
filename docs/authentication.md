# 🔐 Autenticação e Sessões

Esta página explica **como o SanchezDNS autentica os usuários** e, principalmente, **por que foram tomadas as decisões de projeto** por trás disso: por que **sessão em vez de JWT**, o que significa **stateless vs. stateful**, e **qual o papel do Redis** nessa história. Essas são exatamente as perguntas que um TCC precisa justificar, então esta seção é intencionalmente detalhada.

## O fluxo de login, passo a passo

Quando o usuário envia email e senha para `POST /login` ([login.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/controller/login.go)):

1. **Já existe sessão?** Se o cookie `sanchezdns_session_id` já está presente, o login é recusado — evita criar sessões duplicadas.
2. **O cadastro está aprovado?** O backend verifica se há uma **solicitação pendente** para aquele email. Se estiver "em análise", retorna `401` com essa mensagem (o fluxo de aprovação é descrito em [Usuários](/users)).
3. **A senha confere?** `ValidarSenha` compara a senha enviada com o **hash bcrypt** guardado no cadastro. Nunca há senha em texto puro no banco.
4. **Coleta de contexto.** A partir do `User-Agent` e do IP, o backend deriva sistema operacional, navegador e **geolocalização aproximada** — dados que ficam registrados na sessão para auditoria ("de onde e de qual dispositivo esse login veio").
5. **Cria a sessão.** Gera um identificador aleatório (o **SID**) e grava um documento na coleção `sessions` do MongoDB.
6. **Entrega o cookie.** O SID é devolvido em um cookie `HttpOnly`, e o navegador passa a enviá-lo em toda requisição.

### Por que bcrypt para as senhas

As senhas são guardadas como **hash bcrypt**, não em texto puro nem com um hash simples (MD5/SHA). O bcrypt é **propositalmente lento** e usa um _fator de custo_ ajustável, o que torna ataques de força bruta caros mesmo se o banco vazar. Ele também embute um _salt_ por senha, então duas contas com a mesma senha geram hashes diferentes — anulando ataques por tabela pré-computada (rainbow tables).

## O que é o SID e por que não é um JWT

O identificador de sessão é gerado assim ([create-session.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/auth/create-session.go)):

```go
func newSID() (string, error) {
    b := make([]byte, 32)
    _, err := rand.Read(b)              // 32 bytes aleatórios de fonte criptográfica
    return base64.RawURLEncoding.EncodeToString(b), nil
}
```

Isso é um **token opaco**: 256 bits de aleatoriedade, sem nenhum significado embutido. Ele é apenas uma **chave** que aponta para os dados da sessão guardados no servidor. Isso é o oposto de um **JWT**.

### JWT vs. Sessão — a decisão central

Um **JWT (JSON Web Token)** é um token **autocontido**: ele carrega dentro de si os dados do usuário (email, nível, expiração) assinados pelo servidor. O servidor valida o token verificando a assinatura — **sem consultar banco nenhum**. Isso é o que se chama de autenticação **stateless** (sem estado no servidor).

O SanchezDNS escolheu **não usar JWT** e adotar **sessão do lado do servidor** (o token opaco + dados no servidor). Os motivos:

| Aspecto                  | Sessão (escolhido)                                         | JWT                                                                          |
| ------------------------ | ---------------------------------------------------------- | ---------------------------------------------------------------------------- |
| **Revogação**            | Imediata — basta marcar a sessão como inativa no servidor  | Difícil — o token continua válido até expirar; precisa de listas de bloqueio |
| **Conteúdo exposto**     | Nada; o token não diz nada sobre o usuário                 | O payload é legível por qualquer um (só assinado, não secreto)               |
| **Mudança de permissão** | Reflete na hora; o nível é lido do servidor a cada request | Fica "congelada" no token até ele expirar                                    |
| **Tamanho do cookie**    | Pequeno (32 bytes)                                         | Maior, cresce com os dados embutidos                                         |
| **Gestão de segredo**    | Não há segredo de assinatura para rotacionar               | Vazar a chave de assinatura compromete todos os tokens                       |

**Por que a revogação é o argumento decisivo aqui?** Porque o SanchezDNS controla o **DNS de uma organização** — algo sensível. Se uma conta é comprometida, desligada ou tem permissões alteradas, o administrador precisa conseguir **encerrar aquela sessão na hora**. Com JWT stateless isso é notoriamente difícil (o token permanece válido até expirar). Com sessão do lado do servidor, encerrar é trivial: marca-se `active: false` e o próximo request já é barrado. Para um painel administrativo, **controle sobre a sessão vale mais do que a economia de uma consulta**.

Além disso, como o nível de acesso (`admin`/`member`) é **lido a cada requisição** e não embutido num token, promover ou rebaixar um usuário tem efeito imediato — o que combina com o modelo de permissões por zona do projeto.

## Stateless vs. Stateful — o que muda

- **Autenticação stateless:** o servidor **não guarda nada** sobre a sessão. Toda a informação necessária está no token que o cliente apresenta (o caso do JWT). Vantagem: escala horizontalmente sem esforço, qualquer instância valida sozinha. Desvantagem: não dá para revogar nem alterar sem esperar o token expirar.

- **Autenticação stateful:** o servidor **guarda o estado** da sessão (quem é, se está ativa, última atividade). O cliente só carrega uma chave (o SID). Vantagem: controle total — revogar, listar sessões ativas, expirar por inatividade. Desvantagem: cada validação precisa **consultar onde o estado está guardado**.

O SanchezDNS é **stateful** — e é justamente essa "desvantagem" (precisar consultar o estado a cada request) que o **Redis** vem resolver.

## O papel do Redis

Sendo stateful, o backend precisa responder duas perguntas em **toda requisição protegida**:

1. "Essa sessão (`sid`) é válida e está ativa?"
2. "Qual o nível (`admin`/`member`) desse usuário?"

A fonte da verdade dessas respostas é o **MongoDB**. Mas consultar o Mongo a cada requisição — inclusive para carregar cada imagem, cada refresh de tela — seria desperdício. É aí que entra o **Redis** como **cache de leitura rápida**.

Veja como a validação funciona ([validate-session.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/auth/validate-session.go)):

```go
sessionRedis, err := redis.GetSession(ctx, sid)
if err == nil {
    if !sessionRedis.Active { return ..., errors.New("session inactive") }
    email = sessionRedis.Email          // achou no Redis: resposta instantânea
} else {
    // não estava no Redis: cai para o MongoDB (fonte da verdade)
    err := sessionsColl.FindOne(ctx, bson.M{"sid": sid, "active": true}).Decode(&sessionMongo)
    email = sessionMongo.Email
}
level, err := cadastro.GetLevel(ctx, email)   // nível: também cacheado no Redis
```

E o `AdminOnly` ([admin.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/middleware/admin.go)) usa o `level` já resolvido para proteger as rotas administrativas.

### Por que Redis e não só o MongoDB

- **Velocidade.** O Redis mantém os dados **em memória**, com latência de microssegundos. A validação de sessão — a operação mais frequente do sistema — fica praticamente gratuita.
- **Chaves com significado e TTL.** As chaves são organizadas por prefixo (`sanchezdns:session:<sid>`, `sanchezdns:level:<email>`) e podem expirar sozinhas, mantendo o cache "fresco" sem limpeza manual.
- **Suporte a JSON.** A sessão é guardada com `JSON.GET`/`JSON.SET` (módulo RedisJSON), então o objeto de sessão é lido e escrito estruturado, não como uma string qualquer.
- **Degradação graciosa.** Se o Redis estiver indisponível ou não tiver a chave, o código **cai automaticamente para o MongoDB**. O sistema fica mais lento, mas **não quebra**. O Redis é otimização, não ponto único de falha.

### O ponto sutil: continua stateful

Um detalhe conceitual importante para o TCC: **usar Redis não torna a autenticação stateless.** O estado continua existindo do lado do servidor — só está distribuído entre um cache rápido (Redis) e a fonte da verdade (MongoDB). O cliente ainda carrega apenas uma chave opaca. Ou seja, o SanchezDNS combina o **controle do modelo stateful** com a **velocidade que normalmente se associa ao stateless**, sem abrir mão da capacidade de revogar sessões. É o melhor dos dois mundos para um painel administrativo.

## O cookie de sessão

O SID é entregue em um cookie configurado assim:

```go
c.SetCookie("sanchezdns_session_id", sessionID, 34560000, "/", siteURL, secure, true)
```

- **`HttpOnly` (o último `true`)** — o cookie **não é acessível pelo JavaScript** do navegador. Isso protege contra roubo de sessão via XSS: mesmo que um script malicioso rode na página, ele não consegue ler o token. É também o motivo de o frontend ser SSR (Nuxt), como explicado na [Arquitetura](/architecture): quem repassa o cookie para a API é o servidor, não o código do cliente.
- **`secure`** — em produção (`NUXT_PUBLIC_PRODUCTION=true`) o cookie só trafega por HTTPS, evitando que seja capturado em uma conexão sem criptografia.
- **`34560000` segundos (~400 dias)** — validade longa do cookie, pensada para conveniência. Como o controle real é do lado do servidor (a sessão pode ser desativada a qualquer momento), uma expiração longa do cookie não enfraquece a segurança: quem manda é o estado da sessão, não o prazo do cookie.

## Validação em cada requisição protegida

Toda rota autenticada passa antes pelo middleware `ValidateSession` ([validate-session.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/middleware/validate-session.go)):

1. Lê o cookie `sanchezdns_session_id`. Sem cookie → `400`.
2. Valida a sessão (Redis → fallback Mongo) e obtém `email` + `level`.
3. Se a sessão for inválida, **limpa o cookie** (`maxAge -1`) e retorna erro — o navegador esquece o token morto.
4. Se for válida, injeta `email`, `level` e `sid` no contexto, e o controlador segue sabendo quem é o usuário.

Sobre esse contexto injetado é que se aplicam as regras de autorização: `AdminOnly` para rotas administrativas e a checagem das listas `leitura`/`escrita` por zona para operações de registro (ver [Usuários](/users)).

## Resumo das decisões

| Decisão                      | Por quê                                                                              |
| ---------------------------- | ------------------------------------------------------------------------------------ |
| Sessão em vez de JWT         | Poder **revogar na hora** e refletir mudanças de permissão imediatamente             |
| Token opaco de 256 bits      | Não vaza nenhuma informação; é só uma chave                                          |
| Modelo **stateful**          | Controle total sobre sessões (revogar, expirar, auditar)                             |
| **Redis** como cache         | Dar ao modelo stateful a **velocidade** do stateless, sem virar ponto único de falha |
| **bcrypt** nas senhas        | Hash lento com salt, resistente a força bruta e rainbow tables                       |
| Cookie `HttpOnly` + `Secure` | Proteção contra XSS e contra captura em rede sem HTTPS                               |
