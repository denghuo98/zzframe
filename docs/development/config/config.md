# 系统配置

本章节详细介绍了系统中的各个配置，配置可以使用多种类型 json/yaml/toml 等都可以。一个详细的实例配置如下:
```yaml
system:
    superAdmin:
        password: "zzframe@admin"
```


## 系统管理员（超级管理员）

系统在每次启动或者运行的时候都会进行初始化工作，即自动添加超级管理员的**角色**和**用户**。 当在配置文件中修改用户名或者密码时，需要重启才能生效。


## 登录配置
用于配置用户登录后的规则

```yaml
system:
    token:
        # 令牌的过期时间，单位: 秒  604800 默认是 7 天
        expires: 604800
        # 令牌的刷新时间，单位: 秒
        refreshInterval: 3600
        # 令牌的最大刷新次数
        maxRefreshTimes: 10
        # 令牌的密钥
        secretKey: "zzframe"
        # 是否允许多端登录
        multiLogin: true
```