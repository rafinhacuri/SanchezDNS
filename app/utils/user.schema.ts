import { z } from 'zod'

export const InsertUserSchema = z.object({
  id: z.string().optional(),
  idcbpf: z.string().min(1, 'ID CBPF é obrigatório'),
  permissao: z.string().min(1, 'Permissão é obrigatória'),
  zona: z.string().min(1, 'Zona é obrigatória'),
})

export type InsertUserType = z.infer<typeof InsertUserSchema>
