package evaluation

import (
	"../structs"
	"../utils"
	"errors"
	"strconv"
	"strings"
)

func and(expression structs.List, symbols *map[string]rune, functionTable *map[string]structs.Function,
	bindings *map[string]string) (interface{}, error) {
	e := expression.Head
	var last interface{}
	e = e.Next()

	for ; e != nil; e = e.Next() {
		isQuoted := false
		if e.Data == "'" {
			e = e.Next()
			isQuoted = true
		}
		switch e.Data.(type) {
		case string:
			if value := (*bindings)[e.Data.(string)]; value != "" {
				last = value
			} else if e.Data.(string) == "t" || e.Data.(string) == "T" {
				last = e.Data
			} else if e.Data.(string) == "nil" || e.Data.(string) == "NIL" {
				return false, nil
			} else if util.IsAlphabetic(e.Data.(string)) {
				return e.Data, errors.New("Unbound symbol, cannot evaluate")
			} else if util.AnySymbol(e.Data.(string)) {
				return e.Data, errors.New("Cannot evaluate symbolic input")
			} else if util.IsNumber(e.Data.(string)) {
				last = e.Data
			} else {
				return e.Data, errors.New("Invalid argument")
			}
		case structs.List:
			if isQuoted == true {
				last = e.Data.(structs.List)
				continue
			}
			l := e.Data.(structs.List)
			e = l.Head
			if (*symbols)[e.Data.(string)] == 'f' || (*symbols)[e.Data.(string)] == 'c' {
				value, err := EvaluateFunction(l, symbols, functionTable, bindings)
				if err != nil {
					return nil, err
				} else if value == false {
					return false, nil
				} else {
					last = value
				}
			} else {
				return nil, errors.New("Invalid list construction")
			}
		default:
			if e.Data == false {
				return false, nil
			}
			last = e.Data
		}
	}
	return last, nil
}

func defun(expression structs.List, symbols *map[string]rune,
	functions *map[string]structs.Function) (structs.List, error) {
	funct := new(structs.Function)
	e := expression.Head
	e = e.Next()
	funct.Name = e.Data.(string)
	e = e.Next()
	defunParamHelper(e.Data.(structs.List), funct, symbols, functions)
	e = e.Next()
	funct.Body = e.Data.(structs.List)

	(*symbols)[funct.Name] = 'f'
	(*functions)[funct.Name] = *funct

	return expression, nil
}

func defunParamHelper(expression structs.List, funct *structs.Function,
	symbols *map[string]rune, functions *map[string]structs.Function) error {
	for e := expression.Head; e != nil; e = e.Next() {
		funct.Args = append(funct.Args, e.Data.(string))
	}

	return nil
}

