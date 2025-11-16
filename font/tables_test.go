package font

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"
)

// func diffVal(key string, sourceData interface{}, standData interface{}, t *testing.T) (err error) {
// 	if !reflect.DeepEqual(*sourceData, *standData) {
// 		fmt.Println("standData offsetTable", standData)
// 		fmt.Println("source offsetTable", sourceData)
// 		t.Log("offsetTable error")
// 		t.Fail()
// 		err = errors.New("not pass")
// 		return
// 	}
// 	return
// }

func TestAllTable(t *testing.T) {
	var (
		fileByte []byte
		err      error
	)

	fileByte, err = os.ReadFile("../test/HanyiSentyCrayon.ttf")

	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	// test offsetTable
	offsetTable := GetOffsetTable(fileByte)
	sOffsetTable := &OffsetTable{
		"TrueType",
		20,
		256,
		4,
		64,
	}
	if !reflect.DeepEqual(*offsetTable, *sOffsetTable) {
		fmt.Println("standData offsetTable", sOffsetTable)
		fmt.Println("offsetTable", offsetTable)
		t.Log("offsetTable error")
		t.Fail()
		return
	}

	var sTableContent TableContent
	sTableContent = TableContent{
		"COLR": &TagItem{
			3186909341,
			332,
			109354,
		},
		"CPAL": &TagItem{
			2465690482,
			109688,
			90,
		},
		"DSIG": &TagItem{
			1,
			16938676,
			8,
		},
		"GDEF": &TagItem{
			1059537,
			109780,
			22,
		},
		"GPOS": &TagItem{
			2215742747,
			109804,
			112,
		},
		"GSUB": &TagItem{
			3410809594,
			109916,
			298,
		},
		"OS/2": &TagItem{
			1137592327,
			110216,
			96,
		},
		"cmap": &TagItem{
			3086641967,
			110312,
			54588,
		},
		"cvt ": &TagItem{
			372968533,
			16934848,
			52,
		},
		"fpgm": &TagItem{
			2654343626,
			16934900,
			3605,
		},
		"gasp": &TagItem{
			16,
			16934840,
			8,
		},
		"glyf": &TagItem{
			3580581701,
			164900,
			16572502,
		},
		"head": &TagItem{
			473240010,
			16737404,
			54,
		},
		"hhea": &TagItem{
			320219366,
			16737460,
			36,
		},
		"hmtx": &TagItem{
			3408615136,
			16737496,
			43836,
		},
		"loca": &TagItem{
			1056580002,
			16781332,
			43848,
		},
		"maxp": &TagItem{
			802902503,
			16825180,
			32,
		},
		"name": &TagItem{
			298926501,
			16825212,
			1020,
		},
		"post": &TagItem{
			745911922,
			16826232,
			108608,
		},
		"prep": &TagItem{
			1749469340,
			16938508,
			167,
		},
	}

	// read table content
	numTables := int(offsetTable.NumTables)
	tableContent := GetTableContent(numTables, fileByte)

	if !reflect.DeepEqual(sTableContent, tableContent) {
		fmt.Println("standData tableConent", sOffsetTable)
		fmt.Println("tableConent", offsetTable)
		t.Log("tableConent error")
		t.Fail()
		return
	}
	// head table expected values:
	// checkSumAdjustment: 3928619034
	// created: 1553718651
	// flags: 13
	// fontDirectionHint: 1
	// fontRevision: 1
	// glyphDataFormat: 0
	// indexToLocFormat: 1
	// lowestRecPPEM: 6
	// macStyle: 0
	// magicNumber: 1594834165
	// modified: 1554811201
	// unitsPerEm: 2048
	// version: 1
	// xMax: 3128
	// xMin: -183
	// yMax: 1930
	// yMin: -632

	// check head table
	headInfo := tableContent["head"]
	head := GetHead(fileByte, int(headInfo.Offset))
	sHead := &Head{
		Version:            1,
		FontRevision:       1,
		CheckSumAdjustment: 3928619034,
		MagicNumber:        1594834165,
		Flags:              13,
		UnitsPerEm:         2048,
		Created:            1553718651,
		Modified:           1554811201,
		XMin:               -183,
		YMin:               -632,
		XMax:               3128,
		YMax:               1930,
		MacStyle:           0,
		LowestRecPPEM:      6,
		FontDirectionHint:  1,
		IndexToLocFormat:   1,
		GlyphDataFormat:    0,
	}

	if !reflect.DeepEqual(sHead, head) {
		fmt.Println("standData head", sHead)
		fmt.Println("head", head)
		t.Log("head error")
		t.Fail()
		return
	}

	// check maxp table
	maxpInfo := tableContent["maxp"]
	maxp := GetMaxp(fileByte, int(maxpInfo.Offset))
	sMaxp := &Maxp{
		"1.0",
		10961,
		977,
		55,
		0,
		0,
		2,
		152,
		252,
		141,
		0,
		941,
		19736,
		0,
		0,
	}

	if !reflect.DeepEqual(sMaxp, maxp) {
		fmt.Println("standData maxp", sMaxp)
		fmt.Println("maxp", maxp)
		t.Log("maxp error")
		t.Fail()
		return
	}

	// check local table
	locaInfo := tableContent["loca"]
	loca := GetLoca(fileByte, int(locaInfo.Offset), maxp.NumGlyphs, head.IndexToLocFormat)

	sLoca := map[int]int{
		0:     0,
		1:     0,
		2:     0,
		20:    34836,
		30:    60384,
		500:   1106272,
		1000:  2068516,
		10960: 16572436, // lastest
	}

	for key, val := range sLoca {
		if loca[key] != val {
			fmt.Println("standData local "+strconv.Itoa(key), val)
			fmt.Println("local "+strconv.Itoa(key), loca[:21])
			t.Log("local error")
			t.Fail()
			return
		}
	}

	// check glyph table
	glyphInfo := tableContent["glyf"]
	glyphs := GetGlyphs(fileByte, int(glyphInfo.Offset), loca, int(maxp.NumGlyphs))
	if len(glyphs.Compounds) != 0 {
		t.Log("glyphs Compounds len error", len(glyphs.Compounds))
		t.Fail()
		return
	}

	if len(glyphs.Simples) != 10961 {
		t.Log("glyphs Simples len error")
		t.Fail()
		return
	}

	type SGlyphSimple struct {
		GlyphCommon
		EndPtsOfContours []uint16 `json:"endPtsOfContours"`
	}

	type SPoint struct {
		X int
		Y int
	}

	simp1 := glyphs.Simples[0]
	sSimp1 := SGlyphSimple{
		GlyphCommon: GlyphCommon{
			NumberOfContours: 4,
			XMin:             279,
			YMin:             -56,
			XMax:             1221,
			YMax:             1354,
			Type:             "simple",
		},
		EndPtsOfContours: []uint16{58, 169, 173, 193},
	}

	if !reflect.DeepEqual(simp1.GlyphCommon, sSimp1.GlyphCommon) {
		fmt.Println("standData glyph simple GlyphCommon", sSimp1.GlyphCommon)
		fmt.Println("glyph simple GlyphCommon", simp1.GlyphCommon)
		t.Log("glyph simple GlyphCommon error")
		t.Fail()
		return
	}

	if !reflect.DeepEqual(simp1.EndPtsOfContours, sSimp1.EndPtsOfContours) {
		fmt.Println("standData glyph simple EndPtsOfContours", sSimp1.EndPtsOfContours)
		fmt.Println("glyph simple EndPtsOfContours", simp1.EndPtsOfContours)
		t.Log("glyph simple EndPtsOfContours error")
		t.Fail()
		return
	}

	// check point
	if len(simp1.Points) != 194 {
		fmt.Println("glyph simple point length", len(simp1.Points))
		t.Log("glyph simple points length error")
		t.Fail()
		return
	}

	sPoints1 := map[int]Point{
		0: {
			X: 597,
			Y: 1354,
			Flag: &Flag{
				OnCurve:      true,
				Repeat:       false,
				XSame:        false,
				XShortVector: false,
				YSame:        false,
				YShortVector: false,
			},
		},
		1: {
			X: 66,
			Y: -18,
			Flag: &Flag{
				OnCurve:      false,
				Repeat:       false,
				XSame:        true,
				XShortVector: true,
				YSame:        false,
				YShortVector: true,
			},
		},
		2: {
			X: 0,
			Y: -18,
			Flag: &Flag{
				OnCurve:      true,
				Repeat:       false,
				XSame:        true,
				XShortVector: false,
				YSame:        false,
				YShortVector: true,
			},
		},
		3: {
			X: -6,
			Y: -12,
			Flag: &Flag{
				OnCurve:      false,
				Repeat:       false,
				XSame:        false,
				XShortVector: true,
				YSame:        false,
				YShortVector: true,
			},
		},
		4: {
			X: 0,
			Y: -6,
			Flag: &Flag{
				OnCurve:      true,
				Repeat:       false,
				XSame:        true,
				XShortVector: false,
				YSame:        false,
				YShortVector: true,
			},
		},
		5: {
			X: 12,
			Y: -6,
			Flag: &Flag{
				OnCurve:      true,
				Repeat:       false,
				XSame:        true,
				XShortVector: true,
				YSame:        false,
				YShortVector: true,
			},
		},
		6: {
			X: 2,
			Y: 12,
			Flag: &Flag{
				OnCurve:      false,
				Repeat:       false,
				XSame:        true,
				XShortVector: true,
				YSame:        true,
				YShortVector: true,
			},
		},
		7: {
			X: 10,
			Y: 0,
			Flag: &Flag{
				OnCurve:      true,
				Repeat:       true,
				XSame:        true,
				XShortVector: true,
				YSame:        true,
				YShortVector: false,
			},
		},
		8: {
			X: 6,
			Y: 0,
			Flag: &Flag{
				OnCurve:      true,
				Repeat:       true,
				XSame:        true,
				XShortVector: true,
				YSame:        true,
				YShortVector: false,
			},
		},
		9: {
			X: -6,
			Y: -42,
			Flag: &Flag{
				OnCurve:      true,
				Repeat:       false,
				XSame:        false,
				XShortVector: true,
				YSame:        false,
				YShortVector: true,
			},
		},
		10: {
			X: 60,
			Y: -24,
			Flag: &Flag{
				OnCurve:      false,
				Repeat:       false,
				XSame:        true,
				XShortVector: true,
				YSame:        false,
				YShortVector: true,
			},
		},
		120: {
			X: -12,
			Y: 0,
			Flag: &Flag{
				OnCurve:      true,
				Repeat:       false,
				XSame:        false,
				XShortVector: true,
				YSame:        true,
				YShortVector: false,
			},
		},
		144: {
			X: 0,
			Y: 18,
			Flag: &Flag{
				OnCurve:      true,
				Repeat:       false,
				XSame:        true,
				XShortVector: false,
				YSame:        true,
				YShortVector: true,
			},
		},
	}

	for key, val := range sPoints1 {
		p := simp1.Points[key]
		if !reflect.DeepEqual(p.Flag, val.Flag) {
			fmt.Println("glyph simple point index", key)
			b, _ := json.MarshalIndent(val.Flag, "", "  ")
			fmt.Println("glyph simple point flag source", string(b))
			pj, _ := json.MarshalIndent(p.Flag, "", "  ")
			fmt.Println("glyph simple point flag", p, string(pj))
			t.Log("glyph simple points error")
			t.Fail()
			return
		}
	}

	for key, val := range sPoints1 {
		p := simp1.Points[key]
		if p.X != val.X || p.Y != val.Y {
			fmt.Println("glyph simple point index", key)
			b, _ := json.MarshalIndent(val.Flag, "", "  ")
			fmt.Println("glyph simple point source", val, string(b))
			pj, _ := json.MarshalIndent(p.Flag, "", "  ")
			fmt.Println("glyph simple point", p, string(pj))
			t.Log("glyph simple points error")
			t.Fail()
			return
		}
	}
	// test GlyphCompound
	fileByte2, err2 := os.ReadFile("../test/Changa-Regular.ttf")

	if err2 != nil {
		t.Log(err2)
		t.Fail()
		return
	}

	offsetTable2 := GetOffsetTable(fileByte2)
	numTables2 := int(offsetTable2.NumTables)
	tableContent2 := GetTableContent(numTables2, fileByte2)

	maxpInfo2 := tableContent2["maxp"]
	maxp2 := GetMaxp(fileByte2, int(maxpInfo2.Offset))

	headInfo2 := tableContent2["head"]
	head2 := GetHead(fileByte2, int(headInfo2.Offset))

	locaInfo2 := tableContent2["loca"]
	loca2 := GetLoca(fileByte2, int(locaInfo2.Offset), maxp2.NumGlyphs, head2.IndexToLocFormat)

	glyphInfo2 := tableContent2["glyf"]
	glyphs2 := GetGlyphs(fileByte2, int(glyphInfo2.Offset), loca2, int(maxp2.NumGlyphs))
	if len(glyphs2.Compounds) != 448 {
		t.Log("glyphs Compounds len error", len(glyphs.Compounds))
		t.Fail()
		return
	}

	sCompound := map[int]GlyphCompound{
		0: {
			GlyphCommon: GlyphCommon{
				NumberOfContours: -1,
				XMin:             15,
				YMin:             0,
				XMax:             536,
				YMax:             837,
				Type:             "compound",
			},
			Component: []Component{
				{
					// Flags:       34,
					// GlyphIndex: 789,
					Argument1: 15,
					Argument2: 0,
					Scale01:   0,
					Scale10:   0,
					Xscale:    16384,
					Yscale:    16384,
				},
				{
					// Flags:      34,
					// GlyphIndex: 4,
					Argument1: 0,
					Argument2: 0,
					Scale01:   0,
					Scale10:   0,
					Xscale:    16384,
					Yscale:    16384,
				},
				{
					// Flags:      263,
					// GlyphIndex: 717,
					Argument1: 488,
					Argument2: 125,
					Scale01:   0,
					Scale10:   0,
					Xscale:    16384,
					Yscale:    16384,
				},
			},
			InstructionLength: 83,
			Instructions: []uint8{
				179, 27, 1, 2, 72, 75, 176, 50, 80, 88, 64, 26, 6, 1, 5, 0,
				0, 1, 5, 0, 101, 0, 4, 4, 2, 93, 0, 2, 2, 37, 75, 3, 1, 1,
				1, 38, 1, 76, 27, 64, 24, 0, 2, 0, 4, 5, 2, 4, 101, 6, 1, 5,
				0, 0, 1, 5, 0, 101, 3, 1, 1, 1, 41, 1, 76, 89, 64, 14, 9, 9,
				9, 12, 9, 12, 18, 17, 17, 17, 17, 7, 8, 36, 43,
			},
		},
		5: {
			GlyphCommon: GlyphCommon{
				NumberOfContours: -1,
				XMin:             15,
				YMin:             0,
				XMax:             536,
				YMax:             792,
				Type:             "compound",
			},
			Component: []Component{
				{
					// Flags:      34,
					// GlyphIndex: 789,
					Argument1: 15,
					Argument2: 0,
					Scale01:   0,
					Scale10:   0,
					Xscale:    16384,
					Yscale:    16384,
				},
				{
					// Flags:      34,
					// GlyphIndex: 4,
					Argument1: 0,
					Argument2: 0,
					Scale01:   0,
					Scale10:   0,
					Xscale:    16384,
					Yscale:    16384,
				},
				{
					// Flags:      263,
					// GlyphIndex: 725,
					Argument1: 475,
					Argument2: 125,
					Scale01:   0,
					Scale10:   0,
					Xscale:    16384,
					Yscale:    16384,
				},
			},
			InstructionLength: 104,
			Instructions: []uint8{
				75, 176, 50, 80, 88, 64, 35, 0, 6, 9, 1, 7, 2, 6, 7, 101,
				8, 1, 5, 0, 0, 1, 5, 0, 101, 0, 4, 4, 2, 93, 0, 2,
				2, 37, 75, 3, 1, 1, 1, 38, 1, 76, 27, 64, 33, 0, 6, 9,
				1, 7, 2, 6, 7, 101, 0, 2, 0, 4, 5, 2, 4, 101, 8, 1,
				5, 0, 0, 1, 5, 0, 101, 3, 1, 1, 1, 41, 1, 76, 89, 64,
				22, 13, 13, 9, 9, 13, 22, 13, 21, 18, 16, 9, 12, 9, 12, 18,
				17, 17, 17, 17, 10, 8, 36, 43,
			},
		},
		45: {
			GlyphCommon: GlyphCommon{
				NumberOfContours: -1,
				XMin:             95,
				YMin:             -286,
				XMax:             566,
				YMax:             625,
				Type:             "compound",
			},
			Component: []Component{
				{
					// Flags:      34,
					// GlyphIndex: 789,
					Argument1: 95,
					Argument2: 0,
					Scale01:   0,
					Scale10:   0,
					Xscale:    16384,
					Yscale:    16384,
				},
				{
					// Flags:      34,
					// GlyphIndex: 69,
					Argument1: 0,
					Argument2: 0,
					Scale01:   0,
					Scale10:   0,
					Xscale:    16384,
					Yscale:    16384,
				},
				{
					// Flags:      259,
					// GlyphIndex: 727,
					Argument1: 424,
					Argument2: 0,
					Scale01:   0,
					Scale10:   0,
					Xscale:    16384,
					Yscale:    16384,
				},
			},
			InstructionLength: 63,
			Instructions: []uint8{
				180, 33, 20, 2, 4, 71, 75, 176, 50, 80, 88, 64, 18, 0, 4, 2,
				4, 132, 1, 1, 0, 0, 37, 75, 3, 1, 2, 2, 38, 2, 76, 27,
				64, 18, 0, 4, 2, 4, 132, 1, 1, 0, 0, 2, 93, 3, 1, 2,
				2, 41, 2, 76, 89, 183, 39, 22, 17, 22, 17, 5, 8, 36, 43,
			},
		},
		99: {
			GlyphCommon: GlyphCommon{
				NumberOfContours: -1,
				XMin:             60,
				YMin:             -15,
				XMax:             655,
				YMax:             700,
				Type:             "compound",
			},
			Component: []Component{
				{
					// Flags:      34,
					// GlyphIndex: 789,
					Argument1: 60,
					Argument2: 0,
					Scale01:   0,
					Scale10:   0,
					Xscale:    16384,
					Yscale:    16384,
				},
				{
					// Flags:      34,
					// GlyphIndex: 152,
					Argument1: 0,
					Argument2: 0,
					Scale01:   0,
					Scale10:   0,
					Xscale:    16384,
					Yscale:    16384,
				},
				{
					// Flags:      259,
					// GlyphIndex: 719,
					Argument1: 650,
					Argument2: 0,
					Scale01:   0,
					Scale10:   0,
					Xscale:    16384,
					Yscale:    16384,
				},
			},
			InstructionLength: 180,
			Instructions: []uint8{
				64, 17, 31, 19, 1, 3, 2, 5, 44, 18, 2, 3, 2, 4, 1, 0,
				4, 3, 74, 75, 176, 17, 80, 88, 64, 27, 0, 5, 5, 47, 75, 0,
				3, 3, 2, 95, 0, 2, 2, 48, 75, 0, 4, 4, 0, 95, 1, 1,
				0, 0, 38, 0, 76, 27, 75, 176, 25, 80, 88, 64, 31, 0, 5, 5,
				47, 75, 0, 3, 3, 2, 95, 0, 2, 2, 48, 75, 0, 0, 0, 38,
				75, 0, 4, 4, 1, 95, 0, 1, 1, 46, 1, 76, 27, 75, 176, 50,
				80, 88, 64, 31, 0, 5, 2, 5, 131, 0, 3, 3, 2, 95, 0, 2,
				2, 48, 75, 0, 0, 0, 38, 75, 0, 4, 4, 1, 95, 0, 1, 1,
				46, 1, 76, 27, 64, 31, 0, 5, 2, 5, 131, 0, 3, 3, 2, 95,
				0, 2, 2, 48, 75, 0, 0, 0, 41, 75, 0, 4, 4, 1, 95, 0,
				1, 1, 49, 1, 76, 89, 89, 89, 64, 9, 41, 35, 36, 38, 35, 18,
				6, 8, 37, 43,
			},
		},
	}

	for ind, compound := range sCompound {
		curComp := glyphs2.Compounds[ind]

		// check InstructionLength
		if curComp.InstructionLength != compound.InstructionLength {
			t.Log("glyphs compound InstructionLength error index:", ind)
			t.Log("glyphs compound InstructionLength is: ", curComp.InstructionLength)
			t.Fail()
			return
		}

		// check Instructions
		for instructionInd := 0; instructionInd < len(compound.Instructions); instructionInd++ {
			if curComp.Instructions[instructionInd] != compound.Instructions[instructionInd] {
				t.Log("glyphs compound error index:", ind)
				t.Log("glyphs compound Instructions error index:", instructionInd)
				t.Log("glyphs compound Instructions value is: ", curComp.Instructions[instructionInd], compound.Instructions[instructionInd])
				t.Fail()
				return
			}
		}

		// check GlyphCommon
		if !reflect.DeepEqual(curComp.GlyphCommon, compound.GlyphCommon) {
			t.Log("glyphs compound error index:", ind)
			t.Log("glyphs compound GlyphCommon is: ", curComp.GlyphCommon)
		}

		// check Component
		for compInd := 0; compInd < len(curComp.Component); compInd++ {
			if !reflect.DeepEqual(curComp.Component[compInd], compound.Component[compInd]) {
				t.Log("glyphs compound error index:", ind)
				t.Log("glyphs compound Component error index:", compInd)
				t.Log("glyphs compound Component value is:", curComp.Component[compInd])
			}
		}
	}

	// check cmap
	cmapInfo := tableContent["cmap"]
	cmap, cmapErr := GetCmap(fileByte, int(cmapInfo.Offset), int(maxp.NumGlyphs))
	if cmapErr != nil {
		t.Log(cmapErr)
		t.Fail()
		return
	}

	if len(cmap.SubTables) != 3 {
		t.Log("cmap subTables length error, error length is ", len(cmap.SubTables))
		t.Fail()
	}

	if cmap.NumberSubtables != 3 {
		t.Log("cmap NumberSubtables error, NumberSubtables is ", cmap.NumberSubtables)
		t.Fail()
	}

	if cmap.Version != 0 {
		t.Log("cmap Version error, Version is ", cmap.NumberSubtables)
		t.Fail()
	}

	// check cmap format 4
	cmapSubTable0 := cmap.SubTables[0]
	if cmapSubTable0["platformID"].(int) != 0 {
		t.Log("cmap subTable 0 platformID error, platformID is ", cmapSubTable0["platformID"].(int))
		t.Fail()
	}

	if cmapSubTable0["platformSpecificID"].(int) != 3 {
		t.Log("cmap subTable 0 platformSpecificID error, platformSpecificID is ", cmapSubTable0["platformSpecificID"].(int))
		t.Fail()
	}

	if cmapSubTable0["offset"].(uint32) != 28 {
		t.Log("cmap subTable 0 offset error, offset is ", cmapSubTable0["offset"].(int))
		t.Fail()
	}

	if cmapSubTable0["format"].(int) != 4 {
		t.Log("cmap subTablfe 0 format error, format is ", cmapSubTable0["format"].(int))
		t.Fail()
	}

	if cmapSubTable0["length"].(int) != 54298 {
		t.Log("cmap subTablfe 0 length error, length is ", cmapSubTable0["length"].(int))
		t.Fail()
	}

	if cmapSubTable0["language"].(uint16) != 0 {
		t.Log("cmap subTablfe 0 language error, language is ", cmapSubTable0["language"].(uint16))
		t.Fail()
	}

	if cmapSubTable0["segCountX2"].(uint16) != 9486 {
		t.Log("cmap subTablfe 0 segCountX2 error, segCountX2 is ", cmapSubTable0["segCountX2"].(uint16))
		t.Fail()
	}

	if cmapSubTable0["searchRange"].(uint16) != 8192 {
		t.Log("cmap subTablfe 0 searchRange error, searchRange is ", cmapSubTable0["searchRange"].(uint16))
		t.Fail()
	}

	if cmapSubTable0["entrySelector"].(uint16) != 12 {
		t.Log("cmap subTablfe 0 entrySelector error, entrySelector is ", cmapSubTable0["entrySelector"].(uint16))
		t.Fail()
	}

	if cmapSubTable0["rangeShift"].(uint16) != 1294 {
		t.Log("cmap subTablfe 0 rangeShift error, rangeShift is ", cmapSubTable0["rangeShift"].(uint16))
		t.Fail()
	}
	// check cmap format 4 endCode slice
	cmapSubTable0endCode := map[int]uint16{
		0:    0,
		1:    29,
		99:   20186,
		2321: 29822,
		4741: 65509,
		4742: 65535,
	}
	for ind, val := range cmapSubTable0endCode {
		endCode := cmapSubTable0["endCode"].([]uint16)
		if endCode[ind] != val {
			t.Log("cmap subTablfe 0 endCode error, endCode index and value is ", ind, endCode[ind], val)
			t.Fail()
		}
	}

	if cmapSubTable0["reservedPad"].(uint16) != 0 {
		t.Log("cmap subTablfe 0 reservedPad error, reservedPad is ", cmapSubTable0["reservedPad"].(uint16))
		t.Fail()
	}

	// check format 4 startCode slice
	cmapSubtable0StartCode := map[int]uint16{
		0:    0,
		44:   8804,
		1632: 26911,
		1633: 26916,
		3700: 35712,
		3701: 35722,
		4741: 65509,
		4742: 65535,
	}
	for ind, val := range cmapSubtable0StartCode {
		startCode := cmapSubTable0["startCode"].([]uint16)
		if startCode[ind] != val {
			t.Log("cmap subTable 0 startCode error, startCode index and value is ", ind, startCode[ind], val)
			t.Fail()
		}
	}

	// check format 4 idDelta slice
	idDelta := cmapSubTable0["idDelta"].([]uint16)
	if len(idDelta) != 4743 {
		t.Log("cmap subtable 0 idDelta length error, real and current length is ", 4743, len(idDelta))
		t.Fail()
	}

	cmapSubtable0IdDelta := map[int]uint16{
		0:    1,
		1:    65508,
		98:   0,
		99:   0,
		1971: 0,
		1972: 37274,
		4740: 10760,
		4741: 10752,
		4742: 1,
	}
	for ind, val := range cmapSubtable0IdDelta {
		if idDelta[ind] != val {
			t.Log("cmap subTablfe 0 startCode error, startCode index and value is ", ind, idDelta[ind], val)
			t.Fail()
		}
	}

	// check format 4 idRangeOffset slice
	cmapSubtable0IdRangeOffset := map[int]uint16{
		0:    0,
		1:    0,
		98:   9658,
		99:   9670,
		2378: 12804,
		2379: 12806,
		4741: 0,
		4742: 0,
	}
	for ind, val := range cmapSubtable0IdRangeOffset {
		idRangeOffset := cmapSubTable0["idRangeOffset"].([]uint16)
		if idRangeOffset[ind] != val {
			t.Log("cmap subTable 0 idRangeOffset error, idRangeOffset index and value is ", ind, idRangeOffset[ind], val)
			t.Fail()
		}
	}

	if cmapSubTable0["glyphIndexArrayOffset"] != 148300 {
		t.Log("cmap subTablfe 0 glyphIndexArrayOffset error, glyphIndexArrayOffset is ", cmapSubTable0["glyphIndexArrayOffset"].(int))
		t.Fail()
	}

	// check format 4 glyphIndexArray slice
	// check length
	if len(cmapSubTable0["glyphIndexArray"].([]uint16)) != 8169 {
		t.Log("cmap subtable 0 glyphIndexArray error, glyphIndexArray length is ", len(cmapSubTable0["glyphIndexArray"].([]uint16)))
		t.Fail()
	}

	cmapSubTable0GlyphIndexArray := map[int]uint16{
		0:    3,
		1:    10660,
		600:  2558,
		601:  9532,
		2848: 252,
		2849: 6753,
		8167: 10739,
		8168: 10724,
	}
	for ind, val := range cmapSubTable0GlyphIndexArray {
		glyphIndexArray := cmapSubTable0["glyphIndexArray"].([]uint16)
		if glyphIndexArray[ind] != val {
			t.Log("cmap subTable 0 glyphIndexArray error, glyphIndexArray index and value is ", ind, glyphIndexArray[ind], val)
			t.Fail()
		}
	}

	// check cmap format 0
	cmapSubTable1 := cmap.SubTables[1]
	if cmapSubTable1["platformID"].(int) != 1 {
		t.Log("cmap subTable 1 platformID error, platformID is ", cmapSubTable1["platformID"].(int))
		t.Fail()
	}

	if cmapSubTable1["platformSpecificID"].(int) != 0 {
		t.Log("cmap subTable 1 platformSpecificID error, platformSpecificID is ", cmapSubTable1["platformSpecificID"].(int))
		t.Fail()
	}

	if cmapSubTable1["offset"].(uint32) != 54326 {
		t.Log("cmap subTable 1 offset error, offset is ", cmapSubTable1["offset"].(uint32))
		t.Fail()
	}

	if cmapSubTable1["length"].(int) != 262 {
		t.Log("cmap subTable 1 length error, length is ", cmapSubTable1["length"].(int))
		t.Fail()
	}

	if cmapSubTable1["language"].(uint16) != 0 {
		t.Log("cmap subTable 1 language error, language is ", cmapSubTable1["language"].(uint16))
		t.Fail()
	}

	// check format 0 glyphIndexArray slice
	// check length
	if len(cmapSubTable1["glyphIndexArray"].([]uint8)) != 256 {
		t.Log("cmap subtable 1 glyphIndexArray error, glyphIndexArray length is ", len(cmapSubTable1["glyphIndexArray"].([]uint8)))
		t.Fail()
	}

	cmapSubTable1GlyphIndexArray := map[int]uint8{
		0:   1,
		1:   0,
		2:   0,
		180: 0,
		181: 0,
		234: 0,
		235: 0,
		254: 0,
		255: 0,
	}

	for ind, val := range cmapSubTable1GlyphIndexArray {
		glyphIndexArray := cmapSubTable1["glyphIndexArray"].([]uint8)
		if glyphIndexArray[ind] != val {
			t.Log("cmap subTable 1 glyphIndexArray error, glyphIndexArray index and value is ", ind, glyphIndexArray[ind], val)
			t.Fail()
		}
	}

	// check hhea table
	hheaInfo := tableContent["hhea"]
	hhea := GetHhea(fileByte, int(hheaInfo.Offset))
	// Expected hhea values captured from HanyiSentyCrayon.ttf
	if hhea.Version != 1 {
		t.Errorf("hhea Version want 1 got %v", hhea.Version)
	}
	if hhea.Ascent != 1716 {
		t.Errorf("hhea Ascent want 1716 got %d", hhea.Ascent)
	}
	if hhea.Descent != -418 {
		t.Errorf("hhea Descent want -418 got %d", hhea.Descent)
	}
	if hhea.LineGap != 222 {
		t.Errorf("hhea LineGap want 222 got %d", hhea.LineGap)
	}
	if hhea.AdvanceWidthMax != 2239 {
		t.Errorf("hhea AdvanceWidthMax want 2239 got %d", hhea.AdvanceWidthMax)
	}
	if hhea.MinLeftSideBearing != -183 {
		t.Errorf("hhea MinLeftSideBearing want -183 got %d", hhea.MinLeftSideBearing)
	}
	if hhea.MinRightSideBearing != -2309 {
		t.Errorf("hhea MinRightSideBearing want -2309 got %d", hhea.MinRightSideBearing)
	}
	if hhea.XMaxExtent != 3128 {
		t.Errorf("hhea XMaxExtent want 3128 got %d", hhea.XMaxExtent)
	}
	if hhea.CaretSlopeRise != 1 {
		t.Errorf("hhea CaretSlopeRise want 1 got %d", hhea.CaretSlopeRise)
	}
	if hhea.CaretSlopeRun != 0 {
		t.Errorf("hhea CaretSlopeRun want 0 got %d", hhea.CaretSlopeRun)
	}
	if hhea.CaretOffset != 0 {
		t.Errorf("hhea CaretOffset want 0 got %d", hhea.CaretOffset)
	}
	if hhea.MetricDataFormat != 0 {
		t.Errorf("hhea MetricDataFormat want 0 got %d", hhea.MetricDataFormat)
	}
	if hhea.NumOfLongHorMetrics != 10957 {
		t.Errorf("hhea NumOfLongHorMetrics want 10957 got %d", hhea.NumOfLongHorMetrics)
	}
	if hhea.Reserved1 != 0 || hhea.Reserved2 != 0 || hhea.Reserved3 != 0 || hhea.Reserved4 != 0 {
		t.Errorf("hhea Reserved fields expected 0 got %d %d %d %d", hhea.Reserved1, hhea.Reserved2, hhea.Reserved3, hhea.Reserved4)
	}

	// check hmtx table
	hmtxInfo := tableContent["hmtx"]
	numLHM := int(hhea.NumOfLongHorMetrics)
	numGlyphs := int(maxp.NumGlyphs)
	hmtxOffset := int(hmtxInfo.Offset)

	// Compute expected metrics directly from bytes
	type pairMetric struct {
		aw  uint16
		lsb int16
	}
	var expectedFirst []pairMetric
	for i := 0; i < 3 && i < numLHM; i++ {
		base := hmtxOffset + i*4
		aw := binary.BigEndian.Uint16(fileByte[base : base+2])
		lsb := int16(binary.BigEndian.Uint16(fileByte[base+2 : base+4]))
		expectedFirst = append(expectedFirst, pairMetric{aw: aw, lsb: lsb})
	}
	var expectedLastLSB int16
	if numGlyphs-numLHM > 0 {
		lastLSBIndex := (numGlyphs - numLHM - 1)
		pos := hmtxOffset + numLHM*4 + lastLSBIndex*2
		expectedLastLSB = int16(binary.BigEndian.Uint16(fileByte[pos : pos+2]))
	}

	hmtx := GetHmtx(fileByte, hmtxOffset, numLHM, numGlyphs)
	if hmtx == nil {
		t.Fatal("hmtx GetHmtx returned nil")
	}
	if got := len(hmtx.HMetrics); got != numLHM {
		t.Fatalf("hmtx HMetrics length want %d got %d", numLHM, got)
	}
	if got := len(hmtx.LeftSideBearing); got != (numGlyphs - numLHM) {
		t.Fatalf("hmtx LeftSideBearing length want %d got %d", (numGlyphs - numLHM), got)
	}
	for i := range expectedFirst {
		if hmtx.HMetrics[i].AdvanceWidth != expectedFirst[i].aw {
			t.Errorf("hmtx HMetrics[%d].AdvanceWidth want %d got %d", i, expectedFirst[i].aw, hmtx.HMetrics[i].AdvanceWidth)
		}
		if hmtx.HMetrics[i].LeftSideBearing != expectedFirst[i].lsb {
			t.Errorf("hmtx HMetrics[%d].LeftSideBearing want %d got %d", i, expectedFirst[i].lsb, hmtx.HMetrics[i].LeftSideBearing)
		}
	}
	if numGlyphs-numLHM > 0 {
		last := len(hmtx.LeftSideBearing) - 1
		if hmtx.LeftSideBearing[last] != expectedLastLSB {
			t.Errorf("hmtx LeftSideBearing[last] want %d got %d", expectedLastLSB, hmtx.LeftSideBearing[last])
		}
	}

	// Sample additional hmtx HMetrics entries for spot checks
	sHMetrics := map[int]LongHorMetric{
		0:     {AdvanceWidth: 1024, LeftSideBearing: 0},
		1:     {AdvanceWidth: 0, LeftSideBearing: 0},
		2:     {AdvanceWidth: 508, LeftSideBearing: 0},
		3:     {AdvanceWidth: 635, LeftSideBearing: 0},
		5600:  {AdvanceWidth: 1500, LeftSideBearing: 258},
		5601:  {AdvanceWidth: 1500, LeftSideBearing: 321},
		5602:  {AdvanceWidth: 1500, LeftSideBearing: 114},
		5603:  {AdvanceWidth: 1500, LeftSideBearing: 192},
		5698:  {AdvanceWidth: 1500, LeftSideBearing: 228},
		5699:  {AdvanceWidth: 1500, LeftSideBearing: 237},
		7942:  {AdvanceWidth: 1500, LeftSideBearing: 102},
		7943:  {AdvanceWidth: 1500, LeftSideBearing: 84},
		10956: {AdvanceWidth: 1500, LeftSideBearing: 487},
	}
	for ind, expected := range sHMetrics {
		if ind >= len(hmtx.HMetrics) {
			t.Errorf("hmtx HMetrics index %d out of range (length %d)", ind, len(hmtx.HMetrics))
			continue
		}
		actual := hmtx.HMetrics[ind]
		if actual.AdvanceWidth != expected.AdvanceWidth {
			t.Errorf("hmtx HMetrics[%d].AdvanceWidth want %d got %d", ind, expected.AdvanceWidth, actual.AdvanceWidth)
		}
		if actual.LeftSideBearing != expected.LeftSideBearing {
			t.Errorf("hmtx HMetrics[%d].LeftSideBearing want %d got %d", ind, expected.LeftSideBearing, actual.LeftSideBearing)
		}
	}
}

