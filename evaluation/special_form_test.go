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

	expressions, _ = engines.Translate("(and nil)")
	test2, _ := and(expressions, &types, &functions, &bindings)

	expressions, _ = engines.Translate("(and (and 3 t t))")
	test3, _ := and(expressions, &types, &functions, &bindings)

	//expressions, _ = engines.Translate("(and t)")
	//test4, _ := and(expressions, &types, &functions, &bindings)

	//expressions, _ = engines.Translate("(and t)")
	//test5, _ := and(expressions, &types, &functions, &bindings)

	//expressions, _ = engines.Translate("(and t)")
	//test6, _ := and(expressions, &types, &functions, &bindings)

	if test1 != "t" {
		t.Errorf("(and t) = %v; want t", test1)
	}
	if test2 != false {
		t.Errorf("(and nil) = %v; want false", test2)
	}
	if test3 != "t" {
		t.Errorf("(and (and 3 t t)) = %v; want t", test3)
	}
}
