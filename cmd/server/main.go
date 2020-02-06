package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wailorman/idemploader"
)

func main() {
	log.Println("Starting ...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := idemploader.Router(&idemploader.Config{
		S3Host:         os.Getenv("IDEMPLOADER_S3_HOST"),
		S3AccessKey:    os.Getenv("IDEMPLOADER_S3_ACCESS_KEY"),
		S3AccessSecret: os.Getenv("IDEMPLOADER_S3_ACCESS_SECRET"),
		S3Bucket:       os.Getenv("IDEMPLOADER_S3_BUCKET"),
		S3Path:         os.Getenv("IDEMPLOADER_S3_PATH"),
		Host:           os.Getenv("IDEMPLOADER_HOST"),
	})

	r.Run(":" + os.Getenv("IDEMPLOADER_PORT"))
}
