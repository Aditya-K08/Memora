package raft

import (
	"context"
	"net"

	pb "memora/internal/raftpb"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedRaftServer
	node *RaftNode
}

func (s *GRPCServer) RequestVote(
	ctx context.Context,
	req *pb.RequestVoteArgs,
) (*pb.RequestVoteReply, error) {

	reply := s.node.RequestVote(RequestVoteArgs{
		Term:        req.Term,
		CandidateID: int(req.CandidateId),
	})

	return &pb.RequestVoteReply{
		Term:        reply.Term,
		VoteGranted: reply.VoteGranted,
	}, nil
}

func (s *GRPCServer) AppendEntries(
	ctx context.Context,
	req *pb.AppendEntriesArgs,
) (*pb.AppendEntriesReply, error) {

	reply := s.node.AppendEntries(AppendEntriesArgs{
		Term:     req.Term,
		LeaderID: int(req.LeaderId),
	})

	return &pb.AppendEntriesReply{
		Term:    reply.Term,
		Success: reply.Success,
	}, nil
}

func StartGRPCServer(port string, node *RaftNode) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterRaftServer(server, &GRPCServer{node: node})

	return server.Serve(lis)
}
