package font

import (
	"fmt"
	"os"
	"reflect"
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
	/*
	   [{"tag":"COLR","checksum":3186909341,"offset":332,"length":109354,"compression":false},{"tag":"CPAL","checksum":2465690482,"offset":109688,"length":90,"compression":false},{"tag":"DSIG","checksum":1,"offset":16938676,"length":8,"compression":false},{"tag":"GDEF","checksum":1059537,"offset":109780,"length":22,"compression":false},{"tag":"GPOS","checksum":2215742747,"offset":109804,"length":112,"compression":false},{"tag":"GSUB","checksum":3410809594,"offset":109916,"length":298,"compression":false},{"tag":"OS/2","checksum":1137592327,"offset":110216,"length":96,"compression":false},{"tag":"cmap","checksum":3086641967,"offset":110312,"length":54588,"compression":false},{"tag":"cvt ","checksum":372968533,"offset":16934848,"length":52,"compression":false},{"tag":"fpgm","checksum":2654343626,"offset":16934900,"length":3605,"compression":false},{"tag":"gasp","checksum":16,"offset":16934840,"length":8,"compression":false},{"tag":"glyf","checksum":3580581701,"offset":164900,"length":16572502,"compression":false},{"tag":"head","checksum":473240010,"offset":16737404,"length":54,"compression":false},{"tag":"hhea","checksum":320219366,"offset":16737460,"length":36,"compression":false},{"tag":"hmtx","checksum":3408615136,"offset":16737496,"length":43836,"compression":false},{"tag":"loca","checksum":1056580002,"offset":16781332,"length":43848,"compression":false},{"tag":"maxp","checksum":802902503,"offset":16825180,"length":32,"compression":false},{"tag":"name","checksum":298926501,"offset":16825212,"length":1020,"compression":false},{"tag":"post","checksum":745911922,"offset":16826232,"length":108608,"compression":false},{"tag":"prep","checksum":1749469340,"offset":16938508,"length":167,"compression":false}]
	*/
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

	// if *head != *standData.Head {
	// 	fmt.Println("standData head", standData.Head)
	// 	fmt.Println("source head", head)
	// 	if head.Version != standData.Head.Version {
	// 		fmt.Println("standData head Version", standData.Head.Version)
	// 		fmt.Println("source head Version", head.Version)
	// 		t.Log("head Version error")
	// 		t.Fail()
	// 		return
	// 	}

	// 	if head.FontRevision != standData.Head.FontRevision {
	// 		fmt.Println("standData head FontRevision", standData.Head.FontRevision)
	// 		fmt.Println("source head FontRevision", head.FontRevision)
	// 		t.Log("head FontRevision error")
	// 		t.Fail()
	// 		return
	// 	}

	// 	if head.CheckSumAdjustment != standData.Head.CheckSumAdjustment {
	// 		fmt.Println("standData head CheckSumAdjustment", standData.Head.CheckSumAdjustment)
	// 		fmt.Println("source head CheckSumAdjustment", head.CheckSumAdjustment)
	// 		t.Log("head CheckSumAdjustment error")
	// 		t.Fail()
	// 		return
	// 	}

	// 	if head.MagicNumber != standData.Head.MagicNumber {
	// 		fmt.Println("standData head MagicNumber", standData.Head.MagicNumber)
	// 		fmt.Println("source head MagicNumber", head.MagicNumber)
	// 		t.Log("head MagicNumber error")
	// 		t.Fail()
	// 		return
	// 	}

	// 	if head.Flags != standData.Head.Flags {
	// 		fmt.Println("standData head Flags", standData.Head.Flags)
	// 		fmt.Println("source head Flags", head.Flags)
	// 		t.Log("head Flags error")
	// 		t.Fail()
	// 		return
	// 	}

	// 	if head.UnitsPerEm != standData.Head.UnitsPerEm {
	// 		fmt.Println("standData head UnitsPerEm", standData.Head.UnitsPerEm)
	// 		fmt.Println("source head UnitsPerEm", head.UnitsPerEm)
	// 		t.Log("head UnitsPerEm error")
	// 		t.Fail()
	// 		return
	// 	}

	// 	if head.Created != standData.Head.Created/1000 {
	// 		fmt.Println("standData head Created", standData.Head.Created)
	// 		fmt.Println("source head Created", head.Created)
	// 		t.Log("head Created error")
	// 		t.Fail()
	// 		return
	// 	}

	// 	if head.Modified != standData.Head.Modified/1000 {
	// 		fmt.Println("standData head Modified", standData.Head.Modified)
	// 		fmt.Println("source head Modified", head.Modified)
	// 		t.Log("head Modified error")
	// 		t.Fail()
	// 		return
	// 	}

	// 	if head.XMin != standData.Head.XMin {
	// 		fmt.Println("standData head XMin", standData.Head.XMin)
	// 		fmt.Println("source head XMin", head.XMin)
	// 		t.Log("head XMin error")
	// 		t.Fail()
	// 		return
	// 	}
	// 	if head.YMin != standData.Head.YMin {
	// 		fmt.Println("standData head YMin", standData.Head.YMin)
	// 		fmt.Println("source head YMin", head.YMin)
	// 		t.Log("head YMin error")
	// 		t.Fail()
	// 		return
	// 	}

	// 	if head.XMax != standData.Head.XMax {
	// 		fmt.Println("standData head XMax", standData.Head.XMax)
	// 		fmt.Println("source head XMax", head.XMax)
	// 		t.Log("head XMax error")
	// 		t.Fail()
	// 		return
	// 	}

	// 	if head.YMax != standData.Head.YMax {
	// 		fmt.Println("standData head YMax", standData.Head.YMax)
	// 		fmt.Println("source head YMax", head.YMax)
	// 		t.Log("head YMax error")
	// 		t.Fail()
	// 		return
	// 	}

	// 	if head.MacStyle != standData.Head.MacStyle {
	// 		fmt.Println("standData head MacStyle", standData.Head.MacStyle)
	// 		fmt.Println("source head MacStyle", head.MacStyle)
	// 		t.Log("head MacStyle error")
	// 		t.Fail()
	// 		return
	// 	}

	// 	if head.LowestRecPPEM != standData.Head.LowestRecPPEM {
	// 		fmt.Println("standData head LowestRecPPEM", standData.Head.LowestRecPPEM)
	// 		fmt.Println("source head LowestRecPPEM", head.LowestRecPPEM)
	// 		t.Log("head LowestRecPPEM error")
	// 		t.Fail()
	// 		return
	// 	}

	// 	if head.FontDirectionHint != standData.Head.FontDirectionHint {
	// 		fmt.Println("standData head FontDirectionHint", standData.Head.FontDirectionHint)
	// 		fmt.Println("source head FontDirectionHint", head.FontDirectionHint)
	// 		t.Log("head FontDirectionHint error")
	// 		t.Fail()
	// 		return
	// 	}
	// 	if head.IndexToLocFormat != standData.Head.IndexToLocFormat {
	// 		fmt.Println("standData head IndexToLocFormat", standData.Head.IndexToLocFormat)
	// 		fmt.Println("source head IndexToLocFormat", head.IndexToLocFormat)
	// 		t.Log("head IndexToLocFormat error")
	// 		t.Fail()
	// 		return
	// 	}
	// 	if head.GlyphDataFormat != standData.Head.GlyphDataFormat {
	// 		fmt.Println("standData head GlyphDataFormat", standData.Head.GlyphDataFormat)
	// 		fmt.Println("source head GlyphDataFormat", head.GlyphDataFormat)
	// 		t.Log("head GlyphDataFormat error")
	// 		t.Fail()
	// 		return
	// 	}
	// }

	// // check Maxp table
	// standMaxp := &Maxp{"1.0", 10961, 977, 55, 0, 0, 2, 152, 252, 141, 0, 941, 19736, 0, 0}
	// maxpInfo := tableContent["maxp"]
	// maxp := GetMaxp(fileByte[maxpInfo.Offset : maxpInfo.Offset+maxpInfo.Length])
	// if *maxp != *standMaxp {
	// 	fmt.Println("standMaxp", standMaxp)
	// 	fmt.Println("source Maxp", maxp)
	// }

	// // check local table
	// // standLocal := &
}
