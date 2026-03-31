package eventlog

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	dbgen "github.com/squeakycheese75/paytunnel/internal/db"
	"github.com/squeakycheese75/paytunnel/internal/repository"
	"github.com/squeakycheese75/paytunnel/internal/signing"
)

type EventRecord = dbgen.Event

type Repo interface {
	ListEvents(ctx context.Context) ([]repository.Event, error)
	GetEvent(ctx context.Context, deliveryID string) (repository.Event, error)
}

type EventLog struct {
	repo Repo
}

func NewEventLog(repo Repo) *EventLog {
	return &EventLog{
		repo: repo,
	}
}

func (e *EventLog) List(ctx context.Context) ([]repository.Event, error) {
	events, err := e.repo.ListEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	return events, nil
}

func (e *EventLog) ReplayEvent(ctx context.Context, deliveryID string) error {
	event, err := e.repo.GetEvent(ctx, deliveryID)
	if err != nil {
		return fmt.Errorf("get event: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		event.TargetUrl,
		bytes.NewBufferString(event.BodyJson),
	)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("BTCPay-Sig", "sha256="+signing.BTCPaySignature([]byte(event.BodyJson), event.Secret))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send replay: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	log.Printf("replayed %s -> %s\n", deliveryID, resp.Status)
	return nil
}
