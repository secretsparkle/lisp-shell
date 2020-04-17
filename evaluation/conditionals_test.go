package evaluation

import (
	"../engines"
	"../structs"
	"testing"
)

func TestEqual(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]bool{
		"(equal t t)":                       true,
		"(equal 3 3)":                     true,
		"(equal nil nil)":                    true,
		"(equal 3 4)":             false,
		"(equal t nil)":           false,
		"(equal '(1 2 3) (list 1 2 3))": true,
		"(equal '(1 2 3) '(1 2 3))": true,
		"(equal '(1 2) '(1 2 3))": false
		"(equal (equal 3 3) (equal 3 3))":   true,
	}

	unsuccessful_tests := map[string]interface{}{
		"(equal)": nil
		"(equal 1)": nil
		"(equal $test nil)": nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := equal(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %f", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := equal(expression, &types, &functions, &bindings)
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
