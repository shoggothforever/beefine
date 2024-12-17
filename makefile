# 定义变量
OBJ := beefine
LOGO_FILE := internal/data/assets/logo.png
SCRIPTS_DIR := scripts
GEN_SCRIPT := $(SCRIPTS_DIR)/gen.sh
PKG ?= default_arg
PKG_OS := linux
# 默认目标
.PHONY: all
all: build

# 执行 gen 脚本
.PHONY: gen
gen:
	@echo "Running gen script from $(GEN_SCRIPT)..."
	@chmod +x $(GEN_SCRIPT) # 确保脚本有执行权限
	@bash $(GEN_SCRIPT) $(PKG)

.PHON: build
build:
	@echo "Running main project ..."
	@BPF2GO_FLAGS="-O2 -g -Wall -Werror -fbuiltin$(CFLAGS)" go generate  ./bpf/...
	@go build -x -v
	@chmod +x $(OBJ)
	@echo "build $(OBJ) successfully ..."

.PHON: run
run: build
	@sudo ./$(OBJ)

.PHON: package linux windows macos
package:
	@go env -w CGO_ENABLED=1
	@fyne package -os $(PKG_OS) -icon $(LOGO_FILE)

.PHON: linux
linux:
	@go env -w GOOS=linux
	@go env -w GOARCH=amd64
	@PKG_OS=linux

.PHON: windows
windows:
	@go env -w GOOS=windows
	@go env -w GOARCH=amd64
	@PKG_OS=windows

.PHON: macos
macos:
	@go env -w CC=darwin
	@go env -w GOOS=darwin
	@go env -w GOARCH=amd64
	@PKG_OS=darwin
.PHON: clean
clean:
	@rm -r *.tar.xz
# 帮助信息
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        - Default target, runs the gen script."
	@echo "  run-gen    - Runs the gen script located in the scripts directory."
	@echo "  help       - Displays this help message."
