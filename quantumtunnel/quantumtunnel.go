package quantumtunnel

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"sync"

	"github.com/SneakyBeagle/quantum_tunnel/libquantum"

	"golang.org/x/crypto/ssh"
)

// Handle a client connection through the tunnel
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	go func() {
		bytes, err := io.Copy(client, remote)
		fmt.Printf("Wrote %d bytes to client\n", bytes)
		if err != nil {
			log.Println(fmt.Sprintf("Error while copying remote -> local: %s", err))
		}
		chDone <- true
	}()

	go func() {
		bytes, err := io.Copy(remote, client)
		fmt.Printf("Wrote %d bytes to remote\n", bytes)
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
	//fmt.Printf("%x\n", string(key.PublicKey().Marshal()))
	fmt.Printf("%s\n", string(key.PublicKey().Type()))
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

	var wg sync.WaitGroup
	nr_clients := 0
	wg.Add(1)
	// Run loop
	for {
		local, err := net.Dial("tcp", remote.String())
		if err != nil {
			log.Fatalln(fmt.Printf("Dial INTO remote service error: %s", err))
		}

		client, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go func(client net.Conn, local net.Conn) {
			nr_clients++
			fmt.Printf("Number of clients: %d\n", nr_clients)
			wg.Add(1)
			defer wg.Done()
			handleClient(client, local)
			nr_clients--
			fmt.Printf("Number of clients: %d\n", nr_clients)
		}(client, local)
	}

	wg.Wait()
	fmt.Println("Finished")
}
