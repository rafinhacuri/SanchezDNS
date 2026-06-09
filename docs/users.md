# 👥 Usuários

## Modelo atual

O sistema trabalha com dois níveis principais:

- `admin`
- `member`

O primeiro cadastro aprovado vira `admin`. Os próximos cadastros ficam como solicitação até aprovação manual.

## Nível `admin`

O administrador tem acesso às áreas administrativas:

- `Logs`
- `Solicitações`
- `Cadastros`
- criação e remoção de zonas
- atualização de SOA
- gestão de permissões por zona

## Nível `member`

Usuários comuns só enxergam as zonas para as quais foram liberados.

As permissões são controladas por zona e podem ser:

- `leitura`
- `escrita`

## Como a permissão funciona

O backend guarda, por zona, duas listas de emails:

- `leitura`
- `escrita`

Quando o usuário abre uma zona:

- `leitura` permite listar registros;
- `escrita` permite criar, editar e excluir registros;
- sem permissão, a zona fica inacessível.

## Tela de usuários da zona

Na tela da zona administrativa é possível:

- adicionar um email com permissão de leitura ou escrita;
- alterar a permissão de um usuário já associado;
- remover o acesso de um usuário daquela zona.

## Tela de cadastros

A página de cadastros lista as contas já aprovadas:

- permite editar nome, foto e senha;
- permite alterar o nível do cadastro;
- permite excluir o cadastro e limpar acessos associados.

## Observações

- O login usa o email e a senha do cadastro aprovado.
- O sistema também salva a foto de perfil do usuário em storage S3 compatível.

---
