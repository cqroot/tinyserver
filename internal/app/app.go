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
	"fmt"
	"net"
	"net/http"
	"slices"
	"strconv"

	"github.com/cqroot/tinyserver/internal/middleware"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func IsWildcardHosts(host string) bool {
	wildcardHosts := []string{"", "0.0.0.0", "::"}
	return slices.Contains(wildcardHosts, host)
}

func GetLocalIps() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		ip := ipNet.IP
		if ip.IsLoopback() || ip.IsLinkLocalUnicast() {
			continue
		}
		if ip.IsGlobalUnicast() {
			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}

func PrintAvailableAddrs(port int) {
	color.HiBlue("  Available bind addresses")
	ips, err := GetLocalIps()
	if err != nil {
		return
	}

	for _, ip := range ips {
		fmt.Printf("    %s http://%s\n",
			color.HiBlueString("•"), net.JoinHostPort(ip, strconv.Itoa(port)))
	}
}

func Run(bindIp string, bindPort int, whitelist []string) error {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.Use(middleware.WhitelistMiddleware(whitelist))
	app.StaticFS("/", http.Dir("."))

	addr := net.JoinHostPort(bindIp, strconv.Itoa(bindPort))

	color.HiGreen("[TinyServer] Starting TinyServer.")
	fmt.Printf("  %s  %s:%d\n", color.HiBlueString("Bind Addr"), bindIp, bindPort)
	fmt.Printf("  %s  %v\n", color.HiBlueString("Whitelist"), whitelist)
	if IsWildcardHosts(bindIp) {
		PrintAvailableAddrs(bindPort)
	}
	fmt.Println()

	return app.Run(addr)
}
