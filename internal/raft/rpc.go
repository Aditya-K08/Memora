package raft

type RequestVoteArgs struct{
	Term uint64
	CandidateID int
}

type RequestVoteReply struct{
	Term uint64
	VoteGranted bool
}
type AppendEntriesArgs struct {
	Term     uint64
	LeaderID int
}

type AppendEntriesReply struct {
	Term    uint64
	Success bool
}
