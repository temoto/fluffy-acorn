package localbitcoins

import (
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
	Sender          struct {
		FeedbackScore    int       `json:"feedback_score"`
		Name             string    `json:"name"`
		LastOnlineString string    `json:"last_online"`
		LastOnline       time.Time `json:"-"`
		TradeCount       string    `json:"trade_count"`
		Username         string    `json:"username"`
	} `json:"sender"`
}
