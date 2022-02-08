package quantumtunnel

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

func NewTunnelOptions() *TunnelOptions {
	return &TunnelOptions{}
}
