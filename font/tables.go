package font

import "errors"

type OffsetTable struct {
	ScalerType    uint32 `json:"scalerType"`
	NumTables     uint16 `json:"numTables"`
	SearchRange   uint16 `json:"searchRange"`
	EntrySelector uint16 `json:"entrySelector"`
	RangeShift    uint16 `json:"rangeShift"`
}

func GetOffsetTable(data []byte) *OffsetTable {
	return &OffsetTable{
		getUint32(data[:4]),
		getUint16(data[4:6]),
		getUint16(data[6:8]),
		getUint16(data[8:10]),
		getUint16(data[10:12]),
	}
}

func GetTableContent(numTables int, date []byte) map[string]*TagItem {
	tableContent := make(map[string]*TagItem)
	pos := 12
	for i := 0; i < numTables; i++ {
		tagName := getString(date[pos : pos+4])
		pos += 4
		tableContent[tagName] = &TagItem{
			getUint32(date[pos : pos+4]),
			getUint32(date[pos+4 : pos+8]),
			getUint32(date[pos+8 : pos+12]),
		}
		pos += 12
	}
	return tableContent
}

type Head struct {
	Version            float64 `json:"version"`
	FontRevision       float64 `json:"fontRevision"`
	CheckSumAdjustment uint32  `json:"checkSumAdjustment"`
	MagicNumber        uint32  `json:"magicNumber"`
	Flags              uint16  `json:"flags"`
	UnitsPerEm         uint16  `json:"unitsPerEm"`
	Created            int64   `json:"created"`
	Modified           int64   `json:"modified"`
	XMin               int16   `json:"xMin"`
	YMin               int16   `json:"yMin"`
	XMax               int16   `json:"xMax"`
	YMax               int16   `json:"yMax"`
	MacStyle           uint16  `json:"macStyle"`
	LowestRecPPEM      uint16  `json:"lowestRecPpem"`
	FontDirectionHint  int16   `json:"fontDirectionHint"`
	IndexToLocFormat   int16   `json:"indexToLocFormat"`
	GlyphDataFormat    int16   `json:"glyphDataFormat"`
}

func GetHead(data []byte) *Head {
	return &Head{
		getFixed(data[:4]),
		getFixed(data[4:8]),
		getUint32(data[8:12]),
		getUint32(data[12:16]),
		getUint16(data[16:18]),
		getUint16(data[18:20]),
		getLongDateTime(data[20:28]),
		getLongDateTime(data[28:36]),
		getFword(data[36:38]),
		getFword(data[38:40]),
		getFword(data[40:42]),
		getFword(data[42:44]),
		getUint16(data[44:46]),
		getUint16(data[46:48]),
		getInt16(data[48:50]),
		getInt16(data[50:52]),
		getInt16(data[52:54]),
	}
}

type Flag struct {
	OnCurve      bool `json:"onCurve"`
	XShortVector bool `json:"xShortVector`
	YShortVector bool `json:"yShortVector"`
	XSame        bool `json:"xSame"`
	YSame        bool `json:"ySave"`
}
type Point struct {
	X    int   `json:"x"`
	Y    int   `json:"y"`
	Flag *Flag `json:"flag"`
}

type GlyphCommon struct {
	NumberOfContours int16  `json:"numberOfContours"`
	XMin             int16  `json:"xMin"`
	YMin             int16  `json:"yMin"`
	XMax             int16  `json:"xMax"`
	YMax             int16  `json:"yMax"`
	Type             string `json:"type"`
}

type GlyphSimple struct {
	GlyphCommon
	EndPtsOfContours  []uint16 `json:"endPtsOfContours"`
	InstructionLength uint16   `json:"instructionLength"`
	Instructions      []uint8  `json:"instructions"`
	Points            []*Point `json:"points"`
}

type Component struct {
	Flags      uint16  `json:"flags"`
	GlyphIndex uint16  `json:"glyphIndex"`
	Argument1  int     `json:"argument1"`
	Argument2  int     `json:"argument2"`
	Unsign     bool    `json:"unsign"`
	Scale      float32 `json:"scale"`
	Xscale     float32 `json:"xscale"`
	Yscale     float32 `json:"yscale"`
	Scale01    float32 `json:"scale01"`
	Scale10    float32 `json:"scale10"`
}

