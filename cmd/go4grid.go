package cmd

import (
	"github.com/spf13/cobra"
)

var Namespace string

const (
	CHROME    = "chrome"
	FIREFOX   = "firefox"
	HUB       = "hub"
	NAMESPACE = "namespace"
)

var rootCmd = &cobra.Command{
	Use:   "go4grid",
	Short: "Commandline app to spin up Selenium 4 Grid on Kubernetes",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(terminateCmd)
	rootCmd.AddCommand(statusCmd)
}
