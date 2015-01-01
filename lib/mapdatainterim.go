package lib

import ()

// This type is required as the JSON map is not interface->interface.
// The map is then coherced to the common MapData type.
type MapDataInterim map[string]interface{}

func (m MapDataInterim) GetIndex(index interface{}) interface{} {
    return m[index.(string)]
}

func (m MapDataInterim) SetIndex(index interface{}, value interface{}) {
    m[index.(string)] = value
    // return m[index.(string)]
}

func (m MapDataInterim) Keys() []interface{} {
    keys := make([]interface{}, 0, len(m))
    for k, _ := range m {
        keys = append(keys, k)
    }

    return keys
}

func (m MapDataInterim) Length() int {
    return len(m)
}

func (m *MapDataTransition) ToMapData() *MapData {
    var converted_map = make(MapData)
    for key := range m.data {
        conv_key := interface{}(key)
        converted_map[conv_key] = m.data[key]
    }

    return &converted_map
}

type MapDataTransition struct {
    data MapDataInterim
}

func NewMapDataTransition(data MapDataInterim) *MapDataTransition {
    return &MapDataTransition{data}
}
