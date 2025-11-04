package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var Conf *Config

func NewConfig(p string) *viper.Viper {
	envConf := os.Getenv("APP_CONF")
	if envConf == "" {
		envConf = p
	}
	fmt.Println("load conf file:", envConf)
	return getConfig(envConf)
}

func getConfig(path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(path)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return conf
}

// EnvConfig 环境配置（如开发/生产环境）
type EnvConfig struct {
	Env string `mapstructure:"env"` // 环境标识（local/dev/prod等）
}

// HTTPConfig HTTP服务配置
type HTTPConfig struct {
	Host string `mapstructure:"host"` // 服务监听地址
	Port int    `mapstructure:"port"` // 服务监听端口
}

// DingTalkConfig 钉钉配置
type DingTalkConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

// ESConfig Elasticsearch配置
type ESConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// ModelConfig 模型配置
type ModelConfig struct {
	Analysis string `mapstructure:"analysis"` // 分析用模型
	Answer   string `mapstructure:"answer"`   // 回答用模型
	URL      string `mapstructure:"url"`      // 模型接口地址
	Token    string `mapstructure:"token"`    // 模型访问令牌
}

// LogConfig 日志配置
type LogConfig struct {
	LogLevel    string `mapstructure:"log_level"`     // 日志级别（debug/info/warn/error）
	Mode        string `mapstructure:"mode"`          // 输出模式（file/console/both）
	Encoding    string `mapstructure:"encoding"`      // 编码格式（json/console）
	LogFileName string `mapstructure:"log_file_name"` // 日志文件路径
	MaxBackups  int    `mapstructure:"max_backups"`   // 最大备份文件数
	MaxAge      int    `mapstructure:"max_age"`       // 日志文件保留天数
	MaxSize     int    `mapstructure:"max_size"`      // 单个日志文件最大大小（MB）
	Compress    bool   `mapstructure:"compress"`      // 是否压缩备份文件
}

// Config 总配置结构体（聚合所有子配置）
type Config struct {
	Env      string         `mapstructure:"env"`      // 环境标识（直接映射顶级字段）
	HTTP     HTTPConfig     `mapstructure:"http"`     // HTTP服务配置
	DingTalk DingTalkConfig `mapstructure:"dingtalk"` // 钉钉配置
	Es       ESConfig       `mapstructure:"es"`       // ES配置
	Model    ModelConfig    `mapstructure:"model"`    // 模型配置
	Log      LogConfig      `mapstructure:"log"`      // 日志配置
	Mcp      Mcp            `mapstructure:"mcp"`
}

type Mcp struct {
	Servers map[string]Server `mapstructure:"mcp_servers"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
