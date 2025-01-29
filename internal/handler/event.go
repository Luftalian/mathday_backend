package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/labstack/echo/v4"

	// リポジトリとの連携
	"github.com/ras0q/go-backend-template/internal/pkg/config"
	"github.com/ras0q/go-backend-template/internal/repository"
)

// -------------------
// 変換用ヘルパー関数
// -------------------
func convertSpeakers(speakers []Speaker) []repository.Speaker {
	result := make([]repository.Speaker, len(speakers))
	for i, s := range speakers {
		result[i] = repository.Speaker{
			Name:         s.Name,
			Title:        s.Title,
			Organization: s.Organization,
		}
	}
	return result
}

func convertSpeakersToResponse(speakers []repository.Speaker) []Speaker {
	result := make([]Speaker, len(speakers))
	for i, s := range speakers {
		result[i] = Speaker{
			Name:         s.Name,
			Title:        s.Title,
			Organization: s.Organization,
		}
	}
	return result
}

func convertSchedules(schedules []Schedule) []repository.Schedule {
	result := make([]repository.Schedule, len(schedules))
	for i, s := range schedules {
		result[i] = repository.Schedule{
			Time:    s.Time,
			Title:   s.Title,
			Speaker: s.Speaker,
		}
	}
	return result
}

func convertSchedulesToResponse(schedules []repository.Schedule) []Schedule {
	result := make([]Schedule, len(schedules))
	for i, s := range schedules {
		result[i] = Schedule{
			Time:    s.Time,
			Title:   s.Title,
			Speaker: s.Speaker,
		}
	}
	return result
}

// ---------------------
// リクエスト/レスポンス
// ---------------------
type (
	CreateEventRequest struct {
		Title            string     `json:"title"`
		Organizer        string     `json:"organizer"`
		StartDate        string     `json:"startDate"`
		StartTime        string     `json:"startTime"`
		EndDate          string     `json:"endDate"`
		EndTime          string     `json:"endTime"`
		Email            string     `json:"email"`
		Prefecture       *string    `json:"prefecture"`
		EventType        *string    `json:"eventType"`
		IsOnline         bool       `json:"isOnline"`
		IsOffline        bool       `json:"isOffline"`
		OfficialURL      *string    `json:"officialUrl"`
		OnlineLectureURL *string    `json:"onlineLectureUrl"`
		Venue            *string    `json:"venue"`
		Target           *string    `json:"target"`
		Capacity         *string    `json:"capacity"`
		Description      *string    `json:"description"`
		Tags             []string   `json:"tags"`
		Speakers         []Speaker  `json:"speakers"`
		Schedule         []Schedule `json:"schedule"`
	}

	CreateEventResponse struct {
		ID string `json:"id"`
	}

	UpdateEventResponse struct {
		Message string `json:"message"`
	}

	GetEventResponse struct {
		ID               int        `json:"id"`
		Title            string     `json:"title"`
		Organizer        string     `json:"organizer"`
		StartDate        string     `json:"startDate"`
		StartTime        string     `json:"startTime"`
		EndDate          string     `json:"endDate"`
		EndTime          string     `json:"endTime"`
		Email            string     `json:"email"`
		Prefecture       *string    `json:"prefecture"`
		EventType        *string    `json:"eventType"`
		IsOnline         bool       `json:"isOnline"`
		IsOffline        bool       `json:"isOffline"`
		OfficialURL      *string    `json:"officialUrl"`
		OnlineLectureURL *string    `json:"onlineLectureUrl"`
		Venue            *string    `json:"venue"`
		Target           *string    `json:"target"`
		Capacity         *string    `json:"capacity"`
		Description      *string    `json:"description"`
		Tags             []string   `json:"tags"`
		Speakers         []Speaker  `json:"speakers"`
		Schedule         []Schedule `json:"schedule"`
	}

	Speaker struct {
		Name         string `json:"name"`
		Title        string `json:"title"`
		Organization string `json:"organization"`
	}

	Schedule struct {
		Time    string `json:"time"`
		Title   string `json:"title"`
		Speaker string `json:"speaker"`
	}
)

// -------------------
// Slack連携
// -------------------
var slackWebhookURL string

// InitSlack は環境変数からWebhook URLを取得
func InitSlack() {
	slackWebhookURL = os.Getenv("SLACK_WEBHOOK_URL")
}

// postToSlack はSlackへメッセージをPOST送信します
func postToSlack(msg string) error {
	if slackWebhookURL == "" {
		return fmt.Errorf("slackWebhookURLが設定されていません")
	}
	payload := map[string]string{"text": msg}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("payloadのJSON変換に失敗: %w", err)
	}

	req, err := http.NewRequest("POST", slackWebhookURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Slackへのリクエスト作成に失敗しました: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Slackへのリクエスト送信に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slackからの応答が異常です: %s", resp.Status)
	}
	return nil
}

// -------------------
//  ハンドラ実装
// -------------------