func TestGetFvar(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x00, 0x00, // version 1.0
		0x00, 0x00, // offsetToData
		0x00, 0x00, // countSizePairs
		0x00, 0x01, // axisCount
		0x00, 0x14, // axisSize
		0x00, 0x01, // instanceCount
		0x00, 0x0A, // instanceSize
		// axis
		0x12, 0x34, 0x56, 0x78, // tag
		0x00, 0x00, 0x00, 0x00, // min
		0x00, 0x00, 0x00, 0x01, // def
		0x00, 0x00, 0x00, 0x02, // max
		0x00, 0x03, // flags
		0x00, 0x04, // nameID
		// instance
		0x00, 0x05, // nameID
		0x00, 0x06, // flags
		0x00, 0x00, 0x00, 0x07, // coord
		0x00, 0x08, // ps
	}

	fvar := GetFvar(data, 0)
	if fvar == nil {
		t.Fatal("fvar should not be nil")
	}
	if fvar.Version != "1.0" {
		t.Errorf("version want 1.0 got %s", fvar.Version)
	}
	if fvar.AxisCount != 1 || len(fvar.Axis) != 1 {
		t.Errorf("axisCount want 1 got %d", fvar.AxisCount)
	}
	if fvar.InstanceCount != 1 || len(fvar.Instance) != 1 {
		t.Errorf("instanceCount want 1 got %d", fvar.InstanceCount)
	}
	// Assert axis field values
	axis := fvar.Axis[0]
	if axis.AxisTag != 0x12345678 {
		t.Errorf("axis.AxisTag want 0x12345678 got 0x%08X", axis.AxisTag)
	}
	if axis.MinValue != 0x00000000 {
		t.Errorf("axis.MinValue want 0x00000000 got 0x%08X", axis.MinValue)
	}
	if axis.DefaultValue != 0x00000001 {
		t.Errorf("axis.DefaultValue want 0x00000001 got 0x%08X", axis.DefaultValue)
	}
	if axis.MaxValue != 0x00000002 {
		t.Errorf("axis.MaxValue want 0x00000002 got 0x%08X", axis.MaxValue)
	}
	if axis.Flags != 0x0003 {
		t.Errorf("axis.Flags want 0x0003 got 0x%04X", axis.Flags)
	}
	if axis.NameID != 0x0004 {
		t.Errorf("axis.NameID want 0x0004 got 0x%04X", axis.NameID)
	}

	// Assert instance field values
	inst := fvar.Instance[0]
	if inst.NameID != 0x0005 {
		t.Errorf("instance.NameID want 0x0005 got 0x%04X", inst.NameID)
	}
	if inst.Flags != 0x0006 {
		t.Errorf("instance.Flags want 0x0006 got 0x%04X", inst.Flags)
	}
	if len(inst.Coordinates) != 1 || inst.Coordinates[0] != 0x00000007 {
		t.Errorf("instance.Coordinates[0] want 0x00000007 got %v", inst.Coordinates)
	}
	if inst.PsNameID != 0x0008 {
		t.Errorf("instance.PsNameID want 0x0008 got 0x%04X", inst.PsNameID)
	}
}

