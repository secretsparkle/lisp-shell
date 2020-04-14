package evaluation

import (
	"../structs"
	"errors"
	"strconv"
)

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
			subValue, err := EvaluateFunction(e.Data.(structs.List), symbols, functions, bindings)
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
				subValue, err := EvaluateFunction(e.Data.(structs.List), symbols, functions, bindings)
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
			subValue, err := EvaluateFunction(e.Data.(structs.List), symbols, functions, bindings)
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
			subValue, err := EvaluateFunction(e.Data.(structs.List), symbols, functions, bindings)
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
