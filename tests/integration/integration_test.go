//go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/M-kos/anti_brutforce/internal/api"
	"github.com/M-kos/anti_brutforce/internal/api/pb"
	"github.com/M-kos/anti_brutforce/internal/config"
	"github.com/M-kos/anti_brutforce/internal/db"
	"github.com/M-kos/anti_brutforce/internal/lib/logger"
	"github.com/M-kos/anti_brutforce/internal/lib/ratelimit"
	"github.com/M-kos/anti_brutforce/internal/models"
	labellistservice "github.com/M-kos/anti_brutforce/internal/services/label-list-service"
	limitterservice "github.com/M-kos/anti_brutforce/internal/services/limitter-service"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type LimitterSuite struct {
	suite.Suite
	server          *api.ServerAPI
	redis           *db.RedisDB
	whitelabelList  *labellistservice.LabelListService
	blacklabelList  *labellistservice.LabelListService
	limitterService *limitterservice.LimitterService
	ctx             context.Context
}

func (s *LimitterSuite) SetupSuite() {
	// conf := config.LoadConfig()
	file, _ := os.ReadFile("../../configs/config.json")

	var conf config.Config

	_ = json.Unmarshal(file, &conf)

	l := logger.NewLogger()
	redis := db.NewRedisDB(&conf, l)

	s.redis = redis

	loginRateLimitter := ratelimit.NewRateLimiter(conf.LoginNumberAttempts, conf.Timeout, redis)
	passwordRateLimitter := ratelimit.NewRateLimiter(conf.PasswordNumberAttempts, conf.Timeout, redis)
	ipRateLimitter := ratelimit.NewRateLimiter(conf.IPNumberAttempts, conf.Timeout, redis)

	whitelabelList := labellistservice.NewLabelListService(
		labellistservice.NewLabelListRepository(redis, labellistservice.WhitelistKey),
		l,
		labellistservice.WhitelistKey,
	)
	blacklabelList := labellistservice.NewLabelListService(
		labellistservice.NewLabelListRepository(redis, labellistservice.BlacklistKey),
		l,
		labellistservice.BlacklistKey,
	)

	s.whitelabelList = whitelabelList
	s.blacklabelList = blacklabelList

	limitterService := limitterservice.NewLimitterService(loginRateLimitter, passwordRateLimitter, ipRateLimitter, whitelabelList, blacklabelList, l)

	s.limitterService = limitterService

	server := api.NewServerAPI(limitterService, whitelabelList, blacklabelList, l)
	s.server = server

	s.ctx = context.Background()
}

func (s *LimitterSuite) SetupTest() {
	_, err := s.redis.Client.FlushDB(s.ctx).Result()
	s.Require().NoError(err)
}

func (s *LimitterSuite) TestCheckItemWithEmpyDBSuccess() {
	res, err := s.server.CheckCredentials(s.ctx, &pb.CheckCredentialsRequest{Login: "test", Password: "12345", Ip: "127.0.0.1"})

	s.Require().NoError(err)
	s.Require().Equal(res, &pb.OkResponse{Ok: true})

	val := s.getBucketFromDB("test")

	s.Require().Equal(val.Count, uint(1))
	s.Require().Equal(val.Key, "test")
}

func (s *LimitterSuite) TestCheckWhenItemExistSuccess() {
	item := s.defaultBucket()
	item.Count = 1

	data, err := json.Marshal(&item)

	s.Require().NoError(err)

	_, err = s.redis.Client.Set(s.ctx, item.Key, data, time.Hour).Result()
	s.Require().NoError(err)

	res, err := s.server.CheckCredentials(s.ctx, &pb.CheckCredentialsRequest{Login: item.Key, Password: "12345", Ip: "127.0.0.1"})

	s.Require().NoError(err)
	s.Require().Equal(res, &pb.OkResponse{Ok: true})

	val := s.getBucketFromDB(item.Key)

	s.Require().Equal(val.Count, uint(2))
	s.Require().Equal(val.Key, item.Key)
}

func (s *LimitterSuite) TestCheckWhenItemExistAndAtttempsHaveEnded() {
	item := s.addDefaultBucketToDB()

	res, err := s.server.CheckCredentials(s.ctx, &pb.CheckCredentialsRequest{Login: item.Key, Password: "12345", Ip: "127.0.0.1"})

	s.Require().Error(err)
	s.Require().Equal(res, &pb.OkResponse{Ok: false})

	val := s.getBucketFromDB(item.Key)

	s.Require().Equal(val.Count, item.Count)
	s.Require().Equal(val.Key, item.Key)
}

func (s *LimitterSuite) TestCheckWhenItemIsInWhitelabelList() {
	item := s.addDefaultBucketToDB()

	s.addDefaultListToDB(labellistservice.WhitelistKey)

	res, err := s.server.CheckCredentials(s.ctx, &pb.CheckCredentialsRequest{Login: item.Key, Password: "12345", Ip: "127.0.0.1"})
	s.Require().NoError(err)
	s.Require().Equal(res, &pb.OkResponse{Ok: true})

	val := s.getBucketFromDB(item.Key)

	s.Require().Equal(val.Count, item.Count)
	s.Require().Equal(val.Key, item.Key)
}

