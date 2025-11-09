package font

import (
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
}
