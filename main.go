package main

import (
	"./conditionals"
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
		_, s_expressions = parse.Transliterate(s_expressions, args)

		if err = parse.Parse(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Handle the execution of the input.
		if value, err = conditionals.ExecInput(s_expressions, &symbols, &functions, &bindings); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if value != nil {
			switch value.(type) {
			case float64:
				fmt.Println(value)
			case string:
				fmt.Println(value.(string))
			default:
				fmt.Print("(")
				structs.PrintList(value.(structs.List))
				fmt.Println(")")
			}
		}
	}
}
