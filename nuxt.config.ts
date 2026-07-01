import process from 'node:process'

import { name, version } from './package.json'

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
  devtools: { enabled: true },
  app: { head: { templateParams: { separator: '•' } } },
  css: ['~/assets/main.css'],
  site: {
    name: 'Sanchez DNS',
    description: '🗄️ Web application to manage authoritative dns servers using PowerDNS',
  },
  runtimeConfig: {
    public: {
      production: false,
      siteUrl: '',
      name,
      version,
    },
  },
  devServer: {
    host: DEV_URL,
    https: DEV_KEY && DEV_CERT ? { key: DEV_KEY, cert: DEV_CERT } : undefined,
  },
  compatibilityDate: '2026-01-26',
  nitro: {
    preset: 'bun',
    devProxy: {
      '/go': {
        target: 'http://localhost:8080',
        ws: true,
        changeOrigin: true,
      },
    },
  },
  linkChecker: { enabled: false },
  ogImage: {
    enabled: false,
  },
  security: {
    xssValidator: false,
    headers: {
      contentSecurityPolicy: {
        'img-src': ["'self'", 'data:', 'blob:'],
      },
    },
  },
})
