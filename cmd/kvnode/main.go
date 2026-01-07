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

func main() {
	n1 := raft.NewRaftNode(1)
	n2 := raft.NewRaftNode(2)
	n3 := raft.NewRaftNode(3)

	n1.AddPeer(n2)
	n1.AddPeer(n3)
	n2.AddPeer(n1)
	n2.AddPeer(n3)
	n3.AddPeer(n1)
	n3.AddPeer(n2)

	go n1.RunElectionTimer()
	go n2.RunElectionTimer()
	go n3.RunElectionTimer()

	store, _ := storage.Open("data.wal")
	defer store.Close()

	fmt.Println("KV Store started")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		parts := strings.SplitN(scanner.Text(), " ", 3)

		switch strings.ToUpper(parts[0]) {
		case "PUT":
			if !n1.IsLeader() {
				fmt.Println("ERR: not leader")
				continue
			}
			store.Put(parts[1], []byte(parts[2]), 0)
			fmt.Println("OK")

		case "GET":
			v, ok := store.Get(parts[1])
			if !ok {
				fmt.Println("(nil)")
			} else {
				fmt.Println(string(v))
			}

		case "DEL":
			store.Delete(parts[1])
			fmt.Println("OK")

		case "EXIT":
			time.Sleep(100 * time.Millisecond)
			return
		}
	}
}
