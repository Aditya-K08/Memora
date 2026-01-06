package raft

import (
	"sync"
	"time"
)

type RaftNode struct {
	mu sync.Mutex

	id    int
	peers []int

	state State

	currentTerm uint64
	votedFor    int

	electionTimeout time.Duration
	lastHeartbeat   time.Time
}

func NewRaftNode(id int, peers []int) *RaftNode {
	return &RaftNode{
		id:              id,
		peers:           peers,
		state:           Follower,
		votedFor:        -1,
		lastHeartbeat:   time.Now(),
		electionTimeout: 300 * time.Millisecond,
	}
}
