package font

import (
	"encoding/binary"
	"os"
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

type Directory struct {
	OffsetTable  *OffsetTable
	TableContent map[string]*TagItem
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
		tagName := getTag(fileByte[pos : pos+4])
		pos += 4
		tableContent[tagName] = &TagItem{
			getUint32(fileByte[pos : pos+4]),
			getUint32(fileByte[pos+4 : pos+8]),
			getUint32(fileByte[pos+8 : pos+12]),
		}
		pos += 12
	}

	directory = &Directory{offsetTable, tableContent}
	return
}

func getUint32(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

func getUint16(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)
}

func getTag(data []byte) string {
	return string(data)
}
