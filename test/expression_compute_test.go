package test

import (
	"github.com/stanleychan/go-utils/math"
	"testing"
)

func TestExpressionCompute(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{"Simple Addition", "-3+(2.5+3.5)*2", 9, false},
		{"Simple Max", "max(2,3)", 3, false},
		{"Min with Multiple Arguments", "min(1,2)", 1, false},
		{"Nested Max", "max(1,max(2,3))", 3, false},
		{"Complex Expression(include min)", "2+min(3,4)*2", 8, false},

		{"Negative Number Multiple Positive Number Expression ", "3*(-2)", -6, false},
		{"Negative Numbers Multiple Expression", "-3*(-2)+max(5,6)", 12, false},
		{"Max with Negative Numbers", "max(-1,-2)", -1, false},
		{"Min with Negative Numbers", "min(-8,-10)", -10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := math.Computing(tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("Computing() error = %v, wantErr %v", err, tt.hasError)
				return
			}
			if !tt.hasError && result != tt.expected {
				t.Errorf("Computing() = %v, want %v", result, tt.expected)
			}
		})
	}
}
