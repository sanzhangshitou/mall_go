# Mall — Go 商城 API

基于 **Gin + GORM + Redis** 的商城后端 API，采用 Laravel 风格目录结构，适合多人协作的大型项目。

## 技术栈

| 组件 | 选型 |
|---|---|
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) 1.10 |
| ORM | [GORM](https://gorm.io) + MySQL 8.0 |
| 缓存 | [go-redis](https://github.com/go-redis/redis) 8 |
| 日志 | [Zap](https://github.com/uber-go/zap) + [Lumberjack](https://github.com/natefinch/lumberjack) 日志轮转 |
| 配置 | [godotenv](https://github.com/joho/godotenv) `.env` 管理 |

## 快速开始

### 环境要求

- Go 1.22+
- MySQL 8.0+
- Redis 6.0+

### 初始化

```bash
# 1. 导入数据库
mysql -u root -p < database/schema.sql

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env，修改数据库和 Redis 连接信息

# 3. 安装依赖
go mod tidy

# 4. 启动
make dev
```

服务默认监听 `http://localhost:8080`：

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## 项目结构

```
mall/
├── main.go                    # 应用入口，优雅启停
├── bootstrap/app.go           # 启动引导（配置→日志→DB→Redis）
├── routes/api.go              # API 路由注册
│
├── config/                    # 类型化配置（对应 Laravel config/*.php）
│   ├── app.go                 #   AppConfig
│   ├── database.go            #   DatabaseConfig
│   ├── redis.go               #   RedisConfig
│   ├── log.go                 #   LogConfig
│   └── page.go                #   PageConfig
│
├── app/
│   ├── Http/                  # HTTP 层
│   │   ├── Controllers/
│   │   │   ├── Api/           #   H5/小程序 API 控制器
│   │   │   └── Admin/         #   管理后台 API 控制器
│   │   ├── Middleware/        #   CORS / Recovery / 请求日志 / AdminAuth
│   │   ├── Requests/
│   │   │   ├── Api/           #   H5 请求校验 DTO
│   │   │   └── Admin/         #   管理后台请求校验 DTO
│   │   ├── Resources/
│   │   │   ├── Api/           #   H5 响应格式化
│   │   │   └── Admin/         #   管理后台响应格式化
│   ├── Models/                # GORM 模型 + 枚举常量
│   └── Services/              # 业务逻辑层（H5 + Admin 共用）
│
├── app/
│   ├── Http/                  # HTTP 层（控制器/中间件/请求/响应）
│   ├── Models/                # GORM 模型 + 枚举常量
│   ├── Services/              # 业务逻辑层（H5 + Admin 共用）
│   └── Support/               # 基础设施支持包
│       ├── config/            #   .env 加载器
│       ├── database/          #   GORM MySQL 连接池
│       ├── cache/             #   Redis 客户端 + 便捷方法
│       ├── logger/            #   Zap 日志实例
│       ├── response/          #   统一 JSON 响应
│       ├── paginator/         #   分页参数解析
│       └── helper/            #   辅助函数
│
├── database/schema.sql        # 数据库结构 + 初始种子数据
├── storage/logs/              # 日志输出目录
├── .env.example               # 环境变量模板
├── .editorconfig              # 跨编辑器编码风格
├── .gitattributes             # Git 换行符归一化
├── Makefile                   # 构建脚本
└── AGENTS.md                  # 开发者贡献指南
```

## API 文档

统一响应格式：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

### 商品分类

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/api/v1/categories` | 分类列表（树形） |
| `GET` | `/api/v1/categories/:id` | 分类详情 |
| `POST` | `/api/v1/categories` | 创建分类 |
| `PUT` | `/api/v1/categories/:id` | 更新分类 |
| `DELETE` | `/api/v1/categories/:id` | 删除分类（软删除） |

### 商品

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/api/v1/products` | 商品分页列表 |
| `GET` | `/api/v1/products/:id` | 商品详情（含 SKU + 图片） |
| `POST` | `/api/v1/products` | 创建商品 |
| `PUT` | `/api/v1/products/:id` | 更新商品 |
| `DELETE` | `/api/v1/products/:id` | 删除商品（软删除） |

商品列表查询参数：

| 参数 | 类型 | 默认值 | 说明 |
|---|---|---|---|
| `page` | int | 1 | 页码 |
| `page_size` | int | 15 | 每页条数（最大 100） |
| `category_id` | int | — | 按分类筛选 |
| `status` | int | — | 状态：0-草稿 1-上架 2-下架 3-禁用 |
| `keyword` | string | — | 商品名称模糊搜索 |
| `sort_field` | string | — | 排序字段：`sales_count` / `created_at` / `price` |
| `sort_order` | string | asc | 排序方向：`asc` / `desc` |

### 商品 SKU

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/api/v1/products/:id/skus` | 商品下 SKU 列表 |
| `GET` | `/api/v1/products/:id/skus/:id` | SKU 详情 |
| `POST` | `/api/v1/products/:id/skus` | 创建 SKU |
| `PUT` | `/api/v1/products/:id/skus/:id` | 更新 SKU |
| `DELETE` | `/api/v1/products/:id/skus/:id` | 删除 SKU（软删除） |

## 管理后台 API

前缀 `/admin`，需通过 `Authorization` Header 认证（中间件已预留）。响应比 H5 接口多 `is_deleted` 等字段，列表支持 `include_deleted=true` 查看已删除记录。

### Dashboard

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/admin/dashboard` | 仪表盘统计（商品/SKU/分类数量） |

### 分类

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/admin/categories` | 列表（`?include_disabled=true&include_deleted=true`） |
| `GET` | `/admin/categories/:id` | 详情（含已删除） |
| `POST` | `/admin/categories` | 创建 |
| `PUT` | `/admin/categories/:id` | 更新 |
| `DELETE` | `/admin/categories/:id` | 软删除 |
| `PATCH` | `/admin/categories/:id/restore` | 恢复已删除 |

### 商品

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/admin/products` | 分页列表（`?include_deleted=true`） |
| `GET` | `/admin/products/:id` | 详情（含已删除） |
| `POST` | `/admin/products` | 创建 |
| `PUT` | `/admin/products/:id` | 更新 |
| `DELETE` | `/admin/products/:id` | 软删除 |
| `PATCH` | `/admin/products/:id/restore` | 恢复已删除 |
| `PATCH` | `/admin/products/:id/status` | 修改状态 `{"status": 1}` |

### 商品 SKU

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/admin/products/:id/skus` | 列表（含已删除） |
| `GET` | `/admin/products/:id/skus/:skuId` | 详情 |
| `POST` | `/admin/products/:id/skus` | 创建 |
| `PUT` | `/admin/products/:id/skus/:skuId` | 更新 |
| `DELETE` | `/admin/products/:id/skus/:skuId` | 软删除 |
| `PATCH` | `/admin/products/:id/skus/:skuId/restore` | 恢复已删除 |

## 数据库

四张核心表：

| 表 | 说明 |
|---|---|
| `product_category` | 商品分类（树形，parent_id 自关联） |
| `product` | 商品主表（名称、图片、规格、状态） |
| `product_sku` | 商品 SKU（多规格库存、价格） |
| `product_image` | 商品图片（支持按 SKU 关联） |

所有表均使用 `is_deleted` 字段实现软删除。

## 架构设计

### 查询策略

遵循"主表优先、批量补数据"原则：

```
列表：主表分页 → 提取外键集合 → 批量查询关联表 → 内存组装
详情：主表查询 → 补 SKU → 补图片 → 补分类名 → 拼接返回
```

仅当筛选条件、排序、聚合确实依赖关联表字段时才连表。

### 请求生命周期

```
请求 → Middleware → Controller → Request(校验) → Service(业务) → Resource(转换) → JSON 响应
```

### 枚举规范

状态值统一使用 `app/Models/` 下的类型常量，禁止魔法数字：

```go
// ✅ 枚举常量
db.Where("status = ?", models.ProductStatusOnSale)

// ❌ 魔法值
db.Where("status = ?", 1)
```

## 开发命令

```bash
make dev         # 开发模式运行（go run main.go）
make build       # 编译为 ./mall
make tidy        # 整理依赖
make clean       # 清理编译产物和日志

# 提交前检查
go fmt ./... && go vet ./... && go build ./...
```

## 配置项

完整列表见 `.env.example`：

| 变量 | 说明 | 默认值 |
|---|---|---|
| `APP_NAME` | 应用名称 | `Mall` |
| `APP_PORT` | HTTP 端口 | `8080` |
| `APP_DEBUG` | 调试模式（开启 GORM SQL 日志） | `true` |
| `DB_HOST` | 数据库主机 | `127.0.0.1` |
| `DB_PORT` | 数据库端口 | `3306` |
| `DB_DATABASE` | 数据库名 | `my_mall` |
| `DB_USERNAME` | 数据库用户 | `root` |
| `REDIS_HOST` | Redis 主机 | `127.0.0.1` |
| `REDIS_PORT` | Redis 端口 | `6379` |
| `LOG_PATH` | 日志目录 | `storage/logs/` |
| `LOG_LEVEL` | 日志级别 | `debug` |
| `PAGE_DEFAULT_SIZE` | 默认分页大小 | `15` |

## 相关文档

- [AGENTS.md](./AGENTS.md) — 开发者贡献指南

## 许可证

MIT
