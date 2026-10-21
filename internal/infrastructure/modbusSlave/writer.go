package modbusSlave

import (
	"encoding/binary"
	"fmt"
	"middleware/internal/domain/models"
	"strconv"
)

func WriteTagSlave(client modbusSlaveClient, tag *models.Tag, value interface{}) error {
	if tag == nil {
		return fmt.Errorf("tag vazia")
	}

	switch operationIDSlave(tag) {
	case modbusOperationCoilStatus:
		boolValue, err := boolValueSlave(value)
		if err != nil {
			return fmt.Errorf("valor invalido para coil da tag %d: %w", tag.ID, err)
		}
		_, err = client.WriteSingleCoil(uint16(tag.Offset), coilWriteValueSlave(boolValue))
		return err
	case modbusOperationHoldingRegister:
		raw, quantity, err := rawRegistersValueByTypeSlave(value, tag)
		if err != nil {
			return err
		}
		if quantity == 1 {
			_, err = client.WriteSingleRegister(uint16(tag.Offset), binary.BigEndian.Uint16(raw[:2]))
			return err
		}
		_, err = client.WriteMultipleRegisters(uint16(tag.Offset), quantity, raw)
		return err
	case modbusOperationInputStatus:
		return fmt.Errorf("tag %d usa input status/discrete input, memoria somente leitura", tag.ID)
	case modbusOperationInputRegister:
		return fmt.Errorf("tag %d usa input register, memoria somente leitura", tag.ID)
	default:
		return fmt.Errorf("tag %d possui tipo de memoria modbus desconhecido: %d", tag.ID, tag.OperationID)
	}
}

func WriteTagsSlave(client modbusSlaveClient, tags []*models.Tag, value interface{}) error {
	validTags := make([]*models.Tag, 0, len(tags))
	for _, tag := range tags {
		if tag != nil {
			validTags = append(validTags, tag)
		}
	}
	if len(validTags) == 0 {
		return fmt.Errorf("nenhuma tag para escrita")
	}

	valuesByTag, err := writeValuesByTagSlave(validTags, value)
	if err != nil {
		return err
	}

	for operationID, operationTags := range groupTagsByOperationSlave(validTags) {
		switch operationID {
		case modbusOperationCoilStatus:
			if err := writeCoilTagsSlave(client, operationTags, valuesByTag); err != nil {
				return err
			}
		case modbusOperationHoldingRegister:
			if err := writeRegisterTagsSlave(client, operationTags, valuesByTag); err != nil {
				return err
			}
		case modbusOperationInputStatus:
			return fmt.Errorf("input status/discrete input e memoria somente leitura")
		case modbusOperationInputRegister:
			return fmt.Errorf("input register e memoria somente leitura")
		default:
			return fmt.Errorf("tipo de memoria modbus desconhecido: %d", operationID)
		}
	}

	return nil
}

func writeCoilTagsSlave(client modbusSlaveClient, tags []*models.Tag, valuesByTag map[uint]interface{}) error {
	blocks := createReadSlaveInBlocks(tags, maxCoilsPerRead, quantityBitsTagsSlave)
	for _, block := range blocks {
		quantityBits := block.endExclusive - block.start
		if quantityBits <= 0 {
			continue
		}

		boolValues := make([]bool, quantityBits)
		for _, tag := range block.tags {
			if tag == nil {
				continue
			}

			value, ok := valuesByTag[tag.ID]
			if !ok {
				return fmt.Errorf("valor nao informado para tag %d", tag.ID)
			}

			boolValue, err := boolValueSlave(value)
			if err != nil {
				return fmt.Errorf("valor invalido para coil da tag %d: %w", tag.ID, err)
			}

			bitOffset := int(tag.Offset) - block.start
			if bitOffset < 0 || bitOffset >= len(boolValues) {
				return fmt.Errorf("offset invalido para tag %d", tag.ID)
			}
			boolValues[bitOffset] = boolValue
		}

		if quantityBits == 1 {
			_, err := client.WriteSingleCoil(uint16(block.start), coilWriteValueSlave(boolValues[0]))
			if err != nil {
				return err
			}
			continue
		}

		_, err := client.WriteMultipleCoils(uint16(block.start), uint16(quantityBits), packCoilsSlave(boolValues))
		if err != nil {
			return err
		}
	}

	return nil
}

func writeRegisterTagsSlave(client modbusSlaveClient, tags []*models.Tag, valuesByTag map[uint]interface{}) error {
	blocks := createReadSlaveInBlocks(tags, maxRegistersPerRead, quantityRegistersTagsSlave)
	for _, block := range blocks {
		quantityRegisters := block.endExclusive - block.start
		if quantityRegisters <= 0 {
			continue
		}

		rawBlock := make([]byte, quantityRegisters*2)
		for _, tag := range block.tags {
			if tag == nil {
				continue
			}

			value, ok := valuesByTag[tag.ID]
			if !ok {
				return fmt.Errorf("valor nao informado para tag %d", tag.ID)
			}

			raw, _, err := rawRegistersValueByTypeSlave(value, tag)
			if err != nil {
				return err
			}

			offsetBytes := (int(tag.Offset) - block.start) * 2
			if offsetBytes < 0 || offsetBytes+len(raw) > len(rawBlock) {
				return fmt.Errorf("offset invalido para tag %d", tag.ID)
			}
			copy(rawBlock[offsetBytes:], raw)
		}

		if quantityRegisters == 1 {
			_, err := client.WriteSingleRegister(uint16(block.start), binary.BigEndian.Uint16(rawBlock[:2]))
			if err != nil {
				return err
			}
			continue
		}

		_, err := client.WriteMultipleRegisters(uint16(block.start), uint16(quantityRegisters), rawBlock)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeValuesByTagSlave(tags []*models.Tag, value interface{}) (map[uint]interface{}, error) {
	valuesByTag := make(map[uint]interface{}, len(tags))

	switch typedValue := value.(type) {
	case map[uint]interface{}:
		for _, tag := range tags {
			if tag == nil {
				continue
			}
			if value, ok := typedValue[tag.ID]; ok {
				valuesByTag[tag.ID] = value
			}
		}
	case map[string]interface{}:
		for key, value := range typedValue {
			tagID, err := strconv.ParseUint(key, 10, 64)
			if err != nil {
				continue
			}
			valuesByTag[uint(tagID)] = value
		}
	case []interface{}:
		if len(typedValue) != len(tags) {
			return nil, fmt.Errorf("quantidade de valores (%d) diferente da quantidade de tags (%d)", len(typedValue), len(tags))
		}
		for i, tag := range tags {
			if tag != nil {
				valuesByTag[tag.ID] = typedValue[i]
			}
		}
	default:
		if len(tags) != 1 {
			return nil, fmt.Errorf("escrita multipla precisa receber map[tagID]valor ou lista de valores")
		}
		valuesByTag[tags[0].ID] = value
	}

	for _, tag := range tags {
		if tag == nil {
			continue
		}
		if _, ok := valuesByTag[tag.ID]; !ok {
			return nil, fmt.Errorf("valor nao informado para tag %d", tag.ID)
		}
	}

	return valuesByTag, nil
}
