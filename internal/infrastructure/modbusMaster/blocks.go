package modbusmaster

import (
	"middleware/internal/domain/models"
	"sort"
	"strings"
)

func quantityRegistersTagsMaster(tag *models.Tag) uint16 {
	if tag == nil {
		return 1
	}
	switch tag.TypeID {
	case 1, 4:
		return 2
	default:
		return 1
	}
}

func quantityBitsTagsMaster(tag *models.Tag) uint16 {
	return 1
}

func createReadMasterInBlocks(tags []*models.Tag, maxQuantity int, quantityFunc func(*models.Tag) uint16) []readBlock {
	if maxQuantity <= 0 {
		maxQuantity = 1
	}

	ordenedTags := make([]*models.Tag, 0, len(tags))
	for _, tag := range tags {
		if tag != nil {
			ordenedTags = append(ordenedTags, tag)
		}
	}

	sort.Slice(ordenedTags, func(i, j int) bool {
		if ordenedTags[i].Offset == ordenedTags[j].Offset {
			return ordenedTags[i].ID < ordenedTags[j].ID
		}
		return ordenedTags[i].Offset < ordenedTags[j].Offset
	})

	blocks := make([]readBlock, 0, len(ordenedTags))
	var currentBlock readBlock

	for i, tag := range ordenedTags {
		startTag := int(tag.Offset)
		endTag := startTag + int(quantityFunc(tag))

		if i == 0 {
			currentBlock = readBlock{
				start:        startTag,
				endExclusive: endTag,
				tags:         []*models.Tag{tag},
			}
			continue
		}

		newEnd := endTag
		if currentBlock.endExclusive > newEnd {
			newEnd = currentBlock.endExclusive
		}

		isLimit := newEnd-currentBlock.start > maxQuantity
		isGap := startTag > currentBlock.endExclusive

		if isLimit || isGap {
			blocks = append(blocks, currentBlock)
			currentBlock = readBlock{
				start:        startTag,
				endExclusive: endTag,
				tags:         []*models.Tag{tag},
			}
			continue
		}

		currentBlock.tags = append(currentBlock.tags, tag)
		currentBlock.endExclusive = newEnd
	}

	if len(ordenedTags) > 0 {
		blocks = append(blocks, currentBlock)
	}

	return blocks
}

func groupTagsByOperationMaster(tags []*models.Tag) map[uint][]*models.Tag {
	tagsByOperation := make(map[uint][]*models.Tag)
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		operationID := operationIDMaster(tag)
		tagsByOperation[operationID] = append(tagsByOperation[operationID], tag)
	}
	return tagsByOperation
}

func operationIDMaster(tag *models.Tag) uint {
	if tag == nil {
		return modbusOperationHoldingRegister
	}
	if tag.OperationID != 0 {
		return tag.OperationID
	}

	description := strings.ToLower(strings.TrimSpace(tag.OperationType.Description))
	switch {
	case strings.Contains(description, "coil"):
		return modbusOperationCoilStatus
	case strings.Contains(description, "input status"), strings.Contains(description, "discrete"):
		return modbusOperationInputStatus
	case strings.Contains(description, "input register"):
		return modbusOperationInputRegister
	case strings.Contains(description, "holding"):
		return modbusOperationHoldingRegister
	default:
		return modbusOperationHoldingRegister
	}
}
