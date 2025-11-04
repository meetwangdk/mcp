package config

import (
	"code.alipay.com/sigma/eavesdropping-client/go-client/common"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config 应用程序配置结构
type Config struct {
	YuQue             YuQueConfig          `json:"yuque"`
	Hcp               HcpConfigs           `json:"hcp"`
	DingTalk          DingTalkConfig       `json:"dingtalk"`
	Model             ModelConfig          `json:"model"`
	Lunettes          LunettesConfig       `json:"lunettes"`
	Links             LinksConfig          `json:"links"`
	WebkubectlEnvTest WebkubectlConfig     `json:"webkubectlEnvTest"`
	WebkubectlEnvProd WebkubectlConfig     `json:"webkubectlEnvProd"`
	Helper            HelperConfig         `json:"helper"`
	Clusters          map[string][]Cluster `json:"clusters"`
	EsOptions         *common.ESOptions
	Development       bool `json:"development"`
}

// sigma cluster
// from http://istio-ingressgateway-hcs-sit.istio-system.svc.alipay.net:8080/hapis/user.hcs.io/v1alpha1/clusters
// and https://hcs-api.alipay.com/hapis/user.hcs.io/v1alpha1/clusters
type Cluster struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DingTalkConfig 钉钉相关配置
type DingTalkConfig struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"-"`
}

// ModelConfig 模型相关配置
type ModelConfig struct {
	AnalysisModel string `json:"analysisModel"`
	AnswerModel   string `json:"answerModel"`
	URL           string `json:"url"`
	Token         string `json:"-"`
}

type YuQueConfig struct {
	BaseUrl string `json:"baseUrl"`
	Token   string `json:"token"`
}

type HcpConfigs struct {
	DevConfig  HcpConfig `json:"devConfig"`
	ProdConfig HcpConfig `json:"prodConfig"`
}

type HcpConfig struct {
	URL string `json:"hcpBaseURL"`
	AK  string `json:"ak"`
	SK  string `json:"sk"`
}

// LunettesConfig 服务器相关配置
type LunettesConfig struct {
	LunettesdiAPIBaseURL   string `json:"lunettesdiAPIBaseURL"`
	K8sDashboardAPIBaseURL string `json:"k8sDashboardAPIBaseURL"`
	DashboardBaseURL       string `json:"dashboardBaseURL"`
	DashboardDevBaseURL    string `json:"dashboardDevBaseURL"`
}

type WebkubectlConfig struct {
	HCSWebBaseURL         string `json:"hcsWebBaseURL"`
	HCSAPIBaseURL         string `json:"hcsAPIBaseURL"`
	IAMBaseURL            string `json:"iamBaseURL"`
	HCSAgentHelperBaseURL string `json:"hcsAgentHelperBaseURL"`
}

// hcs environment
type ENV string

const (
	ENVTest = "test"
	ENVProd = "prod"
)

// helper env
type HelperConfig struct {
	Env ENV `json:"env"`
}

// 附加链接相关配置，比如使用手册等
type LinksConfig struct {
	UserGuide string `json:"userGuide"`
}

var AppConfig *Config

// LoadConfig 加载配置
func LoadConfig() error {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	// 定义模型相关参数
	analysisFlag := flag.String("model_analysis", getEnvOrDefault("MODEL_ANALYSIS", "Qwen3-32B"), "Analysis model parameter")
	answerFlag := flag.String("model_answer", getEnvOrDefault("MODEL_ANSWER", "Qwen3-32B"), "Answer model parameter")
	modelUrlFlag := flag.String("model_url", getEnvOrDefault("MODEL_URL", "http://localhost:11434/api/chat"), "Model URL")
	tokenFlag := flag.String("model_token", getEnvOrDefault("MODEL_TOKEN", ""), "Model token")

	//HcpConfig
	hcpDevBaseUrl := flag.String("hcp_dev_url", getEnvOrDefault("HCP_DEV_URL", "https://hcp-dev.alipay.com"), "HCP dev URL")
	hcpDevAK := flag.String("hcp_dev_ak", getEnvOrDefault("HCP_DEV_AK", ""), "HCP dev AK")
	hcpDevSK := flag.String("hcp_dev_sk", getEnvOrDefault("HCP_DEV_SK", ""), "HCP dev SK")
	hcpProdBaseUrl := flag.String("hcp_prod_url", getEnvOrDefault("HCP_PROD_URL", "https://hcp.alipay.com"), "HCP prod URL")
	hcpProdAK := flag.String("hcp_prod_ak", getEnvOrDefault("HCP_PROD_AK", ""), "HCP prod AK")
	hcpProdSK := flag.String("hcp_prod_sk", getEnvOrDefault("HCP_PROD_SK", ""), "HCP prod SK")

	//YuQue 配置
	yuQueBaseUrl := flag.String("yuque_base_url", getEnvOrDefault("YUQUE_BASE_URL", "https://yuque.alipay.com"), "YuQue Base URL")
	yuQueToken := flag.String("yuque_token", getEnvOrDefault("YUQUE_TOKEN", ""), "YuQue Token")

	// 钉钉相关配置
	dingtalkClientID := flag.String("dingtalk_client_id", getEnvOrDefault("DINGTALK_CLIENT_ID", ""), "DingTalk Client ID")
	dingtalkClientSecret := flag.String("dingtalk_client_secret", getEnvOrDefault("DINGTALK_CLIENT_SECRET", ""), "DingTalk Client Secret")

	// 服务器API基础URL
	lunettesdiAPIBaseURL := flag.String("lunettesdi_api_base_url", getEnvOrDefault("LUNETTESDI_API_BASE_URL", "http://lunettesdi.hcs.svc.alipay.net:18883"), "Server API Base URL")
	k8sDashboardAPIBaseURL := flag.String("k8s_dashboard_api_base_url", getEnvOrDefault("K8S_DASHBOARD_API_BASE_URL", "http://istio-ingressgateway-sigma-staging.istio-system.svc.alipay.net:8000"), "k8s dashboard api server")
	// 服务器Dashboard URL
	dashboardBaseURL := flag.String("dashboard_base_url", getEnvOrDefault("LUNETTES_DASHBOARD_BASE_URL", "http://lunettes-dashboard.hcs.svc.alipay.net:30030"), "Dashboard URL")
	dashboardDevBaseURL := flag.String("dashboard_dev_base_url", getEnvOrDefault("LUNETTES_DASHBOARD_DEV_BASE_URL", "http://lunettes-dashboard.hcs-dev.svc.alipay.net:30030"), "Dashboard Dev URL")

	// 用户手册链接
	userGuide := flag.String("user_guide", getEnvOrDefault("USER_GUIDE", ""), "User guide link")

	dev := flag.Bool("dev", getBoolEnvOrDefault("DEVELOPMENT", false), "Development mode.")

	// helper env
	helperEnv := flag.String("helper_env", getEnvOrDefault("HELPER_ENV", "test"), "agent helper environment. [test, prod].")

	// webkubectl
	envTestHCSWebBaseURL := flag.String("test_hcs_web_base_url", getEnvOrDefault("ENV_TEST_HCS_WEB_BASE_URL", "http://hcs.test.alipay.net"), "HCS web base url of test environment.")
	envTestHCSAPIBaseURL := flag.String("test_hcs_api_base_url", getEnvOrDefault("ENV_TEST_HCS_API_BASE_URL", "http://istio-ingressgateway-hcs-sit.istio-system.svc.alipay.net:8080"), "HCS api base url of test environment.")
	envTestIAMBaseURL := flag.String("test_iam_base_url", getEnvOrDefault("ENV_TEST_IAM_BASE_URL", "http://iamconsole.test.alipay.net"), "IAM base url of test environment.")
	envTestHCSAgentHelperBaseURL := flag.String("test_hcsagent_helper_base_url", getEnvOrDefault("ENV_TEST_HCSAGENT_HELPER_BASE_URL", "http://hcs-agent-helper.hcs.svc.alipay.net:8000"), "HCS agent helper base url of test environment.")
	envProdHCSWebBaseURL := flag.String("prod_hcs_web_base_url", getEnvOrDefault("ENV_PROD_HCS_WEB_BASE_URL", "https://hcs.alipay.com"), "HCS web base url of prod environment.")
	envProdHCSAPIBaseURL := flag.String("prod_hcs_api_base_url", getEnvOrDefault("ENV_PROD_HCS_API_BASE_URL", "https://hcs-api.alipay.com"), "HCS api base url of prod environment.")
	envProdIAMBaseURL := flag.String("prod_iam_base_url", getEnvOrDefault("ENV_PROD_IAM_BASE_URL", "https://iam.alipay.com"), "IAM base url of prod environment.")
	envProdHCSAgentHelperBaseURL := flag.String("prod_hcsagent_helper_base_url", getEnvOrDefault("ENV_PROD_HCSAGENT_HELPER_BASE_URL", "https://hcs-agent-helper.alipay.com"), "HCS agent helper base url of test environment.")

	// clusters
	clustersJSON := flag.String("clusters", getEnvOrDefault("CLUSTERS", "{}"), `Clusters in JSON format, ex: {"test":[{"name":"sigma-staging"}],"prod":[{"name":"sigma-ea133"},{"name":"sigma-sa128"}]}`)

	//es
	esEndpoint := flag.String("es_endpoint", getEnvOrDefault("ES_ENDPOINT", "http://es-dev.alipay.net:9200"), "Elasticsearch endpoint")
	esUsername := flag.String("es_username", getEnvOrDefault("ES_USERNAME", ""), "Elasticsearch username")
	esPassword := flag.String("es_password", getEnvOrDefault("ES_PASSWORD", ""), "Elasticsearch password")

	// 从命令行解析参数
	flag.Parse()

	config := &Config{
		EsOptions: &common.ESOptions{
			EndPoint: *esEndpoint,
			Username: *esUsername,
			Password: *esPassword,
		},
		Hcp: HcpConfigs{
			DevConfig: HcpConfig{
				AK:  *hcpDevAK,
				SK:  *hcpDevSK,
				URL: *hcpDevBaseUrl,
			},
			ProdConfig: HcpConfig{
				AK:  *hcpProdAK,
				SK:  *hcpProdSK,
				URL: *hcpProdBaseUrl,
			},
		},
		YuQue: YuQueConfig{
			BaseUrl: *yuQueBaseUrl,
			Token:   *yuQueToken,
		},
		DingTalk: DingTalkConfig{
			ClientID:     *dingtalkClientID, // 使用命令行参数
			ClientSecret: *dingtalkClientSecret,
		},
		Model: ModelConfig{
			AnalysisModel: *analysisFlag,
			AnswerModel:   *answerFlag,
			URL:           *modelUrlFlag,
			Token:         *tokenFlag,
		},
		Lunettes: LunettesConfig{
			LunettesdiAPIBaseURL:   *lunettesdiAPIBaseURL,
			K8sDashboardAPIBaseURL: *k8sDashboardAPIBaseURL,
			DashboardBaseURL:       *dashboardBaseURL,
			DashboardDevBaseURL:    *dashboardDevBaseURL,
		},
		Links: LinksConfig{
			UserGuide: *userGuide,
		},
		Development: *dev,
		WebkubectlEnvTest: WebkubectlConfig{
			HCSWebBaseURL:         *envTestHCSWebBaseURL,
			HCSAPIBaseURL:         *envTestHCSAPIBaseURL,
			IAMBaseURL:            *envTestIAMBaseURL,
			HCSAgentHelperBaseURL: *envTestHCSAgentHelperBaseURL,
		},
		WebkubectlEnvProd: WebkubectlConfig{
			HCSWebBaseURL:         *envProdHCSWebBaseURL,
			HCSAPIBaseURL:         *envProdHCSAPIBaseURL,
			IAMBaseURL:            *envProdIAMBaseURL,
			HCSAgentHelperBaseURL: *envProdHCSAgentHelperBaseURL,
		},
		Helper: HelperConfig{
			Env: ENV(*helperEnv),
		},
	}

	err := json.Unmarshal([]byte(*clustersJSON), &config.Clusters)
	if err != nil {
		return fmt.Errorf("failed to parse clusters: %v", err)
	}

	// 验证必需的配置
	if config.DingTalk.ClientID == "" || config.DingTalk.ClientSecret == "" {
		return fmt.Errorf("DINGTALK_CLIENT_ID and DINGTALK_CLIENT_SECRET environment variables are required")
	}

	AppConfig = config
	slog.Info("Configuration loaded successfully")
	return nil
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnvOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		switch strings.ToLower(value) {
		case "true":
			return true
		case "false":
			return false
		}
	}
	return defaultValue
}
