// Package crypto 实现冒险岛 GMS v83 网络封包加解密。
package crypto

import (
	"encoding/binary"

	"github.com/sirupsen/logrus"
)

// Crypto 定义加解密器的公共接口。
//
// 每个 TCP 连接拥有独立的 Crypto 实例（对应一组 IV），
// 握手阶段由 GenServerHello 完成 IV 初始化，后续通信通过 Decrypt/Encrypt 处理封包。
type Crypto interface {
	// GenServerHello 生成服务端握手包（SERVER_HELLO）。
	// 调用后会初始化内部的 send/receive cipher，后续可正常加解密。
	GenServerHello() []byte

	// Decrypt 解密客户端发来的加密包。
	// 输入: 完整加密包（4 字节 header + N 字节 body）
	// 输出: 解密+反混淆后的明文 body
	Decrypt([]byte) ([]byte, error)

	// Encrypt 加密要发给客户端的包体。
	// 输入: 明文包体（不含 header）
	// 输出: header + 加密混淆后的包体，调用者必须使用返回值发送
	Encrypt([]byte) ([]byte, error)
}

// NewCrypto 创建加解密器。
// version 为客户端版本号（如 83）。
func NewCrypto(version uint16) Crypto {
	return &MyCrypto{
		Version: version,
	}
}

// MyCrypto 是 Crypto 接口的具体实现。
type MyCrypto struct {
	Version       uint16
	SendCipher    *MapleAESOFB // 加密发给客户端的数据（iv=sendIv, version=0xFFFF-83）
	ReceiveCipher *MapleAESOFB // 解密客户端发来的数据（iv=recvIv, version=83）
}

// GenServerHello 生成 IV、初始化 cipher，返回 SERVER_HELLO 包。
//
// 这是握手第一步：服务端生成一对 IV，将其嵌入 Hello 包明文发给客户端。
// 客户端收到后，用 recvIv 加密上行数据，用 sendIv 解密下行数据。
//
// 调用时机：TCP 连接建立后立即调用，调用一次即可。
func (s *MyCrypto) GenServerHello() []byte {
	sendIV := genSendIV()
	recvIV := genRecvIV()

	// 发送方向: iv=sendIv, version=0xFFFF-83（与 Java 一致）
	sendCipher, err := NewMapleAESOFB(sendIV, 0xFFFF-s.Version)
	if err != nil {
		logrus.Errorf("[Crypto] 创建 SendCipher 失败: %v", err)
		return nil
	}

	// 接收方向: iv=recvIv, version=83
	recvCipher, err := NewMapleAESOFB(recvIV, s.Version)
	if err != nil {
		logrus.Errorf("[Crypto] 创建 ReceiveCipher 失败: %v", err)
		return nil
	}

	s.SendCipher = sendCipher
	s.ReceiveCipher = recvCipher
	return createHelloPacket(s.Version, sendIV, recvIV)
}

// Decrypt 解密客户端发来的完整加密包。
//
// 解密步骤:
//  1. 读 4 字节包头（BigEndian uint32）
//  2. 校验包头合法性（IsValidHeader，必须在 Crypt 之前）
//  3. 解码包体长度
//  4. AES-OFB 解密（Crypt，会轮转 recvIV）
//  5. 自定义反混淆（DecryptData）
func (s *MyCrypto) Decrypt(packet []byte) ([]byte, error) {
	if len(packet) < 4 {
		return nil, ErrPacketTooShort
	}

	header := binary.BigEndian.Uint32(packet[:4])
	if !s.ReceiveCipher.IsValidHeader(header) {
		return nil, ErrInvalidHeader
	}

	bodyLen := DecodePacketLengthUint(header)
	if len(packet[4:]) < bodyLen {
		return nil, ErrPacketTooShort
	}

	body := make([]byte, bodyLen)
	copy(body, packet[4:4+bodyLen])

	// 第一层：AES-OFB 解密（会调用 updateIv 轮转 recvIV）
	s.ReceiveCipher.Crypt(body)
	// 第二层：自定义反混淆
	DecryptData(body)

	return body, nil
}

// Encrypt 加密要发给客户端的包体。
//
// 加密步骤（顺序必须与 Java GMSV83PacketProtocol.encode() 一致）:
//  1. 构造 4 字节包头（基于当前 sendIV，必须在 Crypt 之前）
//  2. 自定义混淆（EncryptData）
//  3. AES-OFB 加密（Crypt，会调用 updateIv 轮转 sendIV）
//
// 返回值: header(4B) + 加密混淆后的 body。调用者必须使用此返回值发送，原始 body 不会被修改。
func (s *MyCrypto) Encrypt(body []byte) ([]byte, error) {
	bodyCopy := make([]byte, len(body))
	copy(bodyCopy, body)

	// ① 构造包头（必须在 Crypt 之前，Crypt 会轮转 IV）
	header := s.SendCipher.EncodePacketHeader(len(bodyCopy))

	// ② 自定义混淆
	EncryptData(bodyCopy)

	// ③ AES-OFB 加密（会 updateIv）
	s.SendCipher.Crypt(bodyCopy)

	result := make([]byte, 0, 4+len(bodyCopy))
	result = append(result, header...)
	result = append(result, bodyCopy...)
	return result, nil
}
