# 北斗服务端登录流程分析

## 一、整体调用链路

```
客户端                              服务端
  │                                   │
  │──── LOGIN_PASSWORD 包 ──────────→│ LoginPasswordHandler.handlePacket()
  │                                   │   ├─ 读取 login, pwd, hwid
  │                                   │   ├─ Client.login() → DB 校验
  │                                   │   ├─ IP/MAC/TempBan 检查
  │                                   │   ├─ finishLogin() → 标记已登录
  │                                   │   └─ 调用 login(c)
  │                                   │
  │←─── LOGIN_STATUS 包 ─────────────│ PacketCreator.getAuthSuccess()
  │     (加密后)                       │   → [原始数据] → encryptData → AES-OFB → 包头
  │                                   │
  │  如果 PIN/PIC 启用:                │
  │←─── CHECK_PINCODE (mode=4) ──────│ 请求输入 PIN
  │──── PIN 应答 ──────────────────→│
  │←─── CHECK_PINCODE (mode=0) ──────│ PIN 验证通过
  │←─── CHECK_SPW_RESULT ────────────│ 请求输入 PIC (如果启用)
  │──── PIC 应答 ──────────────────→│
  │←─── CHARLIST ────────────────────│ 返回角色列表
```

---

## 二、`getAuthSuccess` 加密前的数据布局

Opcode: `LOGIN_STATUS = 0x0000`

```
偏移    大小    写入方法             值说明
────    ────    ──────────          ────────────
0       2B      opcode (short LE)   LOGIN_STATUS = 0x0000
2       4B      writeInt(0)         状态/填充
6       2B      writeShort(0)       填充
8       4B      writeInt(accID)     账号 ID (DB 主键)
12      1B      writeByte(gender)   性别 (0=男, 1=女, -1=未设)
13      1B      writeBool(isGM)     GM 权限标记 (0 或 1)
14      1B      writeByte(0x80/0)   Admin 字节 (GM 时为 0x80，否则 0)
15      1B      writeByte(0)        国家代码
16      2B+N    writeString(name)   账号名 (2 字节长度 LE + N 字节 ASCII/UTF-8)
18+N    1B      writeByte(0)        填充
19+N    1B      writeByte(0)        IsQuietBan 标志
20+N    8B      writeLong(0)        QuietBan 时间戳
28+N    8B      writeLong(0)        创建时间戳
36+N    4B      writeInt(1)         跳过世界选择 (1=跳过)
40+N    1B      writeByte(0/1)      PIN 设置: 0=启用, 1=禁用
41+N    1B      writeByte(0/1/2)    PIC 设置: 0=注册, 1=询问, 2=禁用
```

---

## 三、Golang 实现

### 3.1 常量与结构体定义

```go
package login

import (
    "encoding/binary"
)

// SendOpcode 发送包操作码
const (
    OpcodeLoginStatus = 0x0000
)

// PIC 模式常量
const (
    PicRegister = 0 // 注册 PIC
    PicAsk      = 1 // 询问 PIC
    PicDisabled = 2 // 禁用 PIC
)

// Account 账号信息——对应 Java Client 中的 getAccID/getGender/getGMLevel/getAccountName/getPic 等
type Account struct {
    AccID       int32
    Gender      byte  // 0=男, 1=女, -1=未设
    GMLevel     int32 // GM 等级，>1 表示管理员
    AccountName string
    Pin         string
    Pic         string
    CanFly      bool // 是否可以飞行 (相当于 Java canFly)
}

// ServerConfig 服务端配置——对应 Java GameConfig
type ServerConfig struct {
    UseEnforceAdminAccount bool // use_enforce_admin_account
    EnablePin              bool // enable_pin
    EnablePic              bool // enable_pic
}

// 全局配置（实际项目中应从配置文件加载）
var config = ServerConfig{
    UseEnforceAdminAccount: false,
    EnablePin:              false,
    EnablePic:              false,
}
```

### 3.2 OutPacket 写入器

```go
package login

import (
    "bytes"
    "encoding/binary"
    "math"
)

// OutPacket 封包写入器——对应 Java ByteBufOutPacket
// 所有多字节整数均使用小端序 (Little Endian)
type OutPacket struct {
    buf bytes.Buffer
}

// NewOutPacket 创建带 opcode 的封包
func NewOutPacket(opcode uint16) *OutPacket {
    p := &OutPacket{}
    p.WriteShort(opcode)
    return p
}

// Bytes 获取原始字节——对应 Java getBytes()
func (p *OutPacket) Bytes() []byte {
    return p.buf.Bytes()
}

func (p *OutPacket) WriteByte(b byte) {
    p.buf.WriteByte(b)
}

func (p *OutPacket) WriteBytes(b []byte) {
    p.buf.Write(b)
}

// WriteShort 写入 2 字节小端
func (p *OutPacket) WriteShort(v uint16) {
    tmp := make([]byte, 2)
    binary.LittleEndian.PutUint16(tmp, v)
    p.buf.Write(tmp)
}

// WriteInt 写入 4 字节小端
func (p *OutPacket) WriteInt(v int32) {
    tmp := make([]byte, 4)
    binary.LittleEndian.PutUint32(tmp, uint32(v))
    p.buf.Write(tmp)
}

// WriteLong 写入 8 字节小端
func (p *OutPacket) WriteLong(v int64) {
    tmp := make([]byte, 8)
    binary.LittleEndian.PutUint64(tmp, uint64(v))
    p.buf.Write(tmp)
}

// WriteBool 写入布尔值 (0 或 1)
func (p *OutPacket) WriteBool(v bool) {
    if v {
        p.WriteByte(1)
    } else {
        p.WriteByte(0)
    }
}

// WriteString 写入字符串: 2 字节长度 (LE) + 字符串字节
func (p *OutPacket) WriteString(v string) {
    data := []byte(v)
    length := uint16(len(data))
    if int(length) > math.MaxUint16 {
        length = math.MaxUint16
    }
    p.WriteShort(length)
    p.WriteBytes(data[:length])
}
```

