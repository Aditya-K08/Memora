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
	n1 := raft.NewRaftNode(1)
	n2 := raft.NewRaftNode(2)
	n3 := raft.NewRaftNode(3)

	n1.AddPeer(n2); n1.AddPeer(n3)
	n2.AddPeer(n1); n2.AddPeer(n3)
	n3.AddPeer(n1); n3.AddPeer(n2)

	go n1.RunElectionTimer()
	go n2.RunElectionTimer()
	go n3.RunElectionTimer()

	store, _ := storage.Open("data.wal")
	defer store.Close()

	fmt.Println("KV Store started")

	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		sc.Scan()
		parts := strings.SplitN(sc.Text(), " ", 3)

		switch strings.ToUpper(parts[0]) {

		case "PUT":
			leader := findLeader(n1, n2, n3)

			if leader == nil {
				fmt.Println("ERR: no leader")
				continue
			}

			ok := leader.Replicate(parts[1], []byte(parts[2]))
			if !ok {
				fmt.Println("ERR: replication failed")
				continue
			}

			store.Put(parts[1], []byte(parts[2]), 0)
			fmt.Println("OK")

		case "CRASH":
			fmt.Println("💥 crashing leader node1")
			n1.Stop()

		case "GET":
			v, ok := store.Get(parts[1])
			if !ok { fmt.Println("(nil)") } else {
				fmt.Println(string(v))
			}

		case "EXIT":
			time.Sleep(100 * time.Millisecond)
			return
		}
	}
}
