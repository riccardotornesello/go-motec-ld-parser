package ld_file

type LdFileChannelMeta struct {
	PreviousMetaPointer uint32
	NextMetaPointer     uint32
	DataPointer         uint32
	DataLength          uint32
	ChannelId           uint16 // 0x2EE1 + n
	DataType            uint16 // 0x07 for float16/32, 0x05 for int32, 0x03 for int16
	DataTypeLength      uint16
	Frequency           uint16
	Shift               int16
	Mul                 int16
	Scale               int16
	DecPlaces           int16
	Name                [32]byte
	ShortName           [8]byte
	Unit                [12]byte
	_                   [40]byte // (40 bytes for ACC, 32 bytes for acti)
}
