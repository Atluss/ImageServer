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

func TestSetup(t *testing.T) {

	path := "settings.json"
	set := NewApiSetup(path)

	log.Printf("%+v", set)
	log.Printf("Name: %s", set.Config.Name)
	log.Printf("Version: %s", set.Config.Version)
	log.Printf("Host: %s", set.Config.Host)
	log.Printf("Port: %s", set.Config.Port)
}
