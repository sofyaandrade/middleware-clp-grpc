package grpcserver

import (
	"context"
	"fmt"
	realtimev1 "middleware/api/realtime/v1"
	"middleware/internal/domain/models"
	"middleware/internal/infrastructure/jobs"
	"reflect"
	"sort"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const defaultSnapshotInterval = time.Second

type snapshotReader func() (map[uint]map[uint]models.TagState, error)

type RealtimeTagServer struct {
	realtimev1.UnimplementedRealtimeTagServiceServer
	readSnapshot snapshotReader
	interval     time.Duration
}

func NewRealtimeTagServer() *RealtimeTagServer {
	return newRealtimeTagServer(jobs.ReadAllTagsRealTime, defaultSnapshotInterval)
}

func newRealtimeTagServer(reader snapshotReader, interval time.Duration) *RealtimeTagServer {
	return &RealtimeTagServer{readSnapshot: reader, interval: interval}
}

func (s *RealtimeTagServer) StreamTagValues(_ *realtimev1.StreamTagValuesRequest, stream realtimev1.RealtimeTagService_StreamTagValuesServer) error {
	lastSent := make(map[tagKey]tagFingerprint)
	if err := s.sendChanges(stream.Context(), stream.Send, lastSent); err != nil {
		return err
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			if err := s.sendChanges(stream.Context(), stream.Send, lastSent); err != nil {
				return err
			}
		}
	}
}

type tagKey struct {
	equipmentID uint
	tagID       uint
}

type tagFingerprint struct {
	value   interface{}
	quality models.TagQuality
}

func (s *RealtimeTagServer) sendChanges(ctx context.Context, send func(*realtimev1.TagValue) error, lastSent map[tagKey]tagFingerprint) error {
	snapshot, err := s.readSnapshot()
	if err != nil {
		return status.Errorf(codes.Unavailable, "falha ao ler cache de tags: %v", err)
	}

	for _, equipmentID := range sortedKeys(snapshot) {
		for _, tagID := range sortedKeys(snapshot[equipmentID]) {
			if ctx.Err() != nil {
				return nil
			}
			state := snapshot[equipmentID][tagID]
			key := tagKey{equipmentID: equipmentID, tagID: tagID}
			fingerprint := tagFingerprint{value: state.Value, quality: state.Quality}
			if previous, exists := lastSent[key]; exists && reflect.DeepEqual(previous, fingerprint) {
				continue
			}
			message, err := newTagValue(equipmentID, tagID, state)
			if err != nil {
				return status.Error(codes.Internal, err.Error())
			}
			if err := send(message); err != nil {
				return err
			}
			lastSent[key] = fingerprint
		}
	}
	return nil
}

type unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64
}

func sortedKeys[V any, K unsigned](values map[K]V) []K {
	keys := make([]K, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func newTagValue(equipmentID, tagID uint, state models.TagState) (*realtimev1.TagValue, error) {
	message := &realtimev1.TagValue{
		EquipmentId:                uint64(equipmentID),
		TagId:                      uint64(tagID),
		Quality:                    protoQuality(state.Quality),
		LastSuccessfulReadUnixNano: state.LastSuccessfulRead.UnixNano(),
	}
	if state.LastSuccessfulRead.IsZero() {
		message.LastSuccessfulReadUnixNano = 0
	}
	value := state.Value
	if value == nil {
		return message, nil
	}

	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Bool:
		message.Value = &realtimev1.TagValue_BoolValue{BoolValue: reflected.Bool()}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		message.Value = &realtimev1.TagValue_IntValue{IntValue: reflected.Int()}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		message.Value = &realtimev1.TagValue_UintValue{UintValue: reflected.Uint()}
	case reflect.Float32:
		message.Value = &realtimev1.TagValue_FloatValue{FloatValue: float32(reflected.Float())}
	case reflect.Float64:
		message.Value = &realtimev1.TagValue_DoubleValue{DoubleValue: reflected.Float()}
	case reflect.String:
		message.Value = &realtimev1.TagValue_StringValue{StringValue: reflected.String()}
	default:
		return nil, fmt.Errorf("tipo de valor nao suportado para equipamento %d, tag %d: %T", equipmentID, tagID, value)
	}
	return message, nil
}

func protoQuality(quality models.TagQuality) realtimev1.TagQuality {
	switch quality {
	case models.TagQualityGood:
		return realtimev1.TagQuality_TAG_QUALITY_GOOD
	case models.TagQualityStale:
		return realtimev1.TagQuality_TAG_QUALITY_STALE
	case models.TagQualityBad:
		return realtimev1.TagQuality_TAG_QUALITY_BAD
	default:
		return realtimev1.TagQuality_TAG_QUALITY_UNSPECIFIED
	}
}
