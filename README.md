# nConnect

[![GitHub license](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/nknorg/nconnect)](https://goreportcard.com/report/github.com/nknorg/nconnect) [![Build Status](https://travis-ci.org/nknorg/nconnect.svg?branch=master)](https://travis-ci.org/nknorg/nconnect) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](#contributing)

nConnect allows you to securely connect to remote machines without the need of
any server, public IP address, or publicly exposed ports. It features end-to-end
encryption for top level security, and multi-path aggregation for maximum
throughput.

nConnect provides several modes. When using the VPN mode, any TCP-based
application that works in the same local network will continue to work remotely
as if those machines are in the same local network. A TUN device mode and a
SOCKS proxy mode are also available for advanced users.

nConnect uses [nkn-tunnel](https://github.com/nknorg/nkn-tunnel) for end to end
tunneling, thus benefits from all the advantages of
[nkn-tunnel](https://github.com/nknorg/nkn-tunnel):

- Network agnostic: Neither sender nor receiver needs to have public IP address
  or port forwarding. NKN tunnel only establishes outbound (websocket)
  connections, so Internet access is all they need on both sides.

- Top level security: All data are end to end authenticated and encrypted. No
  one else in the world except the sender and receiver can see or modify the content
  of the data. The same public key is used for both routing and encryption,
  eliminating the possibility of man in the middle attack.

- Decent performance: By aggregating multiple overlay paths concurrently, one
  can usually get much higher throughput than direct connection. Even using the
  free mode, one can still get <100ms end to end latency and 10+mbps end to end
  throughput.

- Everything is open source and decentralized. The default free mode is,
  suggested by its name, free of charge (If you are curious, node relay traffic
  for clients for free to earn mining rewards in NKN blockchain); while in tuna
  mode, nConnect (server mode) will pay NKN token directly to relay service
  providers.

## Build

```shell
make
```

## Usage

nConnect needs to be started in either server or client mode, server mode allows
incoming connections from client mode.

When started for the first time, nConnect will generate a config file
`config.json` in the current working directory. This file contains your private
key and should not be shared.

When started in server mode, nConnect might generate some data files in the
current working directory. You will also need the directory `web` located in the
current working directory if admin web dashboard is enabled.

### Server Mode

The minimal arguments to start nConnect in server mode is just

```shell
./nConnect -s
```

But most of the time you might want to start nConnect server with a few useful
arguments:

```shell
./nConnect -s --tuna --admin-http 127.0.0.1:8001
```

- `--tuna` enables tuna mode, which gets much better performance but requires
  you to purchase data plan (you can find the link in admin web dashboard). This
  argument is required if you want to be compatible with nConnect mobile and
  desktop clients.

- `--admin-http 127.0.0.1:8001` starts the admin web dashboard at
  `http://127.0.0.1:8001`. You can visit this address in your browser to change
  various config (e.g. access control), bind with nConnect mobile client, etc.
  Do not make this port public as anyone who can access this endpoint can change
  your configuration. If you want the best security, disable the admin dashboard
  once you have done using it.

#### Access Control

Before you can connect from nConnect client mode, you need to add your nConnect
client address (see [Get Your Client Address](#get-your-client-address) for how
to get it) to allowed addresses. You can do it using the admin web dashboard, or
manually edit the `config.json` file, which will be generated after first run.

There are two lists of allowed address:

- Accept address: address in this list will be able to connect to nConnect
  server.

- Admin address: address in this list will be able to connect to nConnect server
  and manage nConnect server config and permissions (import/export account,
  view/change accept addresses and admin addresses, etc).

Items in each list are regular expressions. If you want to add a nConnect client
public key to the list, it is important that you add `$` to the end to match
the public key part. For example,
`ad37e248005113dd42be15a4885e6446e9e23f35537dfa6c584f2563a7e8f96d$`
will allow any address using this public key, such as
`ad37e248005113dd42be15a4885e6446e9e23f35537dfa6c584f2563a7e8f96d`
and
`nkn.ad37e248005113dd42be15a4885e6446e9e23f35537dfa6c584f2563a7e8f96d`.

#### Get Your Server Address

You will need your nConnect server address in order to connect from nConnect client. You can get your server address using:

```shell
./nConnect -s --address
```

which can be passed to the `-a` argument on the nConnect client side.

### Client Mode

Before connecting to nConnect server, you will first need to set up nConnect
server side correctly. Make sure you have done these:

- Add client address or public key to server's allowed list, see [Access
  Control](#access-control).

- Get server address, see [Get Your Server Address](#get-your-server-address).

When starting nConnect in client mode, you have a few sub-modes as options:

- [VPN Mode](#vpn-mode): TCP connections made to nConnect server's local IP
  address will be captured transparently and tunneled to nConnect server. Most
  applications will work without any further configurations.

- [TUN Device Mode](#tun-device-mode): create a TUN device, TCP connections
  routed via this device will be tunneled to nConnect server.

- [SOCKS Proxy Mode](#socks-proxy-mode): create a local SOCKS proxy, TCP
  connections routed through this proxy will be tunneled to nConnect server.

#### VPN Mode

Start nConnect client in VPN mode requires root privilege in most cases:

```shell
sudo ./nConnect -c -a <server-addr> --tuna --vpn
```

Replace `<server-addr>` with the server address you get in [Get Your Server
Address](#get-your-server-address), and add `--tuna` only if nConnect starts
with `--tuna` as well.

In the console you should see one or more `Adding route <local-ip>/32`. You can
then connect to server machine using any one of these local IP addresses as if
they are in the same local network, e.g. `ssh user@<local-ip>`.

By default, all local IP addresses on the server machine will be added to routes,
but you can manually specify which IP or IP range you would like to route
through the VPN using `--vpn-route` arguments. Use `./nConnect -h` for all available arguments.

If you start multiple nConnect clients in VPN mode, make sure to use different
subnets for both `--tun-addr` and `--tun-gateway` (e.g. `10.0.86.X` for one
client, `10.0.87.X` for another client).

If you are using windows, you will need to install the network adaptor driver
and change adaptor info beforehand. The simplest way of doing that is to install
nConnect client for windows before using nConnect command line version.

#### TUN Device Mode

Start nConnect client in TUN mode requires root privilege in most cases:

```shell
sudo ./nConnect -c -a <server-addr> --tuna --tun
```

Replace `<server-addr>` with the server address you get in [Get Your Server
Address](#get-your-server-address), and add `--tuna` only if nConnect starts
with `--tuna` as well.

After nConnect client is started, the TUN device will be up and running. TCP
connections routed via this device will be tunneled to nConnect server. You will
need to modify system routing table yourself to determine what traffic should be
routed through the TUN device.

You can also change the name, IP, gateway, network mask and DNS resolvers of the TUN device. Use `./nConnect -h` for
all available arguments.

If you start multiple nConnect clients in TUN device mode, make sure to use
different subnets for both `--tun-addr` and `--tun-gateway` (e.g. `10.0.86.X`
for one client, `10.0.87.X` for another client).

If you are using windows, you will need to install the network adaptor driver
and change adaptor info beforehand. The simplest way of doing that is to install
nConnect client for windows before using nConnect command line version.

#### SOCKS Proxy Mode

```shell
./nConnect -c -a <server-addr> --tuna
```

Replace `<server-addr>` with the server address you get in [Get Your Server
Address](#get-your-server-address), and add `--tuna` only if nConnect starts
with `--tuna` as well.

After nConnect client is started, a SOCKS proxy will be listening at
`127.0.0.1:1080`. TCP connections routed through this proxy will be tunneled to
nConnect server. You can change the SOCKS proxy listening address using `-l`
argument. Use `./nConnect -h` for all available arguments.

#### Get Your Client Address

You will need your nConnect client address to add to allowed addresses on
nConnect server side. You can get your client address using:

```shell
./nConnect -c --address
```

The address typically contains one or more dot, with the part after last dot
being your client public key.

### UDP support

You can enable UDP support when starting nConnect server with tuna mode, 

```shell
./nConnect -s --tuna --udp
```

### Use nConnect as library

You can also use nConnect as library. Please check [proxy_test.go](tests/proxy_test.go) for usages.

### Use pre-built Docker image

*Pre-requirement*: Have working docker software installed. For help with that
*visit [official docker
*docs](https://docs.docker.com/install/#supported-platforms)

We host the latest Docker image on our official Docker Hub account. You can get
it by

```shell
$ docker pull nknorg/nconnect
```

and run it with

```shell
docker run --rm -it --net=host -v ${PWD}:/nConnect/data nknorg/nconnect
```

followed by the command line argument you want to add.

## nConnect Client Connects to Multi Servers

Now nConnect client can connect to multiple servers. You can edit `config.json` to add multiple servers admin addresses:

```
{
  "Client": true,
  "Server": false,
  "identifier": "alice",
  "seed": "",
  "remoteAdminAddr": [
      "nConnect.bob1.7cafe0ae02789f8eb6b293e46b0ac5cf8f92f73042199c8161e5b5f90b13dcb5",
      "nConnect.bob2.7cafe0ae02789f8eb6b293e46b0ac5cf8f92f73042199c8161e5b5f90b13dcb5",
      "nConnect.bob3.7cafe0ae02789f8eb6b293e46b0ac5cf8f92f73042199c8161e5b5f90b13dcb5",
  ],
  "localSocksAddr": "127.0.0.1:1080",
}

```

After config multi `remoteAdminAddr`, the nConnect client will add routing information to each Server's local IP. So you can access all the servers by their local IP address.

Specifically, the first item in the `remoteAdminAddr` will become the default server which will get the forwarded data whose targets are beyond all these servers' local IP addresses. Such as access to the website by domain, or some other applications.

You can use command argument to connect to multiple servers too. Use multi times argument `-a` to pass multi servers addresses:

```
nConnect -c -a server-address1 -a server-address2 -a server-address3

```

## Use `config.json` to Simplify Command Arguments

You can use `config.json` to simplify command arguments. Copy config.client.json or config.server.json as `config.json` and edit it before starting your nConnect client or server. After saving `config.json`, you can start nConnect simply.

## Set up a Virtual Private Network by nConnect
Yes, nConnect supports setting up a virtual private network. It means many computers can join a nConnect virtual network, and access each other just like all nodes are in a local network no matter where they are.

In nConnect private virtual network, there are two types of nodes:

* **manager node**
The manager node is the network administrative node that configures network parameters and authorizes network members. It is just like a registry portal and privilege management center.

Based on NKN decentralized network, you can set up nConnect manager node anywhere and only need to connect to the internet. Each nConnect manager has an NKN address, which is used to identify this node and used for other members to register into the network.

* **member node**
The member nodes are the members of the network. All member nodes need to register with the network manager first. After the manager node authorizes the member's `join network` request, the member node will have a network-specific IP and mask. All the member nodes (not including the manager node) can communicate by their network-specific IP no matter where they are. And the data transmitted between network members are encrypted, and high secured.

To set up a nConnect network, you need firstly to start a network manager node. 

### Start a network manager
To start the network manager, we copy `config.network.json` to `config.manager.json`:

```
cp config.network.josn config.manager.json

```

Then edit config.json, enable `NetworkManager`, and give a value to "identifier", just like below:

```
{
    "identifier": "manager",
    "AdminHTTPAddr": "127.0.0.1:8000",
}

```

> If you want to access your manager web page from another computer, you may set `AdminHTTPAddr` in a public IP instead of `127.0.0.1`. But it is not safe and not recommanded. After finishing your network configuration, had better change it back to `127.0.0.1`. It means people can access this web page only from this computer.

Then start nConnect as a network manager with parameter `-m -f config.manager.json` :

```
./nConnect -m -f config.manager.json

```

After nConnect network manager starts, you can see a console printed message:

```
nConnect network manager is listening at: manager.0ec192083....
Network manager web serve at:  http://127.0.0.1:8000/network
```

Copy this listening address: `manager.0ec192083....`, it is the manager's address. Other member nodes need this address to join this network.
After the manager starts, you can visit the web service `http://127.0.0.1:8000/network` (default), to config, to manage the network.

If you want to access nConnect manager from a public IP, you may configure `AdminHTTPAddr` with your computer's public IP.  But do remember that other people can access your manager web page too. After configuring your network, you had better disable `AdminHTTPADDR` and set it to "127.0.0.0" or empty.

### Start network member and join the network
On another computer, you can start a network member, and let it join the nConnect which you start above.
First, you copy `config.network.json` to `config.member.json`

```
cp config.network.json config.member.json

```

Then edit `config.member.json` to edit `identifier`, `managerAddress` and `nodeName`.

```
{
    "identifier": "alice",
    "managerAddress": "manager.0ec192083....",
    "nodeName": "alice",
    "seed": "...",
    "AdminHTTPAddr": "127.0.0.1:8000",
}

```

Set `managerAddress` as your network manager's listening address, and identify your node name `nodeName`. Each network member should have a different `nodeName`.
The field `seed` is the seed of the wallet which you use to pay for the `tuna` fee. Please keep it secured. If your wallet has zero balance, then nConnect Server cannot start at `tuna` mode.

> On Unix-like systems you may need to preface commands with `sudo`, while on Windows you will need to use an `administrator-mode` command prompt. 

Then you can start this node to join the network:

```
sudo ./nConnect -n -s -c -f config.member.json --tuna --vpn --udp
```

or 

```
./nConnect.exe -n -s -c -f config.member.json --tuna --vpn --udp
```

`-n` means this is a network member `node`
For a network member, you may start both `-c` client, and `-s` server, which means you can access other nodes, and other nodes can access you too. 
Or you can only set `-c`, which means you can access other nodes, but you don't want other nodes to access you.
Or you can only set `-s`, which means you can only be accessed, and you don't want to access other nodes.

> A nice tips, when you start nConnect with parameters `-s`, `-tuna`, it means you start nConnect Server and connect to `TUNA` service providers, you need make sure your seed's wallet have NKN tokens, which is used for paying `TUNA` service. And don't worry, it's definitely a low cost for data transmitting compare to other type tunneling service.

#### How to join nConnect network without NKN balance

If you only want to join the nConnect network as a client, it means you can access other member nodes, but other nodes needn't access your node. You can start nConnect without parameter `-s`, which means it will not start nConnect server, and won't spend any NKN tokens.

This is especially useful when you only want to test the network functions and works for most network members.

This is to start a node to join the network without starting nConnect server, and needn't spend any NKN tokens.

```
sudo ./nConnect -n -c -f config.member.json --tuna --vpn --udp
```

or 

```
./nConnect.exe -n -c -f config.member.json --tuna --vpn --udp
```

### Manage the network

When a network member starts, it first will send a `JoinNetwork` message to the network manager.
After the network administrator should open the manager's web administrate page `http://127.0.0.1:8000/network` (default), to configure the network name, IP range, netmask, and gateway.

There are two lists on the manager's web page:

* Waiting for Authorization
  This lists all the nodes which are waiting for authorization to join this network. The administrator can accept it or reject it.
  Only authorized nodes can become network members and will get a network-specific IP address.
  When authorizing a node, it will pop up a dialog to set this node's permission to other nodes which decides if all members or only some of them can access this node.

* Network Members
  In the network members list, the administrator can reset nodes' access permission and remove a node from the network (authorization).

If you don't see your node information in `Waiting for Authorization`, please click the `Refresh` button to fetch updated data from the manager.

### Test your network

To test your network, you can run a TCP/UDP server on a member node, and run a TCP/UDP client on another member node to do some echo tests.

* Start a TCP and UDP server, so another node can access your node

```
go run tests/tools/main.go -server
```

* Start a TCP client to access another node if you know his IP, such as making an echo test to `10.0.86.3` node:

```
go run tests/tools/main.go -serverAddr 10.0.86.3
```

* Start a UDP client to access another node if you know his IP, such as making an echo test to `10.0.86.3` node:

```
go run tests/tools/main.go -serverAddr 10.0.86.3 -udp
```

You should see both the server and client's echo test messages.

### Interact with the nConnect node by command line interface

Now we provide a command line interface to interact with the running nConnect process.

```
./nConnect -i <cmd>
```
The <cmd> can be:

```
help:        this help
join:        join network
leave:       leave network
status:      get network status
list:       list nodes I can access and nodes which can access me.
```

You can input these sub-commands to interact with the nConnect network member:

* join: to join a network that is configured with a manager address;
* leave: to leave a network;
* status: to show your nConnect network member information, such as your node's IP information.
* list: to list nodes I can access and nodes that can access me.

## Contributing

**Can I submit a bug, suggestion or feature request?**

Yes. Please open an issue for that.

**Can I contribute patches?**

Yes, we appreciate your help! To make contributions, please fork the repo, push
your changes to the forked repo with signed-off commits, and open a pull request
here.

Please sign off your commit. This means adding a line "Signed-off-by: Name
<email>" at the end of each commit, indicating that you wrote the code and have
the right to pass it on as an open source patch. This can be done automatically
by adding -s when committing:

```shell
git commit -s
```

## Community

- [Forum](https://forum.nkn.org/)
- [Discord](https://discord.gg/c7mTynX)
- [Telegram](https://t.me/nknorg)
- [Reddit](https://www.reddit.com/r/nknblockchain/)
- [Twitter](https://twitter.com/NKN_ORG)
