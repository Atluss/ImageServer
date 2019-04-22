package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/Atluss/ImageServer/lib"
	"github.com/Atluss/ImageServer/lib/config"
	"github.com/Atluss/ImageServer/lib/headers"
	"github.com/Atluss/ImageServer/lib/images"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	settingPath := "settings.json"
	set := config.NewApiSetup(settingPath)

	set.Route.HandleFunc("/form_data", func(w http.ResponseWriter, r *http.Request) {

		headers.SetDefaultHeadersJson(w)
		reply := headers.ReplayStatus{
			Status: http.StatusOK,
			Images: images.GetImages(r),
		}

		if len(reply.Images) == 0 {
			reply.Description = "Images not found :("
		}

		lib.LogOnError(reply.Encode(w), "error")
	})

	set.Route.HandleFunc("/json_image", func(w http.ResponseWriter, r *http.Request) {

	})

	// images public folder
	set.Route.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	lib.FailOnError(http.ListenAndServe(fmt.Sprintf(":%s", set.Config.Port), set.Route), "error")
}

func GetJsonImageBase64(image string) {

	img, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		log.Println(err)
	} else {
		out, err := os.Create("test.jpg")
		if err != nil {
			log.Println(err)
		}
		defer out.Close()

		// Write the body to file
		r := bytes.NewReader(img)
		_, err = io.Copy(out, r)
		if err != nil {
			log.Println(err)
		}
	}
}