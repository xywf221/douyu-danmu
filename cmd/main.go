package main

import (
	"douyuDm/net"
)

func main() {
	ws, err := net.NewDYWebSocket(475252)
	if err != nil {
		panic(err)
	}
	ws.Run()
}
