# internal/api/

REST API 层，基于 [gin](https://github.com/gin-gonic/gin) 提供管理后台接口。

## 文件

| 文件 | 职责 |
|---|---|
| `router.go` | Gin 路由定义 + handler 实现（当前为骨架）：`/health` 健康检查、`/api/v1/accounts` 账号 CRUD、`/api/v1/characters` 角色管理、`/api/v1/server/info` 服务器状态 |

## API 列表（规划）

```
GET    /health                    # 健康检查
GET    /api/v1/accounts           # 账号列表
POST   /api/v1/accounts           # 创建账号
PUT    /api/v1/accounts/:id       # 更新账号
DELETE /api/v1/accounts/:id       # 删除账号
GET    /api/v1/characters         # 角色列表
PUT    /api/v1/characters/:id     # 更新角色
GET    /api/v1/server/info        # 服务器信息（在线人数、运行状态等）
```

## 注意

- API 层对应原项目的 `gms-ui/` Vue.js 管理后台
- 当前 handler 为骨架实现，返回空数据或 "not implemented"
- 后续可在 `handler/` 子目录按模块拆分（`account_api.go`、`character_api.go` 等）
