package ld_file

type LdFileEvent struct {
	Name         [64]byte
	Session      [64]byte
	Comment      [1024]byte
	VenuePointer uint16
}
