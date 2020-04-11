package conditionals

import (
	"../functions"
	"../structs"
	"../utils"
	"errors"
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
			} else if e.Data.(string) == "t" {
				last = e.Data
			} else if e.Data.(string) == "nil" {
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
			} else if (*symbols)[e.Data.(string)] == 'f' || (*symbols)[e.Data.(string)] == 'c' {
				value, err := ExecInput(e.Data.(structs.List), symbols, functionTable, bindings)
				if err != nil {
					return nil, err
				} else if value == false {
					return false, nil
				} else {
					last = value
				}
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
	e := expression.Head
	e = e.Next()
	c := e.Data.(structs.List)
	value, err := ExecInput(c, symbols, functionTable, bindings)
	if err != nil {
		return nil, err
	}
	// must return bool
	// need to change to use any truthy value
	switch value.(type) {
	case bool:
	default:
		value = true
	}

	e = e.Next()
	if e.Data == "'" {
		e = e.Next()
	}
	// need switch
	t := e.Data.(structs.List)
	e = e.Next()
	if e.Data == "'" {
		e = e.Next()
	}
	// need switch
	f := e.Data.(structs.List)

	if value == true {
		// need to check if t is a function or a value
		if retVal, err := ExecInput(t, symbols, functionTable, bindings); err != nil {
			return nil, err
		} else {
			return retVal, nil
		}
	} else {
		// need to check if t is a function or a value
		if retVal, err := ExecInput(f, symbols, functionTable, bindings); err != nil {
			return nil, err
		} else {
			return retVal, nil
		}
	}
}
