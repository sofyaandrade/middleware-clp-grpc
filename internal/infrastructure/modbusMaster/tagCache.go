package modbusmaster

import (
	"middleware/internal/domain/models"
	"middleware/internal/infrastructure/jobs"
	"time"
)

func initializeCLPTagCache(clp models.CLP) {
	tagIDs := make([]uint, 0, len(clp.Tags))
	for _, tag := range clp.Tags {
		if tag != nil {
			tagIDs = append(tagIDs, tag.ID)
		}
	}
	jobs.ApplyPoll(clp.ID, tagIDs, nil, time.Now())
}
