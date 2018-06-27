package utils

import (
	"net/http"
	"fmt"
	"encoding/json"
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
		fmt.Printf("\nJson:\n%s\n\n", data)
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
		panic("unknow type PrintJson")
	}
}