// POST /api/v1/event/new
// 通知まで含めて失敗ならROLLBACKする
func (h *Handler) CreateEvent(c echo.Context) error {
	req := new(CreateEventRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request").SetInternal(err)
	}

	// バリデーション
	if err := vd.ValidateStruct(
		req,
		vd.Field(&req.Title, vd.Required),
		vd.Field(&req.Organizer, vd.Required),
		vd.Field(&req.StartDate, vd.Required),
		vd.Field(&req.StartTime, vd.Required),
		vd.Field(&req.EndDate, vd.Required),
		vd.Field(&req.EndTime, vd.Required),
		vd.Field(&req.Email, vd.Required, is.Email),
	); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Errorf("invalid request body: %w", err))
	}

	ctx := c.Request().Context()

	// 1) トランザクション開始
	tx, err := h.repo.BeginTx(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			"failed to begin transaction").SetInternal(err)
	}
	// deferでROLLBACKを仕込む（後で成功時はCommitする）
	defer tx.Rollback()

	// 2) DB登録
	params := repository.CreateEventParams{
		Title:            req.Title,
		Organizer:        req.Organizer,
		StartDate:        req.StartDate,
		StartTime:        req.StartTime,
		EndDate:          req.EndDate,
		EndTime:          req.EndTime,
		Email:            req.Email,
		Prefecture:       req.Prefecture,
		EventType:        req.EventType,
		IsOnline:         req.IsOnline,
		IsOffline:        req.IsOffline,
		OfficialURL:      req.OfficialURL,
		OnlineLectureURL: req.OnlineLectureURL,
		Venue:            req.Venue,
		Target:           req.Target,
		Capacity:         req.Capacity,
		Description:      req.Description,
		Tags:             req.Tags,
		Speakers:         convertSpeakers(req.Speakers),
		Schedule:         convertSchedules(req.Schedule),
	}

	eventID, authCode, err := h.repo.CreateEventTx(ctx, tx, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			"failed to create event").SetInternal(err)
	}

	// 3) Slack通知
	InitSlack()
	authLink := fmt.Sprintf(config.CORE_BACKEND_URL+"/api/v1/event/update/%d?auth_code=%s",
		eventID, authCode)
	slackMessage := fmt.Sprintf(
		"新しいイベントが作成されました。\nタイトル: %s\nオーガナイザー: %s\n認証リンク: %s",
		req.Title, req.Organizer, authLink,
	)
	if err := postToSlack(slackMessage); err != nil {
		// 失敗→ROLLBACKしてエラー応答
		return echo.NewHTTPError(http.StatusInternalServerError,
			"failed to send Slack notification").SetInternal(err)
	}

	// 4) COMMIT
	if err := tx.Commit(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			"failed to commit transaction").SetInternal(err)
	}

	// 正常レスポンス
	res := CreateEventResponse{
		ID: fmt.Sprintf("%d", eventID),
	}
	return c.JSON(http.StatusOK, res)
}

// GET /api/v1/event/all
func (h *Handler) GetEvents(c echo.Context) error {
	events, err := h.repo.GetEvents(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	eventsResponse := make([]GetEventResponse, len(events))
	for i, event := range events {
		eventsResponse[i] = GetEventResponse{
			ID:               event.ID,
			Title:            event.Title,
			Organizer:        event.Organizer,
			StartDate:        event.StartDate,
			StartTime:        event.StartTime,
			EndDate:          event.EndDate,
			EndTime:          event.EndTime,
			Email:            event.Email,
			Prefecture:       event.Prefecture,
			EventType:        event.EventType,
			IsOnline:         event.IsOnline,
			IsOffline:        event.IsOffline,
			OfficialURL:      event.OfficialURL,
			OnlineLectureURL: event.OnlineLectureURL,
			Venue:            event.Venue,
			Target:           event.Target,
			Capacity:         event.Capacity,
			Description:      event.Description,
			Tags:             event.Tags,
			Speakers:         convertSpeakersToResponse(event.Speakers),
			Schedule:         convertSchedulesToResponse(event.Schedule),
		}
	}
	return c.JSON(http.StatusOK, eventsResponse)
}

// GET /api/v1/event/update/:id
func (h *Handler) UpdateEvent(c echo.Context) error {
	idParam := c.Param("id")
	authCode := c.QueryParam("auth_code")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event ID").SetInternal(err)
	}

	err = h.repo.AuthenticateEvent(c.Request().Context(), id, authCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "authentication failed").SetInternal(err)
	}

	res := UpdateEventResponse{
		Message: "Event authenticated successfully",
	}
	return c.JSON(http.StatusOK, res)
}

// GET /api/v1/event/:id
func (h *Handler) GetEvent(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event ID").SetInternal(err)
	}

	event, err := h.repo.GetEvent(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	if event == nil {
		return echo.NewHTTPError(http.StatusNotFound, "event not found")
	}

	res := GetEventResponse{
		ID:               event.ID,
		Title:            event.Title,
		Organizer:        event.Organizer,
		StartDate:        event.StartDate,
		StartTime:        event.StartTime,
		EndDate:          event.EndDate,
		EndTime:          event.EndTime,
		Email:            event.Email,
		Prefecture:       event.Prefecture,
		EventType:        event.EventType,
		IsOnline:         event.IsOnline,
		IsOffline:        event.IsOffline,
		OfficialURL:      event.OfficialURL,
		OnlineLectureURL: event.OnlineLectureURL,
		Venue:            event.Venue,
		Target:           event.Target,
		Capacity:         event.Capacity,
		Description:      event.Description,
		Tags:             event.Tags,
		Speakers:         convertSpeakersToResponse(event.Speakers),
		Schedule:         convertSchedulesToResponse(event.Schedule),
	}
	return c.JSON(http.StatusOK, res)
}
