package cmd

import (
	"github.com/spf13/cobra"
)

var (
	userFlagQuery string
	pwdFlagQuery  string
	idFlagQuery   int
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Returns the status of a job given its job ID.",
	Long: `This command queries the status of the job. The status code correspoding to the status is as follows:
	Created:   0
	Running:   1
	Completed: 2
	Failed:    3

Usage:
	./client query -u admin -p admin -i 2
	
	   `,
	Run: func(cmd *cobra.Command, args []string) {
		queryRequest(idFlagQuery, userFlagQuery, pwdFlagQuery)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringVarP(&userFlagQuery, "username", "u", "admin", "Username of the command issuing User.")
	queryCmd.Flags().StringVarP(&pwdFlagQuery, "password", "p", "admin", "Password corresponding to the username.")
	queryCmd.Flags().IntVarP(&idFlagQuery, "job-id", "i", -1, "Job Id of the job to be fetched.")

	// queryCmd.MarkFlagRequired("username")
	// queryCmd.MarkFlagRequired("password")
	queryCmd.MarkFlagRequired("job-id")
}
