package api

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/M-kos/anti_brutforce/internal/api/pb" //nolint
	"github.com/M-kos/anti_brutforce/internal/config"
	"github.com/M-kos/anti_brutforce/internal/lib/logger"
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

	logger logger.LoggerProvider
}

func NewServerApi(limitterService LimitterService, whitelabelService LabelService, blacklistService LabelService, logger logger.LoggerProvider) *ServerApi {
	return &ServerApi{
		limitterService:   limitterService,
		whitelabelService: whitelabelService,
		blacklistService:  blacklistService,
		logger:            logger,
	}
}

func Run(conf *config.Config, limitterService LimitterService, whitelabelService LabelService, blacklistService LabelService, logger logger.LoggerProvider) error {
	server := grpc.NewServer()
	api := NewServerApi(limitterService, whitelabelService, blacklistService, logger)
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
		s.logger.Error(ErrLoginIsRequired.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, ErrLoginIsRequired.Error())
	}
	if password == "" {
		s.logger.Error(ErrPasswordIsRequired.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, ErrPasswordIsRequired.Error())
	}
	if net.ParseIP(ip) == nil {
		s.logger.Error(ErrIpIsRequired.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, ErrIpIsRequired.Error())
	}

	if err := s.limitterService.Check(ctx, login, ip, password); err != nil {
		s.logger.Error(err.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.Internal, ErrCheck.Error())
	}

	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) ResetBuckets(ctx context.Context, req *pb.ResetBucketsRequest) (*pb.OkResponse, error) {
	login := req.GetLogin()
	ip := req.GetIp()

	if login == "" {
		s.logger.Error(ErrLoginIsRequired.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, ErrLoginIsRequired.Error())
	}
	if net.ParseIP(ip) == nil {
		s.logger.Error(ErrLoginIsRequired.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, ErrIpIsRequired.Error())
	}

	if err := s.limitterService.Remove(ctx, login, ip); err != nil {
		s.logger.Error(err.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.Internal, ErrRemove.Error())
	}

	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) AddToWhitelist(ctx context.Context, req *pb.WhitelistRequest) (*pb.OkResponse, error) {
	cidr := req.GetCidr()

	if cidr == "" {
		s.logger.Error(ErrCidrRequired.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, ErrCidrRequired.Error())
	}

	if err := s.whitelabelService.Add(ctx, cidr); err != nil {
		s.logger.Error(err.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.Internal, ErrAddToWhitelist.Error())
	}

	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) RemoveFromWhitelist(ctx context.Context, req *pb.WhitelistRequest) (*pb.OkResponse, error) {
	cidr := req.GetCidr()

	if cidr == "" {
		s.logger.Error(ErrCidrRequired.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, ErrCidrRequired.Error())
	}

	if err := s.whitelabelService.Remove(ctx, cidr); err != nil {
		s.logger.Error(err.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.Internal, ErrRemoveFromWhitelist.Error())
	}

	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) AddToBlacklist(ctx context.Context, req *pb.BlacklistRequest) (*pb.OkResponse, error) {
	cidr := req.GetCidr()

	if cidr == "" {
		s.logger.Error(ErrCidrRequired.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, ErrCidrRequired.Error())
	}

	if err := s.blacklistService.Add(ctx, cidr); err != nil {
		s.logger.Error(err.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.Internal, ErrAddToBlacklist.Error())
	}

	return &pb.OkResponse{Ok: true}, nil
}

func (s *ServerApi) RemoveFromBlacklist(ctx context.Context, req *pb.BlacklistRequest) (*pb.OkResponse, error) {
	cidr := req.GetCidr()

	if cidr == "" {
		s.logger.Error(ErrCidrRequired.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.InvalidArgument, ErrCidrRequired.Error())
	}

	if err := s.blacklistService.Remove(ctx, cidr); err != nil {
		s.logger.Error(err.Error())
		return &pb.OkResponse{Ok: false}, status.Error(codes.Internal, ErrRemoveFromBlacklist.Error())
	}

	return &pb.OkResponse{Ok: true}, nil
}
