// maple_crypto.go
// 冒险岛 v83 客户端-服务端加解密的 Go 实现
// 对应 Java 文件:
//   - InitializationVector.java → IV 生成
//   - PacketCreator.getHello()  → Hello 包构造
//   - MapleAESOFB.java          → AES-OFB 加解密核心
//   - MapleCustomEncryption.java → 自定义混淆/反混淆

package main

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"math/big"
)

// ============================================================================
// 1. AES 密钥（与 Java 完全一致）
// ============================================================================

var aesKey = []byte{
	0x13, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00,
	0x06, 0x00, 0x00, 0x00, 0xB4, 0x00, 0x00, 0x00,
	0x1B, 0x00, 0x00, 0x00, 0x0F, 0x00, 0x00, 0x00,
	0x33, 0x00, 0x00, 0x00, 0x52, 0x00, 0x00, 0x00,
}

// ============================================================================
// 2. IV 生成（对应 InitializationVector.java）
// ============================================================================

// GenerateSendIV 生成服务端"发送方向"的 IV（服务端加密→客户端解密）
// Java: byte[] ivSend = {82, 48, 120, getRandomByte()};  // "R0x"
func GenerateSendIV() []byte {
	// return []byte{82, 48, 120, randomByte()} // "R0x" + 随机
	return []byte{82, 48, 120, 193}
}

// GenerateReceiveIV 生成服务端"接收方向"的 IV（客户端加密→服务端解密）
// Java: byte[] ivRecv = {70, 114, 122, getRandomByte()};  // "Frz"
func GenerateReceiveIV() []byte {
	// return []byte{70, 114, 122, randomByte()} // "Frz" + 随机
	return []byte{70, 114, 122, 96}
}

func randomByte() byte {
	n, err := rand.Int(rand.Reader, big.NewInt(255))
	if err != nil {
		return 0x60 // fallback，实际不会发生
	}
	return byte(n.Int64())
}

// ============================================================================
// 3. Hello 包（对应 PacketCreator.getHello）
// ============================================================================

// CreateHelloPacket 构造服务器发给客户端的握手包。
// 这是唯一明文发送的包，后续所有包都用加密通道。
//
// 包格式:
//
//	[2B] opcode = 0x0E
//	[2B] version (LE)
//	[2B] flag = 0x01
//	[1B] subversion = 0x31 ('1')
//	[4B] recvIv  ← 客户端收到后，用这个 IV 来加密发给服务器的数据
//	[4B] sendIv  ← 客户端收到后，用这个 IV 来解密服务器发来的数据
//	[1B] terminator = 0x08
func CreateHelloPacket(version uint16, sendIv, recvIv []byte) []byte {
	// 总长度: 2+2+2+1+4+4+1 = 16 字节
	p := make([]byte, 0, 16)

	p = append(p, 0x0E, 0x00)                           // opcode (short LE)
	p = append(p, byte(version&0xFF), byte(version>>8)) // version (short LE)
	p = append(p, 0x01, 0x00)                           // flag (short LE)
	p = append(p, 0x31)                                 // subversion byte
	p = append(p, recvIv...)                            // recvIv (4 bytes)
	p = append(p, sendIv...)                            // sendIv (4 bytes)
	p = append(p, 0x08)                                 // terminator
	return p
}

// ============================================================================
// 4. MapleAESOFB 加解密核心（对应 MapleAESOFB.java）
// ============================================================================

// funnyBytes 是 256 字节的 S-Box 查找表，Java 里是静态常量。
// 用于 IV 轮转（getNewIv）和自定义混淆（MapleCustomEncryption 不直接用它）。
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

// MapleAESOFB 是冒险岛的 AES-OFB 流密码。
// 加密和解密使用完全相同的 Crypt() 方法（XOR 对称性）。
type MapleAESOFB struct {
	iv           []byte   // 4字节 IV
	mapleVersion uint16   // 字节交换后的版本号（用于包头校验）
	block        aesBlock // AES-256 分组密码实例
}

// aesBlock 是 crypto/aes 的 Cipher 接口，这里只用了 Encrypt 方向
type aesBlock interface {
	Encrypt(dst, src []byte)
	BlockSize() int
}

