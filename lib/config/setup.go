package config

import (
	"github.com/Atluss/ImageServer/lib"
	"github.com/gorilla/mux"
)

func NewApiSetup(settings string) *Setup {
	cnf, err := Config(settings)
	lib.FailOnError(err, "error config file")
	set, err := newSetup(cnf)
	lib.FailOnError(err, "error setup")
	return set
}

func newSetup(cnf *config) (*Setup, error) {
	set := Setup{}
	if err := cnf.validate(); err != nil {
		return &set, err
	}
	set.Config = cnf
	set.Route = mux.NewRouter().StrictSlash(true)
	return &set, nil
}

// setup main setup api struct
type Setup struct {
	Config *config     // api setting
	Route  *mux.Router // mux frontend
}
