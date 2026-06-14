package jobs

import "sync"

var (
	TagsByClpMaster       = make(map[uint]map[uint]interface{})
	filaTags              []uint
	mapaTags              = make(map[uint]float32)
	StatusClpRealTimeSync = sync.Map{}
	MutexFila             sync.Mutex
	MutexMaster           sync.Mutex
	TagsByIDMaster        sync.Map
	TagsCache             sync.Map
)

func ReadAllTagsRealTime() (map[uint]map[uint]interface{}, error) {
	result := make(map[uint]map[uint]interface{})

	MutexMaster.Lock()
	for clpID, tags := range TagsByClpMaster {
		clpMap := make(map[uint]interface{})
		for tagID := range tags {
			if valor, ok := TagsByIDMaster.Load(tagID); ok {
				clpMap[tagID] = valor
			}
		}
		thisValueExists(result, clpID, clpMap)
	}
	MutexMaster.Unlock()

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

func thisValueExists(result map[uint]map[uint]interface{}, clpID uint, clpMap map[uint]interface{}) {
	if exist, ok := result[clpID]; ok {
		for k, v := range clpMap {
			exist[k] = v
		}
	} else {
		result[clpID] = clpMap
	}
}