// NewMapleAESOFB 创建加解密器。
//
// 参数:
//
//	iv           - 4字节初始向量（recvIv 或 sendIv）
//	             - 服务端解密客户端数据用 recvIv
//	             - 服务端加密发给客户端的数据用 sendIv
//	mapleVersion - 客户端版本号，如 83
func NewMapleAESOFB(iv []byte, mapleVersion uint16) (*MapleAESOFB, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	// Java: this.mapleVersion = (short)(((mapleVersion >> 8) & 0xFF) | ((mapleVersion << 8) & 0xFF00))
	// 效果: 交换两个字节，例如 83(0x0053) → 0x5300
	swapped := ((mapleVersion >> 8) & 0xFF) | ((mapleVersion << 8) & 0xFF00)

	// 复制 IV，防止外部修改
	ivCopy := make([]byte, 4)
	copy(ivCopy, iv)

	return &MapleAESOFB{
		iv:           ivCopy,
		mapleVersion: swapped,
		block:        block,
	}, nil
}

// IV 返回当前 IV 的副本
func (m *MapleAESOFB) IV() []byte {
	cp := make([]byte, 4)
	copy(cp, m.iv)
	return cp
}

// multiplyBytes 把 in 的前 count 个字节重复 mul 次。
// 例: multiplyBytes({a,b,c,d}, 4, 4) → {a,b,c,d, a,b,c,d, a,b,c,d, a,b,c,d} (16 bytes)
func multiplyBytes(in []byte, count, mul int) []byte {
	out := make([]byte, count*mul)
	for i := range out {
		out[i] = in[i%count]
	}
	return out
}

// Crypt 对 data 进行原地加/解密。加密和解密调用同一个方法。
//
// 算法流程:
//  1. 把 4 字节 IV 重复 4 遍 → 得到 16 字节的 AES 输入块
//  2. AES 加密这个块 → 得到 16 字节的密钥流
//  3. 密钥流 XOR 数据
//  4. 每 16 字节用上一次 AES 的输出继续生成新密钥流
//  5. 全部处理完后，轮转 IV（为下一个包做准备）
func (m *MapleAESOFB) Crypt(data []byte) {
	remaining := len(data)
	chunkSize := 0x5B0 // 第一个分块大小: 1456
	start := 0

	for remaining > 0 {
		// 初始密钥流块: IV[4] 重复 4 次 = 16 字节
		myIv := multiplyBytes(m.iv, 4, 4)

		if remaining < chunkSize {
			chunkSize = remaining
		}

		for x := start; x < start+chunkSize; x++ {
			// 每 16 字节，用 AES 刷新一次密钥流
			if (x-start)%len(myIv) == 0 {
				// AES 加密 myIv → 作为新的密钥流（同时也是下一次的输入）
				m.block.Encrypt(myIv, myIv)
			}
			// 数据 XOR 密钥流对应字节
			data[x] ^= myIv[(x-start)%len(myIv)]
		}

		start += chunkSize
		remaining -= chunkSize
		chunkSize = 0x5B4 // 后续分块大小: 1460
	}

	// 一个包处理完毕，轮转 IV 状态
	m.updateIv()
}

// ============================================================================
// 5. 包头编码/解码（对应 MapleAESOFB.getPacketHeader / isValidHeader / decodePacketLength）
// ============================================================================

// EncodePacketHeader 加密前：构造 4 字节包头。
// 包头同时编码了包体长度，并作为合法性校验。
func (m *MapleAESOFB) EncodePacketHeader(bodyLength int) []byte {
	// Java:
	// int iiv = (iv[3]) & 0xFF;
	// iiv |= (iv[2] << 8) & 0xFF00;
	// iiv ^= mapleVersion;
	iiv := (uint16(m.iv[3]) & 0xFF) | (uint16(m.iv[2])<<8)&0xFF00
	iiv ^= m.mapleVersion

	// int mlength = ((length << 8) & 0xFF00) | (length >>> 8);
	mlength := ((uint16(bodyLength) << 8) & 0xFF00) | (uint16(bodyLength) >> 8)

	// int xoredIv = iiv ^ mlength;
	xoredIv := iiv ^ mlength

	ret := make([]byte, 4)
	ret[0] = byte((iiv >> 8) & 0xFF)
	ret[1] = byte(iiv & 0xFF)
	ret[2] = byte((xoredIv >> 8) & 0xFF)
	ret[3] = byte(xoredIv & 0xFF)
	return ret
}

// IsValidHeader 校验收到的包头是否合法。
// 合法意味着: header[0] XOR iv[2] == mapleVersion_hi && header[1] XOR iv[3] == mapleVersion_lo
func (m *MapleAESOFB) IsValidHeader(header uint32) bool {
	b0 := byte((header >> 24) & 0xFF)
	b1 := byte((header >> 16) & 0xFF)
	return ((int(b0^m.iv[2]) & 0xFF) == int((m.mapleVersion>>8)&0xFF)) &&
		((int(b1^m.iv[3]) & 0xFF) == int(m.mapleVersion&0xFF))
}

