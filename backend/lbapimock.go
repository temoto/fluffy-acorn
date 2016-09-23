package main

import (
	"net/http"
)

type mapstrany map[string]interface{}

func lbapiMockDashboard(req *http.Request, params interface{}) (interface{}, error) {
	r := struct {
		Data mapstrany `json:"data"`
	}{
		Data: mapstrany{"contact_list": []string{"one", "two"}},
	}
	return r, nil
}
