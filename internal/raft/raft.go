package raft

import (
	"sync"
	"time"
)
type LogEntry struct {
	Term  uint64
	Key   string
	Value []byte
}

type RaftNode struct {
	mu sync.Mutex

	id    int
	state State

	currentTerm uint64
	votedFor    int

	peerAddrs []string // gRPC peers
	log         []LogEntry
	commitIndex int
	lastApplied int

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
		log: make([]LogEntry, 0),
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

	if len(args.Entries) > 0 {
		r.log = append(r.log, args.Entries...)
	}

	return AppendEntriesReply{
		Term:    r.currentTerm,
		Success: true,
	}
}
func (r *RaftNode) Replicate(key string, val []byte) bool {

	r.mu.Lock()
	entry := LogEntry{
		Term:  r.currentTerm,
		Key:   key,
		Value: val,
	}
	r.log = append(r.log, entry)
	index := len(r.log) - 1
	r.mu.Unlock()

	acks := 1
	majority := (len(r.peerAddrs)+1)/2 + 1

	for _, addr := range r.peerAddrs {
		client, err := dial(addr)
		if err != nil {
			continue
		}

		resp, err := client.AppendEntries(
			context.Background(),
			&pb.AppendEntriesArgs{
				Term:     r.currentTerm,
				LeaderId: int32(r.id),
				Entries:  []*pb.LogEntry{
					{
						Term:  entry.Term,
						Key:   entry.Key,
						Value: entry.Value,
					},
				},
			},
		)

		if err == nil && resp.Success {
			acks++
		}
	}

	if acks >= majority {
		r.mu.Lock()
		r.commitIndex = index
		r.mu.Unlock()
		return true
	}

	return false
}

