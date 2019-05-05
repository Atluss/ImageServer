package config

import (
	"encoding/json"
	"fmt"
	"github.com/Atluss/ImageServer/pkg/v1"
	"io/ioutil"
	"os"
)

// config load new config for API
func Config(path string) (*config, error) {
	conf := config{}
	if err := v1.CheckFileExist(path); err != nil {
		return &conf, err
	}
	conf.FilePath = path
	if err := conf.load(); err != nil {
		return &conf, err
	}
	return &conf, nil
}

// config main
type config struct {
	Name     string `json:"Name"`    // API name
	Version  string `json:"Version"` // API version
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	FilePath string `json:"FilePath"` // path to Json settings file
}

// load all settings
func (obj *config) load() error {
	jsonSet, err := os.Open(obj.FilePath)
	defer func() {
		v1.LogOnError(jsonSet.Close(), "warning: Can't close json settings file.")
	}()
	if !v1.LogOnError(err, "Can't open config file") {
		return err
	}
	bytesVal, _ := ioutil.ReadAll(jsonSet)
	err = json.Unmarshal(bytesVal, &obj)
	if !v1.LogOnError(err, "Can't unmarshal json file") {
		return err
	}
	return obj.validate()
}

func (obj *config) validate() error {
	if obj.Name == "" {
		return fmt.Errorf("config miss name")
	}
	if obj.Version == "" {
		return fmt.Errorf("config miss version")
	}
	if obj.Host == "" {
		return fmt.Errorf("config miss host")
	}
	if obj.Port == "" {
		return fmt.Errorf("config miss port")
	}
	return nil
}
