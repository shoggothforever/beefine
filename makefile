# 定义变量
OBJ := fynebpf

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

# 帮助信息
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        - Default target, runs the gen script."
	@echo "  run-gen    - Runs the gen script located in the scripts directory."
	@echo "  help       - Displays this help message."
