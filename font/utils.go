package font

import (
	"encoding/binary"
	"strconv"
	"time"
)

func getUint8(data []byte) uint8 {
	return data[0]
}

func getUint16(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)
}

func getUint32(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

func getUint64(data []byte) uint64 {
	return binary.BigEndian.Uint64(data)
}

func getInt8(data []byte) int8 {
	return int8(data[0])
}

func getInt16(data []byte) int16 {
	return int16(getUint16((data)))
}

func getInt32(data []byte) int32 {
	return int32(getUint32(data))
}

func getInt64(data []byte) int64 {
	return int64(getUint64(data))
}

func getString(data []byte) string {
	return string(data)
}

func getFixed(data []byte) float64 {
	return float64(getInt32(data) / 65535)
}

func getFword(data []byte) int16 {
	return getInt16(data)
}

func get2Dot14(data []byte) float32 {
	return float32(getInt16(data) / 16384)
}

func getLongDateTime(data []byte) string {
	longDateTime := getInt64(data)
	starTime := time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()

	unixTime := longDateTime - starTime

	return time.Unix(unixTime, 0).Local().Format(time.UnixDate)
}

func getVersion(data []byte) string {
	// 32 bytes
	major := strconv.Itoa(int(getUint16(data[:2])))
	minor := strconv.Itoa(int(getUint16(data[2:4])))
	return major + "." + minor
}