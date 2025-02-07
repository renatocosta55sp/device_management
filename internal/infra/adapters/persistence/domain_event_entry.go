package persistence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/modeling/domain"
	"github.com/renatocosta55sp/modeling/eventstore"
	"github.com/renatocosta55sp/modeling/infra"
	"github.com/renatocosta55sp/modeling/infra/bus"
)

const DomainEventTableName = "domain_event_entry"
const DomainEventEntrySeq = "domain_event_entry_seq"

type PersistentEventStore struct {
	store         map[string][]domain.Event
	Conn          *pgxpool.Pool
	eventRegistry *bus.EventRegistry
	DBSchema      string
}

func NewPersistentEventStore(conn *pgxpool.Pool, eventRegistry *bus.EventRegistry, dbSchema string) eventstore.EventStore {
	return &PersistentEventStore{
		store:         make(map[string][]domain.Event),
		Conn:          conn,
		eventRegistry: eventRegistry,
		DBSchema:      dbSchema,
	}
}

func (p *PersistentEventStore) AppendToStream(ctx context.Context, streamID string, newEvents []domain.Event, expectedVersion int) error {

	stream, exists := p.store[streamID]
	if !exists {
		stream = []domain.Event{}
	}

	/*if currentVersion != aggregateVersion {
		return eventstore.ErrConcurrencyConflict
	}*/

	entries := [][]any{}
	columns := []string{"global_index", "aggregate_identifier", "event_identifier", "sequence_number", "time_stamp", "type", "meta_data", "payload"}

	metaDataSerialized, err := infra.Serialize(infra.EventMetadata{
		UserId: "System",
		Source: "WebApi",
	})
	if err != nil {
		return fmt.Errorf("failed to serialize metadata: %w", err)
	}

	for _, event := range newEvents {

		globalIndex, err := GetNextVal(context.Background(), p.Conn, DomainEventEntrySeq)
		if err != nil {
			log.Fatalf("Failed to fetch nextval: %v", err)
		}

		payloadSerialized, err := infra.Serialize(event)
		if err != nil {
			return fmt.Errorf("failed to serialize payload: %w", err)
		}

		entries = append(entries, []any{globalIndex, streamID, event.GetId(), expectedVersion, time.Now().Format("2006-01-02T15:04:05"), event.GetName(), metaDataSerialized, payloadSerialized})
	}

	_, err = p.Conn.CopyFrom(
		ctx,
		pgx.Identifier{DomainEventTableName},
		columns,
		pgx.CopyFromRows(entries),
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" && pgErr.ConstraintName == "domain_event_entry_aggregate_identifier_sequence_number_key" {
			return eventstore.ErrConcurrencyConflict
		}
		return fmt.Errorf("error copying into %s table: %w", DomainEventTableName, err)
	}

	p.store[streamID] = append(stream, newEvents...)

	return nil

}

func (p *PersistentEventStore) ReadStream(ctx context.Context, streamID string) ([]domain.Event, error) {

	rows, err := p.Conn.Query(ctx, "SELECT payload, type from domain_event_entry where aggregate_identifier = $1", streamID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {

		var payload []byte
		var eventType string

		if err := rows.Scan(&payload, &eventType); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		event, err := p.eventRegistry.CreateEvent(eventType, payload)
		if err != nil {
			return nil, fmt.Errorf("failed to create event: %w", err)
		}

		p.store[streamID] = append(p.store[streamID], event)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return p.store[streamID], nil
}

func (s *PersistentEventStore) ReadAllStream(ctx context.Context) ([]domain.Event, error) {
	return nil, nil
}
