// internal/repository/event_test.go
package repository_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql" // MySQL/MariaDBドライバ
	"github.com/stretchr/testify/require"

	"github.com/ras0q/go-backend-template/internal/repository"
)

// テスト用のDBへの接続文字列（環境変数などから取得する想定）
var dsn = os.Getenv("TEST_DB_DSN")

func setupTestDB(t *testing.T) *sql.DB {
	if dsn == "" {
		// 必要に応じて DSN が無い場合はスキップ or fail
		t.Fatal("TEST_DB_DSN is not set")
	}

	db, err := sql.Open("mysql", dsn)
	require.NoError(t, err, "failed to open db")

	// テーブル作成など準備
	// すでに migration が走っているなら不要かもしれません
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS events (
        id INT PRIMARY KEY AUTO_INCREMENT,
        title VARCHAR(255) NOT NULL,
        organizer VARCHAR(255) NOT NULL,
        start_date DATE NOT NULL,
        start_time TIME NOT NULL,
        end_date DATE NOT NULL,
        end_time TIME NOT NULL,
        email VARCHAR(255) NOT NULL,
        prefecture VARCHAR(255),
        event_type VARCHAR(255),
        is_online BOOLEAN DEFAULT FALSE,
        is_offline BOOLEAN DEFAULT FALSE,
        official_url VARCHAR(255),
        online_lecture_url VARCHAR(255),
        venue VARCHAR(255),
        target VARCHAR(255),
        capacity VARCHAR(255),
        description TEXT,
        tags JSON,
        speakers JSON,
        schedule JSON,
        auth_code VARCHAR(36),
        is_authenticated BOOLEAN DEFAULT FALSE
    );
  `)
	require.NoError(t, err, "failed to create table")

	// テーブルの初期化（データ削除など）
	_, err = db.Exec(`DELETE FROM events`)
	require.NoError(t, err)

	return db
}

func TestRepository_CreateEventTx(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.New(db)
	ctx := context.Background()

	// トランザクション開始
	tx, err := repo.BeginTx(ctx)
	require.NoError(t, err)

	// パラメータ準備
	params := repository.CreateEventParams{
		Title:     "Test Event",
		Organizer: "Test Organizer",
		StartDate: "2025-01-01",
		StartTime: "09:00:00",
		EndDate:   "2025-01-01",
		EndTime:   "10:00:00",
		Email:     "test@example.com",
		// ほかのフィールドは省略
	}

	// CreateEventTx 実行
	eventID, authCode, err := repo.CreateEventTx(ctx, tx, params)
	require.NoError(t, err, "CreateEventTx should succeed")
	require.NotZero(t, eventID)
	require.NotEmpty(t, authCode)

	// ロールバックしてみる（本来はコミットする）
	err = tx.Rollback()
	require.NoError(t, err)

	// ロールバックしたので、eventsテーブルにレコードが無いか確認
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM events`).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 0, count, "rolled back transaction => no rows expected")
}

func TestRepository_GetEvents(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.New(db)
	ctx := context.Background()

	// 事前に認証済みレコードを挿入しておく
	_, err := db.Exec(`
    INSERT INTO events (title, organizer, start_date, start_time, end_date, end_time,
                        email, is_authenticated)
    VALUES ('Event1','Org1','2025-01-01','09:00:00','2025-01-01','10:00:00',
            'test1@example.com', TRUE),
           ('Event2','Org2','2025-01-02','10:00:00','2025-01-02','11:00:00',
            'test2@example.com', FALSE)
  `)
	require.NoError(t, err)

	events, err := repo.GetEvents(ctx)
	require.NoError(t, err)

	// is_authenticated=TRUE のレコードのみ返るはず => 1件だけ
	require.Len(t, events, 1)
	require.Equal(t, "Event1", events[0].Title)
}
