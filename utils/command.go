package utils

import (
	"os/exec"

	"github.com/beego/beego/v2/core/logs"
)

// RunCommand runs the command and get the result
func RunCommand(c string) ([]byte, error) {
	logs.Debug("Running command: %s", c)
	command := exec.Command("bash", "-c", c)
	out, err := command.CombinedOutput()
	if err != nil {
		logs.Error("Run command failed!, command: " + c + " out: " + string(out) + " err: " + err.Error())
	}
	return out, err
}
