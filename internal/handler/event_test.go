// internal/handler/event_test.go
package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	// テスト対象のパッケージ
	"github.com/ras0q/go-backend-template/internal/handler"
	"github.com/ras0q/go-backend-template/internal/repository"
)

// TestCreateEvent_Success は正常系のハンドラテスト例
func TestCreateEvent_Success(t *testing.T) {
	// 1) Repositoryモックを用意
	mockRepo := &repository.MockRepository{
		BeginTxFunc: func(ctx context.Context) (*repository.MockTx, error) {
			// ここでは MockTx struct を作って返すイメージ
			// 実装は省略
			return &repository.MockTx{}, nil
		},
		CreateEventTxFunc: func(ctx context.Context, tx *repository.MockTx, params repository.CreateEventParams) (int, string, error) {
			return 123, "abc-auth-code", nil
		},
	}

	// 2) Slack送信をモック
	originalSlackPoster := handler.SlackPoster
	defer func() { handler.SlackPoster = originalSlackPoster }()
	handler.SlackPoster = func(msg string) error {
		// 成功するモック
		return nil
	}

	// 3) Handlerを作り、POSTリクエスト
	h := handler.New(mockRepo)
	e := echo.New()

	// リクエストボディ
	reqBody := `{
        "title":"Test Event",
        "organizer":"Test Org",
        "startDate":"2025-01-01",
        "startTime":"09:00:00",
        "endDate":"2025-01-01",
        "endTime":"10:00:00",
        "email":"test@example.com"
    }`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/event/new", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	// テスト実行
	err := h.CreateEvent(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	// レスポンスをチェック
	var got handler.CreateEventResponse
	err = json.Unmarshal(rec.Body.Bytes(), &got)
	require.NoError(t, err)
	require.Equal(t, "123", got.ID)
}

// TestCreateEvent_SlackFail はSlack通知が失敗するケース
func TestCreateEvent_SlackFail(t *testing.T) {
	mockRepo := &repository.MockRepository{
		BeginTxFunc: func(ctx context.Context) (*repository.MockTx, error) {
			return &repository.MockTx{}, nil
		},
		CreateEventTxFunc: func(ctx context.Context, tx *repository.MockTx, params repository.CreateEventParams) (int, string, error) {
			return 123, "abc-auth-code", nil
		},
	}

	originalSlackPoster := handler.SlackPoster
	defer func() { handler.SlackPoster = originalSlackPoster }()
	handler.SlackPoster = func(msg string) error {
		return assertAnError // 適当なエラーを返す
	}

	h := handler.New(mockRepo)
	e := echo.New()

	reqBody := `{
        "title":"Fail Slack",
        "organizer":"Test Org",
        "startDate":"2025-01-01",
        "startTime":"09:00:00",
        "endDate":"2025-01-01",
        "endTime":"10:00:00",
        "email":"test@example.com"
    }`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/event/new", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	err := h.CreateEvent(c)

	// Slack通知で失敗 → 500エラーを想定
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusInternalServerError, httpErr.Code)
}

// ほかにも、バリデーションエラー、DBエラーなど様々なケースを網羅
