package lib

import (
    "log"
    "reflect"
    "strings"
)

type Selector struct {
    path       string
    tokens     []string
    tokens_acc []string
    match      bool
}

func (s *Selector) Match() bool {
    if DEBUG {
        log.Printf("Comparing tokens: %v == tokens_acc: %v", s.tokens, s.tokens_acc)
    }

    s.match = reflect.DeepEqual(s.tokens, s.tokens_acc)
    return s.match
}

func (s *Selector) MatchFound() bool {
    return s.match
}

// needs ordering with push
func (s *Selector) Push(token string) int {
    if DEBUG {
        log.Printf("Pushing: [%v]", token)
    }

    s.tokens_acc = append((*s).tokens_acc, token)
    return len((*s).tokens_acc)
}

// Pops from tail
func (s *Selector) PopT() int {
    remove := len(s.tokens_acc) - 1

    if DEBUG {
        log.Printf("Popping: [%v]", s.tokens_acc[remove])
    }

    s.tokens_acc = s.tokens_acc[:remove]
    return len(s.tokens_acc)
}

func (s *Selector) Print() {
    log.Printf("Selector tokens accumulated: %v\n", s.tokens_acc)
}

func NewSelector(path string) Selector {
    return Selector{
        path:       path,
        tokens:     strings.Split(path, "."),
        tokens_acc: []string{},
        match:      false,
    }
}
