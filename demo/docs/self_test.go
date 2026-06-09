// self_test.go
// 自测：加密→解密 往返，验证加解密实现正确
// 运行: go run maple_crypto.go self_test.go

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	version := uint16(83)

	sendIv := GenerateSendIV()    // {82, 48, 120, 193}
	recvIv := GenerateReceiveIV() // {70, 114, 122, 96}

	fmt.Printf("sendIv = % X\n", sendIv)
	fmt.Printf("recvIv = % X\n", recvIv)
	fmt.Println()

	// 构造 47 字节测试数据（模拟客户端登陆包体）
	original := make([]byte, 47)
	idx := 0
	original[idx] = 0x01
	idx++
	original[idx] = 0x00
	idx++
	original[idx] = 0x05
	idx++
	original[idx] = 0x00
	idx++
	copy(original[idx:], "admin")
	idx += 5
	original[idx] = 0x06
	idx++
	original[idx] = 0x00
	idx++
	copy(original[idx:], "123456")
	idx += 6
	idx += 6 // padding 6 bytes
	original[idx] = 0xAB
	idx++
	original[idx] = 0xCD
	idx++
	original[idx] = 0xEF
	idx++
	original[idx] = 0x01

	fmt.Printf("原始数据 (%d 字节):\n", len(original))
	fmt.Println(hex.Dump(original))

	failed := false

	// ============================================================
	// 测试1: recvCipher 往返（recvIv + version=83）
	// 模拟: 客户端用 recvIv 加密 → 服务端用 recvIv 解密
	// ============================================================
	fmt.Println("=== 测试1: recvCipher 往返 ===")

	encCipher, _ := NewMapleAESOFB(recvIv, version)
	decCipher, _ := NewMapleAESOFB(recvIv, version)

	encrypted := make([]byte, len(original))
	copy(encrypted, original)

	// 加密顺序（与 Java 一致）: ①包头 → ②混淆 → ③AES
	header := encCipher.EncodePacketHeader(len(encrypted))
	EncryptData(encrypted)
	encCipher.Crypt(encrypted)

	fullPacket := make([]byte, 0, 4+len(encrypted))
	fullPacket = append(fullPacket, header...)
	fullPacket = append(fullPacket, encrypted...)

	// 解密
	hdrU32 := binary.BigEndian.Uint32(fullPacket[:4])
	if !decCipher.IsValidHeader(hdrU32) {
		fmt.Println("❌ 测试1 包头校验失败!")
		failed = true
	} else {
		fmt.Println("✓ 测试1 包头校验通过")

		bodyLen := DecodePacketLengthUint(hdrU32)
		decrypted := make([]byte, bodyLen)
		copy(decrypted, fullPacket[4:4+bodyLen])
		decCipher.Crypt(decrypted)
		DecryptData(decrypted)

		if bytes.Equal(original, decrypted) {
			fmt.Println("✅ 测试1 通过!")
		} else {
			fmt.Println("❌ 测试1 数据不匹配!")
			for i := 0; i < len(original); i++ {
				if original[i] != decrypted[i] {
					fmt.Printf("  [%d] 原始:%02X 解密:%02X\n", i, original[i], decrypted[i])
				}
			}
			failed = true
		}
	}
	fmt.Println()

	// ============================================================
	// 测试2: sendCipher 往返（sendIv + version=0xFFFF-83）
	// ============================================================
	fmt.Println("=== 测试2: sendCipher 往返 ===")

	encCipher2, _ := NewMapleAESOFB(sendIv, 0xFFFF-version)
	decCipher2, _ := NewMapleAESOFB(sendIv, 0xFFFF-version)

	encrypted2 := make([]byte, len(original))
	copy(encrypted2, original)

	header2 := encCipher2.EncodePacketHeader(len(encrypted2))
	EncryptData(encrypted2)
	encCipher2.Crypt(encrypted2)

	fullPacket2 := append(header2, encrypted2...)

	fmt.Println(hex.Dump(fullPacket2))

	hdrU32_2 := binary.BigEndian.Uint32(fullPacket2[:4])
	if !decCipher2.IsValidHeader(hdrU32_2) {
		fmt.Println("❌ 测试2 包头校验失败!")
		failed = true
	} else {
		fmt.Println("✓ 测试2 包头校验通过")

		bodyLen2 := DecodePacketLengthUint(hdrU32_2)
		decrypted2 := make([]byte, bodyLen2)
		copy(decrypted2, fullPacket2[4:4+bodyLen2])
		decCipher2.Crypt(decrypted2)
		DecryptData(decrypted2)

		if bytes.Equal(original, decrypted2) {
			fmt.Println("✅ 测试2 通过!")
		} else {
			fmt.Println("❌ 测试2 数据不匹配!")
			failed = true
		}
	}
	fmt.Println()

	// ============================================================
	// 测试3: ClientCrypto.Decrypt
	// 用 recvCipher 手动构造客户端加密包 → ClientCrypto.Decrypt 解密
	// ============================================================
	fmt.Println("=== 测试3: ClientCrypto.Decrypt 正确性 ===")

	crypto, _ := NewClientCrypto(version, sendIv, recvIv)

	recvEnc, _ := NewMapleAESOFB(recvIv, version)
	body3 := make([]byte, len(original))
	copy(body3, original)
	hdr3 := recvEnc.EncodePacketHeader(len(body3))
	EncryptData(body3)
	recvEnc.Crypt(body3)

	clientPkt := append(hdr3, body3...)

	decrypted3, err := crypto.Decrypt(clientPkt)
	if err != nil {
		fmt.Printf("❌ 测试3 Decrypt 失败: %v\n", err)
		failed = true
	} else if bytes.Equal(original, decrypted3) {
		fmt.Println("✅ 测试3 通过!")
	} else {
		fmt.Println("❌ 测试3 数据不匹配!")
		failed = true
	}
	fmt.Println(hex.Dump(decrypted3))

	if failed {
		fmt.Println("=== 有测试失败! ===")
		os.Exit(1)
	} else {
		fmt.Println("=== 全部测试通过! ===")
	}
}
