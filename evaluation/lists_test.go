package evaluation

import (
	"../engines"
	"../structs"
	"testing"
)

func TestCar(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]interface{}{
		"(car '(1 2 3))":           "1",
		"(car (list 1 2 3))":       "1",
		"(car '(a b c))":           "a",
		"(car '((1 2 3) (4 5 6)))": "(1 2 3)",
	}

	unsuccessful_tests := map[string]interface{}{
		"(car)":           nil,
		"(car 1)":         nil,
		"(car $test nil)": nil,
		"(car 2 3 4)":     nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := car(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %v", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := car(expression, &types, &functions, &bindings)
		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}

func TestCdr(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]interface{}{
		"(cdr '(1 2 3))":           "(2 3)",
		"(cdr (list 1 2 3))":       "(2 3)",
		"(cdr '(a b c))":           "(b c)",
		"(cdr '((1 2 3) (4 5 6)))": "(4 5 6)",
	}

	unsuccessful_tests := map[string]interface{}{
		"(cdr)":           nil,
		"(cdr 1)":         nil,
		"(cdr $test nil)": nil,
		"(cdr 2 3 4)":     nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := cdr(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %v", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := cdr(expression, &types, &functions, &bindings)
		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}

func TestCons(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]interface{}{
		"(cons 1 '(2 3))":            "(1 2 3)",
		"(cons 1 (list 2 3))":        "(1 2 3)",
		"(cons 'a '(b c))":           "(a b c)",
		"(cons 1 '())":               "(1)",
		"(cons '(1 2 3) '((4 5 6)))": "((1 2 3) (4 5 6))",
	}

	unsuccessful_tests := map[string]interface{}{
		"(cons)":           nil,
		"(cons 1)":         nil,
		"(cons $test nil)": nil,
		"(cons 2 3 4)":     nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := cons(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %v", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := cons(expression, &types, &functions, &bindings)
		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}

func TestLast(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]interface{}{
		"(last '(1 2 3))":           "3",
		"(last (list 2 3))":         "3",
		"(last '(a b c))":           "c",
		"(last '((1 2 3) (4 5 6)))": "(4 5 6)",
	}

	unsuccessful_tests := map[string]interface{}{
		"(last)":           nil,
		"(last 1)":         nil,
		"(last $test nil)": nil,
		"(last 2 3 4)":     nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := last(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %v", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := last(expression, &types, &functions, &bindings)
		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}

func TestList(t *testing.T) {
	successful_tests := map[string]interface{}{
		"(list)":                     nil,
		"(list 1)":                   "(1)",
		"(list 2 3 4)":               "(2 3 4)",
		"(list 1 '(2 3))":            "(1 (2 3))",
		"(list 1 (list 2 3))":        "(1 (2 3))",
		"(list 'a '(b c))":           "(a (b c)) ",
		"(list 1 '())":               "(1 ())",
		"(list '(1 2 3) '((4 5 6)))": "((1 2 3) ((4 5 6)))",
	}

	unsuccessful_tests := map[string]interface{}{
		"(list $test nil)": nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := list(expression)
		if result != desired {
			t.Errorf("%s = %v; want %v", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := list(expression)
		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}

func TestReverse(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]interface{}{
		"(reverse '(1 2 3))":           "(3 2 1)",
		"(reverse (list 1 2 3))":       "(3 2 1)",
		"(reverse '(a b c))":           "(c b a)",
		"(reverse '((1 2 3) (4 5 6)))": "((4 5 6) (1 2 3))",
	}

	unsuccessful_tests := map[string]interface{}{
		"(reverse)":           nil,
		"(reverse 1)":         nil,
		"(reverse $test nil)": nil,
		"(reverse 2 3 4)":     nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := reverse(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %v", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := reverse(expression, &types, &functions, &bindings)
		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}
