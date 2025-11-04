package main

import (
	"context"
	"errors"
	"flag"
	"hcs-agent/pkg/config"
	"hcs-agent/pkg/log"
)

func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	newConfig := config.NewConfig(*envConf)
	logger := log.NewLog(newConfig)
	if err := newConfig.Unmarshal(&config.Conf); err != nil {
		panic(errors.New("parse config error"))

	}
	app, cleanup, err := NewChatBoxWire(config.Conf, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
