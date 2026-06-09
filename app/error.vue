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
</script>

<template>
  <UError :error="error" />
</template>
