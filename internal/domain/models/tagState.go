package models

import "time"

type TagQuality string

const (
	TagQualityGood  TagQuality = "GOOD"
	TagQualityStale TagQuality = "STALE"
	TagQualityBad   TagQuality = "BAD"
)

type TagState struct {
	Value              interface{} `json:"value"`
	Quality            TagQuality  `json:"quality"`
	LastSuccessfulRead time.Time   `json:"last_successful_read"`
}
