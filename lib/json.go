package lib

import (
    "encoding/json"
    "fmt"
    "log"
)

type JsonData struct {
    BaseData
}

func (j JsonData) String() string {
    return string(j.data)
}

func (j JsonData) ToMapData() *MapData {
    if DEBUG {
        log.Printf("Mapping this Json Data:\n%s\n", j.data)
    }

    var data_json = make(MapDataInterim)
    err := json.Unmarshal(j.data, &data_json)
    if err != nil {
        err_msg := fmt.Sprintf("Could not read data json, error: %s", err)
        log.Fatal(err_msg)
    }

    return ToMapData(data_json)
}

func NewJsonData(data []byte) *JsonData {
    return &JsonData{BaseData{data: data}}
}
