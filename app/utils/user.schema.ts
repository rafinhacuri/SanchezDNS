import { nonEmpty, object, optional, pipe, string } from 'valibot'
import type { InferInput } from 'valibot'

export const InsertUserSchema = object({
  id: optional(string('ID é uma string')),
  email: pipe(string('Email é uma string'), nonEmpty('Email é obrigatório')),
  permissao: pipe(string('Permissão é uma string'), nonEmpty('Permissão é obrigatória')),
  zona: pipe(string('Zona é uma string'), nonEmpty('Zona é obrigatória')),
})

export type InsertUserType = InferInput<typeof InsertUserSchema>
