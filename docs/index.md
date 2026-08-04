---
layout: home

hero:
  name: 'SanchezDNS'
  tagline: 'Plataforma web para operar zones, records, usuários e auditoria em PowerDNS'
  image:
    src: /logo.png
    alt: 'SanchezDNS'
  actions:
    - theme: brand
      text: Começar
      link: /setup
    - theme: alt
      text: Ver no GitHub
      link: https://github.com/rafinhacuri/SanchezDNS

features:
  - title: ⚙️ Integração com PowerDNS
    details: O backend fala diretamente com a API authoritative para criar zonas, editar SOA e manipular registros sem sincronização manual.
  - title: 🧭 Controle por zona
    details: O acesso é calculado por zona, com permissões de leitura e escrita separadas e nível administrativo global.
  - title: 🔁 Reverse automático
    details: Registros A e AAAA mantêm PTR reverso automaticamente em zonas in-addr.arpa e ip6.arpa.
  - title: 🔎 Conferência de reversos
    details: Encontre de uma vez os A e AAAA sem PTR e os PTR órfãos que sobraram, e corrija em lote.
  - title: 📊 Estatísticas em tempo real
    details: Zonas, registros, uptime e tráfego UDP/TCP são atualizados periodicamente no painel.
  - title: 🪵 Auditoria
    details: Criação, edição e remoção de zonas, registros e permissões geram logs de auditoria.
  - title: 🐳 Docker pronto
    details: O projeto sobe em um stack composto por app, MongoDB, Redis e storage S3 compatível.

footer: |
  © 2026 SanchezDNS — desenvolvido por Rafael Curi.  
  Documentação feita com [VitePress](https://vitepress.dev) e interface em [Nuxt 4](https://nuxt.com).
---
