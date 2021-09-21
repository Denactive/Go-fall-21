package main

import (
	"flag"
	"fmt"
	"os"
	"uniq"
)

// bool to int convertation
var boolToInt = map[bool]int{false: 0, true: 1}

func main() {
	// Flag proccessing
	cFlag := flag.Bool("c", false, "\tprefix lines by the number of occurrences")
	dFlag := flag.Bool("d", false, "\tonly print duplicate lines, one for each group")
	uFlag := flag.Bool("u", false, "\tonly print unique lines")
	iFlag := flag.Bool("i", false, "\tignore differences in case when comparing")
	fFlag := flag.Int("f", 0, "\tavoid comparing the first N fields")
	sFlag := flag.Int("s", 0, "\tavoid comparing the first N characters")
	flag.Parse()

	// safety
	// f & s flags are not negative
	if *fFlag < 0 {
		*fFlag = 0
	}
	if *sFlag < 0 {
		*sFlag = 0
	}

	flags := map[string]int{
		"c": boolToInt[*cFlag],
		"d": boolToInt[*dFlag],
		"u": boolToInt[*uFlag],
		"i": boolToInt[*iFlag],
		"f": *fFlag,
		"s": *sFlag,
	}

	// Input/ Output files Arguments proccessing
	source := os.Stdin
	destination := os.Stdout
	var err error

	sourceFile := flag.Arg(0)
	destinationFile := flag.Arg(1)
	fmt.Println("params: ", sourceFile, destinationFile, flags)

	if sourceFile != "" {
		source, err = os.Open(sourceFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer source.Close()
	}

	if destinationFile != "" {
		// cross-platform file openning
		destination, err = os.OpenFile(destinationFile, os.O_WRONLY|os.O_CREATE, 666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()
	}

	uniq.Uniq(source, destination, flags)
}
