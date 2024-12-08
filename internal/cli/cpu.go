package cli

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type CPUStats struct {
	User   uint64
	System uint64
	Idle   uint64
}

func GetCPUStats() (CPUStats, error) {
	// 读取 /proc/stat 文件
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return CPUStats{}, fmt.Errorf("failed to read /proc/stat: %w", err)
	}

	// 查找 CPU 行
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			user, _ := strconv.ParseUint(fields[1], 10, 64)
			system, _ := strconv.ParseUint(fields[3], 10, 64)
			idle, _ := strconv.ParseUint(fields[4], 10, 64)
			return CPUStats{User: user, System: system, Idle: idle}, nil
		}
	}

	return CPUStats{}, fmt.Errorf("no cpu stats found in /proc/stat")
}

func CalculateCPUUsage(prevStats, currStats CPUStats) float64 {
	// 计算 CPU 使用时间差
	totalPrev := prevStats.User + prevStats.System + prevStats.Idle
	totalCurr := currStats.User + currStats.System + currStats.Idle

	// 计算 CPU 使用率
	totalDiff := totalCurr - totalPrev
	idleDiff := currStats.Idle - prevStats.Idle

	// 防止除以零
	if totalDiff == 0 {
		return 0
	}

	return float64(totalDiff-idleDiff) / float64(totalDiff) * 100
}
