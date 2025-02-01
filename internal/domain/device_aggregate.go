package domain

import (
	"errors"

	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/renatocosta55sp/device_management/internal/events"
	"github.com/renatocosta55sp/modeling/domain"
	"github.com/renatocosta55sp/modeling/slice"
)

type DeviceAggregate struct {
	domain.Aggregate
	Name, Brand string
}

func NewDeviceAggregate(stream []domain.Event) *DeviceAggregate {
	d := &DeviceAggregate{}
	d.hydrate(stream)
	return d
}

func (d *DeviceAggregate) hydrate(stream []domain.Event) {
	for _, e := range stream {
		d.Apply(e)
	}
}

func (d *DeviceAggregate) Apply(event domain.Event) {

	switch e := event.(type) {
	case events.DeviceAdded:
		d.AggregateID = e.AggregateId
		d.Name = e.Name
		d.Brand = e.Brand
	}

}

func (d *DeviceAggregate) Add(cmd commands.AddDeviceCommand) (slice.CommandResult, error) {

	commandResult := slice.CommandResult{
		Identifier:        cmd.AggregateID,
		AggregateSequence: d.Version,
	}

	if cmd.Name == "" {
		return commandResult, ErrEmptyName
	}

	if cmd.Brand == "" {
		return commandResult, ErrEmptyBrand
	}

	event := events.DeviceAdded{
		AggregateId: cmd.AggregateID,
		Name:        cmd.Name,
		Brand:       cmd.Brand,
	}

	d.Version += 1
	d.UncommittedEvents = append(d.UncommittedEvents, event)

	d.Apply(event)

	return commandResult, nil

}

var (
	ErrEmptyName  = errors.New("error.device.name.required")
	ErrEmptyBrand = errors.New("error.device.brand.required")
)
