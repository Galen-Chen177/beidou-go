// example_usage.go
// 演示用 maple_crypto.go 完成三步操作：
//   Step 1: 生成 IV
//   Step 2: 打印 Hello 包
//   Step 3: 加解密客户端数据
//
// 可直接 go run example_usage.go

package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	version := uint16(83) // 冒险岛 v83

	// ========================================
	// Step 1: 生成随机 IV
	// ========================================
	sendIv := GenerateSendIV()    // 例如 {0x52, 0x30, 0x78, 0xC1}  = "R0x?"
	recvIv := GenerateReceiveIV() // 例如 {0x46, 0x72, 0x7A, 0x60} = "Frz?"

	fmt.Printf("sendIv = % X  (%q)\n", sendIv, sendIv)
	fmt.Printf("recvIv = % X  (%q)\n", recvIv, recvIv)
	fmt.Println()

	// ========================================
	// Step 2: 打印 Hello 包（这是明文，直接发 TCP）
	// ========================================
	helloPkt := CreateHelloPacket(version, sendIv, recvIv)
	fmt.Printf("Hello 包 (%d 字节):\n", len(helloPkt))
	fmt.Println(hex.Dump(helloPkt))
	// 输出示例:
	// 0e 00 53 00 01 00 31 46 72 7a 60 52 30 78 c1 08
	// └──┘ └──┘ └──┘ └┘ └─────────┘ └─────────┘ └┘
	// opcode ver flag '1' recvIv       sendIv       term

	// ========================================
	// Step 3: 创建加解密器
	// ========================================
	crypto, err := NewClientCrypto(version, sendIv, recvIv)
	if err != nil {
		panic(err)
	}

	// ========================================
	// Step 3a: 解密客户端发来的包
	// ========================================
	// 假设从 TCP 收到了客户端的加密数据（你之前给的 hex）
	encryptedHex := "971ab81a58ba5bf7adc8b3bac28ee67fbcfa6601d1267e0d7ddbbf70e5746bb459a258db002402ef6df9431fe4b718be318b13"
	encryptedBytes, _ := hex.DecodeString(encryptedHex)

	plainBody, err := crypto.Decrypt(encryptedBytes)
	if err != nil {
		fmt.Printf("解密失败: %v\n", err)
		// 如果 IV 和加密时不匹配，会返回 ErrInvalidHeader
	} else {
		fmt.Printf("解密后包体 (%d 字节): % X\n", len(plainBody), plainBody)
		// 现在 plainBody 就是 LoginPasswordHandler 能解析的格式:
		//   [2B opcode] [string login] [string password] [6B padding] [4B hwid]
	}

	// ========================================
	// Step 3b: 加密要发给客户端的包
	// ========================================
	plainPayload := []byte{
		0x00, 0x00, // opcode (示例: LOGIN_STATUS)
		0x00,                   // status=0 (成功)
		0x00, 0x00, 0x00, 0x00, // accountId
		0x00, // gender
		// ... 其他字段
	}
	encryptedPkt := crypto.Encrypt(plainPayload)
	fmt.Printf("加密后包 (%d 字节): % X\n", len(encryptedPkt), encryptedPkt)
	// 可以直接 Write 到 TCP socket
	_ = encryptedPkt
}
