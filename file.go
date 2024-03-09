package motec_ld_parser

import (
	"bytes"
	"encoding/binary"
	"os"
	"time"

	"riccardotornesello.it/motec-ld-parser/ldfile"
)

/*
	|---------------|
	|	HEAD		|
	|---------------| <- EVENT_POINTER
	|	EVENT		|
	|---------------| <- VENUE_POINTER
	|	VENUE		|
	|---------------| <- VEHICLE_POINTER
	|	VEHICLE		|
	|---------------| <- CHANNELS_META_POINTER
	|	CHANNEL H	|
	|	CHANNEL H	|
	|---------------| <- CHANNELS_DATA_POINTER
	|	CHANNEL D	|
	|	CHANNEL D	|
	|---------------|
*/

type File struct {
	Time         time.Time
	Driver       string
	Vehicle      string
	Venue        string
	ShortComment string

	EventName    string
	EventSession string
	EventComment string

	VehicleId      string
	VehicleWeight  uint32
	VehicleType    string
	VehicleComment string

	Channels []interface{}
}

type Channel[T float32 | int16 | int32] struct {
	Frequency uint16
	Name      string
	ShortName string
	Unit      string
	Data      *[]T
}

func (f *File) Write(fd *os.File) {
	// Calculate pointers
	headerSize := uintptr(binary.Size(ldfile.LdFileHead{}))
	eventSize := uintptr(binary.Size(ldfile.LdFileEvent{}))
	venueSize := uintptr(binary.Size(ldfile.LdFileVenue{}))
	vehicleSize := uintptr(binary.Size(ldfile.LdFileVehicle{}))
	channelMetaSize := uintptr(binary.Size(ldfile.LdFileChannelMeta{}))

	eventPointer := headerSize
	venuePointer := eventPointer + eventSize
	vehiclePointer := venuePointer + venueSize
	channelsMetaPointer := vehiclePointer + vehicleSize
	channelsDataPointer := channelsMetaPointer + channelMetaSize*uintptr(len(f.Channels))

	// Create the file header
	head := ldfile.LdFileHead{
		LDMarker:         0x40,
		Unknown1:         1,
		Unknown2:         0x4240,
		Unknown3:         0xF,
		Unknown4:         0xADB0,
		DeviceSerial:     0x1F44,
		DeviceType:       [8]byte{'A', 'D', 'L', 0, 0, 0, 0, 0},
		DeviceVersion:    420,
		EnableProLogging: 0xC81A4,
		ChannelsCount:    uint32(len(f.Channels)),

		EventPointer:        uint32(eventPointer),
		ChannelsMetaPointer: uint32(channelsMetaPointer),
		ChannelsDataPointer: uint32(channelsDataPointer),
	}

	date := f.Time.Format("02/01/2006")
	hour := f.Time.Format("15:04:05")
	copy(head.Date[:], date)
	copy(head.Time[:], hour)

	copy(head.Driver[:], f.Driver)
	copy(head.Vehicle[:], f.Vehicle)
	copy(head.Venue[:], f.Venue)
	copy(head.ShortComment[:], f.ShortComment)

	// Create the Event
	event := ldfile.LdFileEvent{
		VenuePointer: uint16(venuePointer),
	}

	copy(event.Name[:], f.EventName)
	copy(event.Session[:], f.EventSession)
	copy(event.Comment[:], f.EventComment)

	// Create the Venue
	venue := ldfile.LdFileVenue{
		VehiclePointer: uint16(vehiclePointer),
	}

	copy(venue.Name[:], f.Venue)

	// Create the Vehicle
	vehicle := ldfile.LdFileVehicle{
		Weight: f.VehicleWeight,
	}

	copy(vehicle.Id[:], f.VehicleId)
	copy(vehicle.Type[:], f.VehicleType)
	copy(vehicle.Comment[:], f.VehicleComment)

	// Write to file
	binary.Write(fd, binary.LittleEndian, head)

	fd.Seek(int64(eventPointer), 0)
	binary.Write(fd, binary.LittleEndian, event)

	fd.Seek(int64(venuePointer), 0)
	binary.Write(fd, binary.LittleEndian, venue)

	fd.Seek(int64(vehiclePointer), 0)
	binary.Write(fd, binary.LittleEndian, vehicle)

	// Write channels
	currentDataPointer := channelsDataPointer
	for i, channel := range f.Channels {
		switch any(channel).(type) {
		case *Channel[float32]:
			currentDataPointer = channel.(*Channel[float32]).Write(fd, uint16(i), head.ChannelsCount, channelsMetaPointer, currentDataPointer)
			break
		case *Channel[int16]:
			currentDataPointer = channel.(*Channel[int16]).Write(fd, uint16(i), head.ChannelsCount, channelsMetaPointer, currentDataPointer)
			break
		case *Channel[int32]:
			currentDataPointer = channel.(*Channel[int32]).Write(fd, uint16(i), head.ChannelsCount, channelsMetaPointer, currentDataPointer)
			break
		}
	}
}

