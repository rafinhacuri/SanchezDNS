<script setup lang="ts">
const error = useError()

const statusCode = ref(500)
const message = ref('Something went wrong...')

if(error.value && 'statusCode' in error.value && 'message' in error.value){
  statusCode.value = error.value.statusCode
  message.value = error.value.message
}

useHead({
  title: String(statusCode.value),
})
</script>

<template>
  <UError
    :error="{
      statusCode: statusCode,
      statusMessage: message,
      message: message,
    }"
  />
</template>
