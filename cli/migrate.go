package cli

import (
	"context"
	"fmt"

	"beidou-go/config"
	"beidou-go/internal/store"

	"github.com/urfave/cli/v3"
)

func cmdMigrate() *cli.Command {
	return &cli.Command{
		Name:  "migrate",
		Usage: "执行数据库迁移（创建/更新表结构）",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runMigrate()
		},
	}
}

func runMigrate() error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 仅初始化数据库连接用于迁移
	if err := store.InitDB(cfg.Database); err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	fmt.Println("数据库迁移完成!")
	return nil
}
