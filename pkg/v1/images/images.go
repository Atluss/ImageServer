package images

import (
	"bytes"
	"encoding/base64"
	"fmt"
	v1 "github.com/Atluss/ImageServer/pkg/v1"
	"github.com/Atluss/ImageServer/pkg/v1/headers"

	"github.com/disintegration/imaging"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	ImageFolder   = "images"
	PreviewWidth  = 100
	PreviewHeight = 100
)

var AllowFormats = [3]string{"jpg", "jpeg", "png"}

// GetImagesFormDataAndQuery search all images in multipart/form-data or query request
func GetImagesFormDataAndQuery(r *http.Request) []headers.LoadedImage {
	images := getImagesFormData(r)
	if tmp, err := getImageFromLink(r); err != nil {
		log.Println(err)
	} else {
		images = append(images, tmp)
	}
	return images
}

// createTestFilesDir создаем папку для тестовых файлов отчетов.
func CreateTestFilesDir(dirName string) error {
	dirname := fmt.Sprintf("./%s/", dirName)
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		if err := os.Mkdir(dirname, 0777); err != nil {
			log.Printf("ERROR: can't create dir, %s", err)
			return err
		}
	}
	return nil
}

// GetJsonImageBase64 save image from json base64 image
func GetJsonImageBase64(r *http.Request) (loadedImg headers.LoadedImage, err error) {
	req := headers.RequestCreateImgJsonBase64{}
	if err := req.Decode(r); err != nil {
		return loadedImg, err
	}
	format := ""
	for _, fr := range AllowFormats {
		if strings.Contains(req.Data, fr) {
			format = fr
		}
	}
	if format == "" {
		return loadedImg, fmt.Errorf("format image file not allow")
	}
	if img, err := base64.StdEncoding.DecodeString(req.Body); err != nil {
		return loadedImg, err
	} else {
		for {
			loadedImg = GenerateName(format)
			if err := v1.CheckFileExist(loadedImg.Source); err != nil {
				break
			}
		}
		out, err := os.Create(loadedImg.Source)
		if err != nil {
			log.Println(err)
		}
		defer out.Close()
		// Write the body img to file
		r := bytes.NewReader(img)
		_, err = io.Copy(out, r)
		if err != nil {
			log.Println(err)
		}
		if err := createPreview(loadedImg, PreviewWidth, PreviewHeight); err != nil {
			return loadedImg, err
		}
	}
	return loadedImg, err
}

func getImagesFormData(r *http.Request) (images []headers.LoadedImage) {
	reader, err := r.MultipartReader()
	if err != nil {
		log.Println(err)
		return images
	}
	for {
		part, err_part := reader.NextPart()
		if err_part == io.EOF {
			break
		}
		if strings.Contains(part.Header.Get("Content-Type"), "image") && part.FileName() != "" {
			loadedImg, err := CreateImageName(part.FileName())
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
func getImageFromLink(r *http.Request) (loadedImg headers.LoadedImage, err error) {
	link := r.URL.Query().Get("image")
	if link == "" {
		return loadedImg, fmt.Errorf("no link send")
	}
	loadedImg, err = CreateImageName(link)
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

// CreateImageName generate new name and preview name from source name
func CreateImageName(name string) (newName headers.LoadedImage, err error) {
	arr := strings.Split(name, ".")
	if len(arr) < 2 {
		return newName, err
	}
	format := arr[len(arr)-1]
	// yes it rude but right now perfect no need cycle :)
	if format != AllowFormats[0] && format != AllowFormats[1] && format != AllowFormats[2] {
		return newName, fmt.Errorf("format image file not allow")
	}
	for {
		newName = GenerateName(format)
		if err := v1.CheckFileExist(newName.Source); err != nil {
			break
		}
	}
	return newName, err
}

// GenerateName
func GenerateName(format string) (newName headers.LoadedImage) {
	shortName := uuid.NewV4()
	newName.Source = fmt.Sprintf("%s/%s.%s", ImageFolder, shortName, format)
	newName.Preview = fmt.Sprintf("%s/%s_preview.%s", ImageFolder, shortName, format)
	return newName
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
