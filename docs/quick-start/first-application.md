# 第一个应用

通过一个完整的例子，学习如何使用 ZZFrame 构建应用。我们将创建一个简单的商品管理系统。

## 功能需求

- 商品列表（分页、搜索）
- 新增商品
- 编辑商品
- 删除商品

## 步骤 1：创建数据模型

创建 `internal/model/product.go`：

```go
package model

const (
    TableNameProduct = "zz_product"
)

type Product struct {
    Id          uint64    `json:"id"`
    Name        string    `json:"name"`
    Price       float64   `json:"price"`
    Stock       int       `json:"stock"`
    Description string    `json:"description"`
    Status      int       `json:"status"`
    CreatedAt   string    `json:"createdAt"`
    UpdatedAt   string    `json:"updatedAt"`
}
```

创建 `do/product.go`：

```go
package do

import (
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
)

type Product struct {
    g.Meta       `orm:"table:zz_product"`
    Id           uint64    `json:"id"`
    Name         string    `json:"name"  v:"required#商品名称不能为空"`
    Price        float64   `json:"price" v:"required#价格不能为空"`
    Stock        int       `json:"stock" v:"required#库存不能为空"`
    Description  string    `json:"description"`
    Status       int       `json:"status"`
    CreatedAt    string    `json:"createdAt"`
    UpdatedAt    string    `json:"updatedAt"`
}
```

## 步骤 2：创建 API 定义

创建 `api/product.go`：

```go
package api

// 商品列表请求
type ProductListReq struct {
    g.Meta `path:"/product/list" method:"get" tags:"商品管理"`
    Page   int    `json:"page" v:"required#页码不能为空"`
    Size   int    `json:"size" v:"required#每页数量不能为空"`
    Name   string `json:"name"` // 搜索关键字
}

type ProductListRes struct {
    List  []ProductItem `json:"list"`
    Total int64         `json:"total"`
    Page  int           `json:"page"`
    Size  int           `json:"size"`
}

type ProductItem struct {
    Id          uint64  `json:"id"`
    Name        string  `json:"name"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
    Description string  `json:"description"`
    Status      int     `json:"status"`
}

// 新增商品请求
type ProductAddReq struct {
    g.Meta      `path:"/product/add" method:"post" tags:"商品管理"`
    Name        string  `json:"name" v:"required#商品名称不能为空"`
    Price       float64 `json:"price" v:"required#价格不能为空|gt:0#价格必须大于0"`
    Stock       int     `json:"stock" v:"required#库存不能为空|gte:0#库存不能为负数"`
    Description string  `json:"description"`
}

type ProductAddRes struct {
    Id uint64 `json:"id"`
}

// 编辑商品请求
type ProductEditReq struct {
    g.Meta      `path:"/product/edit" method:"post" tags:"商品管理"`
    Id          uint64  `json:"id" v:"required#商品ID不能为空"`
    Name        string  `json:"name" v:"required#商品名称不能为空"`
    Price       float64 `json:"price" v:"required#价格不能为空|gt:0#价格必须大于0"`
    Stock       int     `json:"stock" v:"required#库存不能为空|gte:0#库存不能为负数"`
    Description string  `json:"description"`
}

type ProductEditRes struct{}

// 删除商品请求
type ProductDeleteReq struct {
    g.Meta `path:"/product/delete" method:"post" tags:"商品管理"`
    Id     uint64 `json:"id" v:"required#商品ID不能为空"`
}

type ProductDeleteRes struct{}
```

## 步骤 3：创建 DAO 层

创建 `dao/product.go`：

```go
package dao

import (
    "context"

    "github.com/gogf/gf/v2/database/gdb"
    "my-admin/internal/model/do"
)

var Product = productDao{}

type productDao struct{}

// List 获取商品列表
func (d *productDao) List(ctx context.Context, condition string, args ...interface{}) ([]*do.Product, error) {
    all, err := Product.Model(ctx).
        Where(condition, args...).
        OrderDesc("id").
        All()
    if err != nil {
        return nil, err
    }
    return all, nil
}

// Count 统计商品数量
func (d *productDao) Count(ctx context.Context, condition string, args ...interface{}) (int64, error) {
    return Product.Model(ctx).Where(condition, args...).Count()
}

// GetById 根据 ID 获取商品
func (d *productDao) GetById(ctx context.Context, id uint64) (*do.Product, error) {
    product, err := Product.Model(ctx).Where("id", id).One()
    if err != nil {
        return nil, err
    }
    return product, nil
}

// Insert 新增商品
func (d *productDao) Insert(ctx context.Context, product *do.Product) (int64, error) {
    result, err := Product.Model(ctx).Data(product).Insert()
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}

// Update 更新商品
func (d *productDao) Update(ctx context.Context, product *do.Product) (int64, error) {
    result, err := Product.Model(ctx).Where("id", product.Id).Data(product).Update()
    if err != nil {
        return 0, err
    }
    return result.RowsAffected()
}

// Delete 删除商品
func (d *productDao) Delete(ctx context.Context, id uint64) (int64, error) {
    result, err := Product.Model(ctx).Where("id", id).Delete()
    if err != nil {
        return 0, err
    }
    return result.RowsAffected()
}
```

## 步骤 4：创建 Service 层

创建 `service/product.go`：

```go
package service

import (
    "context"

    "my-admin/api"
    "my-admin/dao"
    "my-admin/internal/model/do"
)

type sProduct struct{}

var Product = sProduct{}

