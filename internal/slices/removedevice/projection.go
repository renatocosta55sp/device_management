package removedevice

import (
	"context"

	"github.com/renatocosta55sp/device_management/internal/events"
	"github.com/renatocosta55sp/device_management/internal/infra/adapters/persistence"
	"github.com/renatocosta55sp/modeling/domain"
	"github.com/renatocosta55sp/modeling/slice"
)

type DeviceProjection struct {
	repo persistence.RepoDevice
}

func NewDeviceProjection(repo persistence.RepoDevice) slice.EventHandleable {
	return &DeviceProjection{
		repo: repo,
	}
}

func (d DeviceProjection) On(ctx context.Context, event domain.Event) error {
	evt := event.Data.(events.DeviceRemoved)
	return d.repo.Remove(&evt, ctx)
}
