package grpcserver

import (
	realtimev1 "middleware/api/realtime/v1"
	"middleware/internal/domain/models"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

func TestSendChangesSendsInitialSnapshotAndOnlyChangedValues(t *testing.T) {
	snapshot := map[uint]map[uint]models.TagState{
		2: {20: {Value: true, Quality: models.TagQualityGood}},
		1: {11: {Value: float32(12.5), Quality: models.TagQualityGood}, 10: {Value: int16(-7), Quality: models.TagQualityGood}},
	}
	server := newRealtimeTagServer(func() (map[uint]map[uint]models.TagState, error) {
		return snapshot, nil
	}, defaultSnapshotInterval)
	var messages []*realtimev1.TagValue
	send := func(message *realtimev1.TagValue) error {
		messages = append(messages, message)
		return nil
	}
	lastSent := make(map[tagKey]tagFingerprint)

	if err := server.sendChanges(t.Context(), send, lastSent); err != nil {
		t.Fatalf("sendChanges() error = %v", err)
	}
	if len(messages) != 3 {
		t.Fatalf("initial message count = %d, want 3", len(messages))
	}
	if messages[0].GetEquipmentId() != 1 || messages[0].GetTagId() != 10 || messages[0].GetIntValue() != -7 {
		t.Fatalf("first message = %v, want equipment 1 tag 10 value -7", messages[0])
	}
	if messages[1].GetFloatValue() != 12.5 {
		t.Fatalf("second message float value = %v, want 12.5", messages[1].GetFloatValue())
	}
	if !messages[2].GetBoolValue() {
		t.Fatal("third message bool value = false, want true")
	}

	if err := server.sendChanges(t.Context(), send, lastSent); err != nil {
		t.Fatalf("unchanged sendChanges() error = %v", err)
	}
	if len(messages) != 3 {
		t.Fatalf("message count after unchanged snapshot = %d, want 3", len(messages))
	}

	snapshot[1][11] = models.TagState{Value: float32(13.75), Quality: models.TagQualityGood}
	if err := server.sendChanges(t.Context(), send, lastSent); err != nil {
		t.Fatalf("changed sendChanges() error = %v", err)
	}
	if len(messages) != 4 || messages[3].GetEquipmentId() != 1 || messages[3].GetTagId() != 11 || messages[3].GetFloatValue() != 13.75 {
		t.Fatalf("changed message = %v, want equipment 1 tag 11 value 13.75", messages[3])
	}
}

func TestNewTagValueRejectsUnsupportedValue(t *testing.T) {
	_, err := newTagValue(1, 2, models.TagState{Value: struct{}{}, Quality: models.TagQualityGood})
	if err == nil {
		t.Fatal("newTagValue() error = nil, want unsupported type error")
	}
}

func TestNewTagValueIncludesQualityAndLastSuccessfulRead(t *testing.T) {
	readAt := time.Unix(123, 456)
	message, err := newTagValue(1, 2, models.TagState{
		Value:              uint16(7),
		Quality:            models.TagQualityStale,
		LastSuccessfulRead: readAt,
	})
	if err != nil {
		t.Fatal(err)
	}
	if message.GetQuality() != realtimev1.TagQuality_TAG_QUALITY_STALE {
		t.Fatalf("quality = %v, want STALE", message.GetQuality())
	}
	if message.GetLastSuccessfulReadUnixNano() != readAt.UnixNano() {
		t.Fatalf("last successful read = %d, want %d", message.GetLastSuccessfulReadUnixNano(), readAt.UnixNano())
	}
}

func TestGetTagsSnapshotReturnsAllTagsAndGroupedEquipments(t *testing.T) {
	readAt := time.Unix(170, 99)
	server := newRealtimeTagServer(func() (map[uint]map[uint]models.TagState, error) {
		return map[uint]map[uint]models.TagState{
			2: {
				20: {Value: true, Quality: models.TagQualityGood, LastSuccessfulRead: readAt},
			},
			1: {
				10: {Value: int16(-7), Quality: models.TagQualityStale},
				11: {Value: float32(12.5), Quality: models.TagQualityBad, LastSuccessfulRead: readAt},
			},
		}, nil
	}, defaultSnapshotInterval)

	payload, err := server.GetTagsSnapshot(t.Context(), &emptypb.Empty{})
	if err != nil {
		t.Fatalf("GetTagsSnapshot() error = %v", err)
	}

	allTags := payload.Fields["all_tags"].GetListValue().GetValues()
	if len(allTags) != 3 {
		t.Fatalf("all_tags count = %d, want 3", len(allTags))
	}

	equipments := payload.Fields["equipments"].GetListValue().GetValues()
	if len(equipments) != 2 {
		t.Fatalf("equipments count = %d, want 2", len(equipments))
	}

	firstEquipment := equipments[0].GetStructValue().GetFields()
	if firstEquipment["equipment_id"].GetNumberValue() != 1 {
		t.Fatalf("first equipment id = %v, want 1", firstEquipment["equipment_id"].GetNumberValue())
	}
	if firstEquipment["tag_count"].GetNumberValue() != 2 {
		t.Fatalf("first equipment tag_count = %v, want 2", firstEquipment["tag_count"].GetNumberValue())
	}

	firstTag := firstEquipment["tags"].GetListValue().GetValues()[0].GetStructValue().GetFields()
	if firstTag["tag_id"].GetNumberValue() != 10 {
		t.Fatalf("first grouped tag id = %v, want 10", firstTag["tag_id"].GetNumberValue())
	}
	if firstTag["quality"].GetStringValue() != "STALE" {
		t.Fatalf("first grouped tag quality = %q, want STALE", firstTag["quality"].GetStringValue())
	}
}
