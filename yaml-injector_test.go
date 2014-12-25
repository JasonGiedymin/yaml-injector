package main

import (
    "reflect"
    "strings"
    "testing"
)

var TEST_YAML = `---
  a: Easy!
  b:
    c: 2
    d: 
      - 3
      - 4
`

var TEST_DATA = `---
  new_a: Easy2!
  new_b:
    new_c: 22
    new_d: 
      - 32
      - 42
`

var TEST_YAML_PARSER = `---
  new_a: 2
  new_b:
    new_c: 22
    new_d: 
      - 32
      - 42
  new_e: two
  new_f:
    new_f2:
      new_f3: f3_pass!
`

var TEST_DATA_JSON = `
{
  "new_a": "Easy3!",
  "new_b": {
    "new_c" : 22
  },
  "new_e": "two",
  "new_f": {
    "new_f2": {
      "new_f3": "f3_pass!"
    }
  }
}
`

var test_assets = []struct {
    yaml     string
    data     string
    yaml_key string
    data_key string
    expected string
}{
    {
        TEST_YAML,
        TEST_DATA,
        "a",
        "new_a",
        `
a: Easy2!
b:
  c: 2
  d:
  - 3
  - 4
`,
    },
}

var test_assets_json = []struct {
    dest     string
    data     string
    yaml_key string
    data_key string
    expected string
}{
    {
        TEST_YAML,
        TEST_DATA_JSON,
        "a",
        "new_a",
        `
a: Easy3!
b:
  c: 2
  d:
  - 3
  - 4
`,
    },
}

var test_assets_parser_yaml = []struct {
    key      string
    data     string
    expected interface{}
}{
    {
        "new_b.new_c",
        TEST_YAML_PARSER,
        22,
    },
    {
        "new_e",
        TEST_YAML_PARSER,
        "two",
    },
    {
        "new_f.new_f2.new_f3",
        TEST_YAML_PARSER,
        "f3_pass!",
    },
}

var test_assets_parser_json = []struct {
    key      string
    data     string
    expected interface{}
}{
    {
        "new_b.new_c",
        TEST_DATA_JSON,
        float64(22),
    },
    {
        "new_e",
        TEST_DATA_JSON,
        "two",
    },
    {
        "new_f.new_f2.new_f3",
        TEST_DATA_JSON,
        "f3_pass!",
    },
}

func TestYamlInject(t *testing.T) {
    for _, test_assets := range test_assets {
        // t.Logf("\niteration: %d \n data: %s", i, test_assets)

        if result := inject(
            NewYamlData([]byte(test_assets.yaml)),
            NewYamlData([]byte(test_assets.data)),
            test_assets.yaml_key,
            test_assets.data_key); result != strings.TrimLeft(test_assets.expected, "\n") {
            t.Errorf("\nExpected: \n[%s], \ngot: \n[%s]", test_assets.expected, result)
        }   // else all is well
    }
}

func TestJsonInject(t *testing.T) {
    for _, asset := range test_assets_json {
        // t.Logf("\niteration: %d \n data: %s", i, asset)

        if result := inject(
            NewYamlData([]byte(asset.dest)),
            NewJsonData([]byte(asset.data)),
            asset.yaml_key,
            asset.data_key); result != strings.TrimLeft(asset.expected, "\n") {
            t.Errorf("\nExpected: \n[%s], \ngot: \n[%s]", asset.expected, result)
        }   // else all is well
    }
}

func TestGetKey(t *testing.T) {
    // Parse basic dot notation on yaml
    for _, asset := range test_assets_parser_yaml {
        map_data := *NewYamlData([]byte(asset.data)).Map()
        if result, ok := GetKey(strings.Split(asset.key, "."), map_data); ok {
            if result != asset.expected {
                t.Errorf("\nExpected: \n[%v], \ngot: \n[%v]", asset.expected, result)
                t.Errorf("Map data >> %v", map_data)
            }
        }
    }

    // Parse basic dot notation on json
    for _, asset := range test_assets_parser_json {
        map_data := *NewJsonData([]byte(asset.data)).Map()
        if result, ok := GetKey(strings.Split(asset.key, "."), map_data); ok == true {
            if result != asset.expected {
                t.Errorf("\nExpected: \n[%v](%s), \ngot: \n[%v](%s)", asset.expected, reflect.TypeOf(asset.expected), result, reflect.TypeOf(result))
                t.Errorf("Map data >> %v", map_data)
            }   //else {
            //     t.Logf("\nExpected: \n[%v](%s), \ngot: \n[%v](%s)", asset.expected, reflect.TypeOf(asset.expected), result, reflect.TypeOf(result))
            // }
        }
    }
}
