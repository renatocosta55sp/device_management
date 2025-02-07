package main

import (
	"reflect"

	"github.com/renatocosta55sp/device_management/internal/events"
	"github.com/renatocosta55sp/modeling/infra/bus"
)

func configureDomainEventRegistry() *bus.EventRegistry {

	eventRegistry := bus.NewEventRegistry()
	eventRegistry.RegisterEvents(map[string]reflect.Type{
		(&events.DeviceAdded{}).GetName():   reflect.TypeOf(events.DeviceAdded{}),
		(&events.DeviceRemoved{}).GetName(): reflect.TypeOf(events.DeviceRemoved{}),
	})

	return eventRegistry

}
