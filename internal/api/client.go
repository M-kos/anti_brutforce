package api

import (
	"fmt"

	"github.com/M-kos/anti_brutforce/internal/api/pb"
	"github.com/M-kos/anti_brutforce/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientAPI struct {
	Pb   pb.AntiBruteforceClient
	Conn *grpc.ClientConn
}

func NewClient(conf *config.Config) (*ClientAPI, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", "localhost", conf.GRPCPopt),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewAntiBruteforceClient(conn)

	return &ClientAPI{Pb: client, Conn: conn}, nil
}

func (c *ClientAPI) Close() {
	if c.Conn == nil {
		return
	}

	c.Conn.Close()
}
