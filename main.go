package main

import (
	"./parse"
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		// Read the keyboard input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if err = parse.Parse(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Handle the execution of the input.
		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execInput(input string) error {

	// Remove the leading parentheses
	input = strings.TrimPrefix(input, "(")

	// Remove the trailing parenthese and newline character.
	input = strings.TrimSuffix(input, ")\n")

	// Split the input to separate the command and the arguments.
	args := strings.Split(input, " ")

	// Check for built-in commands
	switch args[0] {
	case "+":
		if len(args) == 1 {
			return errors.New("Invalid number of arguments.")
		}
		sum := 0.0
		for _, number := range args[1:] {
			if n, err := strconv.ParseFloat(number, 64); err == nil {
				sum += n
			} else {
				return errors.New("Only numbers can be added.")
			}
		}
		fmt.Println(sum)
		return nil
	case "-":
		if len(args) == 1 {
			return errors.New("Invalid number of arguments.")
		} else if len(args) == 2 {
			difference, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return errors.New("Only numbers can be subtracted.")
			}
			difference = 0 - difference
			fmt.Println(difference)
			return nil
		} else {
			difference, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return errors.New("Only numbers can be subtracted.")
			}
			for _, number := range args[2:] {
				if n, err := strconv.ParseFloat(number, 64); err == nil {
					difference -= n
				} else {
					return errors.New("Only numbers can be subtracted.")
				}
			}
			fmt.Println(difference)
			return nil
		}
	case "*":
		if len(args) == 1 {
			return errors.New("Invalid number of arguments.")
		}
		product := 1.0
		for _, number := range args[1:] {
			if n, err := strconv.ParseFloat(number, 64); err == nil {
				product *= n
			} else {
				return errors.New("Only numbers can be multiplied.")
			}
		}
		fmt.Println(product)
		return nil
	case "/":
		if len(args) < 3 {
			return errors.New("Invalid number of arguments.")
		}
		numerator, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return errors.New("Only numbers can be divided.")
		}
		for _, number := range args[2:] {
			if n, err := strconv.ParseFloat(number, 64); err == nil {
				numerator /= n
			} else {
				return errors.New("Only numbers can be divided.")
			}
		}
		fmt.Println(numerator)
		return nil
	case "cd":
		// 'cd' to home dir with empty path not yet supported.
		if len(args) < 2 {
			return errors.New("path required")
		}
		// Change the directory and return the error.
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	}

	// Pass the program and the arguments separately.
	cmd := exec.Command(args[0], args[1:]...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command
	return cmd.Run()
}
