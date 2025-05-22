package entity

import "google.golang.org/genproto/googleapis/type/decimal"

type ReviewStats struct {
	TotalReviews     int64           `json:"total_reviews"`
	AvgOverall       decimal.Decimal `json:"avg_overall"`
	AvgLocation      decimal.Decimal `json:"avg_location"`
	AvgValue         decimal.Decimal `json:"avg_value"`
	AvgAccuracy      decimal.Decimal `json:"avg_accuracy"`
	AvgCommunication decimal.Decimal `json:"avg_communication"`
	AvgCleanliness   decimal.Decimal `json:"avg_cleanliness"`
	AvgCheckIn       decimal.Decimal `json:"avg_check_in"`
}

type ViewStats struct {
	TotalViews  int64 `json:"total_views"`
	UniqueViews int64 `json:"unique_views"`
}
