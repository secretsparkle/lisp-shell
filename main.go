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

		var new structs.SExpression
		args := strings.Split(input, " ")
		_, new = convertInput(new, args)

		if err = parse.Parse(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		fmt.Println("Valid s-expression")
		// Handle the execution of the input.
		if new, err = execInput(new, &symbols, &functions); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func convertInput(new structs.SExpression, args []string) (int, structs.SExpression) {
	expressionCounter := 0
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
				new.SExpression = append(new.SExpression, strings.TrimRight(strings.TrimLeft(arg, "("), ")\n"))
			} else {
				new.SExpression = append(new.SExpression, strings.TrimRight(arg, ")"))
			}
			break
		} else if strings.Contains(arg, ")") {
			if strings.Contains(arg, "(") {
				new.SExpression = append(new.SExpression, strings.TrimRight(strings.TrimLeft(arg, "("), ")\n"))
			} else {
				new.SExpression = append(new.SExpression, strings.TrimRight(arg, ")"))
			}
			break
		} else if strings.Contains(arg, "'(") && expressionCounter == 0 { // beginning of an s-expression
			new.Data = true
			new.SExpression = append(new.SExpression, arg[2:])
			expressionCounter++
		} else if strings.Contains(arg, "(") && expressionCounter == 0 { // beginning of an s-expression
			new.SExpression = append(new.SExpression, arg[1:])
			expressionCounter++
		} else if strings.Contains(arg, "'(") && expressionCounter > 0 {
			var newIndex int
			var inner structs.SExpression
			inner.Data = true
			newIndex, inner = convertInput(inner, args[index:])
			catchUpIndex = newIndex
			new.SExpression = append(new.SExpression, inner.SExpression)
		} else if strings.Contains(arg, "(") && expressionCounter > 0 {
			var newIndex int
			var inner structs.SExpression
			newIndex, inner = convertInput(inner, args[index:])
			catchUpIndex = newIndex
			new.SExpression = append(new.SExpression, inner.SExpression)
		} else {
			new.SExpression = append(new.SExpression, arg)
		}
	}
	fmt.Println("Settled Input: ", new.SExpression)
	return currIndex, new
}

func execInput(new structs.SExpression, symbols *map[string]rune,
	functionTable *map[string]structs.Function) (structs.SExpression, error) {

	//function := new.SExpression[0].(string)
	// Check for built-in commands
	switch (*symbols)[new.SExpression[0].(string)] {
	case 'c':
		return new, nil
	case 'f':
		return functions.ExecFunction(new, symbols, functionTable, nil)
	default:
		return functions.ExecFunction(new, symbols, functionTable, nil)
	}
}
