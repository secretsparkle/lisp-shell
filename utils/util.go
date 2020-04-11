package util

import "unicode"

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

func IsAlphabetic(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsNumber(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func AnySymbol(s string) bool {
	for _, r := range s {
		if unicode.IsSymbol(r) {
			return true
		}
	}
	return false
}
