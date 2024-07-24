package math

import (
	"errors"
	"github.com/stanleychan/go-utils/data"
	"math"
	"strconv"
	"strings"
	"unicode"
)

const validStr = "x0123456789.+-*/()^%"

func parseStringToF64(raw string) (float64, error) {
	// check the
	length := len(raw)
	if length < 1 {
		return 0, errors.New("no string to convert")
	}
	if raw[length-1] == '%' {
		// it may be a percentage
		f64, err := strconv.ParseFloat(raw[:length-1], 64)
		if err != nil {
			return 0, err
		}
		return f64 / 100, nil
	}

	f64raw, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, err
	}
	return f64raw, nil
}

func convertTypeToF64(raw interface{}) (float64, error) {
	switch raw.(type) {
	case string:
		return parseStringToF64(raw.(string))
	case bool:
		if raw.(bool) {
			return 1, nil
		} else {
			return 0, nil
		}
	case int:
		return (float64)(raw.(int)), nil
	case int8:
		return (float64)(raw.(int8)), nil
	case int16:
		return (float64)(raw.(int16)), nil
	case int32:
		return (float64)(raw.(int32)), nil
	case float32:
		return (float64)(raw.(float32)), nil
	case float64:
		return raw.(float64), nil
	default:
		return 0, errors.New("not support type")
	}
}

func formatParams(exp string) (string, error) {
	if len(exp) == 0 {
		return "", errors.New("no expression")
	}
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
		return "", errors.New("invalid expression")
	}

	exp = infixToPostfix(exp)
	return exp, nil
}

func Computing(exp string) (float64, error) {
	exp, err := formatParams(exp)
	if err != nil {
		return 0, err
	}
	result := calculate(exp)
	if result == nil {
		return 0, errors.New("calculate error")
	}
	return result.(float64), nil
}
func infixToPostfix(exp string) string {
	stack := data.Stack{}
	postfix := ""
	expLen := len(exp)
	for i := 0; i < expLen; i++ {
		char := string(exp[i])
		switch char {
		case " ":
			continue
		case "(":
			stack.Push("(")
		case ")":
			for !stack.IsEmpty() {
				preChar := stack.Top()
				if preChar == "(" {
					stack.Pop() // pop "("
					break
				}
				postfix += preChar.(string) + " "
				stack.Pop()
			}

		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ".", "x":
			j := i
			digit := ""
			for ; j < expLen && (unicode.IsDigit(rune(exp[j])) || exp[j] == '.' || exp[j] == 'x'); j++ {
				digit += string(exp[j])
			}
			postfix += digit + " "
			i = j - 1

		default:
			for !stack.IsEmpty() {
				top := stack.Top()
				if top == "(" || isLower(top.(string), char) {
					break
				}
				postfix += top.(string) + " "
				stack.Pop()
			}
			stack.Push(char)
		}
	}

	for !stack.IsEmpty() {
		postfix += stack.Pop().(string) + " "
	}

	return postfix
}

func isLower(top string, newTop string) bool {
	switch top {
	case "+", "-":
		if newTop == "*" || newTop == "/" || newTop == "^" {
			return true
		}
	case "*", "/", "%":
		if newTop == "^" {
			return true
		}
	case "(":
		return true
	}
	return false
}

func calculate(postfix string) interface{} {
	stack := data.Stack{}
	fixLen := len(postfix)
	strNum := ""
	for i := 0; i < fixLen; i++ {
		nextChar := string(postfix[i])
		if unicode.IsDigit(rune(postfix[i])) || postfix[i] == '.' {
			strNum += nextChar
		} else if postfix[i] == ' ' {
			if len(strNum) > 0 {
				stack.Push(strNum)
			}
			strNum = ""
		} else {
			if stack.Top() == nil {
				return nil
			}
			num1, _ := strconv.ParseFloat(stack.Pop().(string), 64)
			noFirstNum := false
			num2 := float64(0)
			if stack.Top() == nil {
				noFirstNum = true
			} else {
				num2, _ = strconv.ParseFloat(stack.Pop().(string), 64)
			}
			switch nextChar {
			case "+":
				if noFirstNum {
					num2 = 0
				}
				stack.Push(strconv.FormatFloat(num1+num2, 'f', 20, 64))
			case "-":
				if noFirstNum {
					num2 = 0
				}
				stack.Push(strconv.FormatFloat(num2-num1, 'f', 20, 64))
			case "*":
				if noFirstNum {
					return nil
				}
				stack.Push(strconv.FormatFloat(num1*num2, 'f', 20, 64))
			case "/":
				if noFirstNum {
					return nil
				}
				stack.Push(strconv.FormatFloat(num2/num1, 'f', 20, 64))
			case "^":
				if noFirstNum {
					return nil
				}
				stack.Push(strconv.FormatFloat(math.Pow(num2, num1), 'f', 20, 64))
			case "%":
				if noFirstNum {
					return nil
				}
				i64Num1 := int64(num1)
				if i64Num1 == 0 {
					return nil
				}
				i64Num2 := int64(num2)
				stack.Push(strconv.FormatFloat(float64(i64Num2%i64Num1), 'f', 20, 64))
			}

		}
	}
	result, _ := strconv.ParseFloat(stack.Top().(string), 64)
	return result
}
