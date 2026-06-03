package codec

// Packet 封包结构体
// 格式: [header: 4 bytes (小端序, = 包体长度)] [body: opcode(2B) + data(NB)]
type Packet struct {
	Opcode uint16
	Data   []byte
}

// NewPacket 创建一个新封包
func NewPacket(opcode uint16) *Packet {
	return &Packet{
		Opcode: opcode,
		Data:   make([]byte, 0, 64),
	}
}

// Len 返回封包体长度 (opcode 2 字节 + data 字节数)
func (p *Packet) Len() int {
	return 2 + len(p.Data)
}

// Bytes 返回完整封包体 [opcode 2B][data]
func (p *Packet) Bytes() []byte {
	buf := make([]byte, 2+len(p.Data))
	buf[0] = byte(p.Opcode & 0xFF)
	buf[1] = byte((p.Opcode >> 8) & 0xFF)
	copy(buf[2:], p.Data)
	return buf
}
