package cmd

import (
	"github.com/spf13/cobra"
)

var (
	userFlagstop string
	pwdFlagstop  string
	idFlagstop   int
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		stopRequest(idFlagstop, userFlagstart, pwdFlagstart)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().StringVarP(&userFlagstop, "username", "u", "admin", "Username of the command issuing User.")
	stopCmd.Flags().StringVarP(&pwdFlagstop, "password", "p", "admin", "Password corresponding to the username.")
	stopCmd.Flags().IntVarP(&idFlagstop, "job-id", "i", -1, "Job Id of the job to be stopped.")

	// stopCmd.MarkFlagRequired("username")
	// stopCmd.MarkFlagRequired("password")
	stopCmd.MarkFlagRequired("job-id")
}
