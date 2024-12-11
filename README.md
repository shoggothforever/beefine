# beefine: Visualizing Docker Internals with eBPF

## 项目简介

**Beefine** 是一个基于 **Fyne** 和 **Cilium eBPF** 框架开发的工具，旨在通过图形化交互界面（GUI）实时观测 Docker 容器的创建过程，深入理解虚拟化技术的核心理念和实现原理。本项目同时支持加载和管理 eBPF 程序，帮助用户追踪操作系统在 Docker 操作中的行为，后续将扩展到 Kubernetes 集群的 Pod 监控。

---

## 创建初衷

容器化技术已经成为现代软件开发和部署的核心工具，而 Docker 是其中的代表。然而，在使用容器时，很多开发者并不清楚容器创建的具体过程，例如操作系统如何处理镜像解压、文件系统挂载和网络隔离等底层细节。

为了帮助开发者更直观地理解这些细节，**Beefine** 通过 eBPF 技术捕获和分析容器的系统调用行为，并借助 Fyne 提供图形化的实时反馈，最终实现以下目标：

1. **简化学习过程**：

    - 通过实时观测容器创建过程，帮助开发者更好地理解虚拟化技术。
2. **提高可视化交互体验**：

    - 提供直观的图形界面，展示关键的系统行为和资源变化。
3. **开发工具化**：

    - 为学习者和工程师提供一个可以随时实验和验证的工具，减少操作系统实验的门槛。

---

## 项目用途

该项目目前专注于以下场景：

1. **Docker 容器创建观测**：

    - 实时加载 eBPF 程序，捕获 Docker 使用镜像创建容器时的操作系统行为。
    - 包括镜像拉取、文件系统设置、命名空间管理、网络配置等。
    - 为开发者提供实验环境，用于理解容器的实现原理，如 Namespace、Cgroup 和 UnionFS。
2. **动态 eBPF 程序管理**：

    - 提供直观界面加载和管理 eBPF 程序，动态分析系统行为。
3. **图形化分析**：

    - 通过 GUI 展示分析结果，帮助用户直观理解容器的底层运行机制。
4. **未来扩展**：

    - 为 Kubernetes 集群提供监控支持，观测 Pod 的创建、调度和运行。

---

## 项目进度

### 已实现功能

1. **Docker 创建过程观测**：

    - 使用 eBPF 追踪操作系统调用，捕获 Docker 使用镜像创建容器的全过程。
2. **图形化界面**：

    - 使用 Fyne 开发可交互的 GUI，包括日志查看、动态程序加载等功能。
3. **实时 eBPF 程序加载**：

    - 支持用户加载自定义 eBPF 程序，动态分析特定行为。

### 待开发功能

1. **Kubernetes 集群观测**：

    - 设计用于追踪 Kubernetes Pod 调度与运行的功能模块。
2. **镜像与容器的性能分析**：

    - 提供更多统计功能，分析资源使用情况（CPU、内存、I/O 等）。
3. **历史数据管理**：

    - 支持保存和回放观测结果，便于后续分析。
4. **优化 GUI**：

    - 增加更多交互功能，例如高级过滤、实时图表更新。

---

## 功能特性

1. **实时观测**：

    - 动态加载 eBPF 程序，实时分析 Docker 创建容器过程中的关键操作。
2. **图形化界面**：

    - 使用 Fyne 提供直观、易用的 GUI，帮助用户快速理解分析结果。
3. **可扩展性**：

    - 为未来支持 Kubernetes 和其他虚拟化技术的监控分析提供基础。
4. **学习友好**：

    - 为开发者提供操作系统行为的可视化教学工具。

---

## 当前开发进度

1. **功能模块**：

    - [X]  图形界面搭建（基于 Fyne）
    - [X]  支持动态加载 eBPF 程序
    - [X]  基础的 Docker 行为观测功能（镜像拉取、容器创建）
    - [ ]  数据可视化（实时图表展示系统调用频率）
    - [ ]  丰富的日志和分析工具
2. **观测重点**：

    - Docker 使用镜像创建容器的全过程：
        - [X] 镜像文件的拉取与解压。
        - [ ] 文件系统的挂载（OverlayFS）。
        - [ ] 容器隔离环境的设置（Namespace 和 Cgroup）。

- Docker 容器运行中的性能观测

    - [ ] 网络连接
    -  [ ]
- 后续计划扩展到更多场景，包括：

    - [ ] 容器运行时的性能分析。
    - [ ] Kubernetes 集群中的容器行为观测。

## 快速开始文档

### 系统要求

