package utils

import (
	"fmt"
	"testing"
)

func TestRunCommand(t *testing.T) {
	cmd := "echo \"hello world\""
	out, err := RunCommand(cmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(out))

	cmd2 := "echo \"hello world\" | grep -o world"
	out, err = RunCommand(cmd2)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(out))
}
