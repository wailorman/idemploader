package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
