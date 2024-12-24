package commands

import "github.com/google/uuid"

type AddDeviceCommand struct {
	AggregateID uuid.UUID
	Name, Brand string
}
