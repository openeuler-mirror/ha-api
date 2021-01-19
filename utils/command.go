package utils

import (
	"os/exec"

	"github.com/beego/beego/v2/core/logs"
)

// RunCommand runs the command and get the result
func RunCommand(c string) ([]byte, error) {
	logs.Debug("Running command: %s", c)
	command := exec.Command("bash", "-c", c)
	return command.CombinedOutput()
}
