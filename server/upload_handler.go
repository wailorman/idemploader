package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
