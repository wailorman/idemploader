package idemploader

import (
	"errors"
	"fmt"

	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
)

// ErrorResponse _
type ErrorResponse struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

// Config _
type Config struct {
	S3Host         string
	S3AccessKey    string
	S3AccessSecret string
	S3Bucket       string
	S3Path         string
	Host           string
}

// ErrUnknown _
var ErrUnknown = errors.New("UNKNOWN_ERROR")

const multipartFormContentType = "multipart/form-data"

// Router _
func Router(cfg *Config) *gin.Engine {
	router := gin.Default()

	router.Use(nice.Recovery(recoveryHandler))

	router.POST("/api/v1/upload", func(c *gin.Context) {
		if c.ContentType() != multipartFormContentType {
			panic(fmt.Errorf("Content-Type `%s` received. Expected %s", c.ContentType(), multipartFormContentType))
		}

		formFile, _ := c.FormFile("file")

		file, err := NewMultipartFile(formFile)

		if err != nil {
			panic(err)
		}

		uploader, err := NewStorage(&StorageConfig{
			S3Host:         cfg.S3Host,
			S3AccessKey:    cfg.S3AccessKey,
			S3AccessSecret: cfg.S3AccessSecret,
			S3Bucket:       cfg.S3Bucket,
			S3Path:         cfg.S3Path,
			URLBuilder: func(checksum string) string {
				return cfg.Host + "/api/v1/files/" + checksum
			},
		})

		if err != nil {
			panic(err)
		}

		err = uploader.UploadFileIfNotExists(file)

		if err != nil {
			panic(err)
		}

		uploadedFile, err := uploader.GetFile(file)

		if err != nil {
			panic(err)
		}

		c.JSON(200, uploadedFile)
	})

	return router
}

func recoveryHandler(c *gin.Context, err interface{}) {
	var message string

	switch e := err.(type) {
	case error:
		message = e.Error()
	}

	c.JSON(500, &ErrorResponse{
		Code:    ErrUnknown.Error(),
		Message: message,
	})
}
