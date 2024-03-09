package ldfile

type LdFileVehicle struct {
	Id      [64]byte
	_       [128]byte
	Weight  uint32
	Type    [32]byte
	Comment [32]byte
}
