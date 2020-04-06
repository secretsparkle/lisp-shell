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
		_, s_expressions = parse.Transliterate(s_expressions, args, 0)

		if err = parse.Parse(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Handle the execution of the input.
		if value, err = conditionals.ExecInput(s_expressions, &symbols, &functions, &bindings); err != nil {
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
