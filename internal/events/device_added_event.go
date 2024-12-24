package events

import "github.com/google/uuid"

const DeviceAddedEvent = "DeviceAddedEvent"

type DeviceAdded struct {
	AggregateId uuid.UUID
	Name, Brand string
}
