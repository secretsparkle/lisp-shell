package conditionals

import (
	"../functions"
	"../structs"
	"../utils"
	"errors"
	//"fmt"
)

func ExecInput(expression structs.List, symbols *map[string]rune,
	functionTable *map[string]structs.Function, bindings *map[string]string) (interface{}, error) {

	// Check for built-in commands
	switch (*symbols)[expression.Head.Data.(string)] {
	case 'c':
		value, err := EvalConditional(expression, symbols, functionTable, bindings)
		return value, err
	case 'f':
		value, err := functions.ExecFunction(expression, symbols, functionTable, bindings)
		return value, err
	default:
		return expression, nil
	}
}

func EvalConditional(expression structs.List, symbols *map[string]rune,
	functionTable *map[string]structs.Function, bindings *map[string]string) (interface{}, error) {
	switch expression.Head.Data {
	case "and":
		result, err := and(expression, symbols, functionTable, bindings)
		return result, err
	case "cond":
	case "if":
		result, err := if_statement(expression, symbols, functionTable, bindings)
		return result, err
	}
	return nil, nil
}

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
				value, err := ExecInput(l, symbols, functionTable, bindings)
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

func if_statement(expression structs.List, symbols *map[string]rune,
	functionTable *map[string]structs.Function, bindings *map[string]string) (interface{}, error) {
	var condition bool
	e := expression.Head
	e = e.Next()
	cQuoted := false
	tQuoted := false
	fQuoted := false

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
		_, err := ExecInput(value.(structs.List), symbols, functionTable, bindings)
		if err != nil {
			return nil, err
		}
		condition = true
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
				if retVal, err := ExecInput(t.(structs.List), symbols, functionTable, bindings); err == nil {
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
			if retVal, err := ExecInput(t.(structs.List), symbols, functionTable, bindings); err == nil {
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
			if retVal, err := ExecInput(f.(structs.List), symbols, functionTable, bindings); err == nil {
				return retVal, nil
			} else {
				return nil, err
			}
		default:
			return f, errors.New("Unknown error has occured")
		}
	}
}
