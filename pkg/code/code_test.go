package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expect   []byte
	}{
		// OpConstant's operand is two bytes wide meaning 65535 is the highest value that can be represented.
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpPop, []int{}, []byte{byte(OpPop)}},
		{OpSub, []int{}, []byte{byte(OpSub)}},
		{OpDiv, []int{}, []byte{byte(OpDiv)}},
		{OpMul, []int{}, []byte{byte(OpMul)}},
		{OpTrue, []int{}, []byte{byte(OpTrue)}},
		{OpFalse, []int{}, []byte{byte(OpFalse)}},
		{OpEQ, []int{}, []byte{byte(OpEQ)}},
		{OpNEQ, []int{}, []byte{byte(OpNEQ)}},
		{OpGT, []int{}, []byte{byte(OpGT)}},
		{OpMinus, []int{}, []byte{byte(OpMinus)}},
		{OpBang, []int{}, []byte{byte(OpBang)}},
	}

	for _, test := range tests {
		instruction := Make(test.op, test.operands...)

		if len(instruction) != len(test.expect) {
			t.Errorf(
				"instruction has wrong length. expected=%d, got=%d",
				len(instruction),
				len(test.expect),
			)
		}

		for i, b := range test.expect {
			if instruction[i] != b {
				t.Errorf("wrong byte at position %d. expected=%d, got=%d", i, b, instruction[i])
			}
		}
	}

}

func TestInstructionsString(t *testing.T) {
	tests := []struct {
		instructions []Instructions
		expected     string
	}{
		{
			[]Instructions{
				Make(OpConstant, 1),
				Make(OpConstant, 2),
				Make(OpConstant, 65534),
			}, "0000 OpConstant 1\n0003 OpConstant 2\n0006 OpConstant 65534"},
		{
			[]Instructions{
				Make(OpAdd),
				Make(OpConstant, 2),
				Make(OpConstant, 65534),
			}, "0000 OpAdd\n0001 OpConstant 2\n0004 OpConstant 65534"},
	}

	for _, test := range tests {
		concatted := Instructions{}
		for _, instruction := range test.instructions {
			concatted = append(concatted, instruction...)
		}

		if concatted.String() != test.expected {
			t.Fatalf(
				"instructions incorrectly formatted. expected=%s, got=%s",
				test.expected,
				concatted,
			)
		}
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
		{OpAdd, []int{}, 0},
		{OpPop, []int{}, 0},
		{OpMul, []int{}, 0},
		{OpDiv, []int{}, 0},
		{OpSub, []int{}, 0},
		{OpTrue, []int{}, 0},
		{OpFalse, []int{}, 0},
		{OpEQ, []int{}, 0},
		{OpNEQ, []int{}, 0},
		{OpGT, []int{}, 0},
		{OpMinus, []int{}, 0},
		{OpBang, []int{}, 0},
	}

	for _, test := range tests {
		instruction := Make(test.op, test.operands...)

		def, err := Lookup(byte(test.op))
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])
		if n != test.bytesRead {
			t.Fatalf("n wrong. expected=%d, got=%d", test.bytesRead, n)
		}

		for i, expected := range test.operands {
			if operandsRead[i] != expected {
				t.Errorf("operand wrong. expected=%d, got=%d", expected, operandsRead[i])
			}
		}
	}
}