type GlyphCompound struct {
	GlyphCommon
	Component         []Component
	instructionLength int
	instructions      []uint8
}

type Glyphs struct {
	Simples   []GlyphSimple
	Compounds []GlyphCompound
}

const GLYPH_TYPE_SIMPLE, GLYPH_TYPE_COMPOUND = "simple", "compound"

func GetGlyphSimple(data []byte) (simple *GlyphSimple, pos int) {
	simple = new(GlyphSimple)
	simple.Type = GLYPH_TYPE_SIMPLE
	simple.NumberOfContours = getInt16(data[0:2])
	simple.XMin = getFword(data[2:4])
	simple.YMin = getFword(data[4:6])
	simple.XMax = getFword(data[6:8])
	simple.YMax = getFword(data[8:10])

	pos = 10
	// get endPtsOfContours
	for i := 0; i < int(simple.NumberOfContours); i++ {
		simple.EndPtsOfContours = append(simple.EndPtsOfContours, getUint16(data[pos:pos+2]))
		pos += 2
	}

	// get instructionLength
	simple.InstructionLength = getUint16(data[pos : pos+2])
	for i := 0; i < int(simple.InstructionLength); i++ {
		simple.Instructions = append(simple.Instructions, getUint8(data[pos:pos+1]))
		// test pos++
		pos++
	}

	// get points num
	pointsNum := int(simple.EndPtsOfContours[0])

	for _, num := range simple.EndPtsOfContours {
		contoursNum := int(num)
		if contoursNum > pointsNum {
			pointsNum = contoursNum
		}
	}

	var flags []uint8

	// get flags
	for i := 0; i < pointsNum; i++ {
		f := getUint8(data[pos : pos+1])
		flags = append(flags, f)
		pos++

		if f&0x08 == 1 {
			repeatNum := int(getUint8(data[pos : pos+1]))
			pos++

			for j := 0; j < repeatNum; j++ {
				flags = append(flags, getUint8(data[pos:pos+1]))
				pos++
			}
			i += repeatNum
		}
	}

	// should check number of flags same with points

	// get x points
	for i := 0; i < pointsNum; i++ {

		flagBit := flags[i]

		flag := &Flag{
			flagBit&0x01 == 1,
			flagBit&0x02 == 1,
			flagBit&0x04 == 1,
			flagBit&0x10 == 1,
			flagBit&0x20 == 1,
		}

		var point Point
		point.Flag = flag
		if flag.XShortVector {
			point.X = int(getUint8(data[pos : pos+1]))
			pos++
		} else {
			point.X = int(getUint16(data[pos : pos+2]))
			pos += 2
		}

		simple.Points = append(simple.Points, &point)
	}

	// get y points
	for i := 0; i < pointsNum; i++ {
		var y int
		point := simple.Points[i]
		if point.Flag.YShortVector {
			y = int(getUint8(data[pos : pos+1]))
			pos++
		} else {
			y = int(getUint16(data[pos : pos+2]))
			pos += 2
		}
		point.Y = y
	}

	return
}

