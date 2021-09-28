package pica

import (
	"fmt"
	"testing"
)

func TestGenPostmanGet(t *testing.T) {
	generator := NewScriptsGenerator("postman")
	result := generator.Generate("gen_test.json")
	fmt.Println(result)
}
