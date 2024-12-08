package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	containerStateRunning = "running"
	containerStateExited  = "exited"
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
	ctx         = context.Background()
)

func init() {
	cliInstance, _ = initDockerClient()
}

// GetDockerClient 返回单例 Docker 客户端
func GetDockerClient() (*client.Client, error) {
	var err error
	once.Do(func() {
		cliInstance, _ = initDockerClient()
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
func ExecDockerCmd(input string) (string, error) {
	args := strings.Split(input, " ")
	if len(args) < 2 {
		return "", nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to execute docker run: %w", err)
	}
	fmt.Println(string(output))
	return string(output), nil
}

// ParseAndRunDockerRun 解析 JSON 并运行 docker 命令
/*
参考示例配置
`{
    "image": "nginx",
    "name": "my-container",
    "ports": ["80:80", "443:443"],
    "volumes": ["/host/path:/container/path"],
    "env": ["ENV_VAR1=value1", "ENV_VAR2=value2"],
    "detach": true,
    "rm": true
}`
*/
func ParseAndRunDockerRun(jsonConfig string, image string) (string, error) {
	var config DockerRunConfig
	err := json.Unmarshal([]byte(jsonConfig), &config)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}
	// 构造 docker run 命令
	cmdArgs := []string{"run"}
	if len(config.Image) == 0 {
		if len(image) == 0 {
			return "no image specified", nil
		}
		config.Image = image
	}
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
	cmdArgs = append(cmdArgs, "-d")

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
	// 检查镜像名称是否合法
	if imageName == "" {
		return "", fmt.Errorf("image name cannot be empty")
	}
	// 设置拉取选项
	options := image.PullOptions{}
	// 调用 ImagePull 方法
	reader, err := cliInstance.ImagePull(ctx, imageName, options)
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

// ListImage 获取系统中的镜像
func ListImage() ([]image.Summary, error) {
	images, err := cliInstance.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return nil, err
	}
	return images, nil

}

// ListContainer 获取系统中的镜像
func ListContainer() ([]types.Container, error) {
	containers, err := cliInstance.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}
	return containers, nil

}

func GetContainerStat(id string) (types.ContainerJSON, error) {
	return cliInstance.ContainerInspect(ctx, id)
}

// CheckContainerRunningState 如果容器正常运行返回true
func CheckContainerRunningState(status string) bool {
	if status == containerStateRunning {
		return true
	}
	return false
}

func ChangeContainerState(id string, oldState bool) error {
	if oldState {
		return cliInstance.ContainerStop(ctx, id, container.StopOptions{})
	} else {
		return cliInstance.ContainerStart(ctx, id, container.StartOptions{})
	}
}

type DashBoardData struct {
	CpuUsage     float64
	MemUsage     float64
	ContainerLen int
	ImagesLen    int
}

func GetDockerDashBoardData() (*DashBoardData, error) {
	containers, err := cliInstance.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}
	images, err := cliInstance.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return nil, err
	}
	var data DashBoardData
	data.ContainerLen = len(containers)
	data.ImagesLen = len(images)
	for _, v := range containers {
		stats, err := cliInstance.ContainerStats(context.Background(), v.ID, false)
		if err != nil {
			log.Fatalf("Error getting container stats: %v", err)
		}
		defer stats.Body.Close()
		// 解析并显示统计信息
		var statsJSON container.StatsResponse
		if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err != nil {
			log.Fatalf("Error decoding stats: %v", err)
		}
		// 输出容器的 CPU 和内存使用情况
		if statsJSON.CPUStats.SystemUsage != 0 {
			data.CpuUsage += float64(statsJSON.CPUStats.CPUUsage.TotalUsage) / float64(statsJSON.CPUStats.SystemUsage) * 100

		}
		if statsJSON.MemoryStats.Limit != 0 {
			data.MemUsage += float64(statsJSON.MemoryStats.Usage / statsJSON.MemoryStats.Limit / (1024 * 1024)) // 转换为 MB
		}
		//fmt.Printf("CPU Usage: %.2f%%\n", cpuUsage)
		//fmt.Printf("Memory Usage: %.2fMB / %.2fMB\n", float64(memUsage), float64(statsJSON.MemoryStats.Limit)/(1024*1024)) // 显示内存限制

	}

	return &data, nil
}
