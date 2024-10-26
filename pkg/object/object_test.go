package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello"}
	hello2 := &String{Value: "Hello"}
	diff1 := &String{Value: "World"}
	diff2 := &String{Value: "World"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	hello1 := &Boolean{Value: true}
	hello2 := &Boolean{Value: true}
	diff1 := &Boolean{Value: false}
	diff2 := &Boolean{Value: false}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	hello1 := &Integer{Value: 42}
	hello2 := &Integer{Value: 42}
	diff1 := &Integer{Value: 7}
	diff2 := &Integer{Value: 7}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
}

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		obj      Object
		isTruthy bool
	}{
		{&Integer{42}, true},
		{&Integer{0}, false},
		{&Integer{-42}, false},
		{&Boolean{true}, true},
		{&Boolean{false}, false},
		{&Null{}, false},
		{&Null{}, false},
		{&Array{}, false},
		{&Array{Elements: []Object{&Integer{42}}}, true},
		{&Hash{}, false},
		{
			&Hash{
				Pairs: map[HashKey]HashPair{
					{Type: INTEGER_OBJ, Value: 42}: {
						Key:   &Integer{42},
						Value: &Boolean{true},
					},
				},
			},
			true,
		},
	}

	for _, test := range tests {
		if IsTruthy(test.obj) != test.isTruthy {
			t.Errorf("%s is not %ty", test.obj.Inspect(), test.isTruthy)
		}
	}
}
