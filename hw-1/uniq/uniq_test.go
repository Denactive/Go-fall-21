package uniq_test

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"uniq"
)

const root_cases = "test_cases"
const root_examples = "example_cases"
const root_answ = "test_answers"
const read_buffer_size = 128

// utility functions
func empty_syms_splitter(sym rune) bool {
	// важно сохранить в исходной строке пробелы и табуляции
	if sym == 0 || sym == 10 || sym == 13 {
		return true
	}
	return false
}

func string_slice_eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// test utilities
func toggle_flags(filename string, flags map[string]int) {
	flags_set := strings.TrimPrefix(strings.TrimSuffix(filename, ".txt"), "src_")
	for _, flag := range []string{"c", "d", "u", "i", "f", "s"} {
		flags[flag] = 0
		if strings.Contains(flags_set, flag) {
			flags[flag] = 1
			// test exception: skip 2 runes here
			if flag == "s" {
				flags[flag]++
			}
		}
	}
}

func walking_test(path string, info os.FileInfo, err error) error {
	// escaping root directory
	if info.IsDir() {
		return nil
	}

	// prepare. setting neccessaries
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("[TEST FAIL]\t Problem with test test file " + path)
		fmt.Println("[         ]\t", err)
		// continue executing walkFunk
		return nil
	}
	defer file.Close()

	flags := map[string]int{
		"c": 0,
		"d": 0,
		"u": 0,
		"i": 0,
		"f": 0,
		"s": 0,
	}

	r, w, err := os.Pipe()
	r_buffed := bufio.NewReader(r)

	toggle_flags(info.Name(), flags)
	fmt.Println("Test running with params: from [", info.Name(), "] to [", "stdout ]", flags)

	// act. function call
	uniq.Uniq(file, w, flags)

	// checking results. reading output from stdout
	output := make([]byte, read_buffer_size)
	_, err = r_buffed.Read(output)
	if err != nil {
		return err
	}

	// getting correct answer
	answer_filename := root_answ + strings.TrimPrefix(path, root_cases)
	file_answ, err := os.OpenFile(answer_filename, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("[TEST FAIL]\t Problem with test answer file " + answer_filename)
		fmt.Println("[         ]\t", err)
		// continue executing walkFunk
		return nil
	}
	defer file.Close()

	answer := make([]byte, read_buffer_size)
	_, err = file_answ.Read(answer)
	if err != nil {
		return err
	}

	// comparing
	output_str := strings.FieldsFunc(string(output), empty_syms_splitter)
	answer_str := strings.FieldsFunc(string(answer), empty_syms_splitter)

	if !string_slice_eq(output_str, answer_str) {
		return errors.New(
			fmt.Sprintf(
				"\"Uniq\" output is not equal to the correct one\n"+
					"Expected:\n%s\n------------\n"+
					"Your output:\n%s",
				string(answer), string(output),
			),
		)
	}
	return nil
}

// tests
func TestUniq(t *testing.T) {
	if err := filepath.Walk(root_cases, walking_test); err != nil {
		if strings.Contains(err.Error(), "test") {
			t.Fail()
		} else {
			t.Fatal(err)
		}
	}
}
