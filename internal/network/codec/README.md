# internal/network/codec/

封包编解码层，定义冒险岛封包的内存结构和流式读写方法。

## 文件

| 文件 | 职责 |
|---|---|
| `packet.go` | `Packet` 结构体 — `Opcode uint16` + `Data []byte`。提供 `Bytes()` 获取完整包体、`Len()` 获取包体长度 |
| `reader.go` | `Reader` — 从 `[]byte` 顺序读取：`ReadByte/Short/Int/Long/String/Pos` 等。同时提供 `ReadFrom()` 处理 TCP 粘包拆包（4字节小端包头） |
| `writer.go` | `Writer` — 流式构造封包：`WriteByte/Short/Int/Long/String/Pos/Bool` 等。同时提供 `EncodePacket()` 将封包装换为网络格式（4字节包头 + 包体） |

## 封包格式

```
网络格式:  [header: 4 bytes LE] [body: opcode(2B LE) + data(NB)]
内存格式:  Packet{ Opcode uint16, Data []byte }
```

- **包头** = 包体长度（小端序），不包含自身 4 字节
- **包体** = opcode（2 字节小端序）+ 数据
- Reader/Writer 默认操作的是**包体数据**（不含包头），粘包拆包由 `ReadFrom` / `EncodePacket` 处理

## 数据类型约定（MapleStory GMS v0.83）

| 类型 | 字节数 | 说明 |
|---|---|---|
| byte | 1 | 无符号 |
| short | 2 | 小端序 |
| int | 4 | 小端序 |
| long | 8 | 小端序 |
| string | 2 + N | 前 2 字节小端序 = ASCII 长度，后 N 字节 = 内容 |
| padded string | N (固定) | 固定长度字符串（如角色名 13 字节），不足补 0x00 |
| pos | 2 + 2 | x: short, y: short |
| bool | 1 | 0=false, 1=true |
