package config

import (
	"github.com/Atluss/ImageServer/pkg/v1"
	"github.com/gorilla/mux"
)

// NewApiSetup return new api setup with router
func NewApiSetup(settings string) *Setup {
	cnf, err := Config(settings)
	v1.FailOnError(err, "error config file")
	set, err := newSetup(cnf)
	v1.FailOnError(err, "error setup")
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

// Setup main setup api struct
type Setup struct {
	Config *config     // api setting
	Route  *mux.Router // mux frontend
}
