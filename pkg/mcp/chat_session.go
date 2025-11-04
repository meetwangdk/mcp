package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"strings"
)

// ChatSession 聊天会话管理器
type ChatSession struct {
	data     *chatbot.BotCallbackDataModel
	llm      *LLMClient
	mcp      *MCPManager
	messages []openai.ChatCompletionMessage
	tools    []openai.Tool // 缓存的工具列表
}

// NewChatSession 创建聊天会话
func NewChatSession(llm *LLMClient, mcp *MCPManager, data *chatbot.BotCallbackDataModel) *ChatSession {
	return &ChatSession{
		data:     data,
		llm:      llm,
		mcp:      mcp,
		messages: []openai.ChatCompletionMessage{},
	}
}

// Init 初始化会话（加载工具并生成系统提示）
func (c *ChatSession) Init(ctx context.Context) error {
	// 从 MCP 加载所有工具
	tools, err := c.mcp.GetAllTools(ctx)
	if err != nil {
		return fmt.Errorf("加载工具列表失败: %w", err)
	}
	c.tools = tools

	// 生成系统提示（包含工具说明）
	systemMsg := c.buildSystemPrompt()
	c.messages = []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: systemMsg},
	}
	return nil
}

// buildSystemPrompt 构建系统提示
func (c *ChatSession) buildSystemPrompt() string {
	var sb strings.Builder
	sb.WriteString("你可以使用以下工具辅助回答问题：\n\n")
	for _, tool := range c.tools {
		marshal, _ := json.Marshal(tool.Function)
		sb.Write(marshal)
		//fn := tool.Function
		//sb.WriteString(fmt.Sprintf("工具名：%s\n", fn.Name))
		//sb.WriteString(fmt.Sprintf("描述：%s\n", fn.Description))
		//sb.WriteString("参数：\n")
		//
		//// 解析参数信息
		//props, _ := fn.Parameters.(map[string]interface{})["properties"].(map[string]interface{})
		//required, _ := fn.Parameters.(map[string]interface{})["required"].([]string)
		//requiredSet := make(map[string]bool)
		//for _, r := range required {
		//	requiredSet[r] = true
		//}
		//
		//// 格式化参数说明
		//for param, info := range props {
		//	desc := info.(map[string]interface{})["description"].(string)
		//	if requiredSet[param] {
		//		sb.WriteString(fmt.Sprintf("  - %s：%s（必填）\n", param, desc))
		//	} else {
		//		sb.WriteString(fmt.Sprintf("  - %s：%s（可选）\n", param, desc))
		//	}
		//}
		sb.WriteString("\n")
	}
	// 工具调用格式说明
	systemMsg := fmt.Sprintf(`你是一个有帮助的助手，可以使用以下工具:

%s

根据用户的问题选择合适的工具。如果不需要工具，直接回复。

重要提示: 当你需要使用工具时，必须只返回以下格式的JSON对象，不能有其他内容:
{
    "tool": "tool-name",
    "arguments": {
        "argument-name": "value"
    }
}

收到工具的响应后:
1. 将原始数据转换为自然、对话式的响应
2. 保持响应简洁但信息丰富
3. 专注于最相关的信息
4. 使用用户问题中的适当上下文
5. 避免简单重复原始数据

请只使用上面明确定义的工具。`, sb.String())

	return systemMsg
}

// HandleUserInput 处理用户输入并生成响应
func (c *ChatSession) HandleUserInput(ctx context.Context, userInput string) error {
	// 添加用户消息到会话历史
	if userInput != "" {
		c.messages = append(c.messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		})
		fmt.Printf("\n%s\n%s\n", strings.Repeat("-", 10)+"用户输入"+strings.Repeat("-", 10), userInput)
	}

	// 调用 LLM 生成响应
	resp, err := c.llm.CreateCompletion(ctx, c.messages, c.tools)
	if err != nil {
		return fmt.Errorf("LLM 调用失败: %w", err)
	}

	// 处理响应（流式/非流式）
	if c.llm.stream {
		return c.handleStreamResponse(ctx, resp.(*openai.ChatCompletionStream))
	}
	return c.handleFullResponse(ctx, resp.(openai.ChatCompletionResponse))
}