### 3.3 `getAuthSuccess` 核心函数（Go 版本）

```go
package login

// GetAuthSuccess 构建登录成功应答包——对应 Java PacketCreator.getAuthSuccess()
// 返回加密前的原始字节数组
func GetAuthSuccess(acc *Account) []byte {
    p := NewOutPacket(OpcodeLoginStatus) // LOGIN_STATUS = 0x0000

    // ─────────────────────────────────────────────
    // 偏移 2: writeInt(0) — 状态/填充
    // ─────────────────────────────────────────────
    p.WriteInt(0)

    // ─────────────────────────────────────────────
    // 偏移 6: writeShort(0) — 填充
    // ─────────────────────────────────────────────
    p.WriteShort(0)

    // ─────────────────────────────────────────────
    // 偏移 8: writeInt(accID) — 账号 ID
    // ─────────────────────────────────────────────
    p.WriteInt(acc.AccID)

    // ─────────────────────────────────────────────
    // 偏移 12: writeByte(gender) — 性别
    // ─────────────────────────────────────────────
    p.WriteByte(acc.Gender)

    // ─────────────────────────────────────────────
    // 偏移 13-14: GM / Admin 标记
    //   isGM = (use_enforce_admin_account || canFly) && GMLevel > 1
    // ─────────────────────────────────────────────
    isGM := (config.UseEnforceAdminAccount || acc.CanFly) && acc.GMLevel > 1
    p.WriteBool(isGM)             // 偏移 13: 1=GM, 0=普通
    if isGM {
        p.WriteByte(0x80)         // 偏移 14: Admin 字节 (0x80/0x40/0x20)
    } else {
        p.WriteByte(0)
    }

    // ─────────────────────────────────────────────
    // 偏移 15: writeByte(0) — 国家代码
    // ─────────────────────────────────────────────
    p.WriteByte(0)

    // ─────────────────────────────────────────────
    // 偏移 16: writeString(accountName) — 账号名
    //   格式: 2 字节长度 (LE) + N 字节 ASCII/UTF-8
    // ─────────────────────────────────────────────
    p.WriteString(acc.AccountName)

    // ─────────────────────────────────────────────
    // 偏移 16+N+2: writeByte(0) — 填充
    // ─────────────────────────────────────────────
    p.WriteByte(0)

    // ─────────────────────────────────────────────
    // writeByte(0) — IsQuietBan 标志
    // ─────────────────────────────────────────────
    p.WriteByte(0)

    // ─────────────────────────────────────────────
    // writeLong(0) — QuietBan 时间戳
    // ─────────────────────────────────────────────
    p.WriteLong(0)

    // ─────────────────────────────────────────────
    // writeLong(0) — 创建时间戳
    // ─────────────────────────────────────────────
    p.WriteLong(0)

    // ─────────────────────────────────────────────
    // writeInt(1) — 跳过世界选择界面
    //   1 = 不显示 "Select the world you want to play in"
    // ─────────────────────────────────────────────
    p.WriteInt(1)

    // ─────────────────────────────────────────────
    // writeByte(pinSetting) — PIN 设置
    //   0 = PIN 系统启用
    //   1 = PIN 系统禁用
    // ─────────────────────────────────────────────
    if config.EnablePin && !canBypassPin(acc) {
        p.WriteByte(0) // PIN 启用
    } else {
        p.WriteByte(1) // PIN 禁用
    }

    // ─────────────────────────────────────────────
    // writeByte(picSetting) — PIC 设置
    //   0 = 需要注册 PIC
    //   1 = 需要输入 PIC
    //   2 = PIC 禁用
    // ─────────────────────────────────────────────
    if config.EnablePic && !canBypassPic(acc) {
        if acc.Pic == "" {
            p.WriteByte(PicRegister) // 0: 注册 PIC
        } else {
            p.WriteByte(PicAsk)      // 1: 输入 PIC
        }
    } else {
        p.WriteByte(PicDisabled)     // 2: 禁用
    }

    return p.Bytes()
}

// canBypassPin 是否可以跳过 PIN（对应 Java LoginBypassCoordinator）
// 实际项目中需要结合 HWID 和 Session 缓存判断
func canBypassPin(acc *Account) bool {
    // 简化实现：如果 PIN 为空则跳过
    return acc.Pin == ""
}

// canBypassPic 是否可以跳过 PIC
func canBypassPic(acc *Account) bool {
    // 简化实现：如果 PIC 为空则跳过
    return acc.Pic == ""
}
```

### 3.4 加密流水线（Go 版本）

对应 Java `GMSV83PacketProtocol.encode()` 的处理顺序：

```
原始字节
  → ① 生成 4 字节包头（sendCypher.getPacketHeader）
  → ② MapleCustomEncryption.encryptData()  ← 自定义混淆（6 轮）
  → ③ sendCypher.crypt()                   ← AES-OFB 流加密
  → ④ 拼装: [4字节包头] + [加密后的包体] → 发送
```

