package store

import (
	"errors"
	"time"

	"beidou-go/internal/model"

	"gorm.io/gorm"
)

// ErrAccountNotFound 账号不存在
var ErrAccountNotFound = errors.New("account not found")

// AccountStore 账号数据访问层
type AccountStore struct {
	db *gorm.DB
}

// NewAccountStore 创建 AccountStore
func NewAccountStore(db *gorm.DB) *AccountStore {
	return &AccountStore{db: db}
}

// FindByName 按账号名查找
func (s *AccountStore) FindByName(name string) (*model.Account, error) {
	var account model.Account
	result := s.db.Where("name = ?", name).First(&account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, result.Error
	}
	return &account, nil
}

// Create 创建新账号（用于自动注册）
func (s *AccountStore) Create(name, passwordHash string) (*model.Account, error) {
	now := time.Now()
	account := &model.Account{
		Name:           name,
		Password:       passwordHash,
		Gender:         0,    // 默认男
		Characterslots: 3,    // 默认 3 个角色槽位
		Tos:            true, // 自动注册视为已同意
		Language:       0,
		CreatedAt:      &now,
		Lastlogin:      now,
	}
	result := s.db.Create(account)
	if result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}

// UpdatePassword 更新账号密码 hash（用于旧格式迁移到 bcrypt）
func (s *AccountStore) UpdatePassword(id uint64, hash string) error {
	return s.db.Model(&model.Account{}).Where("id = ?", id).Update("password", hash).Error
}

// UpdateLastLogin 更新最后登录时间
func (s *AccountStore) UpdateLastLogin(id uint64, t time.Time) error {
	return s.db.Model(&model.Account{}).Where("id = ?", id).Update("lastlogin", t).Error
}
