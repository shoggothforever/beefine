package logic

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

const hintStr = "eBPF helpers supported for program type"

func ListProbes() map[string][]string {
	// 执行 `bpftool feature probe` 命令
	cmd := exec.Command("bpftool", "feature", "probe")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing bpftool: %v\n", err)
		return nil
	}

	// 解析输出
	helperMap := extractHelpers(string(output))

	// 打印结果
	//fmt.Println("eBPF Helpers by Program Type:")
	//for progType, helpers := range helperMap {
	//	fmt.Printf("%s:\n", progType)
	//	for _, helper := range helpers {
	//		fmt.Printf("  - %s\n", helper)
	//	}
	//}
	return helperMap
}

// extractHelpers 解析 `bpftool feature probe` 输出，提取每种程序类型支持的 helpers
func extractHelpers(output string) map[string][]string {
	helperMap := make(map[string][]string)
	var currentProgType string
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 检查是否是程序类型的起始行
		if strings.HasPrefix(line, hintStr) {
			// 提取程序类型
			parts := strings.Split(line, " ")
			if len(parts) > 6 {
				currentProgType = parts[len(parts)-1]
			} else {
				currentProgType = "unknown"
			}
			helperMap[currentProgType] = []string{} // 初始化该程序类型的 helpers 列表
			continue
		}

		// 如果当前处于一个程序类型 block 中，记录 helpers
		if currentProgType != "" && strings.HasPrefix(line, "- ") {
			helper := strings.TrimPrefix(line, "- ")
			helperMap[currentProgType] = append(helperMap[currentProgType], helper)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading bpftool output: %v\n", err)
	}

	return helperMap
}
