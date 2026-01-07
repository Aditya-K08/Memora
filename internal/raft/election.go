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
	total := len(r.peers) + 1
	majority := total/2 + 1

	for _, p := range r.peers {
		reply := p.RequestVote(RequestVoteArgs{
			Term:        r.currentTerm,
			CandidateID: r.id,
		})
		if reply.VoteGranted {
			votes++
		}
	}

	if votes >= majority {
		r.state = Leader
		r.lastHeartbeat = time.Now()
		fmt.Printf("[RAFT] Leader %d elected (term %d, votes %d)\n",
			r.id, r.currentTerm, votes)
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
		term := r.currentTerm
		r.mu.Unlock()

		for _, peer := range r.peers {
			peer.AppendEntries(AppendEntriesArgs{
				Term:     term,
				LeaderID: r.id,
			})
		}
	}
}
