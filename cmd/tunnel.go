package cmd

import (
	"fmt"
	"log"

	"github.com/SneakyBeagle/quantum_tunnel/libquantum"
	"github.com/SneakyBeagle/quantum_tunnel/quantumtunnel"

	"github.com/spf13/cobra"
)

var cmdTunnel *cobra.Command
var remote *libquantum.Endpoint
var entry *libquantum.Endpoint
var server *libquantum.Endpoint

func runTunnel(cmd *cobra.Command, args []string) error {

	globalOpts, tunnelOpts, err := parseTunnelOptions()

	if err != nil {
		log.Fatalln(fmt.Sprintf("Could not parse options: %s", err))
		return err
	}

	if globalOpts.Verbose {
		fmt.Println("Verbosity is set")
	}

	remote = &libquantum.Endpoint{
		Host: tunnelOpts.RemoteIP,
		Port: tunnelOpts.RemotePort,
	}
	entry = &libquantum.Endpoint{
		Host: tunnelOpts.EntryIP,
		Port: tunnelOpts.EntryPort,
	}
	server = &libquantum.Endpoint{
		Host: tunnelOpts.ServerIP,
		Port: tunnelOpts.ServerPort,
	}
	user := tunnelOpts.ServerUser
	privkey := tunnelOpts.PrivKey

	quantumtunnel.Tunnel(server, entry, remote, user, privkey)

	return nil
}

func parseTunnelOptions() (*libquantum.Options, *quantumtunnel.TunnelOptions, error) {
	globalOpts, err := parseGlobalOptions()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Could not parse global options: %s", err))
		return nil, nil, err
	}

	tunnelOpts := quantumtunnel.NewTunnelOptions()

	tunnelOpts.ServerIP, err = cmdTunnel.Flags().GetString("server-ip")
	if err != nil {
		return nil, nil, err
	}
	tunnelOpts.ServerPort, err = cmdTunnel.Flags().GetInt("server-port")
	if err != nil {
		return nil, nil, err
	}
	tunnelOpts.ServerUser, err = cmdTunnel.Flags().GetString("server-user")
	if err != nil {
		return nil, nil, err
	}

	tunnelOpts.EntryIP, err = cmdTunnel.Flags().GetString("entry-ip")
	if err != nil {
		return nil, nil, err
	}
	tunnelOpts.EntryPort, err = cmdTunnel.Flags().GetInt("entry-port")
	if err != nil {
		return nil, nil, err
	}

	tunnelOpts.RemoteIP, err = cmdTunnel.Flags().GetString("remote-ip")
	if err != nil {
		return nil, nil, err
	}
	tunnelOpts.RemotePort, err = cmdTunnel.Flags().GetInt("remote-port")
	if err != nil {
		return nil, nil, err
	}

	tunnelOpts.PrivKey, err = cmdTunnel.Flags().GetString("identity-file")
	if err != nil {
		return nil, nil, err
	}

	return globalOpts, tunnelOpts, nil
}

func init() {
	cmdTunnel = &cobra.Command{
		Use:   "tunnel",
		Short: "Run the tunnel",
		RunE:  runTunnel,
	}

	// Add flags
	cmdTunnel.Flags().StringP("server-ip", "s", "", "Server IP address or domain name. This is the (probably externally) listening server.")
	cmdTunnel.Flags().IntP("server-port", "p", 22, "Server port. This is the (probably externally) listening servers port.")
	cmdTunnel.Flags().StringP("server-user", "u", "root", "Server username. This is the (probably externally) listening server.")

	cmdTunnel.Flags().StringP("entry-ip", "e", "localhost", "Server IP address to listen on. This is the (probably externally) listening servers entrypoint to the tunnel.")
	cmdTunnel.Flags().IntP("entry-port", "l", 5555, "Server port to listen on. This is the (probably externally) listening servers entrypoint to the tunnel.")

	cmdTunnel.Flags().StringP("remote-ip", "r", "localhost", "Remote IP address or domain name. This is the (probably internal) machine that could otherwise not be reached from outside the network.")
	cmdTunnel.Flags().IntP("remote-port", "a", 22, "Remote port. This is the (probably internal) machine that could otherwise not be reached from outside the network.")

	cmdTunnel.Flags().StringP("identity-file", "i", libquantum.GetHomeDir()+"/.ssh/id_rsa", "SSH identity file (Private Key)")

	rootCmd.AddCommand(cmdTunnel)
}
