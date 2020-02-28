package parse

import (
	"errors"
	"fmt"
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
		fmt.Println(input)
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
	fmt.Println("s-expression")
	fmt.Println(tokens)

	if atom(tokens[0]) || list(tokens) {
		return true
	}
	return false
}

func list(tokens []string) bool {
	fmt.Println("list")
	fmt.Println(tokens)
	var stack []string

	if tokens[0] != "(" {
		return false
	}
	stack = append(stack, "(")
	tokens = tokens[1:]
	for range tokens {
		fmt.Println(stack)
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
	fmt.Println("atom")
	fmt.Println(token)
	var IsAtom = regexp.MustCompile(`^[[:graph:]]+$`).MatchString
	if !IsAtom(token) || token == "(" || token == ")" {
		return false
	}
	return true
}
