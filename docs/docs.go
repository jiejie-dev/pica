package docs

import "github.com/jeremaihloo/pica"

type DocGenerator struct {
	ctx *pica.ApiContext
}

func (g *DocGenerator) Get() (string, error) {
	return "", nil
}
