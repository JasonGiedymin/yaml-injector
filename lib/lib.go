package lib

import (
    "fmt"
    "log"
    "strings"
    // "reflect"
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

func Debug() bool {
    return DEBUG
}

func Test() bool {
    return TEST
}

type Mappable interface {
    GetIndex(index interface{}) interface{}
    SetIndex(index interface{}, value interface{})
    Keys() []interface{}
    Length() int
}

type MapFunc func(interface{}) interface{}

// Copy modify all, treat input as Immutable
func Map(acc Mappable, input_map interface{}, fx MapFunc) Mappable {
    data_map := input_map.(Mappable)
    for _, key := range data_map.Keys() {
        input_value := data_map.GetIndex(key)
        switch input_value.(type) {
        case Mappable:
            result := Map(make(MapData), input_value, fx)
            acc.SetIndex(key, result)
        default:
            value := fx(input_value)
            acc.SetIndex(key, value)
        }
    }

    return acc
}

// In-place modify all
func MapIn(input_map Mappable, fx MapFunc) Mappable {
    data_map := input_map.(Mappable)

    for _, key := range data_map.Keys() {
        input_value := data_map.GetIndex(key)
        switch input_value.(type) {
        case Mappable:
            value := Map(make(MapData), input_value, fx)
            input_map.SetIndex(key, value)
        default:
            value := fx(input_value)
            input_map.SetIndex(key, value)
        }
    }

    return input_map
}

// In-place modify using a Selector which acts as a filter. Rather than
// monadically doing `data.filter(keys_i_want).map(fx)` I try to mutate
// and search in one pass. Here I treat fx as a modifier instead.
func MapInSelect(selector *Selector, input_map Mappable, fx MapFunc) (Mappable, interface{}) {
    var new_value, original_value interface{}
    data_map := input_map.(Mappable)
    keys := data_map.Keys()

    if DEBUG {
        log.Printf("    Fullmap: %v\n", input_map)
        log.Printf("    All Keys: %v\n", keys)
    }

    for _, key := range keys {
        selector.Push(key.(string))

        input_value := data_map.GetIndex(key)

        switch input_value.(type) {
        case Mappable:
            new_value, original_value = MapInSelect(selector, input_value.(Mappable), fx)
            if selector.Match() {
                input_map.SetIndex(key, new_value)
                return input_map, original_value
            }
        default:
            if selector.Match() {
                original_value = fx(input_value)
                input_map.SetIndex(key, fx(input_value))
                return input_map, original_value
            }

            selector.PopT()
        }
    }

    if DEBUG {
        selector.Print()
    }

    selector.PopT()
    return input_map, nil
}

func ToMapData(data MapDataInterim) *MapData {
    return NewMapDataTransition(data).ToMapData()
}

func find(tokens []string) string {
    if DEBUG {
        fmt.Printf(".")
    }

    switch len(tokens) {
    case 0:
        return ""
    case 1:
        return tokens[0]
    default:
        return find(tokens[1:])
    }
}

func GetMapValue(key string, data interface{}) (interface{}, bool) {
    return GetValue(strings.Split(key, "."), data)
}

// Keep it really simple by using dot notation.
// Retain everything else as it will just fail anyway.
func GetValue(tokens []string, data interface{}) (interface{}, bool) {
    // fmt.Printf("\n>>%v, %v\n\n", key, reflect.TypeOf(key))
    // fmt.Printf("\n>>%v, %v\n\n", data, reflect.TypeOf(data))
    // fmt.Printf("\n>>%v, %v\n\n", value, reflect.TypeOf(value))

    var next = func(format string) (interface{}, bool) {
        token := tokens[0]
        if DEBUG {
            fmt.Printf(format, token)
        }

        switch data.(type) {
        case map[string]interface{}:
            value := data.(map[string]interface{})[token]
            // data.(map[string]interface{})[token] = "crap"
            // fmt.Printf("\n**> %v\n", data)
            return GetValue(tokens[1:], value)
        default:
            new_data := data.(MapData)[token]
            return GetValue(tokens[1:], new_data)
        }
    }

    if len(tokens) > 1 {
        return next("(%v).")
    } else if len(tokens) == 1 {
        return next("(%v)")
    } else {
        if DEBUG {
            fmt.Printf(" -> [%v]\n", data)
        }
        // data = "crap"
        return data, true
    }

    return nil, false
}
