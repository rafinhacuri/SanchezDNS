// oxlint-disable typescript/explicit-function-return-type typescript/explicit-module-boundary-types harlanzw/vue-no-faux-composables unicorn/max-nested-calls
import { integer, minValue, nonEmpty, number, object, pipe, string } from 'valibot'
import type { InferInput } from 'valibot'

export function useZonesSchema() {
  const zoneSchema = object({
    domain: pipe(string('Dominio é obrigatório'), nonEmpty('Dominio é obrigatório')),
    type: pipe(string('Tipo é obrigatório'), nonEmpty('Tipo é obrigatório')),
    soa: object({
      startOfAuthority: pipe(
        string('Start of Authority é obrigatório'),
        nonEmpty('Start of Authority é obrigatório'),
      ),
      email: pipe(string('Email é obrigatório'), nonEmpty('Email é obrigatório')),
      refresh: pipe(
        number('Refresh precisa ser um inteiro positivo'),
        integer('Refresh precisa ser um inteiro positivo'),
        minValue(1, 'Refresh precisa ser um inteiro positivo'),
      ),
      retry: pipe(
        number('Retry precisa ser um inteiro positivo'),
        integer('Retry precisa ser um inteiro positivo'),
        minValue(1, 'Retry precisa ser um inteiro positivo'),
      ),
      expire: pipe(
        number('Expire precisa ser um inteiro positivo'),
        integer('Expire precisa ser um inteiro positivo'),
        minValue(1, 'Expire precisa ser um inteiro positivo'),
      ),
      negativeCacheTtl: pipe(
        number('Negative Cache TTL precisa ser um inteiro positivo'),
        integer('Negative Cache TTL precisa ser um inteiro positivo'),
        minValue(1, 'Negative Cache TTL precisa ser um inteiro positivo'),
      ),
    }),
  })

  return { zoneSchema }
}

export type ZoneSchema = InferInput<ReturnType<typeof useZonesSchema>['zoneSchema']>
