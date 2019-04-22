package headers

import (
	"encoding/json"
	"github.com/Atluss/ImageServer/lib"
	"net/http"
)

type LoadedImage struct {
	Source string
	Preview string
}

type ReplayStatus struct {
	Status      int           `json:"Status"`
	Description string        `json:"Description"`
	Images      []LoadedImage `json:"Images"`
}

func (obj *ReplayStatus) Encode(w http.ResponseWriter) error {
	err := json.NewEncoder(w).Encode(&obj)
	if !lib.LogOnError(err, "error: can't encode ReplayStatus") {
		return err
	}
	return nil
}

func SetDefaultHeadersJson(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Content-Type", "application/json")
}

