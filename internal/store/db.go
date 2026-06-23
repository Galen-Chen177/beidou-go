package store

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"beidou-go/config"
	"beidou-go/internal/model"
)

var db *gorm.DB

// DB 返回全局数据库实例
func DB() *gorm.DB {
	return db
}

// InitDB 初始化数据库连接并自动迁移
func InitDB(cfg config.DatabaseConfig) error {
	var err error
	db, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // 生产模式只记录警告
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层 sql.DB 失败: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	return nil
}

// AutoMigrate 自动迁移表结构
func AutoMigrate() error {
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return db.AutoMigrate(
		&model.Account{},
		&model.Alliance{},
		&model.Allianceguild{},
		&model.AreaInfo{},
		&model.AutobanConfig{},
		&model.BbsReply{},
		&model.BbsThread{},
		&model.BosslogDaily{},
		&model.BosslogWeekly{},
		&model.Buddy{},
		&model.Character{},
		&model.CommandInfo{},
		&model.Cooldown{},
		&model.DropData{},
		&model.DropDataGlobal{},
		&model.Dueyitem{},
		&model.Dueypackage{},
		&model.Eventstat{},
		&model.ExtendValue{},
		&model.Famelog{},
		&model.FamilyCharacter{},
		&model.FamilyEntitlement{},
		&model.FlywaySchemaHistory{},
		&model.Fredstorage{},
		&model.GachaponReward{},
		&model.GachaponRewardPool{},
		&model.GameConfig{},
		&model.Gift{},
		&model.Guild{},
		&model.HpMpAlert{},
		&model.Hwidaccount{},
		&model.Hwidban{},
		&model.Inventoryequipment{},
		&model.Inventoryitem{},
		&model.Inventorymerchant{},
		&model.Ipban{},
		&model.Keymap{},
		&model.LangResource{},
		&model.Macban{},
		&model.Macfilter{},
		&model.Makercreatedata{},
		&model.Makerreagentdata{},
		&model.Makerrecipedata{},
		&model.Makerrewarddata{},
		&model.Marriage{},
		&model.Medalmap{},
		&model.ModifiedCashItem{},
		&model.Monsterbook{},
		&model.Monstercarddata{},
		&model.MtsCart{},
		&model.MtsItem{},
		&model.Namechange{},
		&model.Newyear{},
		&model.Note{},
		&model.Nxcode{},
		&model.NxcodeItem{},
		&model.Nxcoupon{},
		&model.Pet{},
		&model.Petignore{},
		&model.Playerdisease{},
		&model.Playernpc{},
		&model.PlayernpcEquip{},
		&model.PlayernpcField{},
		&model.Plife{},
		&model.Questaction{},
		&model.Questprogress{},
		&model.Questrequirement{},
		&model.Queststatus{},
		&model.Quickslotkeymapped{},
		&model.Reactordrop{},
		&model.Report{},
		&model.Response{},
		&model.Ring{},
		&model.Savedlocation{},
		&model.ServerQueue{},
		&model.Shop{},
		&model.Shopitem{},
		&model.Skill{},
		&model.Skillmacro{},
		&model.Specialcashitem{},
		&model.Storage{},
		&model.Trocklocation{},
		&model.Wishlist{},
		&model.Worldtransfer{},
	)
}
