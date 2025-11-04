package mcp

import (
	"context"
	"errors"
	"fmt"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sashabaranov/go-openai"
	"log"
	"sync"
	"time"
)

// MCPClient 单个 MCP 服务器客户端
type MCPClient struct {
	cancel  context.CancelFunc
	name    string // 服务器名称
	host    string
	session *mcp.ClientSession // MCP 会话
	mu      sync.Mutex         // 并发安全锁
}

// NewMCPSeverClient 创建 MCP 客户端
func NewMCPSeverClient(name, host string) *MCPClient {
	return &MCPClient{name: name, host: host}
}

// Connect 连接 MCP 服务器
func (m *MCPClient) Connect(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.session != nil {
		return nil // 已连接
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "mcp-chat-client",
		Version: "1.0.0",
	}, nil)

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer func() {
		if m.session == nil {
			cancel()
		}
	}()

	session, err := client.Connect(ctxWithTimeout, &mcp.StreamableClientTransport{Endpoint: m.host}, nil)
	m.cancel = cancel
	if err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}
	m.session = session
	log.Printf("MCP 服务器[%s]已连接", m.name)
	return nil
}

// ListTools 获取服务器上的工具列表（转换为 OpenAI 格式）
func (m *MCPClient) ListTools(ctx context.Context) ([]openai.Tool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session == nil {
		return nil, errors.New("未连接 MCP 服务器")
	}

	// 调用 MCP 接口获取工具列表
	toolsResp, err := m.session.ListTools(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("获取工具列表失败: %w", err)
	}

	// 转换为 OpenAI Tool 格式
	tools := make([]openai.Tool, 0)
	for _, item := range toolsResp.Tools {
		toolsInputSchema := item.InputSchema.(map[string]any)
		tools = append(tools, openai.Tool{
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

// ExecuteTool 执行服务器上的指定工具
func (m *MCPClient) ExecuteTool(ctx context.Context, toolName string, args map[string]interface{}) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.session == nil {
		return nil, errors.New("未连接 MCP 服务器")
	}

	// 带重试的工具调用（最多3次）
	var result interface{}
	var err error
	for i := 0; i < 3; i++ {
		result, err = m.session.CallTool(ctx, &mcp.CallToolParams{
			Name:      toolName,
			Arguments: args,
		})
		if err == nil {
			return result, nil
		}
		if i < 2 {
			log.Printf("工具[%s]调用失败，重试中（%d/3）: %v", toolName, i+1, err)
			time.Sleep(1 * time.Second)
		}
	}
	return nil, fmt.Errorf("工具调用失败: %w", err)
}

// Close 关闭 MCP 连接
func (m *MCPClient) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	defer m.cancel()
	if m.session != nil {
		if err := m.session.Close(); err != nil {
			return fmt.Errorf("关闭连接失败: %w", err)
		}
		m.session = nil
		log.Printf("MCP 服务器[%s]已断开", m.name)
	}
	return nil
}
