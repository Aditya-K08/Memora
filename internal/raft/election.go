package raft

import (
	"fmt"
	"math/rand"
	"time"
)

func (r *RaftNode) RunElectionTimer() {
	for {
		timeout := time.Duration(150+rand.Intn(150)) * time.Millisecond
		time.Sleep(timeout)

		r.mu.Lock()
		if r.state == Leader {
			r.mu.Unlock()
			continue
		}

		if time.Since(r.lastHeartbeat) >= timeout {
			r.startElection()
		}
		r.mu.Unlock()
	}
}

func (r *RaftNode) startElection() {
	r.state = Candidate
	r.currentTerm++
	r.votedFor = r.id
	votes := 1
	fmt.Println("Node became leader for term", r.currentTerm)

	if votes > len(r.peers)/2 {
		r.state = Leader
		r.lastHeartbeat = time.Now()
	}
}
