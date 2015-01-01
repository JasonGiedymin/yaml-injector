package lib

import (
    "log"
)

// As we cannot work with the pointers within a map inherently.
// Maps vary and shift in size, modifing any pointers you have
// references to. We must manually convert all values to pointers.
// Making a map of pointers rather than values.
type MapDataPointers map[*interface{}]*interface{}

func (mp *MapDataPointers) Print() {
    log.Println("--- Map Listing ---")
    for key := range *mp {
        log.Printf("[%v]:(%v)", *key, *(*mp)[key])
    }
    log.Println("-------------------")
}
