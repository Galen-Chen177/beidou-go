# internal/server/login/

登录服务器，负责客户端连接后的认证和角色选择流程。

## 文件

| 文件 | 职责 | 状态 |
|---|---|---|
| `login_server.go` | `Server` — 启动端口监听、握手 (SERVER_HELLO/CLIENT_HELLO)、加解密通道建立、dispatch() 按 opcode 分发 | ✅ 已实现 |
| `login_handler.go` | `LoginHandler` 接口 — 定义 7 个方法：`HandleCheckPassword`、`HandleServerList`、`HandleServerStatusRequest`、`HandleCharList`、`HandleCheckCharName`、`HandleCharSelect`、`HandleCharCreate` | ✅ 已实现 |
| `session_coordinator.go` | `SessionCoordinator` — 多开检测：同一账号不能同时登录两次 | ✅ 已实现 |

## dispatch 路由

```
0x01 LOGIN_CHECK_PASSWORD  → HandleCheckPassword
0x04/0x0B SERVER_LIST      → HandleServerList
0x05 CHAR_LIST             → HandleCharList
0x06 SERVER_STATUS         → HandleServerStatusRequest
0x13 CHAR_SELECT           → HandleCharSelect
0x15 CHECK_CHAR_NAME       → HandleCheckCharName
0x16 CHAR_CREATE           → HandleCharCreate
0x23 CLIENT_HELLO          → 握手 (IV 同步)
```

## 客户端交互流程

```
客户端连接 (8484)
  → SERVER_HELLO (明文 16B，含 sendIV/recvIV)
  → 客户端回 CLIENT_HELLO (0x23 加密) → IV 同步完成
  → LoginCheckPassword (0x01) → 验证用户名密码 →
    ← LoginStatus (0x00) Success / Failed
  → ServerListRequest (0x0B) →
    ← ServerList (0x0A) ×N → EndOfServerList → LastConnectedWorld → RecommendedWorlds
  → ServerStatusRequest (0x06) →
    ← ServerStatus (0x03)
  → CharListRequest (0x05) →
    ← CharList (0x0B) 含角色列表 + PIC 状态
  → CheckCharName (0x15) →
    ← CharNameResponse (0x0D)
  → CharCreate (0x16) →
    ← AddNewCharEntry (0x0E) / DeleteCharResponse (0x0F)
  → CharSelect (0x13) →
    ← ServerIP (0x0C) 含频道 IP:Port + charID
  → 客户端断开 8484，连接频道端口 7575
```

## handler/ (实现放在 internal/server/handler/)

接口方法的具体实现在 `internal/server/handler/auth.go`（`AuthHandler`），不在本目录。

## 未来规划

- [ ] 真实的服务器负载检测 (`getWorldCapacityStatus`)
- [ ] 角色删除功能 (`0x17`)
- [ ] PIC/PIN 二次密码验证
- [ ] 游客登录 (`0x02`)
- [ ] 频道负载信息从 Channel Server 动态获取
