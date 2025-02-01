package persistence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/modeling/domain"
	"github.com/renatocosta55sp/modeling/eventstore"
	"github.com/renatocosta55sp/modeling/infra"
	"github.com/sirupsen/logrus"
)

const DomainEventTableName = "domain_event_entry"
const DomainEventEntrySeq = "domain_event_entry_seq"

type PersistentEventStore struct {
	store    map[string][]domain.Event
	Conn     *pgxpool.Pool
	DBSchema string
}

func NewPersistentEventStore(conn *pgxpool.Pool, dbSchema string) eventstore.EventStore {
	return &PersistentEventStore{
		Conn:     conn,
		DBSchema: dbSchema,
		store:    make(map[string][]domain.Event),
	}
}

func (p *PersistentEventStore) AppendToStream(ctx context.Context, streamID string, newEvents []domain.Event, nextVersion int) error {
	stream, exists := p.store[streamID]
	if !exists {
		stream = []domain.Event{}
	}

	// Check for version conflict (occ)
	currentVersion := len(stream)
	if currentVersion != nextVersion {
		return eventstore.ErrConcurrencyConflict
	}

	query := fmt.Sprintf(`INSERT INTO %s.%s (global_index, aggregate_identifier, sequence_number, time_stamp, type, meta_data, payload) VALUES (@global_index, @aggregate_identifier, @sequence_number, @time_stamp, @type, @meta_data, @payload)`, p.DBSchema, DomainEventTableName)

	// Use context with timeout for query
	queryCtx, queryCancel := context.WithTimeout(ctx, 5*time.Second)
	defer queryCancel()

	globalIndex, err := GetNextVal(context.Background(), p.Conn, DomainEventEntrySeq)
	if err != nil {
		log.Fatalf("Failed to fetch nextval: %v", err)
	}

	batch := &pgx.Batch{}
	for _, evt := range newEvents {

		metaDataJSON, err := infra.Serialize(infra.EventMetadata{
			UserId: "System",
			Source: "WebApi",
		})
		if err != nil {
			return fmt.Errorf("failed to serialize metadata: %w", err)
		}

		payloadJSON, err := infra.Serialize(evt)
		if err != nil {
			return fmt.Errorf("failed to serialize payload: %w", err)
		}

		args := pgx.NamedArgs{
			"global_index":        globalIndex,
			"aggregateIdentifier": streamID,
			"eventIdentifier":     evt.GetId(),
			"sequence_number":     nextVersion,
			"time_stamp":          time.Now().Format("2006-01-02T15:04:05"),
			"type":                evt.GetName(),
			"meta_data":           metaDataJSON,
			"payload_revision":    1,
			"payload":             payloadJSON,
		}
		logrus.Info("args", args)
		batch.Queue(query, args)
	}

	results := p.Conn.SendBatch(queryCtx, batch)

	var closeErr error

	defer func() {
		if cerr := results.Close(); cerr != nil {
			logrus.Infof("Error closing batch results: %v", cerr)
			closeErr = cerr
		}
	}()

	if err := ExecuteResultSetBatch(results, newEvents, streamID, nextVersion); err != nil {
		return err
	}

	if closeErr != nil {
		return fmt.Errorf("error closing batch results: %w", closeErr)
	}

	return nil

}

func (p *PersistentEventStore) ReadStream(ctx context.Context, streamID string) ([]domain.Event, error) {

	/*result := queryEventStore(aggregateID)
	if result.IsEmpty() {
		return []domain.Event{}, nil
	}*/

	rows, err := p.Conn.Query(ctx, "SELECT payload, type from domain_event_entry where aggregate_identifier = $1", streamID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[[]domain.Event{}])

	var events []domain.Aggregate
	for rows.Next() {

		var e domain.Event
		if err := rows.Scan(&e.AggregateID, &d.Name, &d.Brand); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		devices = append(devices, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return &devices, nil

	return s.store[streamID], nil
}

func (s *PersistentEventStore) ReadAllStream(ctx context.Context) ([]domain.Event, error) {
	return nil, nil
}
