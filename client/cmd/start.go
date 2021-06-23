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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
