# internal/store/

数据访问层（DAO），封装 gorm 数据库操作。

## 文件

| 文件 | 职责 |
|---|---|
| `db.go` | 数据库初始化：连接 MySQL、配置连接池（最大 25 连接 / 空闲 10 连接）、`AutoMigrate()` 自动创建/更新表结构。暴露全局 `DB()` 方法供其他包获取 `*gorm.DB` 实例 |

## 使用方式

```go
import "beidou-go/internal/store"

// 查询
var acc model.Account
store.DB().Where("name = ?", username).First(&acc)

// 创建
store.DB().Create(&model.Character{...})

// 更新
store.DB().Model(&acc).Update("lastlogin", time.Now())
```

## 后续扩展

当前 `store/` 只包含数据库初始化，后续各模块的 DAO 操作（如 `account_store.go`、`character_store.go`）直接放在 `store/` 目录下，不按模块分子目录。
