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

package netutil

import (
	"fmt"
	"net"
	"strconv"
)

// GetLocalIPs returns a list of non-loopback, non-link-local unicast IP addresses
// assigned to the host.
func GetLocalIPs() ([]string, error) {
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

// FindAvailablePort attempts to find an available TCP port starting from startPort
// and trying up to maxAttempts ports. It returns the first available port number,
// or an error if none is found.
func FindAvailablePort(bindIp string, startPort, maxAttempts int) (int, error) {
	for port := startPort; port < startPort+maxAttempts; port++ {
		addr := net.JoinHostPort(bindIp, strconv.Itoa(port))
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found in range %d-%d on %s",
		startPort, startPort+maxAttempts-1, bindIp)
}
