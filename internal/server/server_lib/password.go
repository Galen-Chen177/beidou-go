package server_lib

import (
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// BcryptCost bcrypt 工作因子
const BcryptCost = 10

// VerifyPassword 校验密码明文与数据库中存储的 hash
//
// 支持三种 hash 格式:
//   - bcrypt: $2a$ 或 $2b$ 前缀
//   - SHA-1:   40 字符 hex 字符串
//   - SHA-512: 128 字符 hex 字符串
func VerifyPassword(plain, hash string) bool {
	if len(hash) == 0 {
		return false
	}

	switch {
	case strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$"):
		// bcrypt
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
		return err == nil
	case len(hash) == 40:
		// SHA-1
		return sha1Hex(plain) == strings.ToLower(hash)
	case len(hash) == 128:
		// SHA-512
		return sha512Hex(plain) == strings.ToLower(hash)
	default:
		// 未知格式，拒绝
		return false
	}
}

// HashPassword 使用 bcrypt 加密密码
func HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// NeedsRehash 判断密码 hash 是否需要升级为 bcrypt
// 非 bcrypt 格式（SHA-1、SHA-512）返回 true
func NeedsRehash(hash string) bool {
	return !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$")
}

// sha1Hex 计算 SHA-1 hex 字符串
func sha1Hex(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// sha512Hex 计算 SHA-512 hex 字符串
func sha512Hex(s string) string {
	h := sha512.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
