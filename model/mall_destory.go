package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type FailedOrderResp struct {
	CreatedAt    time.Time       `json:"created_at"`
	SubSystem    string          `json:"sub_system"`
	OrderID      string          `json:"order_id"`
	OriginToken  string          `json:"origin_token"`
	OriginAmount decimal.Decimal `json:"origin_amount"`
	DestroyToken string          `json:"destroy_token"`
	OpenID       string          `json:"open_id"`
	ShopOpenID   string          `json:"shop_open_id"`
	State        string          `json:"state"`
	DestroyState string          `json:"destroy_state"`
	ErrCode      string          `json:"err_code"`
	ErrDesc      string          `json:"err_desc"`
}
