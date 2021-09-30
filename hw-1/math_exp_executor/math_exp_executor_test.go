package mathexecutor_test

import (
	"math"
	"mathexecutor"
	"strconv"
	"strings"
	"testing"
)

type TestCase struct {
	name   string
	input  string
	output string
}

var execMathExpTestingTable = []TestCase{
	TestCase{"simple plus int", "10+12", "22"},
	TestCase{"simple minus int", "10-12", "-2"},
	TestCase{"simple multiplication int", "10*12", "120"},
	TestCase{"simple division int", "9/3", "3"},

	TestCase{"simple plus float", "10.5+12.3", "22.8"},
	TestCase{"simple minus float", "10.5-1.2", "9.3"},
	TestCase{"simple multiplication float", "10.5*12.1", "127.05"},
	TestCase{"simple division float", "5/2", "2.5"},

	TestCase{"negative plus float", "-10.5+12.3", "1.8"},
	TestCase{"negative minus float", "-10.5-1.2", "-11.7"},
	TestCase{"negative multiplication float", "-10.5*12.1", "-127.05"},
	TestCase{"negative division float", "-5/2", "-2.5"},

	TestCase{"negative second operand plus float", "10.5+-12.3", "-1.8"},
	TestCase{"negative second operand minus float", "10.5--1.2", "11.7"},
	TestCase{"negative second operand multiplication float", "10.5*-12.1", "-127.05"},
	TestCase{"negative second operand division float", "5/-2", "-2.5"},

	TestCase{"both negative plus float", "-10.5+-12.3", "-22.8"},
	TestCase{"both negative minus float", "-10.5--1.2", "-9.3"},
	TestCase{"both negative multiplication float", "-10.5*-12.1", "127.05"},
	TestCase{"both negative division float", "-5/-2", "2.5"},

	TestCase{"bracket simple", "(12)", "12"},
	TestCase{"bracket plus", "(10+12)", "22"},
	TestCase{"bracket minus", "10-12", "-2"},
	TestCase{"bracket mul", "(10*5)", "50"},
	TestCase{"bracket div", "(9/2)", "4.5"},
	TestCase{"bracket minus", "(-9)", "-9"},

	TestCase{"pririty +*", "10+12*7", "94"},
	TestCase{"pririty +/", "10+9/2", "14.5"},
	TestCase{"pririty brackets", "(10+12)*7", "154"},
	TestCase{"pririty * combo", "3*4*5", "60"},
	TestCase{"pririty / combo", "12/2/6", "1"},
	TestCase{"pririty + combo", "3+4+5", "12"},
	TestCase{"pririty - combo", "3-4-5.5", "-6.5"},

	TestCase{"zero division", "12/0", "infinity"},
	TestCase{"zero multiplication", "12*0", "0"},

	TestCase{"complex", "(2+1)*(1+9/2+2)", "22.5"},
	TestCase{"bracket spam", "(1)+(2)-(3)+(4)-(5)+(6)-(7)+(8)+(0)-(0)+(0)-(0)+(0)-(0)", "6"},
	TestCase{"inner bracket spam", "(((((((((0)))))))))", "0"},

	TestCase{"error test empty expression", "", ""},
	TestCase{"error test no close bracket", "(12+13", ""},
	TestCase{"error test no open bracket", "12+13)", ""},
	TestCase{"error test no left operand", "*12", ""},
	TestCase{"error test no right operand", "175+", ""},
	TestCase{"error test no operands in operational sequence", "175+*", ""},
	TestCase{"error test no operand between brackets", "17(19+2)", ""},
	TestCase{"error test unknown operation", "2^3", ""},
	TestCase{"error test unknown symbol", "x*x", ""},
	TestCase{"error test complicated 1", "(19+1.0)*12-+----7", ""},
	TestCase{"error test complicated 2", "0+0-000+12*(12-)", ""},
	TestCase{"error test complicated 3", "(125-123*2)-11+", ""},
}

var prepExpTestingTable []TestCase = []TestCase{
	TestCase{"trim spaces", "    10    +     12       ", "10+12"},
	TestCase{"trim tabs", "\t10\t+\t12\t", "10+12"},
	TestCase{"trim '\\r'", "\r10\r+\r12\r", "10+12"},
	TestCase{"trim '\\n\\r'", "\n\r10\n\r+\n\r12\n\r", "10+12"},
	TestCase{"numbers", "0123456789", "0123456789"},
	TestCase{"operations", "+--**///++++", "+--**///++++"},
	TestCase{"brackets", "		(\t)\n)()()()((           )", "())()()()(()"},
	TestCase{"dots", ".    .     .\t.\n.\n\r.", "......"},

	TestCase{"error test letters", "qwertyuiopasdfghjklzxcvbnm", ""},
	TestCase{"error test unsupported opearions", "+///-++++ * * * ^", ""},
	TestCase{"error test letter in equation", "15 * 2 * (15-2.8) * 2^x", ""},
}

const EPS = 1.0e-6

func TestExecMathExp(t *testing.T) {
	for _, testCase := range execMathExpTestingTable {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := mathexecutor.ExecMathExp(testCase.input)
			if err != nil {
				if strings.Contains(testCase.name, "error") {
					return
				}
				t.Error(err)
			}
			if strings.Contains(testCase.name, "error") {
				t.Error("need to raise an error")
			}
			output, _ := strconv.ParseFloat(testCase.output, 64)
			if math.Abs(res-output) > EPS {
				t.Errorf("got %v, want %v", res, testCase.output)
			}
		})
	}
}

func TestPrepareExp(t *testing.T) {
	for _, testCase := range prepExpTestingTable {
		t.Run(testCase.name, func(t *testing.T) {
			output, err := mathexecutor.PrepareExp(testCase.input)
			if err != nil {
				if strings.Contains(testCase.name, "error") {
					return
				}
				t.Error(err)
			}
			if strings.Contains(testCase.name, "error") {
				t.Error("need to raise an error")
			}
			if output != testCase.output {
				t.Errorf("got %v, want %v", output, testCase.output)
			}
		})
	}
}
