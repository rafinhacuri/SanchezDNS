<script setup lang="ts">
  defineProps({
    eyebrow: { type: String, default: '' },
    title: { type: String, required: true },
    description: { type: String, default: '' },
    mono: { type: Boolean, default: false },
  })
</script>

<template>
  <div class="flex flex-col gap-4 pb-5 lg:flex-row lg:items-end lg:justify-between">
    <div class="min-w-0 space-y-1.5">
      <p v-if="eyebrow" class="text-[11px] font-semibold tracking-[0.18em] text-dimmed uppercase">
        {{ eyebrow }}
      </p>

      <h1
        class="truncate text-2xl font-bold text-highlighted sm:text-3xl"
        :class="mono ? 'data tracking-tight' : ''">
        {{ title }}
      </h1>

      <p v-if="description" class="max-w-prose text-sm text-muted">
        {{ description }}
      </p>
    </div>

    <div v-if="$slots.actions" class="flex shrink-0 flex-wrap items-center gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>
