package console

import (
	"github.com/abandon1a2b/kuai/console/cmd"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kuai",
	Short: "kkkkk",
	Long:  `kuai kuai kuai kuai kuai`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	// Run: runWeb,
}

func init() {
	rootCmd.AddCommand(cmd.GetCommands()...)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
