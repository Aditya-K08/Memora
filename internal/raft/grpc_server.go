package raft

import (
	"context"
	"net"

	pb "memora/internal/raftpb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedRaftServer
	node *RaftNode
}

func (s *Server) RequestVote(ctx context.Context, a *pb.RequestVoteArgs)
(*pb.RequestVoteReply, error) {

	r := s.node.RequestVote(RequestVoteArgs{
		Term: a.Term,
		CandidateID: int(a.CandidateId),
	})

	return &pb.RequestVoteReply{
		Term: r.Term,
		VoteGranted: r.VoteGranted,
	}, nil
}

func (s *Server) AppendEntries(ctx context.Context, a *pb.AppendEntriesArgs)
(*pb.AppendEntriesReply, error) {

	r := s.node.AppendEntries(AppendEntriesArgs{
		Term: a.Term,
		LeaderID: int(a.LeaderId),
	})

	return &pb.AppendEntriesReply{
		Term: r.Term,
		Success: r.Success,
	}, nil
}

func StartServer(port string, node *RaftNode) {
	l, _ := net.Listen("tcp", port)
	s := grpc.NewServer()
	pb.RegisterRaftServer(s, &Server{node: node})
	s.Serve(l)
}