func defvar(expression structs.List, symbols *map[string]rune,
	functions *map[string]structs.Function, bindings *map[string]string) (interface{}, error) {
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
		retVal, err := EvaluateFunction(l, symbols, functions, bindings)
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
func if_statement(expression structs.List, symbols *map[string]rune,
	functionTable *map[string]structs.Function, bindings *map[string]string) (interface{}, error) {
	var condition bool
	e := expression.Head
	e = e.Next()
	cQuoted := false
	tQuoted := false
	fQuoted := false

	if expression.Len() != 3 && expression.Len() != 4 {
		return nil, errors.New("invalid number of arguments")
	}
	if e.Data == "'" {
		e = e.Next()
		cQuoted = true
	}

	value := e.Data
	switch value.(type) {
	case bool:
		condition = value.(bool)
	case float64:
		condition = true
	case string:
		if evaluatedSymbol := (*bindings)[value.(string)]; evaluatedSymbol != "" {
			condition = true
		} else if value == "t" || value == "T" {
			condition = true
		} else if value == "nil" || value == "NIL" {
			condition = false
		} else if util.IsAlphabetic(value.(string)) {
			return value, errors.New("Unbound symbol, cannot evaluate")
		} else if util.AnySymbol(value.(string)) {
			return value, errors.New("Cannot evaluate symbolic input")
		} else if util.IsNumber(value.(string)) {
			condition = true
		} else {
			return value, errors.New("Invalid argument")
		}
	case structs.List:
		if cQuoted == true {
			condition = true
			break
		}
		retVal, err := EvaluateFunction(value.(structs.List), symbols, functionTable, bindings)
		if err != nil {
			return nil, err
		}
		switch retVal.(type) {
		case bool:
			condition = retVal.(bool)
		case float64:
			condition = true
		case string:
			if evaluatedSymbol := (*bindings)[retVal.(string)]; evaluatedSymbol != "" {
				condition = true
			} else if retVal == "t" || retVal == "T" {
				condition = true
			} else if retVal == "nil" || retVal == "NIL" {
				condition = false
			} else if util.IsAlphabetic(retVal.(string)) {
				return retVal, errors.New("Unbound symbol, cannot evaluate")
			} else if util.AnySymbol(retVal.(string)) {
				return retVal, errors.New("Cannot evaluate symbolic input")
			} else if util.IsNumber(retVal.(string)) {
				condition = true
			} else {
				return retVal, errors.New("Invalid argument")
			}
		default:
			condition = true
		}
	default:
		return value, errors.New("Unknown error has occured")
	}

	e = e.Next()
	if e.Data == "'" {
		e = e.Next()
		tQuoted = true
	}
	t := e.Data
	e = e.Next()
	if e == nil {
		if condition == true {
			switch t.(type) {
			case bool:
				return t, nil
			case float64:
				return t, nil
			case string:
				if evaluatedSymbol := (*bindings)[t.(string)]; evaluatedSymbol != "" {
					return evaluatedSymbol, nil
				} else if t == "t" || t == "T" {
					return true, nil
				} else if t == "nil" || t == "NIL" {
					return false, nil
				} else if util.IsAlphabetic(t.(string)) {
					return t, errors.New("Unbound symbol, cannot evaluate")
				} else if util.AnySymbol(t.(string)) {
					return t, errors.New("Cannot evaluate symbolic input")
				} else if util.IsNumber(t.(string)) {
					return t, nil
				} else {
					return t, errors.New("Invalid argument")
				}
			case structs.List:
				if tQuoted == true {
					return t, nil
				}
				if retVal, err := EvaluateFunction(t.(structs.List), symbols, functionTable, bindings); err == nil {
					return retVal, nil
				} else {
					return nil, err
				}
			default:
				return t, errors.New("Unknown error has occured")

			}
		} else {
			return false, nil
		}
	}

	if e.Data == "'" {
		e = e.Next()
		fQuoted = true
	}
	// need switch
	f := e.Data

	if condition == true {
		switch t.(type) {
		case bool:
			return t, nil
		case float64:
			return t, nil
		case string:
			if evaluatedSymbol := (*bindings)[t.(string)]; evaluatedSymbol != "" {
				return evaluatedSymbol, nil
			} else if t == "t" || t == "T" {
				return true, nil
			} else if t == "nil" || t == "NIL" {
				return false, nil
			} else if util.IsAlphabetic(t.(string)) {
				return t, errors.New("Unbound symbol, cannot evaluate")
			} else if util.AnySymbol(t.(string)) {
				return t, errors.New("Cannot evaluate symbolic input")
			} else if util.IsNumber(t.(string)) {
				return t, nil
			} else {
				return t, errors.New("Invalid argument")
			}
		case structs.List:
			if tQuoted == true {
				return t, nil
			}
			if retVal, err := EvaluateFunction(t.(structs.List), symbols, functionTable, bindings); err == nil {
				return retVal, nil
			} else {
				return nil, err
			}
		default:
			return t, errors.New("Unknown error has occured")

		}
	} else {
		switch f.(type) {
		case bool:
			return f, nil
		case float64:
			return f, nil
		case string:
			if evaluatedSymbol := (*bindings)[f.(string)]; evaluatedSymbol != "" {
				return evaluatedSymbol, nil
			} else if f == "t" || f == "T" {
				return true, nil
			} else if f == "nil" || f == "NIL" {
				return false, nil
			} else if util.IsAlphabetic(f.(string)) {
				return f, errors.New("Unbound symbol, cannot evaluate")
			} else if util.AnySymbol(f.(string)) {
				return f, errors.New("Cannot evaluate symbolic input")
			} else if util.IsNumber(f.(string)) {
				return f, nil
			} else {
				return f, errors.New("Invalid argument")
			}
		case structs.List:
			if fQuoted == true {
				return f, nil
			}
			if retVal, err := EvaluateFunction(f.(structs.List), symbols, functionTable, bindings); err == nil {
				return retVal, nil
			} else {
				return nil, err
			}
		default:
			return f, errors.New("Unknown error has occured")
		}
	}
}
