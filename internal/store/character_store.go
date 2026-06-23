package store

import (
	"errors"

	"beidou-go/internal/model"

	"gorm.io/gorm"
)

// ErrCharacterNotFound 角色不存在
var ErrCharacterNotFound = errors.New("character not found")

// CharacterStore 角色数据访问层
type CharacterStore struct {
	db *gorm.DB
}

// NewCharacterStore 创建 CharacterStore
func NewCharacterStore(db *gorm.DB) *CharacterStore {
	return &CharacterStore{db: db}
}

// FindByAccountAndWorld 按账号ID+世界ID查询角色列表
func (s *CharacterStore) FindByAccountAndWorld(accountID uint, worldID int) ([]model.Character, error) {
	var chars []model.Character
	result := s.db.Where("accountid = ? AND world = ?", accountID, worldID).Find(&chars)
	if result.Error != nil {
		return nil, result.Error
	}
	return chars, nil
}

// Create 创建新角色
func (s *CharacterStore) Create(chr *model.Character) error {
	return s.db.Create(chr).Error
}

// FindByName 按角色名查找（检查名字是否已被占用）
func (s *CharacterStore) FindByName(name string) (*model.Character, error) {
	var char model.Character
	result := s.db.Where("name = ?", name).First(&char)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrCharacterNotFound
		}
		return nil, result.Error
	}
	return &char, nil
}

// FindByID 按角色ID查询单个角色
func (s *CharacterStore) FindByID(charID int32) (*model.Character, error) {
	var char model.Character
	result := s.db.First(&char, charID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrCharacterNotFound
		}
		return nil, result.Error
	}
	return &char, nil
}
