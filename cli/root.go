package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var (
	log     *logrus.Logger
	cfgPath string
)

// Run 启动 CLI 应用，阻塞直到服务退出
func Run() {
	app := &cli.Command{
		Name:  "beidou-go",
		Usage: "冒险岛 GMS v0.83 服务端 (Go)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       "config/config.yaml",
				Usage:       "配置文件路径",
				Destination: &cfgPath,
			},
		},
		Commands: []*cli.Command{
			cmdServe(),
			cmdMigrate(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
