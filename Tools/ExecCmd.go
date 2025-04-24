package Tools

import (
	"bytes"
	"os/exec"
)

func RunCmd(cmdStr string, params string) (string, error) {

	//为第一个参数 添加一个引号
	//cmdStr = `"`+cmdStr+

	cmd := exec.Command(cmdStr, params)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return stderr.String(), err
	} else {
		return out.String(), nil
	}
}
