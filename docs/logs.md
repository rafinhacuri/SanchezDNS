# 🧾 Logs e Auditoria

A trilha de auditoria é uma das razões de existir do SanchezDNS (ver [Por que SanchezDNS](/reason)). Esta página explica **o que é registrado, como e por quê**.

## O que aparece na página

Cada entrada de log traz:

- **zona** — sobre qual zona (ou registro) a ação incidiu;
- **usuário** — o email de quem executou;
- **ação** — um código do que aconteceu (ex.: `insert_record`, `create_zone`, `update_soa`);
- **detalhes** — uma descrição legível;
- **data de criação** — quando ocorreu.

## Eventos registrados

O sistema gera logs para as operações que realmente importam para auditoria:

- **criação e remoção de zonas** (`create_zone`, `delete_zone`);
- **criação, edição e remoção de registros** (`insert_record`, `edit_record`, `delete_record`);
- **inclusão, edição e remoção de usuários por zona** (`insert_user`, `update_user`, `delete_user`);
- **atualização de SOA** (`update_soa`);
- **criação em lote de reversos ausentes e remoção de reversos órfãos** (`insert_reverses`, `delete_reverse`);
- **aprovação ou rejeição de solicitações** de cadastro.

## Como os logs são gravados — e por que de forma assíncrona

Os logs são escritos no **MongoDB**, mas de um jeito específico: em **segundo plano**, via _goroutine_. Repare no padrão que aparece em todo o backend:

```go
go logs.Insert(name, email, "insert_record", "Criado registro ...")
```

O `go` na frente faz a gravação do log rodar **em paralelo**, sem bloquear a resposta ao usuário. A operação principal (criar o registro, por exemplo) já retorna sucesso enquanto o log é persistido de forma independente.

Essa chamada fica sempre **no controlador**, nunca dentro do pacote de domínio. Os pacotes (`records`, `zonas`, `users`, ...) só executam a operação e devolvem o erro; quem decide o que virou log, com qual texto, é a camada que já conhece o usuário da requisição.

**Por que assíncrono?** Duas razões:

1. **Performance.** O usuário não deve esperar a escrita da auditoria para receber a confirmação da ação. A resposta fica mais rápida.
2. **Robustez.** A auditoria é importante, mas **não deve derrubar a operação principal**. Se, por um instante, a gravação do log falhar, o registro DNS que o usuário criou continua válido. A auditoria é um efeito colateral desejável, não um pré-requisito da ação.

Isso é coerente com o princípio de "automação de melhor esforço" que guia o projeto, o mesmo que rege o PTR reverso (ver [Zonas e Registros](/zones#reverse-automatico-o-ptr-que-se-cuida-sozinho)).

## Comportamento da página

- **busca textual** para filtrar os eventos;
- **paginação** para navegar em grande volume de registros;
- **ordenação** padrão do mais recente para o mais antigo.

## Acesso

Somente **`admin`** acessa a página de logs (rota `GET /logs` protegida por `AdminOnly`). Faz sentido: a auditoria é justamente a ferramenta de quem administra o ambiente.

## Objetivo

O log existe para **auditoria operacional**: permite rastrear **quem fez o quê, em qual zona e quando**, sem depender de memória ou de histórico manual. Em um sistema que controla o DNS de uma organização, essa rastreabilidade é o que torna possível investigar um incidente ou entender uma mudança depois que ela aconteceu.
