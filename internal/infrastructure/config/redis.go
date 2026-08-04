package config

import "fmt"

type Redis struct {
	Host string
	Port int
}

func (redis *Redis) GetAddr() string {
	return fmt.Sprintf("%v:%v", redis.Host, redis.Port)
}
