package pica

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/jeremaihloo/funny/langs"

	"github.com/fatih/color"
)

type HttpClient struct {
	baseUrl string
	client  *http.Client
}

func NewHttpClient(baseUrl string) *HttpClient {
	return &HttpClient{
		baseUrl: baseUrl,
		client:  &http.Client{},
	}
}

func (c *HttpClient) Do(req ApiRequest, vm *langs.Interpreter) (*http.Response, error) {
	var body io.Reader
	switch req.Method {
	case "POST", "PATCH", "PUT":
		body = bytes.NewReader(req.Body)
	}
	targetUrl, query, err := CompileUrl(c.baseUrl+req.Url, vm)
	req.Query = query
	if err != nil {
		return nil, err
	}
	vm.Assign("targetUrl", targetUrl)
	r, err := http.NewRequest(req.Method, targetUrl, body)
	if err != nil {
		return nil, err
	}
	r.Header = req.Headers

	fmt.Printf("%s %s\n", req.Method, targetUrl)
	// print headers
	PrintHeaders(r.Header)

	res, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	fmt.Println(DefaultOutput.L("-"))
	fmt.Printf("\nResponse\n")
	if res.StatusCode == 200 {
		color.Green("Status: %d\n", res.StatusCode)
	} else {
		color.Red("Status: %d\n", res.StatusCode)
	}
	PrintHeaders(res.Header)
	return res, nil
}
