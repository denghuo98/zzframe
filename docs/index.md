---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "ZZFrame"
  text: "简洁高效的 Golang Web 开发框架"
  tagline: 开箱即用的企业级后端开发框架
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/
    - theme: alt
      text: 5分钟搭建
      link: /guide/quick-start

features:
  - title: 开箱即用
    details: 内置用户管理、角色管理、权限控制、登录日志等企业级功能，快速搭建后台管理系统
  - title: 架构清晰
    details: 采用分层架构设计，Controller-Service-Repository 分层明确，代码易于维护和扩展
  - title: 完善的认证授权
    details: 基于 Token 的认证机制，集成了 Casbin 权限控制，支持灵活的权限配置
  - title: 丰富的功能组件
    details: 内置缓存、队列、验证码、加密、日志队列等多种功能组件，满足各种业务需求
  - title: 灵活配置
    details: 支持 YAML、JSON、TOML 等多种配置格式，配置项清晰易懂
  - title: 生产就绪
    details: 完善的日志系统、错误处理、连接池管理，可直接用于生产环境
---
