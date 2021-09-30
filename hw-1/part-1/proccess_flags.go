package flags

import "flag"

func proccessFlags() map[string]int {
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

	return map[string]int{
		"c": boolToInt[*cFlag],
		"d": boolToInt[*dFlag],
		"u": boolToInt[*uFlag],
		"i": boolToInt[*iFlag],
		"f": *fFlag,
		"s": *sFlag,
	}
}
