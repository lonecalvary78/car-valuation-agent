package config

import (
	"fmt"
	"time"
)

type Server struct {
	Host        string
	Port        int
	ReadTimeout time.Duration
	Writetimout time.Duration
}

func (server Server) GetAddr() string {
	return fmt.Sprintf("%v:%v", server.Host, server.Port)
}
