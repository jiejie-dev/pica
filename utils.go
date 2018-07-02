package pica

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/jeremaihloo/funny/langs"
)

func PrintHeaders(headers http.Header) {
	fmt.Printf("Headers:\n")
	for key, _ := range headers {
		fmt.Printf("%s: %s\n", key, headers.Get(key))
	}
}

func PrintJson(obj interface{}) {
	switch obj := obj.(type) {
	case map[string]interface{}:
		PrintJson(&obj)
		break
	case *map[string]interface{}:
		data, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println("\nJson:")
		color.Cyan("%s", data)
		break
	case []byte:
		var newObj map[string]interface{}
		err := json.Unmarshal(obj, &newObj)
		if err != nil {
			panic(err)
		}
		PrintJson(newObj)
		break
	default:
		data, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println("\nJson:")
		color.Cyan("%s", data)
		break

	}
}

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
