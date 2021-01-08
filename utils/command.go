package utils

import (
	"github.com/beego/beego/v2/core/logs"
)

// RunCommand runs the command and get the result
func RunCommand(c string) ([]byte, error) {
	// cmd := exec.Command(c)
	// return cmd.CombinedOutput()

	logs.Debug("Running command: %s", c)

	return nil, nil
}
