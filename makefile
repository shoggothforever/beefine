# 定义变量
OBJ := beefine
ASSETS_DIR := internal/data/assets
LOGO_FILE := logo.png
SCRIPTS_DIR := scripts
GEN_SCRIPT := $(SCRIPTS_DIR)/gen.sh
PKG ?= default_arg
# 默认目标
.PHONY: all
all: build

# 执行 gen 脚本
.PHONY: gen
gen:
	@echo "Running gen script from $(GEN_SCRIPT)..."
	@chmod +x $(GEN_SCRIPT) # 确保脚本有执行权限
	@bash $(GEN_SCRIPT) $(PKG)

.PHON: build run
build:
	@echo "Running main project ..."
	@go generate ./bpf/...
	@go build -x -v
	@chmod +x $(OBJ)
	@echo "build $(OBJ) successfully ..."

.PHON: run
run: build
	@sudo ./$(OBJ)

.PHON: pkg_linux
pkg_linux:
	@go env -w GOOS=linux
	@go env -w GOARCH=amd64
	@fyne package -os linux -icon $(ASSETS_DIR)/$(LOGO_FILE)

.PHON: pkg_windows
pkg_windows:
	@go env -w GOOS=windows
	@go env -w GOARCH=amd64
	@fyne package -os windows -icon $(ASSETS_DIR)/$(LOGO_FILE)

.PHON: pkg_macos
pkg_macos:
	@go env -w GOOS=darwin
	@go env -w GOARCH=amd64
	@fyne package -os darwin -icon $(ASSETS_DIR)/$(LOGO_FILE)
# 帮助信息
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        - Default target, runs the gen script."
	@echo "  run-gen    - Runs the gen script located in the scripts directory."
	@echo "  help       - Displays this help message."
