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

	lenStream := len(stream)

	if lenStream > 0 {
		for _, e := range stream {
			d.Apply(e)
		}
		d.Version = stream[lenStream-1].GetVersion()
	}

}

func (d *DeviceAggregate) Apply(event domain.Event) {

	switch e := event.(type) {
	case *events.DeviceAdded:
		d.AggregateID = e.AggregateId
		d.Name = e.Name
		d.Brand = e.Brand

	case *events.DeviceRemoved:
		d.AggregateID = e.AggregateId
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

	event := &events.DeviceAdded{
		AggregateId: cmd.AggregateID,
		Name:        cmd.Name,
		Brand:       cmd.Brand,
	}

	event.SetVersion(d.Version + 1)

	d.UncommittedEvents = append(d.UncommittedEvents, event)

	d.Apply(event)

	d.Version++

	return commandResult, nil

}

func (d *DeviceAggregate) Remove(cmd commands.RemoveDeviceCommand) (slice.CommandResult, error) {

	commandResult := slice.CommandResult{
		Identifier:        cmd.AggregateID,
		AggregateSequence: d.Version,
	}

	event := &events.DeviceRemoved{
		AggregateId: cmd.AggregateID,
	}

	event.SetVersion(d.Version + 1)

	d.UncommittedEvents = append(d.UncommittedEvents, event)

	d.Apply(event)

	d.Version++

	return commandResult, nil

}

var (
	ErrEmptyName  = errors.New("error.device.name.required")
	ErrEmptyBrand = errors.New("error.device.brand.required")
)
