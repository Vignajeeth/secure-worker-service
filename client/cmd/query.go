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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
