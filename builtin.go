package pica

import (
	"github.com/icrowley/fake"
	"github.com/jerloo/funny"
)

// Email builtin function to fake an email
func Email(interpreter *funny.Interpreter, args []funny.Value) funny.Value {
	return funny.Value(fake.EmailAddress())
}

// Address builtin function to fake an address
func Address(interpreter *funny.Interpreter, args []funny.Value) funny.Value {
	return funny.Value(fake.StreetAddress())
}

// FullName builtin function to fake a fullname
func FullName(interpreter *funny.Interpreter, args []funny.Value) funny.Value {
	return funny.Value(fake.UserName())
}

// Phone builtin function to fake a phone no
func Phone(interpreter *funny.Interpreter, args []funny.Value) funny.Value {
	return funny.Value(fake.Phone())
}

// Words builtin function to fake words
func Words(interpreter *funny.Interpreter, args []funny.Value) funny.Value {
	return funny.Value(fake.Words())
}

// Domain builtin function to fake a domain
func Domain(interpreter *funny.Interpreter, args []funny.Value) funny.Value {
	return funny.Value(fake.DomainZone())
}
