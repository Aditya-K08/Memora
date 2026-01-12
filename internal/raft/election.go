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
	for _, addr := range r.peerAddrs {
		if r.sendVote(addr) {
			votes++
		}
	}

	if votes >= (len(r.peerAddrs)+1)/2+1 {
		r.state = Leader
		fmt.Println("LEADER:", r.id)
		go r.heartbeatLoop()
	}
}


func (r *RaftNode) runHeartbeatLoop() {
	t := time.NewTicker(100 * time.Millisecond)
	for range t.C {
		if r.state != Leader { return }
		for _, a := range r.peerAddrs {
			r.sendHeartbeat(a)
		}
	}
}
