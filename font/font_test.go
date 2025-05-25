package font

// func TestParse(t *testing.T) {
// 	dir, err := DataReader("../test/HanyiSentyCrayon.ttf")
// 	if err == nil {
// 		var dirStr []byte
// 		dirStr, err = json.Marshal(dir)
// 		if err != nil {
// 			t.Log(err)
// 			t.Fail()
// 		}
// 		fmt.Printf("%v", string(dirStr))
// 		return
// 	}

// 	t.Log(err)
// 	t.Fail()

// }

// type File struct {
// 	Pos int `json:"pos"`
// }
// type RespData struct {
// 	File        *File        `json:"file"`
// 	OffsetTable *OffsetTable `json:"offsetTable"`
// 	Head        *Head        `json:"head"`
// }

// func GetRemoteJson() (data *RespData, err error) {
// 	var (
// 		resp *http.Response
// 		body []byte
// 	)
// 	url := "http://127.0.0.1:3111/getfontjson"
// 	resp, err = http.Get(url)
// 	if err != nil {
// 		return
// 	}

// 	body, err = ioutil.ReadAll(resp.Body)

// 	if err != nil {
// 		return
// 	}

// 	err = json.Unmarshal(body, &data)
// 	return
// }

// func TestAllTable(t *testing.T) {
// 	var (
// 		fileByte  []byte
// 		err       error
// 		standData *RespData
// 	)

// 	fileByte, err = os.ReadFile("../test/HanyiSentyCrayon.ttf")

// 	if err != nil {
// 		t.Log(err)
// 		t.Fail()
// 		return
// 	}

// 	standData, err = GetRemoteJson()
// 	if err != nil {
// 		t.Log(err)
// 		t.Fail()
// 		return
// 	}

// 	// test offsetTable
// 	offsetTable := GetOffsetTable(fileByte)

// 	if err != nil {
// 		t.Log(err)
// 		t.Fail()
// 		return
// 	}
// 	if *offsetTable != *standData.OffsetTable {
// 		fmt.Println("standData offsetTable", standData.OffsetTable)
// 		fmt.Println("source offsetTable", offsetTable)
// 		t.Log("offsetTable error")
// 		t.Fail()
// 		return
// 	}

// 	// read table content
// 	numTables := int(offsetTable.NumTables)
// 	tableContent := GetTableContent(numTables, fileByte)

// 	// check head table
// 	headInfo := tableContent["head"]
// 	head := GetHead(fileByte[headInfo.Offset : headInfo.Offset+headInfo.Length])

// 	if *head != *standData.Head {
// 		fmt.Println("standData head", standData.Head)
// 		fmt.Println("source head", head)
// 		if head.Version != standData.Head.Version {
// 			fmt.Println("standData head Version", standData.Head.Version)
// 			fmt.Println("source head Version", head.Version)
// 			t.Log("head Version error")
// 			t.Fail()
// 			return
// 		}

// 		if head.FontRevision != standData.Head.FontRevision {
// 			fmt.Println("standData head FontRevision", standData.Head.FontRevision)
// 			fmt.Println("source head FontRevision", head.FontRevision)
// 			t.Log("head FontRevision error")
// 			t.Fail()
// 			return
// 		}

// 		if head.CheckSumAdjustment != standData.Head.CheckSumAdjustment {
// 			fmt.Println("standData head CheckSumAdjustment", standData.Head.CheckSumAdjustment)
// 			fmt.Println("source head CheckSumAdjustment", head.CheckSumAdjustment)
// 			t.Log("head CheckSumAdjustment error")
// 			t.Fail()
// 			return
// 		}

// 		if head.MagicNumber != standData.Head.MagicNumber {
// 			fmt.Println("standData head MagicNumber", standData.Head.MagicNumber)
// 			fmt.Println("source head MagicNumber", head.MagicNumber)
// 			t.Log("head MagicNumber error")
// 			t.Fail()
// 			return
// 		}

// 		if head.Flags != standData.Head.Flags {
// 			fmt.Println("standData head Flags", standData.Head.Flags)
// 			fmt.Println("source head Flags", head.Flags)
// 			t.Log("head Flags error")
// 			t.Fail()
// 			return
// 		}

