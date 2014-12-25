// yaml-injector.go
//
// Usage:
//
// yaml-injector
//
// COMMANDS:
//    inject, i  Runs the yaml injector replace command
//      [into]   optional data key to inject into
//    help, h  Shows a list of commands or help for one command
//
// GLOBAL OPTIONS:
//    --file, -f     the file you wish to inject yaml into.
//    --using, -u    the yaml datafile with values you wish to use.
//    --key, -k      this is the key you wish to replace or inject into.
//    --debug, -d    debug output.
//    --test, -t     test run the command by printing the output to stdout
//    --help, -h     show help
//    --version, -v  print the version
//
// Stdin (JSON)
// The command can accept stdin JSON. JSON does away with new lines
// and can be read inline.
//
// Examples:
//   - Data file usage:
//     go run yaml-injector.go
//       --debug --file test/input.yaml \
//       --using test/data.yaml \
//       --key "a" \
//       inject into "a"
//
//   - Data file usage but using test run to see expected output on stdout:
//     go run yaml-injector.go
//       --test
//       --debug --file test/input.yaml \
//       --using test/data.yaml \
//       --key "a" \
//       inject into "a"
//
//   - Stdin JSON usage:
//     echo '{"a":1}' | go run yaml-injector.go --debug \
//       --file test/input.yaml \
//       --using test/data.yaml \
//       --key "a" inject into "a"
//

package main

import (
    // "bytes"
    // "errors"
    // "reflect"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"

    "code.google.com/p/go.crypto/ssh/terminal"
    "github.com/codegangsta/cli"
    "gopkg.in/yaml.v2"
)

const (
    APPNAME = "Yaml Injector"
    VERSION = "0.2.0"
)

var (
    debug = true
    test  = false
)

type MapData map[interface{}]interface{}
type MapDataInterim map[string]interface{}

type MapDataTransition struct {
    data MapDataInterim
}

func (m *MapDataTransition) ToMapData() *MapData {
    var converted_map = make(MapData)
    for key := range m.data {
        conv_key := interface{}(key)
        converted_map[conv_key] = m.data[key]
    }

    return &converted_map
}

func NewMapDataTransition(data MapDataInterim) *MapDataTransition {
    return &MapDataTransition{data}
}

func ToMapData(data MapDataInterim) *MapData {
    return NewMapDataTransition(data).ToMapData()
}

type DataReader interface {
    String() string
    Data() []byte
    Map() *MapData
}

type BaseData struct {
    data []byte
}

func (b BaseData) Data() []byte {
    return b.data
}

type JsonData struct {
    BaseData
}

func (j JsonData) String() string {
    return string(j.data)
}

func (j JsonData) Map() *MapData {
    if debug {
        log.Printf("Mapping this Json Data:\n%s\n", j.data)
    }

    var data_json = make(MapDataInterim)
    err := json.Unmarshal(j.data, &data_json)
    if err != nil {
        err_msg := fmt.Sprintf("Could not read data json, error: %s", err)
        log.Fatal(err_msg)
    }

    return NewMapDataTransition(data_json).ToMapData()
}

func NewJsonData(data []byte) *JsonData {
    return &JsonData{BaseData{data: data}}
}

type YamlData struct {
    BaseData
}

func (y YamlData) String() string {
    return string(y.data)
}

