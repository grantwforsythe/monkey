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
		instruction := Make(test.op, test.operands)

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
