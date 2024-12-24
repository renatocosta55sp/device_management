package commands

import "github.com/google/uuid"

type UpdateDeviceCommand struct {
	AggregateID uuid.UUID
	Name, Brand string
}
