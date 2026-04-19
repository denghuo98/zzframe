# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ZZFrame is a Go Web backend framework built on [GoFrame v2](https://goframe.org). It provides a complete admin backend with RBAC authorization (Casbin), JWT authentication, message queues, caching, and structured logging.

## Common Commands

- **Run all services**: `go run examples/basic/main.go`
- **Run HTTP service only**: `go run examples/basic/main.go http`
- **Run queue service only**: `go run examples/basic/main.go queue`
- **Hot reload (requires gf CLI)**: `gf run examples/basic/main.go`
- **Run all tests**: `go test ./...`
- **Run specific package tests**: `go test ./zservice/logic/admin/...`
- **Run a single test**: `go test -run TestAdminMember_Edit ./zservice/logic/admin/...`
- **Docs dev server**: `cd docs && pnpm dev`
- **Docs build**: `cd docs && pnpm build`

## Architecture

### Layered Structure

The codebase follows a strict layered architecture. Data flows in one direction:

```
zapi/          -> API request/response struct definitions (with g.Meta route tags)
zcontroller/   -> HTTP handlers; thin layer that calls services and formats responses
zservice/      -> Service interfaces (IAdminMember, IAdminRole, etc.) and singleton accessors
zservice/logic/-> Service implementations; contains all business logic
internal/dao/  -> Data access objects (GoFrame-generated); wrap `internal/dao/internal/`
internal/model/-> `entity/` (ORM structs) and `do/` (data objects)
web/           -> Framework utilities: zresp, zcontext, ztoken, zcache, zqueue, zcasbin, etc.
zdb/           -> Database utilities and auto-migration
zqueues/       -> Queue consumer topic registrations
zconsts/       -> Application constants and error definitions
zschema/       -> DTOs shared between service and controller layers
zcmd/          -> CLI command definitions (http, queue, all)
```

### Service Registration Pattern

Services use a singleton registration pattern. Each service implementation:

1. Defines an interface in `zservice/` (e.g., `IAdminMember`)
2. Implements it in `zservice/logic/<domain>/` (e.g., `sAdminMember`)
3. Registers itself in an `init()` function:

```go
func init() {
    zservice.RegisterAdminMember(NewAdminMember())
}
```

Access the service elsewhere via `zservice.AdminMember().List(ctx, ...)`. If the singleton is not registered, it panics.

### Request Routing

Controllers are struct-based. GoFrame uses struct tags and `g.Meta` to define routes:

```go
type MemberListReq struct {
    g.Meta `path:"/member/list" method:"get" tags:"SYS-02-用户管理" summary:"获取用户列表"`
    adminSchema.MemberListInput
}
```

Controllers are bound to router groups in `zcontroller/admin.go` with middleware applied per group.

### Middleware Stack

Global middleware (applied to all routes in `zcmd/http.go`):
- `zservice.Middleware().Ctx` — initializes request context
- `zservice.Middleware().CORS` — CORS handling
- `zservice.Middleware().ResponseHandler` — standardizes JSON responses

Auth middleware (applied to `/admin/*` routes in `zcontroller/admin.go`):
- `zservice.Middleware().AdminAuth` — JWT validation + Casbin authorization

### Response Standardization

All JSON responses go through `web/zresp/response.go`. Controllers return `(res *XxxRes, err error)`. The `ResponseHandler` middleware intercepts the handler response, wraps errors into a standard `{code, message, data/error, timestamp, traceId}` JSON envelope, and writes it via `zresp.RJson()`.

### Database & Migrations

The `zdb/zmigrate` package auto-runs on startup (`init()`). It:
- Reads the database config from `config.yaml`
- Falls back to a local SQLite file (`zzframe.db`) if no config is present
- Creates tables if they do not exist, using DB-type-specific SQL (MySQL vs SQLite)

DAOs in `internal/dao/` are thin wrappers around GoFrame-generated `internal/dao/internal/` structs. Prefer `dao.Xxx.Ctx(ctx).Where(...).Scan(&dest)` for queries.

### Authentication & Authorization

- **Auth**: JWT tokens parsed from the `Authorization` header. Token config (secret, expiry) lives in `config.yaml` under `system.token`.
- **Authorization**: Casbin RBAC. The `zweb/zcasbin` package loads policies from `casbin.conf` and the database. Role-permission checks happen in `AdminAuth` middleware.
- **Super admin**: A hardcoded username (`superAdmin`) with a default role key (`superAdmin`) bypasses certain restrictions.

### Queue System

The custom queue abstraction in `web/zqueue/` supports multiple backends (disk default, Redis, RocketMQ, Kafka). Producers and consumers are instantiated per `groupName`. Consumer topics are registered in the `zqueues/` package. The queue is started as a separate goroutine when running the `queue` command or `all` command.

### Testing Conventions

Tests in `zservice/logic/<domain>/zz_*_test.go` use in-memory SQLite. Each test:
- Configures `gdb.SetConfig()` with `sqlite::@file(:memory:)?cache=shared`
- Creates required tables manually via raw SQL
- Seeds necessary reference data (e.g., roles before testing members)
- Cleans up tables in `defer`

Tests use `github.com/gogf/gf/v2/test/gtest` assertions.

## Code Conventions

From `.cursor/rules/project.mdc` and `docs/development-manual/development-standards.md`:

- **Package names**: lowercase, no underscores or mixed case (e.g., `admin`, not `admin_service`)
- **Interface names**: prefix with `I` (e.g., `IAdminMember`)
- **Errors**: wrap with `gerror.Wrap(err, zconsts.ErrorORM)` or `gerror.New("message")` for business errors. Do not return raw errors without context.
- **Comments**: exported functions/types must have comments starting with the name.
- **Git commits**: `<type>(<scope>): <subject>` where type is one of `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`.

## Important Files

- `examples/basic/main.go` — Application entry point
- `config.yaml` — Runtime configuration (DB, cache, token, queue, server port)
- `casbin.conf` — Casbin RBAC model configuration
- `zcmd/cmd.go` — CLI command tree (all, http, queue, help)
- `zcmd/http.go` — HTTP server setup and middleware binding
- `zcontroller/admin.go` — Admin route group registration
- `zservice/logic/logic.go` — Blank import to trigger all service `init()` registrations
