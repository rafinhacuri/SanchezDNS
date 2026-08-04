export default defineAppConfig({
  ui: {
    colors: {
      primary: 'green',
      secondary: 'emerald',
      neutral: 'slate',
      info: 'sky',
      success: 'emerald',
      warning: 'amber',
      error: 'rose',
    },
    button: { defaultVariants: { color: 'primary' } },
    input: { defaultVariants: { color: 'primary' } },
    inputNumber: { defaultVariants: { color: 'primary' } },
    select: { defaultVariants: { color: 'primary' } },
    selectMenu: { defaultVariants: { color: 'primary' } },
    textarea: { defaultVariants: { color: 'primary' } },
    table: {
      slots: {
        th: 'text-xs font-semibold uppercase tracking-wider text-dimmed',
        td: 'text-sm',
      },
    },
  },
})
