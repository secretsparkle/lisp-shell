package parse

import (
	"../structs"
	"errors"
	"regexp"
)

func Parse(input string) error {
	tokens := tokenize(input)
	if !sExpression(tokens) {
		return errors.New("Invalid expression")
	}
	return nil
}

func tokenize(input string) []string {
	var tokens []string
	var atom string

	for range input {
		switch input[0] {
		case '\n':
			break
		case '(':
			if len(atom) > 0 {
				tokens = append(tokens, atom)
				atom = ""
				input = input[1:]
				continue
			}
			tokens = append(tokens, "(")
			input = input[1:]
		case ')':
			if len(atom) > 0 {
				tokens = append(tokens, atom)
				atom = ""
			}
			tokens = append(tokens, ")")
			input = input[1:]
		case ' ':
			if len(atom) > 0 {
				tokens = append(tokens, atom)
				atom = ""
			}
			input = input[1:]
		default:
			atom += string(input[0])
			input = input[1:]
		}
	}
	return tokens
}

func sExpression(tokens []string) bool {
	if atom(tokens[0]) || list(tokens) {
		return true
	}
	return false
}

func list(tokens []string) bool {
	var stack []string
	if tokens[0] != "(" {
		return false
	}
	stack = append(stack, "(")
	tokens = tokens[1:]
	for range tokens {
		if tokens[0] == ")" && len(stack) > 0 {
			stack = stack[1:]
		} else if tokens[0] == ")" && len(stack) == 0 {
			break
		} else if !sExpression(tokens) {
			return false
		}
		tokens = tokens[1:]
	}
	return true
}

func atom(token string) bool {
	var IsAtom = regexp.MustCompile(`^[[:graph:]]+$`).MatchString
	if !IsAtom(token) || token == "(" || token == ")" {
		return false
	}
	return true
}

func Transliterate(list structs.List, args []string, index int) (structs.List, int, error) {
	for {
		token := args[index]
		index++
		if token == "(" {
			var newList structs.List
			if subList, index, err := Transliterate(newList, args[index:], index); err == nil {
				list.PushBack(subList)
			} else {
				return subList, index, err
			}
		} else if token == ")" {
			return list, index, nil
		} else {
			list.PushBack(token)
		}
	}
}

/*
func transliterate(list structs.List, args []string, openParen int) (int, structs.List) {
	expressionCount := 0
	catchUpIndex := 0
	currIndex := 0

	for index, arg := range args {
		currIndex = index
		if catchUpIndex > 0 {
			catchUpIndex--
			continue
		}
		if strings.Contains(arg, ")\n") {
			openParen--
			if strings.Contains(arg, "(") {
				arg = strings.Trim(arg, "()\n")
				list.PushBack(arg)
			} else {
				arg = strings.Trim(arg, ")\n")
				list.PushBack(arg)
			}
			if openParen > 1 {
				continue
			} else {
				break
			}
		} else if strings.Contains(arg, ")") {
			openParen--
			if strings.Contains(arg, "(") {
				arg = strings.Trim(arg, "()\n")
				list.PushBack(arg)
			} else {
				arg = strings.TrimRight(arg, ")")
				list.PushBack(arg)
			}
			if openParen > 1 {
				continue
			} else {
				break
			}
		} else if strings.Contains(arg, "'(") && expressionCount == 0 { // beginning
			list.PushBack("'")
			list.PushBack(arg[2:])
			expressionCount++
		} else if strings.Contains(arg, "(") && expressionCount == 0 { // beginning
			list.PushBack(arg[1:])
			expressionCount++
		} else if strings.Contains(arg, "'(") && expressionCount > 0 {
			var newIndex int
			var innerList structs.List
			list.PushBack("'")
			openParen++
			newIndex, innerList = Transliterate(innerList, args[index:], openParen)
			catchUpIndex = newIndex
			list.PushBack(innerList)
		} else if strings.Contains(arg, "(") && expressionCount > 0 {
			var newIndex int
			var innerList structs.List
			openParen++
			newIndex, innerList = Transliterate(innerList, args[index:], openParen)
			catchUpIndex = newIndex
			list.PushBack(innerList)
		} else {
			if strings.Contains(arg, ")") {
				openParen--
				list.PushBack(strings.TrimRight(arg, ")"))
				if openParen > 1 {
					continue
				} else {
					break
				}
			} else {
				list.PushBack(arg)
			}
		}
	}
	return currIndex, list
}
*/
