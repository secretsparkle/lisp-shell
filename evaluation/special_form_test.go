package evaluation

import (
	"../engines"
	"../structs"
	"testing"
)

func TestAnd(t *testing.T) {
	types, functions, bindings := structs.Maps()

	expressions, _ := engines.Translate("(and t)")
	test1, _ := and(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(and t t)")
	test2, _ := and(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(and t 3)")
	test3, _ := and(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(and (and 3 t t))")
	test4, _ := and(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(and (and 3 t t) t)")
	test5, _ := and(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(and (and 3 t t) (and 3 t t))")
	test6, _ := and(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(and t '(1 2 3))")
	test7, _ := and(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(and nil)")
	test8, _ := and(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(and t nil)")
	test9, _ := and(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(and nil t)")
	test10, _ := and(expressions, &types, &functions, &bindings)

	if test1 != "t" {
		t.Errorf("(and t) = %v; want t", test1)
	}
	if test2 != "t" {
		t.Errorf("(and t t) = %v; want false", test2)
	}
	if test3 != "3" {
		t.Errorf("(and t 3) = %v; want 3", test3)
	}
	if test4 != "t" {
		t.Errorf("(and (and 3 t t)) = %v; want t", test4)
	}
	if test5 != "t" {
		t.Errorf("(and (and 3 t t) t) = %v; want t", test5)
	}
	if test6 != "t" {
		t.Errorf("(and (and 3 t t) (and 3 t t)) = %v; want t", test6)
	}
	if test7 != "(1 2 3)" {
		t.Errorf("(and t '(1 2 3)) = %v; want (1 2 3)", test7)
	}
	if test8 != false {
		t.Errorf("(and nil) = %v; want nil", test8)
	}
	if test9 != false {
		t.Errorf("(and t nil) = %v; want nil", test9)
	}
	if test10 != false {
		t.Errorf("(and nil t) = %v; want nil", test10)
	}
}
