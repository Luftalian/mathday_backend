package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// Event はeventsテーブル1行分の構造体を表します。
type Event struct {
	ID               int        `db:"id"`
	Title            string     `db:"title"`
	Organizer        string     `db:"organizer"`
	StartDate        string     `db:"start_date"`
	StartTime        string     `db:"start_time"`
	EndDate          string     `db:"end_date"`
	EndTime          string     `db:"end_time"`
	Email            string     `db:"email"`
	Prefecture       *string    `db:"prefecture"`
	EventType        *string    `db:"event_type"`
	IsOnline         bool       `db:"is_online"`
	IsOffline        bool       `db:"is_offline"`
	OfficialURL      *string    `db:"official_url"`
	OnlineLectureURL *string    `db:"online_lecture_url"`
	Venue            *string    `db:"venue"`
	Target           *string    `db:"target"`
	Capacity         *string    `db:"capacity"`
	Description      *string    `db:"description"`
	Tags             []string   `db:"tags"`
	Speakers         []Speaker  `db:"speakers"`
	Schedule         []Schedule `db:"schedule"`
	AuthCode         string     `db:"auth_code"`
	IsAuthenticated  bool       `db:"is_authenticated"`
}

// Speaker はスピーカー情報を表します。
type Speaker struct {
	Name         string `json:"name"`
	Title        string `json:"title"`
	Organization string `json:"organization"`
}

// Schedule はスケジュール情報を表します。
type Schedule struct {
	Time    string `json:"time"`
	Title   string `json:"title"`
	Speaker string `json:"speaker"`
}

// CreateEventParams はイベント作成時に必要なパラメータです。
type CreateEventParams struct {
	Title            string
	Organizer        string
	StartDate        string
	StartTime        string
	EndDate          string
	EndTime          string
	Email            string
	Prefecture       *string
	EventType        *string
	IsOnline         bool
	IsOffline        bool
	OfficialURL      *string
	OnlineLectureURL *string
	Venue            *string
	Target           *string
	Capacity         *string
	Description      *string
	Tags             []string
	Speakers         []Speaker
	Schedule         []Schedule
}

