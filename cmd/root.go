package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "quantum_tunnel",
	Short: "Quantum Tunnel is a simple SSH Reverse Forward tunnel",
	Long:  "",

	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("testinnnnggg")
	//},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
}
