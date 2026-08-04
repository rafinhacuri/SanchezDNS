// oxlint-disable typescript/explicit-function-return-type typescript/explicit-module-boundary-types harlanzw/vue-no-faux-composables unicorn/max-nested-calls
import {
  check,
  forward,
  integer,
  minValue,
  nonEmpty,
  number,
  object,
  optional,
  picklist,
  pipe,
  string,
} from 'valibot'
import type { InferInput } from 'valibot'

export function useRecordSchema() {
  const editSOASchema = object({
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
  })

  const recordSchema = pipe(
    object({
      zone: pipe(string('Zone ID é obrigatório'), nonEmpty('Zone ID é obrigatório')),
      type: picklist([
        'A',
        'AAAA',
        'ALIAS',
        'CAA',
        'CNAME',
        'HTTPS',
        'MX',
        'NS',
        'PTR',
        'TXT',
        'SRV',
        'TLSA',
      ]),
      name: optional(string()),
      vl: optional(string()),
      ttl: pipe(
        number('TTL precisa ser no mínimo 60 segundos'),
        minValue(60, 'TTL precisa ser no mínimo 60 segundos'),
      ),
      comment: optional(string()),
      svcPriority: optional(number()),
      targetName: optional(string()),
      svcParams: optional(string()),
      weight: optional(number()),
      port: optional(number()),
      target: optional(string()),
      priority: optional(number()),
    }),
    forward(
      check(
        (data) => ['HTTPS', 'SRV'].includes(data.type) || Boolean(data.vl?.trim()),
        'Value é obrigatório para esse tipo de registro',
      ),
      ['vl'],
    ),
    forward(
      check(
        (data) => data.type !== 'HTTPS' || (data.svcPriority ?? null) !== null,
        'Service Priority é obrigatório para registros HTTPS',
      ),
      ['svcPriority'],
    ),
    forward(
      check(
        (data) => data.type !== 'HTTPS' || Boolean(data.targetName?.trim()),
        'Target Name é obrigatório para registros HTTPS',
      ),
      ['targetName'],
    ),
    forward(
      check(
        (data) => data.type !== 'SRV' || (data.weight ?? null) !== null,
        'Weight é obrigatório para registros SRV',
      ),
      ['weight'],
    ),
    forward(
      check(
        (data) => data.type !== 'SRV' || (data.port ?? null) !== null,
        'Port é obrigatório para registros SRV',
      ),
      ['port'],
    ),
    forward(
      check(
        (data) => data.type !== 'SRV' || Boolean(data.target?.trim()),
        'Target é obrigatório para registros SRV',
      ),
      ['target'],
    ),
    forward(
      check(
        (data) => data.type !== 'SRV' || (data.priority ?? null) !== null,
        'Priority é obrigatório para registros SRV',
      ),
      ['priority'],
    ),
    forward(
      check(
        (data) => data.type !== 'MX' || (data.priority ?? null) !== null,
        'Priority é obrigatório para registros MX',
      ),
      ['priority'],
    ),
  )

  const editRecordSchema = object({
    oldValue: recordSchema,
    newValue: recordSchema,
  })

  return { editSOASchema, recordSchema, editRecordSchema }
}

export type EditSOASchema = InferInput<ReturnType<typeof useRecordSchema>['editSOASchema']>

export type RecordSchema = InferInput<ReturnType<typeof useRecordSchema>['recordSchema']>

export type EditRecordSchema = InferInput<ReturnType<typeof useRecordSchema>['editRecordSchema']>
