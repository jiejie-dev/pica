package langs

import (
	"encoding/base64"
	"fmt"
	"time"
)

type BuiltinFunction = func(interpreter *Interpreter, args []Value) Value

var (
	FUNCTIONS = map[string]BuiltinFunction{
		"echo":         Echo,
		"echoln":       Echoln,
		"now":          Now,
		"base64encode": Base64Encode,
		"base64decode": Base64Decode,
		"assert":       Assert,
		"len":          Len,
	}
)

// ack check function arguments count valid
func ack(args []Value, count int) {
	if len(args) != count {
		panic(fmt.Sprintf("%d arguments required but got %d", count, len(args)))
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
	ack(args, 1)
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
	ack(args, 1)
	if val, ok := args[0].(*List); ok {
		return Value(len(val.Values))
	}
	panic("type error, only support [list]")
}
