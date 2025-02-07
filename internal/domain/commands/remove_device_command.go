package commands

import "github.com/google/uuid"

type RemoveDeviceCommand struct {
	AggregateID uuid.UUID
}
