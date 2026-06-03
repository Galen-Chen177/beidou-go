# internal/script/

JS 脚本引擎层，基于 [goja](https://github.com/dop251/goja) 实现 NPC 对话、任务脚本、事件脚本的执行。

## 文件

| 文件 | 职责 |
|---|---|
| `engine.go` | `Engine` — goja 运行时封装：Runtime 对象池（优化 GC）、`Run()` 执行 JS 脚本、`NewRuntime()` 创建独立运行时（用于需要保留状态的 NPC 对话）。`Bindings` — 桥接上下文，注入 Go 侧实现的 `cm`/`pi`/`ms`/`em` 等对象到 JS 运行时 |
| `loader.go` | `Loader` — 脚本文件加载器：从 `scripts/` 目录按类别( npc/quest/event/reactor )加载 `.js` 文件、检查文件是否存在、列出所有脚本 |

## 桥接对象（待实现）

JS 脚本通过以下全局对象调用 Go 侧方法：

| 对象 | 全称 | 用途 |
|---|---|---|
| `cm` | ConversationManager | NPC 对话核心（100+ 方法）：sendNext、sendSimple、warp、dispose、gainItem 等 |
| `pi` | PlayerInteraction | 玩家交互（组队任务等） |
| `ms` | MapScriptManager | 地图脚本 |
| `em` | EventManager | 事件管理器 |

## 脚本类型

| 类别 | 目录 | 触发方式 |
|---|---|---|
| NPC 脚本 | `scripts/npc/` | 玩家点击 NPC |
| 任务脚本 | `scripts/quest/` | 任务接取/进行中/完成时触发 |
| 事件脚本 | `scripts/event/` | 定时器或特定条件触发 |
| Reactor 脚本 | `scripts/reactor/` | 玩家攻击反应堆（如箱子、药草）时触发 |
