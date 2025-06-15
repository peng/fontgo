package font

import (
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

}
