package jobs

import (
	"middleware/internal/domain/models"
	"sync"
	"time"
)

var (
	StatusClpRealTimeSync = sync.Map{}
	MutexFila             sync.Mutex
	currentTags           = tagStateCache{byCLP: make(map[uint]map[uint]models.TagState)}
)

type tagStateCache struct {
	mu    sync.RWMutex
	byCLP map[uint]map[uint]models.TagState
}

// ApplyPoll atomically publishes one polling cycle. Values contains only tags
// successfully decoded in this cycle; missing tags retain their last valid
// value and timestamp and become STALE (or BAD if they have never succeeded).
func ApplyPoll(clpID uint, tagIDs []uint, values map[uint]interface{}, readAt time.Time) {
	currentTags.mu.Lock()
	defer currentTags.mu.Unlock()

	tags := currentTags.byCLP[clpID]
	if tags == nil {
		tags = make(map[uint]models.TagState, len(tagIDs))
		currentTags.byCLP[clpID] = tags
	}

	active := make(map[uint]struct{}, len(tagIDs))
	for _, tagID := range tagIDs {
		active[tagID] = struct{}{}
		if value, ok := values[tagID]; ok {
			tags[tagID] = models.TagState{
				Value:              value,
				Quality:            models.TagQualityGood,
				LastSuccessfulRead: readAt,
			}
			continue
		}

		state := tags[tagID]
		if state.LastSuccessfulRead.IsZero() {
			state.Quality = models.TagQualityBad
		} else {
			state.Quality = models.TagQualityStale
		}
		tags[tagID] = state
	}

	for tagID := range tags {
		if _, ok := active[tagID]; !ok {
			delete(tags, tagID)
		}
	}
}

// MarkCLPUnavailable changes quality without replacing any last valid value.
func MarkCLPUnavailable(clpID uint) {
	currentTags.mu.Lock()
	defer currentTags.mu.Unlock()

	for tagID, state := range currentTags.byCLP[clpID] {
		if state.LastSuccessfulRead.IsZero() {
			state.Quality = models.TagQualityBad
		} else {
			state.Quality = models.TagQualityStale
		}
		currentTags.byCLP[clpID][tagID] = state
	}
}

func DeleteCLP(clpID uint) {
	currentTags.mu.Lock()
	delete(currentTags.byCLP, clpID)
	currentTags.mu.Unlock()
}

func ReadAllTagsRealTime() (map[uint]map[uint]models.TagState, error) {
	currentTags.mu.RLock()
	defer currentTags.mu.RUnlock()

	result := make(map[uint]map[uint]models.TagState, len(currentTags.byCLP))
	for clpID, tags := range currentTags.byCLP {
		clpMap := make(map[uint]models.TagState, len(tags))
		for tagID, state := range tags {
			clpMap[tagID] = state
		}
		result[clpID] = clpMap
	}
	return result, nil
}

func ReadAllClpsStatus() map[uint]bool {
	result := make(map[uint]bool)

	StatusClpRealTimeSync.Range(func(key, value interface{}) bool {
		clpID, okID := key.(uint)
		status, okStatus := value.(bool)
		if okID && okStatus {
			result[clpID] = status
		}
		return true
	})

	return result
}
