package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var MgitBaseDirectory string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mgit",
	Short: "Apply git commands to multiple repositories at once",
	Long: `mgit lets you execute git status and git pull for all git repositories within a directory (non-recursively).
	
To display the status of all repositories in the parent directory of the current directory:
mgit status -d ..

To pull all repositories in the current directory:
mgit pull

To pull only repositories within indices 2 and 3 in the parent directory:
mgit pull 2 3 -d ..
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&MgitBaseDirectory, "dir", "d", ".", "Specify the directory which contains the repositories. e.g. -d ..")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mgit.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
