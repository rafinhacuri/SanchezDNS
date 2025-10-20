import { z } from 'zod'

export const RecordSchema = z.object({
  zone: z.string().min(1, 'Zone ID is required'),
  type: z.enum([
    'A',
    'AAAA',
    'ALIAS',
    'CAA',
    'CNAME',
    'HTTPS',
    'MX',
    'NS',
    'TXT',
    'SRV',
  ]),
  name: z.string().min(1, 'Name is required'),
  vl: z.string().optional(),
  ttl: z.number().min(60, 'TTL must be at least 60 seconds'),
  comment: z.string().optional(),
  svcPriority: z.number().optional(),
  targetName: z.string().optional(),
  svcParams: z.string().optional(),
  weight: z.number().optional(),
  port: z.number().optional(),
  target: z.string().optional(),
  priority: z.number().optional(),
})
  .refine(data => data.type === 'HTTPS' || data.type === 'SRV' || (data.vl && data.vl.trim() !== ''), {
    message: 'Value is required for this record type',
    path: ['vl'],
  })
  .refine(data => data.type !== 'HTTPS' || (data.svcPriority !== undefined && data.svcPriority !== null), {
    message: 'Service Priority is required for HTTPS records',
    path: ['svcPriority'],
  })
  .refine(data => data.type !== 'HTTPS' || (data.targetName && data.targetName.trim() !== ''), {
    message: 'Target Name is required for HTTPS records',
    path: ['targetName'],
  })
  .refine(data => data.type !== 'SRV' || (data.weight !== undefined && data.weight !== null), {
    message: 'Weight is required for SRV records',
    path: ['weight'],
  })
  .refine(data => data.type !== 'SRV' || (data.port !== undefined && data.port !== null), {
    message: 'Port is required for SRV records',
    path: ['port'],
  })
  .refine(data => data.type !== 'SRV' || (data.target && data.target.trim() !== ''), {
    message: 'Target is required for SRV records',
    path: ['target'],
  })
  .refine(data => data.type !== 'SRV' || (data.priority !== undefined && data.priority !== null), {
    message: 'Priority is required for SRV records',
    path: ['priority'],
  })
  .refine(data => data.type !== 'MX' || (data.priority !== undefined && data.priority !== null), {
    message: 'Priority is required for MX records',
    path: ['priority'],
  })

export type RecordForm = z.infer<typeof RecordSchema>
