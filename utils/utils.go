package utils

func IsInSlice(str string, sli []string) bool {
	//TODO

	for _, item := range sli {
		if item == str {
			return true
		}
	}
	return false
}

// RemoveDupl remove duplicates in string array
func RemoveDupl(strs []string) []string {
	strSet := map[string]bool{}
	for _, v := range strs {
		strSet[v] = true
	}
	strsDupl := []string{}
	for k := range strSet {
		strsDupl = append(strsDupl, k)
	}
	return strsDupl
}
