package adddevice

import (
	"context"

	"github.com/renatocosta55sp/device_management/internal/events"
	"github.com/renatocosta55sp/device_management/internal/infra/adapters/persistence"
	"github.com/renatocosta55sp/modeling/domain"
	"github.com/renatocosta55sp/modeling/infra/bus"
	"github.com/renatocosta55sp/modeling/slice"
)

func WireApp(ctx context.Context, eventBus *bus.EventBus, repo persistence.RepoDevice) (eventRaisedChan chan bus.EventResult) {

	eventChan := make(chan domain.Event)

	eventBus.Subscribe(events.DeviceAddedEvent, eventChan)
	eventBus.Subscribe(events.DeviceUpdatedEvent, eventChan)
	eventBus.Subscribe(events.DeviceRemovedEvent, eventChan)

	eventHandlers := []slice.EventHandler{
		{
			EventName: events.DeviceAddedEvent,
			Handler:   NewDeviceAddedProjection(repo),
		},
		{
			EventName: events.DeviceUpdatedEvent,
			Handler:   NewDeviceUpdatedProjection(repo),
		},
		{
			EventName: events.DeviceRemovedEvent,
			Handler:   NewDeviceRemovedProjection(repo),
		},
	}

	eventRaisedChan = make(chan bus.EventResult)

	eventListener := slice.NewEventListener(eventHandlers, eventBus, eventRaisedChan)
	go eventListener.Listen(ctx, eventChan)

	return
}
