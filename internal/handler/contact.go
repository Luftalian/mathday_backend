package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CreateContact
// name, email, messageを受け取り、slackに通知を送信する

// POST /api/v1/contact
// Request Body
type CreateContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// CreateContact
// slackに通知を送信する
func (h *Handler) CreateContact(c echo.Context) error {
	req := new(CreateContactRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").SetInternal(err)
	}

	// 3) Slack通知
	InitSlack()
	slackMessage := fmt.Sprintf(
		"お問い合わせがありました\n"+
			"名前: %s\n"+
			"メールアドレス: %s\n"+
			"メッセージ: %s",
		req.Name, req.Email, req.Message)

	if err := postToSlack(slackMessage); err != nil {
		// 失敗→ROLLBACKしてエラー応答
		return echo.NewHTTPError(http.StatusInternalServerError,
			"failed to send Slack notification").SetInternal(err)
	}

	return c.JSON(http.StatusOK, "ok")
}
