import { nonEmpty, object, pipe, string } from 'valibot'
import type { InferInput } from 'valibot'

export const IdSchema = object({
  id: pipe(string('ID é uma string'), nonEmpty('ID é obrigatório')),
})

export type Id = InferInput<typeof IdSchema>
