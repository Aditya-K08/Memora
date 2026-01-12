package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"memora/internal/raft"
	"memora/internal/storage"
)

func findLeader(nodes ...*raft.RaftNode) *raft.RaftNode {
    for _, n := range nodes {
        if n.IsLeader() {
            return n
        }
    }
    return nil
}


func main() {
		id := os.Args[1]
		port := os.Args[2]
	
		node := raft.NewRaftNode(id)
		node.peerAddrs = []string{
			"localhost:5001",
			"localhost:5002",
			"localhost:5003",
		}
	
		go raft.StartServer(":"+port, node)
		go node.RunElectionTimer()
	
		select {}
	
	
}
