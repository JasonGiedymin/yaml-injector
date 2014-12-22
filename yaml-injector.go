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

    "github.com/codegangsta/cli"
    "gopkg.in/yaml.v2"
)

const (
    APPNAME = "Amuxbit Yaml Injector"
    VERSION = "1.0.0.0"
)

var (
    debug = false
)

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

func inject(yaml_input []byte, data []byte, yaml_key string, data_key string) string {
    var parsed_yaml = make(map[interface{}]interface{})
    var data_yaml = make(map[interface{}]interface{})

    if debug {
        log.Printf("Input yaml: %s", string(yaml_input))
        log.Printf("Data: %s", string(data))
        log.Printf("Key to replace: %s", yaml_key)
        log.Printf("Data key to use: %s", data_key)
    }

    err := yaml.Unmarshal(yaml_input, &parsed_yaml)
    if err != nil {
        err_msg := fmt.Sprintf("Could not read yaml input, error: %s", err)
        log.Fatal(err_msg)
    }

    data_err := yaml.Unmarshal(data, &data_yaml)
    if err != nil {
        err_msg := fmt.Sprintf("Could not read data yaml, error: %s", err)
        log.Fatal(err_msg)
    }

    if parsed_yaml[yaml_key] != nil && data_yaml[data_key] != nil {
        if debug {
            log.Printf("yaml node: %s\n", parsed_yaml[yaml_key])
            log.Printf("data node: %s\n", data_yaml[data_key])
        }
        // var t string
        // t = "b"
        // log.Printf(parsed_yaml)
        parsed_yaml[yaml_key] = data_yaml[data_key]
    } else if err != nil {
        err_msg := fmt.Sprintf("Key: [%s] not found in yaml input.", yaml_key)
        log.Fatal(err_msg)
    } else if data_err != nil {
        err_msg := fmt.Sprintf("Key: [%s] not found in data yaml input.", data_key)
        log.Fatal(err_msg)
    } else {
        if parsed_yaml[yaml_key] == nil {
            log.Fatalf("Could not find input key: %s", yaml_key)
        }

        if data_yaml[data_key] == nil {
            log.Fatalf("Could not find data key: %s", data_key)
        }
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

                        yaml_input := readYaml(c.GlobalString("file"))
                        data := readYaml(c.GlobalString("using"))
                        data_key := c.GlobalString("key")
                        yaml_key := c.Args().First()
                        inject(yaml_input, data, yaml_key, data_key)
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