func GetGlyphCompound(data []byte) (compound *GlyphCompound, pos int) {
	compound = new(GlyphCompound)
	compound.Type = GLYPH_TYPE_COMPOUND
	const (
		ARG_1_AND_2_ARE_WORDS    = 0x0001
		ARGS_ARE_XY_VALUES       = 0x0002
		ROUND_XY_TO_GRID         = 0x0004
		WE_HAVE_A_SCALE          = 0x0008
		MORE_COMPONENTS          = 0x0020
		WE_HAVE_AN_X_AND_Y_SCALE = 0x0040
		WE_HAVE_A_TWO_BY_TWO     = 0x0080
		WE_HAVE_INSTRUCTIONS     = 0x0100
		USE_MY_METRICS           = 0x0200
		OVERLAP_COMPOUND         = 0x0400
	)

	compound.Type = GLYPH_TYPE_COMPOUND
	compound.NumberOfContours = getInt16(data[0:2])
	compound.XMin = getFword(data[2:4])
	compound.YMin = getFword(data[4:6])
	compound.XMax = getFword(data[6:8])
	compound.YMax = getFword(data[8:10])

	var flags uint16
	pos = 10

	moreComponent := true

	for moreComponent {
		component := new(Component)

		component.Flags = getUint16(data[pos : pos+2])
		pos += 2

		flags = component.Flags

		if flags&ARG_1_AND_2_ARE_WORDS == 1 {
			if flags&ARGS_ARE_XY_VALUES == 1 {
				component.Argument1 = int(getInt16(data[pos : pos+2]))
				pos += 2
				component.Argument2 = int(getInt16(data[pos : pos+2]))
			} else {
				component.Unsign = true
				component.Argument1 = int(getUint16(data[pos : pos+2]))
				pos += 2
				component.Argument2 = int(getUint16(data[pos : pos+2]))
			}
			pos += 2
		} else {
			if flags&ARGS_ARE_XY_VALUES == 1 {
				component.Argument1 = int(getInt8(data[pos : pos+1]))
				pos++
				component.Argument2 = int(getInt8(data[pos : pos+1]))
			} else {
				component.Unsign = true
				component.Argument1 = int(getUint8(data[pos : pos+1]))
				pos++
				component.Argument2 = int(getUint8(data[pos : pos+1]))
			}
			pos++
		}

		if flags&WE_HAVE_A_SCALE == 1 {
			component.Scale = get2Dot14(data[pos : pos+2])
		} else if flags&WE_HAVE_AN_X_AND_Y_SCALE == 1 {
			component.Xscale = get2Dot14((data[pos : pos+2]))
			pos += 2
			component.Yscale = get2Dot14(data[pos : pos+2])
		} else if flags&WE_HAVE_A_TWO_BY_TWO == 1 {
			component.Xscale = get2Dot14(data[pos : pos+2])
			pos += 2
			component.Scale01 = get2Dot14(data[pos : pos+2])
			pos += 2
			component.Scale10 = get2Dot14(data[pos : pos+2])
			pos += 2
			component.Xscale = get2Dot14(data[pos : pos+2])
		}
		pos += 2

		moreComponent = flags&MORE_COMPONENTS == 1
	}

	// can't understand
	if flags&WE_HAVE_INSTRUCTIONS == 1 {
		compound.instructionLength = int(getUint16(data[pos : pos+2]))
		pos += 2

		for i := 0; i < compound.instructionLength; i++ {
			compound.instructions = append(compound.instructions, getUint8(data[pos:pos+1]))
			pos++
		}
	}

	return
}

func GetGlyphs(data []byte, loca []uint16) (glyphs *Glyphs) {
	glyphs = new(Glyphs)

	for i := 0; i < 30; i++ {

		offset := int(loca[i])
		nextOffset := int(loca[i+1])
		// fmt.Printf("innoffset %v", offset)
		// fmt.Printf("innnextoffset %v", nextOffset)
		numberOfContours := getInt16(data[offset : offset+2])

		if offset != nextOffset {
			if numberOfContours >= 0 {
				// simple
				simp, _ := GetGlyphSimple(data[offset:])
				glyphs.Simples = append(glyphs.Simples, *simp)
			} else {
				// compound
				compound, _ := GetGlyphCompound(data[offset:])
				glyphs.Compounds = append(glyphs.Compounds, *compound)
			}
		}
	}
	return
}

type Maxp struct {
	Version               string
	NumGlyphs             uint16
	MaxPoints             uint16
	MaxContours           uint16
	MaxComponentPoints    uint16
	MaxComponentContours  uint16
	MaxZones              uint16
	MaxTwilightPoints     uint16
	MaxStorage            uint16
	MaxFunctionDefs       uint16
	MaxInstructionDefs    uint16
	MaxStackElements      uint16
	MaxSizeOfInstructions uint16
	MaxComponentElements  uint16
	MaxComponentDepth     uint16
}

