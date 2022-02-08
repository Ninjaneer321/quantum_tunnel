package quantumtunnel

import "fmt"

type TunnelOptions struct {
	ServerIP   string
	ServerPort int
	ServerUser string

	EntryIP   string
	EntryPort int

	RemoteIP   string
	RemotePort int

	PrivKey string
}

func (options *TunnelOptions) String() string {
	return fmt.Sprintf("\tServerIP: %s\n\tServerPort: %d\n\tServerUser: %s\n\t"+
		"EntryIP: %s\n\tEntryPort: %d\n\t"+
		"RemoteIP: %s\n\tRemotePort: %d\n\t"+
		"PrivKey: %s\n",
		options.ServerIP, options.ServerPort,
		options.ServerUser, options.EntryIP,
		options.EntryPort, options.RemoteIP,
		options.RemotePort, options.PrivKey)
}

func NewTunnelOptions() *TunnelOptions {
	return &TunnelOptions{}
}
