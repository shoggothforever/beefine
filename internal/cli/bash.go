package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// GetNamespaceID 获取指定进程的 Namespace ID
func GetNamespaceID(pid int, nsType string) (string, error) {
	nsPath := fmt.Sprintf("/proc/%d/ns/%s", pid, nsType)
	link, err := os.Readlink(nsPath)
	if err != nil {
		return "", err
	}

	// 提取 Namespace ID，例如从 "pid:[4026532902]" 提取 4026532902
	parts := strings.Split(link, ":")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid namespace format: %s", link)
	}
	nsID := strings.Trim(parts[1], "[]")
	return nsID, nil
}

// GetPeersInNamespace 获取与指定 Namespace ID 和类型处于同一 Namespace 的进程
func GetPeersInNamespace(discoveredPeers map[int]struct{}, nsID, nsType string) ([]int, error) {
	var peers []int
	procEntries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, entry := range procEntries {
		// 忽略非进程目录
		if !entry.IsDir() || !isNumeric(entry.Name()) {
			continue
		}

		pid, _ := strconv.Atoi(entry.Name())
		if _, ok := discoveredPeers[pid]; ok {
			continue
		}
		discoveredPeers[pid] = struct{}{}
		nsPath := fmt.Sprintf("/proc/%d/ns/%s", pid, nsType)
		link, err := os.Readlink(nsPath)
		if err != nil {
			continue
		}
		// 检查是否属于同一 Namespace
		if strings.Contains(link, nsID) {
			peers = append(peers, pid)
		}
	}
	return peers, nil
}

// GetProcessInfo 获取进程的详细信息
func GetProcessInfo(pid int) (string, error) {
	cmd := fmt.Sprintf("ps -p %d -o pid,ppid,user,comm --no-headers", pid)
	output, err := ExecCommand(cmd)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// execCommand 执行系统命令并返回输出
func ExecCommand(cmd string) (string, error) {
	output, err := exec.Command("bash", "-c", cmd).Output()
	return string(output), err
}

// isNumeric 检查字符串是否为数字
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
