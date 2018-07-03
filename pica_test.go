package pica

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestQuery(t *testing.T) {
	query := Query{
		"name": "jeremaihloo",
		"age":  "20",
	}
	qs, err := query.String()
	if err != nil {
		t.Error(err)
	}
	t.Log(qs)
	queryString := "name=jeremaihloo&age=20"
	assert.Equal(t, queryString, qs)
	newQuery, err := ParseQuery(queryString)
	if err != nil {
		t.Error(err)
	}
	for k, v := range query {
		if newQuery[k] != v {
			t.Error("not eq	")
		}
	}
}
