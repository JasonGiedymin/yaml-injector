package lib

type DataReader interface {
    String() string
    Data() []byte
    ToMapData() *MapData
}

type BaseData struct {
    data []byte
}

func (b BaseData) Data() []byte {
    return b.data
}
