package middleware

import (
	"net/http"
	"strings"

	"github.com/devusSs/crosshairs/api/responses"
	"github.com/devusSs/crosshairs/database"
	"github.com/gin-gonic/gin"
)

var (
	Svc database.Service
)

func VerifyEngineerMiddleware(c *gin.Context) {
	if c.Request.Header.Get("Authorization") == "" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "Missing engineer token header."
		c.AbortWithStatusJSON(resp.Code, resp)
		return
	}

	authHeader := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]

	latestToken, err := Svc.GetLatestEngineerToken()
	if err != nil {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusInternalServerError
		resp.Error.ErrorCode = "internal_error"
		resp.Error.ErrorMessage = "Something went wrong, sorry."
		c.AbortWithStatusJSON(resp.Code, resp)
		return
	}

	if authHeader != latestToken {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "Invalid engineer token."
		c.AbortWithStatusJSON(resp.Code, resp)
		return
	}

	c.Next()
}
