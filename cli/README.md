# cli/

CLI 层，基于 `urfave/cli/v3` 实现命令行参数解析和子命令路由。

## 文件

| 文件 | 职责 |
|---|---|
| `root.go` | 根命令定义：`beidou-go` 应用、全局 `--config` flag、注册 `serve` / `migrate` 子命令 |
| `serve.go` | `serve` 子命令：加载配置 → 初始化日志 → 连接数据库 → 启动 Login + Channel 服务器 → 等待退出信号 |
| `migrate.go` | `migrate` 子命令：加载配置 → 连接数据库 → 执行 gorm AutoMigrate |

## 调用链

```
cmd/server/main.go
  → cli.Run()
    → serve / migrate 子命令
```
