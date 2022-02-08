package quantumtunnel

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"github.com/SneakyBeagle/quantum_tunnel/libquantum"

	"golang.org/x/crypto/ssh"
)

// Handle a client connection through the tunnel
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	go func() {
		bytes, err := io.Copy(client, remote)
		fmt.Println("Wrote %d bytes to client", bytes)
		if err != nil {
			log.Println(fmt.Sprintf("Error while copying remote -> local: %s", err))
		}
		chDone <- true
	}()

	go func() {
		bytes, err := io.Copy(remote, client)
		fmt.Println("Wrote %d bytes to remote", bytes)
		if err != nil {
			log.Println(fmt.Sprintf("Error while copying local -> remote: %s", err))
		}
		chDone <- true
	}()

	<-chDone
}

// Read the private key and generate the public key to use
func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot read SSH private key file %s", file))
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	//fmt.Printf(string(key.PublicKey()))
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot parse SSH private key file %s %s", file, err))
		return nil
	}

	return ssh.PublicKeys(key)
}

// Main run function
// server: The server to connect to and acts as one end of the tunnel
// serverEntry: The server listener host:port to act as entrypoint to the tunnel
// remote: The remote (or local) host for the tunnel traffic to be forwarded to
func Tunnel(server *libquantum.Endpoint, serverEntry *libquantum.Endpoint, remote *libquantum.Endpoint,
	usr_server string, privkey_file string) {

	sshConfig := &ssh.ClientConfig{
		User: usr_server,
		Auth: []ssh.AuthMethod{
			publicKeyFile(privkey_file),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	serverConn, err := ssh.Dial("tcp", server.String(), sshConfig)
	if err != nil {
		log.Fatalln(fmt.Printf("Dial INTO server error: %s", err))
	}

	listener, err := serverConn.Listen("tcp", serverEntry.String())
	if err != nil {
		log.Fatalln(fmt.Printf("Listen open port ON server error: %s", err))
	}
	defer listener.Close()

	// Run loop
	for {
		//fmt.Println("test1")
		local, err := net.Dial("tcp", remote.String())
		//fmt.Println("test2")
		if err != nil {
			log.Fatalln(fmt.Printf("Dial INTO remote service error: %s", err))
		}

		//fmt.Println("test3")
		client, err := listener.Accept()
		//fmt.Println("test4")
		if err != nil {
			log.Fatalln(err)
		}

		//fmt.Println("test5")
		handleClient(client, local)
		//fmt.Println("test6")
	}
}
