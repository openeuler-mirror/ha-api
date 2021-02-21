package utils

import "regexp"

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

// GetNumAndUnitFromStr gets the first number and the unit after this number
// like "20.5min" ==> ["20.5", "min"]
func GetNumAndUnitFromStr(s string) (string, string) {
	r := regexp.MustCompile("[0-9](.*)[0-9]")
	index := r.FindStringIndex(s)
	if len(index) == 0 {
		return s[:1], s[1:]
	}
	return s[:index[1]], s[index[1]:]
}
