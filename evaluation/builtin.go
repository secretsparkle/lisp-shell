package evaluation

import (
	"../structs"
	"errors"
	"os"
	"strings"
)

func cd(expression structs.List, functionList *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]interface{}) (interface{}, error) {
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
		dir, err = EvaluateFunction(e.Data.(structs.List), functionList, functions, bindings)
		if err != nil {
			return nil, err
		}
	}
	// Change the directory and return the error.
	return nil, os.Chdir(dir.(string))
}

func echo(expression structs.List, functionList *map[string]rune, functions *map[string]structs.Function,
	bindings *map[string]interface{}) (interface{}, error) {
	e := expression.Head
	e = e.Next()
	var out []string
	for ; e != nil; e = e.Next() {
		switch e.Data.(type) {
		case string:
			out = append(out, e.Data.(string))
		default:
			value, err := EvaluateFunction(e.Data.(structs.List), functionList, functions, bindings)
			if err != nil {
				return nil, err
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
}
