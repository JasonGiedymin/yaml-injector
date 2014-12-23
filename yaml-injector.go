// yaml-injector.go
//
// Usage:
//   yaml-injector
//      --file <sample.yaml>
//        The file you wish to inject yaml into.
//
//      --data <data.yaml>
//        The yaml datafile with values you wish to use.
//
//      --key some.node OR node
//        This is the key you wish to replace or inject into.
//
//      run [-key]
//      The run command takes in an optional data key you wish to pluck from
//      the data file. If ommitted then the entire file will be used.
//

package main

import (
    // "bytes"
    "fmt"
    "io/ioutil"
    // "errors"
    "log"
    // "net/http"
    "os"
    // "text/template"
    // "time"

    "code.google.com/p/go.crypto/ssh/terminal"
    "github.com/codegangsta/cli"
    "gopkg.in/yaml.v2"
)

const (
    APPNAME = "Yaml Injector"
    VERSION = "0.1.0"
)

var (
    debug = false
)

type MapData map[interface{}]interface{}

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

func readYaml(filename string) []byte {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatalf("Could not read file %s, error: %s", filename, err)
    }

    return data
}

func writeYaml(yaml_data map[interface{}]interface{}) string {
    modified_yaml, err := yaml.Marshal(yaml_data)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    yaml_string := string(modified_yaml)

    fmt.Printf(yaml_string)
    return yaml_string
}

func inject(yaml_input DataReader, data DataReader, yaml_key string, data_key string) string {

    if debug {
        log.Printf("Input yaml: %s", yaml_input)
        log.Printf("Data: %s", data)
        log.Printf("Key to replace: %s", yaml_key)
        log.Printf("Data key to use: %s", data_key)
    }

    parsed_yaml := *yaml_input.Map()
    data_yaml := *data.Map()

    if parsed_yaml[yaml_key] != nil && data_yaml[data_key] != nil {
        if debug {
            log.Printf("yaml node: %s\n", parsed_yaml[yaml_key])
            log.Printf("data node: %s\n", data_yaml[data_key])
        }
        parsed_yaml[yaml_key] = data_yaml[data_key]
    } else if parsed_yaml[yaml_key] == nil {
        log.Fatalf("Could not find input key: %s", yaml_key)
    } else if data_yaml[data_key] == nil {
        log.Fatalf("Could not find data key: %s", data_key)
    }

    return writeYaml(parsed_yaml)
}

func main() {

    app := cli.NewApp()
    app.Name = APPNAME
    app.Usage = "yaml injector"

    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "file, f",
            Usage: "The file you wish to inject yaml into.",
        },
        cli.StringFlag{
            Name:  "using, u",
            Usage: "The yaml datafile with values you wish to use.",
        },
        cli.StringFlag{
            Name:  "key, k",
            Usage: "This is the key you wish to replace or inject into.",
        },
        cli.BoolFlag{
            Name:  "debug, d",
            Usage: "Debug output.",
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
                        if c.GlobalBool("debug") {
                            debug = true
                        }

                        if !terminal.IsTerminal(0) {
                            input, _ := ioutil.ReadAll(os.Stdin)
                            fmt.Print(string(input))
                        } else {
                            fmt.Println("no piped data")
                        }

                        yaml_input := readYaml(c.GlobalString("file"))
                        data := readYaml(c.GlobalString("using"))
                        data_key := c.GlobalString("key")
                        yaml_key := c.Args().First()
                        inject(NewYamlData(yaml_input), NewYamlData(data), yaml_key, data_key)
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