// List 获取商品列表
func (s *sProduct) List(ctx context.Context, req *api.ProductListReq) (res *api.ProductListRes, err error) {
    condition := "1=1"
    args := []interface{}{}

    if req.Name != "" {
        condition += " AND name LIKE ?"
        args = append(args, "%"+req.Name+"%")
    }

    list, err := dao.Product.List(ctx, condition, args...)
    if err != nil {
        return nil, err
    }

    total, err := dao.Product.Count(ctx, condition, args...)
    if err != nil {
        return nil, err
    }

    // 转换为输出格式
    items := make([]api.ProductItem, 0, len(list))
    for _, item := range list {
        items = append(items, api.ProductItem{
            Id:          item.Id,
            Name:        item.Name,
            Price:       item.Price,
            Stock:       item.Stock,
            Description: item.Description,
            Status:      item.Status,
        })
    }

    return &api.ProductListRes{
        List:  items,
        Total: total,
        Page:  req.Page,
        Size:  req.Size,
    }, nil
}

// Add 新增商品
func (s *sProduct) Add(ctx context.Context, req *api.ProductAddReq) (res *api.ProductAddRes, err error) {
    product := &do.Product{
        Name:        req.Name,
        Price:       req.Price,
        Stock:       req.Stock,
        Description: req.Description,
        Status:      1,
    }

    id, err := dao.Product.Insert(ctx, product)
    if err != nil {
        return nil, err
    }

    return &api.ProductAddRes{Id: uint64(id)}, nil
}

// Edit 编辑商品
func (s *sProduct) Edit(ctx context.Context, req *api.ProductEditReq) (res *api.ProductEditRes, err error) {
    product := &do.Product{
        Id:          req.Id,
        Name:        req.Name,
        Price:       req.Price,
        Stock:       req.Stock,
        Description: req.Description,
    }

    _, err = dao.Product.Update(ctx, product)
    if err != nil {
        return nil, err
    }

    return &api.ProductEditRes{}, nil
}

// Delete 删除商品
func (s *sProduct) Delete(ctx context.Context, req *api.ProductDeleteReq) (res *api.ProductDeleteRes, err error) {
    _, err = dao.Product.Delete(ctx, req.Id)
    if err != nil {
        return nil, err
    }

    return &api.ProductDeleteRes{}, nil
}
```

## 步骤 5：创建 Controller 层

创建 `controller/product.go`：

```go
package controller

import (
    "github.com/gogf/gf/v2/frame/g"

    "my-admin/api"
    "my-admin/service"
)

var Product = cProduct{}

type cProduct struct{}

// List 获取商品列表
func (c *cProduct) List(ctx g.Ctx, req *api.ProductListReq) (res *api.ProductListRes, err error) {
    return service.Product.List(ctx, req)
}

// Add 新增商品
func (c *cProduct) Add(ctx g.Ctx, req *api.ProductAddReq) (res *api.ProductAddRes, err error) {
    return service.Product.Add(ctx, req)
}

// Edit 编辑商品
func (c *cProduct) Edit(ctx g.Ctx, req *api.ProductEditReq) (res *api.ProductEditRes, err error) {
    return service.Product.Edit(ctx, req)
}

// Delete 删除商品
func (c *cProduct) Delete(ctx g.Ctx, req *api.ProductDeleteReq) (res *api.ProductDeleteRes, err error) {
    return service.Product.Delete(ctx, req)
}
```

## 步骤 6：配置路由

在 `main.go` 中添加路由：

```go
package main

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"

    "github.com/denghuo98/zzframe/zcmd"
    "github.com/denghuo98/zzframe/zservice"

    _ "github.com/denghuo98/zzframe/zservice/logic"
    _ "github.com/gogf/gf/contrib/drivers/mysql/v2"
    _ "github.com/gogf/gf/contrib/nosql/redis/v2"

    "my-admin/controller"
)

func main() {
    var ctx = gctx.GetInitCtx()

    if err := zservice.SystemConfig().LoadConfig(ctx); err != nil {
        g.Log().Panicf(ctx, "初始化系统配置失败: %v", err)
    }

    // 注册商品管理路由
    s := g.Server()
    s.Group("/api", func(group *ghttp.RouterGroup) {
        group.Bind(
            controller.Product,
        )
    })

    zcmd.Main.Run(ctx)
}
```

## 步骤 7：测试 API

启动服务后，可以通过 Swagger UI 测试：

```bash
# 访问 Swagger UI
http://localhost:9090/swagger
```

或者使用 curl 测试：

```bash
# 新增商品
curl -X POST http://localhost:9090/api/product/add \
  -H "Content-Type: application/json" \
  -d '{"name":"iPhone 15","price":7999,"stock":100,"description":"苹果手机"}'

# 获取商品列表
curl http://localhost:9090/api/product/list?page=1&size=10

# 编辑商品
curl -X POST http://localhost:9090/api/product/edit \
  -H "Content-Type: application/json" \
  -d '{"id":1,"name":"iPhone 15 Pro","price":8999,"stock":50}'

# 删除商品
curl -X POST http://localhost:9090/api/product/delete \
  -H "Content-Type: application/json" \
  -d '{"id":1}'
```

## 总结

通过这个例子，你应该已经掌握了：

1. 如何定义数据模型
2. 如何创建 API 接口定义
3. 如何实现 DAO 层的数据访问
4. 如何实现 Service 层的业务逻辑
5. 如何实现 Controller 层的请求处理
6. 如何配置路由

## 下一步

继续学习 [核心概念](../concepts/) 了解框架的更多特性。