func (f *File) AddChannels(channels ...interface{}) {
	f.Channels = append(f.Channels, channels...)
}

func (c *Channel[T]) Write(
	fd *os.File,
	n uint16,
	channelsCount uint32,
	channelsMetaPointer uintptr,
	currentDataPointer uintptr,
) uintptr {
	var dataType ldfile.DataType
	var previousMetaPointer uintptr = 0
	var nextMetaPointer uintptr = 0

	switch any(c).(type) {
	case *Channel[float32]:
		dataType = ldfile.DataTypeFloat32
		break
	case *Channel[int16]:
		dataType = ldfile.DataTypeInt16
		break
	case *Channel[int32]:
		dataType = ldfile.DataTypeInt32
		break
	}

	if n > 0 {
		previousMetaPointer = channelsMetaPointer + uintptr(binary.Size(ldfile.LdFileChannelMeta{}))*(uintptr(n-1))
	}

	if n < uint16(channelsCount-1) {
		nextMetaPointer = channelsMetaPointer + uintptr(binary.Size(ldfile.LdFileChannelMeta{}))*(uintptr(n+1))
	}

	currentMetaPointer := channelsMetaPointer + uintptr(binary.Size(ldfile.LdFileChannelMeta{}))*uintptr(n)

	channelMeta := ldfile.LdFileChannelMeta{
		PreviousMetaPointer: uint32(previousMetaPointer),
		NextMetaPointer:     uint32(nextMetaPointer),
		DataPointer:         uint32(currentDataPointer),
		DataLength:          uint32(len(*c.Data)),
		ChannelId:           0x2EE1 + n,
		DataType:            dataType.DataType,
		DataTypeLength:      dataType.DataTypeLength,
		Frequency:           c.Frequency,
		Shift:               0,
		Mul:                 1,
		Scale:               1,
		DecPlaces:           0,
	}

	copy(channelMeta.Name[:], c.Name)
	copy(channelMeta.ShortName[:], c.ShortName)
	copy(channelMeta.Unit[:], c.Unit)

	// Convert data to binary slice
	binaryDataWriter := new(bytes.Buffer)
	binary.Write(binaryDataWriter, binary.LittleEndian, c.Data)
	binaryData := binaryDataWriter.Bytes()

	// Write to file
	fd.Seek(int64(currentMetaPointer), 0)
	binary.Write(fd, binary.LittleEndian, channelMeta)

	fd.Seek(int64(currentDataPointer), 0)
	binary.Write(fd, binary.LittleEndian, binaryData)

	// Return next data pointer
	nextDataPointer := currentDataPointer + uintptr(len(binaryData))
	return nextDataPointer
}

func (c *Channel[T]) AddData(data T) {
	*c.Data = append(*c.Data, data)
}
