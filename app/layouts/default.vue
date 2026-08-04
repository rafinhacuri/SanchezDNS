<template>
  <div class="relative flex min-h-dvh flex-col">
    <div class="pointer-events-none absolute inset-x-0 top-0 h-105 console-grid opacity-60" />

    <Header />

    <UMain class="relative flex-1">
      <slot />
    </UMain>

    <Footer />
  </div>
</template>
