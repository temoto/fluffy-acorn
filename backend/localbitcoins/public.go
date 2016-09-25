package localbitcoins

import (
	"fmt"
	"time"
)

// {
//     "actions": {
//         "public_view": "https://localbitcoins.com/ad/357342"
//     },
//     "data": {
//         "ad_id": 357342,
//         "age_days_coefficient_limit": "0.00",
//         "atm_model": null,
//         "bank_name": "Альфа Cash-in",
//         "city": "",
//         "countrycode": "RU",
//         "created_at": "2016-07-17T12:01:54+00:00",
//         "currency": "RUB",
//         "email": null,
//         "first_time_limit_btc": null,
//         "hidden_by_opening_hours": false,
//         "is_local_office": false,
//         "lat": 0.0,
//         "location_string": "Russia",
//         "lon": 0.0,
//         "max_amount": "500000",
//         "max_amount_available": "487001",
//         "min_amount": "50000",
//         "msg": "Пополняете счет через Alfa Cash-in, моментально получаете биткойны\r\n",
//         "online_provider": "SPECIFIC_BANK",
//         "payment_window_minutes": 90,
//         "profile": {
//             "feedback_score": 100,
//             "last_online": "2016-09-24T15:34:40+00:00",
//             "name": "moneyman2k (100+; 100%)",
//             "trade_count": "100+",
//             "username": "moneyman2k"
//         },
//         "require_feedback_score": 99,
//         "require_identification": false,
//         "require_trade_volume": 0.0,
//         "require_trusted_by_advertiser": false,
//         "sms_verification_required": false,
//         "temp_price": "39033.82",
//         "temp_price_usd": "609.82",
//         "trade_type": "ONLINE_SELL",
//         "trusted_required": false,
//         "visible": true,
//         "volume_coefficient_btc": "1.50"
//     }
// },
type PublicAd struct {
	Actions struct {
		PublicView string `json:"public_view"`
	} `json:"actions"`
	Data struct {
		AdId                    int    `json:"ad_id"`
		AgeDaysCoefficientLimit string `json:"age_days_coefficient_limit"` // float
		// AtmModel ?
		BankName                   string    `json:"bank_name"`
		City                       string    `json:"city"`
		CountryCode                string    `json:"countrycode"`
		CreatedAtString            string    `json:"created_at"`
		CreatedAt                  time.Time `json:"-"`
		Currency                   string    `json:"currency"`
		Lat                        float32   `json:"lat"`
		Lon                        float32   `json:"lon"`
		MaxAmount                  string    `json:"max_amount"`           // int
		MaxAmountAvailable         string    `json:"max_amount_available"` // int
		MinAmount                  string    `json:"min_amount"`           // int
		Msg                        string    `json:"msg"`
		OnlineProvider             string    `json:"online_provider"`
		PaymentWindowMinutes       int       `json:"payment_window_minutes"`
		Profile                    ProfileT  `json:"profile"`
		RequireFeedbackScore       int       `json:"require_feedback_score"`
		RequireIdentification      bool      `json:"require_identification"`
		RequireTradeVolume         float32   `json:"require_trade_volume"`
		RequireTrustedByAdvertiser bool      `json:"require_trusted_by_advertiser"`
		SmsVerificationRequired    bool      `json:"sms_verification_required"`
		TempPrice                  string    `json:"temp_price"`     // float
		TempPriceUsd               string    `json:"temp_price_usd"` // float
		TradeType                  string    `json:"trade_type"`
		TrustedRequired            bool      `json:"trusted_required"`
		Visible                    bool      `json:"visible"`
		VolumeCoefficientBtc       string    `json:"volume_coefficient_btc"` // float
	} `json:"data"`
}
type PublicAdsBuyOnlineResponse struct {
	Data struct {
		AdCount int        `json:"ad_count"`
		AdList  []PublicAd `json:"ad_list"`
	} `json:"data"`
	Pagination struct {
		Next string `json:"next"`
	} `json:"pagination"`
}

func (self *Api) PublicAdsBuyOnlineCurrency(currency string) (*PublicAdsBuyOnlineResponse, error) {
	r := new(PublicAdsBuyOnlineResponse)
	err := self.RequestJson("GET", fmt.Sprintf("/buy-bitcoins-online/%s/.json", currency), "", true, r)
	return r, err
}
