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
	Head *Head      `json:"head"`
	Maxp *Maxp      `json:"maxp`
	Loca []uint16   `json:"loca"`
	Cmap *Cmap      `json:"cmap,omitempty"`
	Name *NameTable `json:"name,omitempty"`
	Hhea *Hhea      `json:"hhea,omitempty"`
	Hmtx *Hmtx      `json:"hmtx,omitempty"`
	Kern *Kern      `json:"kern,omitempty"`
	Os2  *OS2       `json:"os2"`
	Post *Post      `json:"post"`
	Fvar *Fvar      `json:"fvar"`
}

type Directory struct {
	OffsetTable  *OffsetTable        `json:"offsetTable"`
	TableContent map[string]*TagItem `json:"tableContent"`
	Tables       *Tables             `json:"tables"`
	Glyphs       *Glyphs             `json:"glyphs"`
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

	headInfo, existHead := tableContent["head"]
	maxpInfo, existMaxp := tableContent["maxp"]
	locaInfo, existLoca := tableContent["loca"]
	cmapInfo, existCmap := tableContent["cmap"]
	nameInfo, existName := tableContent["name"]
	hheaInfo, existHhea := tableContent["hhea"]
	hmtxInfo, existHmtx := tableContent["hmtx"]
	kernInfo, existKern := tableContent["kern"]
	os2Info, existOs2 := tableContent["OS/2"]
	postInfo, existPost := tableContent["post"]
	fvarInfo, existFvar := tableContent["fvar"]
	// add test

	// tables content
	tables := new(Tables)
	if existHead {
		tables.Head = GetHead(fileByte[headInfo.Offset : headInfo.Offset+headInfo.Length])
	}
	if existMaxp {
		tables.Maxp = GetMaxp(fileByte[maxpInfo.Offset : maxpInfo.Offset+maxpInfo.Length])
	}
	if existLoca && tables.Maxp != nil && tables.Head != nil {
		tables.Loca = GetLoca(fileByte[locaInfo.Offset:locaInfo.Offset+locaInfo.Length], tables.Maxp.NumGlyphs, tables.Head.IndexToLocFormat)
	}
	if existCmap && tables.Maxp != nil {
		tables.Cmap, err = GetCmap(fileByte[cmapInfo.Offset:cmapInfo.Offset+cmapInfo.Length], int(cmapInfo.Offset), int(tables.Maxp.NumGlyphs))
	}

	if existName {
		tables.Name = GetName(fileByte, int(nameInfo.Offset))
	}

	if existHhea {
		tables.Hhea = GetHhea(fileByte, int(hheaInfo.Offset))
	}

	if existHmtx && existHhea && existMaxp {
		tables.Hmtx = GetHmtx(fileByte, int(hmtxInfo.Offset), int(tables.Hhea.NumOfLongHorMetrics), int(tables.Maxp.NumGlyphs))
	}

	if existKern {
		tables.Kern, err = GetKern(fileByte, int(kernInfo.Offset))
	}

	if existOs2 {
		tables.Os2 = GetOS2(fileByte, int(os2Info.Offset))
	}

	if existPost {
		tables.Post = GetPost(fileByte, int(postInfo.Offset))
	}

	if existFvar {
		tables.Fvar = GetFvar(fileByte, int(fvarInfo.Offset))
	}

	directory = new(Directory)

	directory.OffsetTable = offsetTable
	directory.TableContent = tableContent
	directory.Tables = tables
	glyfStart := tableContent["glyf"].Offset
	directory.Glyphs = GetGlyphs(fileByte[glyfStart:], tables.Loca)

	return
}
