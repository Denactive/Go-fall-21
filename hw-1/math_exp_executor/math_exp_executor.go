package mathexecutor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// const DEBUG = true
const DEBUG = false

type Result struct {
	acc  float64 // accumulator - buffer for current result
	rest string  // a rest part of expression string which will be proccessed later
}

func PrepareExp(input string) (string, error) {
	// unknown - a list of unsuported symbols. Will be seen in the error
	var unknown string
	var err error
	var trimmed []rune = []rune(strings.Join(strings.Fields(input), ""))

	for _, sym := range trimmed {
		if !unicode.IsDigit(sym) && !strings.ContainsRune("()+-*/.", sym) {
			if !strings.ContainsRune(unknown, sym) {
				unknown += string(sym)
			}
		}
	}
	if unknown != "" {
		err = errors.New("Unknown characters: " + unknown)
		return "", err
	}
	return string(trimmed), err
}

func ExecMathExp(exp string) (float64, error) {
	var res Result
	var err error
	if exp == "" {
		return 0, errors.New("Empty expression")
	}
	res, err = add_sub_handler(exp)

	if err != nil {
		return res.acc, err
	}
	if res.rest != "" {
		return res.acc, errors.New("Cannot fully parse. Stoped here: [ " + res.rest + " ]")
	}

	return res.acc, nil
}

// priority 1
func add_sub_handler(exp string) (Result, error) {
	if DEBUG {
		fmt.Println("add & sub handler:", exp)
	}

	var left_operand Result
	var err error
	var acc_buffer float64

	left_operand, err = mul_div_handler(exp)
	if err != nil {
		return left_operand, err
	}
	acc_buffer = left_operand.acc

	for {
		if DEBUG {
			fmt.Println("add sub acc_buffer:", acc_buffer, "operand:", left_operand)
		}
		if left_operand.rest == "" {
			break
		}
		var sign string = string(left_operand.rest[0])
		// if next sym is '+' or '-' this handler will be recalled
		if !strings.ContainsAny("+-", sign) {
			break
		}

		// get second number if sign is found
		// mul div have more priority
		left_operand, err = mul_div_handler(left_operand.rest[1:])
		if err != nil {
			return left_operand, err
		}
		if sign == "+" {
			acc_buffer += left_operand.acc
		} else {
			acc_buffer -= left_operand.acc
		}
	}
	return Result{acc_buffer, left_operand.rest}, nil
}

// priority 2
func mul_div_handler(exp string) (Result, error) {
	if DEBUG {
		fmt.Println("mul & div handler:", exp)
	}

	var left_operand Result
	var err error
	var acc_buffer float64

	left_operand, err = bracket_handler(exp)
	if err != nil {
		return left_operand, err
	}
	acc_buffer = left_operand.acc

	// iterations on constructions like n*n*n\n...
	for {
		// stop iterations if nothing left
		if left_operand.rest == "" {
			return left_operand, nil
		}

		var sign string
		sign = string(left_operand.rest[0])
		// if all '*', '\' signed are passed -> need to perform, other -> return
		if !strings.ContainsAny("/*", sign) {
			return left_operand, nil
		}

		// get second number
		// bracket expression has more priority
		var right_operand Result
		right_operand, err = bracket_handler(left_operand.rest[1:])
		if err != nil {
			return right_operand, err
		}

		// perform operation
		if sign == "/" {
			acc_buffer /= right_operand.acc
		}
		if sign == "*" {
			acc_buffer *= right_operand.acc
		}

		// continue iterations for another '*' '/' signs
		left_operand = Result{acc_buffer, right_operand.rest}
	}
}

// priority 3
func bracket_handler(exp string) (Result, error) {
	if DEBUG {
		fmt.Println("bracket handler:", exp)
	}

	if exp != "" && exp[0] == '(' {
		// priority reset
		var res Result
		var err error
		res, err = add_sub_handler(string(exp[1:]))
		if err != nil {
			return res, err
		}
		if res.rest != "" && res.rest[0] == ')' {
			res.rest = string(res.rest[1:])
		} else {
			err = errors.New("no close bracket")
		}
		return res, err
	}
	// if it is not a bracket construction then it is a number
	return num_handler(exp)
}

// none priority
func num_handler(exp string) (Result, error) {
	if DEBUG {
		fmt.Println("num handler:", exp)
	}

	var i int
	var exp_runes []rune = []rune(exp)

	// search for not n=[0-9] and '.'
	// case nnnn.nnnn.nn is invalid and ParseFloat will return error
	for ; i < len(exp_runes); i++ {
		var sym rune = exp_runes[i]

		// check for negative
		if sym == '-' && i == 0 {
			continue
		}

		if !unicode.IsDigit(sym) && sym != '.' {
			break
		}
	}

	var num float64
	var err error
	num, err = strconv.ParseFloat(string(exp_runes[:i]), 64)

	if DEBUG {
		fmt.Println("\tnum:", num, "[] error:", err, "]")
	}
	var rest string = string(exp_runes[i:])
	if DEBUG {
		fmt.Printf("\texp_runes: '%s'[%d:] == '%s'\n", string(exp_runes), i, rest)
	}
	return Result{num, rest}, err
}
