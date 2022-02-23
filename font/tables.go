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

func Glyphs(data []byte) {

}