// TestGetFvarMultiAxis verifies parsing of multiple axes and multiple coordinates per instance.
func TestGetFvarMultiAxis(t *testing.T) {
	// Construct spec-like fvar table:
	// Header (16 bytes): version(4) offsetToData(2) countSizePairs(2) axisCount(2) axisSize(2) instanceCount(2) instanceSize(2)
	// Two axes (each 20 bytes): 'wght' and 'wdth'
	// Two instances (each 14 bytes: nameID(2) flags(2) coords(2*4) psNameID(2))
	data := []byte{
		0x00, 0x01, 0x00, 0x00, // version 1.0
		0x00, 0x10, // offsetToData = 16
		0x00, 0x01, // countSizePairs = 1
		0x00, 0x02, // axisCount = 2
		0x00, 0x14, // axisSize = 20
		0x00, 0x02, // instanceCount = 2
		0x00, 0x0E, // instanceSize = 14 (with psNameID)
		// Axis 0 'wght'
		'w', 'g', 'h', 't', // tag
		0x00, 0x64, 0x00, 0x00, // min 100.0 -> 0x00640000
		0x01, 0x90, 0x00, 0x00, // default 400.0 -> 0x01900000
		0x03, 0x84, 0x00, 0x00, // max 900.0 -> 0x03840000
		0x00, 0x00, // flags
		0x00, 0x01, // nameID
		// Axis 1 'wdth'
		'w', 'd', 't', 'h', // tag
		0x00, 0x32, 0x00, 0x00, // min 50.0 -> 0x00320000
		0x00, 0x64, 0x00, 0x00, // default 100.0 -> 0x00640000
		0x00, 0xC8, 0x00, 0x00, // max 200.0 -> 0x00C80000
		0x00, 0x00, // flags
		0x00, 0x02, // nameID
		// Instance 0
		0x00, 0x0A, // nameID
		0x00, 0x00, // flags
		0x01, 0x90, 0x00, 0x00, // wght = 400.0
		0x00, 0x64, 0x00, 0x00, // wdth = 100.0
		0x00, 0x14, // psNameID
		// Instance 1
		0x00, 0x0B, // nameID
		0x00, 0x00, // flags
		0x00, 0xC8, 0x00, 0x00, // wght = 200.0
		0x00, 0x64, 0x00, 0x00, // wdth = 100.0
		0x00, 0x15, // psNameID
	}

	fvar := GetFvar(data, 0)
	if fvar == nil {
		t.Fatal("fvar should not be nil for valid data")
	}
	if fvar.AxisCount != 2 || len(fvar.Axis) != 2 {
		t.Fatalf("expected 2 axes, got %d", len(fvar.Axis))
	}
	if fvar.InstanceCount != 2 || len(fvar.Instance) != 2 {
		t.Fatalf("expected 2 instances, got %d", len(fvar.Instance))
	}
	// Axis 0 checks
	axis0 := fvar.Axis[0]
	if axis0.AxisTag != uint32('w')<<24|uint32('g')<<16|uint32('h')<<8|uint32('t') {
		t.Errorf("axis0 tag mismatch: %08X", axis0.AxisTag)
	}
	if axis0.MinValue != 0x00640000 || axis0.DefaultValue != 0x01900000 || axis0.MaxValue != 0x03840000 {
		t.Errorf("axis0 values mismatch")
	}
	// Axis 1 checks
	axis1 := fvar.Axis[1]
	if axis1.MinValue != 0x00320000 || axis1.DefaultValue != 0x00640000 || axis1.MaxValue != 0x00C80000 {
		t.Errorf("axis1 values mismatch")
	}
	// Instance 0
	inst0 := fvar.Instance[0]
	if len(inst0.Coordinates) != 2 || inst0.Coordinates[0] != 0x01900000 || inst0.Coordinates[1] != 0x00640000 {
		t.Errorf("instance0 coordinates mismatch: %v", inst0.Coordinates)
	}
	if inst0.PsNameID != 0x0014 {
		t.Errorf("instance0 psNameID mismatch: %04X", inst0.PsNameID)
	}
	// Instance 1
	inst1 := fvar.Instance[1]
	if len(inst1.Coordinates) != 2 || inst1.Coordinates[0] != 0x00C80000 || inst1.Coordinates[1] != 0x00640000 {
		t.Errorf("instance1 coordinates mismatch: %v", inst1.Coordinates)
	}
	if inst1.PsNameID != 0x0015 {
		t.Errorf("instance1 psNameID mismatch: %04X", inst1.PsNameID)
	}
}

