package script

import (
	"fmt"
	"sync"

	"github.com/dop251/goja"
)

// Engine goja JS 脚本引擎封装
type Engine struct {
	pool sync.Pool // goja.Runtime 对象池
}

// NewEngine 创建脚本引擎
func NewEngine() *Engine {
	return &Engine{
		pool: sync.Pool{
			New: func() any {
				return goja.New()
			},
		},
	}
}

// Run 在沙箱中执行一段 JS 代码
func (e *Engine) Run(src string) (goja.Value, error) {
	vm := e.pool.Get().(*goja.Runtime)
	defer e.pool.Put(vm)

	// 清理运行时（防止上次执行残留）
	// goja 没有直接的 Reset 方法，但复用是安全的，只要不依赖全局状态

	val, err := vm.RunString(src)
	if err != nil {
		return nil, fmt.Errorf("脚本执行失败: %w", err)
	}
	return val, nil
}

// NewRuntime 获取一个独立的 goja.Runtime（用于需要保留状态的场景，如 NPC 对话）
func (e *Engine) NewRuntime() *goja.Runtime {
	return goja.New()
}

// Bindings 桥接上下文：为指定 runtime 注入 Go 侧提供的 API 对象
// 每个 NPC/任务脚本启动时调用，注入 cm、pi 等对象
type Bindings struct {
	vm *goja.Runtime
}

// NewBindings 创建桥接上下文
func NewBindings(vm *goja.Runtime) *Bindings {
	return &Bindings{vm: vm}
}

// SetCM 注入 ConversationManager (cm) 对象（待实现）
func (b *Bindings) SetCM() {
	// TODO: 实现 cm 对象的方法绑定
	// cm.sendNext(), cm.sendSimple(), cm.warp(), cm.dispose() 等
}

// SetPI 注入 PlayerInteraction (pi) 对象（待实现）
func (b *Bindings) SetPI() {
	// TODO: 实现 pi 对象的方法绑定
}

// Set 注入任意 Go 对象到 JS 运行时
func (b *Bindings) Set(name string, obj any) {
	b.vm.Set(name, obj)
}
