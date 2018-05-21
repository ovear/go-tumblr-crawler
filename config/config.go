package config

import (
	"encoding/json"
	io "io/ioutil"
	"fmt"
	"github.com/go-yaml/yaml"
)

type Config struct {
}

func NewConfig() *Config {
	return &Config{}
}

func (this *Config) Load(filename string, v interface{}) {
	data, err := io.ReadFile(filename)

	if err != nil {
		return
	}

	dataJson := []byte(data)
	err = json.Unmarshal(dataJson, v)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func (this *Config) LoadYaml(filename string, v interface{}) {
	data, err := io.ReadFile(filename)

	if err != nil {
		return
	}

	dataJson := []byte(data)
	err = yaml.Unmarshal(dataJson, v)

	if err != nil {
		fmt.Println(err)
		return
	}
}
