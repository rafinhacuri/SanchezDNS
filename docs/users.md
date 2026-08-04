# 👥 Usuários e Permissões

Esta página descreve **como o SanchezDNS gerencia contas, níveis de acesso e permissões por zona**, e por que o modelo foi desenhado assim. O funcionamento da sessão em si está em [Autenticação](/authentication).

## Dois níveis de acesso

O sistema trabalha com dois níveis:

- **`admin`** — acesso administrativo global.
- **`member`** — usuário comum, com acesso apenas às zonas para as quais foi liberado.

**O primeiro cadastro aprovado vira `admin`.** Os cadastros seguintes entram como **solicitação** e ficam pendentes até um administrador aprovar. Isso resolve o problema do "primeiro acesso" sem senha padrão nem passo manual de bootstrap: quem instala e cria a primeira conta assume o controle, e a partir daí o acesso é sempre por aprovação.

## O fluxo de cadastro e aprovação

1. Uma pessoa envia um pedido de cadastro (`POST /cadastro`).
2. O pedido fica em **solicitações**, com status pendente.
3. Ao tentar logar antes da aprovação, ela recebe "seu cadastro ainda está em análise" (o backend barra no login — ver [Autenticação](/authentication#o-fluxo-de-login-passo-a-passo)).
4. Um administrador **aprova ou rejeita** na página de Solicitações. Aprovar cria o cadastro; a ação é registrada em log.

**Por que aprovação manual?** Porque o SanchezDNS controla o DNS de uma organização — não faz sentido permitir auto-registro aberto. A barreira de aprovação garante que só entra quem um administrador autorizou.

## O que o `admin` pode fazer

O administrador tem acesso às áreas e operações administrativas, todas protegidas no backend pelo middleware `AdminOnly`:

- páginas de **Logs**, **Solicitações** e **Cadastros**;
- **criar e remover zonas**;
- **atualizar o SOA** de uma zona;
- **gerenciar permissões por zona** (quem lê, quem escreve);
- gerenciar cadastros (editar, alterar nível, excluir).

O `admin` **ignora** as restrições de leitura/escrita por zona — ele enxerga e opera todas.

## O que o `member` pode fazer

O usuário comum só enxerga as zonas para as quais recebeu permissão, e o que pode fazer nelas depende do tipo de permissão:

- **`leitura`** — pode **listar** os registros da zona.
- **`escrita`** — pode **criar, editar e remover** registros da zona.
- **sem permissão** — a zona fica inacessível.

## Como a permissão por zona é armazenada e verificada

Para cada zona, o backend guarda na coleção `users` do MongoDB **duas listas de emails**: `leitura` e `escrita`. A checagem fica concentrada em três funções do pacote `users`, e todo controlador que precisa decidir acesso chama uma delas:

| Função                                  | Responde                                                     |
| --------------------------------------- | ------------------------------------------------------------ |
| `PodeLer(ctx, zona, email, level)`      | pode ver os registros desta zona?                            |
| `PodeEscrever(ctx, zona, email, level)` | pode criar, editar ou remover registros desta zona?          |
| `FetchNivel(ctx, zona, email, level)`   | qual o nível (`ADMINISTRADOR`, `ESCRITA`, `LEITURA` ou nada) |

Por dentro, todas seguem o mesmo caminho ([pode-escrever.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/users/pode-escrever.go)):

```go
func PodeEscrever(ctx context.Context, zona, email, level string) (bool, error) {
    if level == "admin" {
        return true, nil          // admin ignora a permissão por zona
    }

    permissao, err := FetchPermissao(ctx, zona, email)
    if err != nil {
        return false, err
    }

    return permissao == "escrita", nil
}
```

E o controlador só orquestra ([insert/record.go](https://github.com/rafinhacuri/SanchezDNS/blob/main/api/controller/insert/record.go)):

```go
pode, err := users.PodeEscrever(ctx, body.Zone, email, c.GetString("level"))
if !pode {
    c.AbortWithStatusJSON(403, gin.H{"message": "Sem permissão para alterar registros dessa zona"})
    return
}
```

Quatro pontos de projeto valem destaque:

- **A verificação é sempre no backend.** A interface esconde botões conforme o nível, mas quem **decide** é o servidor a cada operação. Ocultar um botão nunca é a proteção real.
- **Uma única implementação.** Como a checagem vive no pacote `users`, todos os controladores decidem acesso do mesmo jeito — não há risco de uma rota esquecer um caso que outra trata.
- **Comparação normalizada.** Os emails são comparados em minúsculas e sem espaços, evitando que `Fulano@x.com ` e `fulano@x.com` sejam tratados como pessoas diferentes.
- **Modelo por listas.** Guardar as permissões como listas de emails dentro do documento da zona encaixa no modelo de documentos do MongoDB (ver [Arquitetura](/architecture#mongodb-a-verdade-do-aplicativo)) e torna trivial adicionar/remover acesso.

## Tela de usuários da zona

Em uma zona administrativa é possível:

- **adicionar** um email com permissão de leitura ou escrita;
- **alterar** a permissão de um usuário já associado;
- **remover** o acesso de um usuário àquela zona.

Cada uma dessas ações gera log de auditoria.

## Tela de cadastros

A página de Cadastros lista as contas já aprovadas e permite:

- editar **nome, foto e senha**;
- alterar o **nível** do cadastro (`admin`/`member`) — com efeito imediato, pois o nível é lido a cada requisição;
- **excluir** o cadastro e limpar os acessos associados.

## Fotos de perfil

O login usa o email e a senha do cadastro aprovado (senha em bcrypt). A **foto de perfil** é guardada em um storage S3 compatível, não no banco nem no disco do container — coerente com a ideia de container sem estado (ver [Arquitetura](/architecture#s3-compativel-rustfs-minio-arquivos)).