func GetMaxp(data []byte) *Maxp {
	maxp := new(Maxp)
	maxp.Version = getVersion(data[0:4])
	maxp.NumGlyphs = getUint16(data[4:6])
	if maxp.Version == "1.0" {
		maxp.MaxPoints = getUint16(data[6:8])
		maxp.MaxContours = getUint16(data[8:10])
		maxp.MaxComponentPoints = getUint16(data[10:12])
		maxp.MaxComponentContours = getUint16(data[12:14])
		maxp.MaxZones = getUint16(data[14:16])
		maxp.MaxTwilightPoints = getUint16(data[16:18])
		maxp.MaxStorage = getUint16(data[18:20])
		maxp.MaxFunctionDefs = getUint16(data[20:22])
		maxp.MaxInstructionDefs = getUint16(data[22:24])
		maxp.MaxStackElements = getUint16(data[24:26])
		maxp.MaxSizeOfInstructions = getUint16(data[26:28])
		maxp.MaxComponentElements = getUint16(data[28:30])
		maxp.MaxComponentDepth = getUint16(data[30:32])
	}

	return maxp
}

func GetLoca(data []byte, numGlyphs uint16, indexToLocFormat int16) []uint16 {
	// long version:  otf, ttf is different
	var locations []uint16
	pos := 0
	for i := 0; i < int(numGlyphs)+1; i++ {
		offset := getUint16(data[pos : pos+2])
		if indexToLocFormat == 0 {
			// 0 is short, 1 is long
			offset *= 2
		}
		locations = append(locations, offset)
		pos += 2
	}

	return locations
}

type Cmap struct {
	Version         uint16 `json:"version"`
	NumberSubtables uint16 `json:"numberSubtables"`
	SubTables []map[string]interface{} `json:"subTables"`
}

type CmapChild struct {
	PlatformID         uint16 `json:"platformId"`
	PlatformSpecificID uint16 `json:"platformSpecificId"`
	Offset             uint16 `json:"offset"`
}

type CmapFormat0 struct {
	Format          uint16     `json:"format"`
	Length          uint16     `json:"length"`
	Language        uint16     `json:"language"`
	GlyphIndexArray [256]uint8 `json:"glyphIndexArray"`
}

type CmapFormat2 struct {
	Format          uint16      `json:"format"`
	Length          uint16      `json:"length"`
	Language        uint16      `json:"language"`
	SubHeaderKeys   [256]uint16 `json:"subHeaderKeys"`
	SubHeaders      []uint16    `json:"subHeaders"`
	GlyphIndexArray []uint16    `json:"glyphIndexArray"`
}

type CmapFormat4 struct {
	Format          uint16   `json:"format"`
	Length          uint16   `json:"length"`
	Language        uint16   `json:"language"`
	SegCountX2      uint16   `json:"segCountx2"`
	SearchRange     uint16   `json:"searchRange"`
	EntrySelector   uint16   `json:"entrySelector"`
	RangeShift      uint16   `json:"rangeShift"`
	EndCode         []uint16 `json:"endCode"`
	ReservedPad     uint16   `json:"reservedPad"`
	StartCode       []uint16 `json:"startCode"`
	IdDelta         []uint16 `json:"idDelta"`
	IdRangeOffset   []uint16 `json:"idRangeOffset"`
	GlyphIndexArray []uint16 `json:"glyphIndexArray"`
}

type CmapFormat6 struct {
	Format          uint16   `json:"format"`
	Length          uint16   `json:"length"`
	Language        uint16   `json:"language"`
	FirstCode       uint16   `json:"firstCode"`
	EntryCount      uint16   `json:"entryCount"`
	GlyphIndexArray []uint16 `json:"glyphIndexArray"`
}

type NGrups struct {
	StartCharCode uint32 `json:"startCharCode"`
	EndCharCode uint32 `json:"endCharCode"`
	StartGlyphCode uint32 `json:"startGlyphCode"`
}
type CmapFormat8 struct {
	Format          uint16   `json:"format"`
	Reserved 				uint16 		`json:"reserved"`
	Length          uint16   `json:"length"`
	Language        uint16   `json:"language"`
	Is32       			[65536]uint8   `json:"is32"`
	NGroups      		NGrups   `json:"nGroups"`
}

