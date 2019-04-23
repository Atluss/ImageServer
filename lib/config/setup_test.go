package config

import (
	"log"
	"testing"
)

func TestSetup(t *testing.T) {
	path := "settings.json"
	set := NewApiSetup(path)

	log.Printf("%+v", set)
	log.Printf("Name: %s", set.Config.Name)
	log.Printf("Version: %s", set.Config.Version)
	log.Printf("Host: %s", set.Config.Host)
	log.Printf("Port: %s", set.Config.Port)
}
