# migrations/

数据库迁移脚本，存放 SQL 初始化脚本。

## 文件

| 文件 | 说明 |
|---|---|
| `001_init.sql` | 初始建表 SQL（待从原项目提取） |

## 说明

当前使用 gorm 的 `AutoMigrate` 自动管理表结构。手动 SQL 迁移文件作为备份和参考，从原项目 `BeiDou-Server` 的 SQL 初始化脚本提取。
