package controllers

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statInputRequestCount = stats.Int64(
		"input_requests_count",
		"How many requests happened",
		stats.UnitDimensionless,
	)
	statSuccessedRequestCount = stats.Int64(
		"successed_requests_count",
		"How many requests successed",
		stats.UnitDimensionless,
	)
	statFailedRequestCount = stats.Int64(
		"failed_requests_count",
		"How many requests failed",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statInputRequestCount.Name(),
			Description: statInputRequestCount.Description(),
			Measure:     statInputRequestCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statSuccessedRequestCount.Name(),
			Description: statSuccessedRequestCount.Description(),
			Measure:     statSuccessedRequestCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statFailedRequestCount.Name(),
			Description: statFailedRequestCount.Description(),
			Measure:     statFailedRequestCount,
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}
