package adddevice

import (
	"context"

	"github.com/gookit/event"
	"github.com/renatocosta55sp/device_management/internal/domain"
	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/renatocosta55sp/device_management/internal/events"
	"github.com/renatocosta55sp/device_management/internal/infra/adapters/persistence"
	"github.com/renatocosta55sp/modeling/slice"
)

type CommandGateway struct{}

func (c CommandGateway) Send(ctx context.Context, command commands.AddDeviceCommand, repo persistence.RepoDevice) (commandResult slice.CommandResult, device *domain.DeviceAggregate, err error) {

	deviceReadModel := DeviceReadModel{repo: repo, ctx: ctx}
	event.On(events.DeviceAddedEvent, &deviceReadModel)

	device = domain.NewDevice(command.AggregateID)

	commandResult, err = device.HandleAdd(command)
	if err != nil {
		return commandResult, device, err
	}

	for _, evt := range device.Events {
		eventHandlerResult := event.MustFire(evt.Type, event.M{"data": evt.Data})

		if errValue, ok := eventHandlerResult.Get("error").(error); ok {
			return commandResult, device, errValue
		}
	}

	return commandResult, device, nil
}
