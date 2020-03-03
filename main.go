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
	//strings := make(map[string]string)
	//numbers := make(map[string]float64)
	//bools := make(map[string]bool)

	for {
		fmt.Print("> ")
		// Read the keyboard input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if err = parse.Parse(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		fmt.Println("Valid s-expression")
		// Handle the execution of the input.
		if err = execInput(input, &symbols, &functions); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execInput(input string, symbols *map[string]rune,
	functionTable *map[string]structs.Function) error {

	// Remove the leading parentheses
	input = strings.Replace(input, "(", " ", 1)

	// Remove the trailing parenthese and newline character.
	input = strings.TrimSuffix(input, ")\n")

	// Split the input to separate the command and the arguments.
	args := strings.Split(input, " ")

	// Check for built-in commands
	switch (*symbols)[args[0]] {
	case 'c':
		return nil
	case 'f':
		return functions.ExecFunction(args, symbols, functionTable, nil)
	default:
		return functions.ExecFunction(args, symbols, functionTable, nil)
	}
}
