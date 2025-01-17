package app

import (
	"fmt"
	"net/http"

	"github.com/cqroot/tinyserver/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Run(bindIp string, bindPort int, whitelist []string) error {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.Use(middleware.WhitelistMiddleware(whitelist))
	app.StaticFS("/", http.Dir("."))

	addr := fmt.Sprintf("%s:%d", bindIp, bindPort)
	return app.Run(addr)
}
