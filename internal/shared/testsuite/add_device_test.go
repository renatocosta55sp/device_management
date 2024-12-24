package testsuite

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/renatocosta55sp/device_management/internal/domain"
	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/renatocosta55sp/device_management/internal/events"
	"github.com/renatocosta55sp/device_management/internal/infra/adapters/persistence"
	"github.com/renatocosta55sp/device_management/internal/slices/adddevice"
	"github.com/renatocosta55sp/modeling/infra/bus"
	"github.com/renatocosta55sp/modeling/slice"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

var ag = &bus.AggregateRootTestCase{}
var eventBus = bus.NewEventBus()
var ctx context.Context
var ctxCancFunc context.CancelFunc
var eventResultChan chan bus.EventResult
var pgContainer testcontainers.Container

func init() {

	ctx, ctxCancFunc = context.WithTimeout(context.Background(), 5*time.Second)

	dbConn, container, err := InitTestContainer()
	if err != nil {
		log.Fatalf("Failed to initialize test container: %v", err)
	}
	pgContainer = container

	eventResultChan = adddevice.WireApp(ctx,
		eventBus,
		*persistence.NewDeviceRepository(dbConn, "public"),
	)
}

func runAddCommand() {

	aggregateIdentifier := uuid.New()
	command := commands.AddDeviceCommand{
		AggregateID: aggregateIdentifier,
		Name:        "IOS",
		Brand:       "Apple",
	}

	device := domain.NewDevice(aggregateIdentifier)

	commandResult, err := device.HandleAdd(command)
	if err != nil {
		ag.T.Fatal(err)
	}

	commandResultToCompare := slice.CommandResult{
		Identifier:        aggregateIdentifier,
		AggregateSequence: domain.DeviceAggregateVersion,
	}

	assert.Equal(ag.T, commandResult, commandResultToCompare, "The CommandResult should be equal")

	err = (&slice.CommandExecutionResult{
		EventBus:        eventBus,
		CtxCancFunc:     ctxCancFunc,
		EventResultChan: eventResultChan,
	}).Execute(device.Events)

	if err != nil {
		ag.T.Fatal(err)
	}

}

func TestAddDevice(t *testing.T) {

	// Clean up the pg container
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	ag.T = t

	ag.
		Given(runAddCommand).
		When(eventBus.RaisedEvents()).
		Then(
			events.DeviceAddedEvent,
		).
		Assert()

}
