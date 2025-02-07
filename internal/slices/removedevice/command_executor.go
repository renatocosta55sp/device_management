package removedevice

import (
	"context"

	"github.com/renatocosta55sp/device_management/internal/domain"
	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/renatocosta55sp/modeling/eventstore"
	"github.com/renatocosta55sp/modeling/infra/bus"
	"github.com/renatocosta55sp/modeling/slice"
)

type CommandExecutor struct {
	eventStore eventstore.EventStore
}

func (c CommandExecutor) Send(ctx context.Context, cmd commands.RemoveDeviceCommand) (commandResult slice.CommandResult, device *domain.DeviceAggregate, err error) {
	//Get the current state
	stream, err := c.eventStore.ReadStream(ctx, cmd.AggregateID.String())

	if err != nil {
		return commandResult, nil, err
	}

	deviceAggregate := domain.NewDeviceAggregate(stream)

	commandResult, err = deviceAggregate.Remove(cmd)
	if err != nil {
		return commandResult, nil, err
	}

	dispatcher := bus.NewEventDispatcher()

	deviceReadModel := DeviceReadModel{deviceAggregate: deviceAggregate, eventStore: c.eventStore, ctx: ctx}
	bus.RegisterHandler(dispatcher, deviceReadModel)

	if err := dispatcher.DispatchUncommittedEvents(deviceAggregate.UncommittedEvents); err != nil {
		return commandResult, nil, err
	}

	return commandResult, deviceAggregate, nil
}
