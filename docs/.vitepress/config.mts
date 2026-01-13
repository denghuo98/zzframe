import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  base: '/zzframe/',
  title: "ZZFrame",
  description: "ZZ 科技开发框架",
  ignoreDeadLinks: true,
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: '主页', link: '/' },
      { text: '快速开始', link: '/quick-start' },
      { text: '开发手册', link: '/development-manual' }
    ],

    sidebar: [
      {
        text: '快速开始',
        items: [
          { text: '快速开始', link: '/quick-start' }
        ]
      },
      {
        text: '开发手册',
        items: [
          { text: '开发手册', link: '/development-manual' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/denghuo98/zzframe' }
    ]
  }
})
