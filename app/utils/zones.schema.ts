import { integer, minValue, nonEmpty, number, object, pipe, string } from 'valibot'
import type { InferInput } from 'valibot'

const SOASchema = object({
  startOfAuthority: pipe(
    string('Start of Authority é uma string'),
    nonEmpty('Start of Authority é obrigatório'),
  ),
  email: pipe(string('Email é uma string'), nonEmpty('Email é obrigatório')),
  refresh: pipe(
    number('Refresh é um número'),
    integer('Refresh precisa ser um inteiro positivo'),
    minValue(1, 'Refresh precisa ser um inteiro positivo'),
  ),
  retry: pipe(
    number('Retry é um número'),
    integer('Retry precisa ser um inteiro positivo'),
    minValue(1, 'Retry precisa ser um inteiro positivo'),
  ),
  expire: pipe(
    number('Expire é um número'),
    integer('Expire precisa ser um inteiro positivo'),
    minValue(1, 'Expire precisa ser um inteiro positivo'),
  ),
  negativeCacheTtl: pipe(
    number('Negative Cache TTL é um número'),
    integer('Negative Cache TTL precisa ser um inteiro positivo'),
    minValue(1, 'Negative Cache TTL precisa ser um inteiro positivo'),
  ),
})

export const ZoneSchema = object({
  domain: pipe(string('Domain é uma string'), nonEmpty('Domain é obrigatório')),
  type: pipe(string('Type é uma string'), nonEmpty('Type é obrigatório')),
  soa: SOASchema,
})

export type ZoneSchemaType = InferInput<typeof ZoneSchema>
