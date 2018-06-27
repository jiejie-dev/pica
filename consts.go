package pica

import "github.com/jeremaihloo/funny/langs"

const (
	GENERATE_DESC = "created by pica https://github.com/jeremaihloo/pica"
)

var (
	DefaultHeaders = map[string]langs.Value{
		"Accept":          "* /*",
		"Accept-Language": "en-US,en;q=0.8",
		"Cache-Control":   "max-age=0",
		"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36",
		"Connection":      "keep-alive",
		"Referer":         "http://www.baidu.com/",
		"Content-Type":    "application/json",
	}

	DefaultInitScope = map[string]langs.Value{
		"headers": DefaultHeaders,
	}
)
