package lib

import (
    // "fmt"
    "reflect"
    "strings"
    "testing"
)

func TestSelector(t *testing.T) {
    // SetDebug(true)

    path := "a.b.c"
    result := Selector{
        path:       path,
        tokens:     strings.Split(path, "."),
        tokens_acc: []string{},
        match:      false,
    }

    constructor_result := NewSelector(path)

    if !reflect.DeepEqual(result, constructor_result) {
        t.Errorf("NewSelector() failed, expected constructor result:\n%v\nto match:\n%v\n", constructor_result, result)
    }

    // Mutations follow
    result.tokens_acc = []string{"a", "b", "c"}

    if !result.Match() {
        t.Errorf("Match() failed, expected tokens: %v and tokens_acc: %v, to match\n", result.tokens, result.tokens_acc)
    }

    // Test for no match
    result.tokens_acc = []string{"a", "b"}
    if !result.Match() && result.MatchFound() { //execute match, and test matchFound
        t.Errorf("MatchFound() was expected to be false, instead found true.")
    }

}

func TestSelectorPushPop(t *testing.T) {
    path := "a.b.c"
    selector := Selector{
        path:       path,
        tokens:     strings.Split(path, "."),
        tokens_acc: []string{},
        match:      false,
    }

    expected := []string{"a", "b", "c"}

    selector.Push("a")
    selector.Push("b")
    count := selector.Push("c")

    if count != 3 {
        t.Logf("Push() failed, expected count of %d after pushes, got:%d", 3, count)
    }

    if !reflect.DeepEqual(selector.tokens_acc, expected) {
        t.Logf("Push() method failed, expected: %v, got: %v\n", expected, selector.tokens_acc)
    }

    count = selector.PopT()
    expected = []string{"a", "b"}

    if count != 2 {
        t.Logf("Push() failed, expected count of %d after pushes, got:%d", 2, count)
    }

    if !reflect.DeepEqual(selector.tokens_acc, expected) {
        t.Logf("Push() method failed, expected: %v, got: %v\n", expected, selector.tokens_acc)
    }
}
