package main

import (
	"fmt"
	"github.com/Atluss/ImageServer/lib"
	"github.com/Atluss/ImageServer/lib/config"
	"github.com/Atluss/ImageServer/lib/headers"
	"github.com/Atluss/ImageServer/lib/images"
	"log"
	"net/http"
)

func main() {

	settingPath := "settings.json"
	set := config.NewApiSetup(settingPath)

	lib.FailOnError(images.CreateTestFilesDir(images.ImageFolder), "fatal error")

	set.Route.HandleFunc("/form_data", func(w http.ResponseWriter, r *http.Request) {

		headers.SetDefaultHeadersJson(w)
		reply := headers.ReplayStatus{
			Status: http.StatusOK,
			Images: images.GetImagesFormDataAndQuery(r),
		}

		if len(reply.Images) == 0 {
			reply.Description = "Images not found :("
		}

		lib.LogOnError(reply.Encode(w), "error")
	})

	set.Route.HandleFunc("/json_image", func(w http.ResponseWriter, r *http.Request) {

		headers.SetDefaultHeadersJson(w)
		reply := headers.ReplayStatus{
			Status: http.StatusOK,
			Images: []headers.LoadedImage{},
		}

		if img, err := images.GetJsonImageBase64(r); err != nil {
			log.Printf("error json base 64 img: %s", err)
			reply.Status = http.StatusBadRequest
			reply.Description = http.StatusText(http.StatusBadRequest)
		} else {
			reply.Images = append(reply.Images, img)
		}

		lib.LogOnError(reply.Encode(w), "error")
	})

	// images public folder
	sImgDir := fmt.Sprintf("/%s/", images.ImageFolder)
	pImgDir := fmt.Sprintf("./%s", images.ImageFolder)
	set.Route.PathPrefix(sImgDir).Handler(http.StripPrefix(sImgDir, http.FileServer(http.Dir(pImgDir))))

	lib.FailOnError(http.ListenAndServe(fmt.Sprintf(":%s", set.Config.Port), set.Route), "error")
}
