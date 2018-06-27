package utils

import (
	"github.com/jeremaihloo/funny/langs"
	"net/http"
)

func HttpHeaders2VmMap(httpHeader http.Header) map[string]langs.Value {
	var headers = map[string]langs.Value{}
	for k, _ := range httpHeader {
		headers[k] = httpHeader.Get(k)
	}
	return headers
}
