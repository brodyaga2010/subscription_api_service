package model

import "time"

type SumRequest struct {
	UserID      string `query:"user_id"`
	ServiceName string `query:"service_name"`
	From        string `query:"from"`
	To          string `query:"to"`
	FromTime    time.Time
	ToTime      time.Time
}

type SumResponse struct {
	Total int `json:"total"`
}
