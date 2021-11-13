package billing

type Config struct {
	AppBaseUrl string
	Debug      bool
	Env        string
	billingUrl string
}

var billingConfig Config
var billingUrlMap = map[string]string{
	"sandbox": "http://127.0.0.1:1180/v1",
	"prod":    "https://pay.akouendy.com/v1",
}

// Init to configure the sdk
func Init(config Config) {
	billingUrl, ok := billingUrlMap[config.Env]
	if !ok {
		config.Env = "sandbox"
		billingUrl = billingUrlMap["sandbox"]
	}
	// set billing url
	config.billingUrl = billingUrl
	billingConfig = config
}
