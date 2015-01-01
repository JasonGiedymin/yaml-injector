package main

import (
    "./lib"

    "testing"
)

var bench_test_assets = struct {
    yaml     string
    data     string
    yaml_key string
    data_key string
    expected string
}{
    TEST_YAML,
    TEST_DATA,
    "a",
    "new_a",
    `a: Easy2!
b:
  c: 2
  d:
  - 3
  - 4
`}

var bench_dest_yaml = lib.NewYamlData([]byte(bench_test_assets.yaml))
var bench_data_yaml = lib.NewYamlData([]byte(bench_test_assets.data))

func benchmarkInject() {
    inject(
        bench_dest_yaml,
        bench_data_yaml,
        bench_test_assets.yaml_key,
        bench_test_assets.data_key)
}

func benchmarkInjectHarness(b *testing.B) {
    for n := 0; n < b.N; n++ {
        benchmarkInject()
    }
}

func BenchmarkInject(b *testing.B) {
    SetDebug(false)
    b.ResetTimer()
    benchmarkInjectHarness(b)
}
