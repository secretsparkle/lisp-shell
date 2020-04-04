package main

import (
	"./conditionals"
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
		"car":     'f',
		"cdr":     'f',
		"cond":    'c',
		"cons":    'f',
		"defun":   'f',
		"defvar":  'f',
		"equal":   'f',
		"first":   'f',
		"if":      'c',
		"last":    'f',
		"map":     'f',
		"quote":   'f',
		"rest":    'f',
		"reverse": 'f',
		"=":       'f',
		"+":       'f',
		"-":       'f',
		"*":       'f',
		"/":       'f',
    "<":       'f',
    ">":       'f',
	}

	functions := make(map[string]structs.Function)
	bindings := make(map[string]string)

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
		if value, err = execInput(s_expressions, &symbols, &functions, &bindings); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if value != nil {
			switch value.(type) {
			case bool:
				if value == true {
					fmt.Println("T")
				} else {
					fmt.Println("NIL")
				}
			case float64:
				fmt.Println(value)
			case string:
				fmt.Println(value.(string))
			case structs.List:
				fmt.Print("(")
				structs.PrintList(value.(structs.List))
				fmt.Println(")")
			}
		}
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
			expressionCount++
		} else if strings.Contains(arg, "(") && expressionCount == 0 { // beginning
			list.PushBack(arg[1:])
			expressionCount++
		} else if strings.Contains(arg, "'(") && expressionCount > 0 {
			var newIndex int
			var innerList structs.List
			list.PushBack("'")
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

func ExecInput(expression structs.List, symbols *map[string]rune,
	functionTable *map[string]structs.Function, bindings *map[string]string) (interface{}, error) {

	// Check for built-in commands
	switch (*symbols)[expression.Head.Data.(string)] {
	case 'c':
		value, err := conditionals.EvalConditional(expression, symbols, functionTable, bindings)
		return value, err
	case 'f':
		value, err := functions.ExecFunction(expression, symbols, functionTable, bindings)
		return value, err
	default:
		return expression, nil
	}
}
