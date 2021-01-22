package utils

import (
	"fmt"
	"testing"
)

func TestGetCorosyncConfig(t *testing.T) {
	result, err := getCorosyncConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}