// TestGetFvarMalformed ensures truncated data returns nil.
func TestGetFvarMalformed(t *testing.T) {
	// Truncate after header (no axis/instance data)
	data := []byte{
		0x00, 0x01, 0x00, 0x00,
		0x00, 0x10,
		0x00, 0x01,
		0x00, 0x01,
		0x00, 0x14,
		0x00, 0x01,
		0x00, 0x0A,
	}
	fvar := GetFvar(data, 0)
	if fvar != nil {
		// Expect nil due to missing axis bytes
		t.Errorf("expected nil for malformed/truncated data")
	}
}

// TestGetKernWindows tests Windows format kern table (version 0)
func TestGetKernWindows(t *testing.T) {
	data := getWindowsKernData()
	kern, err := GetKern(data, 0)
	if err != nil {
		t.Fatalf("GetKern failed: %v", err)
	}

	// Verify kern table header
	if kern.Version != 0 {
		t.Errorf("Expected Version=0, got %d", kern.Version)
	}
	if kern.NTables != 1 {
		t.Errorf("Expected NTables=1, got %d", kern.NTables)
	}
	if kern.IsMacNewKern {
		t.Errorf("Expected IsMacNewKern=false for Windows format")
	}

	// Verify subtable headers
	if kern.SubHeaders["version"] != 0 {
		t.Errorf("Expected subtable version=0, got %d", kern.SubHeaders["version"])
	}
	if kern.SubHeaders["length"] != 26 {
		t.Errorf("Expected subtable length=26, got %d", kern.SubHeaders["length"])
	}
	if kern.SubHeaders["coverage"] != 0 {
		t.Errorf("Expected coverage=0, got %d", kern.SubHeaders["coverage"])
	}
	if kern.SubHeaders["format"] != 0 {
		t.Errorf("Expected format=0, got %d", kern.SubHeaders["format"])
	}
	if kern.SubHeaders["nPairs"] != 2 {
		t.Errorf("Expected nPairs=2, got %d", kern.SubHeaders["nPairs"])
	}

	// Verify kerning pairs
	if len(kern.Pairs) != 2 {
		t.Fatalf("Expected 2 kerning pairs, got %d", len(kern.Pairs))
	}

	// Pair 1: A-V = -50
	pair1 := kern.Pairs[0]
	if pair1.Left != 65 || pair1.Right != 86 || pair1.Value != -50 {
		t.Errorf("Pair 1 expected (65,86,-50), got (%d,%d,%d)", pair1.Left, pair1.Right, pair1.Value)
	}

	// Pair 2: F-. = -30
	pair2 := kern.Pairs[1]
	if pair2.Left != 70 || pair2.Right != 46 || pair2.Value != -30 {
		t.Errorf("Pair 2 expected (70,46,-30), got (%d,%d,%d)", pair2.Left, pair2.Right, pair2.Value)
	}

	// Verify using binary.BigEndian for independent validation
	if binary.BigEndian.Uint16(data[0:2]) != 0 {
		t.Error("Version should be 0")
	}
	if binary.BigEndian.Uint16(data[2:4]) != 1 {
		t.Error("NTables should be 1")
	}
	// First pair left glyph
	if binary.BigEndian.Uint16(data[18:20]) != 65 {
		t.Error("First pair left should be 65")
	}
}

