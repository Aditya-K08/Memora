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
	raftNode := raft.NewRaftNode(1, []int{2, 3})
	go raftNode.RunElectionTimer()

	store, err := storage.Open("data.wal")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	fmt.Println("KV Store started")
	fmt.Println("Commands: PUT key value | GET key | DEL key | EXIT")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		parts := strings.SplitN(scanner.Text(), " ", 3)

		switch strings.ToUpper(parts[0]) {
		case "PUT":
			if len(parts) != 3 {
				fmt.Println("Usage: PUT key value")
				continue
			}
			store.Put(parts[1], []byte(parts[2]), 0)
			fmt.Println("OK")

		case "GET":
			val, ok := store.Get(parts[1])
			if !ok {
				fmt.Println("(nil)")
			} else {
				fmt.Println(string(val))
			}

		case "DEL":
			store.Delete(parts[1])
			fmt.Println("OK")

		case "EXIT":
			fmt.Println("Shutting down")
			time.Sleep(100 * time.Millisecond)
			return
		}
	}
}
