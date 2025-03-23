package api

import (
	"fmt"

	"github.com/M-kos/anti_brutforce/internal/api/pb"
	"github.com/M-kos/anti_brutforce/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientApi struct {
	Pb   pb.AntiBruteforceClient
	Conn *grpc.ClientConn
}

func NewClient(conf *config.Config) *ClientApi {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", "localhost", conf.GRPCPopt), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := pb.NewAntiBruteforceClient(conn)

	return &ClientApi{Pb: client, Conn: conn}
}

func (c *ClientApi) Close() {
	c.Conn.Close()
}
