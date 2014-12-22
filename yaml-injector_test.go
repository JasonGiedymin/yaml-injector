package main

import (
    "strings"
    "testing"
)

var testData = []struct {
    yaml     string
    data     string
    yaml_key string
    data_key string
    expected string
}{
    {
        `---
  a: Easy!
  b:
    c: 2
    d: 
      - 3
      - 4
`,
        `---
 new_a: Easy2!
 new_b:
   new_c: 22
   new_d: 
     - 32
     - 42
`,
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

func TestYamlInject(t *testing.T) {
    for _, testData := range testData {
        // t.Logf("\niteration: %d \n data: %s", i, testData)

        if result := inject(
            []byte(testData.yaml),
            []byte(testData.data),
            testData.yaml_key,
            testData.data_key); result != strings.TrimLeft(testData.expected, "\n") {
            t.Errorf("\nExpected: \n[%s], \ngot: \n[%s]", testData.expected, result)
        }   // else all is well
    }
}