// TestGetKernMac tests Mac format kern table (version 1, old format)
func TestGetKernMac(t *testing.T) {
	data := getMacKernData()
	kern, err := GetKern(data, 0)
	if err != nil {
		t.Fatalf("GetKern failed: %v", err)
	}

	// Verify kern table header
	if kern.Version != 1 {
		t.Errorf("Expected Version=1, got %d", kern.Version)
	}
	if kern.NTables != 1 {
		t.Errorf("Expected NTables=1, got %d", kern.NTables)
	}
	if kern.IsMacNewKern {
		t.Errorf("Expected IsMacNewKern=false for old Mac format")
	}

	// Verify subtable headers
	if kern.SubHeaders["length"] != 32 {
		t.Errorf("Expected subtable length=32, got %d", kern.SubHeaders["length"])
	}
	if kern.SubHeaders["coverage"] != 0x8000 {
		t.Errorf("Expected coverage=0x8000, got %d", kern.SubHeaders["coverage"])
	}
	if kern.SubHeaders["nPairs"] != 2 {
		t.Errorf("Expected nPairs=2, got %d", kern.SubHeaders["nPairs"])
	}

	// Verify kerning pairs
	if len(kern.Pairs) != 2 {
		t.Fatalf("Expected 2 kerning pairs, got %d", len(kern.Pairs))
	}

	// Pair 1: T-o = -40
	pair1 := kern.Pairs[0]
	if pair1.Left != 84 || pair1.Right != 111 || pair1.Value != -40 {
		t.Errorf("Pair 1 expected (84,111,-40), got (%d,%d,%d)", pair1.Left, pair1.Right, pair1.Value)
	}

	// Pair 2: W-a = -20
	pair2 := kern.Pairs[1]
	if pair2.Left != 87 || pair2.Right != 97 || pair2.Value != -20 {
		t.Errorf("Pair 2 expected (87,97,-20), got (%d,%d,%d)", pair2.Left, pair2.Right, pair2.Value)
	}

	// Verify using binary.BigEndian for independent validation
	if binary.BigEndian.Uint16(data[0:2]) != 1 {
		t.Error("Version should be 1")
	}
	if binary.BigEndian.Uint16(data[2:4]) != 1 {
		t.Error("NTables should be 1")
	}
}

