package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// Change the default value of username and password to "no".
var (
	userFlagstart string
	pwdFlagstart  string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Creates and starts a given job.",
	Long: `This command creates and starts a job. Some commands may not work unless it is enclosed in double quotes(""). Streaming commands like top doesn't work.

Usage:
	./client start -u admin -p admin "sleep 10 && echo foo"
`,
	Run: func(cmd *cobra.Command, args []string) {
		command := strings.Join(args, " ")
		startRequest(command, userFlagstart, pwdFlagstart)

	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&userFlagstart, "username", "u", "admin", "Username of the command issuing User.")
	startCmd.Flags().StringVarP(&pwdFlagstart, "password", "p", "admin", "Password corresponding to the username.")

	// startCmd.MarkFlagRequired("username")
	// startCmd.MarkFlagRequired("password")

}
