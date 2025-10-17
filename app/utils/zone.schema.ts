import { z } from 'zod'

export const ZoneSchema = z.object({
  domain: z.string().min(1, 'Domain are required'),
  soa: z.object({
    startOfAuthority: z.string().min(1, 'Start of Authority is required'),
    email: z.string().min(1, 'Email is required'),
    refresh: z.number().int().positive('Refresh must be a positive integer'),
    retry: z.number().int().positive('Retry must be a positive integer'),
    expire: z.number().int().positive('Expire must be a positive integer'),
    negativeCacheTtl: z.number().int().positive('Negative Cache TTL must be a positive integer'),
  }),
})

export type ZoneSchemaType = z.infer<typeof ZoneSchema>
