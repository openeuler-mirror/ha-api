package models

import (
	"errors"

	"openkylin.com/ha-api/utils"
)

// DO NOT EDIT ANYWHERE!!!
var commandTable = map[int]string{
	1: "crm_mon -1 -o",
	2: "crm_simulate -Ls",
	3: "pcs config show",
	4: "corosync-cfgtool -s",
}

func GetCommandsList() map[string]interface{} {
	var result map[string]interface{}

	result["action"] = true
	result["data"] = commandTable
	return result
}

func RunBuiltinCommand(cmdID int) (string, error) {
	if _, ok := commandTable[cmdID]; ok {
		cmd := commandTable[cmdID]
		output, err := utils.RunCommand(cmd)
		if err != nil {
			return "", errors.New("Run command failed: " + err.Error())
		}
		return string(output), nil
	}

	return "", errors.New("Invalid command index")
}
