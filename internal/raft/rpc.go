package raft

type RequestVoteArgs struct{
	Term uint64
	CandidateId int
}

type RequestVoteReply struct{
	Term uint64
	VoteGranted bool
}