package main

import (
	"context"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
	Name string `json:"name" jsonschema:"地点名称"`
}

type Output struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

// MockGetWeather 模拟获取天气信息
func MockGetWeather(ctx context.Context, req *mcp.CallToolRequest, input Input) (
	*mcp.CallToolResult,
	map[string]interface{},
	error,
) {
	var temperature int
	switch input.Name {
	case "杭州":
		temperature = 25
	case "北京":
		temperature = 30
	case "上海":
		temperature = 27
	default:
		return nil, map[string]interface{}{
			"error_msg": "无法查询该地址的天气信息",
		}, nil
	}

	return nil, map[string]interface{}{
		"temperature": temperature,
		"location":    input.Name,
		"unit":        "CELSIUS",
	}, nil
}

func main() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	mcp.AddTool[Input, map[string]interface{}](server, &mcp.Tool{Name: "getWeather", Description: "获取天气情况"}, MockGetWeather)
	// Run the server over stdin/stdout, until the client disconnects.
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)
	log.Printf("MCP handler listening at %s", "http://localhost:8001")
	_ = http.ListenAndServe(":8001", handler)
	select {}
}
