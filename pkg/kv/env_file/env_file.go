package env_file

import (
	"fmt"
	"encoding/base64"
	"github.com/starlingbank/vault-unsealer/pkg/kv"
)

type envFile struct {
	properties Properties
}

var _ kv.Service = &envFile{}

func New(filename string) (*envFile, error) {
	properties, err := ReadPropertiesFile(filename)
	if err != nil {
		return nil, err
	}

	return &envFile{
		properties: properties,
		}, nil
}

func (a *envFile) Get(key string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(a.properties[key])
	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}

func (a *envFile) Set(key string, val []byte) error {
	return fmt.Errorf("Not implemented")
}

func (a *envFile) Delete(key string) error {
	return fmt.Errorf("Not implemented")
}

func (g *envFile) Test(key string) error {
	return fmt.Errorf("Not implemented")
}
