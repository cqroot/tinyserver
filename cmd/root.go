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

	"github.com/cqroot/tinyserver/internal/app"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigFile("./tinyserver.yaml")
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func RunRootCmd(cmd *cobra.Command, args []string) {
	bindIp := viper.GetString("bind_ip")
	bindPort := viper.GetInt("bind_port")
	whitelist := viper.GetStringSlice("whitelist")

	color.HiGreen("[TinyServer] Starting TinyServer.")
	fmt.Printf("  %s: %s:%d\n", color.HiBlueString("Bind Addr"), bindIp, bindPort)
	fmt.Printf("  %s: %v\n", color.HiBlueString("Whitelist"), whitelist)
	fmt.Println()

	cobra.CheckErr(app.Run(bindIp, bindPort, whitelist))
}

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "tinyserver",
		Short: "A tiny server",
		Long:  "A tiny server",
		Run:   RunRootCmd,
	}

	rootCmd.PersistentFlags().StringP("bind_ip", "i", "", "bind ip")
	rootCmd.PersistentFlags().IntP("bind_port", "p", 9876, "bind port")
	rootCmd.PersistentFlags().StringArrayP("whitelist", "w", nil, "whitelist")
	cobra.CheckErr(viper.BindPFlag("bind_ip", rootCmd.PersistentFlags().Lookup("bind_ip")))
	cobra.CheckErr(viper.BindPFlag("bind_port", rootCmd.PersistentFlags().Lookup("bind_port")))
	cobra.CheckErr(viper.BindPFlag("whitelist", rootCmd.PersistentFlags().Lookup("whitelist")))

	return &rootCmd
}

func Execute() {
	cobra.CheckErr(NewRootCmd().Execute())
}
