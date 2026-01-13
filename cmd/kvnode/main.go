package main

import (
	"fmt"
	"os"
	"strconv"

	"memora/internal/raft"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <id> <port>")
		return
	}

	id, _ := strconv.Atoi(os.Args[1])
	port := os.Args[2]

	all := map[int]string{
		1: "localhost:5001",
		2: "localhost:5002",
		3: "localhost:5003",
	}

	peers := []string{}
	for k, v := range all {
		if k != id {
			peers = append(peers, v)
		}
	}

	node := raft.NewRaftNode(id, peers)

	go raft.StartGRPCServer(":"+port, node)
	go node.RunElectionTimer()

	select {}
}
