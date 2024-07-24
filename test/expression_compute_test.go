package test

import (
	"fmt"
	"github.com/stanleychan/go-utils/math"
	"testing"
)

func TestExpressionCompute(t *testing.T) {
	expression := "(2.1*(1.4*3+1)^2-4*1.7)/2.3"
	result, err := math.Computing(expression)
	if err != nil {
		fmt.Println("test expression '(2.1*(1.4*3+1)^2-4*1.7)/2.3', expected value is 21.732173913043475..., calculated error is :", err)
	} else {
		fmt.Println("test expression '(2.1*(1.4*3+1)^2-4*1.7)/2.3', expected value is 21.732173913043475..., calculated value is :", result)
	}
}
