package util

// split but keep seperators
func SplitWith(str string, sep []rune) []string {
	var strArr []string
	match := false
	subStr := ""
	for _, c := range str {
		for _, r := range sep {
			match = false
			if c == r {
				match = true
				if subStr != "" {
					strArr = append(strArr, subStr)
				}
				char := string(c)
				strArr = append(strArr, char)
				subStr = ""
				break
			}

		}
		if match == false {
			char := string(c)
			subStr += char
		}
	}
	return strArr
}

func RemoveMember(list []string, remove string) []string {
	var newList []string
	for _, member := range list {
		if member != remove {
			newList = append(newList, member)
		}
	}
	return newList
}

func CheckForSymbol(bindings *map[string]string, value string) string {
	if str := (*bindings)[value]; str == "" {
		return value
	} else {
		return str
	}
}