#### 3.4.1 自定义混淆加密

```go
package encrypt

// MapleEncryptData 自定义混淆加密——对应 Java MapleCustomEncryption.encryptData()
// 在原数组上直接修改，共 6 轮（偶数轮正向、奇数轮反向）
func MapleEncryptData(data []byte) {
    length := len(data)

    for j := 0; j < 6; j++ {
        var remember byte
        dataLength := byte(length & 0xFF)

        if j%2 == 0 {
            // ── 偶数轮：正向遍历 ──
            for i := 0; i < length; i++ {
                cur := data[i]
                cur = rollLeft(cur, 3)
                cur += dataLength
                cur ^= remember
                remember = cur
                cur = rollRight(cur, int(dataLength)&0xFF)
                cur = ^cur // 按位取反
                cur += 0x48
                dataLength--
                data[i] = cur
            }
        } else {
            // ── 奇数轮：反向遍历 ──
            for i := length - 1; i >= 0; i-- {
                cur := data[i]
                cur = rollLeft(cur, 4)
                cur += dataLength
                cur ^= remember
                remember = cur
                cur ^= 0x13
                cur = rollRight(cur, 3)
                dataLength--
                data[i] = cur
            }
        }
    }
}

// rollLeft 循环左移
func rollLeft(b byte, count int) byte {
    tmp := uint32(b & 0xFF)
    tmp = tmp << (count % 8)
    return byte((tmp & 0xFF) | (tmp >> 8))
}

// rollRight 循环右移
func rollRight(b byte, count int) byte {
    tmp := uint32(b & 0xFF)
    tmp = (tmp << 8) >> (count % 8)
    return byte((tmp & 0xFF) | (tmp >> 8))
}
```

#### 3.4.2 AES-OFB 流加密

