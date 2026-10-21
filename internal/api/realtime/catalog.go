package grpcserver

import (
	"context"
	"fmt"
	"middleware/internal/domain/models"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

const RealtimeTagCatalogService_GetTagsSnapshot_FullMethodName = "/middleware.realtime.v1.RealtimeTagCatalogService/GetTagsSnapshot"

type RealtimeTagCatalogServiceServer interface {
	GetTagsSnapshot(context.Context, *emptypb.Empty) (*structpb.Struct, error)
}

func RegisterRealtimeTagCatalogServiceServer(s grpc.ServiceRegistrar, srv RealtimeTagCatalogServiceServer) {
	s.RegisterService(&RealtimeTagCatalogService_ServiceDesc, srv)
}

func _RealtimeTagCatalogService_GetTagsSnapshot_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RealtimeTagCatalogServiceServer).GetTagsSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RealtimeTagCatalogService_GetTagsSnapshot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RealtimeTagCatalogServiceServer).GetTagsSnapshot(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var RealtimeTagCatalogService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "middleware.realtime.v1.RealtimeTagCatalogService",
	HandlerType: (*RealtimeTagCatalogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTagsSnapshot",
			Handler:    _RealtimeTagCatalogService_GetTagsSnapshot_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/api/realtime/catalog.go",
}

func (s *RealtimeTagServer) GetTagsSnapshot(_ context.Context, _ *emptypb.Empty) (*structpb.Struct, error) {
	snapshot, err := s.readSnapshot()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "falha ao ler cache de tags: %v", err)
	}

	payload, err := snapshotToCatalog(snapshot)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "falha ao montar snapshot de tags: %v", err)
	}

	return structpb.NewStruct(payload)
}

func snapshotToCatalog(snapshot map[uint]map[uint]models.TagState) (map[string]any, error) {
	allTags := make([]any, 0)
	equipments := make([]any, 0, len(snapshot))

	for _, equipmentID := range sortedKeys(snapshot) {
		tagStates := snapshot[equipmentID]
		groupedTags := make([]any, 0, len(tagStates))

		for _, tagID := range sortedKeys(tagStates) {
			tagPayload, err := tagStateToCatalogEntry(equipmentID, tagID, tagStates[tagID])
			if err != nil {
				return nil, err
			}
			allTags = append(allTags, tagPayload)
			groupedTags = append(groupedTags, tagPayload)
		}

		equipments = append(equipments, map[string]any{
			"equipment_id": float64(equipmentID),
			"tag_count":    float64(len(groupedTags)),
			"tags":         groupedTags,
		})
	}

	return map[string]any{
		"all_tags":   allTags,
		"equipments": equipments,
	}, nil
}

func tagStateToCatalogEntry(equipmentID, tagID uint, state models.TagState) (map[string]any, error) {
	value, err := normalizeCatalogValue(equipmentID, tagID, state.Value)
	if err != nil {
		return nil, err
	}

	entry := map[string]any{
		"equipment_id": float64(equipmentID),
		"tag_id":       float64(tagID),
		"quality":      catalogQuality(state.Quality),
	}
	if value != nil {
		entry["value"] = value
	}
	if state.LastSuccessfulRead.IsZero() {
		entry["last_successful_read_unix_nano"] = "0"
		entry["last_successful_read"] = ""
		return entry, nil
	}

	entry["last_successful_read_unix_nano"] = fmt.Sprintf("%d", state.LastSuccessfulRead.UnixNano())
	entry["last_successful_read"] = state.LastSuccessfulRead.UTC().Format(time.RFC3339Nano)
	return entry, nil
}

func normalizeCatalogValue(equipmentID, tagID uint, value interface{}) (any, error) {
	if value == nil {
		return nil, nil
	}

	switch typed := value.(type) {
	case bool:
		return typed, nil
	case int:
		return float64(typed), nil
	case int8:
		return float64(typed), nil
	case int16:
		return float64(typed), nil
	case int32:
		return float64(typed), nil
	case int64:
		return float64(typed), nil
	case uint:
		return float64(typed), nil
	case uint8:
		return float64(typed), nil
	case uint16:
		return float64(typed), nil
	case uint32:
		return float64(typed), nil
	case uint64:
		return float64(typed), nil
	case float32:
		return float64(typed), nil
	case float64:
		return typed, nil
	case string:
		return typed, nil
	default:
		return nil, fmt.Errorf("tipo de valor nao suportado para equipamento %d, tag %d: %T", equipmentID, tagID, value)
	}
}

func catalogQuality(quality models.TagQuality) string {
	switch quality {
	case models.TagQualityGood:
		return "GOOD"
	case models.TagQualityStale:
		return "STALE"
	case models.TagQualityBad:
		return "BAD"
	default:
		return "UNSPECIFIED"
	}
}
