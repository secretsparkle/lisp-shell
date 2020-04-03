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
	functions *map[string]structs.Function, bindings *map[string]string) (interface{}, error) {
	switch expression.Head.Data {
	case "'":
		return list(expression)
	case "car":
		return car(expression, symbols, functions, bindings)
	case "cdr":
		return cdr(expression, symbols, functions, bindings)
	case "cons":
		return cons(expression, symbols, functions, bindings)
	case "defun":
		return defun(expression, symbols, functions)
	case "defvar": // redo for nesting
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
	case "reverse": // redo for nesting
		return reverse(expression, symbols, functions, bindings)
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
	case ">":
		result, err := gt_or_lt(expression, symbols, functions, bindings, ">")
		return result, err
	case "<":
		result, err := gt_or_lt(expression, symbols, functions, bindings, "<")
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
			retVal, err := ExecFunction(body, symbols, functions, &function.Bindings)
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
	bindings *map[string]string) (interface{}, error) {
	// honestly will probably need to be reworked in the future
	e := expression.Head
	e = e.Next()
	if e.Data == "'" {
		e = e.Next()
	}
	switch e.Data.(type) {
	case string:
		return nil, errors.New("car requires a list")
	default:
		l := e.Data.(structs.List)
		e = l.Head
		if e.Data == "list" {
			e = l.Head.Next()
		} else if (*symbols)[e.Data.(string)] == 'f' {
			retVal, err := ExecFunction(l, symbols, functions, bindings)
			if err != nil {
				return nil, err
			}
			e = retVal.(structs.List).Head
		}
		return e.Data, nil
	}
}

func cdr(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]string) (interface{}, error) {
	e := expression.Head
	e = e.Next()
	if e.Data == "'" {
		e = e.Next()
	}
	switch e.Data.(type) {
	case string:
		return nil, errors.New("cdr requires a list")
	default:
		//e = e.Next()
		l := e.Data.(structs.List)
		e = l.Head
		if e.Data == "list" {
			e = l.Head.Next()
		} else if (*symbols)[e.Data.(string)] == 'f' {
			retVal, err := ExecFunction(l, symbols, functions, bindings)
			if err != nil {
				return nil, err
			}
			e = retVal.(structs.List).Head
		}
		e = e.Next()
		var rest structs.List
		for ; e != nil; e = e.Next() {
			rest.PushBack(e.Data.(string))
		}
		return rest, nil
	}
}

func cons(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]string) (interface{}, error) {
	var first structs.List
	var second structs.List
	var list structs.List

	e := expression.Head
	e = e.Next()
	if e.Data == "'" {
		e = e.Next()
	}
	placeHolder := e
	switch e.Data.(type) {
	case structs.List:
		l := e.Data.(structs.List)
		e = l.Head
		if (*symbols)[e.Data.(string)] == 'f' {
			value, err := ExecFunction(l, symbols, functions, bindings)
			if err != nil {
				return nil, err
			}
			list.PushBack(value)
		} else {
			for ; e != nil; e = e.Next() {
				first.PushBack(e.Data)
			}
			list.PushBack(first)
		}
	default:
		list.PushBack(e.Data)
	}
	e = placeHolder
	e = e.Next()
	if e.Data == "'" {
		e = e.Next()
	}
	switch e.Data.(type) {
	case structs.List:
		l := e.Data.(structs.List)
		e = l.Head
		if (*symbols)[e.Data.(string)] == 'f' {
			value, err := ExecFunction(l, symbols, functions, bindings)
			if err != nil {
				return nil, err
			}
			list.PushBack(value)
		} else {
			for ; e != nil; e = e.Next() {
				second.PushBack(e.Data)
			}
			list.PushBack(second)
		}
	default:
		list.PushBack(e.Data)
	}
	return list, nil
}

func defun(expression structs.List, symbols *map[string]rune,
	functions *map[string]structs.Function) (structs.List, error) {
	funct := new(structs.Function)
	e := expression.Head
	e = e.Next()
	funct.Name = e.Data.(string)
	e = e.Next()
	params(e.Data.(structs.List), funct, symbols, functions)
	e = e.Next()
	funct.Body = e.Data.(structs.List)

	(*symbols)[funct.Name] = 'f'
	(*functions)[funct.Name] = *funct

	return expression, nil
}

func params(expression structs.List, funct *structs.Function,
	symbols *map[string]rune, functions *map[string]structs.Function) error {
	for e := expression.Head; e != nil; e = e.Next() {
		funct.Args = append(funct.Args, e.Data.(string))
	}

	return nil
}

