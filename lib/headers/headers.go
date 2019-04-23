// various structs for API
package headers

import (
	"encoding/json"
	"fmt"
	"github.com/Atluss/ImageServer/lib"
	"net/http"
)

type LoadedImage struct {
	Source  string
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

type RequestCreateImgJsonBase64 struct {
	Data string `json:"Data"` // Example: data:image/jpeg;
	Body string `json:"Body"` // body Base64
}

// Decode request
func (t *RequestCreateImgJsonBase64) Decode(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return err
	}
	return nil
}

func (t *RequestCreateImgJsonBase64) Validate() error {
	if t.Data == "" || t.Body == "" {
		return fmt.Errorf("error: bad request")
	}
	return nil
}

func SetDefaultHeadersJson(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Content-Type", "application/json")
}
