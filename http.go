package mob_push_sdk

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var gHTTPClient *HTTPClient

func GetHTTPClient(proxyUrl *string) *HTTPClient {
	if gHTTPClient == nil {
		if proxyUrl == nil {
			gHTTPClient = CreateHTTPClient()
		} else {
			gHTTPClient = CreateHTTPClientProxy(*proxyUrl)
		}
	}
	return gHTTPClient
}

type HTTPClient struct {
	client    *http.Client
	Sign      string
	Appkey    string
	AppSecret string
}

func CreateHTTPClient() *HTTPClient {
	return &HTTPClient{client: &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}}
}

func CreateHTTPClientProxy(proxy string) *HTTPClient {
	p, _ := url.Parse(proxy)
	return &HTTPClient{client: &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(p),
		},
	}}
}

func (client *HTTPClient) Request(request *http.Request) ([]byte, error) {
	response, err := client.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	return ioutil.ReadAll(response.Body)
}

func (client *HTTPClient) Get(uri string) ([]byte, error) {
	if request, err := http.NewRequest(http.MethodGet, uri, nil); err == nil {
		return client.Request(request)
	} else {
		return nil, err
	}
}

func (client *HTTPClient) PostJSON(pushClient *PushClient, uri string, values interface{}) ([]byte, error) {

	var body []byte

	if requestData, err := json.Marshal(values); err != nil {
		return nil, err
	} else {
		body = requestData
	}
	client.SetSign(pushClient, body)

	//fmt.Println(uri,string(body))

	if request, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(body)); err == nil {
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("key", client.Appkey)
		request.Header.Add("sign", client.Sign)
		return client.Request(request)
	} else {
		return nil, err
	}
}

func (client *HTTPClient) SetSign(pushClient *PushClient, requestJson []byte) *HTTPClient {

	client.Appkey = pushClient.AppKey
	client.AppSecret = pushClient.AppSecert

	originString := string(requestJson) + client.AppSecret
	client.Sign = getMD5Encode(originString)
	return client
}

func getMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
