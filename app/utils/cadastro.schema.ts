import { email, nonEmpty, object, pipe, string } from 'valibot'
import type { InferInput } from 'valibot'

export const CadastroSchema = object({
  email: pipe(
    string('Email é uma string'),
    email('Email é inválido'),
    nonEmpty('Email é obrigatório'),
  ),
  senha: pipe(string('Senha é uma string'), nonEmpty('Senha é obrigatória')),
  foto: pipe(string('Foto é uma string'), nonEmpty('Foto é obrigatória')),
  nome: pipe(string('Nome é uma string'), nonEmpty('Nome é obrigatório')),
})

export type Cadastro = InferInput<typeof CadastroSchema>
