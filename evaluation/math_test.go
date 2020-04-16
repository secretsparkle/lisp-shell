package evaluation

import (
	"../engines"
	"../structs"
	"testing"
)

func TestPlus(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(+ 1 2)")
	test1, _ := plus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(+ 1 2 3)")
	test2, _ := plus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(+ 1 -4)")
	test3, _ := plus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(+ 1 (+ 1 2))")
	test4, _ := plus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(+ 1 (+ (+ 1 2) 1))")
	test5, _ := plus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(+ 1)")
	test6, _ := plus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(+)")
	_, err7 := plus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(+ $ #)")
	_, err8 := plus(expressions, &types, &functions, &bindings)

	if test1 != 3 {
		t.Errorf("(+ 1 2) = %v; want 3", test1)
	}
	if test2 != 6 {
		t.Errorf("(+ 1 2 3) = %v; want 6", test2)
	}
	if test3 != -3 {
		t.Errorf("(+ 1 -4) = %v; want -3", test3)
	}
	if test4 != 4 {
		t.Errorf("(+ 1 (+ 1 2)) = %v; want 4", test4)
	}
	if test5 != 5 {
		t.Errorf("(+ 1 (+ (+ 1 2) 1)) = %v; want 5", test5)
	}
	if test6 != 1 {
		t.Errorf("(+ 1) = %v; want 1", test6)
	}
	if err7 == nil {
		t.Errorf("(+) = %v; want invalid number of arguments.", err7)
	}
	if err8 == nil {
		t.Errorf("(+ $ #) = %v; want invalid syntax.", err8)
	}
}

func TestMinus(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(- 1)")
	test1, _ := minus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(- 2 1)")
	test2, _ := minus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(- 3 2 1)")
	test3, _ := minus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(- 1 2 3)")
	test4, _ := minus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(- 1 (- 1 2))")
	test5, _ := minus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(- 1 (- (- 1 2) 1))")
	test6, _ := minus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(-)")
	_, err7 := minus(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(- $ #)")
	_, err8 := minus(expressions, &types, &functions, &bindings)

	if test1 != -1 {
		t.Errorf("(- 1) = %v; want -1", test1)
	}
	if test2 != 1 {
		t.Errorf("(- 2 1) = %v; want 1", test2)
	}
	if test3 != 0 {
		t.Errorf("(- 3 2 1) = %v; want 0", test3)
	}
	if test4 != -4 {
		t.Errorf("(- 1 2 3) = %v; want -4", test4)
	}
	if test5 != 2 {
		t.Errorf("(- 1 (- 1 2)) = %v; want 2", test5)
	}
	if test6 != 3 {
		t.Errorf("(- 1 (- (- 1 2) 1)) = %v; want 3.", test6)
	}
	if err7 == nil {
		t.Errorf("(-) = %v; want invalid number of arguments.", err7)
	}
	if err8 == nil {
		t.Errorf("(- $ #) = %v; want invalid syntax.", err8)
	}
}

func TestTimes(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(* 1)")
	test1, _ := times(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(* 1 2)")
	test2, _ := times(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(* 1 2 -3)")
	test3, _ := times(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(* 1 2 0)")
	test4, _ := times(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(* 1 (* 1 2) (* 1 2))")
	test5, _ := times(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(* 1 (* (* 1 2) 2))")
	test6, _ := times(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(*)")
	_, err7 := times(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(* $ #)")
	_, err8 := times(expressions, &types, &functions, &bindings)

	if test1 != 1 {
		t.Errorf("(* 1) = %v; want 1", test1)
	}
	if test2 != 2 {
		t.Errorf("(* 1 2) = %v; want 2", test2)
	}
	if test3 != -6 {
		t.Errorf("(* 1 2 -3) = %v; want -6", test3)
	}
	if test4 != 0 {
		t.Errorf("(* 1 2 0) = %v; want 0", test4)
	}
	if test5 != 4 {
		t.Errorf("(* 1 (* 1 2) (* 1 2)) = %v; want 4", test5)
	}
	if test6 != 4 {
		t.Errorf("(* 1 (* (* 1 2) 2)) = %v; want 4.", test6)
	}
	if err7 == nil {
		t.Errorf("(*) = %v; want invalid number of arguments.", err7)
	}
	if err8 == nil {
		t.Errorf("(* $ #) = %v; want invalid syntax.", err8)
	}
}

func TestDivide(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(/ 1 2)")
	test1, _ := divide(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(/ 2 1)")
	test2, _ := divide(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(/ -3 2 1)")
	test3, _ := divide(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(/ 1 0)")
	_, err4 := divide(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(/ 6 (/ 2 1) (/ 3 1))")
	test5, _ := divide(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(/ 4 (/ (/ 2 1) 2))")
	test6, _ := divide(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(/ 1)")
	_, err7 := divide(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(/ $ #)")
	_, err8 := divide(expressions, &types, &functions, &bindings)

	if test1 != 0.5 {
		t.Errorf("(/ 1 2) = %v; want 0.5", test1)
	}
	if test2 != 2 {
		t.Errorf("(/ 2 1) = %v; want 2", test2)
	}
	if test3 != -1.5 {
		t.Errorf("(/ -3 2 1) = %v; want -1.5", test3)
	}
	if err4 == nil {
		t.Errorf("(/ 1 0) = %v; want cannot divide by zero", err4)
	}
	if test5 != 1 {
		t.Errorf("(/ 6 (/ 2 1) (/ 3 1)) = %v; want 1", test5)
	}
	if test6 != 4 {
		t.Errorf("(/ 4 (/ (/ 2 1) 2)) = %v; want 4.", test6)
	}
	if err7 == nil {
		t.Errorf("(/) = %v; want invalid number of arguments.", err7)
	}
	if err8 == nil {
		t.Errorf("(/ $ #) = %v; want invalid syntax.", err8)
	}
}
