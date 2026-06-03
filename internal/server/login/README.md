# internal/server/login/

登录服务器，负责客户端连接后的认证和角色选择流程。

## 文件

| 文件 | 职责 |
|---|---|
| `login_server.go` | `Server` — 登录服务器主逻辑：`Start()` 注册端口监听，`handleConnection()` 为每个客户端连接启动处理循环（当前为骨架） |

## handler/

| 文件 | 职责 |
|---|---|
| `auth.go` | `AuthHandler` — 认证相关的封包处理（待实现）：密码验证、服务器列表、角色列表/创建/删除/选择 |

## 客户端交互流程（待实现）

```
客户端连接 (8484)
  → 握手 (Hello 封包, IV 协商)
  → LoginCheckPassword  → 验证用户名密码
  → LoginServerList     → 返回世界和频道列表
  → LoginCharListRereq  → 角色列表
  → LoginCharCreate     → 创建角色（名字检查、属性设置）
  → LoginCharSelect     → 选择角色 → 通知频道服务器 → 客户端切换端口 7575
```
