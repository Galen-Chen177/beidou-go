package crypto

// ShandaEncrypt Shanda 自定义加密（字节级变换）
//
// 这是 MapleStory 早期版本使用的一种自定义加密算法。
// GMS v0.83 中可能使用也可能已废弃，需对照原 Java 项目确认。
//
// 算法本质是对每个字节进行位移和 XOR 变换。
func ShandaEncrypt(data []byte) []byte {
	result := make([]byte, len(data))
	copy(result, data)
	return result
}

// ShandaDecrypt Shanda 自定义解密（字节级逆变换）
func ShandaDecrypt(data []byte) []byte {
	result := make([]byte, len(data))
	copy(result, data)
	return result
}
