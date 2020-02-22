package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthConfig _
type AuthConfig struct {
	AllowedAccessToken string
}

// AccessTokenHeaderName _
const AccessTokenHeaderName = "X-Access-Token"

// ErrMissingAccessToken _
var ErrMissingAccessToken = errors.New("MISSING_ACCESS_TOKEN")

// ErrNotAuthorized _
var ErrNotAuthorized = errors.New("NOT_AUTHORIZED")

func authMiddleware(cfg AuthConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.Header.Get(AccessTokenHeaderName) == cfg.AllowedAccessToken {
			c.Next()
			return
		}

		if c.Request.Header.Get(AccessTokenHeaderName) == "" {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, ErrorResponse{
				Code: ErrMissingAccessToken.Error(),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
			Code: ErrNotAuthorized.Error(),
		})
	}
}
