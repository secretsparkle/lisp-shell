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

func ExecFunction(args []string, symbols *map[string]rune,
	functions *map[string]structs.Function, bindings map[string]string) error {
	fmt.Println("ExecFunction")
	fmt.Println(args)
	switch args[0] {
	case "'":
		return list(args)
	case "car":
	case "cdr":
	case "cons":
	case "defun":
		return defun(args, symbols, functions)
	case "defvar":
	case "first":
	case "last":
	case "list":
		return list(args)
	case "quote":
	case "rest":
	case "reverse":
	case "+":
		return plus(args, bindings)
	case "-":
		return minus(args, bindings)
	case "*":
		return times(args, bindings)
	case "/":
		return divide(args, bindings)
	case "cd":
		// 'cd' to home dir with empty path not yet supported.
		if len(args) < 2 {
			return errors.New("path required")
		}
		// Change the directory and return the error.
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	default:
		function := (*functions)[args[0]]
		if function.Name != "" { // User defined function
			fmt.Println("Here!")
			body := function.Body
			function.Bindings = make(map[string]string)
			body[0] = strings.TrimPrefix(body[0], "(")
			body[len(body)-1] = strings.TrimSuffix(body[len(body)-1], ")")
			index := 1
			for _, arg := range function.Args {
				function.Bindings[arg] = strings.Trim(args[index], "()")
				index++
			}
			ExecFunction(body, symbols, functions, function.Bindings)
			return nil
		} else { // UNIX command
			// Pass the program and the arguments separately
			cmd := exec.Command(args[0], args[1:]...)

			//Set the correct output device.
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			// Execute the command
			return cmd.Run()
		}
	}
	return nil
}

func list(args []string) error {
	newList := "("

	for index, token := range args {
		if strings.Contains("list", token) || strings.Contains("'", token) {
			continue
		} else if index == len(args)-1 {
			newList += token
		} else if strings.Contains("(", token) {

		} else {
			newList += token
			newList += " "
		}
	}
	newList += ")"
	fmt.Println(newList)
	return nil
}

func defun(args []string, symbols *map[string]rune,
	functions *map[string]structs.Function) error {
	newFunction := new(structs.Function)
	newFunction.Name = args[1]
	fmt.Println(newFunction.Name)
	i := 2
	for {
		fmt.Println("args[i]: " + args[i])
		if strings.Contains(args[i], "(") {
			args[i] = strings.Trim(args[i], "(")
		}
		if strings.Contains(args[i], ")") {
			args[i] = strings.Trim(args[i], ")")
			newFunction.Args = append(newFunction.Args, args[i])
			fmt.Println(newFunction.Args)
			i++
			break
		}
		newFunction.Args = append(newFunction.Args, args[i])
		i++
	}
	for ; i < len(args); i++ {
		newFunction.Body = append(newFunction.Body, args[i])
		fmt.Println(newFunction.Body)
	}
	(*symbols)[newFunction.Name] = 'f'
	(*functions)[newFunction.Name] = *newFunction
	return nil
}

func plus(args []string, bindings map[string]string) error {
	if len(args) == 1 {
		return errors.New("Invalid number of arguments.")
	}
	sum := 0.0
	for _, number := range args[1:] {
		if bindings != nil {
			number = bindings[number]
		}
		if num, err := strconv.ParseFloat(number, 64); err == nil {
			sum += num
		} else {
			return errors.New("Only numbers can be added.")
		}
	}
	fmt.Println(sum)
	return nil
}

func minus(args []string, bindings map[string]string) error {
	var number string

	if len(args) == 1 {
		return errors.New("Invalid number of arguments.")
	} else if len(args) == 2 {
		if bindings != nil {
			number = bindings[args[1]]
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
			number = bindings[args[1]]
		}
		difference, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return errors.New("Only numbers can be subtracted.")
		}
		for _, value := range args[2:] {
			if bindings != nil {
				number = bindings[value]
				if num, err := strconv.ParseFloat(number, 64); err == nil {
					difference -= num
				} else {
					return errors.New("Only numbers can be subtracted.")
				}
			} else {
				if num, err := strconv.ParseFloat(value, 64); err == nil {
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

func times(args []string, bindings map[string]string) error {
	if len(args) == 1 {
		return errors.New("Invalid number of arguments.")
	}
	product := 1.0
	var number string
	for _, value := range args[1:] {
		if bindings != nil {
			number = bindings[value]
			if num, err := strconv.ParseFloat(number, 64); err == nil {
				product *= num
			} else {
				return errors.New("Only numbers can be multiplied.")
			}
		} else {
			if num, err := strconv.ParseFloat(value, 64); err == nil {
				product *= num
			} else {
				return errors.New("Only numbers can be multiplied.")
			}
		}
	}
	fmt.Println(product)
	return nil
}

func divide(args []string, bindings map[string]string) error {
	var numer, numerator float64
	var err error

	if len(args) < 3 {
		return errors.New("Invalid number of arguments.")
	}
	if bindings != nil {
		numer, err = strconv.ParseFloat(bindings[args[1]], 64)
		if err != nil {
			return errors.New("Only numbers can be divided.")
		}
	} else {
		numer, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			return errors.New("Only numbers can be divided.")
		}
	}
	numerator = numer
	for _, number := range args[2:] {
		if bindings != nil {
			number = bindings[number]
			if num, err := strconv.ParseFloat(number, 64); err == nil {
				numerator /= num
			} else {
				return errors.New("Only numbers can be divided.")
			}
		} else {
			if num, err := strconv.ParseFloat(number, 64); err == nil {
				numerator /= num
			} else {
				return errors.New("Only numbers can be divided.")
			}
		}
	}
	fmt.Println(numerator)
	return nil
}
