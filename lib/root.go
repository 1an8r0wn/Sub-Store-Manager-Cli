package lib

import (
	"os/exec"
	"strings"
)

// CheckPort 检查端口是否可用，可用-true 不可用-false
func CheckPort(port string) bool {
	cmd := exec.Command("docker", "ps", "--format", "{{.Ports}}")
	output, err := cmd.Output()
	if err != nil {
		return false // 命令执行失败，视为端口不可用
	}

	// 检查端口是否在输出中
	ports := strings.Split(string(output), "\n")
	for _, p := range ports {
		if strings.Contains(p, ":"+port+"->") {
			return false // 端口被占用
		}
	}

	return true
}
