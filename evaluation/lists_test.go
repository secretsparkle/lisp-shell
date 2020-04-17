package evaluation

import (
	"../engines"
	"../structs"
	"testing"
)

func TestCar(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(car '(1 2 3))")
	test1, _ := car(expressons, &types, &functions, &bindings)

	if test1 != "t" {
		t.Errorf("(car '(1 2 3)) = %v; want 1", test1)
	}
}

func TestCdr(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(and t)")
	test1, _ := and(expressions, &types, &functions, &bindings)

	if test1 != "t" {
		t.Errorf("(and t) = %v; want t", test1)
	}
}

func TestCons(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(and t)")
	test1, _ := and(expressions, &types, &functions, &bindings)

	if test1 != "t" {
		t.Errorf("(and t) = %v; want t", test1)
	}
}

func TestLast(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(and t)")
	test1, _ := and(expressions, &types, &functions, &bindings)

	if test1 != "t" {
		t.Errorf("(and t) = %v; want t", test1)
	}
}

func TestList(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(and t)")
	test1, _ := and(expressions, &types, &functions, &bindings)

	if test1 != "t" {
		t.Errorf("(and t) = %v; want t", test1)
	}
}

func TestReverse(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(and t)")
	test1, _ := and(expressions, &types, &functions, &bindings)

	if test1 != "t" {
		t.Errorf("(and t) = %v; want t", test1)
	}
}
