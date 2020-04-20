package evaluation

import (
	"../engines"
	"../structs"
	"testing"
)

func TestEqual(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]bool{
		"(equal t t)":                     true,
		"(equal 3 3)":                     true,
		"(equal nil nil)":                 true,
		"(equal 3 4)":                     false,
		"(equal t nil)":                   false,
		"(equal '(1 2 3) (list 1 2 3))":   true,
		"(equal '(1 2 3) '(1 2 3))":       true,
		"(equal '(1 2) '(1 2 3))":         false,
		"(equal (equal 3 3) (equal 3 3))": true,
		"(equal (+ 3 3) (+ 3 3))":         true,
	}

	unsuccessful_tests := map[string]interface{}{
		"(equal)":           nil,
		"(equal 1)":         nil,
		"(equal $test nil)": nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := equal(expression, &types, &functions, &bindings)
		if result != desired {
			t.Errorf("%s = %v; want %t", test, result, desired)
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

	successful_tests := map[string]bool{
		"(> 3 4)":       false,
		"(> 3 3)":       false,
		"(> 4 3)":       true,
		"(> (+ 3 3) 4)": true,
		"(> 4 (+ 3 3))": false,
	}

	unsuccessful_tests := map[string]interface{}{
		"(>)":           nil,
		"(> 1)":         nil,
		"(> $test nil)": nil,
		"(> 2 3 4)":     nil,
		"(> '(1 2 3))":  nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := gt_or_lt(expression, &types, &functions, &bindings, ">")
		if result != desired {
			t.Errorf("%s = %v; want %t", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := gt_or_lt(expression, &types, &functions, &bindings, ">")
		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}

func TestLessThan(t *testing.T) {
	types, functions, bindings := structs.Maps()

	successful_tests := map[string]bool{
		"(< 3 4)":       true,
		"(< 3 3)":       false,
		"(< 4 3)":       false,
		"(< (+ 3 3) 4)": false,
		"(< 4 (+ 3 3))": true,
	}

	unsuccessful_tests := map[string]interface{}{
		"(<)":           nil,
		"(< 1)":         nil,
		"(< $test nil)": nil,
		"(< 2 3 4)":     nil,
		"(< '(1 2 3))":  nil,
	}

	for test, desired := range successful_tests {
		expression, _ := engines.Translate(test)
		result, _ := gt_or_lt(expression, &types, &functions, &bindings, "<")
		if result != desired {
			t.Errorf("%s = %v; want %t", test, result, desired)
		}
	}

	for test, undesired := range unsuccessful_tests {
		expression, _ := engines.Translate(test)
		_, err := gt_or_lt(expression, &types, &functions, &bindings, "<")
		if err == undesired {
			t.Errorf("%s triggered %v; want error triggered.", test, err)
		}
	}
}
