package math

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

const validStr = "max,in0123456789.+-*/()^%"

var precedence = map[string]int{
	"+": 1, "-": 1,
	"*": 2, "/": 2, "%": 2,
	"^":   3,
	"max": 4, "min": 4,
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isOperator(s string) bool {
	//_, ok := precedence[s]
	//return ok
	return s == "+" || s == "-" || s == "*" || s == "/" || s == "^" || s == "%"
}

func isFunction(s string) bool {
	return s == "max" || s == "min"
}

func tokenize(expression string) []string {
	var tokens []string
	var current strings.Builder
	for i := 0; i < len(expression); i++ {
		ch := expression[i]
		switch {
		case unicode.IsDigit(rune(ch)) || ch == '.':
			current.WriteByte(ch)
		case ch == '-' && (i == 0 || expression[i-1] == '(' || isOperator(string(expression[i-1])) || expression[i-1] == ','):
			current.WriteByte(ch) // 负数的处理
		case ch == ',' || ch == '(' || ch == ')':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(ch))
		case !unicode.IsSpace(rune(ch)):
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			if ch == 'm' && i+2 < len(expression) {
				if expression[i:i+3] == "max" {
					tokens = append(tokens, "max")
					i += 2
				} else if expression[i:i+3] == "min" {
					tokens = append(tokens, "min")
					i += 2
				} else {
					tokens = append(tokens, string(ch))
				}
			} else {
				tokens = append(tokens, string(ch))
			}
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}

func formatParams(exp string) ([]string, error) {
	if len(exp) == 0 {
		return nil, errors.New("no expression")
	}
	exp = strings.ToLower(exp)
	isValid := true
	var builder strings.Builder
	for _, r := range exp {
		if r != ' ' {
			builder.WriteRune(r)
		}
	}
	byteExp := []rune(builder.String())
	for i := 0; i < len(byteExp); i++ {
		char := byteExp[i]
		if !strings.ContainsRune(validStr, char) {
			isValid = false
			break
		}
	}
	if !isValid {
		return nil, errors.New("invalid expression")
	}

	tokens, err := infixToPostfix(exp)
	return tokens, err
}

func Computing(exp string) (float64, error) {
	tokens, err := formatParams(exp)
	if err != nil {
		return 0, err
	}
	result, err := calculate(tokens)
	if err != nil {
		return 0, err
	}
	return result, nil
}
func infixToPostfix(expression string) ([]string, error) {
	tokens := tokenize(expression)
	var output []string
	var stack []string
	expectOperand := true
	for _, token := range tokens {
		switch {
		case isNumber(token):
			output = append(output, token)
			expectOperand = false
		case token == "(":
			stack = append(stack, token)
			expectOperand = true
		case token == ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf("括号不匹配")
			}
			stack = stack[:len(stack)-1] // 弹出左括号
			if len(stack) > 0 && isFunction(stack[len(stack)-1]) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			expectOperand = false
		case token == ",":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			expectOperand = true
		case isFunction(token):
			stack = append(stack, token)
			expectOperand = true
		case isOperator(token):
			if expectOperand && token == "-" {
				output = append(output, "0") // 添加一个 0 来处理负数
			}
			for len(stack) > 0 && isOperator(stack[len(stack)-1]) &&
				precedence[stack[len(stack)-1]] >= precedence[token] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
			expectOperand = true
		default:
			return nil, fmt.Errorf("无效的token: %s", token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, fmt.Errorf("括号不匹配")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func isLower(top string, newTop string) bool {
	switch top {
	case "+", "-":
		if newTop == "*" || newTop == "/" || newTop == "^" || newTop == "%" || newTop == "x" || newTop == "n" {
			return true
		}
	case "*", "/", "%":
		if newTop == "^" || newTop == "x" || newTop == "n" {
			return true
		}
	case "(":
		return true
	}
	return false
}

func applyOperator(operator string, a, b float64) (float64, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("除数不能为零")
		}
		return a / b, nil
	case "^":
		return math.Pow(a, b), nil
	case "%":
		return math.Mod(a, b), nil
	default:
		return 0, fmt.Errorf("未知的运算符: %s", operator)
	}
}

func applyFunction(function string, a, b float64) (float64, error) {
	switch function {
	case "max":
		return math.Max(a, b), nil
	case "min":
		return math.Min(a, b), nil
	default:
		return 0, fmt.Errorf("未知的函数: %s", function)
	}
}

func calculate(tokens []string) (float64, error) {
	stack := []float64{}

	for _, token := range tokens {
		switch {
		case isNumber(token):
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		case isOperator(token):
			if len(stack) < 2 {
				return 0, fmt.Errorf("表达式错误：运算符 %s 缺少操作数", token)
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			result, err := applyOperator(token, a, b)
			if err != nil {
				return 0, err
			}
			stack = append(stack, result)
		case isFunction(token):
			if len(stack) < 2 {
				return 0, fmt.Errorf("表达式错误：函数 %s 缺少参数", token)
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			result, err := applyFunction(token, a, b)
			if err != nil {
				return 0, err
			}
			stack = append(stack, result)
		default:
			return 0, fmt.Errorf("未知的token: %s", token)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("表达式错误：操作数过多")
	}

	return stack[0], nil

}
