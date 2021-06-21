package helper

import (
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func UploadImage(file interface{}) (*uploader.UploadResult, error) {
	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))

	if err != nil {
		log.Fatal("Error creating instance cloudinary")
	}
	var ctx = context.Background()

	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{})

	if err != nil {
		return resp, err
	}

	return resp, nil
}
