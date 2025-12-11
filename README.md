# ZZFrame

一个简洁、高效的 Golang Web 开发框架。

## 项目结构

```
zzframe/
├── cmd/server/         # 应用入口
├── internal/           # 私有代码
│   ├── config/         # 配置管理
│   ├── handler/        # 请求处理器
│   ├── service/        # 业务逻辑
│   ├── repository/     # 数据访问
│   └── model/          # 数据模型
├── pkg/utils/          # 公共工具
├── api/                # API 定义
├── configs/            # 配置文件
└── docs/               # 文档
```

## 快速开始

### 环境要求

- Go 1.24+

### 安装

```bash
git clone https://github.com/denghuo98/zzframe.git
cd zzframe
go mod tidy
```

### 运行

```bash
go run cmd/server/main.go
```

### 配置

通过环境变量配置：

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| SERVER_PORT | 服务端口 | 8080 |
| SERVER_MODE | 运行模式 | debug |
| DB_HOST | 数据库地址 | localhost |
| DB_PORT | 数据库端口 | 3306 |
| DB_USER | 数据库用户 | root |
| DB_PASSWORD | 数据库密码 | - |
| DB_NAME | 数据库名 | zzframe |
| LOG_LEVEL | 日志级别 | info |

## 开发规范

详见 [.cursorrules](.cursorrules)

## License

MIT

