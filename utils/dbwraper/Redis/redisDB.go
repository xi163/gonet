package Redis

import (
	"context"
	"time"

	"github.com/cwloo/gonet/logs"
	"github.com/dtm-labs/rockscache"
	go_redis "github.com/go-redis/redis/v8"
)

type DB struct {
	DB go_redis.UniversalClient
	Rc [2]*rockscache.Client
}

func (s *DB) Init(conf Cfg) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.DB = go_redis.NewUniversalClient(&go_redis.UniversalOptions{
		Addrs:    conf.Addr,
		Username: conf.Username,
		Password: conf.Password,
		// MasterName:   conf.Master,
		// DB:           conf.DB,
		PoolSize: 50,
		// DialTimeout:  5 * time.Second,
		// ReadTimeout:  5 * time.Second,
		// WriteTimeout: 5 * time.Second,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
	})
	_, err := s.DB.Ping(ctx).Result()
	if err != nil {
		logs.Fatalf(err.Error())
	}
	s.Rc[0] = rockscache.NewClient(s.DB, rockscache.NewDefaultOptions())
	s.Rc[0].Options.StrongConsistency = true

	s.Rc[1] = rockscache.NewClient(s.DB, rockscache.NewDefaultOptions())
	s.Rc[1].Options.StrongConsistency = false
	logs.Debugf("ok")
}
