package font

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

type TagItem struct {
	CheckSum uint32 `json:"checkSum"`
	Offset   uint32 `json:"offset"`
	Length   uint32 `json:"length"`
}

type Tables struct {
	Head *Head      `json:"head"`
	Maxp *Maxp      `json:"maxp"`
	Loca []int      `json:"loca"`
	Cmap *Cmap      `json:"cmap,omitempty"`
	Name *NameTable `json:"name,omitempty"`
	Hhea *Hhea      `json:"hhea,omitempty"`
	Hmtx *Hmtx      `json:"hmtx,omitempty"`
	Kern *Kern      `json:"kern,omitempty"`
	Os2  *OS2       `json:"os2"`
	Post *Post      `json:"post"`
	Fvar *Fvar      `json:"fvar"`
	Ltag *Ltag      `json:"ltag,omitempty"`
	Meta *Meta      `json:"meta"`
}

type FontInfo struct {
	OffsetTable  *OffsetTable        `json:"offsetTable"`
	TableContent map[string]*TagItem `json:"tableContent"`
	Tables       *Tables             `json:"tables"`
	Glyphs       *Glyphs             `json:"glyphs"`
}

type Font struct {
	fileByte []byte
	filePath string
	fontInfo *FontInfo
}

func ReadFontFile(filePath string) (f *Font, err error) {
	f = &Font{
		filePath: filePath,
	}
	f.fileByte, err = os.ReadFile(filePath)
	return
}

