// internal/repository/mock_repository.go (例)
package repository

import (
	"context"
	"database/sql"
	"errors"
)

// MockRepository はテスト用のモック実装です。
type MockRepository struct {
	BeginTxFunc           func(ctx context.Context) (*sql.Tx, error)
	CreateEventTxFunc     func(ctx context.Context, tx *sql.Tx, params CreateEventParams) (int, string, error)
	GetEventsFunc         func(ctx context.Context) ([]*Event, error)
	AuthenticateEventFunc func(ctx context.Context, id int, authCode string) error
	GetEventFunc          func(ctx context.Context, id int) (*Event, error)
}

func (m *MockRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	if m.BeginTxFunc != nil {
		return m.BeginTxFunc(ctx)
	}
	return nil, errors.New("BeginTx not implemented")
}

func (m *MockRepository) CreateEventTx(ctx context.Context, tx *sql.Tx, params CreateEventParams) (int, string, error) {
	if m.CreateEventTxFunc != nil {
		return m.CreateEventTxFunc(ctx, tx, params)
	}
	return 0, "", errors.New("CreateEventTx not implemented")
}

func (m *MockRepository) GetEvents(ctx context.Context) ([]*Event, error) {
	if m.GetEventsFunc != nil {
		return m.GetEventsFunc(ctx)
	}
	return nil, errors.New("GetEvents not implemented")
}

func (m *MockRepository) AuthenticateEvent(ctx context.Context, id int, authCode string) error {
	if m.AuthenticateEventFunc != nil {
		return m.AuthenticateEventFunc(ctx, id, authCode)
	}
	return errors.New("AuthenticateEvent not implemented")
}

func (m *MockRepository) GetEvent(ctx context.Context, id int) (*Event, error) {
	if m.GetEventFunc != nil {
		return m.GetEventFunc(ctx, id)
	}
	return nil, errors.New("GetEvent not implemented")
}
