# internal/crypto/

加解密层，实现冒险岛 GMS v0.83 的网络封包加密方案。

## 文件

| 文件 | 职责 |
|---|---|
| `aes.go` | `MapleCrypto` — AES/OFB 加解密器。握手阶段从客户端获取 IV，经过版本变换后初始化 OFB 流密码，后续所有封包通过 `Encrypt`/`Decrypt` 进行 XOR 加解密 |
| `shanda.go` | Shanda 自定义加密/解密 — 占位实现（GMS v0.83 可能已不使用，需对照原 Java 项目确认） |

## 加解密流程

```
客户端握手 → 客户端发送 IV
                ↓
          MapleCrypto.New(AES key=IV, version=83)
                ↓
         transformIV(IV, version)  ← MapleStory 自定义 IV 变换
                ↓
         创建 AES/OFB 加密流
                ↓
         正常通信: send: XorEncrypt / recv: XorDecrypt
```

## 注意

- AES/OFB 的 IV 在这里同时充当 AES 密钥，这是 MapleStory 的特殊设计
- IV 变换算法（`transformIV`）需对照原项目验证，不同版本的变换规则可能不同
- Shanda 加密在 v0.83 中可能已弃用，待确认
