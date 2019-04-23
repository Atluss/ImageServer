// api run file
package main

import (
	"fmt"
	"github.com/Atluss/ImageServer/lib"
	"github.com/Atluss/ImageServer/lib/config"
	"github.com/Atluss/ImageServer/lib/headers"
	"github.com/Atluss/ImageServer/lib/images"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	settingPath := "settings.json"
	set := config.NewApiSetup(settingPath)

	// do something if user close program (close DB, or wait running query)
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Exit program...")
		os.Exit(1)
	}()

	lib.FailOnError(images.CreateTestFilesDir(images.ImageFolder), "fatal error")
	// handle multipart/form-data and query
	set.Route.HandleFunc("/form_data", func(w http.ResponseWriter, r *http.Request) {
		headers.SetDefaultHeadersJson(w)
		reply := headers.ReplayStatus{
			Status: http.StatusOK,
			Images: images.GetImagesFormDataAndQuery(r),
		}
		if len(reply.Images) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			reply.Status = http.StatusBadRequest
			reply.Description = "Images not found :("
		}
		lib.LogOnError(reply.Encode(w), "error")
	})
	// handle json base64 image
	set.Route.HandleFunc("/json_image", func(w http.ResponseWriter, r *http.Request) {
		headers.SetDefaultHeadersJson(w)
		reply := headers.ReplayStatus{
			Status: http.StatusOK,
			Images: []headers.LoadedImage{},
		}
		if img, err := images.GetJsonImageBase64(r); err != nil {
			log.Printf("error json base 64 img: %s", err)
			w.WriteHeader(http.StatusBadRequest)
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
