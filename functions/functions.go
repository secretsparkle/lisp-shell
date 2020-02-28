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
	// need to put nuew functions in the symbol table and the functions table
	functions *map[string]structs.Function, bindings map[string]string) error {
	fmt.Println("ExecFunction")
	fmt.Println(args)
	switch args[0] {
	case "defun":
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
	case "cons":
	case "+":
		if len(args) == 1 {
			return errors.New("Invalid number of arguments.")
		}
		sum := 0.0
		for _, number := range args[1:] {
			if bindings != nil {
				number = bindings[number]
			}
			if n, err := strconv.ParseFloat(number, 64); err == nil {
				sum += n
			} else {
				return errors.New("Only numbers can be added.")
			}
		}
		fmt.Println(sum)
		return nil
	case "-":
		if len(args) == 1 {
			return errors.New("Invalid number of arguments.")
		} else if len(args) == 2 {
			difference, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return errors.New("Only numbers can be subtracted.")
			}
			difference = 0 - difference
			fmt.Println(difference)
			return nil
		} else {
			difference, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return errors.New("Only numbers can be subtracted.")
			}
			for _, number := range args[2:] {
				if n, err := strconv.ParseFloat(number, 64); err == nil {
					difference -= n
				} else {
					return errors.New("Only numbers can be subtracted.")
				}
			}
			fmt.Println(difference)
			return nil
		}
	case "*":
		if len(args) == 1 {
			return errors.New("Invalid number of arguments.")
		}
		product := 1.0
		for _, number := range args[1:] {
			if n, err := strconv.ParseFloat(number, 64); err == nil {
				product *= n
			} else {
				return errors.New("Only numbers can be multiplied.")
			}
		}
		fmt.Println(product)
		return nil
	case "/":
		if len(args) < 3 {
			return errors.New("Invalid number of arguments.")
		}
		numerator, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return errors.New("Only numbers can be divided.")
		}
		for _, number := range args[2:] {
			if n, err := strconv.ParseFloat(number, 64); err == nil {
				numerator /= n
			} else {
				return errors.New("Only numbers can be divided.")
			}
		}
		fmt.Println(numerator)
		return nil
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
	}

	// Pass the program and the arguments separately
	cmd := exec.Command(args[0], args[1:]...)

	//Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command
	return cmd.Run()
}