func (s *LimitterSuite) TestCheckWhenItemIsInBlacklabelList() {
	item := s.addDefaultBucketToDB()

	s.addDefaultListToDB(labellistservice.BlacklistKey)

	res, err := s.server.CheckCredentials(s.ctx, &pb.CheckCredentialsRequest{Login: item.Key, Password: "12345", Ip: "127.0.0.1"})
	s.Require().Error(err)
	s.Require().Equal(res, &pb.OkResponse{Ok: false})

	val := s.getBucketFromDB(item.Key)

	s.Require().Equal(val.Count, item.Count)
	s.Require().Equal(val.Key, item.Key)
}

func (s *LimitterSuite) TestResetBucketByLoginAndIpAndItemExist() {
	item := s.addDefaultBucketToDB()

	res, err := s.server.ResetBuckets(s.ctx, &pb.ResetBucketsRequest{Login: item.Key, Ip: "127.0.0.1"})
	s.Require().NoError(err)
	s.Require().Equal(res, &pb.OkResponse{Ok: true})

	v, err := s.redis.Client.Get(s.ctx, item.Key).Result()
	s.Require().ErrorIs(err, redis.Nil)
	s.Require().Equal(v, "")
}

func (s *LimitterSuite) TestResetBucketByLoginAndIpAndItemNotExist() {
	res, err := s.server.ResetBuckets(s.ctx, &pb.ResetBucketsRequest{Login: "test", Ip: "127.0.0.1"})
	s.Require().NoError(err)
	s.Require().Equal(res, &pb.OkResponse{Ok: true})
}

func (s *LimitterSuite) TestAddToWhitelabelList() {
	list := s.defaultList()

	res, err := s.server.AddToWhitelist(s.ctx, &pb.WhitelistRequest{Cidr: "127.0.0.1/24"})

	s.Require().NoError(err)
	s.Require().Equal(res, &pb.OkResponse{Ok: true})

	l := s.getListFromDB(labellistservice.WhitelistKey)
	s.Require().Equal(l.Values, list.Values)
}

func (s *LimitterSuite) TestRemoveFromWhitelabelList() {
	s.addDefaultListToDB(labellistservice.WhitelistKey)

	res, err := s.server.RemoveFromWhitelist(s.ctx, &pb.WhitelistRequest{Cidr: "127.0.0.1/24"})

	s.Require().NoError(err)
	s.Require().Equal(res, &pb.OkResponse{Ok: true})

	l := s.getListFromDB(labellistservice.WhitelistKey)
	s.Require().Equal(l.Values, []string{})
}

func (s *LimitterSuite) TestAddToBlacklabelList() {
	list := s.defaultList()

	res, err := s.server.AddToBlacklist(s.ctx, &pb.BlacklistRequest{Cidr: "127.0.0.1/24"})

	s.Require().NoError(err)
	s.Require().Equal(res, &pb.OkResponse{Ok: true})

	l := s.getListFromDB(labellistservice.BlacklistKey)
	s.Require().Equal(l.Values, list.Values)
}

func (s *LimitterSuite) TestRemoveFromBlacklabelList() {
	s.addDefaultListToDB(labellistservice.BlacklistKey)

	res, err := s.server.RemoveFromBlacklist(s.ctx, &pb.BlacklistRequest{Cidr: "127.0.0.1/24"})

	s.Require().NoError(err)
	s.Require().Equal(res, &pb.OkResponse{Ok: true})

	l := s.getListFromDB(labellistservice.BlacklistKey)
	s.Require().Equal(l.Values, []string{})
}

func (s *LimitterSuite) defaultBucket() models.Bucket {
	return models.Bucket{
		Count:     10,
		Key:       "test",
		StartTime: time.Now(),
	}
}

func (s *LimitterSuite) defaultList() *models.LabelList {
	return &models.LabelList{
		Values: []string{"127.0.0.1/24"},
	}
}

func (s *LimitterSuite) addDefaultBucketToDB() models.Bucket {
	item := s.defaultBucket()

	data, err := json.Marshal(&item)
	s.Require().NoError(err)
	_, err = s.redis.Client.Set(s.ctx, item.Key, data, time.Hour).Result()
	s.Require().NoError(err)

	return item
}

func (s *LimitterSuite) getBucketFromDB(key string) models.Bucket {
	v, err := s.redis.Client.Get(s.ctx, key).Result()
	s.Require().NoError(err)

	var val models.Bucket

	err = json.Unmarshal([]byte(v), &val)
	s.Require().NoError(err)
	return val
}

func (s *LimitterSuite) addDefaultListToDB(key string) {
	list := s.defaultList()

	l, err := json.Marshal(list)
	s.Require().NoError(err)
	_, err = s.redis.Client.Set(s.ctx, key, string(l), time.Hour).Result()
	s.Require().NoError(err)
}

func (s *LimitterSuite) getListFromDB(key string) *models.LabelList {
	v, err := s.redis.Client.Get(s.ctx, key).Result()
	s.Require().NoError(err)

	var val models.LabelList

	err = json.Unmarshal([]byte(v), &val)
	s.Require().NoError(err)

	return &val
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(LimitterSuite))
}
