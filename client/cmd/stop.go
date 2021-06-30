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
	Short: "Stops a job given its job ID.",
	Long: `This command stops a job. 

Usage:
	./client stop -u admin -p admin -i 1
`,
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
