import process from 'node:process'

const { DEV_URL, DEV_KEY, DEV_CERT } = process.env

export default defineNuxtConfig({
  modules: [
    '@nuxt/image',
    '@vueuse/nuxt',
    '@nuxtjs/seo',
    'nuxt-security',
    '@nuxt/ui',
    '@nuxt/hints',
    '@nuxt/a11y',
  ],
  $development: {
    security: { headers: { crossOriginEmbedderPolicy: 'unsafe-none' } },
  },
  routeRules: {
    '/server/**': { proxy: { to: DEV_KEY && DEV_CERT ?`https://${DEV_URL}:8080/**` : 'http://localhost:8080/**' } },
  },
  devtools: { enabled: true },
  app: { head: { templateParams: { separator: '•' } } },
  css: ['~/assets/main.css'],
  site: {
    name: 'Sanchez DNS',
    description: '🗄️ Web application to manage authoritative dns servers using PowerDNS',
  },
  devServer: {
    host: DEV_URL,
    https: DEV_KEY && DEV_CERT ? { key: DEV_KEY, cert: DEV_CERT } : undefined,
  },
  compatibilityDate: '2026-01-26',
  linkChecker: { enabled: false },
  security: {
    headers: {
      contentSecurityPolicy: {
        'img-src': ["'self'", 'https://assets.cbpf.br', 'https://assets.cbpf.dev.br', 'data:'],
      },
    },
  },
})
