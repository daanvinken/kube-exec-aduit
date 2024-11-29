package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kube-exec",
	Short: "A secure and auditable Kubernetes exec tool",
	Long: `kube-exec is a CLI tool designed to securely and efficiently execute commands
in Kubernetes containers while maintaining strict auditing and access control.`,
	// Run: func(cmd *cobra.Command, args []string) {
	//     // Optional: Add default behavior for the root command
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This function is called by main.main() and only needs to happen once.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Define persistent flags (global across all commands)
	// For example, a global configuration file:
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kube-exec.yaml)")

	// Define local flags (specific to the root command)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
