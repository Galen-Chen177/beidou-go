package codec

import (
	"encoding/binary"
)

// Reader 封包读取器，从 []byte 中顺序读取各类型数据
type Reader struct {
	buf []byte
	pos int
}

// NewReader 从字节数组创建读取器
func NewReader(data []byte) *Reader {
	return &Reader{buf: data, pos: 0}
}

// ReadByte 读取 1 字节
func (r *Reader) ReadByte() byte {
	b := r.buf[r.pos]
	r.pos++
	return b
}

// ReadShort 读取 2 字节 (小端序)
func (r *Reader) ReadShort() uint16 {
	v := binary.LittleEndian.Uint16(r.buf[r.pos:])
	r.pos += 2
	return v
}

// ReadInt 读取 4 字节 (小端序)
func (r *Reader) ReadInt() uint32 {
	v := binary.LittleEndian.Uint32(r.buf[r.pos:])
	r.pos += 4
	return v
}

// ReadLong 读取 8 字节 (小端序)
func (r *Reader) ReadLong() uint64 {
	v := binary.LittleEndian.Uint64(r.buf[r.pos:])
	r.pos += 8
	return v
}

// ReadString 读取 MapleStory 风格字符串：2 字节长度 + 内容
func (r *Reader) ReadString() string {
	length := r.ReadShort()
	s := string(r.buf[r.pos : r.pos+int(length)])
	r.pos += int(length)
	return s
}

// ReadPaddedString 读取 MapleStory 风格定长字符串（常用于角色名，固定 13 字节，ASCII）
func (r *Reader) ReadPaddedString(n int) string {
	end := r.pos + n
	for i := r.pos; i < end; i++ {
		if r.buf[i] == 0 {
			s := string(r.buf[r.pos:i])
			r.pos = end
			return s
		}
	}
	s := string(r.buf[r.pos:end])
	r.pos = end
	return s
}

// ReadBytes 读取 n 字节
func (r *Reader) ReadBytes(n int) []byte {
	b := r.buf[r.pos : r.pos+n]
	r.pos += n
	return b
}

// ReadPos 读取位置坐标：[x: 2B][y: 2B]
func (r *Reader) ReadPos() (x, y int16) {
	x = int16(r.ReadShort())
	y = int16(r.ReadShort())
	return
}

// Skip 跳过 n 字节
func (r *Reader) Skip(n int) {
	r.pos += n
}

// Remaining 返回剩余可读字节数
func (r *Reader) Remaining() int {
	return len(r.buf) - r.pos
}
