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
