package parse

import (
	"../structs"
	"errors"
	"regexp"
)

// everything but transliterate needs to be reworked
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
			if subList, newIndex, err := Transliterate(newList, args[index:], 0); err == nil {
				list.PushBack(subList)
				index += newIndex
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
