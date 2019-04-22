package images

import (
	"fmt"
	"github.com/Atluss/ImageServer/lib"
	"github.com/Atluss/ImageServer/lib/headers"
	"github.com/disintegration/imaging"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	ImageFolder = "images"
	PreviewWidth = 100
	PreviewHeight = 100
)

var AllowFormats = [3]string{"jpg", "jpeg", "png"}

func GetImages(r *http.Request) []headers.LoadedImage{

	images := getImagesFormData(r)

	if tmp, err := getImageFromLink(r); err != nil {
		log.Println(err)
	} else {
		images = append(images, tmp)
	}

	return images
}

func getImagesFormData(r *http.Request) (images []headers.LoadedImage){

	reader, err := r.MultipartReader()
	if err != nil {
		log.Println(err)
	}

	for {
		part, err_part := reader.NextPart()
		if err_part == io.EOF {
			break
		}

		if part.FileName() != "" && strings.Contains(part.Header.Get("Content-Type"), "image") {

			loadedImg, err := createImageName(part.FileName())
			if err != nil {
				log.Printf("error data %s", err)
				continue
			}

			// Create the file
			out, err := os.Create(loadedImg.Source)
			if err != nil {
				log.Println(err)
				continue
			}
			defer out.Close()

			// Write the body to file
			_, err = io.Copy(out, part)
			if err != nil {
				log.Println(err)
				continue
			}

			if err := createPreview(loadedImg, PreviewWidth, PreviewHeight); err != nil {
				log.Println(err)
				continue
			}

			images = append(images, loadedImg)

		}
	}

	return images
}

// getImageFromLink search query and download image
func getImageFromLink(r *http.Request) (loadedImg headers.LoadedImage, err error){

	link := r.URL.Query().Get("image")
	if link == "" {
		return loadedImg, fmt.Errorf("no link send")
	}

	log.Println(link)

	loadedImg, err = createImageName(link)
	if err != nil {
		return loadedImg, fmt.Errorf("link not correct, %s", err)
	}

	// Create the file
	out, err := os.Create(loadedImg.Source)
	if err != nil {
		return loadedImg, err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(link)
	if err != nil {
		return loadedImg, err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return loadedImg, err
	}

	if err := createPreview(loadedImg, PreviewWidth, PreviewHeight); err != nil {
		return loadedImg, err
	}

	return loadedImg, err
}

// createImageName generate new name and preview name from source name
func createImageName(name string) (newName headers.LoadedImage, err error) {

	arr := strings.Split(name, ".")

	if len(arr) < 2 {
		return newName, err
	}

	format := arr[len(arr)-1]

	// yes it rude but right now perfect no need cycle :)
	if format != AllowFormats[0] && format != AllowFormats[1] && format != AllowFormats[2]{
		return newName, fmt.Errorf("format image file not allow")
	}

	for {
		newName = generateName(format)
		if err := lib.CheckFileExist(newName.Source); err != nil {
			break
		}
	}

	return newName, err
}

// generateName
func generateName(format string) (newName headers.LoadedImage) {
	shortName := uuid.NewV4()
	newName.Source = fmt.Sprintf("%s/%s.%s", ImageFolder, shortName, format)
	newName.Preview = fmt.Sprintf("%s/%s_preview.%s", ImageFolder, shortName, format)
	return  newName
}

// createPreview source and sizes
func createPreview(newName headers.LoadedImage, width, height int) error {
	src, err := imaging.Open(newName.Source)
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	dstImageFill := imaging.Fill(src, width, height, imaging.Center, imaging.Lanczos)

	if err := imaging.Save(dstImageFill, newName.Preview); err != nil {
		return fmt.Errorf("can't save preview img: %s", err)
	}
	return nil
}