// handleFullResponse 处理非流式响应
func (c *ChatSession) handleFullResponse(ctx context.Context, resp openai.ChatCompletionResponse) error {
	if len(resp.Choices) == 0 {
		return errors.New("LLM 未返回有效结果")
	}

	assistantMsg := resp.Choices[0].Message
	c.messages = append(c.messages, assistantMsg)
	fmt.Printf("\n%s\n%s\n", strings.Repeat("-", 10)+"LLM 响应"+strings.Repeat("-", 10), assistantMsg.Content)
	if resp.Choices[0].FinishReason == openai.FinishReasonStop {
		//发送消息通知
		replier := chatbot.NewChatbotReplier()
		err := replier.SimpleReplyMarkdown(ctx, c.data.SessionWebhook, generateMarkdownTitle([]byte(assistantMsg.Content)), []byte(assistantMsg.Content))
		if err != nil {
			return err
		}
	}
	// 处理工具调用
	if resp.Choices[0].FinishReason == openai.FinishReasonToolCalls && len(assistantMsg.ToolCalls) > 0 {
		return c.handleToolCalls(ctx, assistantMsg.ToolCalls)
	}
	return nil
}

// handleStreamResponse 处理流式响应
func (c *ChatSession) handleStreamResponse(ctx context.Context, stream *openai.ChatCompletionStream) error {
	defer stream.Close()

	var assistantMsg openai.ChatCompletionMessage
	var contentBuf strings.Builder
	fmt.Printf("\n%s\n", strings.Repeat("-", 10)+"LLM 响应（流式）"+strings.Repeat("-", 10))

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("流式响应错误: %w", err)
		}

		if len(msg.Choices) == 0 {
			continue
		}
		delta := msg.Choices[0].Delta

		// 累加内容
		if delta.Content != "" {
			fmt.Print(delta.Content)
			contentBuf.WriteString(delta.Content)
		}
		// 累加工具调用
		assistantMsg.Content = contentBuf.String()
		assistantMsg.ToolCalls = append(assistantMsg.ToolCalls, delta.ToolCalls...)
	}

	c.messages = append(c.messages, assistantMsg)
	fmt.Println() // 换行

	// 处理工具调用
	if len(assistantMsg.ToolCalls) > 0 {
		return c.handleToolCalls(ctx, assistantMsg.ToolCalls)
	}
	return nil
}

// handleToolCalls 处理工具调用并继续会话
func (c *ChatSession) handleToolCalls(ctx context.Context, toolCalls []openai.ToolCall) error {
	for _, call := range toolCalls {
		// 解析工具参数
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
			log.Printf("工具参数解析失败: %v", err)
			continue
		}

		// 调用 MCP 执行工具
		fmt.Printf("\n%s\n调用工具: %s，参数: %v\n", strings.Repeat("-", 10)+"工具调用"+strings.Repeat("-", 10), call.Function.Name, args)
		result, err := c.mcp.ExecuteTool(ctx, call.Function.Name, args)
		if err != nil {
			log.Printf("工具执行失败: %v", err)
			continue
		}

		// 将工具结果添加到会话历史
		resultJSON, _ := json.Marshal(result)
		c.messages = append(c.messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			ToolCallID: call.ID,
			Content:    string(resultJSON),
		})
		fmt.Printf("工具执行结果: %v\n", result)
	}

	// 工具调用后，让 LLM 整理结果
	return c.HandleUserInput(ctx, "")
}

func generateMarkdownTitle(content []byte) []byte {
	idx := bytes.IndexByte(content, '\n')
	if idx == -1 {
		return content
	}
	return content[:idx]
}
