package middleware

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func WhitelistMiddleware(whiteList []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIp := c.ClientIP()
		if !slices.Contains(whiteList, clientIp) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Client IP %s denied", clientIp),
			})
			return
		}
	}
}
