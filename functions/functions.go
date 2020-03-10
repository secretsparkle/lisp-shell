package functions

import (
	"../structs"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func ExecFunction(new structs.SExpression, symbols *map[string]rune,
	functions *map[string]structs.Function, bindings map[string]string) error {
	fmt.Println("ExecFunction")
	fmt.Println(new.SExpression)
	switch new.SExpression[0] {
	case "'":
		return list(new)
	case "car":
	case "cdr":
	case "cons":
	case "defun":
		return defun(new, symbols, functions)
	case "defvar":
	case "first":
	case "last":
	case "list":
		return list(new)
	case "quote":
	case "rest":
	case "reverse":
	case "+":
		return plus(new, bindings)
	case "-":
		return minus(new, bindings)
	case "*":
		return times(new, bindings)
	case "/":
		return divide(new, bindings)
	case "cd":
		// 'cd' to home dir with empty path not yet supported.
		if len(new.SExpression) < 2 {
			return errors.New("path required")
		}
		dir := new.SExpression[1].(string)
		// Change the directory and return the error.
		return os.Chdir(dir)
	case "exit":
		os.Exit(0)
	default:
		command := new.SExpression[0].(string)
		function := (*functions)[command]
		if function.Name != "" { // User defined function
			fmt.Println("Here!")
			body := function.Body
			frstExpression := body.SExpression[0].(string)
			lstExpression := body.SExpression[len(body.SExpression)-1].(string)
			function.Bindings = make(map[string]string)
			body.SExpression[0] = strings.TrimPrefix(frstExpression, "(")
			body.SExpression[len(body.SExpression)-1] = strings.TrimSuffix(lstExpression, ")")
			index := 1
			for _, arg := range function.Args {
				value := new.SExpression[index].(string)
				function.Bindings[arg] = strings.Trim(value, "()")
				index++
			}
			ExecFunction(body, symbols, functions, function.Bindings)
			return nil
		} else { // UNIX command
			var command []string
			for _, arg := range new.SExpression {
				command = append(command, arg.(string))
			}
			// Pass the program and the arguments separately
			cmd := exec.Command(command[0], command[1:]...)

			//Set the correct output device.
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			// Execute the command
			return cmd.Run()
		}
	}
	return nil
}

func list(args structs.SExpression) error {
	newList := "("

	for index, token := range args.SExpression {
		if strings.Contains("list", token.(string)) || strings.Contains("'", token.(string)) {
			continue
		} else if index == len(args.SExpression)-1 {
			newList += token.(string)
		} else if strings.Contains("(", token.(string)) {

		} else {
			newList += token.(string)
			newList += " "
		}
	}
	newList += ")"
	fmt.Println(newList)
	return nil
}

func defun(args structs.SExpression, symbols *map[string]rune,
	functions *map[string]structs.Function) error {
	newFunction := new(structs.Function)
	newFunction.Name = args.SExpression[1].(string)
	fmt.Println(newFunction.Name)
	i := 2
	for {
		fmt.Println("args[i]: " + args.SExpression[i].(string))
		if strings.Contains(args.SExpression[i].(string), "(") {
			args.SExpression[i] = strings.Trim(args.SExpression[i].(string), "(")
		}
		if strings.Contains(args.SExpression[i].(string), ")") {
			args.SExpression[i] = strings.Trim(args.SExpression[i].(string), ")")
			newFunction.Args = append(newFunction.Args, args.SExpression[i].(string))
			fmt.Println(newFunction.Args)
			i++
			break
		}
		newFunction.Args = append(newFunction.Args, args.SExpression[i].(string))
		i++
	}
	for ; i < len(args.SExpression); i++ {
		newFunction.Body.SExpression = append(newFunction.Body.SExpression, args.SExpression[i])
		fmt.Println(newFunction.Body)
	}
	(*symbols)[newFunction.Name] = 'f'
	(*functions)[newFunction.Name] = *newFunction
	return nil
}

func plus(args structs.SExpression, bindings map[string]string) error {
	if len(args.SExpression) == 1 {
		return errors.New("Invalid number of arguments.")
	}
	sum := 0.0
	for _, number := range args.SExpression[1:] {
		if bindings != nil {
			number = bindings[number.(string)]
		}
		if num, err := strconv.ParseFloat(number.(string), 64); err == nil {
			sum += num
		} else {
			return errors.New("Only numbers can be added.")
		}
	}
	fmt.Println(sum)
	return nil
}

func minus(args structs.SExpression, bindings map[string]string) error {
	var number string

	if len(args.SExpression) == 1 {
		return errors.New("Invalid number of arguments.")
	} else if len(args.SExpression) == 2 {
		if bindings != nil {
			number = bindings[args.SExpression[1].(string)]
		}
		difference, err := strconv.ParseFloat(number, 64)
		if err != nil {
			return errors.New("Only numbers can be subtracted.")
		}
		difference = 0 - difference
		fmt.Println(difference)
		return nil
	} else {
		if bindings != nil {
			number = bindings[args.SExpression[1].(string)]
		}
		difference, err := strconv.ParseFloat(args.SExpression[1].(string), 64)
		if err != nil {
			return errors.New("Only numbers can be subtracted.")
		}
		for _, value := range args.SExpression[2:] {
			if bindings != nil {
				number = bindings[value.(string)]
				if num, err := strconv.ParseFloat(number, 64); err == nil {
					difference -= num
				} else {
					return errors.New("Only numbers can be subtracted.")
				}
			} else {
				if num, err := strconv.ParseFloat(value.(string), 64); err == nil {
					difference -= num
				} else {
					return errors.New("Only numbers can be subtracted.")
				}
			}
		}
		fmt.Println(difference)
		return nil
	}
}

func times(args structs.SExpression, bindings map[string]string) error {
	if len(args.SExpression) == 1 {
		return errors.New("Invalid number of arguments.")
	}
	product := 1.0
	var number string
	for _, value := range args.SExpression[1:] {
		if bindings != nil {
			number = bindings[value.(string)]
			if num, err := strconv.ParseFloat(number, 64); err == nil {
				product *= num
			} else {
				return errors.New("Only numbers can be multiplied.")
			}
		} else {
			if num, err := strconv.ParseFloat(value.(string), 64); err == nil {
				product *= num
			} else {
				return errors.New("Only numbers can be multiplied.")
			}
		}
	}
	fmt.Println(product)
	return nil
}

func divide(args structs.SExpression, bindings map[string]string) error {
	var numer, numerator float64
	var err error

	if len(args.SExpression) < 3 {
		return errors.New("Invalid number of arguments.")
	}
	if bindings != nil {
		numer, err = strconv.ParseFloat(bindings[args.SExpression[1].(string)], 64)
		if err != nil {
			return errors.New("Only numbers can be divided.")
		}
	} else {
		numer, err = strconv.ParseFloat(args.SExpression[1].(string), 64)
		if err != nil {
			return errors.New("Only numbers can be divided.")
		}
	}
	numerator = numer
	for _, number := range args.SExpression[2:] {
		if bindings != nil {
			number = bindings[number.(string)]
			if num, err := strconv.ParseFloat(number.(string), 64); err == nil {
				numerator /= num
			} else {
				return errors.New("Only numbers can be divided.")
			}
		} else {
			if num, err := strconv.ParseFloat(number.(string), 64); err == nil {
				numerator /= num
			} else {
				return errors.New("Only numbers can be divided.")
			}
		}
	}
	fmt.Println(numerator)
	return nil
}
