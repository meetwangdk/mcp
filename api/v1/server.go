package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sashabaranov/go-openai"
	"hcs-agent/pkg/config"
	"log"
	"sync"
)

// 服务器管理
type Server struct {
	Name      string
	config    config.Server
	session   *mcp.ClientSession
	mu        sync.Mutex
	exitStack []func() error
}

func NewServer(name string, config config.Server) *Server {
	return &Server{
		Name:   name,
		config: config,
	}
}

func (s *Server) Initialize(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	client := mcp.NewClient(&mcp.Implementation{
		Name:    fmt.Sprintf(" %s-time-client", s.Name),
		Version: "1.0.0",
	}, nil)

	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: s.config.Host}, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	s.session = session
	return nil
}

func (s *Server) ListTools(ctx context.Context) ([]*openai.Tool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.session == nil {
		return nil, errors.New("服务器未初始化")
	}

	toolsResp, err := s.session.ListTools(ctx, nil)
	if err != nil {
		return nil, err
	}

	tools := make([]*openai.Tool, 0)
	for _, item := range toolsResp.Tools {
		toolsInputSchema := item.InputSchema.(map[string]any)
		tools = append(tools, &openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        item.Name,
				Description: item.Description,
				Parameters:  toolsInputSchema,
			},
		})
	}

	return tools, nil
}
