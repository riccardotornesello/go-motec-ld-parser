package ldfile

type LdFileHead struct {
	LDMarker            uint32 // 0x40
	_                   [4]byte
	ChannelsMetaPointer uint32
	ChannelsDataPointer uint32
	_                   [20]byte
	EventPointer        uint32
	_                   [24]byte
	Unknown1            uint16  // 1
	Unknown2            uint16  // 0x4240
	Unknown3            uint16  // 0xF
	DeviceSerial        uint32  // 0x1F44
	DeviceType          [8]byte // "ADL"
	DeviceVersion       uint16  // 420
	Unknown4            uint16  // 0xADB0
	ChannelsCount       uint32
	_                   [4]byte
	Date                [16]byte // "dd/MM/yyyy"
	_                   [16]byte
	Time                [16]byte // "HH:mm:ss"
	_                   [16]byte
	Driver              [64]byte
	Vehicle             [64]byte
	_                   [64]byte
	Venue               [64]byte
	_                   [64]byte
	_                   [1024]byte
	EnableProLogging    uint32 // 0xC81A4
	_                   [66]byte
	ShortComment        [64]byte
	_                   [126]byte
}