type CmapFormat10 struct {
	Format          uint16   `json:"format"`
	Reserved 				uint16 		`json:"reserved"`
	Length          uint32   `json:"length"`
	Language        uint32   `json:"language"`
	StartCharCode 	uint32 `json:"startCharCode"`
	NumChars 				uint32 `json:"numChars"`
	Glyphs  []uint32  `json:"glyphs"`
}

type CmapFormat12 struct {
	Format          uint16   `json:"format"`
	Reserved 				uint16 		`json:"reserved"`
	Length          uint32   `json:"length"`
	Language        uint32   `json:"language"`
	NGroups      		NGrups   `json:"nGroups"`
}

type CmapFormat13 struct {
	StartCharCode uint32 `json:"startCharCode"`
	EndCharCode uint32 `json:"endCharCode"`
	GlyphCode uint32 `json:"glyphCode"`
}

type CmapFormat14 struct {
	Format          uint16   `json:"format"`
	Length          uint32   `json:"length"`
	NumVarSelectorRecords uint32 `json:"numVarSelectorRecords"`
}

type CmapFormat2SubHeader struct {
	FirstCode uint16 `json:"firstCode"`
	EntryCount uint16 `json:"entryCount"`
	IdDelta int16 `json:"idDelta"`
	IdRangeOffset uint16 `json:"idRangeOffset"`
}

type CmapFormat8nGroup struct {
	StartCharCode uint32 `json:"startCharCode"`
	EndCharCode uint32 `json:"endCharCode"`
	StartGlyphCode uint32 `json:"startGlyphCode"`
}

type CmapFormatDefaultUVS struct {
	StartUnicode int `json:"startUnicode"`
	EndUnicode int `json:"endUnicode"`
	VarSelector uint32 `json:"varSelector"`
}

type CmapFormatNonDefaultUVS struct {
	UnicodeValue int `json:"unicodeValue"`
	GlyphID uint16 `json:"glyphId"`
	VarSelector uint32 `json:"varSelector"`
}

