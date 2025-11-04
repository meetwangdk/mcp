package mcp

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/sync/errgroup"
	"hcs-agent/pkg/config"
	"log"
	"sync"
)

// MCPManager MCP 服务器管理器（管理多个客户端）
type MCPManager struct {
	clients      []*MCPClient // 所有 MCP 客户端
	toolMap      sync.Map     // 工具-服务器映射: toolName -> []*MCPClient
	toolMapMutex sync.RWMutex // 映射读写锁
}

// NewMCPManager 从配置文件创建 MCP 管理器
func NewMCPManager() (*MCPManager, error) {
	// 初始化 MCP Server 列表
	var serverClient []*MCPClient
	mcpServers := config.Conf.Mcp.Servers
	for name, srvConfig := range mcpServers {
		serverClient = append(serverClient, NewMCPSeverClient(name, fmt.Sprintf("%s:%s", srvConfig.Host, srvConfig.Port)))
	}
	return &MCPManager{clients: serverClient}, nil
}

// ConnectAll 连接所有 MCP 服务器并构建工具映射
func (m *MCPManager) ConnectAll(ctx context.Context) error {
	// 并发连接所有服务器
	g, errCtx := errgroup.WithContext(ctx)
	for _, client := range m.clients {
		client := client
		g.Go(func() error {
			if err := client.Connect(errCtx); err != nil {
				log.Printf("服务器[%s]连接失败: %v", client.name, err)
				return nil
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	// 连接成功后构建工具-服务器映射
	return m.buildToolMap(ctx)
}

// buildToolMap 构建工具与服务器的映射关系
func (m *MCPManager) buildToolMap(ctx context.Context) error {
	m.toolMapMutex.Lock()
	defer m.toolMapMutex.Unlock()
	m.toolMap = sync.Map{}

	// 遍历所有服务器，收集工具信息
	for _, client := range m.clients {
		tools, err := client.ListTools(ctx)
		if err != nil {
			log.Printf("服务器[%s]获取工具列表失败，跳过: %v", client.name, err)
			continue
		}

		// 记录该服务器提供的所有工具
		for _, tool := range tools {
			toolName := tool.Function.Name
			if val, ok := m.toolMap.Load(toolName); ok {
				// 工具已存在，追加服务器
				servers := val.([]*MCPClient)
				m.toolMap.Store(toolName, append(servers, client))
			} else {
				// 工具不存在，初始化服务器列表
				m.toolMap.Store(toolName, []*MCPClient{client})
			}
		}
	}

	// 打印工具映射（调试用）
	m.toolMap.Range(func(key, value interface{}) bool {
		toolName := key.(string)
		servers := value.([]*MCPClient)
		serverNames := make([]string, len(servers))
		for i, s := range servers {
			serverNames[i] = s.name
		}
		log.Printf("工具[%s]可在服务器执行: %v", toolName, serverNames)
		return true
	})
	return nil
}

// RefreshToolMap 刷新工具-服务器映射（应对动态变化）
func (m *MCPManager) RefreshToolMap(ctx context.Context) error {
	log.Println("刷新工具-服务器映射...")
	return m.buildToolMap(ctx)
}

// GetAllTools 获取所有服务器的工具列表（去重）
func (m *MCPManager) GetAllTools(ctx context.Context) ([]openai.Tool, error) {
	toolSet := make(map[string]openai.Tool) // 去重

	// 遍历所有服务器收集工具
	for _, client := range m.clients {
		tools, err := client.ListTools(ctx)
		if err != nil {
			log.Printf("服务器[%s]获取工具失败，跳过: %v", client.name, err)
			continue
		}
		// 去重存储
		for _, tool := range tools {
			toolSet[tool.Function.Name] = tool
		}
	}

	// 转换为切片返回
	allTools := make([]openai.Tool, 0, len(toolSet))
	for _, tool := range toolSet {
		allTools = append(allTools, tool)
	}
	return allTools, nil
}

// ExecuteTool 只向提供目标工具的服务器发起调用
func (m *MCPManager) ExecuteTool(ctx context.Context, toolName string, args map[string]interface{}) (interface{}, error) {
	// 读取工具映射（读锁，不阻塞其他读取）
	m.toolMapMutex.RLock()
	val, ok := m.toolMap.Load(toolName)
	m.toolMapMutex.RUnlock()

	// 工具不存在于任何服务器
	if !ok {
		return nil, fmt.Errorf("没有服务器提供工具[%s]", toolName)
	}

	// 获取提供该工具的服务器列表
	servers := val.([]*MCPClient)
	var lastErr error

	// 尝试在可用服务器上执行工具
	for _, server := range servers {
		result, err := server.ExecuteTool(ctx, toolName, args)
		if err == nil {
			return result, nil // 执行成功，返回结果
		}
		lastErr = err
		log.Printf("服务器[%s]执行工具[%s]失败: %v", server.name, toolName, err)
	}

	// 所有可用服务器均执行失败
	return nil, fmt.Errorf("工具[%s]执行失败（所有可用服务器均尝试过）: %v", toolName, lastErr)
}

// CloseAll 关闭所有 MCP 连接
func (m *MCPManager) CloseAll() {
	var wg sync.WaitGroup
	for _, client := range m.clients {
		wg.Add(1)
		go func(c *MCPClient) {
			defer wg.Done()
			if err := c.Close(); err != nil {
				log.Printf("服务器[%s]关闭失败: %v", c.name, err)
			}
		}(client)
	}
	wg.Wait()
}
