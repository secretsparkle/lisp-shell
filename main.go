package main

import (
	"./engines"
	"./evaluation"
	"./structs"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	symbols, functions, bindings := structs.Maps()
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if strings.TrimLeft(input, " ") == ("\n") {
			continue
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		expressions, err := engines.Translate(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		value, err := evaluation.EvaluateFunction(expressions, &symbols, &functions, &bindings)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		engines.Output(value)
	}
}
