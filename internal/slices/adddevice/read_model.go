package adddevice

import (
	"context"

	"github.com/gookit/event"
	"github.com/renatocosta55sp/device_management/internal/events"
	"github.com/renatocosta55sp/device_management/internal/infra/adapters/persistence"
)

type DeviceReadModel struct {
	repo persistence.RepoDevice
	ctx  context.Context
}

func (d *DeviceReadModel) Handle(e event.Event) error {
	evt := e.Get("data").(events.DeviceAdded)
	_, err := d.repo.Add(&evt, d.ctx)

	if err != nil {
		e.Set("error", err)
	}

	return nil
}
