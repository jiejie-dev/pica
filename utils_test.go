package pica

import (
	"testing"

	"github.com/magiconair/properties/assert"

	"github.com/jeremaihloo/funny/langs"
)

func TestCompileUrl(t *testing.T) {
	vm := langs.NewInterpreterWithScope(langs.Scope{})
	vm.Assign("user_id", "10")
	url, query, err := CompileUrl("/api/users/<user_id>", vm)
	if err != nil {
		t.Error(err)
	}
	t.Log(url)
	assert.Equal(t, url, "/api/users/10")
	vm.Assign("query", map[string]interface{}{
		"name": "jeremaihloo",
		"age":  "10",
	})
	url, query, err = CompileUrl("/api/users/<user_id>", vm)
	if err != nil {
		t.Error(err)
	}
	t.Log(url)
	assert.Equal(t, url, "/api/users/10?name=jeremaihloo&age=10")
	assert.Equal(t, query["name"], "jeremaihloo")
}
