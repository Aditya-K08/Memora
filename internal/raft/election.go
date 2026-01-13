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

		// Leader does not start election
		if r.state == Leader {
			r.mu.Unlock()
			continue
		}

		// Start election if no heartbeat
		if time.Since(r.lastHeartbeat) >= timeout {
			r.startElection()
		}

		r.mu.Unlock()
	}
}

func (r *RaftNode) startElection() {
	fmt.Println("Node", r.id, "starting election")

	r.state = Candidate
	r.currentTerm++
	r.votedFor = r.id

	votes := 1
	majority := (len(r.peerAddrs)+1)/2 + 1

	for _, addr := range r.peerAddrs {
		if r.sendVote(addr) {
			votes++
		}
	}

	fmt.Println("Node", r.id, "got", votes, "votes")

	if votes >= majority {
		r.state = Leader
		fmt.Println("LEADER:", r.id)
		go r.runHeartbeatLoop()
	}
}


func (r *RaftNode) runHeartbeatLoop() {
	ticker := time.NewTicker(r.heartbeatInterval)
	defer ticker.Stop()

	for range ticker.C {
		r.mu.Lock()
		if r.state != Leader {
			r.mu.Unlock()
			return
		}
		r.mu.Unlock()

		for _, addr := range r.peerAddrs {
			r.sendHeartbeat(addr)
		}
	}
}