func defvar(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]string) (interface{}, error) {
	if expression.Len() != 3 {
		return nil, errors.New("Invalid number of arguments supplied to defvar")
	}
	e := expression.Head
	e = e.Next()
	symbol := e.Data.(string)
	e = e.Next()
	switch e.Data.(type) {
	case string:
		value := e.Data.(string)
		(*bindings)[symbol] = value

		return strings.ToUpper(symbol), nil
	default:
		l := e.Data.(structs.List)
		retVal, err := ExecFunction(l, symbols, functions, bindings)
		if err != nil {
			return nil, err
		}
		switch retVal.(type) {
		case float64:
			value := strconv.FormatFloat(retVal.(float64), 'f', 6, 64)
			(*bindings)[symbol] = value
			return strings.ToUpper(symbol), nil
		case string:
			(*bindings)[symbol] = retVal.(string)
			return strings.ToUpper(symbol), nil
		default:
			return nil, errors.New("For the time being, lists cannot be defined as variable values")
			//(*bindings)[symbol] = retVal.(structs.List)
			//return strings.ToUpper(symbol), nil
		}
	}
}

func gt_or_lt(expression structs.List, symbols *map[string]rune,
	functions *map[string]structs.Function, bindings *map[string]string, fun string) (interface{}, error) {
	var a, b float64
	var err error
	e := expression.Head
	e = e.Next()
	if e.Data == "'" {
		return nil, errors.New("> requires valid numerical values, not lists")
	}
	switch e.Data.(type) {
	case float64:
		a = e.Data.(float64)
	case string:
		a, err = strconv.ParseFloat(e.Data.(string), 64)
		if err != nil {
			return nil, errors.New("Cannot parse invalid value")
		}
	case structs.List:
		d := e
		l := e.Data.(structs.List)
		e = l.Head
		if (*symbols)[e.Data.(string)] == 'f' {
			retVal, err := ExecFunction(l, symbols, functions, bindings)
			if err != nil {
				return nil, err
			}
			a = retVal.(float64)
		} else {
			return nil, errors.New("Lists are not valid numerical values")
		}
		e = d
	default:
		return nil, errors.New("> requires a number value")
	}
	e = e.Next()
	if e.Data == "'" {
		return nil, errors.New("> requires valid numerical values, not lists")
	}
	switch e.Data.(type) {
	case float64:
		b = e.Data.(float64)
	case string:
		b, err = strconv.ParseFloat(e.Data.(string), 64)
		if err != nil {
			return nil, errors.New("Cannot parse invalid value")
		}
	case structs.List:
		d := e
		l := e.Data.(structs.List)
		e = l.Head
		if (*symbols)[e.Data.(string)] == 'f' {
			retVal, err := ExecFunction(l, symbols, functions, bindings)
			if err != nil {
				return nil, err
			}
			b = retVal.(float64)
		} else {
			return nil, errors.New("Lists are not valid numerical values")
		}
		e = d
	default:
		return nil, errors.New("> requires a number value")
	}
	if fun == ">" {
		return a > b, nil
	} else {
		return a < b, nil
	}
}

func last(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]string) (interface{}, error) {
	e := expression.Head
	e = e.Next()
	switch e.Data.(type) {
	case string:
		if e.Data == "'" {
			e = e.Next()
			l := e.Data.(structs.List)
			e = l.Tail
			return e.Data, nil
		} else {
			return nil, errors.New("last requires a list")
		}
	default:
		l := e.Data.(structs.List)
		e = l.Head
		if e.Data == "list" {
			e = l.Tail
		} else if (*symbols)[e.Data.(string)] == 'f' {
			retVal, err := ExecFunction(l, symbols, functions, bindings)
			if err != nil {
				return nil, err
			}
			e = retVal.(structs.List).Tail
		} else {
			e = l.Tail
			return e.Data, nil
		}
		return e.Data, nil
	}
}

// will need to add in symbols, functions and bindings later
func list(expression structs.List) (structs.List, error) {
	var newList structs.List
	for e := expression.Head; e != nil; e = e.Next() {
		if e.Data == "list" {
			continue
		} else if e.Data == "'" {
			continue
		}
		switch e.Data.(type) {
		case string:
			newList.PushBack(e.Data)
		default:
			if subList, err := list(e.Data.(structs.List)); err == nil {
				newList.PushBack(subList)
			} else {
				return subList, err
			}
		}
	}
	return newList, nil
}

