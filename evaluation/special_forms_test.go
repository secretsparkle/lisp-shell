package evaluation

import (
	"../engines"
	"../structs"
	"testing"
)

func TestAnd(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]interface{}{
		"(and t)":                       "t",
		"(and t t)":                     "t",
		"(and t 3)":                     "3",
		"(and (and 3 t t))":             "t",
		"(and (and 3 t t) t)":           "t",
		"(and (and 3 t t) (and 3 t t))": "t",
		"(and t '(1 2 3))":              "(1 2 3)",
		"(and t (+ 1 2 3))":             6,
		"(and nil)":                     false,
		"(and t nil)":                   false,
		"(and nil t)":                   false,
	}

	unsuccessful_tests := map[string]interface{}{
		"(and $test nil)": nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := and(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %v", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := and(expression, &types, &functions, &bindings)
		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}

func TestDefun(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(and t)")
	test1, _ := and(expressions, &types, &functions, &bindings)

	if test1 != "t" {
		t.Errorf("(and t) = %v; want t", test1)
	}
}

func TestDefvar(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(and t)")
	test1, _ := and(expressions, &types, &functions, &bindings)

	if test1 != "t" {
		t.Errorf("(and t) = %v; want t", test1)
	}
}

func TestIfStatement(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]interface{}{
		"(if t 3 4)":           "3",
		"(if nil t 5)":         "5",
		"(if (and 1 2) 4 5)":   "4",
		"(if (and 1 nil) 4 t)": true,
		"(if t t)":             true,
		"(if '(1 2 3) 7 8)":    "7",
		"(if 2 3 4)":           "3",
		"(if (< 2 3) 2)":       "2",
		"(if (equal 3 3) 4 8)": "4",
		"(if (equal 3 4) 4 8)": "8",
	}

	unsuccessful_tests := map[string]interface{}{
		"(if $test nil)":   nil,
		"(if)":             nil,
		"(if 3)":           nil,
		"(if 1 2 3 4 5 6)": nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := if_statement(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %v", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := if_statement(expression, &types, &functions, &bindings)
		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}
