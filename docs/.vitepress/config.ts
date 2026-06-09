import { defineConfig } from 'vitepress'

import { description, version } from '../../package.json'

export default defineConfig({
  cleanUrls: true,
  description:
    description ||
    'Interface web para administrar PowerDNS Authoritative com automações, permissões por zona e auditoria.',
  head: [
    ['meta', { name: 'theme-color', content: '#00c850' }],
    ['meta', { name: 'og:type', content: 'website' }],
    ['meta', { name: 'og:locale', content: 'pt-BR' }],
    ['meta', { name: 'og:site_name', content: 'SanchezDNS' }],
    ['link', { rel: 'icon', href: '/logo.png', type: 'image/png' }],
  ],
  lang: 'pt-BR',
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
      { text: 'Início', link: '/' },
      {
        text: `V${version}`,
        items: [
          {
            text: 'Changelog',
            link: 'https://github.com/rafinhacuri/sanchezdns/releases',
            target: '_blank',
          },
          {
            text: 'Reportar bug',
            link: 'https://github.com/rafinhacuri/sanchezdns/issues',
            target: '_blank',
          },
          { text: 'Apoiar', link: 'https://github.com/sponsors/rafinhacuri', target: '_blank' },
        ],
      },
    ],
    sidebar: [
      {
        text: '📘 Introdução',
        items: [
          { text: 'Por que SanchezDNS', link: '/reason' },
          { text: 'Instalação', link: '/setup' },
        ],
      },
      {
        text: '⚙️ Visão Geral',
        items: [
          { text: 'Configuração', link: '/configuration' },
          { text: 'Zonas', link: '/zones' },
          { text: 'Estatísticas', link: '/statistics' },
          { text: 'Usuários', link: '/users' },
          { text: 'Logs', link: '/logs' },
        ],
      },
      {
        text: '📄 Recursos',
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
      text: 'Sugerir edição no GitHub',
    },
    docFooter: {
      prev: '← Anterior',
      next: 'Próximo →',
    },
    footer: {
      message: 'Projeto open-source sob licença MIT.',
      copyright: '© 2026 Rafael Curi — SanchezDNS',
    },
  },
  title: 'SanchezDNS',
})
