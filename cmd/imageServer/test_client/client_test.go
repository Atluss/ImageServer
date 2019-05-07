package test_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Atluss/ImageServer/pkg/v1/headers"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"testing"
)

// TestSendQueryImage send link to image to get parameter
func TestSendQueryImage(t *testing.T) {
	link := "http://localhost:8080/form_data?image=https://img.zoneofgames.ru/news/2019/04/22/190216-banner_conk_20190422_PaganOnline.jpg"
	resp, err := http.Get(link)
	if err != nil {
		t.Errorf("error api request: %s", err)
		return
	}
	defer resp.Body.Close()
	replySt := headers.ReplayStatus{}
	if err := json.NewDecoder(resp.Body).Decode(&replySt); err != nil {
		t.Errorf("error json decode: %s", err)
	} else {
		fmt.Printf("Answer: %+v", replySt)
	}
}

// TestSendFormData send image in multipart/form-data
func TestSendFormData(t *testing.T) {
	filePath := "test_image.jpg"
	file, _ := os.Open(filePath)
	defer file.Close()

	fieldName := strings.NewReplacer("\\", "\\\\", `"`, "\\\"").Replace(file.Name())
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", fieldName))
	h.Set("Content-Type", "image/jpg")
	part, _ := writer.CreatePart(h)

	if _, err := io.Copy(part, file); err != nil {
		t.Errorf("error can't copy image data: %s", err)
	}
	writer.Close()

	r, _ := http.NewRequest("POST", "http://localhost:8080/form_data", body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	if rest, err := client.Do(r); err != nil {
		t.Errorf("error api request: %s", err)
	} else {
		replySt := headers.ReplayStatus{}
		if err := json.NewDecoder(rest.Body).Decode(&replySt); err != nil {
			t.Errorf("error json decode: %s", err)
		} else {
			fmt.Printf("Answer: %+v", replySt)
		}
	}
}

// TestSendJson send json base64 data
func TestSendJson(t *testing.T) {
	url := "http://localhost:8080/json_image"
	requestBody := headers.RequestCreateImgJsonBase64{
		Data: "data:image/jpeg",
		Body: "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",
	}
	jsonStr, err := json.Marshal(&requestBody)
	if err != nil {
		t.Errorf("error Marhal json^ %s", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("error request: %s", err)
	}
	defer resp.Body.Close()

	replySt := headers.ReplayStatus{}
	if err := json.NewDecoder(resp.Body).Decode(&replySt); err != nil {
		t.Errorf("error json decode: %s", err)
	} else {
		fmt.Printf("Answer: %+v", replySt)
	}
}
