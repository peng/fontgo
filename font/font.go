package font

import (
	"os"
)


type TagItem struct {
	CheckSum uint32 `json:"checkSum"`
	Offset   uint32 `json:"offset"`
	Length   uint32 `json:"length"`
}

type Tables struct {
	Head *Head `json:"head"`
	Maxp *Maxp `json:"maxp`
	Loca []uint16 `json:"loca"`
}

type Directory struct {
	OffsetTable  *OffsetTable `json:"offsetTable"`
	TableContent map[string]*TagItem `json:"tableContent"`
	Tables       *Tables `json:"tables"`
	Glyphs       *Glyphs `json:"glyphs"`
}

func DataReader(filePath string) (directory *Directory, err error) {
	var fileByte []byte
	fileByte, err = os.ReadFile(filePath)

	if err != nil {
		return
	}

	// read offset table
	offsetTable := GetOffsetTable(fileByte)

	// read table content
	numTables := int(offsetTable.NumTables)
	tableContent := GetTableContent(numTables, fileByte)

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
	directory.Glyphs = GetGlyphs(fileByte[glyfStart:], tables.Loca)
	
	return
}
