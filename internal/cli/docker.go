package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// DockerRunConfig 定义 docker run 的 JSON 配置结构
type DockerRunConfig struct {
	Image   string   `json:"image"`             // 镜像名称
	Name    string   `json:"name,omitempty"`    // 容器名称
	Ports   []string `json:"ports,omitempty"`   // 端口映射
	Volumes []string `json:"volumes,omitempty"` // 挂载卷
	Env     []string `json:"env,omitempty"`     // 环境变量
	Detach  bool     `json:"detach,omitempty"`  // 后台运行
	Remove  bool     `json:"rm,omitempty"`      // 自动删除容器
}

// initDockerClient 初始化 Docker 客户端
func initDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

// isConnectionValid 检查 Docker 客户端连接是否有效
func isConnectionValid(cli *client.Client) bool {
	if cli == nil {
		return false
	}
	// 尝试调用 Docker API 检查连接
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := cli.Ping(ctx)
	return err == nil
}

var (
	cliInstance *client.Client
	once        sync.Once
	mu          sync.Mutex // 保证多线程安全
)

// GetDockerClient 返回单例 Docker 客户端
func GetDockerClient() (*client.Client, error) {
	var err error

	once.Do(func() {
		cliInstance, err = initDockerClient()
	})
	// 检查连接是否有效
	if !isConnectionValid(cliInstance) {
		fmt.Println("Docker client connection is invalid. Reinitializing...")
		mu.Lock()
		defer mu.Unlock()
		cliInstance, err = initDockerClient()
		if err != nil {
			return nil, err
		}
	}
	return cliInstance, err
}

// ParseAndRunDockerRun 解析 JSON 并运行 docker 命令
func ParseAndRunDockerRun(jsonConfig string) (string, error) {
	var config DockerRunConfig
	err := json.Unmarshal([]byte(jsonConfig), &config)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}
	// 构造 docker run 命令
	cmdArgs := []string{"run"}

	// 可选参数
	if config.Name != "" {
		cmdArgs = append(cmdArgs, "--name", config.Name)
	}
	if config.Remove {
		cmdArgs = append(cmdArgs, "--rm")
	}
	for _, port := range config.Ports {
		cmdArgs = append(cmdArgs, "-p", port)
	}
	for _, volume := range config.Volumes {
		cmdArgs = append(cmdArgs, "-v", volume)
	}
	for _, env := range config.Env {
		cmdArgs = append(cmdArgs, "-e", env)
	}
	if config.Detach {
		cmdArgs = append(cmdArgs, "-d")
	}

	// 必须的参数（镜像名）
	cmdArgs = append(cmdArgs, config.Image)

	// 打印生成的命令（可选）
	fmt.Println("Executing command:", "docker", strings.Join(cmdArgs, " "))

	// 执行 docker run 命令
	cmd := exec.Command("docker", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to execute docker run: %w", err)
	}

	return string(output), nil
}

// PullDockerImage 使用 Docker CLI 拉取指定镜像
func PullDockerImage(imageName string) (string, error) {
	// 创建上下文
	ctx := context.Background()

	// 检查镜像名称是否合法
	if imageName == "" {
		return "", fmt.Errorf("image name cannot be empty")
	}
	cli, err := GetDockerClient()
	if err != nil {
		return "", err
	}
	// 设置拉取选项
	options := image.PullOptions{}
	// 调用 ImagePull 方法
	reader, err := cli.ImagePull(ctx, imageName, options)
	if err != nil {
		return "", fmt.Errorf("failed to pull image '%s': %w", imageName, err)
	}
	defer reader.Close()
	// 捕获命令输出
	var outBuffer bytes.Buffer

	_, err = io.Copy(&outBuffer, reader)
	if err != nil {
		return "", fmt.Errorf("failed to read image pull output: %w", err)
	}

	// 返回成功信息
	return outBuffer.String(), nil
}
