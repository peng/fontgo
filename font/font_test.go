package font

import (
	"encoding/json"
	"fmt"
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
	File        *File        `json:"file"`
	OffsetTable *OffsetTable `json:"offsetTable"`
	Head        *Head        `json:"head"`
}

func GetRemoteJson() (data *RespData, err error) {
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
		fileByte  []byte
		err       error
		standData *RespData
	)

	fileByte, err = os.ReadFile("../test/HanyiSentyCrayon.ttf")

	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	standData, err = GetRemoteJson()
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	ToffsetTable(t, fileByte, standData)
}

func ToffsetTable(t *testing.T, fileByte []byte, standData *RespData) {
	// test offsetTable
	var (
		err error
	)

	offsetTable := GetOffsetTable(fileByte)

	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	if *offsetTable != *standData.OffsetTable {
		fmt.Println("standData offsetTable", standData.OffsetTable)
		fmt.Println("source offsetTable", offsetTable)
		t.Log("offsetTable error")
		t.Fail()
		return
	}
}
