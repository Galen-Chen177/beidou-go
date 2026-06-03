# internal/model/

数据模型层，定义与数据库表对应的 gorm 结构体。

## 文件

| 文件 | 对应表 | 职责 |
|---|---|---|
| `account.go` | `accounts` | 账号：用户名、密码(SHA-512)、PIN/PIC、GM 等级、封禁状态 |
| `character.go` | `characters` | 角色：等级、职业、属性(STR/DEX/INT/LUK)、HP/MP、经验、金币、外观、地图位置 |
| `item.go` | `inventoryitems` | 背包物品：物品 ID、格子位置、数量、过期时间 |
| `skill.go` | `skills` | 角色技能：技能 ID、等级、最大等级 |

## 约定

- 每个结构体显式指定 `TableName()` 方法
- 字段使用 gorm tag (`column:`, `primaryKey`, `uniqueIndex`, `index` 等)
- 字段名尽量与原项目数据库表对齐
- 后续需要补充的表：`equipments`（装备属性）、`quests`（任务进度）、`buddies`（好友）、`parties`（组队）、`guilds`（公会）等
