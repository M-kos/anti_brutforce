package limittercliservice

import (
	"context"
	"flag"

	"github.com/M-kos/anti_brutforce/internal/api"
	"github.com/M-kos/anti_brutforce/internal/api/pb"
	"github.com/M-kos/anti_brutforce/internal/lib/logger"
)

type LimitterCliService struct {
	Login    string
	Ip       string
	ListName string
	Cidr     string

	resetFlagSet  *flag.FlagSet
	addFlagSet    *flag.FlagSet
	removeFlagSet *flag.FlagSet

	client *api.ClientApi
	log    *logger.Logger
}

const (
	ResetCmd  = "reset"
	AddCmd    = "add"
	RemoveCmd = "remove"
)

func NewLimitterCliService(client *api.ClientApi, log *logger.Logger) *LimitterCliService {
	resetFlagSet := flag.NewFlagSet(ResetCmd, flag.ExitOnError)
	addFlagSet := flag.NewFlagSet(AddCmd, flag.ExitOnError)
	removeFlagSet := flag.NewFlagSet(RemoveCmd, flag.ExitOnError)

	return &LimitterCliService{
		client:        client,
		log:           log,
		resetFlagSet:  resetFlagSet,
		addFlagSet:    addFlagSet,
		removeFlagSet: removeFlagSet,
	}
}

func (l *LimitterCliService) Run() {
	l.resetFlagSet.StringVar(&l.Login, "login", "", "login")
	l.resetFlagSet.StringVar(&l.Ip, "ip", "", "ip")

	l.addFlagSet.StringVar(&l.ListName, "listname", "", "blacklabel or whitelabel")
	l.addFlagSet.StringVar(&l.Cidr, "cidr", "", "cidr (192.168.0.0/24)")

	l.removeFlagSet.StringVar(&l.ListName, "listname", "", "'blacklabel' or 'whitelabel'")
	l.removeFlagSet.StringVar(&l.Cidr, "cidr", "", "cidr, for example '192.168.0.0/24'")

	flag.Parse()
}

func (l *LimitterCliService) Reset(args []string) (bool, error) {
	if err := l.resetFlagSet.Parse(args); err != nil {
		l.log.Error("LimitterCliService: Reset: ", err.Error())
		return false, err
	}

	_, err := l.client.Pb.ResetBuckets(context.Background(), &pb.ResetBucketsRequest{
		Login: l.Login,
		Ip:    l.Ip,
	},
	)
	if err != nil {
		l.log.Error("LimitterCliService: Reset: ResetBuckets: ", err.Error())
		return false, err
	}

	return true, nil
}

func (l *LimitterCliService) Add(args []string) (bool, error) {
	if err := l.addFlagSet.Parse(args); err != nil {
		l.log.Error("LimitterCliService: Add: ", err.Error())
		return false, err
	}

	ctx := context.Background()

	if l.ListName == "whitelabel" && l.Cidr != "" {
		_, err := l.client.Pb.AddToWhitelist(ctx, &pb.WhitelistRequest{
			Cidr: l.Cidr,
		},
		)
		if err != nil {
			l.log.Error("LimitterCliService: Add: Whitelist: ", err.Error())
			return false, err
		}
	}

	if l.ListName == "blacklabel" && l.Cidr != "" {
		_, err := l.client.Pb.AddToBlacklist(context.Background(), &pb.BlacklistRequest{
			Cidr: l.Cidr,
		},
		)
		if err != nil {
			l.log.Error("LimitterCliService: Add: Blacklist: ", err.Error())
			return false, err
		}
	}

	return true, nil
}

func (l *LimitterCliService) Remove(args []string) (bool, error) {
	if err := l.removeFlagSet.Parse(args); err != nil {
		l.log.Error("LimitterCliService: Remove: ", err.Error())
		return false, err
	}

	ctx := context.Background()

	if l.ListName == "whitelabel" && l.Cidr != "" {
		_, err := l.client.Pb.RemoveFromWhitelist(ctx, &pb.WhitelistRequest{
			Cidr: l.Cidr,
		},
		)
		if err != nil {
			l.log.Error("LimitterCliService: Remove: Whitelist: ", err.Error())
			return false, err
		}
	}

	if l.ListName == "blacklabel" && l.Cidr != "" {
		_, err := l.client.Pb.RemoveFromBlacklist(context.Background(), &pb.BlacklistRequest{
			Cidr: l.Cidr,
		},
		)
		if err != nil {
			l.log.Error("LimitterCliService: Remove: Blacklist: ", err.Error())
			return false, err
		}
	}

	return true, nil
}
