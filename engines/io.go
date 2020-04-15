package engines

import (
	"../parse"
	"../structs"
	"../utils"
	"fmt"
)

func Translate(input string) (structs.List, error) {
	var expressions structs.List
	sep := []rune("() ")
	args := util.SplitWith(input, sep)
	args = util.RemoveMember(args, " ")
	args = args[1:] //shave off that opening paren
	expressions, _, err := parse.Transliterate(expressions, args, 0)

	/*
		if err = parse.Parse(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	*/
	return expressions, err
}

func Output(value interface{}) {
	if value != nil {
		switch value.(type) {
		case bool:
			if value == true {
				fmt.Println("T")
			} else {
				fmt.Println("NIL")
			}
		case float64:
			fmt.Println(value)
		case string:
			if value == "t" {
				fmt.Println("T")
			} else {
				fmt.Println(value.(string))
			}
		case structs.List:
			fmt.Print("(")
			structs.PrintList(value.(structs.List))
			fmt.Println(")")
		}
	}
}
