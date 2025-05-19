package httpn

import (
	"fmt"
	"net"
)

func SetupTcp(port int) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	return listener
}
