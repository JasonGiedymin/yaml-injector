package lib

import (
    "fmt"
    "strings"
    "testing"
)

var input_map = MapData{
    "hello":  "hello",
    "hello2": "hello2",
    "hello3": MapData{
        "hello4": "hello4",
    },
}

var deepest_key = "hello3.hello4"

var fx = func(in interface{}) interface{} {
    return fmt.Sprintf("%s world", in)
}

var selector = NewSelector("hello3.hello4") // only modify this key

var test_result, original_value = MapInSelect(&selector, input_map, fx)

func benchmarkMapInPlaceSelector() {
    MapInSelect(&selector, input_map, fx)
}

func benchmarkMapInPlaceSelectorHarness(b *testing.B) {
    for n := 0; n < b.N; n++ {
        benchmarkMapInPlaceSelector()
    }
}

func BenchmarkMapInPlaceSelector(b *testing.B) {
    SetDebug(false)
    b.ResetTimer()
    benchmarkMapInPlaceSelectorHarness(b)
}

// In-place with selector test
func benchmarkGetValue() {
    GetValue(strings.Split(deepest_key, "."), test_result)
}

func benchmarkGetKeyHarness(b *testing.B) {
    for n := 0; n < b.N; n++ {
        benchmarkGetValue()
    }
}

func BenchmarkGetValue(b *testing.B) {
    SetDebug(false)
    b.ResetTimer()
    benchmarkGetKeyHarness(b)
}
