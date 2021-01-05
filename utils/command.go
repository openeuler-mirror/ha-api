package utils

import (
	"os/exec"
)

// RunCommand runs the command and get the result
func RunCommand(c string, args ...string) ([]byte, error) {
	cmd := exec.Command(c, args...)
	return cmd.CombinedOutput()
}
