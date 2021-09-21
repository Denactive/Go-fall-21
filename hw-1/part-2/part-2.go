package main

import (
	"fmt"
	"mathexecutor"
	"os"
)

func main() {
	for _, input := range os.Args[1:] {
		fmt.Print(input + " => ")
		res, err := mathexecutor.ExecMathExp(input)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%v\n", res)
	}
}
