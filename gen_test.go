package pica

import (
	"fmt"
	"testing"
)

func TestGenPostmanGet(t *testing.T) {
	result := GenerateScriptsByPostman("gen_test.json")
	fmt.Println(result)
}
