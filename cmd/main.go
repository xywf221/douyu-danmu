package main

import (
	"douyuDm/net"
)

func main() {
	ws, err := net.NewDYWebSocket(288016)
	if err != nil {
		panic(err)
	}
	ws.Run()
}
