package app

import (
	"fmt"
	"github.com/cqroot/tinyserver/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Run(bindIp string, bindPort int, whitelist []string) error {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.Use(middleware.WhitelistMiddleware(whitelist))
	app.StaticFS("/", http.Dir("."))

	addr := fmt.Sprintf("%s:%d", bindIp, bindPort)
	return app.Run(addr)
}
