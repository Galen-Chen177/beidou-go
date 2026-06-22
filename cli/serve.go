package cli

import (
	"beidou-go/config"
	"beidou-go/internal/network"
	"beidou-go/internal/server/channel"
	"beidou-go/internal/server/login"
	loginhandler "beidou-go/internal/server/login/handler"
	"beidou-go/internal/store"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

func cmdServe() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "启动游戏服务",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runServe()
		},
	}
}

func runServe() error {
	// 加载配置
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 初始化日志
	log = logrus.New()
	log.SetLevel(logrus.DebugLevel) // 骨架阶段用 Debug 级别，方便看 hex dump 对标 Wireshark
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableQuote:    true, // 允许 hex dump 中的 \n 正常换行
	})
	if cfg.Server.LogFile != "" {
		f, err := os.OpenFile(cfg.Server.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			return fmt.Errorf("打开日志文件失败: %w", err)
		}
		log.SetOutput(f)
	}
	log.Infof("加载配置完成: %s", cfgPath)

	// 将 logger 注入 network 包，使 Session 的 hex dump 日志能正常输出
	network.Log = log

	// 初始化数据库
	if err := store.InitDB(cfg.Database); err != nil {
		return fmt.Errorf("数据库初始化失败: %w", err)
	}
	log.Info("数据库连接成功")

	// 自动迁移表结构
	if err := store.AutoMigrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}
	log.Info("数据库迁移完成")

	// 创建 AccountStore
	accountStore := store.NewAccountStore(store.DB())

	// 创建登录相关组件
	coordinator := login.NewSessionCoordinator()
	authHandler := loginhandler.NewAuthHandler(accountStore, coordinator, cfg.Login.AutoRegister, log)

	// 创建 TCP 服务器
	tcpSrv := network.NewTCPServer(cfg.Server.Host)
	log.Infof("TCP 引擎就绪")

	// 启动登录服务器
	loginSrv := login.NewServer(cfg, tcpSrv, log, authHandler)
	go func() {
		if err := loginSrv.Start(); err != nil {
			log.Fatalf("登录服务器启动失败: %v", err)
		}
	}()
	log.Infof("登录服务器启动 @ :%d", cfg.Login.Port)

	// 启动频道服务器
	channelSrv := channel.NewServer(cfg, tcpSrv, log)
	go func() {
		if err := channelSrv.Start(); err != nil {
			log.Fatalf("频道服务器启动失败: %v", err)
		}
	}()
	log.Infof("频道服务器启动 @ :%d", cfg.Channel.Port)

	log.Infof("=== %s 启动完成 ===", cfg.Server.Name)

	// 等待退出信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Info("正在关闭服务器...")
	tcpSrv.Shutdown()
	log.Info("服务器已关闭")
	return nil
}
