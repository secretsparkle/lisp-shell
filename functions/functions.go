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

func ExecFunction(expression structs.List, symbols *map[string]rune,
	functions *map[string]structs.Function, bindings map[string]string) (interface{}, error) {
	switch expression.Head.Data {
	case "'":
		return list(expression)
	case "car":
		return car(expression, symbols, functions, bindings)
	case "cdr":
		return cdr(expression, symbols, functions, bindings)
	case "cons":
	case "defun":
		return defun(expression, symbols, functions)
	case "defvar":
		return defvar(expression, symbols, functions, bindings)
	case "first":
		return car(expression, symbols, functions, bindings)
	case "last":
		return last(expression, symbols, functions, bindings)
	case "list":
		return list(expression)
	case "quote":
	case "map":
	case "rest":
		return cdr(expression, symbols, functions, bindings)
	case "reverse":
	case "+":
		sum, err := plus(expression, symbols, functions, bindings)
		return sum, err
	case "-":
		difference, err := minus(expression, symbols, functions, bindings)
		return difference, err
	case "*":
		product, err := times(expression, symbols, functions, bindings)
		return product, err
	case "/":
		result, err := divide(expression, symbols, functions, bindings)
		return result, err
	case "cd":
		// 'cd' to home dir with empty path not yet supported.
		if expression.Len() < 2 {
			return expression, errors.New("path required")
		}
		var dir interface{}
		var err error
		e := expression.Head
		e = e.Next()
		switch e.Data.(type) {
		case string:
			dir = e.Data.(string)
		default:
			dir, err = ExecFunction(e.Data.(structs.List), symbols, functions, bindings)
			if err != nil {
				return 0.0, err
			}
		}
		// Change the directory and return the error.
		return nil, os.Chdir(dir.(string))
	case "echo":
		e := expression.Head
		e = e.Next()
		var out []string
		for ; e != nil; e = e.Next() {
			switch e.Data.(type) {
			case string:
				out = append(out, e.Data.(string))
			default:
				value, err := ExecFunction(e.Data.(structs.List), symbols, functions, bindings)
				if err != nil {
					return "", err
				}
				out = append(out, value.(string))
			}
		}
		var output string
		for _, val := range out {
			output += val
			output += " "
		}
		output = strings.TrimRight(output, " \n")
		return output, nil
	case "exit":
		os.Exit(0)
	default:
		// user defined function
		command := expression.Head.Data.(string)
		function := (*functions)[command]
		if function.Name != "" { // User defined function
			body := function.Body
			function.Bindings = make(map[string]string)
			e := expression.Head.Next()
			for _, arg := range function.Args {
				value := e.Data.(string)
				function.Bindings[arg] = value
				e = e.Next()
			}
			ExecFunction(body, symbols, functions, function.Bindings)
			return expression, nil
		} else { // UNIX command
			var statement []string
			statement = append(statement, command)
			e := expression.Head.Next()
			for ; e != nil; e = e.Next() {
				switch e.Data.(type) {
				case string:
					statement = append(statement, e.Data.(string))
				default:
					subValue, err := ExecFunction(e.Data.(structs.List), symbols, functions, bindings)
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

func car(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings map[string]string) (interface{}, error) {
	structs.PrintList(expression)
	// honestly will probably need to be reworked in the future
	e := expression.Head
	switch e.Next().Data.(type) {
	case string:
		return nil, errors.New("car requires a list")
	default:
		e = e.Next()
		l := e.Data.(structs.List)
		e = l.Head
		if e.Data == "list" {
			e = l.Head.Next()
		}
		return e.Data, nil
	}
}

func cdr(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings map[string]string) (interface{}, error) {
	e := expression.Head
	switch e.Next().Data.(type) {
	case string:
		return nil, errors.New("cdr requires a list")
	default:
		e = e.Next()
		l := e.Data.(structs.List)
		e = l.Head
		if e.Data == "list" {
			e = l.Head.Next()
		}
		e = e.Next()
		var rest []string
		for ; e != nil; e = e.Next() {
			rest = append(rest, e.Data.(string))
		}
		return rest, nil
	}
}

func defun(llat structs.List, symbols *map[string]rune,
	functions *map[string]structs.Function) (structs.List, error) {
	funct := new(structs.Function)
	a := llat.Head
	a = a.Next()
	funct.Name = a.Data.(string)
	a = a.Next()
	params(a.Data.(structs.List), funct, symbols, functions)
	a = a.Next()
	funct.Body = a.Data.(structs.List)

	(*symbols)[funct.Name] = 'f'
	(*functions)[funct.Name] = *funct

	fmt.Println("New Function Name: ", funct.Name)
	fmt.Println("New Function Args: ", funct.Args)
	fmt.Println("New Function Body: ")
	structs.PrintList(funct.Body)

	return llat, nil
}

func defvar(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings map[string]string) (interface{}, error) {
	if expression.Len() != 3 {
		return nil, errors.New("Invalid number of arguments supplied to defvar")
	}
	e := expression.Head
	e = e.Next()
	symbol := e.Data.(string)
	e = e.Next()
	value := e.Data.(string)
	bindings[symbol] = value

	return strings.ToUpper(symbol), nil
}

func last(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings map[string]string) (interface{}, error) {
	e := expression.Head
	switch e.Next().Data.(type) {
	case string:
		return nil, errors.New("last requires a list")
	default:
		e = e.Next()
		l := e.Data.(structs.List)
		e = l.Tail
		return e.Data, nil
	}
}

// will need to add in symbols, functions and bindings later
func list(expression structs.List) (structs.List, error) {
	var newList structs.List
	newList = *newList.PushBack("(")
	for a := expression.Head; a != nil; a = a.Next() {
		if a.Data == "list" {
			continue
		}
		switch a.Data.(type) {
		case string:
			if a.Next() == nil {
				newList = *newList.PushBack(")")
			} else {
				newList = *newList.PushBack(a.Data)
				newList = *newList.PushBack(" ")
			}

		default:
			if subList, err := list(a.Data.(structs.List)); err == nil {
				newList = *newList.PushBack(subList)
			} else {
				return subList, err
			}
		}
	}
	newList = *newList.PushBack("\n")
	return newList, nil
}

func params(lat structs.List, funct *structs.Function, symbols *map[string]rune,
	functions *map[string]structs.Function) error {
	for a := lat.Head; a != nil; a = a.Next() {
		funct.Args = append(funct.Args, a.Data.(string))
	}

	return nil
}

func plus(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings map[string]string) (float64, error) {
	if expression.Len() == 1 {
		return 0.0, errors.New("Invalid number of arguments.")
	}
	sum := 0.0
	e := expression.Head
	for e = e.Next(); e != nil; e = e.Next() {
		switch e.Data.(type) {
		case string:
			number := e.Data
			if bindings != nil {
				number = bindings[number.(string)]
			}
			if num, err := strconv.ParseFloat(number.(string), 64); err == nil {
				sum += num
			} else {
				return 0.0, errors.New("Only numbers can be added.")
			}
		default:
			subValue, err := ExecFunction(e.Data.(structs.List), symbols, functions, bindings)
			if err != nil {
				return 0.0, err
			}
			sum += subValue.(float64)
		}
	}
	return sum, nil
}

func minus(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings map[string]string) (interface{}, error) {
	var number string
	e := expression.Head
	e = e.Next()

	if expression.Len() == 1 {
		return expression, errors.New("Invalid number of arguments.")
	} else if expression.Len() == 2 {
		if bindings != nil {
			number = bindings[e.Data.(string)]
		}
		difference, err := strconv.ParseFloat(number, 64)
		if err != nil {
			return expression, errors.New("Only numbers can be subtracted.")
		}
		difference = 0 - difference
		return difference, nil
	} else {
		if bindings != nil {
			number = bindings[e.Data.(string)]
		}
		difference, err := strconv.ParseFloat(e.Data.(string), 64)
		if err != nil {
			return expression, errors.New("Only numbers can be subtracted.")
		}
		for e = e.Next(); e != nil; e = e.Next() {
			switch e.Data.(type) {
			case string:
				if bindings != nil {
					number = bindings[e.Data.(string)]
					if num, err := strconv.ParseFloat(number, 64); err == nil {
						difference -= num
					} else {
						return expression, errors.New("Only numbers can be subtracted.")
					}
				} else {
					if num, err := strconv.ParseFloat(e.Data.(string), 64); err == nil {
						difference -= num
					} else {
						return expression, errors.New("Only numbers can be subtracted.")
					}
				}
			default:
				subValue, err := ExecFunction(e.Data.(structs.List), symbols, functions, bindings)
				if err != nil {
					return 0.0, err
				}
				difference -= subValue.(float64)

			}
		}
		return difference, nil
	}
}

func times(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings map[string]string) (interface{}, error) {
	if expression.Len() == 1 {
		return expression, errors.New("Invalid number of arguments.")
	}
	product := 1.0
	var number string
	e := expression.Head
	for e = e.Next(); e != nil; e = e.Next() {
		switch e.Data.(type) {
		case string:
			if bindings != nil {
				number = bindings[e.Data.(string)]
				if num, err := strconv.ParseFloat(number, 64); err == nil {
					product *= num
				} else {
					return expression, errors.New("Only numbers can be multiplied.")
				}
			} else {
				if num, err := strconv.ParseFloat(e.Data.(string), 64); err == nil {
					product *= num
				} else {
					return expression, errors.New("Only numbers can be multiplied.")
				}
			}
		default:
			subValue, err := ExecFunction(e.Data.(structs.List), symbols, functions, bindings)
			if err != nil {
				return 0.0, err
			}
			product *= subValue.(float64)
		}
	}
	return product, nil
}

func divide(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings map[string]string) (interface{}, error) {
	var numer, numerator float64
	var err error
	e := expression.Head
	e = e.Next()

	if expression.Len() < 3 {
		return expression, errors.New("Invalid number of arguments.")
	}
	if bindings != nil {
		numer, err = strconv.ParseFloat(bindings[e.Data.(string)], 64)
		if err != nil {
			return expression, errors.New("Only numbers can be divided.")
		}
	} else {
		numer, err = strconv.ParseFloat(e.Data.(string), 64)
		if err != nil {
			return expression, errors.New("Only numbers can be divided.")
		}
	}
	numerator = numer
	for e = e.Next(); e != nil; e = e.Next() {
		switch e.Data.(type) {
		case string:
			if bindings != nil {
				number := bindings[e.Data.(string)]
				if num, err := strconv.ParseFloat(number, 64); err == nil {
					numerator /= num
				} else {
					return expression, errors.New("Only numbers can be divided.")
				}
			} else {
				if num, err := strconv.ParseFloat(e.Data.(string), 64); err == nil {
					numerator /= num
				} else {
					return expression, errors.New("Only numbers can be divided.")
				}
			}
		default:
			subValue, err := ExecFunction(e.Data.(structs.List), symbols, functions, bindings)
			if err != nil {
				return 0.0, err
			}
			numerator /= subValue.(float64)
		}
	}
	return numerator, nil
}
