package models

import "time"

type TagQuality string

const (
	TagQualityGood  TagQuality = "GOOD"
	TagQualityStale TagQuality = "STALE"
	TagQualityBad   TagQuality = "BAD"
)

// TagState is an immutable snapshot of the latest known state of a tag.
// LastSuccessfulRead is not changed when a poll fails.
type TagState struct {
	Value              interface{} `json:"value"`
	Quality            TagQuality  `json:"quality"`
	LastSuccessfulRead time.Time   `json:"last_successful_read"`
}
