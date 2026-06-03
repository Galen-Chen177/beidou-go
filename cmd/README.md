# cmd/

程序入口目录，每个子目录对应一个可执行文件。

## server/

游戏服务端入口 (`main.go`) — 调用 `cli/` 包启动整个服务进程。整个应用只有一个 `main.go`，Login Server 和 Channel Server 在同一个进程内通过 goroutine 并发运行。
