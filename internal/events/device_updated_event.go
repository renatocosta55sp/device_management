package events

import "github.com/google/uuid"

const DeviceUpdatedEvent = "DeviceUpdatedEvent"

type DeviceUpdated struct {
	AggregateId uuid.UUID
	Name, Brand string
}
