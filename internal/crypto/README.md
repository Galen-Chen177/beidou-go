# internal/crypto

冒险岛 GMS v83 网络封包加解密层，实现客户端与服务端之间的加密通信。

## 文件

| 文件 | 职责 |
|------|------|
| `crypto.go` | `Crypto` 接口 + `MyCrypto` 实现，高层加解密入口 |
| `maple_crypto.go` | `MapleAESOFB` 流密码核心、包头编解码、自定义混淆/反混淆、IV 轮转 |
| `crypto_test.go` | 用真实抓包数据验证解密正确性 |
| `testData.md` | 抓包数据说明（SERVER_HELLO / CLIENT_HELLO / LOGIN） |

## 加密体系

网络包经过 **两层保护**：

```
明文 body
  │
  ├─[1] EncryptData()  自定义 6 轮混淆（MapleCustomEncryption）
  │
  ├─[2] Crypt()        AES-OFB 流密码 XOR（MapleAESOFB）
  │
  └─ 密文 body + Header(4B)
```

解密时逆向操作：先 AES-OFB 解密，再反混淆。

## 握手流程

```
服务端                                    客户端
─────────────────────────────────────────────────
1. NewCrypto(version)
2. GenServerHello()
   ├─ 生成 sendIV = "R0x" + 随机
   ├─ 生成 recvIV = "Frz" + 随机
   ├─ 构造 SERVER_HELLO（明文）───▶  提取 sendIV / recvIV
   └─ 初始化 SendCipher / ReceiveCipher   初始化本地 cipher

                                     ◀─── CLIENT_HELLO（第一个加密包）
3. Decrypt(CLIENT_HELLO)
   ├─ IsValidHeader() 校验通过
   ├─ Crypt() → IV 轮转 ✓
   └─ DecryptData() → 明文 opcode=0x23

后续所有通信包都通过 Decrypt/Encrypt 处理
```

**关键点**: CLIENT_HELLO 是加密包，必须在握手阶段就创建 Crypto 并用 `Decrypt` 处理以同步 IV 状态。

## 使用方式

```go
import "beidou-go/internal/crypto"

// 1. 每个 TCP 连接创建一个 Crypto 实例
c := crypto.NewCrypto(83)           // version = 83

// 2. 生成握手包发送给客户端
hello := c.GenServerHello()
conn.Write(hello)

// 3. 读取第一个包（CLIENT_HELLO），用 Decrypt 解密以同步 IV
rawPkt := readPacket(conn)          // 4B header + body
body, err := c.Decrypt(rawPkt)
// body[0:2] = opcode (LE uint16)

// 4. 后续通信
// 解密客户端→服务端
body, err := c.Decrypt(rawPkt)

// 加密服务端→客户端（注意：必须使用返回值！）
encrypted, err := c.Encrypt(plainBody)
conn.Write(encrypted)
```

## 公共 API

### Crypto 接口

```go
type Crypto interface {
    GenServerHello() []byte                  // 生成握手包，初始化 cipher
    Decrypt([]byte) ([]byte, error)          // 解密客户端包（header+body → 明文 body）
    Encrypt([]byte) ([]byte, error)          // 加密发给客户端的包（body → header+加密body）
}
```

### 底层函数（直接操作）

| 函数 | 说明 |
|------|------|
| `NewMapleAESOFB(iv, version)` | 创建 AES-OFB 流密码实例 |
| `NewMapleAESOFB.Crypt(data)` | 原地加/解密（XOR 流密码） |
| `NewMapleAESOFB.EncodePacketHeader(len)` | 构造 4 字节包头 |
| `NewMapleAESOFB.IsValidHeader(header)` | 校验包头合法性 |
| `DecodePacketLength([]byte)` | 从包头字节切片解码长度 |
| `DecodePacketLengthUint(uint32)` | 从包头 uint32 解码长度 |
| `EncryptData(data)` | 自定义 6 轮混淆 |
| `DecryptData(data)` | 自定义 6 轮逆混淆 |

## IV 生命周期

```
初始 IV（来自 SERVER_HELLO）
  │
  ├── 包1 Crypt() → updateIv() → IV₂
  ├── 包2 Crypt() → updateIv() → IV₃
  ├── 包3 Crypt() → updateIv() → IV₄
  └── ...
```

每次 `Crypt()` 执行完毕自动调用 `getNewIv()` 轮转 IV。同一 key + 同一 IV 产生相同密钥流，因此 IV 必须逐包更新。

## 包头格式

每个加密包前有 4 字节包头，同时编码版本校验信息和包体长度：

```
Encode:
  iiv = (iv[2],iv[3]) XOR mapleVersion
  mlength = 字节交换(bodyLength)
  xoredIv = iiv XOR mlength
  输出: [iiv_hi, iiv_lo, xored_hi, xored_lo]

Decode:
  length = (header_hi16 XOR header_lo16) 字节交换还原
```

包头的生成/校验必须在 `Crypt()` 之前，因为 `Crypt()` 会轮转 IV。

## 测试

```bash
# 运行解密验证测试
go test ./internal/crypto/ -v -run TestDecryptCapturedPackets
```

测试用 `testData.md` 中的真实抓包数据，模拟完整的 SERVER_HELLO → CLIENT_HELLO → LOGIN 三包序列。
