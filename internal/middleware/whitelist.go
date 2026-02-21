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

package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func WhitelistMiddleware(whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIp := c.ClientIP()
		if len(whitelist) == 0 {
			return
		}
		if !slices.Contains(whitelist, clientIp) {
			slog.Warn("Forbidden request.", slog.String("ip", clientIp))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Client IP %s denied", clientIp),
			})
			return
		}
	}
}
