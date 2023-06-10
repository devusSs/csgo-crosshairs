package middleware

import (
	"net/http"
	"strings"

	"github.com/devusSs/crosshairs/api/responses"
	"github.com/gin-gonic/gin"
)

var (
	AdminToken string // Generated / assigned upon starting the app.
)

// This route is being used ADDITIONALLY to the usual admin check in each route via a user's session.
func CheckAdminTokenMiddleware(c *gin.Context) {
	authToken := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]

	if authToken == "" {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "Missing admin token header."
		c.AbortWithStatusJSON(resp.Code, resp)
		return
	}

	if authToken != AdminToken {
		resp := responses.ErrorResponse{}
		resp.Code = http.StatusUnauthorized
		resp.Error.ErrorCode = "unauthorized"
		resp.Error.ErrorMessage = "Invalid admin token header."
		c.AbortWithStatusJSON(resp.Code, resp)
		return
	}

	c.Next()
}
