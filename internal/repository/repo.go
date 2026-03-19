package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/squeakycheese75/paytunnel/internal/db"
)

type EventRepository struct {
	q *db.Queries
}

func NewEventRepository(database *sql.DB) *EventRepository {
	return &EventRepository{
		q: db.New(database),
	}
}

func (r *EventRepository) Create(ctx context.Context, event Event) error {
	return r.q.InsertEvent(ctx, db.InsertEventParams{
		DeliveryID: event.DeliveryID,
		EventName:  event.EventName,
		TargetUrl:  event.TargetUrl,
		BodyJson:   event.BodyJson,
		Secret:     event.Secret,
	})
}

func (r *EventRepository) ListEvents(ctx context.Context) ([]Event, error) {
	rows, err := r.q.ListEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	events := make([]Event, 0, len(rows))
	for _, row := range rows {
		events = append(events, Event{
			DeliveryID: row.DeliveryID,
			EventName:  row.EventName,
			TargetUrl:  row.TargetUrl,
			BodyJson:   row.BodyJson,
			Secret:     row.Secret,
		})
	}

	return events, nil
}

func (r *EventRepository) GetEvent(ctx context.Context, deliveryID string) (Event, error) {
	row, err := r.q.GetEvent(ctx, deliveryID)
	if err != nil {
		return Event{}, fmt.Errorf("get event: %w", err)
	}

	return Event{
		DeliveryID: row.DeliveryID,
		EventName:  row.EventName,
		TargetUrl:  row.TargetUrl,
		BodyJson:   row.BodyJson,
		Secret:     row.Secret,
	}, nil
}
