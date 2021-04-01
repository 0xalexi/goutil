package goutil

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

var muxMu sync.Mutex

func RunRPCServer(handler *rpc.Server, host string, port int) {
	muxMu.Lock()
	defer muxMu.Unlock()
	address := fmt.Sprintf("%s:%d", host, port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("listen(%q): %s\n", address, err)
		return
	}
	fmt.Printf("Server listening on %s\n", ln.Addr())
	go func() {
		for {
			cxn, err := ln.Accept()
			if err != nil {
				log.Printf("listen(%q): %s\n", address, err)
				return
			}
			log.Printf("Server accepted connection to %s from %s\n", cxn.LocalAddr(), cxn.RemoteAddr())
			go handler.ServeConn(cxn)
		}
	}()
}
