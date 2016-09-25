package main

import (
	"github.com/temoto/fluffy-acorn/backend/localbitcoins"
	"log"
	"net/http"
	"os"
	"strconv"
)

const lbKey = "secret"
const lbSecret = "secret"
const lbProxy = "127.0.0.1:1083"

var lbApi *localbitcoins.Api

func configureLbApiClient() {
	if lbApi != nil {
		return
	}
	lbApi = new(localbitcoins.Api)
	lbUrlPrefix := ""
	if os.Getenv("acorn_lb_mock") != "" {
		lbUrlPrefix = AcornUrlPrefix + "/lbapi-mock"
	}
	lbApi.Configure(lbUrlPrefix, lbProxy, lbKey, lbSecret)
}

func wrapLbApi(inner func() (interface{}, error)) apiFunc {
	return func(req *http.Request, params interface{}) (interface{}, error) {
		return inner()
	}
}

func parseInt(s string, def int) int {
	x, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return int(x)
}

func init() {
	configureLbApiClient()
	routes = append(routes,
		Route{"api-dashboard", "GET", "/api/dashboard", wrapJson(func(req *http.Request, params interface{}) (interface{}, error) { return lbApi.Dashboard() })},
		Route{"api-recent-messages", "GET", "/api/recent-messages", wrapJson(func(req *http.Request, params interface{}) (interface{}, error) {
			queryLastSeen := req.URL.Query().Get("last_seen")
			queryLimit := parseInt(req.URL.Query().Get("limit"), 50)
			rm, err := lbApi.RecentMessages()
			if err != nil {
				return nil, err
			}
			log.Printf("backend api-recent-messages: queryLastSeen=%s queryLimit=%d", queryLastSeen, queryLimit)
			if queryLimit > 0 {
				rm.Data.MessageList = rm.Data.MessageList[:queryLimit]
			}
			if queryLastSeen != "" {
				unseen := 0
				for _, m := range rm.Data.MessageList {
					if queryLastSeen < m.CreatedAtString {
						unseen++
					} else {
						break
					}
				}
				rm.Data.MessageList = rm.Data.MessageList[:unseen]
			}
			return rm, err
		})},
		Route{"api-online-ads-buy", "GET", "/api/online-ads-buy", wrapJson(func(req *http.Request, params interface{}) (interface{}, error) {
			return lbApi.PublicAdsBuyOnlineCurrency(req.URL.Query().Get("currency"))
		})},
	)
}
