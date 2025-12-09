package Tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// CmdError 自定义命令执行错误类型
type CmdError struct {
	Command string
	Args    []string
	Err     error
	Stderr  string
}

func (e *CmdError) Error() string {
	return fmt.Sprintf("命令执行失败: %s %v, 错误: %v, 标准错误输出: %s", e.Command, e.Args, e.Err, e.Stderr)
}

func (e *CmdError) Unwrap() error {
	return e.Err
}

// IsSafeCommandPath 检查命令路径是否安全
func IsSafeCommandPath(cmdPath string) bool {
	// 不允许包含危险字符
	dangerousChars := []string{";", "|", "&", "<", ">", "$", "`", "\"", "'"}
	for _, char := range dangerousChars {
		if strings.Contains(cmdPath, char) {
			return false
		}
	}
	return true
}

// IsSafeArgs 检查命令参数是否安全
func IsSafeArgs(args []string) bool {
	// 不允许参数中包含危险的shell命令结构
	dangerousPatterns := []string{"$(", "`", ";", "|"}
	for _, arg := range args {
		for _, pattern := range dangerousPatterns {
			if strings.Contains(arg, pattern) {
				return false
			}
		}
	}
	return true
}

// RunCmd 执行外部命令，带详细错误信息和安全检查
func RunCmd(cmdStr string, params ...string) (string, error) {
	// 1. 基本安全检查：命令路径不包含危险字符
	if !IsSafeCommandPath(cmdStr) {
		return "", fmt.Errorf("命令路径包含危险字符: %s", cmdStr)
	}

	// 2. 命令参数安全检查
	if !IsSafeArgs(params) {
		return "", fmt.Errorf("命令参数包含危险模式: %v", params)
	}

	// 3. 确保命令是可执行文件的绝对路径
	cmdPath, err := exec.LookPath(cmdStr)
	if err != nil {
		return "", fmt.Errorf("找不到命令 '%s': %v", cmdStr, err)
	}

	// 4. 进一步验证命令路径的安全性
	if !filepath.IsAbs(cmdPath) {
		// 如果exec.LookPath返回的不是绝对路径，手动获取绝对路径
		cmdPath, err = filepath.Abs(cmdPath)
		if err != nil {
			return "", fmt.Errorf("无法获取命令绝对路径: %v", err)
		}
	}

	// 5. 检查命令文件是否存在且可执行
	fileInfo, err := os.Stat(cmdPath)
	if err != nil {
		return "", fmt.Errorf("命令文件不存在或无法访问: %v", err)
	}
	if fileInfo.IsDir() {
		return "", fmt.Errorf("命令路径指向目录而非文件: %s", cmdPath)
	}
	// 检查文件是否可执行
	// Windows系统特殊处理：通过扩展名判断可执行性，而不是权限位
	if runtime.GOOS == "windows" {
		// Windows上检查常见的可执行文件扩展名
		ext := strings.ToLower(filepath.Ext(cmdPath))
		if ext != ".exe" && ext != ".bat" && ext != ".cmd" && ext != ".com" {
			return "", fmt.Errorf("命令文件不可执行: %s (Windows上只支持.exe, .bat, .cmd, .com扩展名)", cmdPath)
		}
	} else {
		// Unix/Linux系统通过权限位检查
		if fileInfo.Mode()&0111 == 0 {
			return "", fmt.Errorf("命令文件不可执行: %s", cmdPath)
		}
	}

	// 6. 执行命令（不使用shell，避免命令注入）
	cmd := exec.Command(cmdPath, params...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// 限制命令执行环境
	cmd.Env = []string{}
	cmd.Dir = "."

	err = cmd.Run()
	if err != nil {
		return stderr.String(), &CmdError{
			Command: cmdPath,
			Args:    params,
			Err:     err,
			Stderr:  stderr.String(),
		}
	}

	return out.String(), nil
}
