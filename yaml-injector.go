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
// The command can accept stdin JSON. JSON can be easily read inline and is
// the remote service/api response format of choice. In contrast, YAML is the
// written and recorded format of choice.
//
// Example usage: (may be out of date, see Makefile)
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
    // "encoding/json"
    // "fmt"
    // "reflect"
    "io/ioutil"
    "log"
    "os"

    "./lib"

    "code.google.com/p/go.crypto/ssh/terminal"
    "github.com/codegangsta/cli"
)

const (
    APPNAME = "Yaml Injector"
    VERSION = "0.2.0"
)

var (
    DEBUG = false
    TEST  = false
)

func SetDebug(value bool) {
    DEBUG = value
}

func SetTest(value bool) {
    TEST = value
}

func inject(dest_file lib.DataReader, data lib.DataReader, yaml_key string, data_key string) string {

    if DEBUG {
        log.Printf("Dest yaml file: %s", dest_file)
        log.Printf("Data: %s", data)
        log.Printf("Key to replace: %s", yaml_key)
        log.Printf("Data key to use: %s", data_key)
    }

    parsed_yaml := *dest_file.ToMapData()
    data_file := *data.ToMapData()

    // First check that the designated key exists within the dest yaml file
    if dest_value, ok := lib.GetMapValue(yaml_key, parsed_yaml); !ok {
        log.Fatalf("Could not find any value for dest key: (%s)", yaml_key)
    } else {
        if DEBUG {
            log.Printf("Dest data value to replace: %v", dest_value)
        }
    }

    // Check that the designated key exists within the data file
    // This does multiple passes over the data (validation, then modifier).
    if data_value, ok := lib.GetMapValue(data_key, data_file); !ok {
        log.Fatalf("Could not find any value for data key: (%s)", data_key)
    } else {
        if DEBUG {
            log.Printf("Data value to use in inject/replacement: %v", data_value)
        }

        run := func() {
            // Selector to find the value to replace
            selector := lib.NewSelector(yaml_key)

            // Function to replace the value
            modifier := func(in interface{}) interface{} {
                return data_value
            }

            lib.MapInSelect(&selector, parsed_yaml, modifier)
        }

        // At this point we can modify the map
        run()
    }

    // Single pass dest replacement with validation via selector
    // Not much of a difference than above. GetMapValue adds little overhead.
    // if data_value, ok := lib.GetMapValue(data_key, data_file); !ok {
    //     log.Fatalf("Could not find any value for data key: (%s)", data_key)
    // } else {
    //     if DEBUG {
    //         log.Printf("Data value to use in inject/replacement: %v", data_value)
    //     }

    //     run := func() {
    //         // Selector to find the value to replace
    //         selector := lib.NewSelector(yaml_key)

    //         // Function to replace the value
    //         modifier := func(in interface{}) interface{} {
    //             return data_value
    //         }

    //         _, dest_value := lib.MapInSelect(&selector, parsed_yaml, modifier)

    //         if !selector.MatchFound() {
    //             log.Fatalf("Could not find any value for dest key: (%s)", yaml_key)
    //         } else {
    //             if DEBUG {
    //                 log.Printf("Dest data value to replace: %v", dest_value)
    //             }
    //         }

    //     }

    //     // At this point we can modify the map
    //     run() // if ok
    // }

    return lib.WriteYaml(parsed_yaml)
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
                        var data lib.DataReader

                        if c.GlobalBool("debug") {
                            DEBUG = true
                            lib.SetDebug(DEBUG)
                        }

                        if c.GlobalBool("test") {
                            TEST = true
                            lib.SetTest(TEST)
                        }

                        // If stdin presents data, accept it as JSON and
                        // ignore any yaml data file flats.
                        // Otherwise use the flags given by the user.
                        if !terminal.IsTerminal(0) {
                            if DEBUG {
                                log.Println("stdin detected")
                            }

                            input, _ := ioutil.ReadAll(os.Stdin)
                            data = lib.NewJsonData(input)
                        } else {
                            if DEBUG {
                                log.Println("No stdin, proceeding with filemode")
                            }
                            data = lib.NewYamlData(lib.ReadYaml(c.GlobalString("using")))
                        }

                        dest_file := lib.NewYamlData(lib.ReadYaml(c.GlobalString("file")))
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
