import process from 'node:process'
import { description, version } from './package.json'

const { SITE_URL, DEV_URL, DEV_KEY, DEV_CERT } = process.env

export default defineNuxtConfig({
  modules: ['@nuxt/image', '@nuxt/ui', '@nuxtjs/seo', '@vueuse/nuxt', 'nuxt-security'],
  $development: {
    security: { headers: { crossOriginEmbedderPolicy: 'unsafe-none' } },
  },
  devtools: { enabled: true },
  app: { head: { templateParams: { separator: '•' } } },
  css: ['~/assets/main.css'],
  site: {
    url: SITE_URL,
    name: 'SanchezDNS',
    description,
  },
  runtimeConfig: {
    public: { version },
  },
  routeRules: {
    '/server/**': { proxy: { to: 'http://localhost:8080/**' } },
  },
  devServer: {
    host: DEV_URL,
    https: DEV_KEY && DEV_CERT ? { key: DEV_KEY, cert: DEV_CERT } : undefined,
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
