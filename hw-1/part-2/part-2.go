package main

import (
	"bufio"
	"fmt"
	"io"
	"mathexecutor"
	"os"
)

func wellcome() (string, error) {
	fmt.Println("Enter expression >>")
	var input string
	var err error
	input, err = bufio.NewReader(os.Stdin).ReadString('\n')

	switch {
	case err == io.EOF:
		return "", err
	case err != nil:
		fmt.Println(err)
		return "", err
	default:
		return mathexecutor.PrepareExp(input)
	}
}

func main() {
	for input, err := wellcome(); err != io.EOF; input, err = wellcome() {
		fmt.Println(input)
		if input != "" {
			res, err := mathexecutor.ExecMathExp(input)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Result: %v\n", res)
		}
	}
}
