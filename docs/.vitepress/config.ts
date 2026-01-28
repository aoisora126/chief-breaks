import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Chief',
  description: 'Autonomous PRD Agent',
  base: '/chief/',

  themeConfig: {
    siteTitle: 'Chief',

    nav: [
      { text: 'Home', link: '/' },
      { text: 'Docs', link: '/guide/' }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/minicodemonkey/chief' }
    ]
  }
})
