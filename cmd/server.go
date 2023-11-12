/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	commonflags "github.com/seba-ban/seen-places/commonFlags"
	"github.com/seba-ban/seen-places/server"
	"github.com/spf13/cobra"
)

var cfg = &server.ServerConfig{}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: move to flags
		cfg.TemplatesPath = "templates"
		cfg.TemplatesExt = ".tmpl"
		cfg.Run()
	},
}

func init() {
	commonflags.AddDbConfigFlags(serverCmd)
	rootCmd.AddCommand(serverCmd)
}
