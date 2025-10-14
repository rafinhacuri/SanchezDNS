import { z } from 'zod'

export const GroupSchema = z.object({
  name: z.string().min(3, 'At least 3 characters').max(30, 'Maximum 30 characters'),
  description: z.string().max(100, 'Maximum 100 characters'),
  conection: z.string().min(3, 'At least 3 characters'),
  members: z.array(z.string()),
})

export type Group = z.infer<typeof GroupSchema>
