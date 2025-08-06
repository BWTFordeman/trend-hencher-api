package models

type IntradayData struct {
	Timestamp int64   `json:"timestamp"`
	GmtOffset int     `json:"gmtoffset"`
	Datetime  string  `json:"datetime"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int     `json:"volume"`
}
