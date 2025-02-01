package events

import "github.com/google/uuid"

type DeviceAdded struct {
	AggregateId uuid.UUID
	Name, Brand string
}

func (e DeviceAdded) GetName() string {
	return "DeviceAddedEvent"
}

func (e DeviceAdded) GetId() uuid.UUID {
	return uuid.New()
}
