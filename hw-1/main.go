package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

// Main logic function
func uniq(source *os.File, destination *os.File) {
	scanner := bufio.NewScanner(source)
	writer := bufio.NewWriter(destination)
	prev_word := ""

	for scanner.Scan() {
		if scanner.Text() != prev_word {
			_, err := writer.Write(append(scanner.Bytes(), 10))
			if err != nil {
				fmt.Println("w failed to write:", err)
			}
			err = writer.Flush()
			if err != nil {
				fmt.Println("w failed to flush:", err)
			}
			prev_word = scanner.Text()
		}
	}
}

func main() {
	// Flag proccessing
	c_flg := flag.Bool("c", false, "\tprefix lines by the number of occurrences")
	d_flg := flag.Bool("d", false, "\tonly print duplicate lines, one for each group")
	u_flg := flag.Bool("u", false, "\tonly print unique lines")
	i_flg := flag.Bool("i", false, "\tignore differences in case when comparing")
	f_flg := flag.Int("f", 0, "\tavoid comparing the first N fields")
	s_flg := flag.Int("s", 0, "\tavoid comparing the first N characters")
	flag.Parse()
	fmt.Println("c flag: ", *c_flg)
	fmt.Println("d flag: ", *d_flg)
	fmt.Println("u flag: ", *u_flg)
	fmt.Println("i flag: ", *i_flg)
	fmt.Println("f flag: ", *f_flg)
	fmt.Println("s flag: ", *s_flg)

	// Input/ Output files Arguments proccessing
	source := os.Stdin
	destination := os.Stdout
	var err error

	source_file := flag.Arg(0)
	destination_file := flag.Arg(1)
	fmt.Println("params: ", source_file, destination_file)

	if source_file != "" {
		source, err = os.Open(source_file)
		if err != nil {
			panic(err)
		}
		defer source.Close()
	}

	if destination_file != "" {
		// cross-platform file openning
		destination, err = os.OpenFile(destination_file, os.O_WRONLY|os.O_CREATE, 666)
		if err != nil {
			panic(err)
		}
		defer destination.Close()
	}

	uniq(source, destination)
}
