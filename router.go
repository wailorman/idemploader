package idemploader

import (
	"errors"
	"fmt"
	"net/http"

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

	storage, err := NewStorage(StorageConfig{
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

	router.Use(nice.Recovery(recoveryHandler))

	router.POST("/api/v1/upload", uploadFileHandler(storage))
	router.GET("/api/v1/files/:checksum", getFileHandler(storage))
	router.GET("/api/v1/files/:checksum/info", getFileInfoHandler(storage))

	return router
}

func uploadFileHandler(storage *Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		if c.ContentType() != multipartFormContentType {
			panic(fmt.Errorf("Content-Type `%s` received. Expected %s", c.ContentType(), multipartFormContentType))
		}

		formFile, _ := c.FormFile("file")

		file, err := NewMultipartFile(formFile)

		if err != nil {
			panic(err)
		}

		err = storage.UploadFileIfNotExists(file)

		if err != nil {
			panic(err)
		}

		uploadedFile, err := storage.GetFile(file)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, uploadedFile)
	}
}

func getFileHandler(storage *Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		checksum := c.Param("checksum")

		file, err := storage.GetFileByChecksum(checksum)

		if err != nil {
			if err == ErrFileNotFound {
				c.JSON(http.StatusNotFound, &ErrorResponse{
					Code: err.Error(),
				})
			} else {
				panic(err)
			}
		}

		c.DataFromReader(http.StatusOK, int64(file.Size), file.MimeType, file, map[string]string{})
	}
}

func getFileInfoHandler(storage *Storage) func(*gin.Context) {
	return func(c *gin.Context) {
		checksum := c.Param("checksum")

		file, err := storage.GetFileByChecksum(checksum)

		if err != nil {
			if err == ErrFileNotFound {
				c.JSON(http.StatusNotFound, ErrorResponse{
					Code: ErrFileNotFoundCode,
				})

				return
			}

			panic(err)
		}

		c.JSON(http.StatusOK, file)
	}
}

func recoveryHandler(c *gin.Context, err interface{}) {
	var message string

	switch e := err.(type) {
	case error:
		message = e.Error()
	}

	c.JSON(http.StatusInternalServerError, &ErrorResponse{
		Code:    ErrUnknown.Error(),
		Message: message,
	})
}
