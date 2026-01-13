package raft

import (
	"context"
	"time"

	pb "memora/internal/raftpb"
	"google.golang.org/grpc"
)

func dial(addr string) (pb.RaftClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewRaftClient(conn), nil
}

func (r *RaftNode) sendVote(addr string) bool {

	client, err := dial(addr)
	if err != nil {
		return false // peer unreachable
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resp, err := client.RequestVote(ctx,
		&pb.RequestVoteArgs{
			Term:        r.currentTerm,
			CandidateId: int32(r.id),
		})

	if err != nil {
		return false
	}

	return resp.VoteGranted
}

func (r *RaftNode) sendHeartbeat(addr string) {

	client, err := dial(addr)
	if err != nil {
		return // peer down
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	client.AppendEntries(ctx, &pb.AppendEntriesArgs{
		Term:     r.currentTerm,
		LeaderId: int32(r.id),
	})
}
