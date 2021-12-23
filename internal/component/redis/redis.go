// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation and Dapr Contributors.
// Licensed under the MIT License.
// ------------------------------------------------------------

package redis

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	ClusterType = "cluster"
	NodeType    = "node"
)

func ParseClientFromProperties(properties map[string]string, defaultSettings *Settings) (client redis.UniversalClient, settings *Settings, err error) {
	if defaultSettings == nil {
		settings = &Settings{}
	} else {
		settings = defaultSettings
	}
	err = settings.Decode(properties) // 编码
	if err != nil {
		return nil, nil, fmt.Errorf("redis client configuration error: %w", err)
	}
	if settings.Failover {
		return newFailoverClient(settings), settings, nil
	}

	return newClient(settings), settings, nil
}

// redis 使用哨兵进行故障转移
func newFailoverClient(s *Settings) redis.UniversalClient {
	if s == nil {
		return nil
	}
	opts := &redis.FailoverOptions{
		DB:                 s.DB,
		MasterName:         s.SentinelMasterName,
		SentinelAddrs:      []string{s.Host},
		Password:           s.Password,
		Username:           s.Username,
		MaxRetries:         s.RedisMaxRetries,
		MaxRetryBackoff:    time.Duration(s.RedisMaxRetryInterval),
		MinRetryBackoff:    time.Duration(s.RedisMinRetryInterval),
		DialTimeout:        time.Duration(s.DialTimeout),
		ReadTimeout:        time.Duration(s.ReadTimeout),
		WriteTimeout:       time.Duration(s.WriteTimeout),
		PoolSize:           s.PoolSize,
		MaxConnAge:         time.Duration(s.MaxConnAge),
		MinIdleConns:       s.MinIdleConns,
		PoolTimeout:        time.Duration(s.PoolTimeout),
		IdleCheckFrequency: time.Duration(s.IdleCheckFrequency),
		IdleTimeout:        time.Duration(s.IdleTimeout),
	}

	/* #nosec */
	if s.EnableTLS {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: s.EnableTLS,
		}
	}

	if s.RedisType == ClusterType {
		opts.SentinelAddrs = strings.Split(s.Host, ",")

		return redis.NewFailoverClusterClient(opts)
	}

	return redis.NewFailoverClient(opts)
}

// redis 使用不使用哨兵进行故障转移
func newClient(s *Settings) redis.UniversalClient {
	if s == nil {
		return nil
	}
	if s.RedisType == ClusterType {
		options := &redis.ClusterOptions{
			Addrs:              strings.Split(s.Host, ","),
			Password:           s.Password,
			Username:           s.Username,
			MaxRetries:         s.RedisMaxRetries,
			MaxRetryBackoff:    time.Duration(s.RedisMaxRetryInterval),
			MinRetryBackoff:    time.Duration(s.RedisMinRetryInterval),
			DialTimeout:        time.Duration(s.DialTimeout),
			ReadTimeout:        time.Duration(s.ReadTimeout),
			WriteTimeout:       time.Duration(s.WriteTimeout),
			PoolSize:           s.PoolSize,
			MaxConnAge:         time.Duration(s.MaxConnAge),
			MinIdleConns:       s.MinIdleConns,
			PoolTimeout:        time.Duration(s.PoolTimeout),
			IdleCheckFrequency: time.Duration(s.IdleCheckFrequency),
			IdleTimeout:        time.Duration(s.IdleTimeout),
		}
		/* #nosec */
		if s.EnableTLS {
			options.TLSConfig = &tls.Config{
				InsecureSkipVerify: s.EnableTLS,
			}
		}

		return redis.NewClusterClient(options)
	}

	options := &redis.Options{
		Addr:               s.Host,
		Password:           s.Password,
		Username:           s.Username,
		DB:                 s.DB,
		MaxRetries:         s.RedisMaxRetries,
		MaxRetryBackoff:    time.Duration(s.RedisMaxRetryInterval),
		MinRetryBackoff:    time.Duration(s.RedisMinRetryInterval),
		DialTimeout:        time.Duration(s.DialTimeout),
		ReadTimeout:        time.Duration(s.ReadTimeout),
		WriteTimeout:       time.Duration(s.WriteTimeout),
		PoolSize:           s.PoolSize,
		MaxConnAge:         time.Duration(s.MaxConnAge),
		MinIdleConns:       s.MinIdleConns,
		PoolTimeout:        time.Duration(s.PoolTimeout),
		IdleCheckFrequency: time.Duration(s.IdleCheckFrequency),
		IdleTimeout:        time.Duration(s.IdleTimeout),
	}

	/* #nosec */
	if s.EnableTLS {
		options.TLSConfig = &tls.Config{
			InsecureSkipVerify: s.EnableTLS,
		}
	}

	return redis.NewClient(options)
}
// redis.UniversalClient
// 是一个抽象的客户端，根据所提供的选项
// 可以连接到集群，或有哨兵支持的故障转移实例
// 或简单的单实例服务器。这对于测试
// 集群特定的应用程序。