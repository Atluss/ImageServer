package config

import (
	"github.com/Atluss/ImageServer/lib"
	"log"
	"testing"
)

func TestConfig(t *testing.T) {

	path := "settings.json"

	cnf, err := Config(path)
	lib.FailOnError(err, "Test error")

	log.Printf("%+v", cnf)
	log.Printf("Name: %s", cnf.Name)
	log.Printf("Version: %s", cnf.Version)
	log.Printf("Host: %s", cnf.Host)
	log.Printf("Port: %s", cnf.Port)
}
