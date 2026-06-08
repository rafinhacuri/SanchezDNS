import { email, nonEmpty, object, picklist, pipe, string } from 'valibot'
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

export const StatusSchema = object({
  id: pipe(string('ID é uma string'), nonEmpty('ID é obrigatório')),
  status: pipe(
    string('Status é uma string'),
    nonEmpty('Status é obrigatório'),
    picklist(['aprovada', 'rejeitada'], 'Status é inválido'),
  ),
})

export type Status = InferInput<typeof StatusSchema>
