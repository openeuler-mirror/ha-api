package utils

import (
	"os/exec"
	"strings"

	"errors"

	"github.com/beego/beego/v2/core/logs"
)

// RunCommand runs the command and get the result
func RunCommand(c string) ([]byte, error) {
	// cmd := exec.Command(c)
	// return cmd.CombinedOutput()
	c = strings.Trim(c, " ")
	index := strings.Index(c, " ")
	if index < 0 {
		return nil, errors.New("invalid command")
	}

	cmd := c[0:index]
	parameter := c[index:]
	logs.Debug("Running command: %s", c)
	command := exec.Command(cmd, parameter)
	return command.CombinedOutput()
}
