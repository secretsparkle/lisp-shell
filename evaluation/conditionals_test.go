package evaluation

import (
	"../engines"
	"../structs"
	"testing"
)

func TestEqual(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(and t)")
	test1, _ := and(expressions, &types, &functions, &bindings)

	if test1 != "t" {
		t.Errorf("(and t) = %v; want t", test1)
	}
	successful_tests := map[string]interface{}{
		"(and t)":                       "t",
		"(and t t)":                     "t",
		"(and t 3)":                     "3",
		"(and (and 3 t t))":             "t",
		"(and (and 3 t t) t)":           "t",
		"(and (and 3 t t) (and 3 t t))": "t",
		//"(and t '(1 2 3))": "(1 2 3)",
		"(and nil)":   false,
		"(and t nil)": false,
		"(and nil t)": false,
	}

	unsuccessful_tests := map[string]interface{}{
		"(and $test nil)": nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := and(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %f", test, result, desired)
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

func TestGreaterThan(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(and t)")
	test1, _ := and(expressions, &types, &functions, &bindings)

	if test1 != "t" {
		t.Errorf("(and t) = %v; want t", test1)
	}
}

func TestLessThan(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(and t)")
	test1, _ := and(expressions, &types, &functions, &bindings)

	if test1 != "t" {
		t.Errorf("(and t) = %v; want t", test1)
	}
}
