package main

import (
	"./conditionals"
	"./parse"
	"./structs"
	"./utils"
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	// should this be global?
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
		sep := []rune("() ")
		args := util.SplitWith(input, sep)
		args = util.RemoveMember(args, " ")
		args = args[1:] //shave off that opening paren
		s_expressions, _, err = parse.Transliterate(s_expressions, args, 0)

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
