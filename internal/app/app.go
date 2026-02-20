/*
Copyright (C) 2025 Keith Chu <cqroot@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package app

import (
	"net"
	"net/http"
	"slices"
	"strconv"

	"github.com/cqroot/tinyserver/internal/middleware"
	"github.com/gin-gonic/gin"
)

func IsWildcardHosts(host string) bool {
	wildcardHosts := []string{"", "0.0.0.0", "::"}
	return slices.Contains(wildcardHosts, host)
}

func Run(bindIp string, bindPort int, whitelist []string) error {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.Use(middleware.WhitelistMiddleware(whitelist))
	app.StaticFS("/", http.Dir("."))

	addr := net.JoinHostPort(bindIp, strconv.Itoa(bindPort))

	PrintAppInfo(bindIp, bindPort, whitelist)
	return app.Run(addr)
}
