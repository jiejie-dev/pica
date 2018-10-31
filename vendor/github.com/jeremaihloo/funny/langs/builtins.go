package langs

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

// BuiltinFunction function handler
type BuiltinFunction func(interpreter *Interpreter, args []Value) Value

var (
	// FUNCTIONS all builtin functions
	FUNCTIONS = map[string]BuiltinFunction{
		"echo":         Echo,
		"echoln":       Echoln,
		"now":          Now,
		"base64encode": Base64Encode,
		"base64decode": Base64Decode,
		"assert":       Assert,
		"len":          Len,
		"hash":         Md5,
		"max":          Max,
		"typeof":       Typeof,
		"uuid":         UUID,
	}
)

// ackEq check function arguments count valid
func ackEq(args []Value, count int) {
	if len(args) != count {
		panic(fmt.Sprintf("%d arguments required but got %d", count, len(args)))
	}
}

// ackGt check function arguments count valid
func ackGt(args []Value, count int) {
	if len(args) <= count {
		panic(fmt.Sprintf("greater than %d arguments required but got %d", count, len(args)))
	}
}

// Echo builtin function echos one or every item in a array
func Echo(interpreter *Interpreter, args []Value) Value {
	fmt.Sprint(interpreter.Vars)
	for _, item := range args {
		fmt.Print(item)
	}
	return nil
}

// Echoln builtin function echos one or every item in a array
func Echoln(interpreter *Interpreter, args []Value) Value {
	fmt.Sprint(interpreter.Vars)
	for index, item := range args {
		fmt.Print(item)
		if index == len(args)-1 {
			fmt.Print("\n")
		}
	}
	return nil
}

// Now builtin function return now time
func Now(interpreter *Interpreter, args []Value) Value {
	return Value(time.Now())
}

// Base64Encode return base64 encoded string
func Base64Encode(interpreter *Interpreter, args []Value) Value {
	base64encode := func(val string) string {
		return base64.StdEncoding.EncodeToString([]byte(val))
	}
	if len(args) == 1 {
		return Value(base64encode(args[0].(string)))
	}
	var results []string
	for _, item := range args {
		results = append(results, base64encode(item.(string)))
	}
	return Value(results)
}

// Base64Decode return base64 decoded string
func Base64Decode(interpreter *Interpreter, args []Value) Value {
	base64decode := func(val string) string {
		sb, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			panic(err)
		}
		return string(sb)
	}
	if len(args) == 1 {
		return Value(base64decode(args[0].(string)))
	}
	var results []string
	for _, item := range args {
		results = append(results, base64decode(item.(string)))
	}
	return Value(results)
}

// Assert return the value that has been given
func Assert(interpreter *Interpreter, args []Value) Value {
	ackEq(args, 1)
	if val, ok := args[0].(bool); ok {
		if val {
			return Value(args[0])
		}
		panic("assert false")
	}
	panic("type error, only support [bool]")
}

// Len return then length of the given list
func Len(interpreter *Interpreter, args []Value) Value {
	ackEq(args, 1)
	if val, ok := args[0].(*List); ok {
		return Value(len(val.Values))
	}
	panic("type error, only support [list]")
}

// Md5 return then length of the given list
func Md5(interpreter *Interpreter, args []Value) Value {
	ackEq(args, 1)
	switch v := args[0].(type) {
	case string:
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(v))
		return hex.EncodeToString(md5Ctx.Sum(nil))
	default:
		break
	}
	panic("type error, only support [string]")
}

// Max return then length of the given list
func Max(interpreter *Interpreter, args []Value) Value {
	ackGt(args, 1)
	switch v := args[0].(type) {
	case int:
		flag := v
		for _, item := range args[1:] {
			if val, ok := item.(int); ok {
				if val > flag {
					flag = val
				}
			}
			break
		}
		return Value(flag)
	case *List:
		flag := interpreter.EvalExpression(v.Values[0])
		if flagA, ok := flag.(int); ok {
			for _, item := range v.Values {
				val := interpreter.EvalExpression(item)
				if val, ok := val.(int); ok {
					if val > flagA {
						flagA = val
					}
				}
				break
			}
			return Value(flagA)
		}
	default:
		break
	}
	panic("type error, only support [int]")
}

// Typeof builtin function echos one or every item in a array
func Typeof(interpreter *Interpreter, args []Value) Value {
	ackEq(args, 1)
	return Typing(args[0])
}

// UUID builtin function return a uuid string value
func UUID(interpreter *Interpreter, args []Value) Value {
	ackEq(args, 0)
	u1 := uuid.NewV4()
	return Value(fmt.Sprintf("%s", u1))
}
