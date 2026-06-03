# internal/network/

TCP 网络层，负责端口监听、客户端连接的会话管理、封包在 TCP 流上的收发。

## 文件

| 文件 | 职责 |
|---|---|
| `tcp_server.go` | `TCPServer` — 核心网络引擎。支持多端口同时监听，每个端口注册独立的 `SessionHandler`。管理所有活跃 `Session`，提供 `Shutdown` 优雅关闭 |
| `session.go` | `Session` — 单个客户端连接的状态封装：TCP 连接、加解密器、关联的账号/角色信息、线程安全的 Send/Close |

## 子目录

| 目录 | 职责 |
|---|---|
| `codec/` | 封包编解码 — Packet 结构体、Reader/Wrtier 流式读写、TCP 粘包拆包处理 |

## 设计要点

- 每个客户端连接 = 1 个 goroutine，不设连接池（Go 的 goroutine 足够轻量）
- `TCPServer` 负责连接管理和分发，`Session` 负责单个连接的状态
- 网络层不关心业务协议，只负责字节流收发和加解密
