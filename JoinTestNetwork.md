# Join nConnect Test Nework

* The nConnect manager address: 
`manager.7cafe0ae02789f8eb6b293e46b0ac5cf8f92f73042199c8161e5b5f90b13dcb5`
* The nConnect manager web url: 
`http://147.182.218.42:8000/network`
* The test nework IP scope: 
`10.0.86.2 ~ 10.0.86.254`

Please follow these steps to join or manage the test network.

> On Unix-like systems you may need to preface commands with sudo, while on Windows you will need to use an administrator-mode command prompt. 

## 1. Clone nConnect Repository:
```
git clone -b network https://github.com/billfort/nconnect
```

## 2. Compile nConnect

* Linux or MacOs:
```
make
```

* Windows:
```
go build -o nConnect.exe bin/main.go
```

## 3. Start nConnect Manager [optional]
The manager node usually is running. If nConnect network manager is not available, you can log into `147.182.218.42` to start it

```
./nConnect -m -f config.manager.json
```

## 4. Conifg Network Member Node

```
cp config.network.json config.member.json
```

```
vi config.member.json
```

Edit `identifier`, `managerAddress`, `nodeName` and `seed`.

```
{
    "identifier": "alice",  // your nkn client identifier, such alice, bruce, bill, max, ...
    "managerAddress": "manager.7cafe0ae02789f8eb6b293e46b0ac5cf8f92f73042199c8161e5b5f90b13dcb5", // This is testing manager address.
    "nodeName": "alice",    // It can be as same as identifier, or different. it is used only in network node naming.
    "seed": "",		        // If you want other nodes can access your node, you need start nConnect server, make sure your wallet have NKN balance.
}
```

## 5. Start nConnect Network Member Node

* Linux or Mac OS in root

```
./nConnect -n -s -c -f config.member.json --tuna --vpn --udp
```

* Windows Powershell Administrator

```
./nConnect.exe -n -s -c -f config.member.json --tuna --vpn --udp
```

## 6. Wait for nConnect Manager to Authorize Your Joining. 

Open the network manager web page `http://147.182.218.42:8000/network`, and refresh `Wait for Authorization` section, should see the new member node. Click "Accept" to authorize the new member. You can check `Accept All Member` to set your node is accessible for all members.

After the manager authorize your node joining, you should see a console message printed, such as:

```
Congratulations!!! nConnect network member authorized, IP: 10.0.86.xxx, mask: 255.255.255.0
```

The `ip 10.0.86.xxx` is your IP in this private network.

## 7. Test newtork access:

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
