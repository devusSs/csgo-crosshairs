package middleware

import (
	"github.com/devusSs/crosshairs/stats"
	"github.com/gin-gonic/gin"
)

func CountRequestsMiddleware(c *gin.Context) {
	stats.RequestsInLast24Hours++
}
