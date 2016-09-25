package localbitcoins

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/net/proxy"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type mapstrany map[string]interface{}

type Api struct {
	c         *http.Client
	urlPrefix string
	key       string
	secret    string
}

func (self *Api) Configure(url, socksAddr, key, secret string) {
	if url == "" {
		url = "https://localbitcoins.com"
	}
	self.urlPrefix = url
	self.key = key
	self.secret = secret
	tr := &http.Transport{
		ResponseHeaderTimeout: 41 * time.Second,
	}
	if socksAddr != "" {
		socks, err := proxy.SOCKS5("tcp", socksAddr, nil, proxy.Direct)
		if err != nil {
			panic(err)
		}
		tr.Dial = socks.Dial
	}
	self.c = &http.Client{Transport: tr}
}

func (self *Api) sign(req *http.Request, params string) {
	nonce := time.Now().Unix()
	nonceString := strconv.FormatInt(nonce, 10)
	signMessage := nonceString + self.key + req.URL.Path + params
	h := hmac.New(sha256.New, []byte(self.secret))
	h.Write([]byte(signMessage))
	signature := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	log.Printf("localbitcoins.sign: message='%s' signature=%s", signMessage, signature)
	req.Header.Set("Apiauth-Key", self.key)
	req.Header.Set("Apiauth-Nonce", nonceString)
	req.Header.Set("Apiauth-Signature", signature)
}

func (self *Api) MakeRequest(method, path, params string, public bool) (*http.Request, error) {
	req, err := http.NewRequest(method, self.urlPrefix+path, strings.NewReader(params))
	if err != nil {
		return nil, err
	}
	if !public {
		self.sign(req, params)
	}
	return req, err
}

func (self *Api) RequestJson(method, path, params string, public bool, result interface{}) error {
	req, err := self.MakeRequest(method, path, params, public)
	if err != nil {
		return err
	}
	log.Printf("localbitcoins.RequestJson: method=%s path=%s begin url=%s", method, path, req.URL.String())

	t1 := time.Now()
	resp, err := self.c.Do(req)
	td := time.Since(t1)
	log.Printf("localbitcoins.RequestJson: method=%s path=%s time=%s error=%v", method, path, td, err)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	return err
}

func (self *Api) Dashboard() (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("GET", "/api/dashboard/", "", false, &r)
}

func (self *Api) Feedback(username, feedback string) (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("POST", fmt.Sprintf("/api/feedback/%s/", username), "", false, &r)
}

// TODO sort this out
type ProfileT struct {
	FeedbackScore    int       `json:"feedback_score"`
	Name             string    `json:"name"`
	LastOnlineString string    `json:"last_online"`
	LastOnline       time.Time `json:"-"`
	TradeCount       string    `json:"trade_count"`
	Username         string    `json:"username"`
}
