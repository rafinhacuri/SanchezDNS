<script setup lang="ts">
  import type { NuxtError } from '#app'

  const props = defineProps({
    error: {
      type: Object as PropType<NuxtError>,
      default() {
        return {
          status: 500,
          message: 'Server error',
        }
      },
    },
  })

  useHead({ title: String(props.error.status || 500) })

  const status = computed(() => props.error.status || 500)

  const label = computed(() =>
    status.value === 404 ? 'Página não encontrada' : 'Erro no servidor',
  )
</script>

<template>
  <div class="relative flex min-h-dvh items-center justify-center overflow-hidden px-4 py-10">
    <div class="pointer-events-none absolute inset-0 console-grid opacity-70" />

    <div class="relative w-full max-w-md text-center">
      <p class="data text-7xl font-bold tracking-tight text-primary">{{ status }}</p>

      <h1 class="mt-4 text-xl font-bold text-highlighted">{{ label }}</h1>

      <p v-if="error.message" class="mt-2 data text-sm wrap-break-word text-muted">
        {{ error.message }}
      </p>

      <UButton
        class="mt-8"
        icon="i-lucide-arrow-left"
        label="Voltar para as zonas"
        @click="clearError({ redirect: '/' })" />
    </div>
  </div>
</template>
