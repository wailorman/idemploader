package server

import (
	"errors"
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
	S3Host             string
	S3AccessKey        string
	S3AccessSecret     string
	S3Bucket           string
	S3Path             string
	Host               string
	AllowedAccessToken string
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
			return cfg.Host + "/files/" + checksum
		},
	})

	if err != nil {
		panic(err)
	}

	router.Use(nice.Recovery(recoveryHandler))

	router.POST(
		"/files",
		authMiddleware(AuthConfig{AllowedAccessToken: cfg.AllowedAccessToken}),
		uploadFileHandler(storage),
	)

	router.GET(
		"/files/:checksum",
		serveFileHandler(storage),
	)

	router.GET(
		"/files/:checksum/info",
		authMiddleware(AuthConfig{AllowedAccessToken: cfg.AllowedAccessToken}),
		getFileInfoHandler(storage),
	)

	router.GET(
		"/healthcheck",
		healthcheckHandler,
	)

	return router
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
