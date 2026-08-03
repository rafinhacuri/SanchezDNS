// oxlint-disable typescript/explicit-function-return-type typescript/explicit-module-boundary-types harlanzw/vue-no-faux-composables unicorn/max-nested-calls
import { email, minLength, nonEmpty, object, picklist, pipe, string } from 'valibot'
import type { InferInput } from 'valibot'

export function useCadastroSchema() {
  const authSchema = object({
    email: pipe(
      string('Email é obrigatório'),
      nonEmpty('Email é obrigatório'),
      email('Email é inválido'),
    ),
    senha: pipe(string('Senha é obrigatória'), nonEmpty('Senha é obrigatória')),
  })

  const cadastroSchema = object({
    nome: pipe(string('Nome é obrigatório'), nonEmpty('Nome é obrigatório')),
    email: pipe(
      string('Email é obrigatório'),
      nonEmpty('Email é obrigatório'),
      email('Email é inválido'),
    ),
    senha: pipe(
      string('Senha é obrigatória'),
      nonEmpty('Senha é obrigatória'),
      minLength(8, 'A senha precisa ter no mínimo 8 caracteres'),
    ),
    foto: pipe(string('Foto é obrigatória'), nonEmpty('Foto é obrigatória')),
  })

  const cadastroEditSchema = object({
    id: pipe(string('ID é obrigatório'), nonEmpty('ID é obrigatório')),
    nome: pipe(string('Nome é obrigatório'), nonEmpty('Nome é obrigatório')),
    foto: pipe(string('Foto é obrigatória'), nonEmpty('Foto é obrigatória')),
    senha: string('Senha é uma string'),
  })

  const statusSchema = object({
    id: pipe(string('ID é obrigatório'), nonEmpty('ID é obrigatório')),
    status: picklist(['aprovada', 'rejeitada'], 'Status é inválido'),
  })

  const idSchema = object({
    id: pipe(string('ID é obrigatório'), nonEmpty('ID é obrigatório')),
  })

  return { authSchema, cadastroSchema, cadastroEditSchema, statusSchema, idSchema }
}

export type AuthSchema = InferInput<ReturnType<typeof useCadastroSchema>['authSchema']>

export type CadastroSchema = InferInput<ReturnType<typeof useCadastroSchema>['cadastroSchema']>

export type CadastroEditSchema = InferInput<
  ReturnType<typeof useCadastroSchema>['cadastroEditSchema']
>

export type StatusSchema = InferInput<ReturnType<typeof useCadastroSchema>['statusSchema']>

export type IdSchema = InferInput<ReturnType<typeof useCadastroSchema>['idSchema']>
