.PHONY: build run clean test

# 二进制输出目录
BIN_DIR := bin
BIN_NAME := beidou-go

# 默认配置
CONFIG ?= config/config.yaml

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BIN_NAME) ./cmd/server/

run: build
	$(BIN_DIR)/$(BIN_NAME) serve --config=$(CONFIG)

clean:
	rm -rf $(BIN_DIR)

test:
	go test ./...

# 开发模式：文件监听自动重启（需安装 air）
dev:
	air -- --config=$(CONFIG)

# 下载依赖
deps:
	go mod tidy
	go mod download
