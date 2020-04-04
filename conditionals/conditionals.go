package conditionals

import (
	"../functions"
	"../structs"
	//"errors"
	"fmt"
)

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
	value, err := functions.ExecFunction(c, symbols, functionTable, bindings)
	if err != nil {
		return nil, err
	}

	e = e.Next()
	t := e.Data.(structs.List)
	e = e.Next()
	f := e.Data.(structs.List)

	if value == true {
		// need to check if t is a function or a value
		if retVal, err := main.ExecInput(t, symbols, functionTable, bindings); err != nil {
			return nil, err
		} else {
			return retVal, nil
		}
	} else {
		// need to check if t is a function or a value
		if retVal, err := main.ExecInput(f, symbols, functionTable, bindings); err != nil {
			return nil, err
		} else {
			return retVal, nil
		}
	}
}
