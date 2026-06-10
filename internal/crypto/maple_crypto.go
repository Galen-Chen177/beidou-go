// Package crypto 实现冒险岛 GMS v83 网络封包加解密。
//
// 对应 Java 服务端文件:
//   - MapleAESOFB.java            AES-OFB 流密码核心
//   - MapleCustomEncryption.java  自定义 6 轮混淆/反混淆
//   - InitializationVector.java   IV 生成
//   - PacketCreator.getHello()    Hello 包构造
//
// 包体保护分两层: AES-OFB 流密码 + 自定义混淆。加密和解密都须严格按序操作。
package crypto

import (
	"crypto/aes"
	"crypto/rand"
	"errors"
	"math/big"
)

// =============================================================================
// 全局常量 & 错误定义
// =============================================================================

// aesKey 是 AES-256 的 32 字节硬编码密钥，与 Java 服务端完全一致。
var aesKey = []byte{
	0x13, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00,
	0x06, 0x00, 0x00, 0x00, 0xB4, 0x00, 0x00, 0x00,
	0x1B, 0x00, 0x00, 0x00, 0x0F, 0x00, 0x00, 0x00,
	0x33, 0x00, 0x00, 0x00, 0x52, 0x00, 0x00, 0x00,
}

var (
	ErrPacketTooShort = errors.New("packet too short")
	ErrInvalidHeader  = errors.New("invalid packet header")
)

// =============================================================================
// IV 生成
// =============================================================================

func randomByte() byte {
	n, err := rand.Int(rand.Reader, big.NewInt(255))
	if err != nil {
		return 0x60
	}
	return byte(n.Int64())
}

// genSendIV 生成服务端发送方向的 IV（服务端加密 → 客户端解密）。
// Java: {82, 48, 120, randomByte} → ASCII "R0x" + 随机字节
func genSendIV() []byte {
	return []byte{82, 48, 120, 193}
}

// genRecvIV 生成服务端接收方向的 IV（客户端加密 → 服务端解密）。
// Java: {70, 114, 122, randomByte} → ASCII "Frz" + 随机字节
func genRecvIV() []byte {
	return []byte{70, 114, 122, 96}
}

// =============================================================================
// Hello 包（唯一明文发送的包，用于交换 IV）
// =============================================================================

// createHelloPacket 构造服务端发给客户端的握手包。
//
// 包格式（16 字节）:
//
//	Offset  Size  Value
//	 0       2    opcode = 0x000E (LE)
//	 2       2    version (LE)
//	 4       2    flag = 0x0001
//	 6       1    subversion = 0x31 ('1')
//	 7       4    recvIv — 客户端用它来加密发给服务端的数据
//	11       4    sendIv — 客户端用它来解密服务端发来的数据
//	15       1    terminator = 0x08
//
// IV 的方向对应:
//
//	服务端 sendIv ───加密──▶ 客户端用 sendIv 解密
//	服务端 recvIv ◀──解密─── 客户端用 recvIv 加密
func createHelloPacket(version uint16, sendIv, recvIv []byte) []byte {
	p := make([]byte, 0, 16)
	p = append(p, 0x0E, 0x00)                           // opcode (short LE)
	p = append(p, byte(version&0xFF), byte(version>>8)) // version (short LE)
	p = append(p, 0x01, 0x00)                           // flag (short LE)
	p = append(p, 0x31)                                 // subversion
	p = append(p, recvIv...)                            // recvIv (4B)
	p = append(p, sendIv...)                            // sendIv (4B)
	p = append(p, 0x08)                                 // terminator
	return p
}

// =============================================================================
// S-Box 查找表（256 字节，用于 IV 轮转）
// =============================================================================

// funnyBytes 是 Java 服务端的静态 S-Box，用于 getNewIv 中的非线性变换。
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

// =============================================================================
// MapleAESOFB — AES-OFB 流密码核心
// =============================================================================

// aesBlock 是 crypto/aes.Cipher 的接口子集，只用 Encrypt 方向。
type aesBlock interface {
	Encrypt(dst, src []byte)
	BlockSize() int
}

