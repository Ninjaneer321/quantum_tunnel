package cmd

import (
	"fmt"

	"github.com/SneakyBeagle/quantum_tunnel/libquantum"

	"github.com/spf13/cobra"
)

var cmdVersion *cobra.Command

func runVersion(cmd *cobra.Command, args []string) error {
	verbose, err := rootCmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}
	if verbose == true {
		fmt.Println(libquantum.NAME + " - v" + libquantum.VERSION + " by " + libquantum.AUTHOR)
	} else {
		fmt.Println(libquantum.VERSION)
	}
	return nil
}

func init() {

	cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		RunE:  runVersion,
	}
	rootCmd.AddCommand(cmdVersion)
}
