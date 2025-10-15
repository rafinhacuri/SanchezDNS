import { z } from 'zod'

export const AddUserSchema = z.object({
  email: z.email('Invalid email'),
  connection: z.string().min(1, 'Connection is required'),
})

export type AddUser = z.infer<typeof AddUserSchema>
