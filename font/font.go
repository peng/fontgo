package font

import (
	"encoding/binary"
	"os"
	"strconv"
	"time"
)

type OffsetTable struct {
	ScalerType    uint32
	NumTables     uint16
	SearchRange   uint16
	EntrySelector uint16
	RangeShift    uint16
}

type TagItem struct {
	CheckSum uint32
	Offset   uint32
	Length   uint32
}

type Tables struct {
	Head *Head
	Maxp *Maxp
	Loca []uint16
}

type Directory struct {
	OffsetTable  *OffsetTable
	TableContent map[string]*TagItem
	Tables       *Tables
	Glyphs       *Glyphs
}

func DataReader(filePath string) (directory *Directory, err error) {
	var fileByte []byte
	fileByte, err = os.ReadFile(filePath)

	if err != nil {
		return
	}

	// read offset table
	offsetTable := &OffsetTable{
		getUint32(fileByte[:4]),
		getUint16(fileByte[4:6]),
		getUint16(fileByte[6:8]),
		getUint16(fileByte[8:10]),
		getUint16(fileByte[10:12]),
	}

	// read table content
	tableContent := make(map[string]*TagItem)
	numTables := int(offsetTable.NumTables)
	pos := 12
	for i := 0; i < numTables; i++ {
		tagName := getString(fileByte[pos : pos+4])
		pos += 4
		tableContent[tagName] = &TagItem{
			getUint32(fileByte[pos : pos+4]),
			getUint32(fileByte[pos+4 : pos+8]),
			getUint32(fileByte[pos+8 : pos+12]),
		}
		pos += 12
	}

	headInfo := tableContent["head"]
	maxpInfo := tableContent["maxp"]
	locaInfo := tableContent["loca"]

	// tables content
	tables := new(Tables)
	tables.Head = GetHead(fileByte[headInfo.Offset : headInfo.Offset+headInfo.Length])
	tables.Maxp = GetMaxp(fileByte[maxpInfo.Offset : maxpInfo.Offset+maxpInfo.Length])
	tables.Loca = GetLoca(fileByte[locaInfo.Offset:locaInfo.Offset+locaInfo.Length], tables.Maxp.NumGlyphs, tables.Head.IndexToLocFormat)

	directory = new(Directory)

	directory.OffsetTable = offsetTable
	directory.TableContent = tableContent
	directory.Tables = tables
	glyfStart := tableContent["glyf"].Offset
	// glyfEnd := glyfStart + tableContent["glyf"].Length
	directory.Glyphs = GetGlyphs(fileByte[glyfStart:], tables.Loca)

	// directory = &Directory{offsetTable, tableContent, tables}
	return
}

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
