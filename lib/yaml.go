package lib

import (
    "fmt"
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v2"
)

type YamlData struct {
    BaseData
}

func (y YamlData) String() string {
    return string(y.data)
}

func (y YamlData) ToMapData() *MapData {
    var data_yaml = make(MapData)
    err := yaml.Unmarshal(y.data, &data_yaml)
    if err != nil {
        err_msg := fmt.Sprintf("Could not read data yaml, error: %s", err)
        log.Fatal(err_msg)
    }

    return &data_yaml
}

func NewYamlData(data []byte) *YamlData {
    return &YamlData{BaseData{data: data}}
}

func ReadYaml(filename string) []byte {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatalf("Could not read file %s, error: %s", filename, err)
    }

    return data
}

func WriteYaml(yaml_data MapData) string {
    modified_yaml, err := yaml.Marshal(yaml_data)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    yaml_string := string(modified_yaml)

    if TEST {
        fmt.Printf("%s\n", yaml_string)
    }
    return yaml_string
}
