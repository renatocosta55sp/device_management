package domain

import (
	"errors"

	"github.com/google/uuid"
	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/renatocosta55sp/device_management/internal/events"
	"github.com/renatocosta55sp/modeling/domain"
	"github.com/renatocosta55sp/modeling/slice"
)

type DeviceAggregate struct {
	domain.AggregateRoot
	Name, Brand string
}

func NewDevice(
	aggregateId uuid.UUID,
) *DeviceAggregate {

	device := &DeviceAggregate{
		AggregateRoot: domain.AggregateRoot{AggregateID: aggregateId, Version: DeviceAggregateVersion},
	}

	return device
}

func (d *DeviceAggregate) HandleAdd(command commands.AddDeviceCommand) (slice.CommandResult, error) {

	if command.Name == "" {
		return slice.CommandResult{
			Identifier:        command.AggregateID,
			AggregateSequence: DeviceAggregateVersion,
		}, ErrEmptyName
	}

	if command.Brand == "" {
		return slice.CommandResult{
			Identifier:        command.AggregateID,
			AggregateSequence: DeviceAggregateVersion,
		}, ErrEmptyBrand
	}

	d.AggregateRoot.RecordThat(
		domain.Event{
			Type: events.DeviceAddedEvent,
			Data: events.DeviceAdded{
				AggregateId: command.AggregateID,
				Name:        command.Name,
				Brand:       command.Brand,
			},
		},
	)

	return slice.CommandResult{
		Identifier:        command.AggregateID,
		AggregateSequence: DeviceAggregateVersion,
	}, nil

}

func (d *DeviceAggregate) HandleUpdate(command commands.UpdateDeviceCommand) (slice.CommandResult, error) {

	if command.Name == "" {
		return slice.CommandResult{
			Identifier:        command.AggregateID,
			AggregateSequence: DeviceAggregateVersion,
		}, ErrEmptyName
	}

	if command.Brand == "" {
		return slice.CommandResult{
			Identifier:        command.AggregateID,
			AggregateSequence: DeviceAggregateVersion,
		}, ErrEmptyBrand
	}

	d.AggregateRoot.RecordThat(
		domain.Event{
			Type: events.DeviceUpdatedEvent,
			Data: events.DeviceUpdated{
				AggregateId: command.AggregateID,
				Name:        command.Name,
				Brand:       command.Brand,
			},
		},
	)

	return slice.CommandResult{
		Identifier:        command.AggregateID,
		AggregateSequence: DeviceAggregateVersion,
	}, nil

}

func (d *DeviceAggregate) HandleDelete(command commands.RemoveDeviceCommand) (slice.CommandResult, error) {

	d.AggregateRoot.RecordThat(
		domain.Event{
			Type: events.DeviceRemovedEvent,
			Data: events.DeviceRemoved{
				AggregateId: command.AggregateID,
				Name:        command.Name,
				Brand:       command.Brand,
			},
		},
	)

	return slice.CommandResult{
		Identifier:        command.AggregateID,
		AggregateSequence: DeviceAggregateVersion,
	}, nil

}

var DeviceAggregateVersion = int8(1)

var (
	ErrEmptyName  = errors.New("error.device.name.required")
	ErrEmptyBrand = errors.New("error.device.brand.required")
)
