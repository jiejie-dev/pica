package langs

import (
	"fmt"
	"reflect"
)

// Typing return the type name of one object
func Typing(data interface{}) string {
	t := reflect.TypeOf(data)
	if t == nil {
		return "nil"
	}
	return fmt.Sprintf("%s", t.String())
}
