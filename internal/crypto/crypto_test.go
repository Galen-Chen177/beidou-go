// crypto_test.go
// 用真实客户端抓包数据验证加解密正确性
//
// 测试数据来自 Wireshark 抓包（见 demo/docs/xx.md）:
//   SERVER_HELLO → CLIENT_HELLO → LOGIN 三包序列
//
// 运行: go test ./internal/crypto/ -v -run TestDecryptCapturedPackets

package crypto

import (
	"encoding/hex"
	"testing"
)

func TestDecryptCapturedPackets(t *testing.T) {
	version := uint16(83)

	// ---- 真实抓包数据 ----
	// CLIENT_HELLO（6 字节 = 4B header + 2B body，客户端发来的第一个加密包）
	clientHello := []byte{0x29, 0x60, 0x2b, 0x60, 0x03, 0xa2}

	// LOGIN（51 字节 = 4B header + 47B body，登录包）
	loginPkt := []byte{
		0x97, 0x1a, 0xb8, 0x1a, 0x58, 0xba, 0x5b, 0xf7,
		0xad, 0xc8, 0xb3, 0xba, 0xc2, 0x8e, 0xe6, 0x7f,
		0xbc, 0xfa, 0x66, 0x01, 0xd1, 0x26, 0x7e, 0x0d,
		0x7d, 0xdb, 0xbf, 0x70, 0xe5, 0x74, 0x6b, 0xb4,
		0x59, 0xa2, 0x58, 0xdb, 0x00, 0x24, 0x02, 0xef,
		0x6d, 0xf9, 0x43, 0x1f, 0xe4, 0xb7, 0x18, 0xbe,
		0x31, 0x8b, 0x13,
	}

	// 1. 模拟握手
	myCrypto := NewCrypto(version)
	serverHello := myCrypto.GenServerHello()
	if len(serverHello) == 0 {
		t.Fatal("GenServerHello 返回空包")
	}
	t.Logf("SERVER_HELLO (%d bytes):\n%s", len(serverHello), hex.Dump(serverHello))

	// 2. 解密 CLIENT_HELLO（第一个加密包，解密后 IV 与客户端同步）
	decryptedHello, err := myCrypto.Decrypt(clientHello)
	if err != nil {
		t.Fatalf("CLIENT_HELLO 解密失败: %v", err)
	}
	t.Logf("CLIENT_HELLO 解密后 (%d bytes):\n%s", len(decryptedHello), hex.Dump(decryptedHello))

	// 3. 解密 LOGIN
	decryptedLogin, err := myCrypto.Decrypt(loginPkt)
	if err != nil {
		t.Fatalf("LOGIN 解密失败: %v", err)
	}
	t.Logf("LOGIN 解密后 (%d bytes):\n%s", len(decryptedLogin), hex.Dump(decryptedLogin))

	// 4. 验证 LOGIN 包体格式
	if len(decryptedLogin) < 2 {
		t.Fatal("LOGIN 解密后太短")
	}

	opcode := uint16(decryptedLogin[0]) | uint16(decryptedLogin[1])<<8
	if opcode != 0x01 {
		t.Errorf("期望 opcode=0x01 (LoginCheckPassword)，实际=0x%04X", opcode)
	}
	t.Logf("opcode = 0x%04X (LoginCheckPassword) ✓", opcode)

	data := decryptedLogin[2:]
	pos := 0

	// 读账号名
	if pos+2 > len(data) {
		t.Fatal("LOGIN 包体太短，无法读账号名长度")
	}
	nameLen := int(data[pos]) | int(data[pos+1])<<8
	pos += 2
	if pos+nameLen > len(data) {
		t.Fatalf("账号名长度 %d 超出包体", nameLen)
	}
	accountName := string(data[pos : pos+nameLen])
	pos += nameLen
	t.Logf("账号名: %q", accountName)

	// 读密码
	if pos+2 > len(data) {
		t.Fatal("LOGIN 包体太短，无法读密码长度")
	}
	pwdLen := int(data[pos]) | int(data[pos+1])<<8
	pos += 2
	if pos+pwdLen > len(data) {
		t.Fatalf("密码长度 %d 超出包体", pwdLen)
	}
	password := string(data[pos : pos+pwdLen])
	t.Logf("密码: %q", password)

	t.Logf("✓ 解密成功: 账号=%q 密码=%q", accountName, password)
}
