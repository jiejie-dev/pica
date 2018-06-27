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

func VmMap2HttpHeaders(vmMap map[string]langs.Value) http.Header {
	headers := http.Header{}
	for k, v := range vmMap {
		headers.Set(k, v.(string))
	}
	return headers
}
