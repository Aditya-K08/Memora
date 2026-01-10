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

	id      int
	state   State
	stopped bool

	log         []LogEntry
	commitIndex int

	currentTerm uint64
	votedFor    int

	peers map[int]*RaftNode

	lastHeartbeat     time.Time
	electionTimeout   time.Duration
	heartbeatInterval time.Duration
}

func NewRaftNode(id int) *RaftNode {
	return &RaftNode{
		id:                id,
		state:             Follower,
		votedFor:          -1,
		log:               []LogEntry{},
		peers:             make(map[int]*RaftNode),
		lastHeartbeat:     time.Now(),
		electionTimeout:   300 * time.Millisecond,
		heartbeatInterval: 100 * time.Millisecond,
	}
}

func (r *RaftNode) AddPeer(p *RaftNode) {
	r.peers[p.id] = p
}

func (r *RaftNode) Stop() {
	r.mu.Lock()
	r.stopped = true
	r.mu.Unlock()
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
		return RequestVoteReply{Term: r.currentTerm}
	}

	r.currentTerm = args.Term
	r.votedFor = args.CandidateID
	r.lastHeartbeat = time.Now()

	return RequestVoteReply{Term: r.currentTerm, VoteGranted: true}
}

func (r *RaftNode) AppendEntries(args AppendEntriesArgs) AppendEntriesReply {
	r.mu.Lock()
	defer r.mu.Unlock()

	if args.Term < r.currentTerm {
		return AppendEntriesReply{Term: r.currentTerm}
	}

	r.currentTerm = args.Term
	r.state = Follower
	r.lastHeartbeat = time.Now()

	if len(args.Entries) > 0 {
		r.log = append(r.log, args.Entries...)
	}

	return AppendEntriesReply{Term: r.currentTerm, Success: true}
}

func (r *RaftNode) Replicate(key string, val []byte) bool {
	r.mu.Lock()
	entry := LogEntry{Term: r.currentTerm, Key: key, Value: val}
	r.log = append(r.log, entry)
	r.mu.Unlock()

	acks := 1
	majority := (len(r.peers)+1)/2 + 1

	for _, p := range r.peers {
		if p.AppendEntries(AppendEntriesArgs{
			Term:     r.currentTerm,
			LeaderID: r.id,
			Entries:  []LogEntry{entry},
		}).Success {
			acks++
		}
	}

	return acks >= majority
}
