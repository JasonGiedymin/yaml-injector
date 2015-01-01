package lib

import (
    "fmt"
    "reflect"
    "strings"
    "testing"
)

func TestMapFunc(t *testing.T) {
    append := func(in interface{}) interface{} {
        return fmt.Sprintf("%s world", in)
    }

    exec := func(fx MapFunc) {
        if fx("hello") != "hello world" {
            t.Error("Result of MapFunc incorrect.")
        }
    }

    exec(append)
}

func TestMapSimple(t *testing.T) {
    //func Map(acc *interface{}, list Mappable, modifier MapFunc) interface{} {
    var acc = make(MapData)

    var input_map = MapData{
        "hellokey":  "hello",
        "hello2key": "hello2",
        "hello3key": MapData{
            "hello4key": "hello4",
        },
    }

    fx := func(in interface{}) interface{} {
        return fmt.Sprintf("%s world", in)
    }

    result := Map(acc, input_map, fx)

    test_data := []struct {
        key      string
        expected interface{}
    }{
        {
            "hellokey",
            "hello world",
        },
        {
            "hello2key",
            "hello2 world",
        },
        {
            "hello3key",
            map[interface{}]interface{}{
                "hello4key": "hello4 world",
            },
        },
    }

    // t.Logf("Full result: %v\n", result)

    for _, test := range test_data {
        test_result := result.GetIndex(test.key)

        switch test_result.(type) {
        case Mappable:
            if reflect.DeepEqual(test_result, test.expected) {
                t.Errorf("\nExpected: %v\nGot: %v\n", test.expected, test_result)
                t.Errorf("Map() result: %v", result)
            }
        default:
            if test_result != test.expected {
                t.Errorf("\nExpected: %v\nGot: %v\n", test.expected, test_result)
                t.Errorf("Map() result: %v", result)
            }
        }
    }

}

func TestMapInPlace(t *testing.T) {
    var input_map = MapData{
        "hellokey":  "hello",
        "hello2key": "hello2",
        "hello3key": MapData{
            "hello4key": "hello4",
        },
    }

    fx := func(in interface{}) interface{} {
        return fmt.Sprintf("%s world", in)
    }

    result := MapIn(input_map, fx)

    test_data := []struct {
        key      string
        expected interface{}
    }{
        {
            "hellokey",
            "hello world",
        },
        {
            "hello2key",
            "hello2 world",
        },
        {
            "hello3key",
            map[interface{}]interface{}{
                "hello4key": "hello4 world",
            },
        },
    }

    // t.Logf("Full result: %v\n", result)

    for _, test := range test_data {
        test_result := result.GetIndex(test.key)

        switch test_result.(type) {
        case Mappable:
            if reflect.DeepEqual(test_result, test.expected) {
                t.Errorf("\nExpected: %v\nGot: %v\n", test.expected, test_result)
                t.Errorf("Map() result: %v", result)
            }
        default:
            if test_result != test.expected {
                t.Errorf("\nExpected: %v\nGot: %v\n", test.expected, test_result)
                t.Errorf("Map() result: %v", result)
            }
        }
    }
}

// In-place with selector test
func MapInPlaceSelector(run int, t *testing.T) {
    var input_map = MapData{
        "hellokey":  "hello",
        "hello2key": "hello2",
        "hello3key": MapData{
            "hello4key": "hello4",
        },
    }

    fx := func(in interface{}) interface{} {
        return fmt.Sprintf("%s world", in)
    }

    selector := NewSelector("hello3key.hello4key") // only modify this key

    test_data := []struct {
        key      string
        expected interface{}
    }{
        {
            "hello3key.hello4key",
            "hello4 world",
        },
        {
            "hellokey",
            "hello",
        },
    }

    MapInSelect(&selector, input_map, fx)

    if DEBUG {
        t.Logf("Run: [%d], In-place result: %v\n", run, input_map)
    }

    for _, test := range test_data {

        test_result, ok := GetValue(strings.Split(test.key, "."), input_map)
        if !ok {
            t.Errorf("Expected to obtain a value in map with key (%v), but got %v", test.key, test_result)
        }

        // fmt.Printf("GetKey: (%v) -> %v, %v\n", test.key, test_result, ok)

        switch test_result.(type) {
        case Mappable:
            if reflect.DeepEqual(test_result, test.expected) {
                t.Errorf("\nExpected: %v\nGot: %v\n", test.expected, test_result)
                t.Errorf("Map() result: %v", input_map)
            }
        default:
            if test_result != test.expected {
                t.Errorf("\nExpected: %v\nGot: %v\n", test.expected, test_result)
                t.Errorf("Map() result: %v", input_map)
            }
        }

    }
}

func TestMapInPlaceSelector(t *testing.T) {
    runs := 100

    if DEBUG {
        t.Logf("----------------- In-place testing runs [%d] -----------", runs)
    }

    for i := 0; i <= runs; i++ {
        MapInPlaceSelector(i, t)
    }

    if DEBUG {
        t.Logf("----------------- In-place testing runs complete -----------")
    }
}
