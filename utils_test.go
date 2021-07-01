package pica

import (
	"testing"

	"github.com/jerloo/funny"
	"github.com/magiconair/properties/assert"
)

func TestCompileUrl(t *testing.T) {
	vm := funny.NewInterpreterWithScope(funny.Scope{})
	vm.Assign("user_id", "10")
	url, query, err := CompileURL("/api/users/<user_id>", vm)
	if err != nil {
		t.Error(err)
	}
	t.Log(url)
	assert.Equal(t, url, "/api/users/10")
	vm.Assign("query", map[string]interface{}{
		"name": "jerloo",
		"age":  "10",
	})
	url, query, err = CompileURL("/api/users/<user_id>", vm)
	if err != nil {
		t.Error(err)
	}
	t.Log(url)
	assert.Equal(t, url, "/api/users/10?name=jerloo&age=10")
	assert.Equal(t, query["name"], "jerloo")
}
