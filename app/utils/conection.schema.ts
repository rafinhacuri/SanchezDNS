import { z } from 'zod'

export const ConnectionSchema = z.object({
  name: z.string().min(3, 'Name must be at least 3 characters long'),
  host: z.string().min(3, 'Host must be at least 3 characters long'),
  apiKey: z.string().min(10, 'API Key must be at least 10 characters long'),
  serverId: z.string().default('localhost'),
  users: z.array(z.string()),
})

export type ConnectionType = z.infer<typeof ConnectionSchema>
