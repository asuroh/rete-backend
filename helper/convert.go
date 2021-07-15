package helper

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

type UUID struct {
	MSB uint64
	LSB uint64
}

func Uint64ToBytes(n uint64) []byte {
	bytes := make([]byte, 8)
	bytes[0] = byte((n >> 56) & 0xFF)
	bytes[1] = byte((n >> 48) & 0xFF)
	bytes[2] = byte((n >> 40) & 0xFF)
	bytes[3] = byte((n >> 32) & 0xFF)
	bytes[4] = byte((n >> 24) & 0xFF)
	bytes[5] = byte((n >> 16) & 0xFF)
	bytes[6] = byte((n >> 8) & 0xFF)
	bytes[7] = byte(n & 0xFF)

	return bytes
}

func (id UUID) String() string {
	msb := Uint64ToBytes(id.MSB)
	lsb := Uint64ToBytes(id.LSB)
	uid, err := uuid.FromBytes(append(msb, lsb...))
	if err != nil {
		return ""
	}

	return uid.String()
}

// SplitDate ...
func SplitDate(data string) string {
	dataArr := strings.Split(data, "-")
	if len(dataArr) != 3 {
		return ""
	}

	return dataArr[0] + dataArr[1] + dataArr[2]
}
