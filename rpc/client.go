package rpc

import (
	"fmt"
	"log"
	"net/rpc"
)

type SumArgs struct {
	A, B int
}

func StartClient() {
	client, err := rpc.Dial("tcp", "localhost:4444")

	if err != nil {
		log.Fatal(err)
	}

	//in := bufio.NewReader(os.Stdin)

	//for {
	//line, _, err := in.ReadLine()
	//if err != nil {
	//	log.Fatal(err)
	//}
	var reply int
	args := &SumArgs{10, 20}
	err = client.Call("Listener.Add", args, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(reply)
	//}
}
