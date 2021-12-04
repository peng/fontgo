package font

import (
	"bytes"
	"encoding/binary"
	"os"
)

type offsetTableStruct struct {
	scalerType uint32
	numTables uint16
	searchRange uint16
	entrySelector uint16
	rangeShift uint16
}

type tagItemStruct struct {
	checkSum uint32
	offset uint32
	length uint32
}

type DirectoryStruct struct {
	offsetTable offsetTableStruct
	tableContent map[string] tagItemStruct
}

func DataReader (filePath string) (directory DirectoryStruct, readErr error) {
	fileByte, err := os.ReadFile(filePath)

	if err != nil {
		readErr = err
		return
	}

	buf := bytes.NewBuffer(fileByte)

	// read offset table
	var offsetTable offsetTableStruct
	// offsetTable := new(offsetTableStruct)
	offsetTable.scalerType = getUint32(buf.Next(4))
	offsetTable.numTables = getUint16(buf.Next(2))
	offsetTable.searchRange = getUint16(buf.Next(2))
	offsetTable.entrySelector = getUint16(buf.Next(2))
	offsetTable.rangeShift = getUint16(buf.Next(2))

	// read table content

	tableContent := make(map[string]tagItemStruct)
	numTables := int(offsetTable.numTables)
	for i := 0; i < numTables; i++ {
		tagName := getTag(buf.Next(4))
		tableContent[tagName] = tagItemStruct{
			getUint32(buf.Next(4)),
			getUint32(buf.Next(4)),
			getUint32(buf.Next(4)),
		}
	}

	directory = DirectoryStruct{offsetTable, tableContent}
	// fmt.Printf("%v \n", directory)
	return
}

func getUint32(data []byte) uint32{
	return binary.BigEndian.Uint32(data)
}

func getUint16(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)	
}

func getTag(data []byte) string {
	return string(data)
}