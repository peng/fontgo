package font

import (
	"errors"
	"strconv"
)

type OffsetTable struct {
	ScalerType    string `json:"scalerType"`
	NumTables     uint16 `json:"numTables"`
	SearchRange   uint16 `json:"searchRange"`
	EntrySelector uint16 `json:"entrySelector"`
	RangeShift    uint16 `json:"rangeShift"`
}

func GetScalerType(data []byte) string {
	n := int(getUint32(data[0:4]))
	if n == 65536 || n == 1953658213 {
		return "TrueType"
	} else if n == 1954115633 {
		return "typ1"
	} else if n == 1330926671 {
		return "OTTO"
	}
	return ""
}

func GetOffsetTable(data []byte) *OffsetTable {
	return &OffsetTable{
		GetScalerType(data[0:4]),
		getUint16(data[4:6]),
		getUint16(data[6:8]),
		getUint16(data[8:10]),
		getUint16(data[10:12]),
	}
}

type TableContent map[string]*TagItem

func GetTableContent(numTables int, date []byte) TableContent {
	tableContent := make(TableContent)
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

func GetHead(data []byte, pos int) *Head {
	return &Head{
		getFixed(data[pos : pos+4]),
		getFixed(data[pos+4 : pos+8]),
		getUint32(data[pos+8 : pos+12]),
		getUint32(data[pos+12 : pos+16]),
		getUint16(data[pos+16 : pos+18]),
		getUint16(data[pos+18 : pos+20]),
		getLongDateTime(data[pos+20 : pos+28]),
		getLongDateTime(data[pos+28 : pos+36]),
		getFWord(data[pos+36 : pos+38]),
		getFWord(data[pos+38 : pos+40]),
		getFWord(data[pos+40 : pos+42]),
		getFWord(data[pos+42 : pos+44]),
		getUint16(data[pos+44 : pos+46]),
		getUint16(data[pos+46 : pos+48]),
		getInt16(data[pos+48 : pos+50]),
		getInt16(data[pos+50 : pos+52]),
		getInt16(data[pos+52 : pos+54]),
	}
}

type Flag struct {
	OnCurve      bool `json:"onCurve"`
	XShortVector bool `json:"xShortVector`
	YShortVector bool `json:"yShortVector"`
	Repeat       bool `json:"repeat"`
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
	Flags     uint16 `json:"flags"`
	Argument1 int    `json:"argument1"`
	Argument2 int    `json:"argument2"`
	// Unsign     bool    `json:"unsign"`
	// Scale   float32 `json:"scale"`
	Xscale  float32 `json:"xscale"`
	Yscale  float32 `json:"yscale"`
	Scale01 float32 `json:"scale01"`
	Scale10 float32 `json:"scale10"`
}

type GlyphCompound struct {
	GlyphCommon
	Component         []Component `json:"component"`
	InstructionLength int         `json:"instructionLength"`
	Instructions      []uint8     `json:"instructions"`
}

type Glyphs struct {
	Simples   []GlyphSimple   `json:"simples"`
	Compounds []GlyphCompound `json:"compounds"`
}

const GLYPH_TYPE_SIMPLE, GLYPH_TYPE_COMPOUND = "simple", "compound"

func GetGlyphSimple(data []byte, pos int) (simple *GlyphSimple) {
	simple = new(GlyphSimple)
	simple.GlyphCommon.Type = GLYPH_TYPE_SIMPLE
	simple.GlyphCommon.NumberOfContours = getInt16(data[pos : pos+2])
	simple.GlyphCommon.XMin = getFWord(data[pos+2 : pos+4])
	simple.GlyphCommon.YMin = getFWord(data[pos+4 : pos+6])
	simple.GlyphCommon.XMax = getFWord(data[pos+6 : pos+8])
	simple.GlyphCommon.YMax = getFWord(data[pos+8 : pos+10])

	pos += 10
	// get endPtsOfContours
	for i := 0; i < int(simple.GlyphCommon.NumberOfContours); i++ {
		simple.EndPtsOfContours = append(simple.EndPtsOfContours, getUint16(data[pos:pos+2]))
		pos += 2
	}

	// get instructionLength
	simple.InstructionLength = getUint16(data[pos : pos+2])
	pos += 2
	for i := 0; i < int(simple.InstructionLength); i++ {
		simple.Instructions = append(simple.Instructions, getUint8(data[pos:pos+1]))
		// test pos++
		pos++
	}

	// get points num
	pointsNum := 0

	for _, num := range simple.EndPtsOfContours {
		contoursNum := int(num)
		if contoursNum > pointsNum {
			pointsNum = contoursNum
		}
	}
	pointsNum += 1
	var flags []*Flag

	// get flags
	for i := 0; i < pointsNum; i++ {
		f := getUint8(data[pos : pos+1])
		flag := &Flag{
			(f & 0x01) == 0x01,
			(f & 0x02) == 0x02,
			(f & 0x04) == 0x04,
			(f & 0x08) == 0x08,
			(f & 0x10) == 0x10,
			(f & 0x20) == 0x20,
		}
		flags = append(flags, flag)
		pos++

		if flag.Repeat {
			repeatNum := int(getUint8(data[pos : pos+1]))
			pos++

			for j := 0; j < repeatNum; j++ {
				flags = append(flags, flag)
			}
			i += repeatNum
		}
	}

	// should check number of flags same with points

	// get x points
	for i := 0; i < pointsNum; i++ {

		flag := flags[i]

		var (
			point Point
			x     int
		)
		point.Flag = flag
		if flag.XShortVector {
			x = int(getUint8(data[pos : pos+1]))
			if !flag.XSame {
				x *= -1
			}
			pos++
		} else if flag.XSame {
			x = 0
		} else {
			x = int(getInt16(data[pos : pos+2]))
			pos += 2
		}
		point.X = x

		simple.Points = append(simple.Points, &point)
	}

	// get y points
	for i := 0; i < pointsNum; i++ {
		var y int
		point := simple.Points[i]
		flag := point.Flag
		if flag.YShortVector {
			y = int(getUint8(data[pos : pos+1]))
			if !flag.YSame {
				y *= -1
			}
			pos++
		} else if flag.YSame {
			y = 0
		} else {
			y = int(getUint16(data[pos : pos+2]))
			pos += 2
		}
		point.Y = y
	}

	return
}

func GetGlyphCompound(data []byte, pos int) (compound *GlyphCompound) {
	compound = new(GlyphCompound)
	const (
		ARG_1_AND_2_ARE_WORDS    uint16 = 0x0001
		ARGS_ARE_XY_VALUES       uint16 = 0x0002
		ROUND_XY_TO_GRID         uint16 = 0x0004
		WE_HAVE_A_SCALE          uint16 = 0x0008
		MORE_COMPONENTS          uint16 = 0x0020
		WE_HAVE_AN_X_AND_Y_SCALE uint16 = 0x0040
		WE_HAVE_A_TWO_BY_TWO     uint16 = 0x0080
		WE_HAVE_INSTRUCTIONS     uint16 = 0x0100
		USE_MY_METRICS           uint16 = 0x0200
		OVERLAP_COMPOUND         uint16 = 0x0400
	)

	compound.Type = GLYPH_TYPE_COMPOUND
	compound.NumberOfContours = getInt16(data[pos : pos+2])
	compound.XMin = getFWord(data[pos+2 : pos+4])
	compound.YMin = getFWord(data[pos+4 : pos+6])
	compound.XMax = getFWord(data[pos+6 : pos+8])
	compound.YMax = getFWord(data[pos+8 : pos+10])

	var flags uint16
	pos += 10
	moreComponent := true

	for moreComponent {
		component := new(Component)

		flags = getUint16(data[pos : pos+2])
		component.Flags = flags
		pos += 2

		if (flags & ARG_1_AND_2_ARE_WORDS) == ARG_1_AND_2_ARE_WORDS {
			if (flags & ARGS_ARE_XY_VALUES) == ARGS_ARE_XY_VALUES {
				component.Argument1 = int(getInt16(data[pos : pos+2]))
				pos += 2
				component.Argument2 = int(getInt16(data[pos : pos+2]))
			} else {
				// component.Unsign = true
				component.Argument1 = int(getUint16(data[pos : pos+2]))
				pos += 2
				component.Argument2 = int(getUint16(data[pos : pos+2]))
			}
			pos += 2
		} else {
			if flags&ARGS_ARE_XY_VALUES == ARGS_ARE_XY_VALUES {
				component.Argument1 = int(getInt8(data[pos : pos+1]))
				pos++
				component.Argument2 = int(getInt8(data[pos : pos+1]))
			} else {
				// component.Unsign = true
				component.Argument1 = int(getUint8(data[pos : pos+1]))
				pos++
				component.Argument2 = int(getUint8(data[pos : pos+1]))
			}
			pos++
		}

		if flags&WE_HAVE_A_SCALE == WE_HAVE_A_SCALE {
			// component.Scale = get2Dot14(data[pos : pos+2])
			component.Xscale = get2Dot14(data[pos : pos+2])
			component.Yscale = component.Xscale
		} else if flags&WE_HAVE_AN_X_AND_Y_SCALE == WE_HAVE_AN_X_AND_Y_SCALE {
			component.Xscale = get2Dot14((data[pos : pos+2]))
			pos += 2
			component.Yscale = get2Dot14(data[pos : pos+2])
		} else if flags&WE_HAVE_A_TWO_BY_TWO == WE_HAVE_A_TWO_BY_TWO {
			component.Xscale = get2Dot14(data[pos : pos+2])
			pos += 2
			component.Scale01 = get2Dot14(data[pos : pos+2])
			pos += 2
			component.Scale10 = get2Dot14(data[pos : pos+2])
			pos += 2
			component.Xscale = get2Dot14(data[pos : pos+2])
		}
		pos += 2

		compound.Component = append(compound.Component, *component)

		moreComponent = flags&MORE_COMPONENTS == MORE_COMPONENTS
	}

	// can't understand
	if flags&WE_HAVE_INSTRUCTIONS == WE_HAVE_INSTRUCTIONS {
		compound.InstructionLength = int(getUint16(data[pos : pos+2]))
		pos += 2

		for i := 0; i < compound.InstructionLength; i++ {
			compound.Instructions = append(compound.Instructions, getUint8(data[pos:pos+1]))
			pos++
		}
	}

	return
}

func GetGlyphs(data []byte, pos int, loca []int, numGlyphs int) (glyphs *Glyphs) {
	glyphs = new(Glyphs)

	for i := 0; i < numGlyphs; i++ {

		offset := loca[i]
		// fmt.Printf("innoffset %v", offset)
		// fmt.Printf("innnextoffset %v", nextOffset)
		inPos := offset + pos
		numberOfContours := getInt16(data[inPos : inPos+2])

		if numberOfContours >= 0 {
			// fmt.Println("numberOfContours", numberOfContours)
			// simple
			simp := GetGlyphSimple(data, inPos)
			glyphs.Simples = append(glyphs.Simples, *simp)
		} else {
			// compound
			compound := GetGlyphCompound(data, inPos)
			glyphs.Compounds = append(glyphs.Compounds, *compound)
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

func GetMaxp(data []byte, pos int) *Maxp {
	maxp := new(Maxp)
	maxp.Version = getVersion(data[pos : pos+4])
	maxp.NumGlyphs = getUint16(data[pos+4 : pos+6])
	if maxp.Version == "1.0" {
		maxp.MaxPoints = getUint16(data[pos+6 : pos+8])
		maxp.MaxContours = getUint16(data[pos+8 : pos+10])
		maxp.MaxComponentPoints = getUint16(data[pos+10 : pos+12])
		maxp.MaxComponentContours = getUint16(data[pos+12 : pos+14])
		maxp.MaxZones = getUint16(data[pos+14 : pos+16])
		maxp.MaxTwilightPoints = getUint16(data[pos+16 : pos+18])
		maxp.MaxStorage = getUint16(data[pos+18 : pos+20])
		maxp.MaxFunctionDefs = getUint16(data[pos+20 : pos+22])
		maxp.MaxInstructionDefs = getUint16(data[pos+22 : pos+24])
		maxp.MaxStackElements = getUint16(data[pos+24 : pos+26])
		maxp.MaxSizeOfInstructions = getUint16(data[pos+26 : pos+28])
		maxp.MaxComponentElements = getUint16(data[pos+28 : pos+30])
		maxp.MaxComponentDepth = getUint16(data[pos+30 : pos+32])
	}

	return maxp
}

func GetLoca(data []byte, pos int, numGlyphs uint16, indexToLocFormat int16) (locations []int) {
	// long version:  otf, ttf is different
	offsetFn := func(data []byte, pos int) (offset int, nextPos int) {
		size := 2
		ratio := 2
		if indexToLocFormat == 0 {
			offset = int(getUint16(data[pos:pos+size])) * ratio
			nextPos = pos + size
			return
		}
		size = 4
		ratio = 1
		offset = int(getUint32(data[pos:pos+size])) * ratio
		nextPos = pos + size
		return
	}
	var (
		offset int
	)
	for i := 0; i < int(numGlyphs); i++ {
		offset, pos = offsetFn(data, pos)
		locations = append(locations, offset)
	}

	return locations
}

type Cmap struct {
	Version         uint16                   `json:"version"`
	NumberSubtables uint16                   `json:"numberSubtables"`
	SubTables       []map[string]interface{} `json:"subTables"`
	WindowsCode     map[int]int
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
	StartCharCode  uint32 `json:"startCharCode"`
	EndCharCode    uint32 `json:"endCharCode"`
	StartGlyphCode uint32 `json:"startGlyphCode"`
}
type CmapFormat8 struct {
	Format   uint16       `json:"format"`
	Reserved uint16       `json:"reserved"`
	Length   uint16       `json:"length"`
	Language uint16       `json:"language"`
	Is32     [65536]uint8 `json:"is32"`
	NGroups  NGrups       `json:"nGroups"`
}

type CmapFormat10 struct {
	Format        uint16   `json:"format"`
	Reserved      uint16   `json:"reserved"`
	Length        uint32   `json:"length"`
	Language      uint32   `json:"language"`
	StartCharCode uint32   `json:"startCharCode"`
	NumChars      uint32   `json:"numChars"`
	Glyphs        []uint32 `json:"glyphs"`
}

type CmapFormat12 struct {
	Format   uint16 `json:"format"`
	Reserved uint16 `json:"reserved"`
	Length   uint32 `json:"length"`
	Language uint32 `json:"language"`
	NGroups  NGrups `json:"nGroups"`
}

type CmapFormat13 struct {
	StartCharCode uint32 `json:"startCharCode"`
	EndCharCode   uint32 `json:"endCharCode"`
	GlyphCode     uint32 `json:"glyphCode"`
}

type CmapFormat14 struct {
	Format                uint16 `json:"format"`
	Length                uint32 `json:"length"`
	NumVarSelectorRecords uint32 `json:"numVarSelectorRecords"`
}

type CmapFormat2SubHeader struct {
	FirstCode     uint16 `json:"firstCode"`
	EntryCount    uint16 `json:"entryCount"`
	IdDelta       int16  `json:"idDelta"`
	IdRangeOffset uint16 `json:"idRangeOffset"`
}

type CmapFormat8nGroup struct {
	StartCharCode  uint32 `json:"startCharCode"`
	EndCharCode    uint32 `json:"endCharCode"`
	StartGlyphCode uint32 `json:"startGlyphCode"`
}

type CmapFormatDefaultUVS struct {
	StartUnicode int    `json:"startUnicode"`
	EndUnicode   int    `json:"endUnicode"`
	VarSelector  uint32 `json:"varSelector"`
}

type CmapFormatNonDefaultUVS struct {
	UnicodeValue int    `json:"unicodeValue"`
	GlyphID      uint16 `json:"glyphId"`
	VarSelector  uint32 `json:"varSelector"`
}

func readWindowsCode(subTables []map[string]interface{}, maxpNumGlyphs int) (code map[int]int, err error) {
	code = make(map[int]int)

	var format0, format2, format4, format12, format14 map[string]interface{}

	for _, val := range subTables {
		formatSource, exist := val["format"]
		platformIDSource, exist2 := val["platformID"]
		platformSpecificIDSource, exist3 := val["platformSpecificID"]

		if !exist || !exist2 || !exist3 {
			err = errors.New("Read platformID or platformSpecificID error")
			return
		}
		format := formatSource.(int)
		platformID := platformIDSource.(int)
		platformSpecificID := platformSpecificIDSource.(int)
		// https://learn.microsoft.com/en-us/typography/opentype/spec/recom#cmap-table
		if format == 0 {
			format0 = val
		} else if format == 2 && platformID == 3 && platformSpecificID == 3 {
			format2 = val
		} else if format == 4 && platformID == 3 && platformSpecificID == 1 {
			format4 = val
		} else if format == 12 && platformID == 3 && platformSpecificID == 10 {
			format12 = val
		} else if format == 14 && platformID == 0 && platformSpecificID == 5 {
			format14 = val
		}
	}

	if len(format0) > 0 {
		g, exist := format0["glyphIndexArray"]
		if exist {
			glyphIndexArray := g.([]uint8)
			for ind, val := range glyphIndexArray {
				if int(val) != 0 {
					code[ind] = int(val)
				}
			}
		}
	}

	if len(format14) > 0 {
		groupsSource, exist := format14["groups"]
		if exist {
			groups := groupsSource.([]interface{})
			for _, sVal := range groups {
				val := sVal.([]interface{})
				typeVal := val[0].(int)
				if typeVal == 1 {
					nonDefaultUVS := val[1].(CmapFormatNonDefaultUVS)
					code[nonDefaultUVS.UnicodeValue] = int(nonDefaultUVS.GlyphID)
				}
			}
		}
	}

	if len(format12) > 0 {
		gSource, exist := format12["nGroups"]
		if !exist {
			err = errors.New("Read format12 nGroups error")
			return
		}
		groups := gSource.([]*CmapFormat8nGroup)
		for _, val := range groups {
			startCharCode := int(val.StartCharCode)
			endCharCode := int(val.EndCharCode)
			startGlyphCode := int(val.StartGlyphCode)

			for startCharCode <= endCharCode {
				code[startCharCode] = startGlyphCode
				startCharCode++
				startGlyphCode++
			}
		}
	} else if len(format4) > 0 {
		segCountX2, exist1 := format4["segCountX2"]
		startCodeSource, exist2 := format4["startCode"]
		endCodeSource, exist3 := format4["endCode"]
		idRangeOffsetSource, exist4 := format4["idRangeOffset"]
		idDeltaSource, exist5 := format4["idDelta"]
		glyphIndexArraySource, exist6 := format4["glyphIndexArray"]
		glyphIndexArrayOffsetSource, exist7 := format4["glyphIndexArrayOffset"]
		idRangeOffsetOffsetSource, exist8 := format4["idRangeOffsetOffset"]

		if !exist1 || !exist2 || !exist3 || !exist4 || !exist5 || !exist6 || !exist7 || !exist8 {
			err = errors.New("Read format4 map error")
			return
		}

		segCount := int(segCountX2.(uint16)) / 2
		startCode := startCodeSource.([]uint16)
		endCode := endCodeSource.([]uint16)
		idRangeOffset := idRangeOffsetSource.([]uint16)
		idDelta := idDeltaSource.([]uint16)
		glyphIndexArray := glyphIndexArraySource.([]uint16)
		glyphIndexArrayOffset := glyphIndexArrayOffsetSource.(int)
		idRangeOffsetOffset := idRangeOffsetOffsetSource.(int)

		// Calculate graphIdArrayIndexOffset like in JavaScript
		graphIdArrayIndexOffset := (glyphIndexArrayOffset - idRangeOffsetOffset) / 2

		for i := 0; i < segCount; i++ {
			for start, end := int(startCode[i]), int(endCode[i]); start <= end; start++ {
				if int(idRangeOffset[i]) == 0 {
					code[start] = (start + int(idDelta[i])) % 0x10000
				} else {
					// Calculate index exactly as in JavaScript reference code
					index := i + int(idRangeOffset[i])/2 + (start - int(startCode[i])) - graphIdArrayIndexOffset

					// Boundary check to prevent panic
					if index >= 0 && index < len(glyphIndexArray) {
						glyphIndex := int(glyphIndexArray[index])

						if glyphIndex != 0 {
							code[start] = (glyphIndex + int(idDelta[i])) % 0x10000
						} else {
							code[start] = 0
						}
					}
				}
			}
		}

		// wip 65535
		delete(code, 65535)
	} else if len(format2) > 0 {
		subHeaderKeysS, exist1 := format2["subHeaderKeys"]
		subHeadersS, exist2 := format2["subHeaders"]
		glyphIndexArrayS, exist3 := format2["glyphIndexArray"]
		maxPosS, exist4 := format2["maxPos"]

		if !exist1 || !exist2 || !exist3 || !exist4 {
			err = errors.New("Read format2 error")
			return
		}

		subHeaderKeys := subHeaderKeysS.([]uint16)
		subHeaders := subHeadersS.([]*CmapFormat2SubHeader)
		glyphIndexArray := glyphIndexArrayS.([]uint16)
		maxPos := maxPosS.(int)

		index := 0

		for i := 0; i < 256; i++ {
			if i >= len(subHeaderKeys) {
				break
			}
			k := int(subHeaderKeys[i])
			if k == 0 {
				if len(subHeaders) == 0 {
					continue
				}

				idxPos := int(subHeaders[0].IdRangeOffset) + (i - int(subHeaders[0].FirstCode))
				if i >= maxPos || i < int(subHeaders[0].FirstCode) || i >= int(subHeaders[0].FirstCode+subHeaders[0].EntryCount) || idxPos < 0 || idxPos >= len(glyphIndexArray) {
					index = 0
				} else {
					index = int(glyphIndexArray[idxPos])
					if index != 0 {
						index = index + int(subHeaders[0].IdDelta)
					}
				}

				if index != 0 && index < maxpNumGlyphs {
					code[i] = index
				}

			} else {
				if k >= len(subHeaders) {
					continue
				}
				entryCount := int(subHeaders[k].EntryCount)
				for j := 0; j < entryCount; j++ {
					idxPos := int(subHeaders[k].IdRangeOffset) + j

					if idxPos < 0 || idxPos >= len(glyphIndexArray) {
						index = 0
					} else {
						index = int(glyphIndexArray[idxPos])
						if index != 0 {
							index = index + int(subHeaders[k].IdDelta)
						}
					}

					if index != 0 && index < maxpNumGlyphs {
						unicode := ((i << 8) | (j + int(subHeaders[k].FirstCode))) % 0xffff
						code[unicode] = index
					}
				}
			}
		}
	}

	return
}

func GetCmap(data []byte, pos int, maxpNumGlyphs int) (cmap *Cmap, err error) {
	cmap = new(Cmap)
	startPos := pos
	cmap.Version = getUint16(data[pos : pos+2])
	pos += 2
	cmap.NumberSubtables = getUint16(data[pos : pos+2])
	pos += 2

	for i := 0; i < int(cmap.NumberSubtables); i++ {
		subTable := make(map[string]interface{})
		subTable["platformID"] = int(getUint16(data[pos : pos+2]))
		pos += 2
		subTable["platformSpecificID"] = int(getUint16(data[pos : pos+2]))
		pos += 2
		subTable["offset"] = getUint32(data[pos : pos+4])
		pos += 4

		// startPos := pos
		curPos := startPos + int(subTable["offset"].(uint32))
		startOffset := curPos

		subTable["format"] = int(getUint16(data[curPos : curPos+2]))
		curPos += 2
		format, ok := subTable["format"].(int)
		if !ok {
			err = errors.New("cmap format int error")
			return
		}
		if format == 0 {
			subTable["length"] = int(getUint16(data[curPos : curPos+2]))
			curPos += 2
			subTable["language"] = getUint16(data[curPos : curPos+2])
			curPos += 2

			var glyphIndexArray []uint8
			for i := 0; i < 256; i++ {
				glyphIndexArray = append(glyphIndexArray, getUint8(data[curPos:curPos+1]))
				curPos++
			}
			subTable["glyphIndexArray"] = glyphIndexArray

		} else if format == 2 {
			subTable["length"] = int(getUint16(data[curPos : curPos+2]))
			curPos += 2
			subTable["language"] = getUint16(data[curPos : curPos+2])
			curPos += 2

			var subHeaderKeys []uint16
			maxSubHeaderKey := 0
			maxPos := -1
			for i := 0; i < 256; i++ {
				sourceVal := getUint16(data[curPos : curPos+2])
				subHeaderKeys = append(subHeaderKeys, sourceVal)
				curPos += 2
				val := int(sourceVal) / 8

				if val > maxSubHeaderKey {
					maxSubHeaderKey = val
					maxPos = i
				}
			}
			subTable["subHeaderKeys"] = subHeaderKeys
			subTable["maxPos"] = maxPos

			var subHeaders []*CmapFormat2SubHeader
			for k := 0; k < maxSubHeaderKey; k++ {
				subHeaders = append(subHeaders, &CmapFormat2SubHeader{
					getUint16(data[curPos : curPos+2]),
					getUint16(data[curPos+2 : curPos+4]),
					getInt16(data[curPos+4 : curPos+6]),
					getUint16(data[curPos+6 : curPos+8]),
				})
				curPos += 8
			}
			subTable["subHeaders"] = subHeaders
			subTableLen, ok := subTable["length"].(int)
			if !ok {
				return
			}

			glyphIndexArrayLen := (startPos + subTableLen - curPos) / 2
			var glyphIndexArray []uint16
			for i := 0; i < glyphIndexArrayLen; i++ {
				glyphIndexArray = append(glyphIndexArray, getUint16(data[curPos:curPos+2]))
				curPos += 2
			}
			subTable["glyphIndexArray"] = glyphIndexArray

		} else if format == 4 {
			subTable["length"] = int(getUint16(data[curPos : curPos+2]))
			curPos += 2
			subTable["language"] = getUint16(data[curPos : curPos+2])
			curPos += 2
			subTable["segCountX2"] = getUint16(data[curPos : curPos+2])
			curPos += 2
			subTable["searchRange"] = getUint16(data[curPos : curPos+2])
			curPos += 2
			subTable["entrySelector"] = getUint16(data[curPos : curPos+2])
			curPos += 2
			subTable["rangeShift"] = getUint16(data[curPos : curPos+2])
			curPos += 2

			segCount := int(subTable["segCountX2"].(uint16)) / 2

			var endCode []uint16
			for i := 0; i < segCount; i++ {
				endCode = append(endCode, getUint16(data[curPos:curPos+2]))
				curPos += 2
			}
			subTable["endCode"] = endCode
			subTable["reservedPad"] = getUint16(data[curPos : curPos+2])
			curPos += 2

			var startCode []uint16
			for i := 0; i < segCount; i++ {
				startCode = append(startCode, getUint16(data[curPos:curPos+2]))
				curPos += 2
			}
			subTable["startCode"] = startCode

			var idDelta []uint16
			for i := 0; i < segCount; i++ {
				idDelta = append(idDelta, getUint16(data[curPos:curPos+2]))
				curPos += 2
			}
			subTable["idDelta"] = idDelta

			// Save the position where idRangeOffset array starts
			idRangeOffsetOffset := curPos
			subTable["idRangeOffsetOffset"] = idRangeOffsetOffset

			var idRangeOffset []uint16
			for i := 0; i < segCount; i++ {
				idRangeOffset = append(idRangeOffset, getUint16(data[curPos:curPos+2]))
				curPos += 2
			}
			subTable["idRangeOffset"] = idRangeOffset

			subTable["glyphIndexArrayOffset"] = curPos

			// The remaining is glyphIndexArray length
			glyphLen := (subTable["length"].(int) - (curPos - startOffset)) / 2
			var glyphIndexArray []uint16
			for i := 0; i < glyphLen; i++ {
				glyphIndexArray = append(glyphIndexArray, getUint16(data[curPos:curPos+2]))
				curPos += 2
			}
			subTable["glyphIndexArray"] = glyphIndexArray

		} else if format == 6 {
			subTable["length"] = int(getUint16(data[curPos : curPos+2]))
			curPos += 2
			subTable["language"] = getUint16(data[curPos : curPos+2])
			curPos += 2
			subTable["firstCode"] = getUint16(data[curPos : curPos+2])
			curPos += 2
			subTable["entryCount"] = getUint16(data[curPos : curPos+2])
			curPos += 2

			var glyphIndexArray []uint16
			entryCount := subTable["entryCount"].(int)
			for i := 0; i < entryCount; i++ {
				glyphIndexArray = append(glyphIndexArray, getUint16(data[curPos:curPos+2]))
				curPos += 2
			}
			subTable["glyphIndexArray"] = glyphIndexArray
		} else if format == 8 {
			subTable["reserved"] = getUint16(data[curPos : curPos+2])
			curPos += 2
			subTable["length"] = int(getUint32(data[curPos : curPos+4]))
			curPos += 4
			subTable["language"] = getUint32(data[curPos : curPos+4])
			curPos += 4
			var is32 []uint8
			for i := 0; i < 65536; i++ {
				is32 = append(is32, getUint8(data[curPos:curPos+1]))
				curPos++
			}
			subTable["is32"] = is32

			// n := (subTable["length"].(int) - (pos - startPos))/12
			subTable["nGroups"] = getUint32(data[curPos : curPos+4])
			curPos += 4
			n := subTable["nGroups"].(int)
			var groups []*CmapFormat8nGroup
			for i := 0; i < n; i++ {
				groups = append(groups, &CmapFormat8nGroup{
					getUint32(data[curPos : curPos+4]),
					getUint32(data[curPos+4 : curPos+8]),
					getUint32(data[curPos+8 : curPos+12]),
				})
				curPos += 12
			}
			subTable["groups"] = groups

		} else if format == 10 {
			subTable["reserved"] = getUint16(data[curPos : curPos+2])
			curPos += 2
			subTable["length"] = int(getUint32(data[curPos : curPos+4]))
			curPos += 4
			subTable["language"] = getUint32(data[curPos : curPos+4])
			curPos += 4
			subTable["startCharCode"] = getUint32(data[curPos : curPos+4])
			curPos += 4
			subTable["numChars"] = getUint32(data[curPos : curPos+4])
			curPos += 4
			numChars := subTable["numChars"].(int)

			var glyphs []uint16
			for i := 0; i < numChars; i++ {
				glyphs = append(glyphs, getUint16(data[curPos:curPos+2]))
				curPos += 2
			}
			subTable["glyphs"] = glyphs
		} else if format == 12 || format == 13 {
			subTable["reserved"] = getUint16(data[curPos : curPos+2])
			curPos += 2
			subTable["length"] = int(getUint32(data[curPos : curPos+4]))
			curPos += 4
			subTable["language"] = getUint32(data[curPos : curPos+4])
			curPos += 4
			subTable["nGroups"] = getUint32(data[curPos : curPos+4])
			curPos += 4

			n := subTable["nGroups"].(int)
			var groups []*CmapFormat8nGroup
			for i := 0; i < n; i++ {
				groups = append(groups, &CmapFormat8nGroup{
					getUint32(data[curPos : curPos+4]),
					getUint32(data[curPos+4 : curPos+8]),
					getUint32(data[curPos+8 : curPos+12]),
				})
				curPos += 12
			}
			subTable["groups"] = groups

		} else if format == 14 {
			subTable["length"] = int(getUint32(data[curPos : curPos+4]))
			curPos += 4
			subTable["numVarSelectorRecords"] = getUint32(data[curPos : curPos+4])
			curPos += 4

			n := subTable["numVarSelectorRecords"].(int)
			var groups []interface{}
			for i := 0; i < n; i++ {
				var varSelector uint32
				varSelector, err = getUint24(data[curPos : curPos+3])
				curPos += 3
				if err != nil {
					return
				}
				defaultUVSOffset := int(getUint32(data[curPos : curPos+4]))
				curPos += 4
				nonDefaultUVSOffset := int(getUint32(data[curPos : curPos+4]))
				curPos += 4

				if defaultUVSOffset != 0 {
					numUnicodeValueRanges := int(getUint32(data[curPos+defaultUVSOffset : curPos+defaultUVSOffset+4]))

					for i := 0; i < numUnicodeValueRanges; i++ {
						var startUnicode uint32
						startUnicode, err = getUint24(data[curPos : curPos+3])
						curPos += 3
						if err != nil {
							return
						}
						start := int(startUnicode)
						additionalCount := int(getUint8(data[curPos : curPos+1]))
						curPos++
						end := start + additionalCount
						var res []interface{}
						res = append(res, 0)
						res = append(res, &CmapFormatDefaultUVS{
							start,
							end,
							varSelector,
						})
						groups = append(groups, res)
					}
				}

				if nonDefaultUVSOffset != 0 {
					numUVSMappings := int(getUint32(data[startPos+nonDefaultUVSOffset : startPos+nonDefaultUVSOffset+4]))

					for i := 0; i < numUVSMappings; i++ {
						var v uint32
						v, err = getUint24(data[curPos : curPos+3])
						curPos += 3
						if err != nil {
							return
						}
						var res []interface{}
						res = append(res, 1)
						res = append(res, &CmapFormatNonDefaultUVS{
							int(v),
							getUint16(data[curPos : curPos+4]),
							varSelector,
						})
						groups = append(groups, res)
					}
				}

			}

			subTable["groups"] = groups
		} else {
			println("Warning: format " + strconv.Itoa(format) + " not support!")
		}
		cmap.SubTables = append(cmap.SubTables, subTable)
	}

	// Read Windows support
	cmap.WindowsCode, err = readWindowsCode(cmap.SubTables, maxpNumGlyphs)

	return
}

type NameRecord struct {
	PlatformID         uint16 `json:"platformId"`
	PlatformSpecificID uint16 `json:"platformSpecificId"`
	LanguageID         uint16 `json:"languageId"`
	NameID             uint16 `json:"nameId"`
	Length             uint16 `json:"length"`
	Offset             uint16 `json:"offset"`
}

type NameTable struct {
	Format       uint16                       `json:"format`
	Count        uint16                       `json:"count"`
	StringOffset uint16                       `json:"stringOffset"`
	NameRecord   []*NameRecord                `json:"nameRecord"`
	LangTagCount uint16                       `json:"langTagCount"`
	Info         map[string]map[string]string `json:"info"`
}

// macos languages
var macLang = map[int]string{
	0:   "en",
	1:   "fr",
	2:   "de",
	3:   "it",
	4:   "nl",
	5:   "sv",
	6:   "es",
	7:   "da",
	8:   "pt",
	9:   "no",
	10:  "he",
	11:  "ja",
	12:  "ar",
	13:  "fi",
	14:  "el",
	15:  "is",
	16:  "mt",
	17:  "tr",
	18:  "hr",
	19:  "zh-Hant",
	20:  "ur",
	21:  "hi",
	22:  "th",
	23:  "ko",
	24:  "lt",
	25:  "pl",
	26:  "hu",
	27:  "es",
	28:  "lv",
	29:  "se",
	30:  "fo",
	31:  "fa",
	32:  "ru",
	33:  "zh",
	34:  "nl-BE",
	35:  "ga",
	36:  "sq",
	37:  "ro",
	38:  "cz",
	39:  "sk",
	40:  "si",
	41:  "yi",
	42:  "sr",
	43:  "mk",
	44:  "bg",
	45:  "uk",
	46:  "be",
	47:  "uz",
	48:  "kk",
	49:  "az-Cyrl",
	50:  "az-Arab",
	51:  "hy",
	52:  "ka",
	53:  "mo",
	54:  "ky",
	55:  "tg",
	56:  "tk",
	57:  "mn-CN",
	58:  "mn",
	59:  "ps",
	60:  "ks",
	61:  "ku",
	62:  "sd",
	63:  "bo",
	64:  "ne",
	65:  "sa",
	66:  "mr",
	67:  "bn",
	68:  "as",
	69:  "gu",
	70:  "pa",
	71:  "or",
	72:  "ml",
	73:  "kn",
	74:  "ta",
	75:  "te",
	76:  "si",
	77:  "my",
	78:  "km",
	79:  "lo",
	80:  "vi",
	81:  "id",
	82:  "tl",
	83:  "ms",
	84:  "ms-Arab",
	85:  "am",
	86:  "ti",
	87:  "om",
	88:  "so",
	89:  "sw",
	90:  "rw",
	91:  "rn",
	92:  "ny",
	93:  "mg",
	94:  "eo",
	128: "cy",
	129: "eu",
	130: "ca",
	131: "la",
	132: "qu",
	133: "gn",
	134: "ay",
	135: "tt",
	136: "ug",
	137: "dz",
	138: "jv",
	139: "su",
	140: "gl",
	141: "af",
	142: "br",
	143: "iu",
	144: "gd",
	145: "gv",
	146: "ga",
	147: "to",
	148: "el-polyton",
	149: "kl",
	150: "az",
	151: "nn",
}

// Windows Languages
var windowsLang = map[int]string{
	1:     "ar",             // 0x0001
	2:     "bg",             // 0x0002
	3:     "ca",             // 0x0003
	4:     "zh-Hans",        // 0x0004
	5:     "cs",             // 0x0005
	6:     "da",             // 0x0006
	7:     "de",             // 0x0007
	8:     "el",             // 0x0008
	9:     "en",             // 0x0009
	10:    "es",             // 0x000a
	11:    "fi",             // 0x000b
	12:    "fr",             // 0x000c
	13:    "he",             // 0x000d
	14:    "hu",             // 0x000e
	15:    "is",             // 0x000f
	16:    "it",             // 0x0010
	17:    "ja",             // 0x0011
	18:    "ko",             // 0x0012
	19:    "nl",             // 0x0013
	20:    "no",             // 0x0014
	21:    "pl",             // 0x0015
	22:    "pt",             // 0x0016
	23:    "rm",             // 0x0017
	24:    "ro",             // 0x0018
	25:    "ru",             // 0x0019
	26:    "hr",             // 0x001a (bs, hr, or sr - using hr from original)
	27:    "sk",             // 0x001b
	28:    "sq",             // 0x001c
	29:    "sv",             // 0x001d
	30:    "th",             // 0x001e
	31:    "tr",             // 0x001f
	32:    "ur",             // 0x0020
	33:    "id",             // 0x0021
	34:    "uk",             // 0x0022
	35:    "be",             // 0x0023
	36:    "sl",             // 0x0024
	37:    "et",             // 0x0025
	38:    "lv",             // 0x0026
	39:    "lt",             // 0x0027
	40:    "tg",             // 0x0028
	41:    "fa",             // 0x0029
	42:    "vi",             // 0x002a
	43:    "hy",             // 0x002b
	44:    "az",             // 0x002c
	45:    "eu",             // 0x002d
	46:    "hsb",            // 0x002e (dsb or hsb - using hsb from original)
	47:    "mk",             // 0x002f
	48:    "st",             // 0x0030
	49:    "ts",             // 0x0031
	50:    "tn",             // 0x0032
	51:    "ve",             // 0x0033
	52:    "xh",             // 0x0034
	53:    "zu",             // 0x0035
	54:    "af",             // 0x0036
	55:    "ka",             // 0x0037
	56:    "fo",             // 0x0038
	57:    "hi",             // 0x0039
	58:    "mt",             // 0x003a
	59:    "se",             // 0x003b
	60:    "ga",             // 0x003c
	61:    "yi",             // 0x003d
	62:    "ms",             // 0x003e
	63:    "kk",             // 0x003f
	64:    "ky",             // 0x0040
	65:    "sw",             // 0x0041
	66:    "tk",             // 0x0042
	67:    "uz",             // 0x0043
	68:    "tt",             // 0x0044
	69:    "bn",             // 0x0045
	70:    "pa",             // 0x0046
	71:    "gu",             // 0x0047
	72:    "or",             // 0x0048
	73:    "ta",             // 0x0049
	74:    "te",             // 0x004a
	75:    "kn",             // 0x004b
	76:    "ml",             // 0x004c
	77:    "as",             // 0x004d
	78:    "mr",             // 0x004e
	79:    "sa",             // 0x004f
	80:    "mn",             // 0x0050
	81:    "bo",             // 0x0051
	82:    "cy",             // 0x0052
	83:    "km",             // 0x0053
	84:    "lo",             // 0x0054
	85:    "my",             // 0x0055
	86:    "gl",             // 0x0056
	87:    "kok",            // 0x0057
	88:    "mni",            // 0x0058
	89:    "sd",             // 0x0059
	90:    "syr",            // 0x005a
	91:    "si",             // 0x005b
	92:    "chr",            // 0x005c
	93:    "iu",             // 0x005d
	94:    "am",             // 0x005e
	95:    "tzm",            // 0x005f
	96:    "ks",             // 0x0060
	97:    "ne",             // 0x0061
	98:    "fy",             // 0x0062
	99:    "ps",             // 0x0063
	100:   "fil",            // 0x0064
	101:   "dv",             // 0x0065
	102:   "bin",            // 0x0066
	103:   "ff",             // 0x0067
	104:   "ha",             // 0x0068
	105:   "ibb",            // 0x0069
	106:   "yo",             // 0x006a
	107:   "quz",            // 0x006b
	108:   "nso",            // 0x006c
	109:   "ba",             // 0x006d
	110:   "lb",             // 0x006e
	111:   "kl",             // 0x006f
	112:   "ig",             // 0x0070
	113:   "kr",             // 0x0071
	114:   "om",             // 0x0072
	115:   "ti",             // 0x0073
	116:   "gn",             // 0x0074
	117:   "haw",            // 0x0075
	118:   "la",             // 0x0076
	119:   "so",             // 0x0077
	120:   "ii",             // 0x0078
	121:   "pap",            // 0x0079
	122:   "arn",            // 0x007a
	124:   "moh",            // 0x007c
	126:   "br",             // 0x007e
	128:   "ug",             // 0x0080
	129:   "mi",             // 0x0081
	130:   "oc",             // 0x0082
	131:   "co",             // 0x0083
	132:   "gsw",            // 0x0084
	133:   "sah",            // 0x0085
	134:   "qut",            // 0x0086
	135:   "rw",             // 0x0087
	136:   "wo",             // 0x0088
	140:   "prs",            // 0x008c
	145:   "gd",             // 0x0091
	146:   "ku",             // 0x0092
	147:   "quc",            // 0x0093
	1025:  "ar-SA",          // 0x0401
	1026:  "bg-BG",          // 0x0402
	1027:  "ca-ES",          // 0x0403
	1028:  "zh-TW",          // 0x0404
	1029:  "cs-CZ",          // 0x0405
	1030:  "da-DK",          // 0x0406
	1031:  "de-DE",          // 0x0407
	1032:  "el-GR",          // 0x0408
	1033:  "en-US",          // 0x0409
	1034:  "es",             // 0x040a (es-ES_tradnl)
	1035:  "fi-FI",          // 0x040b
	1036:  "fr-FR",          // 0x040c
	1037:  "he-IL",          // 0x040d
	1038:  "hu-HU",          // 0x040e
	1039:  "is-IS",          // 0x040f
	1040:  "it-IT",          // 0x0410
	1041:  "ja-JP",          // 0x0411
	1042:  "ko-KR",          // 0x0412
	1043:  "nl-NL",          // 0x0413
	1044:  "nb-NO",          // 0x0414
	1045:  "pl-PL",          // 0x0415
	1046:  "pt-BR",          // 0x0416
	1047:  "rm-CH",          // 0x0417
	1048:  "ro-RO",          // 0x0418
	1049:  "ru-RU",          // 0x0419
	1050:  "hr-HR",          // 0x041a
	1051:  "sk-SK",          // 0x041b
	1052:  "sq-AL",          // 0x041c
	1053:  "sv-SE",          // 0x041d
	1054:  "th-TH",          // 0x041e
	1055:  "tr-TR",          // 0x041f
	1056:  "ur-PK",          // 0x0420
	1057:  "id-ID",          // 0x0421
	1058:  "uk-UA",          // 0x0422
	1059:  "be-BY",          // 0x0423
	1060:  "sl-SI",          // 0x0424
	1061:  "et-EE",          // 0x0425
	1062:  "lv-LV",          // 0x0426
	1063:  "lt-LT",          // 0x0427
	1064:  "tg-Cyrl-TJ",     // 0x0428
	1065:  "fa-IR",          // 0x0429
	1066:  "vi-VN",          // 0x042a
	1067:  "hy-AM",          // 0x042b
	1068:  "az-Latn-AZ",     // 0x042c
	1069:  "eu-ES",          // 0x042d
	1070:  "hsb-DE",         // 0x042e
	1071:  "mk-MK",          // 0x042f
	1072:  "st-ZA",          // 0x0430
	1073:  "ts-ZA",          // 0x0431
	1074:  "tn-ZA",          // 0x0432
	1075:  "ve-ZA",          // 0x0433
	1076:  "xh-ZA",          // 0x0434
	1077:  "zu-ZA",          // 0x0435
	1078:  "af-ZA",          // 0x0436
	1079:  "ka-GE",          // 0x0437
	1080:  "fo-FO",          // 0x0438
	1081:  "hi-IN",          // 0x0439
	1082:  "mt-MT",          // 0x043a
	1083:  "se-NO",          // 0x043b
	1085:  "ms-MY",          // 0x043e
	1086:  "kk-KZ",          // 0x043f
	1087:  "ky-KG",          // 0x0440
	1088:  "sw-KE",          // 0x0441
	1089:  "tk-TM",          // 0x0442
	1090:  "uz-Latn-UZ",     // 0x0443
	1091:  "tt-RU",          // 0x0444
	1092:  "bn-IN",          // 0x0445
	1093:  "pa-IN",          // 0x0446
	1094:  "gu-IN",          // 0x0447
	1095:  "or-IN",          // 0x0448
	1096:  "ta-IN",          // 0x0449
	1097:  "te-IN",          // 0x044a
	1098:  "kn-IN",          // 0x044b
	1099:  "ml-IN",          // 0x044c
	1100:  "as-IN",          // 0x044d
	1101:  "mr-IN",          // 0x044e
	1102:  "sa-IN",          // 0x044f
	1103:  "mn-MN",          // 0x0450
	1104:  "bo-CN",          // 0x0451
	1105:  "cy-GB",          // 0x0452
	1106:  "km-KH",          // 0x0453
	1107:  "lo-LA",          // 0x0454
	1108:  "my-MM",          // 0x0455
	1109:  "gl-ES",          // 0x0456
	1110:  "kok-IN",         // 0x0457
	1111:  "mni-IN",         // 0x0458
	1112:  "sd-Deva-IN",     // 0x0459
	1113:  "syr-SY",         // 0x045a
	1114:  "si-LK",          // 0x045b
	1115:  "chr-Cher-US",    // 0x045c
	1116:  "iu-Cans-CA",     // 0x045d
	1117:  "am-ET",          // 0x045e
	1118:  "tzm-Arab-MA",    // 0x045f
	1120:  "ks-Arab",        // 0x0460
	1121:  "ne-NP",          // 0x0461
	1122:  "fy-NL",          // 0x0462
	1123:  "ps-AF",          // 0x0463
	1124:  "fil-PH",         // 0x0464
	1125:  "dv-MV",          // 0x0465
	1126:  "bin-NG",         // 0x0466
	1127:  "fuv-NG",         // 0x0467
	1128:  "ha-Latn-NG",     // 0x0468
	1129:  "ibb-NG",         // 0x0469
	1130:  "yo-NG",          // 0x046a
	1131:  "quz-BO",         // 0x046b
	1132:  "nso-ZA",         // 0x046c
	1133:  "ba-RU",          // 0x046d
	1134:  "lb-LU",          // 0x046e
	1135:  "kl-GL",          // 0x046f
	1136:  "ig-NG",          // 0x0470
	1137:  "kr-NG",          // 0x0471
	1138:  "om-ET",          // 0x0472
	1139:  "ti-ET",          // 0x0473
	1140:  "gn-PY",          // 0x0474
	1141:  "haw-US",         // 0x0475
	1142:  "la-Latn",        // 0x0476
	1143:  "so-SO",          // 0x0477
	1144:  "ii-CN",          // 0x0478
	1145:  "pap-029",        // 0x0479
	1146:  "arn-CL",         // 0x047a
	1148:  "moh-CA",         // 0x047c
	1150:  "br-FR",          // 0x047e
	1152:  "ug-CN",          // 0x0480
	1153:  "mi-NZ",          // 0x0481
	1154:  "oc-FR",          // 0x0482
	1155:  "co-FR",          // 0x0483
	1156:  "gsw-FR",         // 0x0484
	1157:  "sah-RU",         // 0x0485
	1158:  "qut-GT",         // 0x0486
	1159:  "rw-RW",          // 0x0487
	1160:  "wo-SN",          // 0x0488
	1164:  "prs-AF",         // 0x048c
	1168:  "zh-yue-HK",      // 0x0490
	1170:  "ku-Arab-IQ",     // 0x0492
	1171:  "quc-CO",         // 0x0493
	2049:  "ar-IQ",          // 0x0801
	2051:  "ca-ES-valencia", // 0x0803
	2052:  "zh",             // 0x0804 (zh-CN in input, but zh in original)
	2055:  "de-CH",          // 0x0807
	2057:  "en-GB",          // 0x0809
	2058:  "es-MX",          // 0x080a
	2060:  "fr-BE",          // 0x080c
	2064:  "it-CH",          // 0x0810
	2067:  "nl-BE",          // 0x0813
	2068:  "nn-NO",          // 0x0814
	2070:  "pt-PT",          // 0x0816
	2072:  "ro-MD",          // 0x0818
	2073:  "ru-MD",          // 0x0819
	2074:  "sr-Latn",        // 0x081a (sr-Latn-CS in input)
	2077:  "sv-FI",          // 0x081d
	2080:  "ur-IN",          // 0x0820
	2092:  "az-Cyrl",        // 0x082c
	2094:  "dsb",            // 0x082e (dsb-DE in input)
	2106:  "se-SE",          // 0x083b
	2108:  "ga-IE",          // 0x083c
	2110:  "ms-BN",          // 0x083e
	2115:  "uz-Cyrl",        // 0x0843
	2117:  "bn",             // 0x0845 (bn-BD in input)
	2128:  "mn-Cyrl",        // 0x0850
	2141:  "iu-Latn",        // 0x085d
	2143:  "tzm-Latn-DZ",    // 0x085f
	2171:  "quz-EC",         // 0x086b
	3073:  "ar-EG",          // 0x0c01
	3076:  "zh-HK",          // 0x0c04
	3079:  "de-AT",          // 0x0c07
	3081:  "en-AU",          // 0x0c09
	3082:  "es",             // 0x0c0a (es-ES in input)
	3084:  "fr-CA",          // 0x0c0c
	3098:  "sr",             // 0x0c1a (sr-Cyrl-CS in input)
	3131:  "se-FI",          // 0x0c3b
	3152:  "mn-Mong",        // 0x0c50
	3163:  "quz",            // 0x0c6b (quz-PE in input)
	4097:  "ar-LY",          // 0x1001
	4100:  "zh-SG",          // 0x1004
	4103:  "de-LU",          // 0x1007
	4105:  "en-CA",          // 0x1009
	4106:  "es-GT",          // 0x100a
	4108:  "fr-CH",          // 0x100c
	4122:  "hr-BA",          // 0x101a
	4155:  "smj-NO",         // 0x103b
	4191:  "tzm-Tfng-MA",    // 0x105f
	5121:  "ar-DZ",          // 0x1401
	5124:  "zh-MO",          // 0x1404
	5127:  "de-LI",          // 0x1407
	5129:  "en-NZ",          // 0x1409
	5130:  "es-CR",          // 0x140a
	5132:  "fr-LU",          // 0x140c
	5146:  "bs",             // 0x141a (bs-Latn-BA in input)
	5179:  "smj-SE",         // 0x143b
	6145:  "ar-MA",          // 0x1801
	6153:  "en-IE",          // 0x1809
	6154:  "es-PA",          // 0x180a
	6156:  "fr-MC",          // 0x180c
	6170:  "sr-Latn-BA",     // 0x181a
	6203:  "sma-NO",         // 0x183b
	7169:  "ar-TN",          // 0x1c01
	7177:  "en-ZA",          // 0x1c09
	7178:  "es-DO",          // 0x1c0a
	7194:  "sr-Cyrl-BA",     // 0x1c1a
	7227:  "sma-SE",         // 0x1c3b
	8193:  "ar-OM",          // 0x2001
	8201:  "en-JM",          // 0x2009
	8202:  "es-VE",          // 0x200a
	8218:  "bs-Cyrl",        // 0x201a
	8251:  "sms-FI",         // 0x203b
	9217:  "ar-YE",          // 0x2401
	9225:  "en-029",         // 0x2409
	9226:  "es-CO",          // 0x240a
	9242:  "sr-Latn-RS",     // 0x241a
	9275:  "smn-FI",         // 0x243b
	10241: "ar-SY",          // 0x2801
	10249: "en-BZ",          // 0x2809
	10250: "es-PE",          // 0x280a
	10266: "sr-Cyrl-RS",     // 0x281a
	11265: "ar-JO",          // 0x2c01
	11273: "en-TT",          // 0x2c09
	11274: "es-AR",          // 0x2c0a
	11290: "sr-Latn-ME",     // 0x2c1a
	12289: "ar-LB",          // 0x3001
	12297: "en-ZW",          // 0x3009
	12298: "es-EC",          // 0x300a
	12314: "sr-Cyrl-ME",     // 0x301a
	13313: "ar-KW",          // 0x3401
	13321: "en-PH",          // 0x3409
	13322: "es-CL",          // 0x340a
	14337: "ar-AE",          // 0x3801
	14346: "es-UY",          // 0x380a
	15361: "ar-BH",          // 0x3c01
	15370: "es-PY",          // 0x3c0a
	16385: "ar-QA",          // 0x4001
	16393: "en-IN",          // 0x4009
	16394: "es-BO",          // 0x400a
	17417: "en-MY",          // 0x4409
	17418: "es-SV",          // 0x440a
	18441: "en-SG",          // 0x4809
	18442: "es-HN",          // 0x480a
	19466: "es-NI",          // 0x4c0a
	20490: "es-PR",          // 0x500a
	21514: "es-US",          // 0x540a
	25626: "bs",             // 0x641a
	26650: "bs-Latn",        // 0x681a
	27674: "sr-Cyrl",        // 0x6c1a
	28698: "sr-Latn",        // 0x701a
	28731: "smn",            // 0x703b
	29740: "az-Cyrl",        // 0x742c
	30203: "sms",            // 0x743b
	30724: "zh",             // 0x7804
	30740: "nn",             // 0x7814
	30746: "bs",             // 0x781a
	30764: "az-Latn",        // 0x782c
	30779: "sma",            // 0x783b
	30819: "uz-Cyrl",        // 0x7843
	30832: "mn-Cyrl",        // 0x7850
	30845: "iu-Cans",        // 0x785d
	30847: "tzm-Tfng",       // 0x785f
	31748: "zh-Hant",        // 0x7c04
	31764: "nb",             // 0x7c14
	31770: "sr",             // 0x7c1a
	31800: "tg-Cyrl",        // 0x7c28
	31822: "dsb",            // 0x7c2e
	31867: "smj",            // 0x7c3b
	31875: "uz-Latn",        // 0x7c43
	31878: "pa-Arab",        // 0x7c46
	31888: "mn-Mong",        // 0x7c50
	31897: "sd-Arab",        // 0x7c59
	31900: "chr-Cher",       // 0x7c5c
	31901: "iu-Latn",        // 0x7c5d
	31903: "tzm-Latn",       // 0x7c5f
	31911: "ff-Latn",        // 0x7c67
	31912: "ha-Latn",        // 0x7c68
	31954: "ku-Arab",        // 0x7c92
}

// NameIDs for the name table.
var nameTableNames = [23]string{
	"copyright",              // 0
	"fontFamily",             // 1
	"fontSubfamily",          // 2
	"uniqueID",               // 3
	"fullName",               // 4
	"version",                // 5
	"postScriptName",         // 6
	"trademark",              // 7
	"manufacturer",           // 8
	"designer",               // 9
	"description",            // 10
	"manufacturerURL",        // 11
	"designerURL",            // 12
	"license",                // 13
	"licenseURL",             // 14
	"reserved",               // 15
	"preferredFamily",        // 16
	"preferredSubfamily",     // 17
	"compatibleFullName",     // 18
	"sampleText",             // 19
	"postScriptFindFontName", // 20
	"wwsFamily",              // 21
	"wwsSubfamily",           // 22
}

var macLanguageEncodings = map[int]string{
	15:  "x-mac-icelandic", // langIcelandic
	17:  "x-mac-turkish",   // langTurkish
	18:  "x-mac-croatian",  // langCroatian
	24:  "x-mac-ce",        // langLithuanian
	25:  "x-mac-ce",        // langPolish
	26:  "x-mac-ce",        // langHungarian
	27:  "x-mac-ce",        // langEstonian
	28:  "x-mac-ce",        // langLatvian
	30:  "x-mac-icelandic", // langFaroese
	37:  "x-mac-romanian",  // langRomanian
	38:  "x-mac-ce",        // langCzech
	39:  "x-mac-ce",        // langSlovak
	40:  "x-mac-ce",        // langSlovenian
	143: "x-mac-inuit",     // langInuktitut
	146: "x-mac-gaelic",    // langIrishGaelicScript
}

var macScriptEncodings = map[int]string{
	0:  "macintosh",         // smRoman
	1:  "x-mac-japanese",    // smJapanese
	2:  "x-mac-chinesetrad", // smTradChinese
	3:  "x-mac-korean",      // smKorean
	6:  "x-mac-greek",       // smGreek
	7:  "x-mac-cyrillic",    // smCyrillic
	9:  "x-mac-devanagai",   // smDevanagari
	10: "x-mac-gurmukhi",    // smGurmukhi
	11: "x-mac-gujarati",    // smGujarati
	12: "x-mac-oriya",       // smOriya
	13: "x-mac-bengali",     // smBengali
	14: "x-mac-tamil",       // smTamil
	15: "x-mac-telugu",      // smTelugu
	16: "x-mac-kannada",     // smKannada
	17: "x-mac-malayalam",   // smMalayalam
	18: "x-mac-sinhalese",   // smSinhalese
	19: "x-mac-burmese",     // smBurmese
	20: "x-mac-khmer",       // smKhmer
	21: "x-mac-thai",        // smThai
	22: "x-mac-lao",         // smLao
	23: "x-mac-georgian",    // smGeorgian
	24: "x-mac-armenian",    // smArmenian
	25: "x-mac-chinesesimp", // smSimpChinese
	26: "x-mac-tibetan",     // smTibetan
	27: "x-mac-mongolian",   // smMongolian
	28: "x-mac-ethiopic",    // smEthiopic
	29: "x-mac-ce",          // smCentralEuroRoman
	30: "x-mac-vietnamese",  // smVietnamese
	31: "x-mac-extarabic",   // smExtArabic
}

func getLangCode(platformID int, languageID int) string {
	if platformID == 0 {
		if languageID == int(0xFFFF) {
			return "und"
		}
		// need ltag
	} else if platformID == 1 {
		return macLang[languageID]
	} else if platformID == 3 {
		return windowsLang[languageID]
	}
	return ""
}

const eumnUtf16 = "utf-16"

// platformSpecificID is same with encoding
func getPlatformSpecific(platformID int, platformSpecificID int, languageID int) string {
	if platformID == 0 {
		return eumnUtf16
	} else if platformID == 1 {
		code1, exist1 := macLanguageEncodings[languageID]
		if exist1 {
			return code1
		}
		code2, exist2 := macScriptEncodings[platformSpecificID]
		if exist2 {
			return code2
		}
	} else if platformID == 3 && (platformSpecificID == 1 || platformSpecificID == 10) {
		return eumnUtf16
	}

	return ""
}

func GetName(data []byte, pos int) (nameTable *NameTable) {
	nameTable = new(NameTable)
	nameTable.Format = getUint16(data[pos : pos+2])
	nameTable.Count = getUint16(data[pos+2 : pos+4])
	nameTable.StringOffset = getUint16(data[pos+4 : pos+6])
	pos += 6

	count := int(nameTable.Count)

	stringOffset := pos + int(nameTable.StringOffset)
	info := make(map[string]map[string]string)
	for i := 0; i < count; i++ {
		nameRecord := &NameRecord{
			getUint16(data[pos : pos+2]),
			getUint16(data[pos+2 : pos+4]),
			getUint16(data[pos+4 : pos+6]),
			getUint16(data[pos+6 : pos+8]),
			getUint16(data[pos+8 : pos+10]),
			getUint16(data[pos+10 : pos+12]),
		}
		nameTable.NameRecord = append(nameTable.NameRecord, nameRecord)
		pos += 12

		property := nameTableNames[int(nameRecord.NameID)]
		language := getLangCode(int(nameRecord.PlatformID), int(nameRecord.LanguageID))
		platformSpecifi := getPlatformSpecific(int(nameRecord.LanguageID), int(nameRecord.PlatformSpecificID), int(nameRecord.LanguageID))

		if platformSpecifi != "" && language != "" {
			var text string
			if platformSpecifi == eumnUtf16 {
				text = DecodeUTF16(data, stringOffset+int(nameRecord.Offset), int(nameRecord.Length))
			} else {
				text = DecodeMACSTRING(data, stringOffset+int(nameRecord.Offset), int(nameRecord.Length), platformSpecifi)
			}

			if text != "" {
				_, exist := info[property]
				if !exist {
					info[property] = make(map[string]string)
				}
				info[property][language] = text
			}
		}
	}
	nameTable.Info = info

	if int(nameTable.Format) == 1 {
		// Windows langTagRecord is not finish
		nameTable.LangTagCount = getUint16(data[pos : pos+2])
		pos += 2
	}
	return
}

type Hhea struct {
	Version             float64 `json:"version"`
	Ascent              int16   `json:"ascent"`
	Descent             int16   `json:"descent"`
	LineGap             int16   `json:"lineGap"`
	AdvanceWidthMax     uint16  `json:"advanceWidthMax"`
	MinLeftSideBearing  int16   `json:"minLeftSideBearing"`
	MinRightSideBearing int16   `json:"minRightSideBearing"`
	XMaxExtent          int16   `json:"xMaxExtent"`
	CaretSlopeRise      int16   `json:"caretSlopeRise"`
	CaretSlopeRun       int16   `json:"caretSlopeRun"`
	CaretOffset         int16   `json:"caretOffset"`
	Reserved1           int16   `json:"reserved1"`
	Reserved2           int16   `json:"reserved2"`
	Reserved3           int16   `json:"reserved3"`
	Reserved4           int16   `json:"reserved4"`
	MetricDataFormat    int16   `json:"metricDataFormat"`
	NumOfLongHorMetrics uint16  `json:"numOfLongHorMetrics"`
}

func GetHhea(data []byte, pos int) (hhea *Hhea) {
	hhea = &Hhea{
		getFixed(data[pos : pos+2]),
		getFWord(data[pos+2 : pos+4]),
		getFWord(data[pos+4 : pos+6]),
		getFWord(data[pos+6 : pos+8]),
		getUFWord(data[pos+8 : pos+10]),
		getFWord(data[pos+10 : pos+12]),
		getFWord(data[pos+12 : pos+14]),
		getFWord(data[pos+14 : pos+16]),
		getInt16(data[pos+16 : pos+18]),
		getInt16(data[pos+18 : pos+20]),
		getFWord(data[pos+20 : pos+22]),
		getInt16(data[pos+22 : pos+24]),
		getInt16(data[pos+24 : pos+26]),
		getInt16(data[pos+26 : pos+28]),
		getInt16(data[pos+28 : pos+30]),
		getInt16(data[pos+30 : pos+32]),
		getUint16(data[pos+32 : pos+34]),
	}
	return
}

type LongHorMetric struct {
	AdvanceWidth    uint16 `json:"advanceWidth"`
	LeftSideBearing int16  `json:"leftSideBearing"`
}

type Hmtx struct {
	HMetrics        []*LongHorMetric `json:"hMetrics"`
	LeftSideBearing []int16          `json:"leftSideBearing"`
}

func GetHmtx(data []byte, pos int, numOfLongHorMetrics int, numGlyph int) (hmtx *Hmtx) {
	hmtx = new(Hmtx)

	for i := 0; i < numOfLongHorMetrics; i++ {
		hmtx.HMetrics = append(hmtx.HMetrics, &LongHorMetric{
			getUint16(data[pos : pos+2]),
			getInt16(data[pos+2 : pos+4]),
		})
		pos += 4
	}

	for i := 0; i < (numGlyph - numOfLongHorMetrics); i++ {
		hmtx.LeftSideBearing = append(hmtx.LeftSideBearing, getInt16(data[pos:pos+2]))
		pos += 2
	}
	return
}

type nPairs struct {
	Left  uint16 `json:"left"`
	Right uint16 `json:"right"`
	Value int16  `json:"value"`
}

func getWindowsKernTable(data []byte, pos int) (subHeaders map[string]int, kernPairs []*nPairs) {
	subHeaders["version"] = int(getUint16(data[pos : pos+2]))
	subHeaders["length"] = int(getUint16(data[pos+2 : pos+4]))
	subHeaders["coverage"] = int(getUint16(data[pos+4 : pos+6]))
	// Missing format 2
	if subHeaders["version"] == 0 {
		subHeaders["nPairs"] = int(getUint16(data[pos+6 : pos+8]))
		subHeaders["searchRange"] = int(getUint16(data[pos+8 : pos+10]))
		subHeaders["entrySelector"] = int(getUint16(data[pos+10 : pos+12]))
		subHeaders["rangeShift"] = int(getUint16(data[pos+12 : pos+14]))
		pos += 14
		nP := int(subHeaders["nPairs"])
		for i := 0; i < nP; i++ {
			kernPairs = append(kernPairs, &nPairs{
				getUint16(data[pos : pos+2]),
				getUint16(data[pos+2 : pos+4]),
				getFWord(data[pos+4 : pos+6]),
			})
			pos += 6
		}
	}
	return
}

func getMacKernTable(data []byte, pos int, nTables int) (subHeaders map[string]int, kernPairs []*nPairs) {
	subHeaders["length"] = int(getUint32(data[pos : pos+4]))
	subHeaders["coverage"] = int(getUint16(data[pos+4 : pos+6]))
	tupleIndex := getUint16(data[pos+6 : pos+8])
	subHeaders["tupleIndex"] = int(tupleIndex)
	subHeaders["nPairs"] = int(getUint16(data[pos+8 : pos+10]))
	subHeaders["searchRange"] = int(getUint16(data[pos+10 : pos+12]))
	subHeaders["entrySelector"] = int(getUint16(data[pos+12 : pos+14]))
	subHeaders["rangeShift"] = int(getUint16(data[pos+14 : pos+16]))
	subHeaders["version"] = int(tupleIndex & 0x00FF)
	pos += 16
	if nTables == 1 {
		// Not support table 2 3
		if subHeaders["version"] == 0 {
			for i := 0; i < subHeaders["nPairs"]; i++ {
				kernPairs = append(kernPairs, &nPairs{
					getUint16(data[pos : pos+2]),
					getUint16(data[pos+2 : pos+4]),
					getFWord(data[pos+4 : pos+6]),
				})
				pos += 6
			}
		}
	}
	return
}

type Kern struct {
	Version      int
	NTables      int
	SubHeaders   map[string]int
	Pairs        []*nPairs
	IsMacNewKern bool `json:"isNewKern,omitempty"`
}

func GetKern(data []byte, pos int) (kern *Kern, err error) {
	kern = new(Kern)
	version := int(getUint16(data[pos : pos+2]))
	nTables := int(getUint16(data[pos+2 : pos+4]))
	kern.Version = version
	kern.NTables = nTables
	if version == 0 {
		pos += 4
		kern.SubHeaders, kern.Pairs = getWindowsKernTable(data, pos)
		return
	} else if version == 1 {
		var isNewKern bool
		pos += 4
		// If nTables is 0, use new Mac kern header.
		if nTables == 0 {
			isNewKern = true
			nTables = int(getUint32(data[pos : pos+4]))
			pos += 4
		}
		kern.IsMacNewKern = isNewKern
		kern.SubHeaders, kern.Pairs = getMacKernTable(data, pos, nTables)
		return
	}

	err = errors.New("Unsupported kern table version:" + strconv.Itoa(version))
	return
}

type FastSetKVOpt struct {
	data    []byte
	pos     int
	key     string
	valType string
	toInt   bool
}

type OS2 struct {
	Version             uint16   `json:"version"`
	XAvgCharWidth       int16    `json:"xAvgCharWidth"`
	UsWeightClass       uint16   `json:"usWeightClass"`
	UsWidthClass        uint16   `json:"usWidthClass"`
	FsType              int16    `json:"fsType"`
	YSubscriptXSize     int16    `json:"ySubscriptXSize"`
	YSubscriptYSize     int16    `json:"ySubscriptYSize"`
	YSubscriptXOffset   int16    `json:"ySubscriptXOffset"`
	YSubscriptYOffset   int16    `json:"ySubscriptYOffset"`
	YSuperscriptXSize   int16    `json:"ySuperscriptXSize"`
	YSuperscriptYSize   int16    `json:"ySuperscriptYSize"`
	YSuperscriptXOffset int16    `json:"ySuperscriptXOffset"`
	YSuperscriptYOffset int16    `json:"ySuperscriptYOffset"`
	YStrikeoutSize      int16    `json:"yStrikeoutSize"`
	YStrikeoutPosition  int16    `json:"yStrikeoutPosition"`
	SFamilyClass        int16    `json:"sFamilyClass"`
	Panose              []uint8  `json:"panose"`
	UlUnicodeRange      []uint32 `json:"ulUnicodeRange"`
	AchVendID           string   `json:"achVendId"`
	FsSelection         uint16   `json:"fsSelection"`
	FsFirstCharIndex    uint16   `json:"fsFirstCharIndex"`
	FsLastCharIndex     uint16   `json:"fsLastCharIndexm"`
	STypoAscender       int16    `json:"sTypoAscender"`
	STypoDescender      int16    `json:"sTypoDescender"`
	STypoLineGap        int16    `json:"sTypoLineGap"`
	UsWinAscent         uint16   `json:"usWinAscent"`
	UsWinDescent        uint16   `json:"usWinDescent"`
	UlCodePageRange     []uint32 `json:"ulCodePageRange,omitempty"`
	SxHeight            int16    `json:"sxHeight,omitempty"`
	SCapHeight          int16    `json:"sCapHeight,omitempty"`
	UsDefaultChar       uint16   `json:"usDefaultChar,omitempty"`
	UsBreakChar         uint16   `json:"usBreakChar,omitempty"`
	UsMaxContext        uint16   `json:"usMaxContext,omitempty"`
	UsLowerPointSize    uint16   `json:"usLowerPointSize,omitempty"`
	UsUpperPointSize    uint16   `json:"usUpperPointSize,omitempty"`
}

func GetOS2(data []byte, pos int) (os2 *OS2) {
	os2.Version = getUint16(data[pos : pos+2])
	os2.XAvgCharWidth = getInt16(data[pos+2 : pos+4])
	os2.UsWeightClass = getUint16(data[pos+4 : pos+6])
	os2.UsWidthClass = getUint16(data[pos+6 : pos+8])
	os2.FsType = getInt16(data[pos+8 : pos+10])
	os2.YSubscriptXSize = getInt16(data[pos+10 : pos+12])
	os2.YSubscriptYSize = getInt16(data[pos+12 : pos+14])
	os2.YSubscriptXOffset = getInt16(data[pos+14 : pos+16])
	os2.YSubscriptYOffset = getInt16(data[pos+16 : pos+18])
	os2.YSuperscriptXSize = getInt16(data[pos+18 : pos+20])
	os2.YSuperscriptYSize = getInt16(data[pos+20 : pos+22])
	os2.YSuperscriptXOffset = getInt16(data[pos+22 : pos+24])
	os2.YSuperscriptYOffset = getInt16(data[pos+24 : pos+26])
	os2.YStrikeoutSize = getInt16(data[pos+26 : pos+28])
	os2.YStrikeoutPosition = getInt16(data[pos+28 : pos+30])
	os2.SFamilyClass = getInt16(data[pos+30 : pos+32])
	pos += 32
	for i := 0; i < 10; i++ {
		os2.Panose = append(os2.Panose, getUint8(data[pos:pos+1]))
		pos++
	}
	for i := 0; i < 4; i++ {
		os2.UlUnicodeRange = append(os2.UlUnicodeRange, getUint32(data[pos:pos+4]))
		pos += 4
	}

	// achVendIDSlice
	var achVendIDSl []int
	for i := 0; i < 4; i++ {
		achVendIDSl = append(achVendIDSl, int(getInt8(data[pos:pos+1])))
		pos++
	}
	os2.AchVendID = FromCharCode(achVendIDSl)
	os2.FsSelection = getUint16(data[pos : pos+2])
	os2.FsFirstCharIndex = getUint16(data[pos+2 : pos+4])
	os2.FsLastCharIndex = getUint16(data[pos+4 : pos+6])
	os2.STypoAscender = getInt16(data[pos+6 : pos+8])
	os2.STypoDescender = getInt16(data[pos+8 : pos+10])
	os2.STypoLineGap = getInt16(data[pos+10 : pos+12])
	os2.UsWinAscent = getUint16(data[pos+12 : pos+14])
	os2.UsWinDescent = getUint16(data[pos+14 : pos+16])
	pos += 16

	version := int(os2.Version)
	if version >= 1 {
		os2.UlCodePageRange = append(os2.UlCodePageRange, getUint32(data[pos:pos+4]))
		pos += 4
	}

	if version >= 2 {
		os2.SxHeight = getInt16(data[pos : pos+2])
		os2.SCapHeight = getInt16(data[pos+2 : pos+4])
		os2.UsDefaultChar = getUint16(data[pos+4 : pos+6])
		os2.UsBreakChar = getUint16(data[pos+6 : pos+8])
		os2.UsMaxContext = getUint16(data[pos+8 : pos+10])
	}
	pos += 10
	if version == 5 {
		os2.UsLowerPointSize = getUint16(data[pos : pos+2])
		os2.UsUpperPointSize = getUint16(data[pos+2 : pos+4])
	}

	return
}

var standardNames = []string{
	".notdef", ".null", "nonmarkingreturn", "space", "exclam", "quotedbl", "numbersign", "dollar", "percent",
	"ampersand", "quotesingle", "parenleft", "parenright", "asterisk", "plus", "comma", "hyphen", "period", "slash",
	"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "colon", "semicolon", "less",
	"equal", "greater", "question", "at", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O",
	"P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "bracketleft", "backslash", "bracketright",
	"asciicircum", "underscore", "grave", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o",
	"p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "braceleft", "bar", "braceright", "asciitilde",
	"Adieresis", "Aring", "Ccedilla", "Eacute", "Ntilde", "Odieresis", "Udieresis", "aacute", "agrave",
	"acircumflex", "adieresis", "atilde", "aring", "ccedilla", "eacute", "egrave", "ecircumflex", "edieresis",
	"iacute", "igrave", "icircumflex", "idieresis", "ntilde", "oacute", "ograve", "ocircumflex", "odieresis",
	"otilde", "uacute", "ugrave", "ucircumflex", "udieresis", "dagger", "degree", "cent", "sterling", "section",
	"bullet", "paragraph", "germandbls", "registered", "copyright", "trademark", "acute", "dieresis", "notequal",
	"AE", "Oslash", "infinity", "plusminus", "lessequal", "greaterequal", "yen", "mu", "partialdiff", "summation",
	"product", "pi", "integral", "ordfeminine", "ordmasculine", "Omega", "ae", "oslash", "questiondown",
	"exclamdown", "logicalnot", "radical", "florin", "approxequal", "Delta", "guillemotleft", "guillemotright",
	"ellipsis", "nonbreakingspace", "Agrave", "Atilde", "Otilde", "OE", "oe", "endash", "emdash", "quotedblleft",
	"quotedblright", "quoteleft", "quoteright", "divide", "lozenge", "ydieresis", "Ydieresis", "fraction",
	"currency", "guilsinglleft", "guilsinglright", "fi", "fl", "daggerdbl", "periodcentered", "quotesinglbase",
	"quotedblbase", "perthousand", "Acircumflex", "Ecircumflex", "Aacute", "Edieresis", "Egrave", "Iacute",
	"Icircumflex", "Idieresis", "Igrave", "Oacute", "Ocircumflex", "apple", "Ograve", "Uacute", "Ucircumflex",
	"Ugrave", "dotlessi", "circumflex", "tilde", "macron", "breve", "dotaccent", "ring", "cedilla", "hungarumlaut",
	"ogonek", "caron", "Lslash", "lslash", "Scaron", "scaron", "Zcaron", "zcaron", "brokenbar", "Eth", "eth",
	"Yacute", "yacute", "Thorn", "thorn", "minus", "multiply", "onesuperior", "twosuperior", "threesuperior",
	"onehalf", "onequarter", "threequarters", "franc", "Gbreve", "gbreve", "Idotaccent", "Scedilla", "scedilla",
	"Cacute", "cacute", "Ccaron", "ccaron", "dcroat",
}

type Post struct {
	Format             float64  `json:"format"`
	ItalicAngle        float64  `json:"italicAngle"`
	UnderlinePosition  int16    `json:"underlinePosition"`
	UnderlineThickness int16    `json:"underlineThickness"`
	IsFixedPitch       uint32   `json:"isFixedPitch"`
	MinMemType42       uint32   `json:"minMemType42"`
	MaxMemType42       uint32   `json:"maxMemType42"`
	MinMemType1        uint32   `json:"minMemType1"`
	MaxMemType1        uint32   `json:"maxMemType1"`
	Names              []string `json:"names,omitempty"`
	NumberOfGlyphs     uint16   `json:"numberOfGlyphs,omitempty"`
	GlyphNameIndex     []uint16 `json:"glyphNameIndex,omitempty"`
	Offset             []int8   `json:"offset,omitempty"`
}

func GetPost(data []byte, pos int) (post *Post) {
	post.Format = getFixed(data[pos : pos+4])
	post.ItalicAngle = getFixed(data[pos+4 : pos+8])
	post.UnderlinePosition = getFWord(data[pos+8 : pos+10])
	post.UnderlineThickness = getFWord(data[pos+10 : pos+12])
	post.IsFixedPitch = getUint32(data[pos+12 : pos+16])
	post.MinMemType42 = getUint32(data[pos+16 : pos+20])
	post.MaxMemType42 = getUint32(data[pos+20 : pos+24])
	post.MinMemType1 = getUint32(data[pos+24 : pos+28])
	post.MaxMemType1 = getUint32(data[pos+28 : pos+32])
	pos += 32

	format := post.Format

	if format == 1 {
		post.Names = standardNames
	} else if format == 2 {
		post.NumberOfGlyphs = getUint16(data[pos : pos+2])
		pos += 2
		numberOfGlyphs := int(post.NumberOfGlyphs)
		for i := 0; i < numberOfGlyphs; i++ {
			post.GlyphNameIndex = append(post.GlyphNameIndex, getUint16(data[pos:pos+2]))
			pos += 2
		}

		for i := 0; i < numberOfGlyphs; i++ {
			// post.Names = append(post.Names, get)
			if int(post.GlyphNameIndex[i]) >= len(standardNames) {
				nameLen := int(getInt8(data[pos : pos+1]))
				pos++

				post.Names = append(post.Names, FromCharCodeByte(data[pos:pos+nameLen]))
				pos += nameLen
			}
		}
	} else if format == 2.5 {
		post.NumberOfGlyphs = getUint16(data[pos : pos+2])
		pos += 2
		numberOfGlyphs := int(post.NumberOfGlyphs)
		for i := 0; i < numberOfGlyphs; i++ {
			post.Offset = append(post.Offset, getInt8(data[pos:pos+1]))
			pos++
		}
	}
	return
}

type SfntVariationAxis struct {
	AxisTag      uint32 `json:"axisTag"`
	MinValue     uint32 `json:"minValue"`
	DefaultValue uint32 `json:"defaultValue"`
	MaxValue     uint32 `json:"maxValue"`
	Flags        uint16 `json:"flags"`
	NameID       uint16 `json:"nameId"`
}
type SfntInstance struct {
	NameID   uint16 `json:"nameId"`
	Flags    uint16 `json:"flags"`
	Coord    uint32 `json:"coord"`
	PsNameID uint16 `json:"psNameId"`
}
type Fvar struct {
	Version        string               `json:"version"`
	OffsetToData   uint16               `json:"offsetToData"`
	CountSizePairs uint16               `json:"countSizePairs"`
	AxisCount      uint16               `json:"axisCount"`
	AxisSize       uint16               `json:"axisSize"`
	InstanceCount  uint16               `json:"instanceCount"`
	InstanceSize   uint16               `json:"instanceSize"`
	Axis           []*SfntVariationAxis `json:"axis"`
	Instance       []*SfntInstance      `json:"instance"`
}

func GetFvar(data []byte, pos int) (fvar *Fvar) {
	majorVersion := strconv.Itoa(int(getUint16(data[pos : pos+2])))
	minorVersion := strconv.Itoa(int(getUint16(data[pos+2 : pos+4])))
	fvar.Version = (majorVersion + minorVersion)
	pos += 4
	fvar.OffsetToData = getUint16(data[pos : pos+2])
	fvar.CountSizePairs = getUint16(data[pos+2 : pos+4])
	fvar.AxisCount = getUint16(data[pos+6 : pos+8])
	fvar.AxisSize = getUint16(data[pos+10 : pos+12])
	fvar.InstanceCount = getUint16(data[pos+12 : pos+14])
	fvar.InstanceSize = getUint16(data[pos+14 : pos+16])
	pos += 16

	axisCount := int(fvar.AxisCount)
	for i := 0; i < axisCount; i++ {
		fvar.Axis = append(fvar.Axis, &SfntVariationAxis{
			getUint32(data[pos : pos+4]),
			getFixed32(data[pos+4 : pos+8]),
			getFixed32(data[pos+8 : pos+12]),
			getFixed32(data[pos+12 : pos+16]),
			getUint16(data[pos+16 : pos+18]),
			getUint16(data[pos+18 : pos+20]),
		})
		pos += 20
	}

	instanceCount := int(fvar.InstanceCount)
	for i := 0; i < instanceCount; i++ {
		fvar.Instance = append(fvar.Instance, &SfntInstance{
			getUint16(data[pos : pos+2]),
			getUint16(data[pos+2 : pos+4]),
			getFixed32(data[pos+4 : pos+8]),
			getUint16(data[pos+8 : pos+10]),
		})
		pos += 10
	}

	return
}

type Itag struct {
	Version  uint32   `json:"version"`
	NumTags  uint32   `json:"numTags"`
	TagRange []string `json:"tagRange"`
}

func GetItag(data []byte, pos int) (itag *Itag, err error) {
	start := pos
	itag.Version = getUint32(data[pos : pos+4])

	if int(itag.Version) != 1 {
		err = errors.New("Unsupported ltag table version.")
		return
	}

	// skip flags
	itag.NumTags = getUint32(data[pos+8 : pos+12])
	pos += 12
	num := int(itag.NumTags)

	for i := 0; i < num; i++ {
		offset := start + int(getUint16(data[pos:pos+2]))
		len := int(getUint16(data[pos+2 : pos+4]))
		pos += 4
		tag := FromCharCodeByte(data[offset : offset+len])
		itag.TagRange = append(itag.TagRange, tag)
	}
	return
}

type Meta struct {
	Version     uint32            `json:"version"`
	Flags       uint32            `json:"flags"`
	DataOffset  uint32            `json:"dataOffset"`
	NumDataMaps uint32            `json:"numDataMaps"`
	Tags        map[string]string `json:"tags"`
}

func GetMeta(data []byte, pos int) (meta *Meta, err error) {
	start := pos
	meta.Version = getUint32(data[pos : pos+4])
	pos += 4
	if int(meta.Version) != 1 {
		err = errors.New("Unsupported META table version.")
		return
	}
	meta.Flags = getUint32(data[pos : pos+4])
	meta.DataOffset = getUint32(data[pos+4 : pos+8])
	meta.NumDataMaps = getUint32(data[pos+8 : pos+12])
	pos += 12

	num := int(meta.NumDataMaps)
	var tags map[string]string
	for i := 0; i < num; i++ {
		tag := FromCharCodeByte(data[pos : pos+4])
		offset := getUint32(data[pos+4 : pos+8])
		len := getUint32(data[pos+8 : pos+12])
		pos += 12
		textS := start + int(offset)
		text := FromCharCodeByte(data[textS : textS+int(len)])
		tags[tag] = text
	}
	meta.Tags = tags
	return
}