// DecodePacketLength 从 4 字节包头中解码出包体长度。
// 包头可以是 []byte 或 uint32。
func DecodePacketLength(header []byte) int {
	// Java: return (((header[1] ^ header[3]) & 0xFF) << 8) | ((header[0] ^ header[2]) & 0xFF)
	return (int(header[1]^header[3]&0xFF) << 8) | int(header[0]^header[2]&0xFF)
}

// DecodePacketLengthUint 从 uint32 包头解码出包体长度。
func DecodePacketLengthUint(header uint32) int {
	// Java:
	// int length = ((header >>> 16) ^ (header & 0xFFFF));
	// length = ((length << 8) & 0xFF00) | ((length >>> 8) & 0xFF);
	length := (header >> 16) ^ (header & 0xFFFF)
	length = ((length << 8) & 0xFF00) | ((length >> 8) & 0x00FF)
	return int(length)
}

// ============================================================================
// 6. IV 轮转（对应 MapleAESOFB.getNewIv / funnyShit）
// ============================================================================

// updateIv 每次 Crypt() 执行完毕后调用，轮转 IV 状态。
func (m *MapleAESOFB) updateIv() {
	m.iv = getNewIv(m.iv)
}

// getNewIv 用固定的种子 {0xf2, 0x53, 0x50, 0xc6} 和当前 IV 进行4轮混合。
func getNewIv(oldIv []byte) []byte {
	in := []byte{0xF2, 0x53, 0x50, 0xC6}
	for x := 0; x < 4; x++ {
		funnyShit(oldIv[x], in)
	}
	return in
}

// funnyShit 是冒险岛自创的字节混合函数，完全照搬 Java 逻辑。
// 它用 funnyBytes 查找表对 in 数组进行非线性变换。
func funnyShit(inputByte byte, in []byte) {
	elina := uint16(in[1])
	anna := uint16(inputByte)

	moritz := uint16(funnyBytes[elina&0xFF])
	moritz -= uint16(inputByte)
	in[0] += byte(moritz)

	moritz = uint16(in[2])
	moritz ^= uint16(funnyBytes[anna&0xFF])
	elina -= moritz & 0xFF
	in[1] = byte(elina)

	moritz = uint16(in[3])
	elina = moritz
	elina -= uint16(in[0]) & 0xFF
	moritz = uint16(funnyBytes[moritz&0xFF])
	moritz += uint16(inputByte)
	moritz ^= uint16(in[2])
	in[2] = byte(moritz)

	elina += uint16(funnyBytes[anna&0xFF]) & 0xFF
	in[3] = byte(elina)

	merry := uint32(in[0]) & 0xFF
	merry |= (uint32(in[1]) << 8) & 0xFF00
	merry |= (uint32(in[2]) << 16) & 0xFF0000
	merry |= (uint32(in[3]) << 24) & 0xFF000000

	retVal := merry >> 0x1d
	merry = merry << 3
	retVal = retVal | merry

	in[0] = byte(retVal & 0xFF)
	in[1] = byte((retVal >> 8) & 0xFF)
	in[2] = byte((retVal >> 16) & 0xFF)
	in[3] = byte((retVal >> 24) & 0xFF)
}

// ============================================================================
// 7. 自定义混淆/反混淆（对应 MapleCustomEncryption.java）
//    AES-OFB 解密后，还要经过这一层才是真正的明文
// ============================================================================

// rollLeft 循环左移
func rollLeft(in byte, count int) byte {
	tmp := uint16(in) & 0xFF
	tmp = tmp << (count % 8)
	return byte((tmp & 0xFF) | (tmp >> 8))
}

// rollRight 循环右移
func rollRight(in byte, count int) byte {
	tmp := uint16(in) & 0xFF
	tmp = (tmp << 8) >> (count % 8)
	return byte((tmp & 0xFF) | (tmp >> 8))
}

// DecryptData 自定义解密（6轮逆混淆），在 AES-OFB 解密之后调用。
// 对应 Java 的 MapleCustomEncryption.decryptData()。
func DecryptData(data []byte) {
	for j := 1; j <= 6; j++ {
		var remember byte
		dataLength := byte(len(data) & 0xFF)
		var nextRemember byte

		if j%2 == 0 {
			// 正向遍历
			for i := 0; i < len(data); i++ {
				cur := data[i]
				cur -= 0x48
				cur = ^cur // 按位取反
				cur = rollLeft(cur, int(dataLength)&0xFF)
				nextRemember = cur
				cur ^= remember
				remember = nextRemember
				cur -= dataLength
				cur = rollRight(cur, 3)
				data[i] = cur
				dataLength--
			}
		} else {
			// 反向遍历
			for i := len(data) - 1; i >= 0; i-- {
				cur := data[i]
				cur = rollLeft(cur, 3)
				cur ^= 0x13
				nextRemember = cur
				cur ^= remember
				remember = nextRemember
				cur -= dataLength
				cur = rollRight(cur, 4)
				data[i] = cur
				dataLength--
			}
		}
	}
}

