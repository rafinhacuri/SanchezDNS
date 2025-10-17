import { z } from 'zod'

export const ZoneSchema = z.object({
  domain: z.string().min(1, 'Domain are required'),
  kind: z.enum(['Native', 'Primary', 'Secondary']),
  soa_edit_api: z.enum(['DEFAULT', 'INCREASE', 'EPOCH', 'OFF']),
  masters: z.array(z.string()),
  also_notify: z.array(z.string()),
})

export type ZoneSchemaType = z.infer<typeof ZoneSchema>
