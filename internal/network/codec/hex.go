package codec

import (
	"fmt"
	"strings"
)

// maxHexDumpLen 单次 hex dump 最大字节数，超过则截断避免日志爆炸
const maxHexDumpLen = 1024

// HexDump 按 Wireshark 风格格式化字节数组为 hex dump 字符串
//
// 输出示例:
//
//	00000000  01 00 61 64 6D 69 6E 00  00 00 00 00 00 00 00 00  |..admin.........|
//	00000010  31 32 33 34 35 36                                 |123456          |
//
// 用于和 Wireshark 抓包对照：Wireshark TCP payload 的 hex 就是服务端收发的加密字节。
func HexDump(data []byte) string {
	length := len(data)
	truncated := false
	if length > maxHexDumpLen {
		data = data[:maxHexDumpLen]
		truncated = true
	}

	var sb strings.Builder
	for i := 0; i < len(data); i += 16 {
		// 偏移地址
		sb.WriteString(fmt.Sprintf("  %08X  ", i))

		// Hex 部分（中间多一个空格分隔）
		for j := 0; j < 16; j++ {
			if j == 8 {
				sb.WriteString(" ")
			}
			if i+j < len(data) {
				sb.WriteString(fmt.Sprintf("%02X ", data[i+j]))
			} else {
				sb.WriteString("   ")
			}
		}

		// ASCII 可读部分
		sb.WriteString(" |")
		for j := 0; j < 16; j++ {
			if i+j < len(data) {
				b := data[i+j]
				if b >= 32 && b < 127 {
					sb.WriteByte(b)
				} else {
					sb.WriteByte('.')
				}
			} else {
				sb.WriteByte(' ')
			}
		}
		sb.WriteString("|\n")
	}

	if truncated {
		sb.WriteString(fmt.Sprintf("  ... (截断, 共 %d 字节)\n", length))
	}

	return sb.String()
}
