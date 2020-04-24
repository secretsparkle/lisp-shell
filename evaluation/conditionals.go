package evaluation

import (
	"../structs"
	"../utils"
	"errors"
	"strconv"
)

func equal(expression structs.List, functionList *map[string]rune,
	functions *map[string]structs.Function, bindings *map[string]interface{}) (interface{}, error) {
	var a, b interface{}

	if expression.Len() != 3 {
		return nil, errors.New("Invalid equal expression.")
	}

	e := expression.Head
	e = e.Next()
	if e.Data == "'" {
		e = e.Next()
	}

	switch e.Data.(type) {
	case bool:
		a = e.Data.(bool)
	case float64:
		a = e.Data.(float64)
	case string:
		if a = (*bindings)[e.Data.(string)]; a != "" {
			a, _ = strconv.ParseFloat(a.(string), 64)
		} else if e.Data.(string) == "t" || e.Data.(string) == "T" {
			a = e.Data
		} else if e.Data.(string) == "nil" || e.Data.(string) == "NIL" {
			a = false
		} else if util.IsAlphabetic(e.Data.(string)) {
			return e.Data, errors.New("Unbound symbol, cannot evaluate")
		} else if util.AnySymbol(e.Data.(string)) {
			return e.Data, errors.New("Cannot evaluate symbolic input")
		} else if util.IsNumber(e.Data.(string)) {
			a, _ = strconv.ParseFloat(e.Data.(string), 64)
		} else {
			return e.Data, errors.New("Invalid argument")
		}
	case structs.List:
		d := e
		l := e.Data.(structs.List)
		e = l.Head
		if (*functionList)[e.Data.(string)] == 'f' {
			retVal, err := EvaluateFunction(l, functionList, functions, bindings)
			if err != nil {
				return nil, err
			}
			a = retVal
		} else {
			a = l
		}
		e = d
	}
	e = e.Next()
	if e.Data == "'" {
		e = e.Next()
	}
	switch e.Data.(type) {
	case bool:
		b = e.Data.(bool)
	case float64:
		b = e.Data.(float64)
	case string:
		if b = (*bindings)[e.Data.(string)]; b != "" {
			b, _ = strconv.ParseFloat(b.(string), 64)
		} else if e.Data.(string) == "t" || e.Data.(string) == "T" {
			b = e.Data
		} else if e.Data.(string) == "nil" || e.Data.(string) == "NIL" {
			b = false
		} else if util.IsAlphabetic(e.Data.(string)) {
			return e.Data, errors.New("Unbound symbol, cannot evaluate")
		} else if util.AnySymbol(e.Data.(string)) {
			return e.Data, errors.New("Cannot evaluate symbolic input")
		} else if util.IsNumber(e.Data.(string)) {
			b, _ = strconv.ParseFloat(e.Data.(string), 64)
		} else {
			return e.Data, errors.New("Invalid argument")
		}
	case structs.List:
		l := e.Data.(structs.List)
		e = l.Head
		if (*functionList)[e.Data.(string)] == 'f' {
			retVal, err := EvaluateFunction(l, functionList, functions, bindings)
			if err != nil {
				return nil, err
			}
			b = retVal
		} else {
			b = l
		}
	}
	switch a.(type) {
	case bool:
		if a == b {
			return true, nil
		} else {
			return false, nil
		}
	case float64:
		if a == b {
			return true, nil
		} else {
			return false, nil
		}
	case string:
		if a == b {
			return true, nil
		} else {
			return false, nil
		}
	case structs.List:
		l := a.(structs.List)
		m := b.(structs.List)
		if l.Len() != m.Len() {
			return false, nil
		}
		d := l.Head
		e := m.Head
		for ; d != nil && e != nil; d, e = d.Next(), e.Next() {
			if d.Data == e.Data {
				continue
			} else {
				return false, nil
			}
		}
		return true, nil
	}
	if a == b {
		return true, nil
	} else {
		return false, nil
	}
}

func gt_or_lt(expression structs.List, functionList *map[string]rune,
	functions *map[string]structs.Function, bindings *map[string]interface{}, fun string) (interface{}, error) {
	var a, b float64
	var err error
	var str interface{}

	if expression.Len() != 3 {
		return nil, errors.New("Invalid number of arguments")
	}

	e := expression.Head
	e = e.Next()
	if e.Data == "'" {
		return nil, errors.New("> requires valid numerical values, not lists")
	}
	switch e.Data.(type) {
	case float64:
		a = e.Data.(float64)
	case string:
		if str = (*bindings)[e.Data.(string)]; str != "" {
			a, err = strconv.ParseFloat(str.(string), 64)
			if err != nil {
				return nil, errors.New("Cannot parse invalid value")
			}
		} else {
			a, err = strconv.ParseFloat(e.Data.(string), 64)
			if err != nil {
				return nil, errors.New("Cannot parse invalid value")
			}
		}
	case structs.List:
		d := e
		l := e.Data.(structs.List)
		e = l.Head
		if (*functionList)[e.Data.(string)] == 'f' {
			retVal, err := EvaluateFunction(l, functionList, functions, bindings)
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
		if str = (*bindings)[e.Data.(string)]; str != "" {
			b, err = strconv.ParseFloat(str.(string), 64)
			if err != nil {
				return nil, errors.New("Cannot parse invalid value")
			}
		} else {
			b, err = strconv.ParseFloat(e.Data.(string), 64)
			if err != nil {
				return nil, errors.New("Cannot parse invalid value")
			}
		}
	case structs.List:
		d := e
		l := e.Data.(structs.List)
		e = l.Head
		if (*functionList)[e.Data.(string)] == 'f' {
			retVal, err := EvaluateFunction(l, functionList, functions, bindings)
			if err != nil {
				return nil, err
			}
			// need exception code here
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
