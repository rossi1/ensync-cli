package api

import (
	"context"

	"github.com/rossi1/ensync-cli/internal/domain"
)

type APIClient interface {
	EventService
	AccessKeyService
	Ping(ctx context.Context) error
}

type EventService interface {
	ListEvents(ctx context.Context, params *ListParams) (*domain.EventList, error)
	GetEvent(ctx context.Context, id string) (*domain.Event, error)
	CreateEvent(ctx context.Context, event *domain.Event) error
	UpdateEvent(ctx context.Context, event *domain.Event) error
}

type AccessKeyService interface {
	ListAccessKeys(ctx context.Context, params *ListParams) (*domain.AccessKeyList, error)
	CreateAccessKey(ctx context.Context, key *domain.AccessKey) error
	VerifyAccessKey(ctx context.Context, key string) (bool, error)
}
