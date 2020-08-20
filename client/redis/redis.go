package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Owen-Zhang/zsf/conf"
	"github.com/Owen-Zhang/zsf/logger"
	"github.com/go-redis/redis/v8"
)

//Redis redis包装实体
type Redis struct {
	config *Config
	client redis.Cmdable
}

//New 实例化redis客户端
func New() (*Redis, error) {
	cnf := defaultRedisConfig()
	if err := conf.UnmarshalFile("redis.yaml", cnf); err != nil {
		logger.FrameLog.Errorf("读取日志配制信息出错:%v", err)
		return nil, err
	}
	if cnf.Addrs == nil || len(cnf.Addrs) == 0 {
		return nil, fmt.Errorf("redis.yaml -> addrs数组不能为空")
	}
	clusterClient := redis.NewClusterClient(
		&redis.ClusterOptions{
			Addrs:        cnf.Addrs,
			Password:     cnf.Password,
			ReadTimeout:  time.Duration(cnf.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cnf.WriteTimeout) * time.Second,
		})
	if err := clusterClient.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("连接redis cluster出现错误:%+v", err)
	}

	return &Redis{
		config: cnf,
		client: clusterClient,
	}, nil
}