// EncryptData 自定义加密（6轮混淆），在 AES-OFB 加密之前调用。
// 对应 Java 的 MapleCustomEncryption.encryptData()。
func EncryptData(data []byte) {
	for j := 0; j < 6; j++ {
		var remember byte
		dataLength := byte(len(data) & 0xFF)

		if j%2 == 0 {
			// 正向遍历
			for i := 0; i < len(data); i++ {
				cur := data[i]
				cur = rollLeft(cur, 3)
				cur += dataLength
				cur ^= remember
				remember = cur
				cur = rollRight(cur, int(dataLength)&0xFF)
				cur = ^cur
				cur += 0x48
				dataLength--
				data[i] = cur
			}
		} else {
			// 反向遍历
			for i := len(data) - 1; i >= 0; i-- {
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

// ============================================================================
// 8. 高层接口：完整的解密/加密流水线
// ============================================================================

// ClientCrypto 封装客户端的完整加解密上下文。
// 一个 TCP 连接对应一个 ClientCrypto 实例。
type ClientCrypto struct {
	SendCipher    *MapleAESOFB // 服务端加密→客户端（用 sendIv 初始化）
	ReceiveCipher *MapleAESOFB // 客户端加密→服务端（用 recvIv 初始化）
}

// NewClientCrypto 创建加解密上下文。
//
// 参数:
//
//	version - 版本号 (83)
//	sendIv  - 服务端发送方向的 IV
//	recvIv  - 服务端接收方向的 IV
func NewClientCrypto(version uint16, sendIv, recvIv []byte) (*ClientCrypto, error) {
	// 服务端发送方向: iv=sendIv, version=0xFFFF-83
	sendCipher, err := NewMapleAESOFB(sendIv, 0xFFFF-version)
	if err != nil {
		return nil, err
	}
	// 服务端接收方向: iv=recvIv, version=83
	recvCipher, err := NewMapleAESOFB(recvIv, version)
	if err != nil {
		return nil, err
	}
	return &ClientCrypto{
		SendCipher:    sendCipher,
		ReceiveCipher: recvCipher,
	}, nil
}

// Decrypt 解密客户端发来的一个完整包（4字节header + body）。
// 返回解密并反混淆后的明文 body。
func (c *ClientCrypto) Decrypt(packet []byte) ([]byte, error) {
	if len(packet) < 4 {
		return nil, ErrPacketTooShort
	}

	header := binary.BigEndian.Uint32(packet[:4])
	if !c.ReceiveCipher.IsValidHeader(header) {
		return nil, ErrInvalidHeader
	}

	bodyLen := DecodePacketLengthUint(header)
	if len(packet[4:]) < bodyLen {
		return nil, ErrPacketTooShort
	}

	body := make([]byte, bodyLen)
	copy(body, packet[4:4+bodyLen])

	// 第一层：AES-OFB 解密
	c.ReceiveCipher.Crypt(body)

	// 第二层：自定义反混淆
	DecryptData(body)

	return body, nil
}

// Encrypt 加密一个要发给客户端的包体，返回完整的加密包（4字节header + body）。
//
// 顺序必须与 Java GMSV83PacketProtocol.encode() 一致:
//
//	① 先算包头（用当前 IV）
//	② 自定义混淆
//	③ AES-OFB 加密（会更新 IV）
func (c *ClientCrypto) Encrypt(body []byte) []byte {
	bodyCopy := make([]byte, len(body))
	copy(bodyCopy, body)

	// ① 构造包头（必须在 Crypt 之前，因为 Crypt 会更新 IV）
	header := c.SendCipher.EncodePacketHeader(len(bodyCopy))

	// ② 自定义混淆
	EncryptData(bodyCopy)

	// ③ AES-OFB 加密（此调用会 updateIv）
	c.SendCipher.Crypt(bodyCopy)

	result := make([]byte, 0, 4+len(bodyCopy))
	result = append(result, header...)
	result = append(result, bodyCopy...)
	return result
}

// ============================================================================
// 9. 错误定义
// ============================================================================

var (
	ErrPacketTooShort = errors.New("packet too short")
	ErrInvalidHeader  = errors.New("invalid packet header")
)
