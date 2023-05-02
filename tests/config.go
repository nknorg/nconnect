package tests

const (
	rounds = 10

	seedHex = "e68e046d13dd911594576ba0f4a196e9666790dc492071ad9ea5972c0b940435"

	proxyAddr  = "127.0.0.1:1080"
	socksProxy = "socks5://" + proxyAddr

	serverAddr     = "127.0.0.1:54321"
	httpServiceUrl = "http://" + serverAddr + "/httpEcho"

	tunaNodeStarted = "tuna node is started"

	webServerIsReady = "web server is ready"
	webServerExited  = "web server exited"
	webClientExited  = "web client exited"

	tcpServerIsReady = "tcp server is ready"
	tcpServerExited  = "tcp server exited"
	tcpClientExited  = "tcp client exited"

	udpServerIsReady = "udp server is ready"
	udpServerExited  = "udp server exited"
	udpClientExited  = "udp client exited"

	exited = "exited"
)
