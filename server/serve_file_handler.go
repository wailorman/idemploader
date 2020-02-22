package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func serveFileHandler(storage *Storage) func(*gin.Context) {
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
