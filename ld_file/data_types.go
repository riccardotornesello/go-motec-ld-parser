package ld_file

type DataType struct {
	DataType       uint16
	DataTypeLength uint16
}

var (
	DataTypeFloat16 = DataType{0x07, 2}
	DataTypeFloat32 = DataType{0x07, 4}
	DataTypeInt16   = DataType{0x03, 2}
	DataTypeInt32   = DataType{0x05, 4}
)
