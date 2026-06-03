package opcode

// 登录服务器 接收 (Client → Server)
const (
	LoginCheckPassword   uint16 = 0x01 // 密码验证请求
	LoginGuestLogin      uint16 = 0x02 // 游客登录
	LoginServerListRereq uint16 = 0x04 // 重新请求服务器列表
	LoginCharListReq     uint16 = 0x05 // 请求角色列表
	LoginCharSelect      uint16 = 0x13 // 选择角色（进入游戏）
	LoginCharCreate      uint16 = 0x16 // 创建角色
	LoginCharDelete      uint16 = 0x17 // 删除角色
	LoginCheckDuplicated uint16 = 0x19 // 检查角色名重复
	LoginCreateSecurity  uint16 = 0x20 // 设置 PIC/PIN
	LoginAuthHttp        uint16 = 0x1A // HTTP 认证
)

// 登录服务器 发送 (Server → Client)
const (
	LoginStatus          uint16 = 0x00 // 登录状态返回
	LoginServerList      uint16 = 0x0A // 服务器列表
	LoginCharList        uint16 = 0x08 // 角色列表
	LoginCharInfo        uint16 = 0x09 // 角色信息
	LoginCharNameResp    uint16 = 0x0D // 角色名检查结果
	LoginSecondPassword  uint16 = 0x04 // 二次密码请求（PIC）
	LoginAuthHttpResp    uint16 = 0x00 // HTTP 认证响应（暂未确定 opcode）
)

// 登录状态码
const (
	LoginOK             byte = 0x00
	LoginBanned         byte = 0x03
	LoginWrongPassword  byte = 0x04
	LoginNotRegistered  byte = 0x05
	LoginDBFail         byte = 0x07
	LoginAlreadyLogin   byte = 0x07
	LoginTooManyConn    byte = 0x0A
	LoginTempBlocked    byte = 0x06
	LoginMasterIP       byte = 0x0B
	LoginWrongGateway   byte = 0x0C
)
