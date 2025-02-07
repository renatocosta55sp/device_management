package adddevice

import (
	"context"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/renatocosta55sp/device_management/internal/events"
	"github.com/renatocosta55sp/device_management/internal/infra/adapters/persistence"
	"github.com/renatocosta55sp/device_management/internal/infra/testsuite"
	"github.com/renatocosta55sp/modeling/infra/bus"
	"github.com/testcontainers/testcontainers-go"
)

var ag = &bus.AggregateRootTestCase{}
var ctx context.Context
var raisedEvents map[string]string
var ctxCancFunc context.CancelFunc
var dbConn *pgxpool.Pool
var pgContainer testcontainers.Container
var container testcontainers.Container
var err error

func init() {

	ctx, ctxCancFunc = context.WithTimeout(context.Background(), 5*time.Second)
	raisedEvents = make(map[string]string)

	dbConn, container, err = testsuite.InitTestContainer()
	if err != nil {
		log.Fatalf("Failed to initialize test container: %v", err)
	}
	pgContainer = container

}

func runAddCommand() {

	aggregateIdentifier := uuid.New()

	evt := events.DeviceAdded{}
	eventRegistry := bus.NewEventRegistry()
	eventRegistry.RegisterEvents(map[string]reflect.Type{
		evt.GetName(): reflect.TypeOf(events.DeviceAdded{}),
	})

	_, deviceAggregate, err := CommandExecutor{
		persistence.NewPersistentEventStore(dbConn, eventRegistry, "public"),
	}.Send(
		ctx,
		commands.AddDeviceCommand{
			AggregateID: aggregateIdentifier,
			Name:        "IOS",
			Brand:       "Apple",
		},
	)

	if err != nil {
		ag.T.Fatal(err)
	}

	for _, evt := range deviceAggregate.UncommittedEvents {
		raisedEvents[evt.GetName()] = evt.GetName()
	}

	/*commandResultToCompare := slice.CommandResult{
		Identifier: aggregateIdentifier,
	}

	assert.Equal(ag.T, commandResult, commandResultToCompare, "The CommandResult should be equal")*/

}

func TestAddDevice(t *testing.T) {

	// Clean up the pg container
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	ag.T = t

	evt := events.DeviceAdded{}
	ag.
		Given(runAddCommand).
		When(raisedEvents).
		Then(
			evt.GetName(),
		).
		Assert()

}
