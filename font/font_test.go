package font

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

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

type File struct {
	Pos int `json:"pos"`
}
type RespData struct {
	File File `json:"file"`
	OffsetTable OffsetTable `json:"offsetTable"`
}

func GetRemoteJson() (data RespData, err error){
	var (
		resp *http.Response
		body []byte
	)
	url := "http://127.0.0.1:3111/getfontjson"
	resp, err = http.Get(url)
	if err != nil {
		return
	}
	
	body, err = ioutil.ReadAll(resp.Body)
	
	if err != nil {
		return
	}
	
	err = json.Unmarshal(body, &data)
	return
}

func TestGetOffsetTable(t *testing.T) {
	var (
		fileByte []byte
		err error
		standData RespData
	)

	fileByte, err = os.ReadFile("../test/HanyiSentyCrayon.ttf")

	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	offsetTable := GetOffsetTable(fileByte)
	
	standData, err = GetRemoteJson()
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	if (standData.OffsetTable.EntrySelector != offsetTable.EntrySelector) {
		t.Log("EntrySelector error")
		t.Fail()
		return
	}

	if (standData.OffsetTable.NumTables != offsetTable.NumTables) {
		t.Log("NumTables error")
		t.Fail()
		return
	}

	if (standData.OffsetTable.RangeShift != offsetTable.RangeShift) {
		t.Log("RangeShift error")
		t.Fail()
		return
	}

	if (standData.OffsetTable.SearchRange != offsetTable.SearchRange) {
		t.Log("SearchRange error")
		t.Fail()
		return
	}

	if (standData.OffsetTable.ScalerType != offsetTable.ScalerType) {
		t.Log("ScalerType error", standData.OffsetTable.ScalerType, offsetTable.ScalerType)
		t.Fail()
		return
	}
}