```go
package encrypt

import (
    "crypto/aes"
)

// 固定 AES 密钥——对应 Java MapleAESOFB.skey
var aesKey = []byte{
    0x13, 0x00, 0x00, 0x00,
    0x08, 0x00, 0x00, 0x00,
    0x06, 0x00, 0x00, 0x00,
    0xB4, 0x00, 0x00, 0x00,
    0x1B, 0x00, 0x00, 0x00,
    0x0F, 0x00, 0x00, 0x00,
    0x33, 0x00, 0x00, 0x00,
    0x52, 0x00, 0x00, 0x00,
}

// MapleAESOFB AES-OFB 加解密器——对应 Java MapleAESOFB
type MapleAESOFB struct {
    iv           []byte       // 当前 IV (4 字节)
    mapleVersion uint16       // 版本号（已做高低字节交换）
    cipher       *aesCipher   // AES 加密器
}

// NewMapleAESOFB 创建 AES-OFB 加解密器
// sendIv: 发送方向 IV (4 字节)
// mapleVersion: 客户端版本号
func NewMapleAESOFB(sendIv []byte, mapleVersion uint16) (*MapleAESOFB, error) {
    // 版本号高低字节交换
    swappedVersion := ((mapleVersion >> 8) & 0xFF) | ((mapleVersion << 8) & 0xFF00)

    return &MapleAESOFB{
        iv:           sendIv,
        mapleVersion: swappedVersion,
    }, nil
}

// Crypt 加密/解密（对称操作）——对应 Java MapleAESOFB.crypt()
// 使用 AES-OFB 模式，每 0x5B0 字节更新一次 IV
func (m *MapleAESOFB) Crypt(data []byte) {
    remaining := len(data)
    llength := 0x5B0
    start := 0

    for remaining > 0 {
        // 扩展 IV 到 16 字节 (4x4)
        myIv := multiplyBytes(m.iv, 4, 4)

        if remaining < llength {
            llength = remaining
        }

        for x := start; x < (start + llength); x++ {
            // 每 16 字节对 myIv 做一次 AES 加密
            if (x-start)%len(myIv) == 0 {
                newIv := aesEncrypt(myIv)
                copy(myIv, newIv)
            }
            data[x] ^= myIv[(x-start)%len(myIv)]
        }

        start += llength
        remaining -= llength
        llength = 0x5B4
    }

    // 更新 IV（对应 Java updateIv → getNewIv）
    m.updateIv()
}

// GetPacketHeader 生成 4 字节包头——对应 Java MapleAESOFB.getPacketHeader()
// 包头 = IV ^ mapleVersion ^ packetLength
func (m *MapleAESOFB) GetPacketHeader(length int) []byte {
    iiv := uint16(m.iv[3])&0xFF | (uint16(m.iv[2])<<8)&0xFF00
    iiv ^= m.mapleVersion

    mlength := uint16((length<<8)&0xFF00 | (length>>8)&0xFF)

    xoredIv := iiv ^ mlength

    return []byte{
        byte((iiv >> 8) & 0xFF),
        byte(iiv & 0xFF),
        byte((xoredIv >> 8) & 0xFF),
        byte(xoredIv & 0xFF),
    }
}

// updateIv 更新 IV——对应 Java MapleAESOFB.getNewIv()
func (m *MapleAESOFB) updateIv() {
    in := []byte{0xf2, 0x53, 0x50, 0xc6}
    for x := 0; x < 4; x++ {
        funnyShit(m.iv[x], in)
    }
    m.iv = in
}

// aesEncrypt AES ECB 加密一个 16 字节块
func aesEncrypt(block []byte) []byte {
    c, err := aes.NewCipher(aesKey)
    if err != nil {
        panic(err)
    }
    result := make([]byte, 16)
    c.Encrypt(result, block)
    return result
}

// multiplyBytes 重复复制字节数组到指定长度——对应 Java multiplyBytes
func multiplyBytes(in []byte, count int, mul int) []byte {
    size := count * mul
    ret := make([]byte, size)
    for x := 0; x < size; x++ {
        ret[x] = in[x%count]
    }
    return ret
}

// ─── funnyShit / funnyBytes 表（对应 Java） ───

var funnyBytes = [256]byte{
    0xEC, 0x3F, 0x77, 0xA4, 0x45, 0xD0, 0x71, 0xBF,
    0xB7, 0x98, 0x20, 0xFC, 0x4B, 0xE9, 0xB3, 0xE1,
    0x5C, 0x22, 0xF7, 0x0C, 0x44, 0x1B, 0x81, 0xBD,
    0x63, 0x8D, 0xD4, 0xC3, 0xF2, 0x10, 0x19, 0xE0,
    0xFB, 0xA1, 0x6E, 0x66, 0xEA, 0xAE, 0xD6, 0xCE,
    0x06, 0x18, 0x4E, 0xEB, 0x78, 0x95, 0xDB, 0xBA,
    0xB6, 0x42, 0x7A, 0x2A, 0x83, 0x0B, 0x54, 0x67,
    0x6D, 0xE8, 0x65, 0xE7, 0x2F, 0x07, 0xF3, 0xAA,
    0x27, 0x7B, 0x85, 0xB0, 0x26, 0xFD, 0x8B, 0xA9,
    0xFA, 0xBE, 0xA8, 0xD7, 0xCB, 0xCC, 0x92, 0xDA,
    0xF9, 0x93, 0x60, 0x2D, 0xDD, 0xD2, 0xA2, 0x9B,
    0x39, 0x5F, 0x82, 0x21, 0x4C, 0x69, 0xF8, 0x31,
    0x87, 0xEE, 0x8E, 0xAD, 0x8C, 0x6A, 0xBC, 0xB5,
    0x6B, 0x59, 0x13, 0xF1, 0x04, 0x00, 0xF6, 0x5A,
    0x35, 0x79, 0x48, 0x8F, 0x15, 0xCD, 0x97, 0x57,
    0x12, 0x3E, 0x37, 0xFF, 0x9D, 0x4F, 0x51, 0xF5,
    0xA3, 0x70, 0xBB, 0x14, 0x75, 0xC2, 0xB8, 0x72,
    0xC0, 0xED, 0x7D, 0x68, 0xC9, 0x2E, 0x0D, 0x62,
    0x46, 0x17, 0x11, 0x4D, 0x6C, 0xC4, 0x7E, 0x53,
    0xC1, 0x25, 0xC7, 0x9A, 0x1C, 0x88, 0x58, 0x2C,
    0x89, 0xDC, 0x02, 0x64, 0x40, 0x01, 0x5D, 0x38,
    0xA5, 0xE2, 0xAF, 0x55, 0xD5, 0xEF, 0x1A, 0x7C,
    0xA7, 0x5B, 0xA6, 0x6F, 0x86, 0x9F, 0x73, 0xE6,
    0x0A, 0xDE, 0x2B, 0x99, 0x4A, 0x47, 0x9C, 0xDF,
    0x09, 0x76, 0x9E, 0x30, 0x0E, 0xE4, 0xB2, 0x94,
    0xA0, 0x3B, 0x34, 0x1D, 0x28, 0x0F, 0x36, 0xE3,
    0x23, 0xB4, 0x03, 0xD8, 0x90, 0xC8, 0x3C, 0xFE,
    0x5E, 0x32, 0x24, 0x50, 0x1F, 0x3A, 0x43, 0x8A,
    0x96, 0x41, 0x74, 0xAC, 0x52, 0x33, 0xF0, 0xD9,
    0x29, 0x80, 0xB1, 0x16, 0xD3, 0xAB, 0x91, 0xB9,
    0x84, 0x7F, 0x61, 0x1E, 0xCF, 0xC5, 0xD1, 0x56,
    0x3D, 0xCA, 0xF4, 0x05, 0xC6, 0xE5, 0x08, 0x49,
}

func funnyShit(inputByte byte, in []byte) {
    elina := in[1]
    anna := inputByte
    moritz := funnyBytes[int(elina)&0xFF]
    moritz -= inputByte
    in[0] += moritz
    moritz = in[2]
    moritz ^= funnyBytes[int(anna)&0xFF]
    elina -= moritz
    in[1] = elina
    elina = in[3]
    moritz = elina
    elina -= in[0]
    moritz = funnyBytes[int(moritz)&0xFF]
    moritz += inputByte
    moritz ^= in[2]
    in[2] = moritz
    elina += funnyBytes[int(anna)&0xFF]
    in[3] = elina

    merry := uint32(in[0]) & 0xFF
    merry |= (uint32(in[1]) << 8) & 0xFF00
    merry |= (uint32(in[2]) << 16) & 0xFF0000
    merry |= (uint32(in[3]) << 24) & 0xFF000000
    retVal := merry >> 0x1D
    merry = merry << 3
    retVal = retVal | merry

    in[0] = byte(retVal & 0xFF)
    in[1] = byte((retVal >> 8) & 0xFF)
    in[2] = byte((retVal >> 16) & 0xFF)
    in[3] = byte((retVal >> 24) & 0xFF)
}
```

