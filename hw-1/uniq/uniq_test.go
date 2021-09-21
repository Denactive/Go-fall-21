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

const rootCases = "test_cases"
const rootExamples = "example_cases"
const rootAnswers = "test_answers"
const readBufferSize = 128

// utility functions
func emptySymsSplitter(sym rune) bool {
	// важно сохранить в исходной строке пробелы и табуляции
	if sym == 0 || sym == 10 || sym == 13 {
		return true
	}
	return false
}

func stringSliceEQ(a, b []string) bool {
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
func toggleFlags(filename string, flags map[string]int) {
	flagsSet := strings.TrimPrefix(strings.TrimSuffix(filename, ".txt"), "src_")
	for _, flag := range []string{"c", "d", "u", "i", "f", "s"} {
		flags[flag] = 0
		if strings.Contains(flagsSet, flag) {
			flags[flag] = 1
			// test exception: skip 2 runes here
			if flag == "s" {
				flags[flag]++
			}
		}
	}
}

func walkingTest(path string, info os.FileInfo, err error) error {
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
	rBuffed := bufio.NewReader(r)

	toggleFlags(info.Name(), flags)
	fmt.Println("Test running with params: from [", info.Name(), "] to [", "stdout ]", flags)

	// act. function call
	uniq.Uniq(file, w, flags)

	// checking results. reading output from stdout
	output := make([]byte, readBufferSize)
	_, err = rBuffed.Read(output)
	if err != nil {
		return err
	}

	// getting correct answer
	answerFilename := rootAnswers + strings.TrimPrefix(path, rootCases)
	answerFile, err := os.OpenFile(answerFilename, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("[TEST FAIL]\t Problem with test answer file " + answerFilename)
		fmt.Println("[         ]\t", err)
		// continue executing walkFunk
		return nil
	}
	defer file.Close()

	answer := make([]byte, readBufferSize)
	_, err = answerFile.Read(answer)
	if err != nil {
		return err
	}

	// comparing
	outputStr := strings.FieldsFunc(string(output), emptySymsSplitter)
	answerStr := strings.FieldsFunc(string(answer), emptySymsSplitter)

	if !stringSliceEQ(outputStr, answerStr) {
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
	if err := filepath.Walk(rootCases, walkingTest); err != nil {
		if strings.Contains(err.Error(), "test") {
			t.Fail()
		} else {
			t.Fatal(err)
		}
	}
}
