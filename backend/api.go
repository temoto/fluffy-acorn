package main

import (
	"github.com/temoto/fluffy-acorn/backend/localbitcoins"
	"net/http"
	"os"
)

const lbKey = "EDIT_THIS"
const lbSecret = "EDIT_THIS"

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
	lbApi.Configure(lbUrlPrefix, lbKey, lbSecret)
}

func wrapLbApi(inner func() (interface{}, error)) apiFunc {
	return func(req *http.Request, params interface{}) (interface{}, error) {
		return inner()
	}
}

func init() {
	configureLbApiClient()
	routes = append(routes,
		Route{"api-dashboard", "GET", "/api/dashboard", wrapJson(func(req *http.Request, params interface{}) (interface{}, error) { return lbApi.Dashboard() })},
		Route{"api-recent-messages", "GET", "/api/recent-messages", wrapJson(func(req *http.Request, params interface{}) (interface{}, error) { return lbApi.RecentMessages() })},
	)
}
