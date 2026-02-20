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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var usedConfig string

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigFile("./tinyserver.yaml")
	if err := viper.ReadInConfig(); err == nil {
		usedConfig = viper.ConfigFileUsed()
	}
}

func RunRootCmd(cmd *cobra.Command, args []string) {
	if usedConfig != "" {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	bindIp := viper.GetString("bind_ip")
	bindPort := viper.GetInt("bind_port")
	whitelist := viper.GetStringSlice("whitelist")

	cobra.CheckErr(app.Run(bindIp, bindPort, whitelist))
}

func NewRootCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "tinyserver",
		Short: "A tiny server",
		Long:  "A tiny server",
		Run:   RunRootCmd,
	}

	c.PersistentFlags().StringP("bind_ip", "i", "", "bind ip")
	c.PersistentFlags().IntP("bind_port", "p", 9876, "bind port")
	c.PersistentFlags().StringArrayP("whitelist", "w", nil, "whitelist")
	cobra.CheckErr(viper.BindPFlag("bind_ip", c.PersistentFlags().Lookup("bind_ip")))
	cobra.CheckErr(viper.BindPFlag("bind_port", c.PersistentFlags().Lookup("bind_port")))
	cobra.CheckErr(viper.BindPFlag("whitelist", c.PersistentFlags().Lookup("whitelist")))

	c.AddCommand(NewDumpConfigCmd())

	return &c
}

func Execute() {
	cobra.CheckErr(NewRootCmd().Execute())
}
