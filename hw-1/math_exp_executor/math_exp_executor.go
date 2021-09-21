package mathexecutor

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type Result struct {
	accumulator float64 // buffer for current result
	rest        string  // a rest part of expression string which will be proccessed later
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

// recursive trigger algorithm implementation
func ExecMathExp(exp string) (float64, error) {
	var res Result
	var err error
	if exp == "" {
		return 0, errors.New("Empty expression")
	}
	res, err = addSubHandler(exp)

	if err != nil {
		return res.accumulator, err
	}
	if res.rest != "" {
		return res.accumulator, errors.New("Cannot fully parse. Stoped here: [ " + res.rest + " ]")
	}

	return res.accumulator, nil
}

// priority 1
func addSubHandler(exp string) (Result, error) {
	var leftOperand Result
	var err error
	var accumulatorBuffer float64

	leftOperand, err = mulDivHandler(exp)
	if err != nil {
		return leftOperand, err
	}
	accumulatorBuffer = leftOperand.accumulator

	for {
		if leftOperand.rest == "" {
			break
		}
		var sign string = string(leftOperand.rest[0])
		// if next sym is '+' or '-' this handler will be recalled
		if !strings.ContainsAny("+-", sign) {
			break
		}

		// get second number if sign is found
		// mul div have more priority
		leftOperand, err = mulDivHandler(leftOperand.rest[1:])
		if err != nil {
			return leftOperand, err
		}
		if sign == "+" {
			accumulatorBuffer += leftOperand.accumulator
			continue
		}
		accumulatorBuffer -= leftOperand.accumulator
	}
	return Result{accumulatorBuffer, leftOperand.rest}, nil
}

// priority 2
func mulDivHandler(exp string) (Result, error) {
	var leftOperand Result
	var err error
	var accumulatorBuffer float64

	leftOperand, err = bracketHandler(exp)
	if err != nil {
		return leftOperand, err
	}
	accumulatorBuffer = leftOperand.accumulator

	// iterations on constructions like n*n*n\n...
	for {
		// stop iterations if nothing left
		if leftOperand.rest == "" {
			return leftOperand, nil
		}

		var sign string
		sign = string(leftOperand.rest[0])
		// if all '*', '\' signed are passed -> need to perform, other -> return
		if !strings.ContainsAny("/*", sign) {
			return leftOperand, nil
		}

		// get second number
		// bracket expression has more priority
		var rightOperand Result
		rightOperand, err = bracketHandler(leftOperand.rest[1:])
		if err != nil {
			return rightOperand, err
		}

		// perform operation
		if sign == "/" {
			accumulatorBuffer /= rightOperand.accumulator
		}
		if sign == "*" {
			accumulatorBuffer *= rightOperand.accumulator
		}

		// continue iterations for another '*' '/' signs
		leftOperand = Result{accumulatorBuffer, rightOperand.rest}
	}
}

// priority 3
func bracketHandler(exp string) (Result, error) {
	if exp != "" && exp[0] == '(' {
		// priority reset
		var res Result
		var err error
		res, err = addSubHandler(string(exp[1:]))
		if err != nil {
			return res, err
		}
		if res.rest != "" && res.rest[0] == ')' {
			res.rest = string(res.rest[1:])
			return res, err
		}
		err = errors.New("no close bracket")
		return res, err
	}
	// if it is not a bracket construction then it is a number
	return numHandler(exp)
}

// none priority
func numHandler(exp string) (Result, error) {
	var i int
	var expRunes []rune = []rune(exp)

	// search for not n=[0-9] and '.'
	// case nnnn.nnnn.nn is invalid and ParseFloat will return error
	for ; i < len(expRunes); i++ {
		var sym rune = expRunes[i]

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
	num, err = strconv.ParseFloat(string(expRunes[:i]), 64)

	return Result{num, string(expRunes[i:])}, err
}
