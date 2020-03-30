package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wailorman/idemploader/server"
)

func main() {
	log.Println("Starting ...")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err.Error())
	}

	r := server.Router(&server.Config{
		S3Host:             os.Getenv("IDEMPLOADER_S3_HOST"),
		S3AccessKey:        os.Getenv("IDEMPLOADER_S3_ACCESS_KEY"),
		S3AccessSecret:     os.Getenv("IDEMPLOADER_S3_ACCESS_SECRET"),
		S3Bucket:           os.Getenv("IDEMPLOADER_S3_BUCKET"),
		S3Path:             os.Getenv("IDEMPLOADER_S3_PATH"),
		Host:               os.Getenv("IDEMPLOADER_HOST"),
		AllowedAccessToken: os.Getenv("IDEMPLOADER_ALLOWED_ACCESS_TOKEN"),
	})

	r.Run(":" + os.Getenv("PORT"))
}
