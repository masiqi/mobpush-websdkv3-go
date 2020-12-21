package mob_push_sdk

const (
	BASE_URL = "https://api.push.mob.com/"
)

type PushClient struct {
	AppKey    string
	AppSecert string
	BaseUrl   string
	ProxyUrl  string
}

func NewPushClient(appKey, appSecret, proxyUrl string) *PushClient {
	return &PushClient{appKey, appSecret, BASE_URL, proxyUrl}
}
