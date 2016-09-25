package localbitcoins

import (
	"fmt"
	"time"
)

// {
//     "contact_id": 6817778.0,
//     "created_at": "2016-09-20T15:03:39+00:00",
//     "is_admin": false,
//     "msg": "ок! удачи)",
//     "sender": {
//         "feedback_score": 100,
//         "last_online": "2016-09-20T15:15:29+00:00",
//         "name": "rangerka (18; 100%)",
//         "trade_count": "18",
//         "username": "rangerka"
//     }
// }
type Message struct {
	ContactId       float32   `json:"contact_id"`
	CreatedAtString string    `json:"created_at"`
	CreatedAt       time.Time `json:"-"`
	IsAdmin         bool      `json:"is_admin"`
	Msg             string    `json:"msg"`
	Sender          ProfileT  `json:"sender"`
}

func (self *Api) ContactMessages(contactId int) (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("GET", fmt.Sprintf("/api/contact_messages/%d/", contactId), "", false, &r)
}

type RecentMessagesResponse struct {
	Data struct {
		MessageList []Message `json:"message_list"`
	} `json:"data"`
	RetryAfter int `json:"retry_after"`
}

func (self *Api) RecentMessages() (*RecentMessagesResponse, error) {
	r := new(RecentMessagesResponse)
	// r.CheckInterval = 34.1
	err := self.RequestJson("GET", "/api/recent_messages/", "", false, r)
	return r, err
}

func (self *Api) ContactMessagePost(contactId int, msg string) (interface{}, error) {
	var r interface{}
	return r, self.RequestJson("POST", fmt.Sprintf("/api/contact_message_post/%d/", contactId), "", false, &r)
}
