package api

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/M-kos/anti_brutforce/internal/api/pb" //nolint
	"github.com/M-kos/anti_brutforce/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LimitterService interface {
	Check(ctx context.Context, login string, ip string, password string) error
	Remove(ctx context.Context, login string, ip string) error
}
type LabelService interface {
	Add(ctx context.Context, subnet string) error
	Remove(ctx context.Context, subnet string) error
}

type ServerApi struct {
	pb.UnimplementedAntiBruteforceServer

	limitterService   LimitterService
	whitelabelService LabelService
	blacklistService  LabelService
}

func NewServerApi(limitterService LimitterService, whitelabelService LabelService, blacklistService LabelService) *ServerApi {
	return &ServerApi{
		limitterService:   limitterService,
		whitelabelService: whitelabelService,
		blacklistService:  blacklistService,
	}
}

func Run(conf *config.Config, limitterService LimitterService, whitelabelService LabelService, blacklistService LabelService) error {
	server := grpc.NewServer()
	api := NewServerApi(limitterService, whitelabelService, blacklistService)
	pb.RegisterAntiBruteforceServer(server, api)

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
	login := req.GetLogin()
	password := req.GetPassword()
	ip := req.GetIp()

	if login == "" {
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, "login is required")
	}
	if password == "" {
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, "password is required")
	}
	if net.ParseIP(ip) == nil {
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, "ip is required")
	}

	if err := s.limitterService.Check(ctx, login, ip, password); err != nil {
		return &pb.OkResponse{Ok: false}, status.Error(codes.Internal, err.Error())
	}

	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) ResetBuckets(ctx context.Context, req *pb.ResetBucketsRequest) (*pb.OkResponse, error) {
	login := req.GetLogin()
	ip := req.GetIp()

	if login == "" {
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, "login is required")
	}
	if net.ParseIP(ip) == nil {
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, "ip is required")
	}

	if err := s.limitterService.Remove(ctx, login, ip); err != nil {
		return &pb.OkResponse{Ok: false}, status.Error(codes.Internal, err.Error())
	}

	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) AddToWhitelist(ctx context.Context, req *pb.WhitelistRequest) (*pb.OkResponse, error) {
	subnet := req.GetSubnet()

	if subnet == "" {
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, "subnet is required")
	}

	if err := s.whitelabelService.Add(ctx, subnet); err != nil {
		return &pb.OkResponse{Ok: false}, status.Error(codes.Internal, err.Error())
	}

	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) RemoveFromWhitelist(ctx context.Context, req *pb.WhitelistRequest) (*pb.OkResponse, error) {
	subnet := req.GetSubnet()

	if subnet == "" {
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, "subnet is required")
	}

	if err := s.whitelabelService.Remove(ctx, subnet); err != nil {
		return &pb.OkResponse{Ok: false}, status.Error(codes.Internal, err.Error())
	}

	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) AddToBlacklist(ctx context.Context, req *pb.BlacklistRequest) (*pb.OkResponse, error) {
	subnet := req.GetSubnet()

	if subnet == "" {
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, "subnet is required")
	}

	if err := s.blacklistService.Add(ctx, subnet); err != nil {
		return &pb.OkResponse{Ok: false}, status.Error(codes.Internal, err.Error())
	}

	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) RemoveFromBlacklist(ctx context.Context, req *pb.BlacklistRequest) (*pb.OkResponse, error) {
	subnet := req.GetSubnet()

	if subnet == "" {
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, "subnet is required")
	}

	if err := s.blacklistService.Remove(ctx, subnet); err != nil {
		return &pb.OkResponse{Ok: false}, status.Error(codes.Internal, err.Error())
	}

	return &pb.OkResponse{Ok: true}, nil
}
