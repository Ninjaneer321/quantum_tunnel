package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/SneakyBeagle/quantum_tunnel/libquantum"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "quantum_tunnel",
	Short: "Quantum Tunnel is a simple SSH Reverse Forward tunnel",
	Long:  "",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseGlobalOptions() (*libquantum.Options, error) {
	globalOpts := libquantum.NewOptions()

	verbose, err := rootCmd.Flags().GetBool("verbose")
	if err != nil {
		log.Fatalln(fmt.Sprintf("Could not get verbose flag: %s", err))
	}
	globalOpts.Verbose = verbose

	return globalOpts, nil
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
}
