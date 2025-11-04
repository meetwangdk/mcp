//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"hcs-agent/internal/biz"
	"hcs-agent/internal/data"
	"hcs-agent/internal/handler"
	"hcs-agent/internal/router"
	"hcs-agent/internal/server"
	service "hcs-agent/internal/serviece"
	"hcs-agent/pkg/app"
	"hcs-agent/pkg/config"
	"hcs-agent/pkg/log"
	"hcs-agent/pkg/server/http"
	"hcs-agent/pkg/sid"
)

func newApp(chatBoxServer *server.ChatBoxServer, httpServer *http.Server) *app.Client {
	return app.NewApp(
		app.WithServer(chatBoxServer),
		//app.WithServer(chatBoxServer, httpServer),
		app.WithName("chat box"),
	)
}

func NewChatBoxWire(*config.Config, *log.Logger) (*app.Client, func(), error) {
	panic(wire.Build(
		data.Set,
		biz.Set,
		wire.Struct(new(router.RouterDeps), "*"),
		service.ProviderSet,
		server.ProviderSet,
		handler.Set,
		sid.NewSid,
		newApp,
	))
}
