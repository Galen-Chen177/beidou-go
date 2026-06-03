# internal/

内部包，不对外暴露（Go 的 `internal` 机制阻止外部模块引用）。

## 目录

| 目录 | 职责 | 详见 |
|---|---|---|
| `network/` | TCP 网络层：服务器监听、会话管理、封包编解码 | [network/README.md](network/README.md) |
| `crypto/` | 加解密层：AES/OFB + Shanda | [crypto/README.md](crypto/README.md) |
| `opcode/` | 封包操作码（opcode）常量定义 | [opcode/README.md](opcode/README.md) |
| `model/` | 数据模型（对应数据库表） | [model/README.md](model/README.md) |
| `store/` | 数据访问层：数据库初始化、gorm 操作 | [store/README.md](store/README.md) |
| `server/` | 游戏服务器实例（Login + Channel） | [server/README.md](server/README.md) |
| `script/` | JS 脚本引擎：goja 封装 + 桥接层 | [script/README.md](script/README.md) |
| `service/` | 业务逻辑层（占位，待实现） | |
| `api/` | REST API：管理后台接口（gin） | [api/README.md](api/README.md) |

## 依赖方向（从上到下）

```
server/   ← 入口（处理客户端连接）
  ├── service/    ← 业务逻辑
  ├── script/     ← 脚本引擎
  ├── store/      ← 数据访问
  │     └── model/  ← 数据模型
  ├── opcode/     ← 封包定义
  ├── crypto/     ← 加解密
  └── network/    ← TCP/封包
```
