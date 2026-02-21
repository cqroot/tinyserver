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
	"log/slog"
	"net"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

func LogAppInfo(bindIp string, bindPort int, whitelist []string) {
	slog.Info("Starting TinyServer.",
		slog.String("bind_addr", net.JoinHostPort(bindIp, strconv.Itoa(bindPort))),
		slog.String("whitelist", strings.Join(whitelist, ", ")),
	)

	if IsWildcardHosts(bindIp) {
		ips, err := GetLocalIps()
		if err != nil {
			return
		}

		slog.Info("Available bind addresses:")
		for _, ip := range ips {
			slog.Info(fmt.Sprintf("  %s http://%s",
				lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render("•"),
				net.JoinHostPort(ip, strconv.Itoa(bindPort)),
			))
		}
	}
}
