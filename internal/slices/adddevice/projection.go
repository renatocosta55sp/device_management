package adddevice

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
type DeviceAddedProjection struct {
	repo persistence.RepoDevice
}
type DeviceUpdatedProjection struct {
	repo persistence.RepoDevice
}

type DeviceRemovedProjection struct {
	repo persistence.RepoDevice
}

func NewDeviceAddedProjection(repo persistence.RepoDevice) slice.EventHandleable {
	return &DeviceAddedProjection{
		repo: repo,
	}
}

func (d DeviceAddedProjection) On(ctx context.Context, event domain.Event) error {
	evt := event.Data.(events.DeviceAdded)
	_, err := d.repo.Add(&evt, ctx)
	return err
}

func NewDeviceUpdatedProjection(repo persistence.RepoDevice) slice.EventHandleable {
	return &DeviceUpdatedProjection{repo: repo}
}

func (d DeviceUpdatedProjection) On(ctx context.Context, event domain.Event) error {
	evt := event.Data.(events.DeviceUpdated)
	return d.repo.Update(&evt, ctx)
}

func NewDeviceRemovedProjection(repo persistence.RepoDevice) slice.EventHandleable {
	return &DeviceRemovedProjection{repo: repo}
}

func (d DeviceRemovedProjection) On(ctx context.Context, event domain.Event) error {
	evt := event.Data.(events.DeviceRemoved)
	return d.repo.Remove(&evt, ctx)
}
