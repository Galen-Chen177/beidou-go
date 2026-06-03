# scripts/

游戏脚本文件目录，存放冒险岛的 NPC 对话脚本、任务脚本、事件脚本、Reactor 脚本。

所有脚本为 **JavaScript (ES5.1)** 格式，由 goja 引擎执行。

## 目录

| 目录 | 内容 | 脚本命名规则 |
|---|---|---|
| `npc/` | NPC 对话脚本 | `{npc_id}.js`，如 `1002000.js` |
| `quest/` | 任务脚本 | 按任务 ID 组织 |
| `event/` | 事件脚本 | 按事件名组织 |
| `reactor/` | Reactor（反应堆）脚本 | `{reactor_id}.js` |

## 来源

脚本文件从原项目 [BeiDou-Server](https://github.com/BeiDouMS/BeiDou-Server) 的 `script-zh-CN/` 目录复制。

当前目录为空（仅占位），待第一期骨架搭建完成后再导入。
