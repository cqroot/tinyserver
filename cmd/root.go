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
	"log/slog"
	"path/filepath"

	"github.com/cqroot/tinyserver/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunRootCmd(cmd *cobra.Command, args []string) {
	workDir := viper.GetString("work_dir")

	bindIp := viper.GetString("bind_ip")
	if cmd.Flags().Changed("bind_ip") {
		bindIp, _ = cmd.Flags().GetString("bind_ip")
	}

	bindPort := viper.GetInt("bind_port")
	if cmd.Flags().Changed("bind_port") {
		bindPort, _ = cmd.Flags().GetInt("bind_port")
	}

	whitelist := viper.GetStringSlice("whitelist")
	if cmd.Flags().Changed("whitelist") {
		whitelist, _ = cmd.Flags().GetStringArray("whitelist")
	}

	a := app.New(workDir)
	cobra.CheckErr(a.Run(bindIp, bindPort, whitelist))
}

func NewRootCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "tinyserver",
		Short: "A tiny server",
		Long:  "A tiny server",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			workDir, err := cmd.Flags().GetString("work_dir")
			cobra.CheckErr(err)

			configPath := filepath.Join(workDir, "tinyserver.yaml")
			viper.SetConfigFile(configPath)
			if err := viper.ReadInConfig(); err == nil {
				slog.Info("Load config file.", slog.String("file", viper.ConfigFileUsed()))
			}

			viper.SetDefault("bind_ip", "")
			viper.SetDefault("bind_port", 9876)
			viper.SetDefault("whitelist", []string{})

			viper.Set("work_dir", workDir)
		},
		Run: RunRootCmd,
	}

	c.PersistentFlags().StringP("work_dir", "d", ".", "working directory")
	c.PersistentFlags().StringP("bind_ip", "i", "", "bind ip")
	c.PersistentFlags().IntP("bind_port", "p", 9876, "bind port")
	c.PersistentFlags().StringArrayP("whitelist", "w", nil, "whitelist")

	c.AddCommand(NewDumpConfigCmd())
	return &c
}

func Execute() {
	cobra.CheckErr(NewRootCmd().Execute())
}
