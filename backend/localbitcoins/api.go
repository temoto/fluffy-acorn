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

func (self *Api) Configure(url, key, secret string) {
	if url == "" {
		url = "https://localbitcoins.com"
	}
	self.urlPrefix = url
	self.key = key
	self.secret = secret
	tr := &http.Transport{}
	socks, err := proxy.SOCKS5("tcp", "127.0.0.1:1083", nil, proxy.Direct)
	if err != nil {
		panic(err)
	}
	tr.Dial = socks.Dial
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
	log.Printf("localbitcoins.RequestJson: url=%s", req.URL.String())

	resp, err := self.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	return err
}

func (self *Api) Ads() (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("GET", "/api/ads/", "", false, &r)
}

func (self *Api) AdGet([]int) (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("GET", "/api/ad-get/", "", false, &r)
}

func (self *Api) AdEdit(id int) (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("POST", fmt.Sprintf("/api/ad/%d/", id), "", false, &r)
}

func (self *Api) AdCreate() (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("POST", "/api/ad-create/", "", false, &r)
}

func (self *Api) AdDelete(id int) (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("POST", fmt.Sprintf("/api/ad-delete/%d/", id), "", false, &r)
}

func (self *Api) Dashboard() (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("GET", "/api/dashboard/", "", false, &r)
}

func (self *Api) ContactMessages(contactId int) (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("GET", fmt.Sprintf("/api/contact_messages/%d/", contactId), "", false, &r)
}

func (self *Api) RecentMessages() (interface{}, error) {
	type RT struct {
		Data struct {
			MessageCount  int       `json:"message_count"`
			MessageList   []Message `json:"message_list"`
			CheckInterval float32   `json:"check_interval"`
		} `json:"data"`
	}
	var r RT
	// r.Data.CheckInterval = 34.1
	return &r.Data, self.RequestJson("GET", "/api/recent_messages/", "", false, &r)
}

func (self *Api) ContactMessagePost(contactId int, msg string) (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("POST", fmt.Sprintf("/api/contact_message_post/%d/", contactId), "", false, &r)
}

func (self *Api) Feedback(username, feedback string) (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("POST", fmt.Sprintf("/api/feedback/%s/", username), "", false, &r)
}

func (self *Api) PublicAdsOnlineBuy(currency string) (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("GET", fmt.Sprintf("/buy-bitcoins-online/%s/.json", currency), "", true, &r)
}
