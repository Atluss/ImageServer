package config

import (
	"log"
	"testing"
)

func TestConfig(t *testing.T) {

	path := "settings.json"

	cnf, err := Config(path)
	if err != nil {
		t.Errorf("Test error: %s", err)
	}

	log.Printf("%+v", cnf)
	log.Printf("Name: %s", cnf.Name)
	log.Printf("Version: %s", cnf.Version)
	log.Printf("Host: %s", cnf.Host)
	log.Printf("Port: %s", cnf.Port)
}
