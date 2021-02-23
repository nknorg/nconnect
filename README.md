# nConnect

[![GitHub license](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/nknorg/nconnect)](https://goreportcard.com/report/github.com/nknorg/nconnect) [![Build Status](https://travis-ci.org/nknorg/nconnect.svg?branch=master)](https://travis-ci.org/nknorg/nconnect) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](#contributing)

nConnect allows you to securely connect to remote machines without the need of
any server, public IP address, or publicly exposed ports. It features end to end
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
  or port forwarding. NKN tunnel only establish outbound (websocket)
  connections, so Internet access is all they need on both side.

- Top level security: All data are end to end authenticated and encrypted. No
  one else in the world except sender and receiver can see or modify the content
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

But most of the times you might want to start nConnect server with a few useful
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
  your configuration. If you want best security, disable the admin dashboard
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
will allow any address using public key
`ad37e248005113dd42be15a4885e6446e9e23f35537dfa6c584f2563a7e8f96d`,
such as
`ad37e248005113dd42be15a4885e6446e9e23f35537dfa6c584f2563a7e8f96d`
and
`nkn.ad37e248005113dd42be15a4885e6446e9e23f35537dfa6c584f2563a7e8f96d`.
If you want to add a specific address to the list, you should add `^` in the front, and `$` to the end. For example,
`^nkn.ad37e248005113dd42be15a4885e6446e9e23f35537dfa6c584f2563a7e8f96d$`
will only match address
`nkn.ad37e248005113dd42be15a4885e6446e9e23f35537dfa6c584f2563a7e8f96d`
but not
`ad37e248005113dd42be15a4885e6446e9e23f35537dfa6c584f2563a7e8f96d`
or
`nConnect.ad37e248005113dd42be15a4885e6446e9e23f35537dfa6c584f2563a7e8f96d`.

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

By default all local IP addresses on the server machine will be added to routes,
but you can manually specify which IP or IP range you would like to route
through the VPN using `--vpn-route` arguments. Use `./nConnect -h` for all available arguments.

Note that if you are using windows, you would need to install the network
adaptor driver and change adaptor info beforehand. The simplest way of doing
that is to install nConnect client for windows before using nConnect command
line version.

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

You can also change the name, IP, etc of the TUN device. Use `./nConnect -h` for
all available arguments.

Note that if you are using windows, you would need to install the network
adaptor driver and change adaptor info beforehand. The simplest way of doing
that is to install nConnect client for windows before using nConnect command
line version.

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
