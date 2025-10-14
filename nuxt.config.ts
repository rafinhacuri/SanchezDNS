import { description, version } from './package.json'

export default defineNuxtConfig({
  modules: ['@nuxt/image', '@nuxt/ui', '@nuxtjs/seo', '@vueuse/nuxt', 'nuxt-security'],
  $development: {
    security: { headers: { crossOriginEmbedderPolicy: 'unsafe-none' } },
  },
  devtools: { enabled: true },
  app: { head: { templateParams: { separator: '•' } } },
  css: ['~/assets/main.css'],
  site: {
    name: 'SanchezDNS',
    description,
  },
  runtimeConfig: {
    public: { version },
  },
  routeRules: {
    '/server/**': { proxy: { to: 'http://localhost:8080/**' } },
  },
  compatibilityDate: '2025-07-15',
  linkChecker: { enabled: false },
  security: {
    headers: {
      contentSecurityPolicy: { 'default-src': ['\'self\''], 'img-src': ['\'self\'', 'data:', 'blob:'] },
    },
  },
  sitemap: { enabled: false },
})