// BeginTx は新たにトランザクションを開始します。
func (r *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// CreateEventTx はトランザクション内でイベントをINSERTし、生成されたIDと認証コードを返します。
func (r *Repository) CreateEventTx(ctx context.Context, tx *sql.Tx, params CreateEventParams) (int, string, error) {
	authCode := uuid.New().String()

	tagsJSON, err := json.Marshal(params.Tags)
	if err != nil {
		return 0, "", fmt.Errorf("タグのシリアライズに失敗: %w", err)
	}
	speakersJSON, err := json.Marshal(params.Speakers)
	if err != nil {
		return 0, "", fmt.Errorf("スピーカーのシリアライズに失敗: %w", err)
	}
	scheduleJSON, err := json.Marshal(params.Schedule)
	if err != nil {
		return 0, "", fmt.Errorf("スケジュールのシリアライズに失敗: %w", err)
	}

	query := `
		INSERT INTO events (
			title, organizer, start_date, start_time, end_date, end_time, email,
			prefecture, event_type, is_online, is_offline, official_url,
			online_lecture_url, venue, target, capacity, description, tags,
			speakers, schedule, auth_code, is_authenticated
		) VALUES (
			?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,FALSE
		)
	`
	result, err := tx.ExecContext(ctx, query,
		params.Title,
		params.Organizer,
		params.StartDate,
		params.StartTime,
		params.EndDate,
		params.EndTime,
		params.Email,
		params.Prefecture,
		params.EventType,
		params.IsOnline,
		params.IsOffline,
		params.OfficialURL,
		params.OnlineLectureURL,
		params.Venue,
		params.Target,
		params.Capacity,
		params.Description,
		tagsJSON,
		speakersJSON,
		scheduleJSON,
		authCode,
	)
	if err != nil {
		return 0, "", fmt.Errorf("イベントの挿入に失敗: %w", err)
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		return 0, "", fmt.Errorf("最後の挿入IDの取得に失敗: %w", err)
	}

	return int(eventID), authCode, nil
}

// GetEvents は認証済みのイベント一覧を取得します。
func (r *Repository) GetEvents(ctx context.Context) ([]*Event, error) {
	query := `
		SELECT id, title, organizer, start_date, start_time, end_date, end_time, email,
		       prefecture, event_type, is_online, is_offline, official_url,
		       online_lecture_url, venue, target, capacity, description, tags,
		       speakers, schedule, is_authenticated
		FROM events
		WHERE is_authenticated = TRUE
		ORDER BY start_date, start_time
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("イベントの取得に失敗: %w", err)
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var event Event
		var tagsJSON, speakersJSON, scheduleJSON []byte

		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Organizer,
			&event.StartDate,
			&event.StartTime,
			&event.EndDate,
			&event.EndTime,
			&event.Email,
			&event.Prefecture,
			&event.EventType,
			&event.IsOnline,
			&event.IsOffline,
			&event.OfficialURL,
			&event.OnlineLectureURL,
			&event.Venue,
			&event.Target,
			&event.Capacity,
			&event.Description,
			&tagsJSON,
			&speakersJSON,
			&scheduleJSON,
			&event.IsAuthenticated,
		); err != nil {
			return nil, fmt.Errorf("イベントのスキャンに失敗: %w", err)
		}

		// JSONフィールドのデシリアライズ
		if err := json.Unmarshal(tagsJSON, &event.Tags); err != nil {
			return nil, fmt.Errorf("タグのデシリアライズに失敗: %w", err)
		}
		if err := json.Unmarshal(speakersJSON, &event.Speakers); err != nil {
			return nil, fmt.Errorf("スピーカーのデシリアライズに失敗: %w", err)
		}
		if err := json.Unmarshal(scheduleJSON, &event.Schedule); err != nil {
			return nil, fmt.Errorf("スケジュールのデシリアライズに失敗: %w", err)
		}

		events = append(events, &event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("イベント行の処理中にエラー発生: %w", err)
	}
	return events, nil
}

// AuthenticateEvent は id と auth_code が一致するイベントを認証済みに更新します。
func (r *Repository) AuthenticateEvent(ctx context.Context, id int, authCode string) error {
	query := `
		UPDATE events
		SET is_authenticated = TRUE
		WHERE id = ? AND auth_code = ?
	`
	result, err := r.db.ExecContext(ctx, query, id, authCode)
	if err != nil {
		return fmt.Errorf("イベント認証の更新に失敗: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("更新された行数の取得に失敗: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("指定されたIDと認証コードに一致するイベントが見つかりません")
	}
	return nil
}

// GetEvent は指定されたIDの認証済みイベントを1件取得します。
func (r *Repository) GetEvent(ctx context.Context, id int) (*Event, error) {
	query := `
		SELECT id, title, organizer, start_date, start_time, end_date, end_time, email,
		       prefecture, event_type, is_online, is_offline, official_url,
		       online_lecture_url, venue, target, capacity, description, tags,
		       speakers, schedule, is_authenticated
		FROM events
		WHERE id = ? AND is_authenticated = TRUE
	`
	var event Event
	var tagsJSON, speakersJSON, scheduleJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Organizer,
		&event.StartDate,
		&event.StartTime,
		&event.EndDate,
		&event.EndTime,
		&event.Email,
		&event.Prefecture,
		&event.EventType,
		&event.IsOnline,
		&event.IsOffline,
		&event.OfficialURL,
		&event.OnlineLectureURL,
		&event.Venue,
		&event.Target,
		&event.Capacity,
		&event.Description,
		&tagsJSON,
		&speakersJSON,
		&scheduleJSON,
		&event.IsAuthenticated,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("イベントの取得に失敗: %w", err)
	}

	// JSONの復元
	if err := json.Unmarshal(tagsJSON, &event.Tags); err != nil {
		return nil, fmt.Errorf("タグのデシリアライズに失敗: %w", err)
	}
	if err := json.Unmarshal(speakersJSON, &event.Speakers); err != nil {
		return nil, fmt.Errorf("スピーカーのデシリアライズに失敗: %w", err)
	}
	if err := json.Unmarshal(scheduleJSON, &event.Schedule); err != nil {
		return nil, fmt.Errorf("スケジュールのデシリアライズに失敗: %w", err)
	}

	return &event, nil
}
