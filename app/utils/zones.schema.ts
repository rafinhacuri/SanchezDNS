import { z } from 'zod'

export const ZoneSchema = z.object({
  domain: z.string().min(1, 'Dominio é obrigatório'),
  type: z.string().min(1, 'Tipo é obrigatório'),
  soa: z.object({
    startOfAuthority: z.string().min(1, 'Start of Authority é obrigatório'),
    email: z.string().min(1, 'Email é obrigatório'),
    refresh: z.number().int().positive('Refresh precisa ser um inteiro positivo'),
    retry: z.number().int().positive('Retry precisa ser um inteiro positivo'),
    expire: z.number().int().positive('Expire precisa ser um inteiro positivo'),
    negativeCacheTtl: z
      .number()
      .int()
      .positive('Negative Cache TTL precisa ser um inteiro positivo'),
  }),
})

export type ZoneSchemaType = z.infer<typeof ZoneSchema>
