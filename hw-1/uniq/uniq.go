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
	prev_word := ""
	counter := 0

	for scanner.Scan() {
		if uniq_ne(scanner.Text(), prev_word, flags) {
			if flags["c"] == 0 && flags["d"] == 0 && flags["u"] == 0 {
				// standart case (pre-printing)
				write_uniq_msg(writer, []byte(scanner.Text()), 0)
			} else {
				// case flags["d"] == 1 || flags["u"] == 1 || flags["c"] == 1:
				// post-printing
				post_write(writer, flags, prev_word, counter)
				counter = 1
			}
			prev_word = scanner.Text()
		} else {
			counter++
		}
	}

	// -c, -d, -u => delayed print
	post_write(writer, flags, prev_word, counter)

	writer.Write([]byte(NEWLINE_SEPARATOR))
	writer.Flush()
}

// returns true if lines are not equal considering -i, -f, -s flags
func uniq_ne(str1 string, str2 string, flags map[string]int) bool {
	if flags["f"] != 0 {
		words1 := strings.Fields(str1)
		if len(words1) == 0 {
			str1 = ""
		} else {
			// skip all symbols till [f] fields
			invisible_word1 := words1[flags["f"]-1]
			str1 = string(str1[strings.Index(str1, invisible_word1)+len(invisible_word1):])
		}
		words2 := strings.Fields(str2)
		if len(words2) == 0 {
			str2 = ""
		} else {
			invisible_word2 := (strings.Fields(str2))[flags["f"]-1]
			str2 = string(str2[strings.Index(str2, invisible_word2)+len(invisible_word2):])
		}
	}

	// cut off [s] symbols
	if flags["s"] != 0 {
		if flags["s"] >= len(str1) {
			str1 = ""
		} else {
			str1 = string(([]rune(str1))[flags["s"]:])
		}
		if flags["s"] >= len(str2) {
			str2 = ""
		} else {
			str2 = string(([]rune(str2))[flags["s"]:])
		}
	}

	if flags["i"] == 1 {
		str1 = strings.ToLower(str1)
		str2 = strings.ToLower(str2)
	}

	// fmt.Println(" => after: |", str1, "|", str2)

	return str1 != str2
}

// prints saved word if -c or -u or -d are set
func post_write(w *bufio.Writer, flags map[string]int, buffer string, counter int) {
	if uniq_flag_process(flags, counter) {
		if flags["c"] == 1 {
			write_uniq_msg(w, []byte(buffer), counter)
		} else {
			write_uniq_msg(w, []byte(buffer), 0)
		}
	}
}

// returns true if flags and counter combination is appropriate and -c or -u or -d are set
func uniq_flag_process(flags map[string]int, counter int) bool {
	// d - print if there is at least 1 repeat
	// u - print if counter is 1
	// in cases -du and -ud no print
	if !(flags["d"] == 1 && flags["u"] == 1) && (flags["d"] == 1 && counter > 1 || flags["u"] == 1 && counter == 1 || flags["c"] == 1 && flags["d"] == 0 && flags["u"] == 0) {
		return true
	}
	return false
}

// writes msg to io.Writer as original uniq util
func write_uniq_msg(w *bufio.Writer, str []byte, counter int) {
	// empty line case
	if len(str) == 0 {
		return
	}
	msg := create_uniq_msg(str, counter)
	w.Write(msg)
	w.Flush()
}

// returns msg as original uniq util
func create_uniq_msg(str []byte, counter int) []byte {
	msg := append(str, 10) // '\n'
	if counter != 0 {
		msg = append([]byte("\t"+strconv.Itoa(counter)+" "), msg...)
	}
	return msg
}
