import { email, nonEmpty, object, pipe, string } from 'valibot'
import type { InferInput } from 'valibot'

export const AuthSchema = object({
  email: pipe(
    string('Email é uma string'),
    email('Email é inválido'),
    nonEmpty('Email é obrigatório'),
  ),
  senha: pipe(string('Senha é uma string'), nonEmpty('Senha é obrigatória')),
})

export type Auth = InferInput<typeof AuthSchema>
