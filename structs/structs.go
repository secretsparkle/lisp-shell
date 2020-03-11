package structs

import (
	"fmt"
)

type Control struct {
	Args []string
	Body []string
}

type Function struct {
	Name     string
	Args     []string
	Bindings map[string]string
	Body     SExpression
}

type SExpression struct {
	Data        bool
	SExpression []interface{}
}

func PPSExpression(expression []interface{}) {
	fmt.Println(expression)
	fmt.Print("(")
	for i, s := range expression {
		switch s.(type) {
		case string:
			fmt.Print(s)
		default:
			fmt.Print("Found interface: ")
			fmt.Print(expression[i])
			//PPSExpression(expression[i])
		}
		fmt.Print(" ")
	}
	fmt.Println(")")
}
