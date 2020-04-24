package evaluation

import (
	"../structs"
	"fmt"
	"os"
	"os/exec"
)

func EvaluateFunction(expression structs.List, functionList *map[string]rune,
	functions *map[string]structs.Function, bindings *map[string]interface{}) (interface{}, error) {
	switch expression.Head.Data {
	case "'":
		return list(expression)
	case "and":
		return and(expression, functionList, functions, bindings)
	case "car":
		return car(expression, functionList, functions, bindings)
	case "cdr":
		return cdr(expression, functionList, functions, bindings)
	case "cond":
	case "cons":
		return cons(expression, functionList, functions, bindings)
	case "define":
		return define(expression, functionList, functions, bindings)
		/*
			case "defun":
				return defun(expression, symbols, functions)
			case "defvar": // redo for nesting
				return defvar(expression, symbols, functions, bindings)
		*/
	case "equal":
		return equal(expression, functionList, functions, bindings)
	case "first":
		return car(expression, functionList, functions, bindings)
	case "if":
		return if_statement(expression, functionList, functions, bindings)
	case "last":
		return last(expression, functionList, functions, bindings)
	case "list":
		return list(expression)
	case "quote":
	case "map":
	case "rest":
		return cdr(expression, functionList, functions, bindings)
	case "reverse": // redo for nesting
		return reverse(expression, functionList, functions, bindings)
	case "=":
		return equal(expression, functionList, functions, bindings)
	case "+":
		sum, err := plus(expression, functionList, functions, bindings)
		return sum, err
	case "-":
		difference, err := minus(expression, functionList, functions, bindings)
		return difference, err
	case "*":
		product, err := times(expression, functionList, functions, bindings)
		return product, err
	case "/":
		result, err := divide(expression, functionList, functions, bindings)
		return result, err
	case ">":
		result, err := gt_or_lt(expression, functionList, functions, bindings, ">")
		return result, err
	case "<":
		result, err := gt_or_lt(expression, functionList, functions, bindings, "<")
		return result, err
	case "cd":
		return cd(expression, functionList, functions, bindings)
	case "echo":
	case "exit":
	default:
		// user defined function
		command := expression.Head.Data.(string)
		function := (*functions)[command]
		if function.Name != "" { // User defined function
			body := function.Body
			function.Bindings = make(map[string]interface{})
			e := expression.Head.Next()
			for _, arg := range function.Args {
				value := e.Data.(string)
				function.Bindings[arg] = value
				e = e.Next()
			}
			retVal, err := EvaluateFunction(body, functionList, functions, &function.Bindings)
			if err != nil {
				return nil, err
			}
			return retVal, nil
		} else { // UNIX command
			var statement []string
			statement = append(statement, command)
			e := expression.Head.Next()
			for ; e != nil; e = e.Next() {
				switch e.Data.(type) {
				case string:
					statement = append(statement, e.Data.(string))
				default:
					subValue, err := EvaluateFunction(e.Data.(structs.List), functionList, functions, bindings)
					fmt.Println("UNIX return: ", subValue)
					if err != nil {
						return 0.0, err
					}
					statement = append(statement, subValue.(string))
				}
			}
			// Pass the program and the arguments separately
			cmd := exec.Command(statement[0], statement[1:]...)

			//Set the correct output device.
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			// Execute the command
			return cmd.Run(), nil
		}
	}
	return nil, nil
}
