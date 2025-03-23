package main

import (
	"fmt"
	"os"

	"github.com/M-kos/anti_brutforce/internal/api"
	"github.com/M-kos/anti_brutforce/internal/config"
	"github.com/M-kos/anti_brutforce/internal/lib/logger"
	limittercliservice "github.com/M-kos/anti_brutforce/internal/services/limitter-cli-service"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("command must be passed")
		return
	}

	conf := config.LoadConfig()
	l := logger.NewLogger()
	client := api.NewClient(conf)

	defer client.Close()

	limitterCLiService := limittercliservice.NewLimitterCliService(client, l)

	switch os.Args[1] {
	case limittercliservice.ResetCmd:
		l.Info("Args >> ", os.Args[2:])
		ok, err := limitterCLiService.Reset(os.Args[2:])
		if err != nil {
			fmt.Println("bucket reset error")
		}

		if ok {
			fmt.Println("bucket reseted")
		}

	case limittercliservice.AddCmd:
		l.Info("Args >> ", os.Args[2:])
		ok, err := limitterCLiService.Add(os.Args[2:])
		if err != nil {
			fmt.Println("add to list error")
		}

		if ok {
			fmt.Println("added to list")
		}

	case limittercliservice.RemoveCmd:
		l.Info("Args >> ", os.Args[2:])
		ok, err := limitterCLiService.Remove(os.Args[2:])
		if err != nil {
			fmt.Println("remove from list error")
		}

		if ok {
			fmt.Println("removed from list")
		}
	}
}
