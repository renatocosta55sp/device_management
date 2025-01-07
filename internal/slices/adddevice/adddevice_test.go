package adddevice

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/device_management/internal/domain"

	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/renatocosta55sp/device_management/internal/events"
	"github.com/renatocosta55sp/device_management/internal/infra/adapters/persistence"
	"github.com/renatocosta55sp/device_management/internal/infra/testsuite"
	"github.com/renatocosta55sp/modeling/infra/bus"
	"github.com/renatocosta55sp/modeling/slice"
	"github.com/stretchr/testify/assert"
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
	command := commands.AddDeviceCommand{
		AggregateID: aggregateIdentifier,
		Name:        "IOS",
		Brand:       "Apple",
	}

	commandResult, device, err := CommandGateway(ctx,
		command,
		*persistence.NewDeviceRepository(dbConn, "public"),
	)

	if err != nil {
		ag.T.Fatal(err)
	}

	for _, evt := range device.Events {
		raisedEvents[evt.Type] = evt.Type
	}

	commandResultToCompare := slice.CommandResult{
		Identifier:        aggregateIdentifier,
		AggregateSequence: domain.DeviceAggregateVersion,
	}

	assert.Equal(ag.T, commandResult, commandResultToCompare, "The CommandResult should be equal")

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
		When(raisedEvents).
		Then(
			events.DeviceAddedEvent,
		).
		Assert()

}
