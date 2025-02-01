package adddevice

import (
	"context"

	"github.com/renatocosta55sp/device_management/internal/domain"
	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/renatocosta55sp/modeling/eventstore"
	"github.com/renatocosta55sp/modeling/infra/bus"
	"github.com/renatocosta55sp/modeling/slice"
)

type CommandExecutor struct {
	store eventstore.EventStore
}

func (c CommandExecutor) Send(ctx context.Context, cmd commands.AddDeviceCommand, eventStore eventstore.EventStore) (commandResult slice.CommandResult, device *domain.DeviceAggregate, err error) {

	//Get the current state
	stream, err := c.store.ReadStream(ctx, cmd.AggregateID.String())

	if err != nil {
		return commandResult, device, err
	}

	deviceAggregate := domain.NewDeviceAggregate(stream)

	commandResult, err = deviceAggregate.Add(cmd)
	if err != nil {
		return commandResult, device, err
	}

	dispatcher := bus.NewEventDispatcher()

	deviceReadModel := DeviceReadModel{deviceAggregate: deviceAggregate, eventStore: eventStore, ctx: ctx}
	bus.RegisterHandler(dispatcher, deviceReadModel)

	if err := dispatcher.DispatchUncommittedEvents(device.UncommittedEvents); err != nil {
		return commandResult, device, err
	}

	return commandResult, device, nil
}
