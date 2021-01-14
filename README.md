# nConnect

[![GitHub license](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/nknorg/nconnect)](https://goreportcard.com/report/github.com/nknorg/nconnect) [![Build Status](https://travis-ci.org/nknorg/nconnect.svg?branch=master)](https://travis-ci.org/nknorg/nconnect) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](#contributing)

nConnect allows you to access any TCP-based applications running anywhere
through a SOCKS proxy. nConnect uses
[nkn-tunnel](https://github.com/nknorg/nkn-tunnel) for end to end tunneling,
thus benefits from all the advantages of
[nkn-tunnel](https://github.com/nknorg/nkn-tunnel):

- Network agnostic: Neither sender nor receiver needs to have public IP address
  or port forwarding. NKN tunnel only establish outbound (websocket)
  connections, so Internet access is all they need on both side.

- Top level security: All data are end to end authenticated and encrypted. No
  one else in the world except sender and receiver can see or modify the content
  of the data. The same public key is used for both routing and encryption,
  eliminating the possibility of man in the middle attack.

- Decent performance: By aggregating multiple overlay paths concurrently, one
  can get ~100ms end to end latency and 10+mbps end to end throughput between
  international devices using the default NKN client mode, or much lower latency
  and higher throughput using Tuna mode.

- Everything is open source and decentralized. The default NKN client mode is
  free (If you are curious, node relay traffic for clients for free to earn
  mining rewards in NKN blockchain), while Tuna mode requires listener to pay
  NKN token directly to Tuna service providers.

A diagram of how nConnect connects applications:

```
Application <--> Socks proxy client <--> NKN Tunnel <--> Socks proxy server <--> Application
```

## Build

```shell
make
```

## Usage

### Quick Start

The easiest way to use nConnect is using various pre-built nConnect clients. You
just need to deploy nConnect in server mode to the device you want to connect
to using the following argument:

```shell
nConnect -s --tuna --admin-http :8001
```

Then you can visit `http://<device-local-ip-address>:8001` in browser to access
the admin dashboard and follow the guide there.

### Advanced Usage

You can use nConnect in client mode if the pre-built nConnect clients do not
work for you. When started in client mode, nConnect will create a local SOCKS
proxy. Any connection made through that proxy will be routed to nConnect server
side first and use the server side as exit. To see available arguments, run
`nConnect -h`.

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
