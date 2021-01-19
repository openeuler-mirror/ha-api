package models

func GetCommandsList() map[string]interface{} {
	var result map[string]interface{}

	commandTable := map[int]string{
		1: "crm_mon -1 -o",
		2: "crm_simulate -Ls",
		3: "pcs config show",
		4: "corosync-cfgtool -s",
	}
	result["action"] = true
	result["data"] = commandTable
	return result
}
