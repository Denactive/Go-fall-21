package uniq

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Main logic function
func Uniq(source *os.File, destination *os.File, flags map[string]int) {
	scanner := bufio.NewScanner(source)
	writer := bufio.NewWriter(destination)
	prevWord := ""
	counter := 0

	for scanner.Scan() {
		if uniqNE(scanner.Text(), prevWord, flags) {
			switch {
			case flags["c"] == 0 && flags["d"] == 0 && flags["u"] == 0:
				// standart case (pre-printing)
				writeUniqMsg(writer, []byte(scanner.Text()), 0)
			default:
				// case flags["d"] == 1 || flags["u"] == 1 || flags["c"] == 1:
				// post-printing
				postWrite(writer, flags, prevWord, counter)
				counter = 1
			}
			prevWord = scanner.Text()
			continue
		}
		counter++
	}

	// -c, -d, -u => delayed print
	postWrite(writer, flags, prevWord, counter)

	writer.Write([]byte(NEWLINE_SEPARATOR))
	writer.Flush()
}

// returns true if lines are not equal considering -i, -f, -s flags
func uniqNE(str1 string, str2 string, flags map[string]int) bool {
	if flags["f"] != 0 {
		words1 := strings.Fields(str1)
		if len(words1) == 0 {
			str1 = ""
		}
		if len(words1) != 0 {
			// skip all symbols till [f] fields
			invisibleWord1 := words1[flags["f"]-1]
			str1 = string(str1[strings.Index(str1, invisibleWord1)+len(invisibleWord1):])
		}
		words2 := strings.Fields(str2)
		if len(words2) == 0 {
			str2 = ""
		}
		if len(words2) != 0 {
			invisibleWord2 := (strings.Fields(str2))[flags["f"]-1]
			str2 = string(str2[strings.Index(str2, invisibleWord2)+len(invisibleWord2):])
		}
	}

	// cut off [s] symbols
	if flags["s"] != 0 {
		if flags["s"] >= len(str1) {
			str1 = ""
		}
		if flags["s"] < len(str1) {
			str1 = string(([]rune(str1))[flags["s"]:])
		}
		if flags["s"] >= len(str2) {
			str2 = ""
		}
		if flags["s"] < len(str2) {
			str2 = string(([]rune(str2))[flags["s"]:])
		}
	}

	if flags["i"] == 1 {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	return str1 != str2
}

// prints saved word if -c or -u or -d are set
func postWrite(w *bufio.Writer, flags map[string]int, buffer string, counter int) {
	if uniqFlagProcess(flags, counter) {
		if flags["c"] == 1 {
			writeUniqMsg(w, []byte(buffer), counter)
			return
		}
		writeUniqMsg(w, []byte(buffer), 0)
	}
}

// returns true if flags and counter combination is appropriate and -c or -u or -d are set
func uniqFlagProcess(flags map[string]int, counter int) bool {
	// d - print if there is at least 1 repeat
	// u - print if counter is 1
	// in cases -du and -ud no print
	if !(flags["d"] == 1 && flags["u"] == 1) && (flags["d"] == 1 && counter > 1 ||
		flags["u"] == 1 && counter == 1 ||
		flags["c"] == 1 && flags["d"] == 0 && flags["u"] == 0) {
		return true
	}
	return false
}

// writes msg to io.Writer as original uniq util
func writeUniqMsg(w *bufio.Writer, str []byte, counter int) {
	// empty line case
	if len(str) == 0 {
		return
	}
	msg := createUniqMsg(str, counter)
	w.Write(msg)
	w.Flush()
}

// returns msg as original uniq util
func createUniqMsg(str []byte, counter int) []byte {
	msg := append(str, 10) // '\n'
	if counter != 0 {
		msg = append([]byte("\t"+strconv.Itoa(counter)+" "), msg...)
	}
	return msg
}
