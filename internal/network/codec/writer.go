package codec

import (
	"encoding/binary"
)

// Writer 封包构造器，流式写入各类型数据
type Writer struct {
	buf []byte
}

// NewWriter 创建写入器
func NewWriter() *Writer {
	return &Writer{buf: make([]byte, 0, 128)}
}

// NewWriterWithOpcode 创建写入器并写入 opcode
func NewWriterWithOpcode(opcode uint16) *Writer {
	w := NewWriter()
	w.WriteShort(opcode)
	return w
}

// WriteByte 写入 1 字节
func (w *Writer) WriteByte(b byte) {
	w.buf = append(w.buf, b)
}

// WriteShort 写入 2 字节 (小端序)
func (w *Writer) WriteShort(v uint16) {
	w.buf = append(w.buf, byte(v&0xFF), byte((v>>8)&0xFF))
}

// WriteInt 写入 4 字节 (小端序)
func (w *Writer) WriteInt(v uint32) {
	w.buf = binary.LittleEndian.AppendUint32(w.buf, v)
}

// WriteLong 写入 8 字节 (小端序)
func (w *Writer) WriteLong(v uint64) {
	w.buf = binary.LittleEndian.AppendUint64(w.buf, v)
}

// WriteString 写入 MapleStory 风格字符串：2 字节长度 (小端序) + ASCII 内容
func (w *Writer) WriteString(s string) {
	w.WriteShort(uint16(len(s)))
	w.buf = append(w.buf, []byte(s)...)
}

// WritePaddedString 写入定长字符串（常用于角色名，固定 13 字节）
func (w *Writer) WritePaddedString(s string, n int) {
	b := make([]byte, n)
	copy(b, s)
	w.buf = append(w.buf, b...)
}

// WriteBytes 写入原始字节
func (w *Writer) WriteBytes(b []byte) {
	w.buf = append(w.buf, b...)
}

// WritePos 写入位置坐标 [x: 2B][y: 2B]
func (w *Writer) WritePos(x, y int16) {
	w.WriteShort(uint16(x))
	w.WriteShort(uint16(y))
}

// WriteBool 写入布尔值 (1 字节: 0 或 1)
func (w *Writer) WriteBool(v bool) {
	if v {
		w.WriteByte(1)
	} else {
		w.WriteByte(0)
	}
}

// Bytes 返回已写入的完整数据 (不含包头)
func (w *Writer) Bytes() []byte {
	return w.buf
}

// Packet 返回封包对象（opcode = 前 2 字节）
func (w *Writer) Packet() *Packet {
	if len(w.buf) < 2 {
		return &Packet{Opcode: 0, Data: w.buf}
	}
	opcode := binary.LittleEndian.Uint16(w.buf[:2])
	return &Packet{Opcode: opcode, Data: w.buf[2:]}
}

