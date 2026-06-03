package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

// MapleCrypto 冒险岛 AES/OFB 加解密器
//
// GMS v0.83 使用 AES/OFB 模式加密网络封包。
// 握手阶段客户端发送 IV，服务端用此 IV 初始化加密/解密流。
type MapleCrypto struct {
	sendIV  []byte
	recvIV  []byte
	sendCipher cipher.Stream
	recvCipher cipher.Stream
}

// NewMapleCrypto 创建加解密器
// sendIV: 用于加密（发送给客户端的数据）
// recvIV: 用于解密（从客户端接收的数据）
// mapleVersion: 冒险岛版本号，用于 IV 变换（v0.83 = 83）
func NewMapleCrypto(sendIV, recvIV []byte, mapleVersion uint16) (*MapleCrypto, error) {
	m := &MapleCrypto{
		sendIV: make([]byte, len(sendIV)),
		recvIV: make([]byte, len(recvIV)),
	}
	copy(m.sendIV, sendIV)
	copy(m.recvIV, recvIV)

	// 对 IV 进行 MapleStory 版本变换
	m.transformIV(m.sendIV, mapleVersion)
	m.transformIV(m.recvIV, mapleVersion)

	// 创建 AES cipher block
	sendBlock, err := aes.NewCipher(m.sendIV)
	if err != nil {
		return nil, err
	}
	recvBlock, err := aes.NewCipher(m.recvIV)
	if err != nil {
		return nil, err
	}

	// AES/OFB 模式：IV 即密钥本身
	m.sendCipher = cipher.NewOFB(sendBlock, m.sendIV)
	m.recvCipher = cipher.NewOFB(recvBlock, m.recvIV)

	return m, nil
}

// transformIV 对 IV 进行版本变换
// 这是 MapleStory 自定义的 IV 变换算法
func (m *MapleCrypto) transformIV(iv []byte, version uint16) {
	versionBytes := []byte{
		byte((version >> 8) & 0xFF),
		byte(version & 0xFF),
	}
	// 对前 4 字节进行变换
	for i := 0; i < 4; i++ {
		iv[i] ^= versionBytes[i%2]
	}
}

// Encrypt 加密数据（就地修改）
// 用于发送给客户端的数据
func (m *MapleCrypto) Encrypt(data []byte) {
	m.sendCipher.XORKeyStream(data, data)
}

// Decrypt 解密数据（就地修改）
// 用于从客户端接收的数据
func (m *MapleCrypto) Decrypt(data []byte) {
	m.recvCipher.XORKeyStream(data, data)
}

// GetSendIV 返回发送端 IV（用于客户端握手协商）
func (m *MapleCrypto) GetSendIV() []byte {
	iv := make([]byte, len(m.sendIV))
	copy(iv, m.sendIV)
	return iv
}

// GetRecvIV 返回接收端 IV
func (m *MapleCrypto) GetRecvIV() []byte {
	iv := make([]byte, len(m.recvIV))
	copy(iv, m.recvIV)
	return iv
}
