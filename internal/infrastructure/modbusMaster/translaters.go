package modbusmaster

import (
	"encoding/binary"
	"fmt"
	"math"
	"middleware/internal/domain/models"
	"strconv"
	"strings"
)

func applySwapMaster(raw []byte, swap string) []byte {
	if len(raw) < 4 {
		return raw
	}

	switch strings.ToUpper(strings.TrimSpace(swap)) {
	case "BADC":
		return []byte{raw[1], raw[0], raw[3], raw[2]}
	case "CDAB":
		return []byte{raw[2], raw[3], raw[0], raw[1]}
	case "DCBA":
		return []byte{raw[3], raw[2], raw[1], raw[0]}
	default:
		return []byte{raw[0], raw[1], raw[2], raw[3]}
	}
}

func parseValueByTypeMaster(raw []byte, tag *models.Tag) (interface{}, error) {
	if tag == nil {
		return nil, fmt.Errorf("tag vazia")
	}

	switch tag.TypeID {
	case 1:
		if len(raw) < 4 {
			return nil, fmt.Errorf("tag %d real precisa de 4 bytes", tag.ID)
		}
		value := applySwapMaster(raw[:4], tag.Swap.Description)
		return math.Float32frombits(binary.BigEndian.Uint32(value)), nil
	case 2:
		if len(raw) < 2 {
			return nil, fmt.Errorf("tag %d int precisa de 2 bytes", tag.ID)
		}
		return int16(binary.BigEndian.Uint16(raw[:2])), nil
	case 3:
		if len(raw) < 2 {
			return nil, fmt.Errorf("tag %d bool precisa de 2 bytes", tag.ID)
		}
		return binary.BigEndian.Uint16(raw[:2]) != 0, nil
	case 4:
		if len(raw) < 4 {
			return nil, fmt.Errorf("tag %d dword precisa de 4 bytes", tag.ID)
		}
		value := applySwapMaster(raw[:4], tag.Swap.Description)
		return binary.BigEndian.Uint32(value), nil
	default:
		if len(raw) < 2 {
			return nil, fmt.Errorf("tag %d precisa de 2 bytes", tag.ID)
		}
		return binary.BigEndian.Uint16(raw[:2]), nil
	}
}

func rawRegistersValueByTypeMaster(value interface{}, tag *models.Tag) ([]byte, uint16, error) {
	if tag == nil {
		return nil, 0, fmt.Errorf("tag vazia")
	}

	switch tag.TypeID {
	case 1:
		floatValue, err := float64ValueMaster(value)
		if err != nil {
			return nil, 0, fmt.Errorf("valor invalido para real da tag %d: %w", tag.ID, err)
		}
		raw := make([]byte, 4)
		binary.BigEndian.PutUint32(raw, math.Float32bits(float32(floatValue)))
		return applySwapMaster(raw, tag.Swap.Description), 2, nil
	case 2:
		intValue, err := int64ValueMaster(value)
		if err != nil {
			return nil, 0, fmt.Errorf("valor invalido para int da tag %d: %w", tag.ID, err)
		}
		if intValue < -32768 || intValue > 32767 {
			return nil, 0, fmt.Errorf("valor %d fora do limite int16 da tag %d", intValue, tag.ID)
		}
		raw := make([]byte, 2)
		binary.BigEndian.PutUint16(raw, uint16(int16(intValue)))
		return raw, 1, nil
	case 3:
		boolValue, err := boolValueMaster(value)
		if err != nil {
			return nil, 0, fmt.Errorf("valor invalido para bool da tag %d: %w", tag.ID, err)
		}
		raw := make([]byte, 2)
		if boolValue {
			binary.BigEndian.PutUint16(raw, 1)
		}
		return raw, 1, nil
	case 4:
		uintValue, err := uint64ValueMaster(value)
		if err != nil {
			return nil, 0, fmt.Errorf("valor invalido para dword da tag %d: %w", tag.ID, err)
		}
		if uintValue > 4294967295 {
			return nil, 0, fmt.Errorf("valor %d fora do limite uint32 da tag %d", uintValue, tag.ID)
		}
		raw := make([]byte, 4)
		binary.BigEndian.PutUint32(raw, uint32(uintValue))
		return applySwapMaster(raw, tag.Swap.Description), 2, nil
	default:
		uintValue, err := uint64ValueMaster(value)
		if err != nil {
			return nil, 0, fmt.Errorf("valor invalido para register da tag %d: %w", tag.ID, err)
		}
		if uintValue > 65535 {
			return nil, 0, fmt.Errorf("valor %d fora do limite uint16 da tag %d", uintValue, tag.ID)
		}
		raw := make([]byte, 2)
		binary.BigEndian.PutUint16(raw, uint16(uintValue))
		return raw, 1, nil
	}
}

func packCoilsMaster(values []bool) []byte {
	raw := make([]byte, (len(values)+7)/8)
	for i, value := range values {
		if value {
			raw[i/8] |= 1 << uint(i%8)
		}
	}
	return raw
}

func coilWriteValueMaster(value bool) uint16 {
	if value {
		return 0xFF00
	}
	return 0x0000
}

func boolValueMaster(value interface{}) (bool, error) {
	switch typedValue := value.(type) {
	case bool:
		return typedValue, nil
	case string:
		normalized := strings.ToLower(strings.TrimSpace(typedValue))
		switch normalized {
		case "true", "1", "on", "yes", "sim":
			return true, nil
		case "false", "0", "off", "no", "nao":
			return false, nil
		default:
			return false, fmt.Errorf("valor booleano invalido: %s", typedValue)
		}
	default:
		number, err := float64ValueMaster(value)
		if err != nil {
			return false, err
		}
		return number != 0, nil
	}
}

func int64ValueMaster(value interface{}) (int64, error) {
	number, err := float64ValueMaster(value)
	if err != nil {
		return 0, err
	}
	return int64(number), nil
}

func uint64ValueMaster(value interface{}) (uint64, error) {
	number, err := float64ValueMaster(value)
	if err != nil {
		return 0, err
	}
	if number < 0 {
		return 0, fmt.Errorf("valor negativo")
	}
	return uint64(number), nil
}

func float64ValueMaster(value interface{}) (float64, error) {
	switch typedValue := value.(type) {
	case int:
		return float64(typedValue), nil
	case int8:
		return float64(typedValue), nil
	case int16:
		return float64(typedValue), nil
	case int32:
		return float64(typedValue), nil
	case int64:
		return float64(typedValue), nil
	case uint:
		return float64(typedValue), nil
	case uint8:
		return float64(typedValue), nil
	case uint16:
		return float64(typedValue), nil
	case uint32:
		return float64(typedValue), nil
	case uint64:
		return float64(typedValue), nil
	case float32:
		return float64(typedValue), nil
	case float64:
		return typedValue, nil
	case string:
		return strconv.ParseFloat(strings.TrimSpace(typedValue), 64)
	default:
		return 0, fmt.Errorf("tipo %T nao suportado", value)
	}
}

func normalizeSlaveIDMaster(idPlc int) byte {
	if idPlc <= 0 {
		return 1
	}
	if idPlc > 255 {
		return 255
	}
	return byte(idPlc)
}
