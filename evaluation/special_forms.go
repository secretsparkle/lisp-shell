package evaluation

import (
	"../structs"
	"../utils"
	"errors"
	"strconv"
	"strings"
)

func and(expression structs.List, functionList *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]interface{}) (interface{}, error) {
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
			if (*functionList)[e.Data.(string)] == 'f' {
				value, err := EvaluateFunction(l, functionList, functions, bindings)
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

func define(expression structs.List, functionList *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]interface{}) (interface{}, error) {
	if expression.Len() != 3 {
		return nil, errors.New("invalid number of arguments")
	}
	var symbol string
	quoted := false
	e := expression.Head
	e = e.Next()

	switch e.Data.(type) {
	case string:
		if util.IsValidSymbol(e.Data.(string)) {
			if (*bindings)[symbol] != nil {
				return nil, errors.New("symbol has already been defined")
			}
			symbol = e.Data.(string)
		} else {
			return e.Data, errors.New("invalid starting character on symbol")
		}
		e = e.Next()
		if e.Data == "'" {
			quoted = true
			e = e.Next()
		}
		switch e.Data.(type) {
		case string:
			// should probably get rid of this else if tree
			// any string should be valid as a symbol
			// although numbers should be converted
			value := e.Data.(string)
			if quoted == true {
				(*bindings)[symbol] = value
				return strings.ToUpper(symbol), nil
			} else if util.IsNumber(value) {
				(*bindings)[symbol], _ = strconv.ParseFloat(value, 64)
				return strings.ToUpper(symbol), nil
			} else {
				(*bindings)[symbol] = value
				return strings.ToUpper(symbol), nil
			}
		case structs.List:
			value := e.Data.(structs.List)
			if quoted == true {
				(*bindings)[symbol] = value
				return strings.ToUpper(symbol), nil
			}
			l := e.Data.(structs.List)
			e = l.Head
			if (*functionList)[e.Data.(string)] == 'f' {
				retVal, err := EvaluateFunction(l, functionList, functions, bindings)
				if err != nil {
					return nil, err
				}
				switch retVal.(type) {
				case float64:
					(*bindings)[symbol] = retVal
					return strings.ToUpper(symbol), nil
				case string:
					(*bindings)[symbol] = retVal
					return strings.ToUpper(symbol), nil
				case structs.List:
					(*bindings)[symbol] = retVal
					return strings.ToUpper(symbol), nil
				}
			} else {
				return nil, errors.New("invalid function call")
			}
		default:
			return e.Data, errors.New("cannot parse this value as a symbol")
		}
	case structs.List:
		// congrats, its a function
		funct := new(structs.Function)
		symbol := e.Data.(structs.List)
		f := symbol.Head
		funct.Name = f.Data.(string)
		f = f.Next()
		// for loop for function arguments
		for ; f != nil; f = f.Next() {
			funct.Args = append(funct.Args, f.Data.(string))
		}
		// move to function body
		e = e.Next()
		switch e.Data.(type) {
		case structs.List:
			// this list needs to be parsed to ensure it is a valid form
			// but this will be a future addition, when I determine what that
			// form will be
			funct.Body = e.Data.(structs.List)
		default:
			return e.Data, errors.New("invalid function type")
		}
		(*functionList)[funct.Name] = 'f'
		(*functions)[funct.Name] = *funct
		return strings.ToUpper(funct.Name), nil
	default:
		return nil, errors.New("invalid argument supplied to define")
	}

	return nil, nil
}

func if_statement(expression structs.List, functionList *map[string]rune,
	functions *map[string]structs.Function, bindings *map[string]interface{}) (interface{}, error) {
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
		retVal, err := EvaluateFunction(value.(structs.List), functionList, functions, bindings)
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
				if retVal, err := EvaluateFunction(t.(structs.List), functionList, functions, bindings); err == nil {
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
			if retVal, err := EvaluateFunction(t.(structs.List), functionList, functions, bindings); err == nil {
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
			if retVal, err := EvaluateFunction(f.(structs.List), functionList, functions, bindings); err == nil {
				return retVal, nil
			} else {
				return nil, err
			}
		default:
			return f, errors.New("Unknown error has occured")
		}
	}
}
