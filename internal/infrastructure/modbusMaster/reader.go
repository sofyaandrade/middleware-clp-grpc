package modbusmaster

import (
	"fmt"
	"middleware/internal/domain/models"

	"github.com/goburrow/modbus"
)

func ReadTagsMaster(client modbus.Client, tags []*models.Tag, previousValue map[uint]interface{}) (map[uint]interface{}, error) {
	tagsMap := make(map[uint]interface{})
	tagsByOperation := groupTagsByOperationMaster(tags)
	if len(tagsByOperation) == 0 {
		return tagsMap, nil
	}

	var firstError error

	for operationID, operationTags := range tagsByOperation {
		var values map[uint]interface{}
		var err error

		switch operationID {
		case modbusOperationCoilStatus:
			values, err = readBitTagsMaster(client.ReadCoils, operationTags, previousValue, "coils")
		case modbusOperationInputStatus:
			values, err = readBitTagsMaster(client.ReadDiscreteInputs, operationTags, previousValue, "discrete inputs")
		case modbusOperationInputRegister:
			values, err = readRegisterTagsMaster(client.ReadInputRegisters, operationTags, previousValue, "input registers")
		default:
			values, err = readRegisterTagsMaster(client.ReadHoldingRegisters, operationTags, previousValue, "holding registers")
		}

		for tagID, value := range values {
			tagsMap[tagID] = value
		}
		if err != nil && firstError == nil {
			firstError = err
		}
	}

	if firstError != nil {
		return tagsMap, firstError
	}
	return tagsMap, nil
}

func readRegisterTagsMaster(read readRegisterFunc, tags []*models.Tag, previousValue map[uint]interface{}, memoryName string) (map[uint]interface{}, error) {
	tagsMap := make(map[uint]interface{})
	blocks := createReadMasterInBlocks(tags, maxRegistersPerRead, quantityRegistersTagsMaster)
	if len(blocks) == 0 {
		return tagsMap, nil
	}

	var firstError error

	for _, block := range blocks {
		quantityRegisters := block.endExclusive - block.start
		if quantityRegisters <= 0 {
			continue
		}

		results, err := read(uint16(block.start), uint16(quantityRegisters))
		if err != nil {
			if firstError == nil {
				firstError = fmt.Errorf("falha ao ler %s start=%d quantity=%d: %w", memoryName, block.start, quantityRegisters, err)
			}
			for _, tag := range block.tags {
				if tag == nil {
					continue
				}
				if value, ok := previousValue[tag.ID]; ok {
					tagsMap[tag.ID] = value
				}
			}
			continue
		}

		for _, tag := range block.tags {
			if tag == nil {
				continue
			}

			offsetBytes := (int(tag.Offset) - block.start) * 2
			sizeBytes := int(quantityRegistersTagsMaster(tag) * 2)
			if offsetBytes < 0 || offsetBytes+sizeBytes > len(results) {
				if value, ok := previousValue[tag.ID]; ok {
					tagsMap[tag.ID] = value
				}
				continue
			}

			raw := results[offsetBytes : offsetBytes+sizeBytes]
			value, err := parseValueByTypeMaster(raw, tag)
			if err != nil {
				if firstError == nil {
					firstError = err
				}
				if previous, ok := previousValue[tag.ID]; ok {
					tagsMap[tag.ID] = previous
				}
				continue
			}
			tagsMap[tag.ID] = value
		}
	}

	if firstError != nil {
		return tagsMap, firstError
	}
	return tagsMap, nil
}

func readBitTagsMaster(read readBitsFunc, tags []*models.Tag, previousValue map[uint]interface{}, memoryName string) (map[uint]interface{}, error) {
	tagsMap := make(map[uint]interface{})
	blocks := createReadMasterInBlocks(tags, maxCoilsPerRead, quantityBitsTagsMaster)
	if len(blocks) == 0 {
		return tagsMap, nil
	}

	var firstError error

	for _, block := range blocks {
		quantityBits := block.endExclusive - block.start
		if quantityBits <= 0 {
			continue
		}

		results, err := read(uint16(block.start), uint16(quantityBits))
		if err != nil {
			if firstError == nil {
				firstError = fmt.Errorf("falha ao ler %s start=%d quantity=%d: %w", memoryName, block.start, quantityBits, err)
			}
			for _, tag := range block.tags {
				if tag == nil {
					continue
				}
				if value, ok := previousValue[tag.ID]; ok {
					tagsMap[tag.ID] = value
				}
			}
			continue
		}

		for _, tag := range block.tags {
			if tag == nil {
				continue
			}

			bitOffset := int(tag.Offset) - block.start
			value, ok := readBitFromBytesMaster(results, bitOffset)
			if !ok {
				if previous, exists := previousValue[tag.ID]; exists {
					tagsMap[tag.ID] = previous
				}
				continue
			}

			tagsMap[tag.ID] = value
		}
	}

	if firstError != nil {
		return tagsMap, firstError
	}
	return tagsMap, nil
}

func readBitFromBytesMaster(raw []byte, bitOffset int) (bool, bool) {
	if bitOffset < 0 {
		return false, false
	}

	byteOffset := bitOffset / 8
	if byteOffset >= len(raw) {
		return false, false
	}

	bit := uint(bitOffset % 8)
	return raw[byteOffset]&(1<<bit) != 0, true
}
