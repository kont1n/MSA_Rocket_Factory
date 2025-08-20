package redis

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
)

func (c *client) SAdd(ctx context.Context, key, value string) error {
	return c.withConn(ctx, func(ctx context.Context, conn redigo.Conn) error {
		_, err := conn.Do("SADD", key, value)
		return err
	})
}

func (c *client) SRem(ctx context.Context, key, value string) error {
	return c.withConn(ctx, func(ctx context.Context, conn redigo.Conn) error {
		_, err := conn.Do("SREM", key, value)
		return err
	})
}

func (c *client) SIsMember(ctx context.Context, key, value string) (bool, error) {
	var isMember bool
	err := c.withConn(ctx, func(ctx context.Context, conn redigo.Conn) error {
		result, err := redigo.Int(conn.Do("SISMEMBER", redigo.Args{key}.Add(value)...))
		if err != nil {
			return err
		}

		isMember = result > 0
		return nil
	})
	if err != nil {
		return false, err
	}

	return isMember, nil
}

func (c *client) SMembers(ctx context.Context, key string) ([]string, error) {
	var members []string
	err := c.withConn(ctx, func(ctx context.Context, conn redigo.Conn) error {
		result, err := redigo.Strings(conn.Do("SMEMBERS", redigo.Args{key}...))
		if err != nil {
			return err
		}

		members = result
		return nil
	})
	if err != nil {
		return nil, err
	}

	return members, nil
}
