package cmd

import (
	"github.com/cqroot/tinyserver/internal/app"
	"github.com/spf13/cobra"
)

var (
	BindIp    string
	BindPort  int
	Whitelist []string
)

func RunRootCmd(cmd *cobra.Command, args []string) {
	cobra.CheckErr(app.Run(BindIp, BindPort, Whitelist))
}

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "tinyserver",
		Short: "A tiny server",
		Long:  "A tiny server",
		Run:   RunRootCmd,
	}

	rootCmd.PersistentFlags().StringVarP(&BindIp, "bind_ip", "i", "", "bind ip")
	rootCmd.PersistentFlags().IntVarP(&BindPort, "bind_port", "p", 9876, "bind port")
	rootCmd.PersistentFlags().StringArrayVarP(&Whitelist, "whitelist", "w", nil, "whitelist")

	return &rootCmd
}

func Execute() {
	cobra.CheckErr(NewRootCmd().Execute())
}
