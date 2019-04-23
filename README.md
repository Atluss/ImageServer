[![GoDoc](https://godoc.org/github.com/Atluss/ImageServer?status.svg)](https://godoc.org/github.com/Atluss/ImageServer)
## ImageServer

In this server you can send image to api and server save this image and creating preview image:

Work schemes:
1. link to image;
2. multipart/form-data;
3. image in base64.

Allows formats: **jpg, jpeg, png**. All images save in `images` folder and if you want you can open it.

Project use [Go modules](https://github.com/golang/go/wiki/Modules) to download dependencies.
You can see dependencies in `go.mod`

Test requests you can find in folder: `test_client` run it, when api is running, don't forget change ports.

## How to run

Don't forget add your user for docker group `sudo usermod -aG docker $USER`. Server run on port `4080`

 1. [Install Docker-CE (ubuntu)](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
 2. [Install Docker compose](https://docs.docker.com/compose/install/)
 3. `sudo docker-compose up`

## Endpoints

### http://localhost:4080/form_data

This endpoint allowed to work in two modes at once: `query` and `multipart/form-data`

1. This endpoint can work with query argument `image` example: 
`http://localhost:4080/form_data?image=https://img.zoneofgames.ru/news/2019/04/22/190216-banner_conk_20190422_PaganOnline.jpg`

2. **multipart/form-data** automatically determined if the argument is a valid format picture.

### http://localhost:4080/json_image

This endpoint only for json base64 image, the body of request, reply have only one image:
1. **Data** - format img base;
2. **Body** - body of image(in the example it is trimmed).

```json
{
	"Data": "data:image/jpeg",
	"Body": "/9j/4AAQSkZJRgABAQEAYABgAAD/"
}
```

### Answer for all endpoints
```json
{
    "Status": 200,
    "Description": "",
    "Images": [
        {
            "Source": "images/74b51062-6b62-481d-acf1-b44d4e08ad91.jpeg",
            "Preview": "images/74b51062-6b62-481d-acf1-b44d4e08ad91_preview.jpeg"
        },
        {
            "Source": "images/a2d3ba2e-5a9c-4840-ade8-8d5c9c81d4c2.jpg",
            "Preview": "images/a2d3ba2e-5a9c-4840-ade8-8d5c9c81d4c2_preview.jpg"
        }
    ]
}
```

### API setup file

Api setup file in [Json RFC7159](https://tools.ietf.org/html/rfc7159) format, all arguments is necessary:
```json
{
  "name": "Image Server",
  "version": "1.0.0",
  "host" : "localhost",
  "port" : "8080"
}
```