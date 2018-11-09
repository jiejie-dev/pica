package pica

import (
	"github.com/icrowley/fake"
	"github.com/jeremaihloo/funny/langs"
)

func Email(interpreter *langs.Interpreter, args []langs.Value) langs.Value {
	return langs.Value(fake.EmailAddress())
}

func Address(interpreter *langs.Interpreter, args []langs.Value) langs.Value {
	return langs.Value(fake.StreetAddress())
}

func FullName(interpreter *langs.Interpreter, args []langs.Value) langs.Value {
	return langs.Value(fake.UserName())
}

func Phone(interpreter *langs.Interpreter, args []langs.Value) langs.Value {
	return langs.Value(fake.Phone())
}

func Words(interpreter *langs.Interpreter, args []langs.Value) langs.Value {
	return langs.Value(fake.Words())
}

func Domain(interpreter *langs.Interpreter, args []langs.Value) langs.Value {
	return langs.Value(fake.DomainZone())
}