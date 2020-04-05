package conditionals

import (
	"../functions"
	"../structs"
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
	case "cond":
	case "if":
		result, err := if_statement(expression, symbols, functionTable, bindings)
		return result, err
	}
	return nil, nil
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
	switch value.(type) {
	case bool:
	default:
		return nil, errors.New("The if conditional requires a boolean value")
	}

	e = e.Next()
	t := e.Data.(structs.List)
	e = e.Next()
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
