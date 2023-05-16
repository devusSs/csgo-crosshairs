package routes

import (
	"fmt"
	"net/http"

	"github.com/devusSs/crosshairs/api/responses"
	"github.com/devusSs/crosshairs/config"
	"github.com/devusSs/crosshairs/database"
	"github.com/gin-gonic/gin"
)

var (
	Svc database.Service
	CFG *config.Config
)

func NotFoundRoute(c *gin.Context) {
	var resp responses.ErrorResponse
	resp.Code = http.StatusNotFound
	resp.Error.ErrorCode = "not_found"
	resp.Error.ErrorMessage = fmt.Sprintf("route %s does not exist", c.Request.URL)
	c.JSON(resp.Code, resp)
}

func MethodNotAllowedRoute(c *gin.Context) {
	var resp responses.ErrorResponse
	resp.Code = http.StatusMethodNotAllowed
	resp.Error.ErrorCode = "method_invalid"
	resp.Error.ErrorMessage = fmt.Sprintf("method %s not allowed on %s", c.Request.Method, c.Request.URL)
	c.JSON(resp.Code, resp)
}

func HomeRoute(c *gin.Context) {
	var resp responses.SuccessResponse
	resp.Code = http.StatusOK
	resp.Data = gin.H{"message": "Welcome to the Crosshairs API!"}
	c.JSON(resp.Code, resp)
}
