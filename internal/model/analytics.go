package model

import "time"

const (
	GroupByDay       = "day"
	GroupByMonth     = "month"
	GroupByUserAgent = "useragent"
)

type AnalyticsRequest struct { // используется для получения от клиента фильтров аналитики
	Shortkey string
	Start    *time.Time `json:"start"`
	End      *time.Time `json:"end"`
	GroupBy  string     `json:"groupby"` // будут 3 значения: useragent/day/month
}

// Для отдачи результата аггрегации аналитики:

type AnalyticsResponse struct {
	GroupBy string `json:"groupby"`
	Total   int    `json:"total"`
	List    []Item `json:"items"`
}

type Item struct {
	Line  string `json:"line"`
	Count int    `json:"count"`
}
