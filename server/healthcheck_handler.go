package server

import "github.com/gin-gonic/gin"

func healthcheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "OK",
		"service": "idemploader",
	})
}
