package parse

import (
	"../structs"
	"errors"
	"regexp"
	"strings"
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

func Transliterate(list structs.List, args []string, openParen int) (int, structs.List) {
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
			if strings.Contains(arg, "(") {
				arg = strings.Trim(arg, "()\n")
				list.PushBack(arg)
			} else {
				arg = strings.Trim(arg, ")\n")
				list.PushBack(arg)
			}
			break
		} else if strings.Contains(arg, ")") {
			if strings.Contains(arg, "(") {
				arg = strings.Trim(arg, "()\n")
				list.PushBack(arg)
			} else {
				arg = strings.TrimRight(arg, ")")
				list.PushBack(arg)
			}
			break
		} else if strings.Contains(arg, "'(") && expressionCount == 0 { // beginning
			list.PushBack("'")
			list.PushBack(arg[2:])
			openParen++
			expressionCount++
		} else if strings.Contains(arg, "(") && expressionCount == 0 { // beginning
			list.PushBack(arg[1:])
			openParen++
			expressionCount++
		} else if strings.Contains(arg, "'(") && expressionCount > 0 {
			var newIndex int
			var innerList structs.List
			list.PushBack("'")
			openParen++
			newIndex, innerList = Transliterate(innerList, args[index:], openParen)
			openParen--
			catchUpIndex = newIndex
			list.PushBack(innerList)
		} else if strings.Contains(arg, "(") && expressionCount > 0 {
			var newIndex int
			var innerList structs.List
			openParen++
			newIndex, innerList = Transliterate(innerList, args[index:], openParen)
			openParen--
			catchUpIndex = newIndex
			list.PushBack(innerList)
		} else {
			if strings.Contains(arg, ")") {
				list.PushBack(strings.TrimRight(arg, ")"))
			} else {
				list.PushBack(arg)
			}
		}
	}
	return currIndex, list
}
