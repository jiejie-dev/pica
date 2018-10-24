package langs

import (
	"fmt"
	"reflect"
)

func Typing(data interface{}) string {
	t := reflect.TypeOf(data)
	if t == nil {
		return "nil"
	}
	return fmt.Sprintf("%s", t.String())
}
