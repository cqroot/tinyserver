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

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func RunDumpConfigCmd(cmd *cobra.Command, args []string) {
	fmt.Println(`bind_ip: 127.0.0.1
bind_port: 9876
whitelist:
  - 192.168.0.10
  - 192.168.0.11
	`)
}

func NewDumpConfigCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "dump-config",
		Short: "Dump default configuration",
		Long:  "Dump default configuration",
		Run:   RunDumpConfigCmd,
	}

	return &c
}
