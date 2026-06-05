package network

import "github.com/sirupsen/logrus"

// Log 是网络层的包级 logger，由上层在启动时注入。
// 默认使用 logrus 全局 logger，避免 nil panic。
var Log = logrus.StandardLogger()
