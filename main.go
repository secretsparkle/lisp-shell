package main

import (
	"./functions"
	"./parse"
	"./structs"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	symbols := map[string]rune{
		"if":      'c',
		"cond":    'c',
		"defvar":  'f',
		"defun":   'f',
		"+":       'f',
		"-":       'f',
		"*":       'f',
		"/":       'f',
		"quote":   'f',
		"cons":    'f',
		"car":     'f',
		"cdr":     'f',
		"first":   'f',
		"rest":    'f',
		"last":    'f',
		"reverse": 'f',
	}

	functions := make(map[string]structs.Function)

	for {
		fmt.Print("> ")
		// Read the keyboard input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		var s_expressions structs.List
		var value interface{}
		args := strings.Split(input, " ")
		_, s_expressions = transliterate(s_expressions, args)

		if err = parse.Parse(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Handle the execution of the input.
		if value, err = execInput(s_expressions, &symbols, &functions); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(value)
	}
}

func transliterate(list structs.List, args []string) (int, structs.List) {
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
				arg = strings.Replace(arg, "(", "", -1)
				arg = strings.Replace(arg, ")", "", -1)
				arg = strings.Replace(arg, "\n", "", -1)
				list.PushBack(arg)
			} else {
				arg = strings.Replace(arg, ")\n", "", -1)
				list.PushBack(strings.Replace(arg, ")", "", -1))
			}
			break
		} else if strings.Contains(arg, ")") {
			if strings.Contains(arg, "(") {
				list.PushBack(strings.TrimRight(strings.TrimLeft(arg, "("), ")\n"))
			} else {
				list.PushBack(strings.TrimRight(arg, ")"))
			}
			break
		} else if strings.Contains(arg, "'(") && expressionCount == 0 { // beginning
			list.PushBack(arg[2:])
			expressionCount++
		} else if strings.Contains(arg, "(") && expressionCount == 0 { // beginning
			list.PushBack(arg[1:])
			expressionCount++
		} else if strings.Contains(arg, "'(") && expressionCount > 0 {
			var newIndex int
			var innerList structs.List
			newIndex, innerList = transliterate(innerList, args[index:])
			catchUpIndex = newIndex
			list.PushBack(innerList)
		} else if strings.Contains(arg, "(") && expressionCount > 0 {
			var newIndex int
			var innerList structs.List
			newIndex, innerList = transliterate(innerList, args[index:])
			catchUpIndex = newIndex
			list.PushBack(innerList)
		} else {
			list.PushBack(strings.TrimRight(arg, ")"))
		}
	}
	return currIndex, list
}

func execInput(expression structs.List, symbols *map[string]rune,
	functionTable *map[string]structs.Function) (interface{}, error) {

	//function := new.SExpression[0].(string)
	// Check for built-in commands
	switch expression.Head.Data {
	case 'c':
		return expression, nil
	case 'f':
		value, err := functions.ExecFunction(expression, symbols, functionTable, nil)
		return value, err
	default:
		value, err := functions.ExecFunction(expression, symbols, functionTable, nil)
		return value, err
	}
}
