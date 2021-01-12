package models

import (
	"openkylin.com/ha-api/utils"
)

func GenerateLog() map[string]interface{} {
	var result map[string]interface{}
	var file map[string]string

	var out []byte
	var err error
	if out, err = utils.RunCommand("/usr/share/heartbeat-gui/ha-api/loggen.sh"); err != nil {
		result["action"] = false
		result["error"] = "Get neokylinha log failed"
		return result
	}

	file_path := string(out)
	file["filepath"] = file_path
	result["action"] = true
	result["data"] = file
	return result
}
