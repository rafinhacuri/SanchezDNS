// oxlint-disable typescript/explicit-function-return-type typescript/explicit-module-boundary-types harlanzw/vue-no-faux-composables unicorn/max-nested-calls
import { email, nonEmpty, object, optional, picklist, pipe, string } from 'valibot'
import type { InferInput } from 'valibot'

export function useUserSchema() {
  const insertUserSchema = object({
    id: optional(string()),
    email: pipe(
      string('Email é obrigatório'),
      nonEmpty('Email é obrigatório'),
      email('Email é inválido'),
    ),
    permissao: picklist(['leitura', 'escrita'], 'Permissão é obrigatória'),
    zona: pipe(string('Zona é obrigatória'), nonEmpty('Zona é obrigatória')),
  })

  return { insertUserSchema }
}

export type InsertUserSchema = InferInput<ReturnType<typeof useUserSchema>['insertUserSchema']>
