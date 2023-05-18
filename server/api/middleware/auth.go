package middleware

import (
	"net/http"

	"github.com/devusSs/crosshairs/api/responses"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		var resp responses.ErrorResponse
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "Please login first."
		c.AbortWithStatusJSON(resp.Code, resp)
		return
	}
	c.Next()
}
