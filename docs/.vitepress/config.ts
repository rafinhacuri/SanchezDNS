import { defineConfig } from 'vitepress'

import { description, version } from '../../package.json'

export default defineConfig({
  cleanUrls: true,
  description:
    description ||
    'A modern web interface for PowerDNS with real-time control, automation, and security.',
  head: [
    ['meta', { name: 'theme-color', content: '#00c850' }],
    ['meta', { name: 'og:type', content: 'website' }],
    ['meta', { name: 'og:locale', content: 'en-US' }],
    ['meta', { name: 'og:site_name', content: 'SanchezDNS' }],
    ['link', { rel: 'icon', href: '/logo.png', type: 'image/png' }],
  ],
  lang: 'en-US',
  lastUpdated: true,
  sitemap: {
    hostname: 'https://sanchezdns.curi.dev.br',
  },
  themeConfig: {
    siteTitle: 'SanchezDNS',
    logo: '/logo.png',
    search: {
      provider: 'local',
    },
    socialLinks: [{ icon: 'github', link: 'https://github.com/rafinhacuri/sanchezdns' }],
    nav: [
      { text: 'Home', link: '/' },
      {
        text: `V${version}`,
        items: [
          {
            text: 'Changelog',
            link: 'https://github.com/rafinhacuri/sanchezdns/releases',
            target: '_blank',
          },
          {
            text: 'Report a Bug',
            link: 'https://github.com/rafinhacuri/sanchezdns/issues',
            target: '_blank',
          },
          { text: 'Sponsor', link: 'https://github.com/sponsors/rafinhacuri', target: '_blank' },
        ],
      },
    ],
    sidebar: [
      {
        text: '📘 Introduction',
        items: [
          { text: 'Why SanchezDNS', link: '/reason' },
          { text: 'Setup', link: '/setup' },
        ],
      },
      {
        text: '⚙️ Interface Overview',
        items: [
          { text: 'Configuration', link: '/configuration' },
          { text: 'Zones', link: '/zones' },
          { text: 'Statistics', link: '/statistics' },
          { text: 'Users', link: '/users' },
          { text: 'Logs', link: '/logs' },
        ],
      },
      {
        text: '📄 Resources',
        items: [
          {
            text: 'License',
            link: 'https://github.com/rafinhacuri/sanchezdns/blob/main/LICENSE',
            target: '_blank',
          },
        ],
      },
    ],
    editLink: {
      pattern: 'https://github.com/rafinhacuri/sanchezdns/edit/main/docs/:path',
      text: 'Suggest an edit on GitHub',
    },
    docFooter: {
      prev: '← Previous',
      next: 'Next →',
    },
    footer: {
      message: 'Open‑source project licensed under MIT.',
      copyright: '© 2025 Rafael Curi — SanchezDNS',
    },
  },
  title: 'SanchezDNS',
})