- **操作系统**：Ubuntu 20.04 或更高版本
- **架构**：AMD64
- **依赖工具**：

    - **Go**：版本 >= 1.18
    - **Docker**：需要 Docker Daemon 运行
    - **权限**：eBPF 程序需要管理员权限运行

---

## 1. 编译源码

以下是编译源码的方法，适用于开发者或需要自定义功能的用户。

### 安装依赖

在编译源码前，请确保系统已安装以下依赖：

1. **安装 Go** 如果未安装 Go，可以通过以下命令安装：
2. bash
3. 复制代码
4. `sudo apt update ``sudo apt install -y golang`
5. 检查安装是否成功：
6. `go version`
7. 确保版本 >= 1.18。
8. **安装 Docker** 如果未安装 Docker，可以参考以下命令安装：
9. bash
10. 复制代码
11. `sudo apt update sudo apt install -y docker.io ``sudo systemctl enable --now docker`
12. **确保管理员权限** eBPF 程序需要管理员权限，确保当前用户在 `sudo` 或 `root` 下运行。

### 克隆项目源码

`git clone ``https://github.com/shoggothforever/beefine.git`
`cd beefine`

### 编译程序

运行以下命令编译项目源码：

`go build -o beefine`

如果需要静态编译以避免依赖动态库，可以使用以下命令：

`go build -ldflags "-s -w" -o beefine`

### 运行程序

编译成功后，执行以下命令运行：

`sudo ./beefine`

程序将启动图形化界面。

---

## 2. 运行二进制应用

以下是直接运行已编译二进制文件的方法，适用于不需要修改源码的用户。

### 下载二进制文件

从项目的 [Release 页面](https://github.com/shoggothforever/beefine/releases) 下载对应版本的二进制文件。

以 Ubuntu 20.04 为例：

bash

复制代码

`wget https://github.com/your-repo/beefine/releases/download/v1.0.0/beefine-linux-amd64 -O beefine`
`chmod +x beefine`

### 安装依赖

确保以下依赖已安装：

1. **Docker**： 安装命令参考上述 **编译源码** 部分。
2. **权限设置**： 确保当前用户拥有运行 eBPF 程序的权限。

### 运行程序

运行以下命令启动程序：

`sudo ./beefine`

---

## 文件结构 `

``
├── bpf/                     # 主程序入口 ` `
│   ├──*/                    # bpf2go与libbpf结合bpf程序 ` `
│   └──vmlinux.h             # bpf 的 btf文件 ` `
├── pkg/                     # 可以导出的包,提供可以复用的组件和逻辑代码 ` `
│   │── gui/                 # Fyne GUI 相关代码 ` `
│   │  ├── themes/           # 存放fyneUI 设计的功能代码 ` `
│   │  └── watcher/          # Docker 功能界面 ` `
│   └── component/           # 存放自定义的的fyne组件 ` `
├── internal/                # Docker 操作工具 ` `
│   ├── cli/                 # Docker以及脚本交互管理 ` `
│   ├── data/                # fyne资源管理 ` `
│   └── helper               # 通用辅助函数 ` `
├── configs/                 # 项目配置  ` `
├── scripts/                 # 存放脚本文件 ` `
├── test/                    # 存放测试用例以及脚本 ` `
├── main.go                  # GO 程序入口
├── go.mod                   # Go 模块文件 ` `
├── makefile                 # 项目编译脚本
├── license                  # 证书文件
└── README.md                # 项目文档 `

---

## 使用方法

### 1. 加载 eBPF 程序

在 GUI 中选择 `Load eBPF Program` 功能，上传您的 eBPF 程序文件（如 ELF 格式），程序将自动加载并运行。

### 2. 观测 Docker 操作

在 GUI 中选择 `Observe Docker`，输入需要观测的 Docker 操作（如容器创建），实时查看捕获的系统行为和分析结果。

---

## 未来计划

1. **Kubernetes 支持**：

    - 增加对 Kubernetes 集群中 Pod 创建过程的监控功能。
2. **实时性能监控**：

    - 使用 eBPF 追踪资源使用（CPU、内存、网络流量等），分析容器性能。
3. **数据可视化**：

    - 增加更多实时图表和分析报告，提升用户体验。
4. **跨平台支持**：

    - 提供更多平台的支持，兼容 Windows 和 macOS。

---

## 贡献指南

欢迎对本项目感兴趣的开发者贡献代码！贡献步骤如下：

1. Fork 本项目到您的 GitHub 账户。
2. 提交 Pull Request 前，请确保所有功能通过测试。
3. 在 Pull Request 中详细描述您的修改内容和目的。

---

## 开源协议

本项目采用 [MIT License](https://opensource.org/licenses/MIT) 开源协议，详情请参见 LICENSE 文件。
