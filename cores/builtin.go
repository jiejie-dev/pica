package pica

import (
	"github.com/icrowley/fake"
	"github.com/jeremaihloo/funny/langs"
)

// Email builtin function to fake an email
func Email(interpreter *langs.Interpreter, args []langs.Value) langs.Value {
	return langs.Value(fake.EmailAddress())
}

// Address builtin function to fake an address
func Address(interpreter *langs.Interpreter, args []langs.Value) langs.Value {
	return langs.Value(fake.StreetAddress())
}

// FullName builtin function to fake a fullname
func FullName(interpreter *langs.Interpreter, args []langs.Value) langs.Value {
	return langs.Value(fake.UserName())
}

// Phone builtin function to fake a phone no
func Phone(interpreter *langs.Interpreter, args []langs.Value) langs.Value {
	return langs.Value(fake.Phone())
}

// Words builtin function to fake words
func Words(interpreter *langs.Interpreter, args []langs.Value) langs.Value {
	return langs.Value(fake.Words())
}

// Domain builtin function to fake a domain
func Domain(interpreter *langs.Interpreter, args []langs.Value) langs.Value {
	return langs.Value(fake.DomainZone())
}
