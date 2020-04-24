package evaluation

import (
	"../structs"
	"errors"
)

func car(expression structs.List, functionList *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]interface{}) (interface{}, error) {

	if expression.Len() != 2 {
		return nil, errors.New("Invalid number of arguments")
	}
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
			e = e.Next()
		} else if e.Data == "'" {
			e = e.Next()
		} else if (*functionList)[e.Data.(string)] == 'f' {
			retVal, err := EvaluateFunction(l, functionList, functions, bindings)
			if err != nil {
				return nil, err
			}
			e = retVal.(structs.List).Head
		}
		return e.Data, nil
	}
}

func cdr(expression structs.List, functionList *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]interface{}) (interface{}, error) {

	if expression.Len() != 2 {
		return nil, errors.New("Invalid expression length")
	}
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
			e = e.Next()
		} else if e.Data == "'" {
			e = e.Next()
		} else if (*functionList)[e.Data.(string)] == 'f' {
			retVal, err := EvaluateFunction(l, functionList, functions, bindings)
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

func cons(expression structs.List, functionList *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]interface{}) (interface{}, error) {
	var first structs.List
	var second structs.List
	var list structs.List

	if expression.Len() != 3 {
		return nil, errors.New("Invalid number of arguments")
	}
	e := expression.Head
	e = e.Next()
	if e.Data == "'" {
		e = e.Next()
	}
	placeHolder := e
	switch e.Data.(type) {
	case string:
		a := (*bindings)[e.Data.(string)]
		if a == "" {
			a = e.Data.(string)
		}
		list.PushBack(a)
	case structs.List:
		l := e.Data.(structs.List)
		e = l.Head
		if e.Data == "'" {
			e = e.Next()
		}
		if (*functionList)[e.Data.(string)] == 'f' {
			value, err := EvaluateFunction(l, functionList, functions, bindings)
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
	case string:
		a := (*bindings)[e.Data.(string)]
		if a == "" {
			a = e.Data.(string)
		}
		list.PushBack(a)
	case structs.List:
		l := e.Data.(structs.List)
		e = l.Head
		if e.Data == "'" {
			e = e.Next()
		}
		if (*functionList)[e.Data.(string)] == 'f' {
			value, err := EvaluateFunction(l, functionList, functions, bindings)
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

func last(expression structs.List, functionList *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]interface{}) (interface{}, error) {

	if expression.Len() != 2 {
		return nil, errors.New("Invalid number of expressions")
	}
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
		} else if (*functionList)[e.Data.(string)] == 'f' {
			retVal, err := EvaluateFunction(l, functionList, functions, bindings)
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

func reverse(expression structs.List, functionList *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]interface{}) (interface{}, error) {
	var reversed structs.List

	if expression.Len() != 2 {
		return nil, errors.New("Invalid number of arguments")
	}
	e := expression.Head
	e = e.Next()
	if e.Data == "'" {
		e = e.Next()
	}
	switch e.Data.(type) {
	case structs.List:
		retVal := e.Data.(structs.List)
		e = retVal.Head
		if e.Data == "list" {
			e = retVal.Head.Next()
		} else if (*functionList)[e.Data.(string)] == 'f' {
			l, err := EvaluateFunction(retVal, functionList, functions, bindings)
			if err != nil {
				return nil, err
			}
			retVal = l.(structs.List)
		}
		e = retVal.Tail
		for ; e != nil; e = e.Prev() {
			if e.Data == "'" || e.Data == "list" {
				continue
			}
			reversed.PushBack(e.Data)
		}
		return reversed, nil
	default:
		return nil, errors.New("reverse requires a list")
	}
}