#### 3.4.3 编码入口（组装加密流水线）

```go
package encrypt

// EncodePacket 完整编码一个发送包——对应 Java GMSV83PacketProtocol.encode()
// 返回可直接写入 TCP Socket 的完整字节
func EncodePacket(packet []byte, sendCypher *MapleAESOFB) []byte {
    // ── ① 生成 4 字节包头 ──
    header := sendCypher.GetPacketHeader(len(packet))

    // ── ② 自定义混淆加密（原地修改） ──
    MapleEncryptData(packet)

    // ── ③ AES-OFB 流加密（原地修改） ──
    sendCypher.Crypt(packet)

    // ── ④ 拼装: 包头 + 加密后的包体 ──
    result := make([]byte, len(header)+len(packet))
    copy(result[:4], header)
    copy(result[4:], packet)
    return result
}
```

### 3.5 完整使用示例

```go
package main

import (
    "fmt"
)

func main() {
    // 假设从数据库查出的账号信息
    acc := &Account{
        AccID:       12345,
        Gender:      0,          // 男
        GMLevel:     0,          // 普通玩家
        AccountName: "testuser",
        Pin:         "1234",
        Pic:         "5678",
        CanFly:      false,
    }

    // ── 1. 构建原始包（加密前） ──
    rawPacket := GetAuthSuccess(acc)
    fmt.Printf("原始包长度: %d bytes\n", len(rawPacket))
    fmt.Printf("原始包: % X\n", rawPacket)

    // ── 2. 初始化加密器 ──
    sendIv := []byte{0x01, 0x02, 0x03, 0x04} // 握手阶段协商的发送 IV
    sendCypher, err := NewMapleAESOFB(sendIv, 83) // GMS v83 版本
    if err != nil {
        panic(err)
    }

    // ── 3. 编码+加密 → 可直接发送给客户端 ──
    wireData := EncodePacket(rawPacket, sendCypher)
    fmt.Printf("发送数据长度: %d bytes (包头4 + 包体%d)\n", len(wireData), len(rawPacket))

    // conn.Write(wireData)  // 写入 TCP 连接
}
```

---

## 四、关键对照表

| Java 类/方法 | Go 函数/结构体 | 说明 |
|---|---|---|
| `PacketCreator.getAuthSuccess()` | `GetAuthSuccess()` | 构建登录成功包 |
| `ByteBufOutPacket` | `OutPacket` | 包写入器 |
| `OutPacket.writeShort(v)` | `WriteShort(v)` | 小端序 2 字节 |
| `OutPacket.writeInt(v)` | `WriteInt(v)` | 小端序 4 字节 |
| `OutPacket.writeLong(v)` | `WriteLong(v)` | 小端序 8 字节 |
| `OutPacket.writeBool(v)` | `WriteBool(v)` | 0 或 1 |
| `OutPacket.writeString(v)` | `WriteString(v)` | 长度(2B LE) + 内容 |
| `MapleCustomEncryption.encryptData()` | `MapleEncryptData()` | 自定义混淆 6 轮 |
| `MapleAESOFB.crypt()` | `(*MapleAESOFB).Crypt()` | AES-OFB 流加密 |
| `MapleAESOFB.getPacketHeader()` | `(*MapleAESOFB).GetPacketHeader()` | 生成 4 字节包头 |
| `GMSV83PacketProtocol.encode()` | `EncodePacket()` | 编码流水线入口 |

---

## 五、加密流程图

```
GetAuthSuccess() 返回的原始字节 (约 42 + len(accountName) 字节)
    │
    ▼
┌─────────────────────────────────────────────┐
│  ① GetPacketHeader(length)                  │
│     生成 4 字节包头                          │
│     header = IV[2:4] ^ mapleVersion ^ length │
└─────────────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────────────┐
│  ② MapleEncryptData(packet)  ← 自定义混淆   │
│     共 6 轮 (偶数正向, 奇数反向)              │
│     rollLeft/rollRight + XOR + 0x48/0x13     │
└─────────────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────────────┐
│  ③ sendCypher.Crypt(packet)  ← AES-OFB      │
│     固定 AES 密钥                            │
│     IV 每 0x5B0 字节做一次 AES 加密更新       │
│     包结束后 IV 通过 getNewIv() 更新          │
└─────────────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────────────┐
│  ④ 拼装发送                                  │
│     [4字节包头] + [加密后的包体] → TCP write  │
└─────────────────────────────────────────────┘
```

---

## 六、角色列表包 (`getCharList`) Golang 实现

