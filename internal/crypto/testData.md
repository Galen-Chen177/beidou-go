# 客户端抓包数据

Wireshark 抓取的真实客户端前三包，用于验证加解密实现。

## 连接信息

- 客户端: 冒险岛 GMS v83 韩服客户端（本地登录器）
- 服务端: beidou-go LoginServer
- 抓包工具: Wireshark（TCP payload 十六进制 dump）

## 数据流

```
TCP 三次握手
  │
  ├─[1] SERVER_HELLO  服务端→客户端  明文，携带 sendIv / recvIv
  │
  ├─[2] CLIENT_HELLO  客户端→服务端  加密，opcode=0x23
  │
  └─[3] LOGIN         客户端→服务端  加密，opcode=0x01 (LoginCheckPassword)
```

## 包详情

### 1. SERVER_HELLO (16 bytes, 明文)

```
0000   0e 00 53 00 01 00 31 46 72 7a 60 52 30 78 c1 08
```

| 偏移 | 字节 | 含义 |
|------|------|------|
| 0 | 0E 00 | opcode = 0x000E |
| 2 | 53 00 | version = 0x0053 = 83 |
| 4 | 01 00 | flag = 1 |
| 6 | 31 | subversion = '1' |
| 7 | 46 72 7A 60 | recvIv = "Frz\`" → 客户端用它加密上行数据 |
| 11 | 52 30 78 C1 | sendIv = "R0xÁ" → 客户端用它解密下行数据 |
| 15 | 08 | terminator |

### 2. CLIENT_HELLO (6 bytes, 加密)

```
0000   29 60 2b 60 03 a2
```

- 前 4 字节 (29 60 2B 60) = 加密包头
- 后 2 字节 (03 A2) = 加密包体
- 解密后: `23 00` → opcode=0x0023（客户端握手响应）

### 3. LOGIN (51 bytes, 加密)

```
0000   97 1a b8 1a 58 ba 5b f7 ad c8 b3 ba c2 8e e6 7f
0010   bc fa 66 01 d1 26 7e 0d 7d db bf 70 e5 74 6b b4
0020   59 a2 58 db 00 24 02 ef 6d f9 43 1f e4 b7 18 be
0030   31 8b 13
```

- 前 4 字节 (97 1A B8 1A) = 加密包头
- 后 47 字节 = 加密包体
- 解密后: opcode=0x0001，账号="a123456"，密码="a123456"
