package raft

import (
	"context"
	"time"

	pb "memora/internal/raftpb"
	"google.golang.org/grpc"
)

func dial(addr string) pb.RaftClient {
	c, _ := grpc.Dial(addr, grpc.WithInsecure())
	return pb.NewRaftClient(c)
}

func (r *RaftNode) sendVote(addr string) bool {
	c := dial(addr)

	resp, _ := c.RequestVote(context.Background(),
		&pb.RequestVoteArgs{
			Term: r.currentTerm,
			CandidateId: int32(r.id),
		})

	return resp.VoteGranted
}

func (r *RaftNode) sendHeartbeat(addr string) {
	c := dial(addr)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	c.AppendEntries(ctx, &pb.AppendEntriesArgs{
		Term: r.currentTerm,
		LeaderId: int32(r.id),
	})
}
