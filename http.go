package pica

import (
	"net/http"
	"bytes"
	"io"
)

type HttpClient struct {
	baseUrl string
	client  *http.Client
}

func NewHttpClient(baseUrl string) *HttpClient {
	return &HttpClient{
		baseUrl: baseUrl,
		client: &http.Client{

		},
	}
}

func (c *HttpClient) Do(req ApiRequest) (*http.Response, error) {
	var body io.Reader
	switch req.Method {
	case "POST", "PATCH", "PUT":
		body = bytes.NewReader(req.Body)
	}
	r, err := http.NewRequest(req.Method, c.baseUrl+req.Url, body)
	if err != nil {
		return nil, err
	}
	r.Header = req.Headers
	res, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
