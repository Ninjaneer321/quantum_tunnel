# Quantum Tunnel

A simple Go tool to create a reverse forward SSH tunnel (and teach me Go in the process). Allows users to get around firewall/NAT restrictions, by creating an outgoing connection, instead of an incoming one. Supports multiple connections.

- [Idea](#idea)
- [Usage](#usage)

## Idea
The tunnel is very similar to the regular SSH reverse forward tunnel which can be created with the following command:
```
ssh -R 5555:localhost:22 root@example.com -p 22
```
The setup will consist of at least two machines, a server and a host in the target network.
The servers ssh port should be reachable for the host in the target network and should accept the public key of that host.
The host in the target network should be able to make an outgoing request to the server and should be able to perform TCP forwarding.

It is also possible to allow the host to act as a proxy and create a tunnel between the server and another machine that would otherwise not be reachable for the server.

### Example:
The host will run Quantum Tunnel and forward its own port 22 (localhost:22) to the server, where the server will create an opening to the tunnel at port 5555 (localhost:5555). Logging into the server and connecting to localhost:5555 will forward the traffic to the host in the network.

## Usage
You run these commands from the machine that is not reachable from outside the network, but has outgoing access.
- Get a help message for the tunnel
```
./quantum_tunnel tunnel --help
```

- Use the default settings:
```
./quantum_tunnel tunnel --server-ip example.com
```
The ssh equivalent of this is:
```
ssh -R 5555:localhost:22 root@example.com -p 22
```

- Create a tunnel between example.com:443 to 192.168.100.100:80 using the local machine as the "proxy" (and using a different private key). This requires example.com to listen for ssh connections on port 443
```
./quantum_tunnel tunnel --server-ip example.com --server-port 443 --server-user admin --remote-ip 192.168.100.100 --remote-port 80 --identity-file ./id_rsa
```
Or the short version:
```
./quantum_tunnel tunnel -s example.com -p 443 -u admin -r 192.168.100.100 -a 80 -i ./id_rsa
```
The ssh equivalent of this is:
```
ssh -R 5555:192.168.100.100:80 admin@example.com -p 443
```
