package server

import (
	"hcs-agent/pkg/config"
	"hcs-agent/pkg/log"
	"hcs-agent/pkg/server/chatbox"
	"sync"
)

var (
	chatBoxServerOnce sync.Once
	chatBoxS          *chatbox.Server
)

func NewChatBoxServer(conf *config.Config, logger *log.Logger) *chatbox.Server {
	chatBoxServerOnce.Do(func() {
		chatBoxS = &chatbox.Server{
			Logger: logger,
			Conf:   conf,
		}
	})
	return chatBoxS
}
