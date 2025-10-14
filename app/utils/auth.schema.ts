import { z } from 'zod'

export const AuthSchema = z.object({
  email: z.email('Invalid email'),
  password: z
    .string()
    .min(8, 'At least 8 characters - Requirement not met')
    .regex(/\d/, 'At least 1 number - Requirement not met')
    .regex(/[a-z]/, 'At least 1 lowercase letter - Requirement not met')
    .regex(/[A-Z]/, 'At least 1 uppercase letter - Requirement not met'),
})

export type Auth = z.infer<typeof AuthSchema>
