package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"

	"golang.org/x/crypto/ssh"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Println(fmt.Sprintf("Error while copying remote -> local: %s", err))
		}
		chDone <- true
	}()

	go func() {
		test, err := io.Copy(remote, client)
		fmt.Println(test)
		if err != nil {
			log.Println(fmt.Sprintf("Error while copying local -> remote: %s", err))
		}
		chDone <- true
	}()

	<-chDone
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot read SSH public key file %s", file))
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	fmt.Printf(string(buffer))
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot parse SSH public key file %s %s", file, err))
		return nil
	}
	return ssh.PublicKeys(key)
}

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func main() {

	starttunnel_addr := flag.String("t", "localhost", "The address of the target host to be forwarded")
	starttunnel_port := flag.Int("tp", 22, "The target port to be forwarded")

	server_addr := flag.String("s", "", "The address of the  host to be forwarded (REQUIRED)")
	server_port := flag.Int("sp", 22, "The port to be forwarded")

	endtunnel_addr := flag.String("e", "localhost", "The address the server will listen to to forward through the tunnel")
	endtunnel_port := flag.Int("ep", 5555, "The port to access the tunnel from server side")

	privkey_file := flag.String("i", getHomeDir()+"/.ssh/id_rsa", "The identity file to use")

	if len(*server_addr) == 0 {
		fmt.Println("Usage: quantum -s <server>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()

	var localEndpoint = Endpoint{
		Host: *starttunnel_addr,
		Port: *starttunnel_port,
	}

	var serverEndpoint = Endpoint{
		Host: *server_addr,
		Port: *server_port,
	}

	var remoteEndpoint = Endpoint{
		Host: *endtunnel_addr,
		Port: *endtunnel_port,
	}
	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			publicKeyFile(*privkey_file),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	serverConn, err := ssh.Dial("tcp", serverEndpoint.String(), sshConfig)
	if err != nil {
		log.Fatalln(fmt.Printf("Dial INTO remote server error: %s", err))
	}

	listener, err := serverConn.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		log.Fatalln(fmt.Printf("Listen open port ON remote server error: %s", err))
	}
	defer listener.Close()

	for {
		local, err := net.Dial("tcp", localEndpoint.String())
		if err != nil {
			log.Fatalln(fmt.Printf("Dial INTO local service error: %s", err))
		}

		client, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		handleClient(client, local)
	}
}
