package pica

import (
	"github.com/jeremaihloo/funny/langs"
	"fmt"
)

const (
	GENERATE_DESC = "created by pica https://github.com/jeremaihloo/pica"
)

var (
	DefaultHeaders = map[string]langs.Value{
		"Accept":          "* /*",
		"Accept-Language": "en-US,en;q=0.8",
		"Cache-Control":   "max-age=0",
		"User-Agent":      fmt.Sprintf("Pica Api Test Client/%s https://github.com/jeremaihloo/pica", Version),
		"Connection":      "keep-alive",
		"Referer":         "http://www.baidu.com/",
		"Content-Type":    "application/json",
	}

	DefaultInitScope = map[string]langs.Value{
		"headers": DefaultHeaders,
	}
)
