package utils

import (
	"errors"
	"runtime"

	"github.com/msteinert/pam"
)

func CheckAuth(username, password string) bool {
	t, err := pam.StartFunc("login", username, func(s pam.Style, msg string) (string, error) {
		switch s {
		case pam.PromptEchoOff:
			return password, nil
		case pam.PromptEchoOn:
			return "", nil
		case pam.ErrorMsg:
			return "", nil
		case pam.TextInfo:
			return "", nil
		}
		return "", errors.New("Unrecognized message style")
	})
	defer func() {
		runtime.GC()
	}()

	if err != nil {
		logs.Error("pam start error:", err)
		return false
	}

	err = t.Authenticate(0)
	if err != nil {
		logs.Error("pam auth error:", err)
		return false
	}

	return true
}
