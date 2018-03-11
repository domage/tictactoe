package sockets

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func SocketClient(ip string, port int) error {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
		return err
	}

	conn.Write([]byte(InitCommand))
	log.Printf("Send: %s", InitCommand)
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	log.Printf("Received:\n%s", buff[:n])
	player := string(buff[:n])

	log.Printf("You play as %s", player)
	var s string
	for {
		fmt.Scanf("%s", &s)
		_, err = conn.Write([]byte(s + player))
		if err != nil {
			log.Fatalln(err)
			return err
		}

		log.Printf("Send: %s", s+player)
		if s == StopCommand {
			return nil
		}

		buff = make([]byte, 1024)
		n, err = conn.Read(buff)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		log.Printf("Received:\n%s", buff[:n])
		if strings.HasSuffix(string(buff[:n]), StopCommand) {
			return nil
		}
	}

}

func StartClient() error {

	var (
		ip   = "127.0.0.1"
		port = 3333
	)

	return SocketClient(ip, port)
}
