package xtime

import (
	"encoding/json"
	"testing"
	"time"
)

type Order struct {
	OrderNumber int    `json:"orderNumber"`
	CreateTime  Time   `json:"createTime"`
	Other       string `json:"-"`
}

func Test_MarshallJson(t *testing.T) {
	data, err := json.Marshal(&Order{
		OrderNumber: 123456,
		CreateTime:  Time{time.Now()},
		Other:       "test",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(data))
}

func Test_UmMarshallJson(t *testing.T) {
	jsonData := `{"orderNumber":789456,"createTime":"2020-10-17 17:58:57"}`
	var o Order
	if err := json.Unmarshal([]byte(jsonData), &o); err != nil {
		t.Error(err)
		return
	}
	t.Log(o)
}
