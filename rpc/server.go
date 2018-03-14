package rpc

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type Listener int

func (l *Listener) GetLine(line []byte, ack *bool) error {
	fmt.Println(string(line))
	return nil
}

func (l *Listener) Add(a *SumArgs, res *int) error {
	*res = a.A + a.B
	return nil
}

func StartServer() {
	addy, err := net.ResolveTCPAddr("tcp", "localhost:4444")

	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		log.Fatal(err)
	}

	listener := new(Listener)
	rpc.Register(listener)
	rpc.Accept(inbound)
}
