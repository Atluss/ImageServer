// Package headers various structs for API
package headers

import (
	"encoding/json"
	"fmt"
	"github.com/Atluss/ImageServer/pkg/v1"
	"net/http"
)

// LoadedImage struct for image Source is full image and Preview
type LoadedImage struct {
	Source  string
	Preview string
}

// ReplayStatus for complete requests
type ReplayStatus struct {
	Status      int           `json:"Status"`
	Description string        `json:"Description"`
	Images      []LoadedImage `json:"Images"`
}

// Encode to json
func (obj *ReplayStatus) Encode(w http.ResponseWriter) error {
	err := json.NewEncoder(w).Encode(&obj)
	if !v1.LogOnError(err, "error: can't encode ReplayStatus") {
		return err
	}
	return nil
}

// RequestCreateImgJsonBase64 struct for base64 images
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

// Validate request
func (t *RequestCreateImgJsonBase64) Validate() error {
	if t.Data == "" || t.Body == "" {
		return fmt.Errorf("error: bad request")
	}
	return nil
}

// SetDefaultHeadersJson for json pages
func SetDefaultHeadersJson(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Content-Type", "application/json")
}
