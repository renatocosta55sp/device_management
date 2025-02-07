package events

import "github.com/google/uuid"

type DeviceRemoved struct {
	AggregateId uuid.UUID
	Version     int
}

func (e *DeviceRemoved) GetName() string {
	return "DeviceRemoved"
}

func (e *DeviceRemoved) GetId() uuid.UUID {
	return uuid.New()
}

func (e *DeviceRemoved) GetVersion() int {
	return e.Version
}

func (e *DeviceRemoved) SetVersion(number int) {
	e.Version = number
}
