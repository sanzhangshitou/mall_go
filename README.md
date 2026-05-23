# Mall — Go 商城 API

Gin + GORM + Redis 商城后端，Laravel 风格目录，H5 与 Admin 双入口。

## 技术栈

**Gin** · **GORM** · **MySQL** · **go-redis** · **Zap** · **godotenv**

## 快速开始

```bash
# 1. 导入数据库
mysql -u root -p < database/schema.sql

# 2. 配置
cp .env.example .env   # 修改数据库/Redis 连接

# 3. 启动
go run main.go         # http://localhost:8080
curl localhost:8080/health
```

## 项目结构

```
mall/
├── main.go
├── bootstrap/app.go               # 启动：配置→日志→DB→Redis
├── routes/
│   ├── api.go                     # /api/v1（H5）
│   └── admin.go                   # /admin（管理后台）
├── config/                        # 类型化配置（app / database / redis / log / page）
├── app/
│   ├── Http/
│   │   ├── Controllers/
│   │   │   ├── Api/               # H5 控制器
│   │   │   └── Admin/             # Admin 控制器
│   │   ├── Middleware/            # CORS / Recovery / RequestLog / AdminAuth
│   │   ├── Requests/
│   │   │   ├── Api/               # H5 请求校验
│   │   │   └── Admin/             # Admin 请求校验
│   │   └── Resources/
│   │       ├── Api/               # H5 响应格式化
│   │       └── Admin/             # Admin 响应格式化
│   ├── Models/                    # GORM 模型 + 枚举
│   ├── Services/                  # 业务逻辑（H5 + Admin 共用）
│   └── Support/                   # 基础设施
│       ├── config/                # .env 加载
│       ├── database/              # MySQL 连接池
│       ├── cache/                 # Redis 客户端
│       ├── logger/                # Zap 日志
│       ├── response/              # 统一 JSON 响应
│       └── paginator/             # 分页
├── database/schema.sql
└── storage/logs/
```

## H5 API — `/api/v1`

统一响应：`{"code":0,"message":"ok","data":{}}`

### 分类

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/categories` | 树形列表 |
| GET | `/api/v1/categories/:id` | 详情 |
| POST | `/api/v1/categories` | 创建 |
| PUT | `/api/v1/categories/:id` | 更新 |
| DELETE | `/api/v1/categories/:id` | 软删除 |

### 商品

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/products` | 分页列表 |
| GET | `/api/v1/products/:id` | 详情（含SKU+图片） |
| POST | `/api/v1/products` | 创建 |
| PUT | `/api/v1/products/:id` | 更新 |
| DELETE | `/api/v1/products/:id` | 软删除 |

查询参数：`page` / `page_size` / `category_id` / `status` / `keyword` / `sort_field` / `sort_order`

### SKU

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/products/:id/skus` | SKU 列表 |
| GET | `/api/v1/products/:id/skus/:id` | SKU 详情 |
| POST | `/api/v1/products/:id/skus` | 创建 |
| PUT | `/api/v1/products/:id/skus/:id` | 更新 |
| DELETE | `/api/v1/products/:id/skus/:id` | 软删除 |

## Admin API — `/admin`

需 `Authorization` Header（中间件已预留），`/admin/login` 不受限。

### 认证

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/admin/login` | 登录 `{"username":"admin","password":"admin123"}` |

### 仪表盘

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/admin/dashboard` | 商品/SKU/分类数量统计 |

### 分类

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/admin/categories` | 列表 `?include_disabled=true&include_deleted=true` |
| GET | `/admin/categories/:id` | 详情（含已删除） |
| POST | `/admin/categories` | 创建 |
| PUT | `/admin/categories/:id` | 更新 |
| DELETE | `/admin/categories/:id` | 软删除 |
| PATCH | `/admin/categories/:id/restore` | 恢复 |

### 商品

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/admin/products` | 分页列表 `?include_deleted=true` |
| GET | `/admin/products/:id` | 详情 |
| POST | `/admin/products` | 创建 |
| PUT | `/admin/products/:id` | 更新 |
| DELETE | `/admin/products/:id` | 软删除 |
| PATCH | `/admin/products/:id/restore` | 恢复 |
| PATCH | `/admin/products/:id/status` | 改状态 `{"status":1}` |

### SKU

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/admin/products/:id/skus` | 列表（含已删除） |
| GET | `/admin/products/:id/skus/:skuId` | 详情 |
| POST | `/admin/products/:id/skus` | 创建 |
| PUT | `/admin/products/:id/skus/:skuId` | 更新 |
| DELETE | `/admin/products/:id/skus/:skuId` | 软删除 |
| PATCH | `/admin/products/:id/skus/:skuId/restore` | 恢复 |

## 数据库

| 表 | 说明 |
|---|---|
| `product_category` | 分类（树形，parent_id 自关联） |
| `product` | 商品主表 |
| `product_sku` | 多规格库存/价格 |
| `product_image` | 图片（支持按 SKU 关联） |

全部使用 `is_deleted` 软删除。

## 架构

```
请求 → Middleware → Controller → Request(校验) → Service(业务) → Resource(转换) → JSON
```

查询策略：主表分页 → 提取外键集合 → 批量查关联表 → 内存组装。仅筛选/排序/聚合依赖关联表时才连表。

枚举使用 `app/Models/product_enums.go`，禁止魔法值：

```go
db.Where("status = ?", models.ProductStatusOnSale)   // ✅
db.Where("status = ?", 1)                             // ❌
```

## 开发命令

```bash
go run main.go                 # 启动
go build -o mall.exe main.go   # 编译
go fmt ./...                   # 格式化
goimports -w -local mall ./    # 整理 import
go vet ./...                   # 静态检查
go mod tidy                    # 整理依赖
```

## 配置

完整见 `.env.example`。关键项：

| 变量 | 默认值 |
|---|---|
| `APP_PORT` | `8080` |
| `APP_DEBUG` | `true` |
| `DB_HOST` / `DB_PORT` / `DB_DATABASE` | `127.0.0.1:3306/my_mall` |
| `REDIS_HOST` / `REDIS_PORT` | `127.0.0.1:6379` |
| `LOG_PATH` / `LOG_LEVEL` | `storage/logs/` / `debug` |
| `PAGE_DEFAULT_SIZE` | `15` |

## 相关文档

- [AGENTS.md](./AGENTS.md) — 开发者规范

## 许可证

[MIT](./LICENSE)
