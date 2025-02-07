package adddevice

import (
	"context"

	"github.com/renatocosta55sp/device_management/internal/domain"
	"github.com/renatocosta55sp/device_management/internal/events"
	"github.com/renatocosta55sp/modeling/eventstore"
)

type DeviceReadModel struct {
	deviceAggregate *domain.DeviceAggregate
	eventStore      eventstore.EventStore
	ctx             context.Context
}

func (d DeviceReadModel) Handle(event *events.DeviceAdded) error {

	if err := d.eventStore.AppendToStream(d.ctx,
		event.AggregateId.String(),
		d.deviceAggregate.UncommittedEvents,
		d.deviceAggregate.Version); err != nil {
		return err
	}

	return nil
}
