package events

import "github.com/google/uuid"

const DeviceRemovedEvent = "DeviceRemovedEvent"

type DeviceRemoved struct {
	AggregateId uuid.UUID
	Name, Brand string
}
