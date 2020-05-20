package googleapi

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleStruct struct{}

type GoogleInterface interface {
	ConfigFromJSON(jsonKey []byte, scope ...string) (*oauth2.Config, error)
}

var (
	GoogleStruct GoogleInterface = &googleStruct{}
)

func (gs *googleStruct) ConfigFromJSON(jsonKey []byte, scope ...string) (*oauth2.Config, error) {
	return google.ConfigFromJSON(jsonKey, scope...)
}
