package ncm

import (
	"encoding/binary"
	"errors"
	"io"
)

var ErrLength = errors.New("Length of Integer must be 1,2,4 or 8")

func readBytes(r io.Reader, length int) (interface{}, error) {
	switch length {
	case 1:
		var tmp uint8
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return 0, err
		}
		return tmp, nil
	case 2:
		var tmp uint16
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return 0, err
		}
		return tmp, nil
	case 4:
		var tmp uint32
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return 0, err
		}
		return tmp, nil
	case 8:
		var tmp uint64
		if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
			return 0, err
		}
		return tmp, nil
	default:
		return 0, ErrLength
	}
}

func readUint8(r io.Reader) (uint8, error) {
	tmp, err := readBytes(r, 1)
	res, _ := tmp.(uint8)
	return res, err
}

func readUint16(r io.Reader) (uint16, error) {
	tmp, err := readBytes(r, 2)
	res, _ := tmp.(uint16)
	return res, err
}

func readUint32(r io.Reader) (uint32, error) {
	tmp, err := readBytes(r, 4)
	res, _ := tmp.(uint32)
	return res, err
}

func readUint64(r io.Reader) (uint64, error) {
	tmp, err := readBytes(r, 8)
	res, _ := tmp.(uint64)
	return res, err
}