对应 Java [PacketCreator.java:893-905](gms-server/src/main/java/org/gms/util/PacketCreator.java#L893-L905) `getCharList()` 及内部调用的 `addCharEntry` / `addCharStats` / `addCharLook` / `addCharEquips`。

### 6.1 角色列表数据布局（加密前）

Opcode: `CHARLIST`

```
偏移    大小    写入方法             值说明
────    ────    ──────────          ────────────
0       2B      opcode (short LE)   CHARLIST
2       1B      writeByte(status)   状态码 (正常=0)
3       1B      writeByte(count)    本区角色数量 N
        ── 以下每个角色循环 N 次 ──
        ...     addCharStats()      角色属性 (见 6.2)
        ...     addCharLook()       角色外观 (见 6.3)
        1B      writeByte(0)        分隔符
        1B      writeByte(1)        world rank 启用标记
        4B      writeInt(rank)      world 总排名
        4B      writeInt(rankMove)  排名变动 (负=下降)
        4B      writeInt(jobRank)   职业排名
        4B      writeInt(jobRankMove) 职业排名变动
        ── 角色循环结束 ──
        1B      writeByte(picSetting) PIC 设置: 0=注册, 1=询问, 2=禁用
        4B      writeInt(slots)     角色槽位总数
```

### 6.2 `addCharStats` 单角色属性布局

```
偏移    大小    说明
────    ────    ────
0       4B      writeInt(id)              角色ID
4       13B     writeFixedString(name)    角色名 (固定13字节，'\0'填充)
17      1B      writeByte(gender)         性别 (0=男,1=女)
18      1B      writeByte(skinColor)      肤色
19      4B      writeInt(face)            脸型
23      4B      writeInt(hair)            发型
27      8B*3    writeLong(petId)          宠物唯一ID ×3 (无宠物=0)
51      1B      writeByte(level)          等级
52      2B      writeShort(job)           职业
54      2B      writeShort(str)           力量
56      2B      writeShort(dex)           敏捷
58      2B      writeShort(int_)          智力
60      2B      writeShort(luk)           运气
62      2B      writeShort(hp)            当前HP
64      2B      writeShort(maxHp)         最大HP
66      2B      writeShort(mp)            当前MP
68      2B      writeShort(maxMp)         最大MP
70      2B      writeShort(remainingAp)   剩余AP
72      1B+N    addRemainingSkillInfo()   SP信息 (有SP表时) 或 2B writeShort(remainingSp)
        4B      writeInt(exp)             经验值
        2B      writeShort(fame)          人气
        4B      writeInt(gachaExp)        转蛋经验
        4B      writeInt(mapId)           当前地图
        1B      writeByte(spawnPoint)     出生点
        4B      writeInt(0)               填充
```

### 6.3 `addCharLook` + `addCharEquips` 外观布局

```
偏移    大小    说明
────    ────    ────
0       1B      writeByte(gender)         性别
1       1B      writeByte(skinColor)      肤色
2       4B      writeInt(face)            脸型
6       1B      writeBool(!mega)          非mega=1
7       4B      writeInt(hair)            发型
        ── addCharEquips ──
        循环: 每个装备 1B(pos) + 4B(itemId)
        1B      0xFF                      结束标记
        循环: 每个遮盖装备 1B(pos) + 4B(itemId)
        1B      0xFF                      结束标记
        4B      writeInt(cashWeapon)      现金武器 (0=无)
        4B*3    writeInt(petItemId)       宠物道具ID ×3 (0=无)
```

### 6.4 Golang 数据结构

```go
package charlist

// ─── 角色外观相关枚举 ───

// SkinColor 肤色 (0-4)
type SkinColor struct{ ID byte }

// ─── 宠物 ───

type Pet struct {
    UniqueID int64
    ItemID   int32
    Name     string
    Level    byte
    Tameness int16
    Fullness byte
}

// ─── 装备物品 ───

type EquipItem struct {
    ItemID   int32
    Position int16 // 负值表示装备栏位，如 -1=帽子, -111=武器
}

// ─── 角色结构体 ───

type Character struct {
    ID               int32
    Name             string
    Gender           byte   // 0=男, 1=女
    SkinColor        byte
    Face             int32
    Hair             int32
    Pets             [3]*Pet
    Level            byte
    Job              int16
    Str              int16
    Dex              int16
    Int_             int16
    Luk              int16
    HP               int16
    MaxHP            int16
    MP               int16
    MaxMP            int16
    RemainingAp      int16
    RemainingSp      int16           // 无 SP 表的职业
    RemainingSps     []byte          // 有 SP 表的职业 (各技能行剩余SP)
    HasSPTable       bool            // 是否有 SP 表 (如海盗、骑士团等)
    Exp              int32
    Fame             int16
    GachaExp         int32
    MapID            int32
    InitialSpawnPoint byte
    Rank             int32           // world 总排名
    RankMove         int32           // 排名变动
    JobRank          int32           // 职业排名
    JobRankMove      int32           // 职业排名变动
    IsGM             bool
    IsGmJob          bool
    Gender2          byte            // addCharLook 用的性别 (通常与 Gender 相同)
    Equips           []EquipItem     // 已装备物品
    MaskedEquips     []EquipItem     // 被遮盖的装备
    CashWeapon       int32           // 现金武器道具ID
}
```

### 6.5 Golang 核心实现

#### 6.5.1 `GetCharList` — 入口函数

```go
package charlist

const OpcodeCharlist = 0x000B // CHARLIST opcode (具体值取决于版本)

// GetCharList 构建角色列表应答包——对应 Java PacketCreator.getCharList()
// serverId: 服务器ID
// status:   状态码 (0=正常)
// chars:    该区角色列表
// pic:      账号 PIC
// picBypass: 是否跳过 PIC
// enablePic: 是否启用 PIC
func GetCharList(serverID int, status byte, chars []*Character,
    pic string, picBypass bool, enablePic bool, charSlots int32) []byte {

    p := NewOutPacket(OpcodeCharlist)

    // ── writeByte(status) ──
    p.WriteByte(status)

    // ── writeByte(charCount) ──
    p.WriteByte(byte(len(chars)))

    // ── 遍历每个角色 ──
    for _, chr := range chars {
        AddCharEntry(p, chr, false)
    }

    // ── PIC 设置 ──
    // 0 = Register PIC, 1 = Ask for PIC, 2 = Disabled
    if enablePic && !picBypass {
        if pic == "" {
            p.WriteByte(0) // 注册 PIC
        } else {
            p.WriteByte(1) // 输入 PIC
        }
    } else {
        p.WriteByte(2) // 禁用
    }

    // ── 角色槽位数 ──
    p.WriteInt(charSlots)

    return p.Bytes()
}
```

#### 6.5.2 `AddCharEntry` — 单个角色条目

```go
// AddCharEntry 写入单个角色条目——对应 Java PacketCreator.addCharEntry()
func AddCharEntry(p *OutPacket, chr *Character, viewall bool) {
    // ── 角色属性 ──
    AddCharStats(p, chr)

    // ── 角色外观 ──
    AddCharLook(p, chr, false)

    // ── 分隔符 ──
    if !viewall {
        p.WriteByte(0)
    }

    // GM 角色不写排名
    if chr.IsGM || chr.IsGmJob {
        p.WriteByte(0)
        return
    }

    // ── 排名信息 ──
    p.WriteByte(1) // world rank enabled
    p.WriteInt(chr.Rank)       // world 排名
    p.WriteInt(chr.RankMove)   // 排名变动 (负数=下降)
    p.WriteInt(chr.JobRank)    // 职业排名
    p.WriteInt(chr.JobRankMove) // 职业排名变动
}
```

#### 6.5.3 `AddCharStats` — 角色属性

```go
import "strings"

// AddCharStats 写入角色属性——对应 Java PacketCreator.addCharStats()
func AddCharStats(p *OutPacket, chr *Character) {
    // ── 角色 ID ──
    p.WriteInt(chr.ID)

    // ── 角色名 (固定 13 字节，'\0' 填充) ──
    name := chr.Name
    if len(name) > 13 {
        name = name[:13]
    }
    padded := name + strings.Repeat("\x00", 13-len(name))
    p.WriteFixedString(padded, 13)

    // ── 性别 ──
    p.WriteByte(chr.Gender)

    // ── 肤色 ──
    p.WriteByte(chr.SkinColor)

    // ── 脸型 ──
    p.WriteInt(chr.Face)

    // ── 发型 ──
    p.WriteInt(chr.Hair)

    // ── 宠物 ID × 3 ──
    for i := 0; i < 3; i++ {
        if chr.Pets[i] != nil {
            p.WriteLong(chr.Pets[i].UniqueID)
        } else {
            p.WriteLong(0)
        }
    }

    // ── 等级 ──
    p.WriteByte(chr.Level)

    // ── 职业 ──
    p.WriteShort(uint16(chr.Job))

    // ── 五维属性 ──
    p.WriteShort(uint16(chr.Str))
    p.WriteShort(uint16(chr.Dex))
    p.WriteShort(uint16(chr.Int_))
    p.WriteShort(uint16(chr.Luk))

    // ── HP/MP ──
    p.WriteShort(uint16(chr.HP))
    p.WriteShort(uint16(chr.MaxHP))
    p.WriteShort(uint16(chr.MP))
    p.WriteShort(uint16(chr.MaxMP))

    // ── 剩余 AP ──
    p.WriteShort(uint16(chr.RemainingAp))

    // ── 剩余 SP ──
    if chr.HasSPTable {
        AddRemainingSkillInfo(p, chr)
    } else {
        p.WriteShort(uint16(chr.RemainingSp))
    }

    // ── 经验值 ──
    p.WriteInt(chr.Exp)

    // ── 人气 ──
    p.WriteShort(uint16(chr.Fame))

    // ── 转蛋经验 ──
    p.WriteInt(chr.GachaExp)

    // ── 当前地图 ──
    p.WriteInt(chr.MapID)

    // ── 出生点 ──
    p.WriteByte(chr.InitialSpawnPoint)

    // ── 填充 ──
    p.WriteInt(0)
}

// AddRemainingSkillInfo 写入 SP 表——对应 Java PacketCreator.addRemainingSkillInfo()
// 用于有 SP 表 (HasSPTable=true) 的职业 (如海盗、骑士团)
func AddRemainingSkillInfo(p *OutPacket, chr *Character) {
    // 统计非零 SP 的技能行数
    effectiveLength := 0
    for _, sp := range chr.RemainingSps {
        if sp > 0 {
            effectiveLength++
        }
    }

    p.WriteByte(byte(effectiveLength))
    for i, sp := range chr.RemainingSps {
        if sp > 0 {
            p.WriteByte(byte(i + 1))
            p.WriteByte(sp)
        }
    }
}
```

#### 6.5.4 `AddCharLook` + `AddCharEquips` — 角色外观

```go
// AddCharLook 写入角色外观——对应 Java PacketCreator.addCharLook()
func AddCharLook(p *OutPacket, chr *Character, mega bool) {
    p.WriteByte(chr.Gender2)
    p.WriteByte(chr.SkinColor)
    p.WriteInt(chr.Face)
    p.WriteBool(!mega) // 非 mega = 1
    p.WriteInt(chr.Hair)
    AddCharEquips(p, chr)
}

// AddCharEquips 写入装备外观——对应 Java PacketCreator.addCharEquips()
//
// 装备分两组:
//   myEquip:     正面可见装备 (pos < 100)
//   maskedEquip: 被遮盖的装备 (pos > 100, 现金装备覆盖普通装备)
//
// 写入格式:
//   循环 myEquip:   1B(pos) + 4B(itemId)
//   0xFF            结束标记
//   循环 maskedEquip: 1B(pos) + 4B(itemId)
//   0xFF            结束标记
//   4B              现金武器 itemId (0=无)
//   4B×3            宠物道具 itemId (0=无)
func AddCharEquips(p *OutPacket, chr *Character) {
    // 构造 myEquip 和 maskedEquip 映射
    myEquip := make(map[int16]int32)
    maskedEquip := make(map[int16]int32)

    for _, item := range chr.Equips {
        pos := item.Position * -1 // 转为负值表示装备栏位

        if pos < 100 && pos != -111 {
            if _, exists := myEquip[pos]; !exists {
                myEquip[pos] = item.ItemID
                continue
            }
        }

        if pos > 100 && pos != 111 {
            realPos := pos - 100
            if existing, exists := myEquip[realPos]; exists {
                maskedEquip[realPos] = existing
            }
            myEquip[realPos] = item.ItemID
        } else if _, exists := myEquip[pos]; exists {
            maskedEquip[pos] = item.ItemID
        }
    }

    // ── 写入 myEquip ──
    for pos, itemID := range myEquip {
        p.WriteByte(byte(pos))
        p.WriteInt(itemID)
    }
    p.WriteByte(0xFF) // 结束标记

    // ── 写入 maskedEquip ──
    for pos, itemID := range maskedEquip {
        p.WriteByte(byte(pos))
        p.WriteInt(itemID)
    }
    p.WriteByte(0xFF) // 结束标记

    // ── 现金武器 ──
    p.WriteInt(chr.CashWeapon)

    // ── 宠物道具 ID × 3 ──
    for i := 0; i < 3; i++ {
        if chr.Pets[i] != nil {
            p.WriteInt(chr.Pets[i].ItemID)
        } else {
            p.WriteInt(0)
        }
    }
}
```

### 6.6 使用示例

```go
package main

func main() {
    // 构造角色列表
    chars := []*Character{
        {
            ID:                1,
            Name:              "TestHero",
            Gender:            0,     // 男
            SkinColor:         0,     // 正常肤色
            Face:              20000, // 脸型ID
            Hair:              30000, // 发型ID
            Level:             30,
            Job:               100,   // 战士
            Str:               100,
            Dex:               50,
            Int_:              4,
            Luk:               4,
            HP:                1500,
            MaxHP:             1500,
            MP:                200,
            MaxMP:             200,
            RemainingAp:       5,
            RemainingSp:       10,
            HasSPTable:        false,
            Exp:               12345,
            Fame:              10,
            MapID:             100000000, // 新手村
            InitialSpawnPoint: 0,
            Rank:              1,
            RankMove:          0,
            JobRank:           1,
            JobRankMove:       0,
            Equips: []EquipItem{
                {ItemID: 1000000, Position: -1},   // 帽子
                {ItemID: 1040000, Position: -5},   // 上衣
                {ItemID: 1060000, Position: -6},   // 裤子
                {ItemID: 1320000, Position: -11},  // 武器
                {ItemID: 1070000, Position: -7},   // 鞋子
            },
        },
    }

    // 构建角色列表包
    rawPacket := GetCharList(
        0,      // serverId
        0,      // status: 正常
        chars,
        "5678", // pic
        false,  // 不跳过 PIC
        true,   // 启用 PIC
        3,      // 角色槽位数
    )

    _ = rawPacket // 之后经过 encryptData → AES-OFB → 发送

    // 加密发送流程同 getAuthSuccess:
    // wireData := encrypt.EncodePacket(rawPacket, sendCypher)
    // conn.Write(wireData)
}
```

### 6.7 Java ↔ Go 对照表（角色列表相关）

| Java 类/方法 | Go 函数 | 说明 |
|---|---|---|
| `PacketCreator.getCharList(c, serverId, status)` | `GetCharList(serverID, status, chars, ...)` | 构建角色列表包 |
| `PacketCreator.addCharEntry(p, chr, viewall)` | `AddCharEntry(p, chr, viewall)` | 单角色条目 |
| `PacketCreator.addCharStats(p, chr)` | `AddCharStats(p, chr)` | 角色属性 |
| `PacketCreator.addCharLook(p, chr, mega)` | `AddCharLook(p, chr, mega)` | 角色外观 |
| `PacketCreator.addCharEquips(p, chr)` | `AddCharEquips(p, chr)` | 装备外观 |
| `PacketCreator.addRemainingSkillInfo(p, chr)` | `AddRemainingSkillInfo(p, chr)` | SP 表技能信息 |
| `StringUtil.getRightPaddedStr(name, '\0', 13)` | `name + strings.Repeat("\x00", 13-len(name))` | 13字节固定长度名称 |
| `c.loadCharacters(serverId)` | (DAO 层, 从 DB 加载) | 加载角色列表 |
| `c.getPic()` | `pic` 参数 | 获取 PIC |
| `c.canBypassPic()` | `picBypass` 参数 | 是否跳过 PIC |
| `c.getCharacterSlots()` | `charSlots` 参数 | 角色槽位数 |
