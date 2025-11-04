package handler

import (
	"bytes"
	"context"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"hcs-agent/pkg/mcp"
	"hcs-agent/pkg/utils"
	"k8s.io/klog/v2"
	"strings"
)

// OnChatBotMessageReceived 全局处理器函数
func OnChatBotMessageReceived(ctx context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
	klog.Infof("UserRequest:%s", utils.Dumps(data))
	mcpManager, err := mcp.NewMCPManager()
	if err != nil {
		klog.Errorf("init es mcpManager error %v", err)
		return []byte{}, err
	}
	defer mcpManager.CloseAll()
	if err = mcpManager.ConnectAll(ctx); err != nil {
		klog.Errorf("connect mcp server error %v", err)
		return []byte{}, err
	}
	mcp.NewLLMClientOnce()
	chatSession := mcp.NewChatSession(mcp.LLMClientO, mcpManager, data)
	if err = chatSession.Init(ctx); err != nil {
		klog.Errorf("init chat session error %v", err)
		return []byte{}, err
	}
	question := strings.TrimSpace(data.Text.Content)
	if err = chatSession.HandleUserInput(ctx, question); err != nil {
		klog.Errorf("handle user input error %v", err)
		return []byte{}, err
	}
	return []byte{}, nil
}

func generateMarkdownTitle(content []byte) []byte {
	idx := bytes.IndexByte(content, '\n')
	if idx == -1 {
		return content
	}
	return content[:idx]
}
