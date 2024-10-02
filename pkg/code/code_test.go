package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expect   []byte
	}{
		// OpConstant's operand is two bytes wide meaning 65535 is the highest value that can be represented.
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
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
	instruction := []Instructions{
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65534),
	}

	expected := "0000 OpConstant 1\n0003 OpConstant 2\n0006 OpConstant 65534"

	concatted := Instructions{}
	for _, instruction := range instruction {
		concatted = append(concatted, instruction...)
	}

	if concatted.String() != expected {
		t.Fatalf("instructions incorrectly formatted. expected=%s, got=%s", expected, concatted)
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
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