func plus(expression structs.List, symbols *map[string]rune,
	functions *map[string]structs.Function, bindings *map[string]string) (float64, error) {
	if expression.Len() == 1 {
		return 0.0, errors.New("Invalid number of arguments.")
	}
	var sum float64
	e := expression.Head
	for e = e.Next(); e != nil; e = e.Next() {
		switch e.Data.(type) {
		case string:
			number := e.Data
			number = (*bindings)[number.(string)]
			if number == "" {
				number = e.Data.(string)
			}
			if f, err := strconv.ParseFloat(number.(string), 64); err == nil {
				sum += f
			} else {
				return 0.0, err
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
	bindings *map[string]string) (interface{}, error) {
	var number string
	var difference float64
	var num_expr int
	var err error
	e := expression.Head
	e = e.Next()

	if expression.Len() == 1 {
		return expression, errors.New("Invalid number of arguments.")
	} else if expression.Len() == 2 {
		number = (*bindings)[e.Data.(string)]
		fmt.Println("NUM: ", number)
		if number == "" {
			number = e.Data.(string)
		}
		difference, err := strconv.ParseFloat(number, 64)
		if err != nil {
			return expression, err
		}
		difference = 0 - difference
		return difference, nil
	} else {
		for ; e != nil; e = e.Next() {
			num_expr++
			switch e.Data.(type) {
			case string:
				if num_expr == 1 {
					number = (*bindings)[e.Data.(string)]
					if number == "" {
						difference, err = strconv.ParseFloat(e.Data.(string), 64)
						continue
					} else {
						difference, err = strconv.ParseFloat(number, 64)
						continue
					}
				}
				if err != nil {
					return expression, err
				}
				number = (*bindings)[e.Data.(string)]
				if number == "" {
					if num, err := strconv.ParseFloat(e.Data.(string), 64); err == nil {
						difference -= num
					} else {
						return expression, err
					}
				} else {
					if num, err := strconv.ParseFloat(number, 64); err == nil {
						difference -= num
					} else {
						return expression, err
					}
				}
			default:
				subValue, err := ExecFunction(e.Data.(structs.List), symbols, functions, bindings)
				if err != nil {
					return 0.0, err
				}
				if num_expr == 1 {
					difference = subValue.(float64)
				} else {
					difference -= subValue.(float64)
				}

			}
		}
		return difference, nil
	}
}

func reverse(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]string) (interface{}, error) {
	var reversed structs.List

	e := expression.Head
	e = e.Next()
	if e.Data == "'" {
		e = e.Next()
	}
	switch e.Data.(type) {
	case string:
		return nil, errors.New("reverse requires a list")
	default:
		retVal := e.Data.(structs.List)
		e = retVal.Head
		if e.Data == "list" {
			e = retVal.Head.Next()
		} else if (*symbols)[e.Data.(string)] == 'f' {
			l, err := ExecFunction(retVal, symbols, functions, bindings)
			if err != nil {
				return nil, err
			}
			retVal = l.(structs.List)
		}
		e = retVal.Tail
		for ; e != nil; e = e.Prev() {
			reversed.PushBack(e.Data)
		}
		return reversed, nil
	}
}

func times(expression structs.List, symbols *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]string) (interface{}, error) {
	if expression.Len() == 1 {
		return expression, errors.New("Invalid number of arguments.")
	}
	product := 1.0
	var number string
	e := expression.Head
	for e = e.Next(); e != nil; e = e.Next() {
		switch e.Data.(type) {
		case string:
			number = (*bindings)[e.Data.(string)]
			if number == "" {
				if num, err := strconv.ParseFloat(e.Data.(string), 64); err == nil {
					product *= num
				} else {
					return expression, errors.New("Only numbers can be multiplied.")
				}
			} else {
				if num, err := strconv.ParseFloat(number, 64); err == nil {
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
	bindings *map[string]string) (interface{}, error) {
	var numerator float64
	var err error
	numExpr := 0
	e := expression.Head
	e = e.Next()

	if expression.Len() < 3 {
		return expression, errors.New("Invalid number of arguments.")
	}
	for ; e != nil; e = e.Next() {
		numExpr++
		switch e.Data.(type) {
		case string:
			numStr := (*bindings)[e.Data.(string)]
			if numStr != "" && numExpr == 1 {
				if numerator, err = strconv.ParseFloat(numStr, 64); err == nil {
					continue
				} else {
					return numStr, err
				}
			} else if numExpr == 1 {
				if numerator, err = strconv.ParseFloat(e.Data.(string), 64); err == nil {
					continue
				} else {
					return e.Data.(string), err
				}
			}
			if number := (*bindings)[e.Data.(string)]; number != "" {
				if num, err := strconv.ParseFloat(number, 64); err == nil {
					numerator /= num
				} else {
					return e.Data.(string), errors.New("Only numbers can be divided")
				}
			} else {
				if num, err := strconv.ParseFloat(e.Data.(string), 64); err == nil {
					numerator /= num
				} else {
					return e.Data.(string), errors.New("Only numbers can be divided")
				}
			}
		default:
			subValue, err := ExecFunction(e.Data.(structs.List), symbols, functions, bindings)
			if err != nil {
				return 0.0, err
			}
			if numExpr == 1 {
				numerator = subValue.(float64)
			} else {
				numerator /= subValue.(float64)
			}
		}
	}
	return numerator, nil
}
