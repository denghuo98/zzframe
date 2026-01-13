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
      { text: '快速开始', link: '/quick-start/' },
      { text: '开发手册', link: '/development-manual/' },
      { text: '技术方案', link: '/technical-solutions/' }
    ],

    sidebar: {
      '/quick-start/': [
        {
          text: '快速开始',
          items: [
            { text: '简介', link: '/quick-start/' },
            { text: '环境安装', link: '/quick-start/installation' },
            { text: '项目配置', link: '/quick-start/configuration' },
            { text: '第一个应用', link: '/quick-start/first-application' },
            { text: '项目结构', link: '/quick-start/project-structure' }
          ]
        }
      ],
      '/development-manual/': [
        {
          text: '开发手册',
          items: [
            { text: '开发规范', link: '/development-manual/' },
            { text: 'Controller 开发', link: '/development-manual/controller' },
            { text: 'Service 开发', link: '/development-manual/service' },
            { text: '数据库操作', link: '/development-manual/database' },
            { text: '配置管理', link: '/development-manual/configuration' },
            { text: 'API 参考', link: '/development-manual/api' }
          ]
        }
      ],
      '/technical-solutions/': [
        {
          text: '技术方案',
          items: [
            { text: '技术方案', link: '/technical-solutions/' },
            { text: '架构设计', link: '/technical-solutions/architecture' },
            { text: '认证授权', link: '/technical-solutions/authentication' },
            { text: '权限管理', link: '/technical-solutions/authorization' },
            { text: '缓存方案', link: '/technical-solutions/cache' },
            { text: '队列方案', link: '/technical-solutions/queue' },
            { text: '日志方案', link: '/technical-solutions/logging' }
          ]
        }
      ]
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/denghuo98/zzframe' }
    ]
  }
})
