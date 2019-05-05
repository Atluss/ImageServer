package images

import (
	"fmt"
	"log"
	"testing"
)

// TestGenerateName
func TestGenerateName(t *testing.T) {
	st := GenerateName(AllowFormats[0])
	log.Printf("Source: %s", st.Source)
	log.Printf("Preview: %s", st.Preview)
}

func TestCreateImageName(t *testing.T) {
	if st, err := CreateImageName(fmt.Sprintf("asdasd.%s", AllowFormats[0])); err != nil {
		t.Errorf("Error: %s", err)
	} else {
		log.Printf("Source: %s", st.Source)
		log.Printf("Preview: %s", st.Preview)
	}
}