func GetCmap(data []byte) (cmap *Cmap,err error) {
	pos := 0
	cmap = new(Cmap)
	cmap.Version = getUint16(data[0:pos+2])
	pos += 2
	cmap.NumberSubtables = getUint16(data[pos:pos+2])
	pos += 2

	for i := 0; i < int(cmap.NumberSubtables); i++ {
		subTable := make(map[string] interface{})
		subTable["platformID"] = getUint16(data[pos:pos+2])
		pos += 2
		subTable["platformSpecificID"] = getUint16(data[pos:pos+2])
		pos+=2
		subTable["offset"] = getUint32(data[pos:pos+4])
		pos+=4

		startPos := pos

		subTable["format"] = getUint16(data[pos:pos+2])
		pos+=2
		format, ok := subTable["format"].(int)
		if !ok {
			err = errors.New("cmap format int error")
			return
		}
		if format == 0 {
			subTable["length"] = getUint16(data[pos:pos+2])
			pos+=2
			subTable["language"] = getUint16(data[pos:pos+2])
			pos+=2

			len, ok := subTable["length"].(int)
			if !ok {
				err = errors.New("cmap format0 length error")
				return
			}
			var glyphIndexArray []uint8
			for i := 0; i < len; i++ {
				glyphIndexArray = append(glyphIndexArray, getUint8(data[pos:pos+1]))
				pos++
			}
			subTable["glyphIndexArray"] = glyphIndexArray
			
		} else if format == 2 {
			subTable["length"] = getUint16(data[pos:pos+2])
			pos+=2
			subTable["language"] = getUint16(data[pos:pos+2])
			pos+=2

			var subHeaderKeys []uint16
			maxSubHeaderKey := 0
			for i := 0; i < 256; i++ {
				sourceVal := getUint16(data[pos:pos+2])
				subHeaderKeys = append(subHeaderKeys, sourceVal)
				pos+=2
				val := int(sourceVal)/8

				if (val > maxSubHeaderKey) {
					maxSubHeaderKey = val
				}
			}
			subTable["subHeaderKeys"] = subHeaderKeys

			var subHeaders []*CmapFormat2SubHeader
			for k := 0; k < maxSubHeaderKey; k++ {
				subHeaders = append(subHeaders, &CmapFormat2SubHeader{
					getUint16(data[pos:pos+2]),
					getUint16(data[pos+2:pos+4]),
					getInt16(data[pos+4:pos+6]),
					getUint16(data[pos+6:pos+8]),
				})
				pos+=8
			}
			subTable["subHeaders"] = subHeaders
			subTableLen, ok := subTable["length"].(int)
			if !ok {
				return
			}

			glyphIndexArrayLen := (startPos + subTableLen - pos)/2
			var glyphIndexArray []uint16
			for i := 0; i < glyphIndexArrayLen; i++ {
				glyphIndexArray = append(glyphIndexArray, getUint16(data[pos:pos+2]))
				pos += 2
			}
			subTable["glyphIndexArray"] = glyphIndexArray

		} else if format == 4 {
			subTable["length"] = getUint16(data[pos:pos+2])
			pos+=2
			subTable["language"] = getUint16(data[pos:pos+2])
			pos+=2
			subTable["segCountX2"] = getUint16(data[pos:pos+2])
			pos+=2
			subTable["searchRange"] = getUint16(data[pos:pos+2])
			pos+=2
			subTable["entrySelector"] = getUint16(data[pos:pos+2])
			pos+=2
			subTable["rangeShift"] = getUint16(data[pos:pos+2])
			pos+=2

			segCount := subTable["segCountX2"].(int) / 2

			var endCode []uint16
			for i := 0; i < segCount; i++ {
				endCode = append(endCode, getUint16(data[pos:pos+2]))
				pos+=2
			}
			subTable["endCode"] = endCode
			subTable["reservedPad"] = getUint16(data[pos:pos+2])
			pos+=2

			var startCode []uint16
			for i := 0; i < segCount; i++ {
				startCode = append(startCode, getUint16(data[0:pos+2]))
				pos+=2
			}
			subTable["startCode"] = startCode

			var idDelta []uint16
			for i := 0; i < segCount; i++ {
				idDelta = append(idDelta, getUint16(data[0:2]))
				pos+=2
			}
			subTable["idDelta"] = idDelta

			var idRangeOffset []uint16
			for i := 0; i < segCount; i++ {
				idRangeOffset = append(idRangeOffset, getUint16(data[pos:pos+2]))
				pos+=2
			}
			subTable["idRangeOffset"] = idRangeOffset

			// The remaining is glyphIndexArray length
			glyphLen := (subTable["length"].(int) - (pos - startPos))/2
			var glyphIndexArray []uint16
			for i := 0; i < glyphLen; i++ {
				glyphIndexArray = append(glyphIndexArray, getUint16(data[pos:pos+2]))
				pos+=2
			}
			subTable["glyphIndexArray"] = glyphIndexArray

		} else if format == 6 {
			subTable["length"] = getUint16(data[pos:pos+2])
			pos+=2
			subTable["language"] = getUint16(data[pos:pos+2])
			pos+=2
			subTable["firstCode"] = getUint16(data[pos:pos+2])
			pos+=2
			subTable["entryCount"] = getUint16(data[pos:pos+2])
			pos+=2

			var glyphIndexArray []uint16
			entryCount := subTable["entryCount"].(int)
			for i := 0; i < entryCount; i++ {
				glyphIndexArray = append(glyphIndexArray, getUint16(data[pos:pos+2]))
				pos+=2
			}
			subTable["glyphIndexArray"] = glyphIndexArray
		} else if format == 8 {
			subTable["reserved"] = getUint16(data[pos:pos+2])
			pos+=2
			subTable["length"] = getUint32(data[pos:pos+4])
			pos+=4
			subTable["language"] = getUint32(data[pos:pos+4])
			pos+=4
			var is32 []uint8
			for i := 0; i < 65536; i++ {
				is32 = append(is32, getUint8(data[pos:pos+1]))
				pos++
			}
			subTable["is32"] = is32

			// n := (subTable["length"].(int) - (pos - startPos))/12
			subTable["nGroups"] = getUint32(data[pos:pos+4])
			pos+=4
			n := subTable["nGroups"].(int)
			var groups []*CmapFormat8nGroup
			for i := 0; i < n; i++ {
				groups = append(groups, &CmapFormat8nGroup{
					getUint32(data[pos:pos+4]),
					getUint32(data[pos+4:pos+8]),
					getUint32(data[pos+8:pos+12]),
				})
				pos+=12
			}
			subTable["groups"] = groups

		} else if format == 10 {
			subTable["reserved"] = getUint16(data[pos:pos+2])
			pos+=2
			subTable["length"] = getUint32(data[pos:pos+4])
			pos+=4
			subTable["language"] = getUint32(data[pos:pos+4])
			pos+=4
			subTable["startCharCode"] = getUint32(data[pos:pos+4])
			pos+=4
			subTable["numChars"] = getUint32(data[pos:pos+4])
			pos+=4
			numChars := subTable["numChars"].(int)

			var glyphs []uint16
			for i := 0; i < numChars; i++ {
				glyphs = append(glyphs, getUint16(data[pos:pos+2]))
				pos+=2
			}
			subTable["glyphs"] = glyphs
		} else if format == 12 || format == 13 {
			subTable["reserved"] = getUint16(data[pos:pos+2])
			pos+=2
			subTable["length"] = getUint32(data[pos:pos+4])
			pos+=4
			subTable["language"] = getUint32(data[pos:pos+4])
			pos+=4
			subTable["nGroups"] = getUint32(data[pos:pos+4])
			pos+=4

			n := subTable["nGroups"].(int)
			var groups []*CmapFormat8nGroup
			for i := 0; i < n; i++ {
				groups = append(groups, &CmapFormat8nGroup{
					getUint32(data[pos:pos+4]),
					getUint32(data[pos+4:pos+8]),
					getUint32(data[pos+8:pos+12]),
				})
				pos+=12
			}
			subTable["groups"] = groups

		} else if format == 14 {
			subTable["length"] = getUint32(data[pos:pos+4])
			pos+=4
			subTable["numVarSelectorRecords"] = getUint32(data[pos:pos+4])
			pos+=4

			n := subTable["numVarSelectorRecords"].(int)
			var groups []interface{}
			for i := 0; i < n; i++ {
				var varSelector uint32
				varSelector, err = getUint24(data[pos:pos+3])
				pos+=3
				if err != nil {
					return
				}
				defaultUVSOffset := int(getUint32(data[pos:pos+4]))
				pos+=4
				nonDefaultUVSOffset := int(getUint32(data[pos:pos+4]))
				pos+=4

				if defaultUVSOffset != 0 {
					numUnicodeValueRanges := int(getUint32(data[pos+defaultUVSOffset:pos+defaultUVSOffset+4]))

					for i := 0; i < numUnicodeValueRanges; i++ {
						var startUnicode uint32
						startUnicode, err = getUint24(data[pos:pos+3])
						pos+=3
						if err != nil {
							return
						}
						start := int(startUnicode)
						additionalCount := int(getUint8(data[pos:pos+1]))
						pos++
						end := start+additionalCount
						groups = append(groups, &CmapFormatDefaultUVS{
							start,
							end,
							varSelector, 
						})
					}
				}

				if nonDefaultUVSOffset != 0 {
					numUVSMappings := int(getUint32(data[startPos+nonDefaultUVSOffset:startPos+nonDefaultUVSOffset+4]))

					for i := 0; i < numUVSMappings; i++ {
						var v uint32
						v, err = getUint24(data[pos:pos+3])
						pos+=3
						if err != nil {
							return
						}
						groups = append(groups, &CmapFormatNonDefaultUVS{
							int(v),
							getUint16(data[pos:pos+4]),
							varSelector,
						})
					}
				}


			}

			subTable["groups"] = groups
		} else {
			err = errors.New("format not support!")
			return
		}
		cmap.SubTables = append(cmap.SubTables, subTable)
	}

	return
}