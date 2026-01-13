package raft

import (
	"sync"
	"time"
)

type RaftNode struct {
	mu sync.Mutex

	id    int
	state State

	currentTerm uint64
	votedFor    int

	peerAddrs []string // gRPC peers

	lastHeartbeat     time.Time
	electionTimeout   time.Duration
	heartbeatInterval time.Duration
}

func NewRaftNode(id int, peers []string) *RaftNode {
	return &RaftNode{
		id:                id,
		state:             Follower,
		votedFor:          -1,
		peerAddrs:         peers,
		lastHeartbeat:     time.Now(),
		electionTimeout:   300 * time.Millisecond,
		heartbeatInterval: 100 * time.Millisecond,
	}
}

func (r *RaftNode) IsLeader() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state == Leader
}

func (r *RaftNode) RequestVote(args RequestVoteArgs) RequestVoteReply {
	r.mu.Lock()
	defer r.mu.Unlock()

	if args.Term < r.currentTerm {
		return RequestVoteReply{
			Term:        r.currentTerm,
			VoteGranted: false,
		}
	}

	if args.Term > r.currentTerm {
		r.currentTerm = args.Term
		r.votedFor = -1
		r.state = Follower
	}

	if r.votedFor == -1 || r.votedFor == args.CandidateID {
		r.votedFor = args.CandidateID
		r.lastHeartbeat = time.Now()

		return RequestVoteReply{
			Term:        r.currentTerm,
			VoteGranted: true,
		}
	}

	return RequestVoteReply{
		Term:        r.currentTerm,
		VoteGranted: false,
	}
}

func (r *RaftNode) AppendEntries(args AppendEntriesArgs) AppendEntriesReply {
	r.mu.Lock()
	defer r.mu.Unlock()

	if args.Term < r.currentTerm {
		return AppendEntriesReply{
			Term:    r.currentTerm,
			Success: false,
		}
	}

	r.currentTerm = args.Term
	r.state = Follower
	r.lastHeartbeat = time.Now()

	return AppendEntriesReply{
		Term:    r.currentTerm,
		Success: true,
	}
}
