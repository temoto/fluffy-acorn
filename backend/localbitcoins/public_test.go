package localbitcoins

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func stringHandler(s string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(s))
	})
}

func TestAdsBuyOnline1(t *testing.T) {
	const input = `{
    "data": {
        "ad_count": 1,
        "ad_list": [
            {
                "actions": {
                    "public_view": "https://localbitcoins.com/ad/357342"
                },
                "data": {
                    "ad_id": 357342,
                    "age_days_coefficient_limit": "0.00",
                    "atm_model": null,
                    "bank_name": "Альфа Cash-in",
                    "city": "",
                    "countrycode": "RU",
                    "created_at": "2016-07-17T12:01:54+00:00",
                    "currency": "RUB",
                    "email": null,
                    "first_time_limit_btc": null,
                    "hidden_by_opening_hours": false,
                    "is_local_office": false,
                    "lat": 0.0,
                    "location_string": "Russia",
                    "lon": 0.0,
                    "max_amount": "500000",
                    "max_amount_available": "487001",
                    "min_amount": "50000",
                    "msg": "Пополняете счет через Alfa Cash-in, моментально получаете биткойны\r\n",
                    "online_provider": "SPECIFIC_BANK",
                    "payment_window_minutes": 90,
                    "profile": {
                        "feedback_score": 100,
                        "last_online": "2016-09-24T15:34:40+00:00",
                        "name": "moneyman2k (100+; 100%)",
                        "trade_count": "100+",
                        "username": "moneyman2k"
                    },
                    "require_feedback_score": 99,
                    "require_identification": false,
                    "require_trade_volume": 0.0,
                    "require_trusted_by_advertiser": false,
                    "sms_verification_required": false,
                    "temp_price": "39033.82",
                    "temp_price_usd": "609.82",
                    "trade_type": "ONLINE_SELL",
                    "trusted_required": false,
                    "visible": true,
                    "volume_coefficient_btc": "1.50"
                }
            }
        ]
    },
    "pagination": {
        "next": "https://localbitcoins.com/buy-bitcoins-online/rub/.json?page=2"
    }
}`
	ts := httptest.NewServer(stringHandler(input))
	lb := new(Api)
	lb.Configure(ts.URL, "", "key", "secret")
	r, err := lb.PublicAdsBuyOnlineCurrency("rub")
	if err != nil {
		t.Fatal(err)
	}
	if r.Data.AdCount != 1 {
		t.Fatalf("invalid Data.AdCount: %v", r.Data.AdCount)
	}
	if r.Data.AdList[0].Data.MaxAmountAvailable != "487001" {
		t.Fatalf("invalid Data.AdList[0].Data.MaxAmountAvailable: %v", r.Data.AdList[0].Data.MaxAmountAvailable)
	}
	if r.Data.AdList[0].Data.Profile.Username != "moneyman2k" {
		t.Fatalf("invalid Data.AdList[0].Data.Profile.Username: %v", r.Data.AdList[0].Data.Profile.Username)
	}
}
