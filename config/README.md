# config/

配置层，定义配置结构体和 YAML 文件加载逻辑。

## 文件

| 文件 | 职责 |
|---|---|
| `config.go` | 定义 `Config` 结构体及其子结构（Server / Database / Login / Channel / Script），提供 `Load()` 和 `Default()` 函数 |
| `config.yaml` | 默认配置文件模板，包含 MySQL 连接、各端口、脚本路径等参数 |

## 设计

- 使用 `gopkg.in/yaml.v3` 直接 Unmarshal，不依赖 viper
- `Default()` 返回内置默认值，`Load()` 在此基础上覆盖 YAML 中的字段
- 没有 "魔法合并"（环境变量覆盖、多源合并等），一切显式可控

## 配置结构

```
Config
├── ServerConfig      # 服务器名称、监听地址、日志文件
├── DatabaseConfig    # MySQL 主机/端口/用户/密码/库名
├── LoginConfig       # 登录服务器端口
├── ChannelConfig     # 频道服务器端口、最大玩家数
└── ScriptConfig      # 脚本目录路径、热加载开关
```
