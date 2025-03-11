package api

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/M-kos/anti_brutforce/internal/api/pb" //nolint
	"github.com/M-kos/anti_brutforce/internal/config"
	"google.golang.org/grpc"
)

type ServerApi struct {
	pb.UnimplementedAntiBruteforceServer
}

func Run(conf *config.Config) error {
	server := grpc.NewServer()
	pb.RegisterAntiBruteforceServer(server, &ServerApi{})

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GRPCPopt))
	if err != nil {
		return err
	}

	log.Printf("starting server on %s", listen.Addr()) // TODO: add logger

	if err = server.Serve(listen); err != nil {
		return err
	}

	return nil
}

func (s *ServerApi) CheckCredentials(ctx context.Context, req *pb.CheckCredentialsRequest) (*pb.OkResponse, error) {
	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) ResetBuckets(ctx context.Context, req *pb.ResetBucketsRequest) (*pb.OkResponse, error) {
	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) AddToWhitelist(ctx context.Context, req *pb.WhitelistRequest) (*pb.OkResponse, error) {
	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) RemoveFromWhitelist(ctx context.Context, req *pb.WhitelistRequest) (*pb.OkResponse, error) {
	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) AddToBlacklist(ctx context.Context, req *pb.BlacklistRequest) (*pb.OkResponse, error) {
	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) RemoveFromBlacklist(ctx context.Context, req *pb.BlacklistRequest) (*pb.OkResponse, error) {
	return &pb.OkResponse{Ok: true}, nil
}
