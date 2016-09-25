package localbitcoins

import (
	"fmt"
)

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
