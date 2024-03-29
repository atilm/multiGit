package cmd

import (
	mgit "atilm/mgit/domain"
	util "atilm/mgit/utilities"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Calls git status on each git repository found in subdirectories of the given path",
	Long: `output format:
<repo-index>: <repo name> (<branch name>) (* indicates uncommitted changes) [ok or count of commits to push / pull]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		error := mgit.ReportStatus(MgitBaseDirectory, util.NewLiveConsolePrinter())
		return error
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
