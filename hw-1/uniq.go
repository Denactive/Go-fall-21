package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Main logic function
func Uniq(source *os.File, destination *os.File, flags map[string]int) {
	scanner := bufio.NewScanner(source)
	writer := bufio.NewWriter(destination)
	prev_word := ""
	counter := 0

	for scanner.Scan() {
		if scanner.Text() != prev_word {
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
	if !(flags["d"] == 1 && flags["u"] == 1) && (flags["d"] == 1 && counter > 1 || flags["u"] == 1 && counter == 1 || flags["c"] == 1) {
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

	Uniq(source, destination, flags)
}
