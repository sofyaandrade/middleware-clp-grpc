package jobs

import (
	"middleware/internal/domain/models"
	"sync"
	"testing"
	"time"
)

func resetTagStateCache() {
	currentTags.mu.Lock()
	currentTags.byCLP = make(map[uint]map[uint]models.TagState)
	currentTags.mu.Unlock()
}

func TestFailedPollRetainsLastValidValueAndMarksItStale(t *testing.T) {
	resetTagStateCache()
	t.Cleanup(resetTagStateCache)

	readAt := time.Unix(123, 456)
	ApplyPoll(1, []uint{10}, map[uint]interface{}{10: int16(42)}, readAt)
	ApplyPoll(1, []uint{10}, nil, readAt.Add(time.Second))

	snapshot, err := ReadAllTagsRealTime()
	if err != nil {
		t.Fatal(err)
	}
	state := snapshot[1][10]
	if state.Value != int16(42) {
		t.Fatalf("value = %v, want last valid value 42", state.Value)
	}
	if state.Quality != models.TagQualityStale {
		t.Fatalf("quality = %q, want STALE", state.Quality)
	}
	if !state.LastSuccessfulRead.Equal(readAt) {
		t.Fatalf("last successful read = %v, want %v", state.LastSuccessfulRead, readAt)
	}
}

func TestFailedFirstPollMarksTagBad(t *testing.T) {
	resetTagStateCache()
	t.Cleanup(resetTagStateCache)

	ApplyPoll(1, []uint{10}, nil, time.Now())
	result, _ := ReadAllTagsRealTime()
	state := result[1][10]
	if state.Quality != models.TagQualityBad || state.Value != nil || !state.LastSuccessfulRead.IsZero() {
		t.Fatalf("state = %#v, want BAD without a fabricated value or timestamp", state)
	}
}

func TestTagStateCacheConcurrentPollAndSnapshot(t *testing.T) {
	resetTagStateCache()
	t.Cleanup(resetTagStateCache)

	const iterations = 1000
	var wg sync.WaitGroup
	for worker := uint(1); worker <= 8; worker++ {
		worker := worker
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				ApplyPoll(worker, []uint{worker}, map[uint]interface{}{worker: i}, time.Unix(0, int64(i+1)))
			}
		}()
	}
	for reader := 0; reader < 4; reader++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				if _, err := ReadAllTagsRealTime(); err != nil {
					t.Errorf("ReadAllTagsRealTime() error = %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()
}