func (y YamlData) Map() *MapData {
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

func find(tokens []string) string {
    fmt.Printf(".")

    switch len(tokens) {
    case 0:
        return ""
    case 1:
        return tokens[0]
    default:
        return find(tokens[1:])
    }
}

// Keep it really simple by using dot notation.
// Retain everything else as it will just fail anyway.
func GetKey(tokens []string, data interface{}) (interface{}, bool) {
    // fmt.Printf("\n>>%v, %v\n\n", key, reflect.TypeOf(key))
    // fmt.Printf("\n>>%v, %v\n\n", data, reflect.TypeOf(data))
    // fmt.Printf("\n>>%v, %v\n\n", value, reflect.TypeOf(value))

    var next = func(format string) (interface{}, bool) {
        token := tokens[0]
        if debug {
            fmt.Printf(format, token)
        }

        switch data.(type) {
        case map[string]interface{}:
            value := data.(map[string]interface{})[token]
            // fmt.Printf("\n==> %s\n", value)
            return GetKey(tokens[1:], value)
        default:
            return GetKey(tokens[1:], data.(MapData)[token])
        }
    }

    if len(tokens) > 1 {
        return next("(%v).")
    } else if len(tokens) == 1 {
        return next("(%v)")
    } else {
        if debug {
            fmt.Printf(" -> [%v]\n", data)
        }

        return data, true
    }

    return nil, false
}

func readYaml(filename string) []byte {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatalf("Could not read file %s, error: %s", filename, err)
    }

    return data
}

func writeYaml(yaml_data MapData) string {
    modified_yaml, err := yaml.Marshal(yaml_data)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    yaml_string := string(modified_yaml)

    if test {
        fmt.Printf("%s\n", yaml_string)
    }
    return yaml_string
}

func inject(dest_file DataReader, data DataReader, yaml_key string, data_key string) string {

    if debug {
        log.Printf("Dest yaml file: %s", dest_file)
        log.Printf("Data: %s", data)
        log.Printf("Key to replace: %s", yaml_key)
        log.Printf("Data key to use: %s", data_key)
    }

    parsed_yaml := *dest_file.Map()
    data_yaml := *data.Map()

    if parsed_yaml[yaml_key] != "" && data_yaml[data_key] != "" {
        if debug {
            log.Printf("yaml node: %s\n", parsed_yaml[yaml_key])
            log.Printf("data node: %s\n", data_yaml[data_key])
        }
        parsed_yaml[yaml_key] = data_yaml[data_key]
    } else if parsed_yaml[yaml_key] == "" {
        log.Fatalf("Could not find input key: %s", yaml_key)
    } else if data_yaml[data_key] == "" {
        log.Fatalf("Could not find data key: %s", data_key)
    }

    return writeYaml(parsed_yaml)
}

func main() {

    app := cli.NewApp()
    app.Name = APPNAME
    app.Usage = "yaml injector"
    app.Version = VERSION

    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "file, f",
            Usage: "the file you wish to inject yaml into.",
        },
        cli.StringFlag{
            Name:  "using, u",
            Usage: "the yaml datafile with values you wish to use.",
        },
        cli.StringFlag{
            Name:  "key, k",
            Usage: "this is the key you wish to replace or inject into.",
        },
        cli.BoolFlag{
            Name:  "debug, d",
            Usage: "debug output.",
        },
        cli.BoolFlag{
            Name: "test, t",
            Usage: "test run the command by printing the output " +
                "to stdout, also useful for pipeing out data without " +
                "having to write out to a file",
        },
    }

    app.Commands = []cli.Command{
        {
            Name:      "inject",
            ShortName: "i",
            Usage:     "Runs the yaml injector replace command",
            Subcommands: []cli.Command{
                {
                    Name:  "into",
                    Usage: "optional data key to inject into",
                    Action: func(c *cli.Context) {
                        var data DataReader

                        if c.GlobalBool("debug") {
                            debug = true
                        }

                        if c.GlobalBool("test") {
                            test = true
                        }

                        // If stdin presents data, accept it as JSON and
                        // ignore any yaml data file flats.
                        // Otherwise use the flags given by the user.
                        if !terminal.IsTerminal(0) {
                            input, _ := ioutil.ReadAll(os.Stdin)
                            data = NewJsonData(input)
                        } else {
                            data = NewYamlData(readYaml(c.GlobalString("using")))
                        }

                        dest_file := NewYamlData(readYaml(c.GlobalString("file")))
                        data_key := c.GlobalString("key")
                        yaml_key := c.Args().First()
                        inject(dest_file, data, yaml_key, data_key)
                    },
                },
            },
        },
    }

    app.Action = func(c *cli.Context) {
        cli.ShowAppHelp(c)
        return
    }

    app.Run(os.Args)
}
