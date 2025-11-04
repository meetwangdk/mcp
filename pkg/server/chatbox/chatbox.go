package chatbox

import (
	"context"
	"flag"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/event"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/payload"
	"hcs-agent/internal/handler"
	"hcs-agent/pkg/config"
	"hcs-agent/pkg/log"
	"k8s.io/klog/v2"
)

type Server struct {
	Conf   *config.Config
	Logger *log.Logger
	stop   func()
}

func OnEventReceived(ctx context.Context, df *payload.DataFrame) (frameResp *payload.DataFrameResponse, err error) {
	eventHeader := event.NewEventHeaderFromDataFrame(df)

	logger.GetLogger().Infof("received event, eventId=[%s] eventBornTime=[%d] eventCorpId=[%s] eventType=[%s] eventUnifiedAppId=[%s] data=[%s]",
		eventHeader.EventId,
		eventHeader.EventBornTime,
		eventHeader.EventCorpId,
		eventHeader.EventType,
		eventHeader.EventUnifiedAppId,
		df.Data)

	//TODO 处理事件
	frameResp = payload.NewSuccessDataFrameResponse()
	if err := frameResp.SetJson(event.NewEventProcessResultSuccess()); err != nil {
		return nil, err
	}

	return
}

func (c *Server) Start(ctx context.Context) error {
	klog.InitFlags(flag.CommandLine)
	flag.Parse()
	cli := client.NewStreamClient(
		client.WithAppCredential(client.NewAppCredentialConfig(
			c.Conf.DingTalk.ClientID,
			c.Conf.DingTalk.ClientSecret)),
		client.WithUserAgent(client.NewDingtalkGoSDKUserAgent()),
	)
	cli.RegisterAllEventRouter(OnEventReceived)
	cli.RegisterChatBotCallbackRouter(handler.OnChatBotMessageReceived)
	c.stop = cli.Close
	klog.Info(" Start chatbox Server...")
	return cli.Start(ctx)
}
func (c *Server) Stop(ctx context.Context) error {
	c.Logger.Sugar().Info("Shutting down chatbox Server...")
	if c.stop != nil {
		c.stop()
	}
	return nil
}
