package main

import (
	"flag"
	"fmt"
	"os"
	"uniq"
)

// bool to int convertation
var conv_b2i = map[bool]int{false: 0, true: 1}

func main() {
	// Flag proccessing
	c_flg := flag.Bool("c", false, "\tprefix lines by the number of occurrences")
	d_flg := flag.Bool("d", false, "\tonly print duplicate lines, one for each group")
	u_flg := flag.Bool("u", false, "\tonly print unique lines")
	i_flg := flag.Bool("i", false, "\tignore differences in case when comparing")
	f_flg := flag.Int("f", 0, "\tavoid comparing the first N fields")
	s_flg := flag.Int("s", 0, "\tavoid comparing the first N characters")
	flag.Parse()

	// safety
	if *f_flg < 0 {
		*f_flg = 0
	}
	if *s_flg < 0 {
		*s_flg = 0
	}

	flags := map[string]int{
		"c": conv_b2i[*c_flg],
		"d": conv_b2i[*d_flg],
		"u": conv_b2i[*u_flg],
		"i": conv_b2i[*i_flg],
		"f": *f_flg,
		"s": *s_flg,
	}

	// Input/ Output files Arguments proccessing
	source := os.Stdin
	destination := os.Stdout
	var err error

	source_file := flag.Arg(0)
	destination_file := flag.Arg(1)
	fmt.Println("params: ", source_file, destination_file, flags)

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

	uniq.Uniq(source, destination, flags)
}
