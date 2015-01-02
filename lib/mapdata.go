package lib

import (
    "log"
    "reflect"
)

type MapData map[interface{}]interface{}

func (m MapData) GetIndex(index interface{}) interface{} {
    return m[index]
}

func (m MapData) SetIndex(index interface{}, value interface{}) {
    m[index] = value
    // log.Printf("SetIndex(): key= %v, entry= %v, full map: %v\n", index, m[index], m)
}

func (m MapData) Keys() []interface{} {
    keys := make([]interface{}, 0, len(m))
    for k, _ := range m {
        //log.Printf("************* (%v)->(%v) \n", k, v)
        keys = append(keys, k)
    }

    return keys
}

func (m MapData) Length() int {
    return len(m)
}

func (m *MapData) ToMapDataPointers() MapDataPointers {
    var converted_map = make(MapDataPointers)

    for key := range *m {
        new_key := key     // key reference, key is in-place
        value := (*m)[key] // value reference

        if DEBUG {
            log.Printf(">> entry: %v (%v) [%v], value: %v (%v) [%v]", new_key, &new_key, reflect.TypeOf(new_key), value, &value, reflect.TypeOf(value))
        }

        converted_map[&new_key] = &value
    }

    return converted_map
}