// 		if head.UnitsPerEm != standData.Head.UnitsPerEm {
// 			fmt.Println("standData head UnitsPerEm", standData.Head.UnitsPerEm)
// 			fmt.Println("source head UnitsPerEm", head.UnitsPerEm)
// 			t.Log("head UnitsPerEm error")
// 			t.Fail()
// 			return
// 		}

// 		if head.Created != standData.Head.Created/1000 {
// 			fmt.Println("standData head Created", standData.Head.Created)
// 			fmt.Println("source head Created", head.Created)
// 			t.Log("head Created error")
// 			t.Fail()
// 			return
// 		}

// 		if head.Modified != standData.Head.Modified/1000 {
// 			fmt.Println("standData head Modified", standData.Head.Modified)
// 			fmt.Println("source head Modified", head.Modified)
// 			t.Log("head Modified error")
// 			t.Fail()
// 			return
// 		}

// 		if head.XMin != standData.Head.XMin {
// 			fmt.Println("standData head XMin", standData.Head.XMin)
// 			fmt.Println("source head XMin", head.XMin)
// 			t.Log("head XMin error")
// 			t.Fail()
// 			return
// 		}
// 		if head.YMin != standData.Head.YMin {
// 			fmt.Println("standData head YMin", standData.Head.YMin)
// 			fmt.Println("source head YMin", head.YMin)
// 			t.Log("head YMin error")
// 			t.Fail()
// 			return
// 		}

// 		if head.XMax != standData.Head.XMax {
// 			fmt.Println("standData head XMax", standData.Head.XMax)
// 			fmt.Println("source head XMax", head.XMax)
// 			t.Log("head XMax error")
// 			t.Fail()
// 			return
// 		}

// 		if head.YMax != standData.Head.YMax {
// 			fmt.Println("standData head YMax", standData.Head.YMax)
// 			fmt.Println("source head YMax", head.YMax)
// 			t.Log("head YMax error")
// 			t.Fail()
// 			return
// 		}

// 		if head.MacStyle != standData.Head.MacStyle {
// 			fmt.Println("standData head MacStyle", standData.Head.MacStyle)
// 			fmt.Println("source head MacStyle", head.MacStyle)
// 			t.Log("head MacStyle error")
// 			t.Fail()
// 			return
// 		}

// 		if head.LowestRecPPEM != standData.Head.LowestRecPPEM {
// 			fmt.Println("standData head LowestRecPPEM", standData.Head.LowestRecPPEM)
// 			fmt.Println("source head LowestRecPPEM", head.LowestRecPPEM)
// 			t.Log("head LowestRecPPEM error")
// 			t.Fail()
// 			return
// 		}

// 		if head.FontDirectionHint != standData.Head.FontDirectionHint {
// 			fmt.Println("standData head FontDirectionHint", standData.Head.FontDirectionHint)
// 			fmt.Println("source head FontDirectionHint", head.FontDirectionHint)
// 			t.Log("head FontDirectionHint error")
// 			t.Fail()
// 			return
// 		}
// 		if head.IndexToLocFormat != standData.Head.IndexToLocFormat {
// 			fmt.Println("standData head IndexToLocFormat", standData.Head.IndexToLocFormat)
// 			fmt.Println("source head IndexToLocFormat", head.IndexToLocFormat)
// 			t.Log("head IndexToLocFormat error")
// 			t.Fail()
// 			return
// 		}
// 		if head.GlyphDataFormat != standData.Head.GlyphDataFormat {
// 			fmt.Println("standData head GlyphDataFormat", standData.Head.GlyphDataFormat)
// 			fmt.Println("source head GlyphDataFormat", head.GlyphDataFormat)
// 			t.Log("head GlyphDataFormat error")
// 			t.Fail()
// 			return
// 		}
// 	}

// 	// check Maxp table
// 	standMaxp := &Maxp{"1.0", 10961, 977, 55, 0, 0, 2, 152, 252, 141, 0, 941, 19736, 0, 0}
// 	maxpInfo := tableContent["maxp"]
// 	maxp := GetMaxp(fileByte[maxpInfo.Offset : maxpInfo.Offset + maxpInfo.Length])
// 	if *maxp != *standMaxp {
// 		fmt.Println("standMaxp", standMaxp)
// 		fmt.Println("source Maxp", maxp)
// 	}

// 	// check local table
// 	// standLocal := &
// }