// MapleAESOFB 是冒险岛的 AES-OFB 流密码。
//
// 它把 4 字节 IV 重复 4 次得到 16 字节块，反复喂给 AES 生成密钥流，然后与数据 XOR。
// 加密和解密是同一操作（XOR 对称性），都调用 Crypt()。
// 每次 Crypt 执行完毕后会自动轮转 IV，保证同内容不同包的密文不同。
type MapleAESOFB struct {
	iv           []byte   // 4 字节当前 IV（每次 Crypt 后轮转）
	mapleVersion uint16   // 字节交换后的版本号，用于包头校验
	block        aesBlock // AES-256 分组密码实例
}

// NewMapleAESOFB 创建加解密器。
//
//	iv           — 4 字节初始向量
//	  - 服务端解密用 recvIv（客户端加密 → 服务端解密）
//	  - 服务端加密用 sendIv（服务端加密 → 客户端解密）
//	mapleVersion — 客户端版本号，如 83
func NewMapleAESOFB(iv []byte, mapleVersion uint16) (*MapleAESOFB, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	// 字节交换: 83(0x0053) → 0x5300
	// Java: (short)(((mapleVersion >> 8) & 0xFF) | ((mapleVersion << 8) & 0xFF00))
	swapped := ((mapleVersion >> 8) & 0xFF) | ((mapleVersion << 8) & 0xFF00)

	ivCopy := make([]byte, 4)
	copy(ivCopy, iv)

	return &MapleAESOFB{
		iv:           ivCopy,
		mapleVersion: swapped,
		block:        block,
	}, nil
}

// IV 返回当前 IV 的副本（只读）。
func (m *MapleAESOFB) IV() []byte {
	cp := make([]byte, 4)
	copy(cp, m.iv)
	return cp
}

// Crypt 对 data 进行原地加/解密（XOR 流密码）。
//
// 算法:
//  1. IV[4] 重复 4 遍 → 16 字节初始块
//  2. AES 加密该块 → 得到 16 字节密钥流
//  3. 密钥流 XOR 数据
//  4. 将 AES 输出再次加密 → 下一段 16 字节密钥流（OFB 链式）
//  5. 数据全部处理完后，调用 updateIv() 轮转 IV
//
// 分块大小: 首块 0x5B0 (1456)，后续块 0x5B4 (1460)
func (m *MapleAESOFB) Crypt(data []byte) {
	remaining := len(data)
	chunkSize := 0x5B0 // 1456
	start := 0

	for remaining > 0 {
		myIv := multiplyBytes(m.iv, 4, 4) // IV 重复 4 遍 → 16B

		if remaining < chunkSize {
			chunkSize = remaining
		}

		for x := start; x < start+chunkSize; x++ {
			if (x-start)%len(myIv) == 0 {
				m.block.Encrypt(myIv, myIv) // AES → 新密钥流
			}
			data[x] ^= myIv[(x-start)%len(myIv)]
		}

		start += chunkSize
		remaining -= chunkSize
		chunkSize = 0x5B4 // 1460
	}

	m.updateIv()
}

// =============================================================================
// 包头编解码
//
// 每个加密包前有 4 字节包头，编码了版本校验信息和包体长度。
// 包头必须在 Crypt 之前生成/校验，因为 Crypt 会轮转 IV。
// =============================================================================

// EncodePacketHeader 构造 4 字节包头（加密前调用，在 Crypt 之前）。
//
// 包头格式（大端序输出）:
//
//	[0] = iiv_hi     [1] = iiv_lo
//	[2] = xored_hi   [3] = xored_lo
//
// 其中 iiv = (iv[2],iv[3] 组合) XOR mapleVersion
// xoredIv = iiv XOR mlength（mlength 是长度的字节交换）
func (m *MapleAESOFB) EncodePacketHeader(bodyLength int) []byte {
	// iiv = (iv[3] | iv[2]<<8) ^ mapleVersion
	iiv := (uint16(m.iv[3]) & 0xFF) | (uint16(m.iv[2])<<8)&0xFF00
	iiv ^= m.mapleVersion

	// mlength = 字节交换后的长度
	mlength := ((uint16(bodyLength) << 8) & 0xFF00) | (uint16(bodyLength) >> 8)

	xoredIv := iiv ^ mlength

	ret := make([]byte, 4)
	ret[0] = byte((iiv >> 8) & 0xFF)
	ret[1] = byte(iiv & 0xFF)
	ret[2] = byte((xoredIv >> 8) & 0xFF)
	ret[3] = byte(xoredIv & 0xFF)
	return ret
}

