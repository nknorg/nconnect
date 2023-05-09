package tests

var port int = 1080
var proxyAddr string = "127.0.0.1:1080"

const (
	rounds = 10

	seedHex = "e68e046d13dd911594576ba0f4a196e9666790dc492071ad9ea5972c0b940435"

	tcpServerAddr  = "127.0.0.1:54321"
	httpServerAddr = "127.0.0.1:54322"
	udpServerAddr  = "127.0.0.1:54323"
	httpServiceUrl = "http://" + httpServerAddr + "/httpEcho"

	tunaNodeStarted = "tuna node is started"
)
