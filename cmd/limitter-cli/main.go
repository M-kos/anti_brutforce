package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/M-kos/anti_brutforce/internal/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ResetCmd  = "reset"
	AddCmd    = "add"
	RemoveCmd = "remove"
)

func main() {
	var (
		login    string
		ip       string
		listName string
		cidr     string
	)
	resetFlagSet := flag.NewFlagSet(ResetCmd, flag.ExitOnError)
	resetFlagSet.StringVar(&login, "login", "", "login")
	resetFlagSet.StringVar(&ip, "ip", "", "ip")

	addFlagSet := flag.NewFlagSet(AddCmd, flag.ExitOnError)
	addFlagSet.StringVar(&listName, "listname", "", "blacklabel or whitelabel")
	addFlagSet.StringVar(&cidr, "cidr", "", "cidr (192.168.0.0/24)")

	removeFlagSet := flag.NewFlagSet(RemoveCmd, flag.ExitOnError)
	removeFlagSet.StringVar(&listName, "listname", "", "blacklabel or whitelabel")
	removeFlagSet.StringVar(&cidr, "cidr", "", "cidr (192.168.0.0/24)")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("command must be passed")
		return
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err) // TODO: handle error
	}

	defer conn.Close()

	cli := pb.NewAntiBruteforceClient(conn)

	switch os.Args[1] {
	case ResetCmd:
		if err := resetFlagSet.Parse(os.Args[2:]); err != nil {
			panic(err) // TODO: handle error
		}

		fmt.Println("Login >> ", login)
		fmt.Println("IP >> ", ip)
		res, err := cli.ResetBuckets(context.Background(), &pb.ResetBucketsRequest{
			Login: login,
			Ip:    ip,
		},
		)
		if err != nil {
			panic(err) // TODO: handle error
		}

		fmt.Println("Reset Response >> ", res)

	case AddCmd:
		if err := addFlagSet.Parse(os.Args[2:]); err != nil {
			panic(err) // TODO: handle error
		}

		fmt.Println("List name >> ", listName)
		fmt.Println("Cidr >> ", cidr)

		if listName == "whitelabel" && cidr != "" {
			res, err := cli.AddToWhitelist(context.Background(), &pb.WhitelistRequest{
				Cidr: cidr,
			},
			)
			if err != nil {
				panic(err) // TODO: handle error
			}

			fmt.Println("AddCmd Response >> ", res)
		}

		if listName == "blacklabel" && cidr != "" {
			res, err := cli.AddToBlacklist(context.Background(), &pb.BlacklistRequest{
				Cidr: cidr,
			},
			)
			if err != nil {
				panic(err) // TODO: handle error
			}

			fmt.Println("AddCmd Response >> ", res)
		}

	case RemoveCmd:
		if err := removeFlagSet.Parse(os.Args[2:]); err != nil {
			panic(err) // TODO: handle error
		}

		fmt.Println("List name >> ", listName)
		fmt.Println("Cidr >> ", cidr)

		if listName == "whitelabel" && cidr != "" {
			res, err := cli.RemoveFromWhitelist(context.Background(), &pb.WhitelistRequest{
				Cidr: cidr,
			},
			)
			if err != nil {
				panic(err) // TODO: handle error
			}

			fmt.Println("RemoveCmd Response >> ", res)
		}

		if listName == "blacklabel" && cidr != "" {
			res, err := cli.RemoveFromBlacklist(context.Background(), &pb.BlacklistRequest{
				Cidr: cidr,
			},
			)
			if err != nil {
				panic(err) // TODO: handle error
			}

			fmt.Println("RemoveCmd Response >> ", res)
		}
	}
}