// IsValidHeader 校验收到的 4 字节包头是否合法（解密前调用，在 Crypt 之前）。
//
// 合法性条件:
//
//	header[0] XOR iv[2] == mapleVersion_hi
//	header[1] XOR iv[3] == mapleVersion_lo
//
// 校验失败说明 IV 不同步（如握手阶段处理不正确）。
func (m *MapleAESOFB) IsValidHeader(header uint32) bool {
	b0 := byte((header >> 24) & 0xFF)
	b1 := byte((header >> 16) & 0xFF)
	return ((int(b0^m.iv[2]) & 0xFF) == int((m.mapleVersion>>8)&0xFF)) &&
		((int(b1^m.iv[3]) & 0xFF) == int(m.mapleVersion&0xFF))
}

// DecodePacketLength 从 4 字节包头解码包体长度（[]byte 版本）。
func DecodePacketLength(header []byte) int {
	return (int((header[1]^header[3])&0xFF) << 8) | int((header[0]^header[2])&0xFF)
}

// DecodePacketLengthUint 从 uint32 包头解码包体长度。
func DecodePacketLengthUint(header uint32) int {
	length := (header >> 16) ^ (header & 0xFFFF)
	length = ((length << 8) & 0xFF00) | ((length >> 8) & 0x00FF)
	return int(length)
}

// =============================================================================
// IV 轮转
// =============================================================================

// updateIv 在每次 Crypt 结束后轮转 IV，保证连续包的密钥流不同。
func (m *MapleAESOFB) updateIv() {
	m.iv = getNewIv(m.iv)
}

// getNewIv 用固定种子 {0xf2,0x53,0x50,0xc6} 和当前 IV 进行 4 轮混合，
// 产生下一轮的 IV。
func getNewIv(oldIv []byte) []byte {
	in := []byte{0xF2, 0x53, 0x50, 0xC6}
	for x := range 4 {
		funnyShit(oldIv[x], in)
	}
	return in
}

// funnyShit 是冒险岛自创的字节混合函数，用 S-Box 对 in 做非线性变换。
// 变量名沿用 Java 原始代码，保持可对照性。
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

// =============================================================================
// 自定义混淆/反混淆（MapleCustomEncryption — 6 轮字节变换）
//
// AES-OFB 解密后还要经过这一层才是真正的明文。
// 加密时先混淆再 AES，解密时先 AES 再反混淆。
// =============================================================================

// rollLeft 循环左移。
func rollLeft(in byte, count int) byte {
	tmp := uint16(in) & 0xFF
	tmp = tmp << (count % 8)
	return byte((tmp & 0xFF) | (tmp >> 8))
}

// rollRight 循环右移。
func rollRight(in byte, count int) byte {
	tmp := uint16(in) & 0xFF
	tmp = (tmp << 8) >> (count % 8)
	return byte((tmp & 0xFF) | (tmp >> 8))
}

// DecryptData 自定义解密（6 轮逆混淆），在 AES-OFB 解密之后调用。
//
//	j 为奇数 → 反向遍历；j 为偶数 → 正向遍历
func DecryptData(data []byte) {
	for j := 1; j <= 6; j++ {
		var remember byte
		dataLength := byte(len(data) & 0xFF)
		var nextRemember byte

		if j%2 == 0 {
			for i := range data {
				cur := data[i]
				cur -= 0x48
				cur = ^cur
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

// EncryptData 自定义加密（6 轮混淆），在 AES-OFB 加密之前调用。
//
//	与 DecryptData 互为逆操作。
func EncryptData(data []byte) {
	for j := range 6 {
		var remember byte
		dataLength := byte(len(data) & 0xFF)

		if j%2 == 0 {
			for i := range data {
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

// =============================================================================
// 工具函数
// =============================================================================

// multiplyBytes 把 in 的前 count 字节重复 mul 次。
// 例: multiplyBytes({a,b,c,d}, 4, 4) → {a,b,c,d, a,b,c,d, a,b,c,d, a,b,c,d}
func multiplyBytes(in []byte, count, mul int) []byte {
	out := make([]byte, count*mul)
	for i := range out {
		out[i] = in[i%count]
	}
	return out
}
