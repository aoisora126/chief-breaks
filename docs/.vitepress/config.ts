import { defineConfig } from 'vitepress'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  title: 'Chief',
  description: 'Autonomous PRD Agent',
  base: '/chief/',

  // Force dark mode only
  appearance: 'force-dark',

  vite: {
    plugins: [tailwindcss()]
  },

  markdown: {
    theme: 'tokyo-night'
  },

  themeConfig: {
    siteTitle: 'Chief',

    nav: [
      { text: 'Home', link: '/' },
      { text: 'Docs', link: '/guide/quick-start' },
      { text: 'GitHub', link: 'https://github.com/minicodemonkey/chief' }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/minicodemonkey/chief' }
    ],

    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Quick Start', link: '/guide/quick-start' },
          { text: 'Installation', link: '/guide/installation' }
        ]
      },
      {
        text: 'Concepts',
        items: [
          { text: 'How Chief Works', link: '/concepts/how-it-works' },
          { text: 'The Ralph Loop', link: '/concepts/ralph-loop' },
          { text: 'PRD Format', link: '/concepts/prd-format' },
          { text: 'The .chief Directory', link: '/concepts/chief-directory' }
        ]
      },
      {
        text: 'Reference',
        items: [
          { text: 'CLI Commands', link: '/reference/cli' },
          { text: 'Configuration', link: '/reference/configuration' },
          { text: 'PRD Schema', link: '/reference/prd-schema' }
        ]
      },
      {
        text: 'Troubleshooting',
        items: [
          { text: 'Common Issues', link: '/troubleshooting/common-issues' },
          { text: 'FAQ', link: '/troubleshooting/faq' }
        ]
      }
    ]
  }
})
