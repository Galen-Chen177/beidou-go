package script

import (
	"fmt"
	"os"
	"path/filepath"
)

// Loader 脚本文件加载器
type Loader struct {
	basePath string
}

// NewLoader 创建加载器
// basePath: 脚本根目录（如 "scripts"）
func NewLoader(basePath string) *Loader {
	return &Loader{basePath: basePath}
}

// Load 加载指定脚本文件的内容
// category: 子目录 (npc/quest/event/reactor)
// name: 脚本文件名（不含路径）
func (l *Loader) Load(category, name string) (string, error) {
	path := filepath.Join(l.basePath, category, name)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("加载脚本 %s/%s 失败: %w", category, name, err)
	}
	return string(data), nil
}

// Exists 检查脚本文件是否存在
func (l *Loader) Exists(category, name string) bool {
	path := filepath.Join(l.basePath, category, name)
	_, err := os.Stat(path)
	return err == nil
}

// List 列出指定类别下的所有脚本文件
func (l *Loader) List(category string) ([]string, error) {
	dir := filepath.Join(l.basePath, category)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("读取脚本目录 %s 失败: %w", dir, err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".js" {
			names = append(names, e.Name())
		}
	}
	return names, nil
}
