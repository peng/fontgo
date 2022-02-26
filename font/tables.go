package font

// 首字符小写无法JSON解析
type Head struct {
	Version            string
	FontRevision       float64
	CheckSumAdjustment uint32
	MagicNumber        uint32
	Flags              uint16
	UnitsPerEm         uint16
	Created            string
	Modified           string
	XMin               int32
	YMin               int32
	XMax               int32
	YMax               int32
	MacStyle           uint16
	LowestRecPPEM      uint16
	FontDirectionHint  int16
	IndexToLocFormat   int16
	GlyphDataFormat    int16
}

func GetHead(data []byte) *Head {
	return &Head{
		getVersion(data[:4]),
		getFixed(data[4:8]),
		getUint32(data[8:12]),
		getUint32(data[12:16]),
		getUint16(data[16:18]),
		getUint16(data[18:20]),
		getLongDateTime(data[20:28]),
		getLongDateTime(data[28:36]),
		getInt32(data[36:40]),
		getInt32(data[40:44]),
		getInt32(data[44:48]),
		getInt32(data[52:56]),
		getUint16(data[56:58]),
		getUint16(data[58:60]),
		getInt16(data[60:62]),
		getInt16(data[62:64]),
		getInt16(data[64:66]),
	}
}

type GlyphCommon struct {
	numberOfContours int16
	xMin             int32
	yMin             int32
	xMax             int32
	yMax             int32
}

type GlyphSimple struct {
	GlyphCommon
	endPtsOfContours  []uint16
	instructionLength uint16
	instructions      []uint8
	flags             []uint8
	xCoordinates      []int
	yCoordinates      []int
}

type GlyphCompound struct {
	GlyphCommon
	flags      uint16
	glyphIndex uint16
	argument1  int
	argument2  int
}

// func GetGlyphSimple(data []byte) (glyph *GlyphSimple) {
// 	return &GlyphSimple{
// 		getInt16(data[0:2]),
// 		getFword(data[2:6]),
// 		getFword(data[6:10]),
// 		getFword(data[10:14]),
// 		getFword(data[14:18]),
// 	}
// 	pos := 0
// 	numberOfContours := getInt16(data[pos:2])
// 	xMin := getFword(data[2:6])
// 	yMin := getFword(data[6:10])
// 	xMax := getFword(data[10:14])
// 	yMax := getFword(data[14:18])
// }

// func GetGlyphs(data []byte) {
// 	sinpLen, compoundLen := 0,10
// 	pos := 0

// 	numberOfContours := getInt16(data[pos:pos+2])

// 	if numberOfContours >= 0 {
// 		// simple
// 	} else {
// 		// compound
// 	}
// }

type Maxp struct {
	Version string
	NumGlyphs uint16
	MaxPoints uint16
	MaxContours uint16
	MaxComponentPoints uint16
	MaxComponentContours uint16
	MaxZones uint16
	MaxTwilightPoints uint16
	MaxStorage uint16
	MaxFunctionDefs uint16
	MaxInstructionDefs uint16
	MaxStackElements uint16
	MaxSizeOfInstructions uint16
	MaxComponentElements uint16
	MaxComponentDepth uint16
}

func GetMaxp (data []byte) *Maxp {
	maxp := new(Maxp)
	maxp.Version = getVersion(data[0:4])
	maxp.NumGlyphs = getUint16(data[4:6])
	if (maxp.Version == "1.0") {
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

func GetLoca (data []byte, numGlyphs uint16, indexToLocFormat int16) []uint16 {
	// long version:  otf, ttf is different
	var locations []uint16
	pos := 0
	for i := 0; i < int(numGlyphs) + 1; i++ {
		offset := getUint16(data[pos:pos+2])
		if (indexToLocFormat == 0) {
			// 0 is short, 1 is long
			offset *= 2
		}
		locations = append(locations, offset)
		pos += 2
	}

	return locations
}