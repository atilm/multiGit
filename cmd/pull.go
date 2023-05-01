/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	mgit "atilm/mgit/domain"
	"atilm/mgit/utilities"

	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "performs git pull on all git subdirectories",
	Long: `
When called without arguments all repositories are pulled.
When called with a list of repo-indices only the specified repos will be pulled. E.g.
mgit pull 3 5 6
pulls the repositories with the indices 3, 5 and 6 as indicated by mgit status.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		error := mgit.Pull(MgitBaseDirectory, args, utilities.NewLiveConsolePrinter())
		return error
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