// TestGetKernMacNew tests Mac format kern table (version 1, new format)
func TestGetKernMacNew(t *testing.T) {
	data := getMacNewKernData()
	kern, err := GetKern(data, 0)
	if err != nil {
		t.Fatalf("GetKern failed: %v", err)
	}

	// Verify kern table header
	if kern.Version != 1 {
		t.Errorf("Expected Version=1, got %d", kern.Version)
	}
	if kern.NTables != 0 { // nTables is 0 in header for new format
		t.Errorf("Expected NTables=0 (new format indicator), got %d", kern.NTables)
	}
	if !kern.IsMacNewKern {
		t.Errorf("Expected IsMacNewKern=true for new Mac format")
	}

	// Verify subtable headers
	if kern.SubHeaders["length"] != 32 {
		t.Errorf("Expected subtable length=32, got %d", kern.SubHeaders["length"])
	}
	if kern.SubHeaders["nPairs"] != 1 {
		t.Errorf("Expected nPairs=1, got %d", kern.SubHeaders["nPairs"])
	}

	// Verify kerning pairs
	if len(kern.Pairs) != 1 {
		t.Fatalf("Expected 1 kerning pair, got %d", len(kern.Pairs))
	}

	// Pair 1: L-Y = -80
	pair1 := kern.Pairs[0]
	if pair1.Left != 76 || pair1.Right != 89 || pair1.Value != -80 {
		t.Errorf("Pair expected (76,89,-80), got (%d,%d,%d)", pair1.Left, pair1.Right, pair1.Value)
	}

	// Verify using binary.BigEndian for independent validation
	if binary.BigEndian.Uint16(data[0:2]) != 1 {
		t.Error("Version should be 1")
	}
	if binary.BigEndian.Uint16(data[2:4]) != 0 {
		t.Error("NTables should be 0 (new format)")
	}
	// Actual nTables (32-bit)
	if binary.BigEndian.Uint32(data[4:8]) != 1 {
		t.Error("Actual nTables should be 1")
	}
}

