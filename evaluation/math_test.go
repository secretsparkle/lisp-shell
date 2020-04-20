package evaluation

import (
	"../engines"
	"../structs"
	"testing"
)

func TestPlus(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]float64{
		"(+ 1 2)":             3,
		"(+ 1 2 3)":           6,
		"(+ 1 -4)":            -3,
		"(+ 1 (+ 1 2))":       4,
		"(+ 1 (+ (+ 1 2) 1))": 5,
		"(+ 1)":               1,
	}

	unsuccessful_tests := map[string]interface{}{
		"(+)":     nil,
		"(+ $ #)": nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := plus(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %f", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := plus(expression, &types, &functions, &bindings)

		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}

func TestMinus(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]float64{
		"(- 1)":               -1,
		"(- 2 1)":             1,
		"(- 3 2 1)":           0,
		"(- 1 2 3)":           -4,
		"(- 1 (- 1 2))":       2,
		"(- 1 (- (- 1 2) 1))": 3,
	}

	unsuccessful_tests := map[string]interface{}{
		"(-)":     nil,
		"(- $ #)": nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := minus(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %f", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := minus(expression, &types, &functions, &bindings)
		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}

func TestTimes(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]float64{
		"(* 1)":                 1,
		"(* 1 2)":               2,
		"(* 1 2 -3)":            -6,
		"(* 1 2 0)":             0,
		"(* 1 (* 1 2) (* 1 2))": 4,
		"(* 1 (* (* 1 2) 2))":   4,
	}

	unsuccessful_tests := map[string]interface{}{
		"(*)":     nil,
		"(* $ #)": nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := times(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %f", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := times(expression, &types, &functions, &bindings)
		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}

func TestDivide(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]float64{
		"(/ 1 2)":               0.5,
		"(/ 2 1)":               2,
		"(/ -3 2 1)":            -1.5,
		"(/ 6 (/ 2 1) (/ 3 1))": 1,
		"(/ 4 (/ (/ 2 1) 2))":   4,
	}

	unsuccessful_tests := map[string]interface{}{
		"(/ 1 0)": nil,
		"(/)":     nil,
		"(/ $ #)": nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := divide(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %f", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := divide(expression, &types, &functions, &bindings)
		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}