func (f *Font) GetFontInfo() (fontInfo *FontInfo, err error) {
	fileByte := f.fileByte
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
	itagInfo, existLtag := tableContent["Ltag"]
	metaInfo, existMeta := tableContent["meta"]
	// add test

	// tables content
	tables := new(Tables)
	if existHead {
		tables.Head = GetHead(fileByte, int(headInfo.Offset))
	}
	if existMaxp {
		tables.Maxp = GetMaxp(fileByte, int(maxpInfo.Offset))
	}
	if existLoca && tables.Maxp != nil && tables.Head != nil {
		tables.Loca = GetLoca(fileByte, int(locaInfo.Offset), tables.Maxp.NumGlyphs, tables.Head.IndexToLocFormat)
	}
	if existCmap && tables.Maxp != nil {
		tables.Cmap, err = GetCmap(fileByte, int(cmapInfo.Offset), int(tables.Maxp.NumGlyphs))
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

	if existLtag {
		ltag, ltagErr := GetLtag(fileByte, int(itagInfo.Offset))
		if ltagErr == nil {
			tables.Ltag = ltag
		}
	}

	if existMeta {
		meta, metaErr := GetMeta(fileByte, int(metaInfo.Offset))

		if metaErr == nil {
			tables.Meta = meta
		}
	}

	fontInfo = new(FontInfo)

	fontInfo.OffsetTable = offsetTable
	fontInfo.TableContent = tableContent
	fontInfo.Tables = tables
	glyfStart := tableContent["glyf"].Offset
	fontInfo.Glyphs = GetGlyphs(fileByte, int(glyfStart), tables.Loca, int(tables.Maxp.NumGlyphs))

	f.fontInfo = fontInfo

	return
}

func (f *Font) Write(filePath string) (err error) {
	ext := filepath.Ext(filePath)
	if ext != ".ttf" {
		err = errors.New("Not support format!")
		return
	}
	supportTable := []string{"cmap", "fvar", "glyf", "head", "hhea", "hmtx", "kern", "Ltag", "loca", "maxp", "meta", "name", "OS/2", "post"}

	fontInfo := f.fontInfo
	data := []byte{}
	pad4 := func(n int) int { return (n + 3) &^ 3 }
	computeCheckSum := func(b []byte) uint32 {
		padded := make([]byte, pad4(len(b)))
		copy(padded, b)
		var sum uint32
		for i := 0; i < len(padded); i += 4 {
			sum += getUint32(padded[i : i+4])
		}
		return sum
	}

	// rewrite offset table
	cpOffsetTable := &OffsetTable{}
	cpOffsetTable.ScalerType = fontInfo.OffsetTable.ScalerType
	cpOffsetTable.NumTables = uint16(len(supportTable))

	maxPowerOf2 := uint16(1)
	entrySelector := uint16(0)
	for maxPowerOf2*2 <= cpOffsetTable.NumTables {
		maxPowerOf2 *= 2
		entrySelector++
	}
	cpOffsetTable.SearchRange = maxPowerOf2 * 16
	cpOffsetTable.EntrySelector = entrySelector
	cpOffsetTable.RangeShift = (cpOffsetTable.NumTables*16 - cpOffsetTable.SearchRange)

	offsetTableData := WriteOffsetTable(cpOffsetTable)
	data = append(data, offsetTableData...)

	// rewrite tableContent table directory
	// write support tables first
	tablesData := map[string][]byte{}
	tablesData["cmap"], err = WriteCmap(fontInfo.Tables.Cmap)
	if err != nil {
		return
	}
	tablesData["fvar"] = WriteFvar(fontInfo.Tables.Fvar)
	tablesData["glyf"], err = WriteGlyphs(fontInfo.Glyphs)
	if err != nil {
		return
	}
	head := *fontInfo.Tables.Head
	head.CheckSumAdjustment = 0
	tablesData["head"] = WriteHead(&head)
	tablesData["hhea"] = WriteHhea(fontInfo.Tables.Hhea)
	tablesData["hmtx"] = WriteHmtx(fontInfo.Tables.Hmtx)
	tablesData["kern"] = WriteKern(fontInfo.Tables.Kern)
	tablesData["Ltag"] = WriteLtag(fontInfo.Tables.Ltag)
	tablesData["loca"] = WriteLoca(fontInfo.Tables.Loca, fontInfo.Tables.Head.IndexToLocFormat)
	tablesData["maxp"] = WriteMaxp(fontInfo.Tables.Maxp)
	tablesData["meta"] = WriteMeta(fontInfo.Tables.Meta)
	tablesData["name"] = WriteName(fontInfo.Tables.Name)
	tablesData["OS/2"] = WriteOS2(fontInfo.Tables.Os2)
	tablesData["post"] = WritePost(fontInfo.Tables.Post)

	// table directory size
	tableDirSize := len(supportTable) * 16
	nextOffset := uint32(len(offsetTableData) + tableDirSize)

	// Build the directory entries and output using WriteTableContent
	cpTableContent := make(map[string]*TagItem, len(supportTable))
	for _, tag := range supportTable {
		td := tablesData[tag]
		length := uint32(len(td))
		checkSum := computeCheckSum(td)
		cpTableContent[tag] = &TagItem{
			CheckSum: checkSum,
			Offset:   nextOffset,
			Length:   length,
		}
		nextOffset += uint32(pad4(int(length)))
	}

	data = append(data, WriteTableContent(cpTableContent)...)

	// append table payloads with padding, regenerate and write data according to cpTableContent
	for _, tag := range supportTable {
		var td []byte
		switch tag {
		case "cmap":
			td, _ = WriteCmap(fontInfo.Tables.Cmap)
		case "fvar":
			td = WriteFvar(fontInfo.Tables.Fvar)
		case "glyf":
			td, _ = WriteGlyphs(fontInfo.Glyphs)
		case "head":
			head := *fontInfo.Tables.Head
			head.CheckSumAdjustment = 0
			td = WriteHead(&head)
		case "hhea":
			td = WriteHhea(fontInfo.Tables.Hhea)
		case "hmtx":
			td = WriteHmtx(fontInfo.Tables.Hmtx)
		case "kern":
			td = WriteKern(fontInfo.Tables.Kern)
		case "Ltag":
			td = WriteLtag(fontInfo.Tables.Ltag)
		case "loca":
			td = WriteLoca(fontInfo.Tables.Loca, fontInfo.Tables.Head.IndexToLocFormat)
		case "maxp":
			td = WriteMaxp(fontInfo.Tables.Maxp)
		case "meta":
			td = WriteMeta(fontInfo.Tables.Meta)
		case "name":
			td = WriteName(fontInfo.Tables.Name)
		case "OS/2":
			td = WriteOS2(fontInfo.Tables.Os2)
		case "post":
			td = WritePost(fontInfo.Tables.Post)
		default:
			log.Printf("[WARN] table %s data missing during payload write, skip", tag)
			td = []byte{}
		}
		data = append(data, td...)
		padLen := pad4(len(td)) - len(td)
		if padLen > 0 {
			data = append(data, make([]byte, padLen)...)
		}
	}

	if err = os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return
	}
	return os.WriteFile(filePath, data, 0o755)
}