// TestGetKernUnsupportedVersion tests unsupported kern table version
func TestGetKernUnsupportedVersion(t *testing.T) {
	// Create kern data with unsupported version 2
	data := []byte{
		0x00, 0x02, // version = 2 (unsupported)
		0x00, 0x00, // nTables = 0
	}
	kern, err := GetKern(data, 0)
	if err == nil {
		t.Error("Expected error for unsupported kern version")
	}
	if kern == nil {
		t.Error("Expected kern struct to be non-nil even on error")
	}
}

// TestGetKernWindowsFormat2 tests Windows format 2 kern table (n×m array)
func TestGetKernWindowsFormat2(t *testing.T) {
	data := getWindowsFormat2KernData()
	kern, err := GetKern(data, 0)
	if err != nil {
		t.Fatalf("GetKern failed: %v", err)
	}

	// Verify kern table header
	if kern.Version != 0 {
		t.Errorf("Expected Version=0, got %d", kern.Version)
	}
	if kern.NTables != 1 {
		t.Errorf("Expected NTables=1, got %d", kern.NTables)
	}

	// Verify subtable headers
	if kern.SubHeaders["format"] != 2 {
		t.Errorf("Expected format=2, got %d", kern.SubHeaders["format"])
	}
	if kern.SubHeaders["coverage"] != 2 {
		t.Errorf("Expected coverage=2, got %d", kern.SubHeaders["coverage"])
	}

	// Verify format 2 data exists
	if kern.Format2 == nil {
		t.Fatal("Expected Format2 data to be non-nil")
	}

	// Verify format 2 structure
	if kern.Format2.RowWidth != 4 {
		t.Errorf("Expected rowWidth=4, got %d", kern.Format2.RowWidth)
	}
	if kern.Format2.LeftOffsetTable != 14 {
		t.Errorf("Expected leftOffsetTable=14, got %d", kern.Format2.LeftOffsetTable)
	}
	if kern.Format2.RightOffsetTable != 22 {
		t.Errorf("Expected rightOffsetTable=22, got %d", kern.Format2.RightOffsetTable)
	}

	// Verify left class table
	if kern.Format2.LeftClassTable == nil {
		t.Fatal("Expected LeftClassTable to be non-nil")
	}
	if kern.Format2.LeftClassTable.FirstGlyph != 65 {
		t.Errorf("Expected left firstGlyph=65, got %d", kern.Format2.LeftClassTable.FirstGlyph)
	}
	if kern.Format2.LeftClassTable.NGlyphs != 2 {
		t.Errorf("Expected left nGlyphs=2, got %d", kern.Format2.LeftClassTable.NGlyphs)
	}

	// Verify right class table
	if kern.Format2.RightClassTable == nil {
		t.Fatal("Expected RightClassTable to be non-nil")
	}
	if kern.Format2.RightClassTable.FirstGlyph != 86 {
		t.Errorf("Expected right firstGlyph=86, got %d", kern.Format2.RightClassTable.FirstGlyph)
	}
	if kern.Format2.RightClassTable.NGlyphs != 2 {
		t.Errorf("Expected right nGlyphs=2, got %d", kern.Format2.RightClassTable.NGlyphs)
	}
}

