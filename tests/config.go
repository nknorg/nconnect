package tests

var port int = 1080
var proxyAddr string = "127.0.0.1:1080"

const (
	numMsgs = 10

	seedHex = "e68e046d13dd911594576ba0f4a196e9666790dc492071ad9ea5972c0b940435"

	tcpPort  = ":54321"
	httpPort = ":54322"
	udpPort  = ":54323"

	tunaNodeStarted = "tuna node is started"
)

var servers = []string{"127.0.0.1"} // {"10.10.0.15", "10.136.0.10"}
