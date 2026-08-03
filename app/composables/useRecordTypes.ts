interface RecordTypeMeta {
  chip: string
  hint: string
}

const META: Record<string, RecordTypeMeta> = {
  A: {
    chip: 'bg-blue-500/10 text-blue-600 ring-blue-500/25 dark:text-blue-400',
    hint: 'Aponta o nome para um endereço IPv4',
  },
  AAAA: {
    chip: 'bg-violet-500/10 text-violet-600 ring-violet-500/25 dark:text-violet-400',
    hint: 'Aponta o nome para um endereço IPv6',
  },
  ALIAS: {
    chip: 'bg-teal-500/10 text-teal-600 ring-teal-500/25 dark:text-teal-400',
    hint: 'Apelido resolvido no servidor, válido na raiz da zona',
  },
  CAA: {
    chip: 'bg-orange-500/10 text-orange-600 ring-orange-500/25 dark:text-orange-400',
    hint: 'Define quais autoridades podem emitir certificados',
  },
  CNAME: {
    chip: 'bg-cyan-500/10 text-cyan-600 ring-cyan-500/25 dark:text-cyan-400',
    hint: 'Apelido que aponta para outro nome',
  },
  HTTPS: {
    chip: 'bg-emerald-500/10 text-emerald-600 ring-emerald-500/25 dark:text-emerald-400',
    hint: 'Parâmetros de conexão HTTPS do serviço',
  },
  MX: {
    chip: 'bg-amber-500/10 text-amber-600 ring-amber-500/25 dark:text-amber-400',
    hint: 'Servidor de e-mail responsável pelo domínio',
  },
  NS: {
    chip: 'bg-indigo-500/10 text-indigo-600 ring-indigo-500/25 dark:text-indigo-400',
    hint: 'Servidor de nomes autoritativo da zona',
  },
  PTR: {
    chip: 'bg-pink-500/10 text-pink-600 ring-pink-500/25 dark:text-pink-400',
    hint: 'Resolução reversa: do IP para o nome',
  },
  TXT: {
    chip: 'bg-slate-500/10 text-slate-600 ring-slate-500/25 dark:text-slate-300',
    hint: 'Texto livre: SPF, DKIM, verificações de domínio',
  },
  SRV: {
    chip: 'bg-fuchsia-500/10 text-fuchsia-600 ring-fuchsia-500/25 dark:text-fuchsia-400',
    hint: 'Host e porta de um serviço específico',
  },
  TLSA: {
    chip: 'bg-red-500/10 text-red-600 ring-red-500/25 dark:text-red-400',
    hint: 'Associa um certificado TLS ao nome (DANE)',
  },
}

const FALLBACK: RecordTypeMeta = {
  chip: 'bg-slate-500/10 text-slate-600 ring-slate-500/25 dark:text-slate-300',
  hint: '',
}

export function useRecordTypes(): {
  recordTypes: ComputedRef<string[]>
  typeMeta: (type: string | undefined) => RecordTypeMeta
} {
  const recordTypes = computed(() => Object.keys(META))

  function typeMeta(type: string | undefined): RecordTypeMeta {
    if (!type) return FALLBACK
    return META[type] ?? FALLBACK
  }

  return { recordTypes, typeMeta }
}