// TestGetKernMacFormat2 tests Mac format 2 kern table (n×m array)
func TestGetKernMacFormat2(t *testing.T) {
	data := getMacFormat2KernData()
	kern, err := GetKern(data, 0)
	if err != nil {
		t.Fatalf("GetKern failed: %v", err)
	}

	// Verify kern table header
	if kern.Version != 1 {
		t.Errorf("Expected Version=1, got %d", kern.Version)
	}
	if kern.NTables != 1 {
		t.Errorf("Expected NTables=1, got %d", kern.NTables)
	}

	// Verify subtable headers
	if kern.SubHeaders["format"] != 2 {
		t.Errorf("Expected format=2, got %d", kern.SubHeaders["format"])
	}
	if kern.SubHeaders["coverage"] != 0x8002 {
		t.Errorf("Expected coverage=0x8002, got %d", kern.SubHeaders["coverage"])
	}

	// Verify format 2 data exists
	if kern.Format2 == nil {
		t.Fatal("Expected Format2 data to be non-nil")
	}

	// Verify format 2 structure
	if kern.Format2.RowWidth != 4 {
		t.Errorf("Expected rowWidth=4, got %d", kern.Format2.RowWidth)
	}

	// Verify left class table
	if kern.Format2.LeftClassTable == nil {
		t.Fatal("Expected LeftClassTable to be non-nil")
	}
	if kern.Format2.LeftClassTable.FirstGlyph != 84 {
		t.Errorf("Expected left firstGlyph=84 ('T'), got %d", kern.Format2.LeftClassTable.FirstGlyph)
	}
	if kern.Format2.LeftClassTable.NGlyphs != 2 {
		t.Errorf("Expected left nGlyphs=2, got %d", kern.Format2.LeftClassTable.NGlyphs)
	}

	// Verify right class table
	if kern.Format2.RightClassTable == nil {
		t.Fatal("Expected RightClassTable to be non-nil")
	}
	if kern.Format2.RightClassTable.FirstGlyph != 111 {
		t.Errorf("Expected right firstGlyph=111 ('o'), got %d", kern.Format2.RightClassTable.FirstGlyph)
	}
}

// TestGetKernMacFormat3 tests Mac format 3 kern table (n×m index array)
func TestGetKernMacFormat3(t *testing.T) {
	data := getMacFormat3KernData()
	kern, err := GetKern(data, 0)
	if err != nil {
		t.Fatalf("GetKern failed: %v", err)
	}

	// Verify kern table header
	if kern.Version != 1 {
		t.Errorf("Expected Version=1, got %d", kern.Version)
	}
	if kern.NTables != 1 {
		t.Errorf("Expected NTables=1, got %d", kern.NTables)
	}

	// Verify subtable headers
	if kern.SubHeaders["format"] != 3 {
		t.Errorf("Expected format=3, got %d", kern.SubHeaders["format"])
	}

	// Verify format 3 data exists
	if kern.Format3 == nil {
		t.Fatal("Expected Format3 data to be non-nil")
	}

	// Verify format 3 structure
	if kern.Format3.GlyphCount != 4 {
		t.Errorf("Expected glyphCount=4, got %d", kern.Format3.GlyphCount)
	}
	if kern.Format3.KernValueCount != 2 {
		t.Errorf("Expected kernValueCount=2, got %d", kern.Format3.KernValueCount)
	}
	if kern.Format3.LeftClassCount != 2 {
		t.Errorf("Expected leftClassCount=2, got %d", kern.Format3.LeftClassCount)
	}
	if kern.Format3.RightClassCount != 2 {
		t.Errorf("Expected rightClassCount=2, got %d", kern.Format3.RightClassCount)
	}

	// Verify kern values
	if len(kern.Format3.KernValues) != 2 {
		t.Fatalf("Expected 2 kern values, got %d", len(kern.Format3.KernValues))
	}
	if kern.Format3.KernValues[0] != -50 {
		t.Errorf("Expected kernValue[0]=-50, got %d", kern.Format3.KernValues[0])
	}
	if kern.Format3.KernValues[1] != -20 {
		t.Errorf("Expected kernValue[1]=-20, got %d", kern.Format3.KernValues[1])
	}

	// Verify left class array
	if len(kern.Format3.LeftClass) != 4 {
		t.Errorf("Expected 4 left class entries, got %d", len(kern.Format3.LeftClass))
	}

	// Verify right class array
	if len(kern.Format3.RightClass) != 4 {
		t.Errorf("Expected 4 right class entries, got %d", len(kern.Format3.RightClass))
	}

	// Verify kern index array
	expectedIndexCount := 2 * 2 // leftClassCount × rightClassCount
	if len(kern.Format3.KernIndex) != expectedIndexCount {
		t.Errorf("Expected %d kern index entries, got %d", expectedIndexCount, len(kern.Format3.KernIndex))
	}
}
