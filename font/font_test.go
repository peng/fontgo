package font

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	dir, err := DataReader("../test/HanyiSentyCrayon.ttf")
	if err == nil {
		var dirStr []byte
		dirStr, err = json.Marshal(dir)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		fmt.Printf("%v", string(dirStr))
		return
	}

	t.Log(err)
	t.Fail()

}
