import {
  check,
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

export const EditSOASchema = object({
  startOfAuthority: pipe(
    string('Start of Authority é uma string'),
    nonEmpty('Start of Authority é obrigatório'),
  ),
  email: pipe(string('Email é uma string'), nonEmpty('Email é obrigatório')),
  refresh: pipe(
    number('Refresh é um número'),
    integer('Refresh must be a positive integer'),
    minValue(1, 'Refresh must be a positive integer'),
  ),
  retry: pipe(
    number('Retry é um número'),
    integer('Retry must be a positive integer'),
    minValue(1, 'Retry must be a positive integer'),
  ),
  expire: pipe(
    number('Expire é um número'),
    integer('Expire must be a positive integer'),
    minValue(1, 'Expire must be a positive integer'),
  ),
  negativeCacheTtl: pipe(
    number('Negative Cache TTL é um número'),
    integer('Negative Cache TTL must be a positive integer'),
    minValue(1, 'Negative Cache TTL must be a positive integer'),
  ),
})

export type EditSOASchemaType = InferInput<typeof EditSOASchema>

export const RecordSchema = pipe(
  object({
    zone: pipe(string('Zone ID é uma string'), nonEmpty('Zone ID é obrigatório')),
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
    name: optional(string('Name é uma string')),
    vl: optional(string('Value é uma string')),
    ttl: pipe(number('TTL é um número'), minValue(60, 'TTL must be at least 60 seconds')),
    comment: optional(string('Comment é uma string')),
    svcPriority: optional(number('Service Priority é um número')),
    targetName: optional(string('Target Name é uma string')),
    svcParams: optional(string('Service Params é uma string')),
    weight: optional(number('Weight é um número')),
    port: optional(number('Port é um número')),
    target: optional(string('Target é uma string')),
    priority: optional(number('Priority é um número')),
  }),
  check(
    (data) => ['HTTPS', 'SRV'].includes(data.type) || Boolean(data.vl?.trim()),
    'Value is required for this record type',
  ),
  check(
    (data) => data.type !== 'HTTPS' || data.svcPriority !== undefined,
    'Service Priority is required for HTTPS records',
  ),
  check(
    (data) => data.type !== 'HTTPS' || Boolean(data.targetName?.trim()),
    'Target Name is required for HTTPS records',
  ),
  check(
    (data) => data.type !== 'SRV' || data.weight !== undefined,
    'Weight is required for SRV records',
  ),
  check(
    (data) => data.type !== 'SRV' || data.port !== undefined,
    'Port is required for SRV records',
  ),
  check(
    (data) => data.type !== 'SRV' || Boolean(data.target?.trim()),
    'Target is required for SRV records',
  ),
  check(
    (data) => data.type !== 'SRV' || data.priority !== undefined,
    'Priority is required for SRV records',
  ),
  check(
    (data) => data.type !== 'MX' || data.priority !== undefined,
    'Priority is required for MX records',
  ),
)

export type RecordForm = InferInput<typeof RecordSchema>

export const EditRecordSchema = object({
  oldValue: RecordSchema,
  newValue: RecordSchema,
})

export type EditRecordForm = InferInput<typeof EditRecordSchema>
