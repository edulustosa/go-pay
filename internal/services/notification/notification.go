package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/edulustosa/go-pay/internal/database/models"
	"github.com/edulustosa/go-pay/internal/dtos"
)

var ErrNotificationUnavailable = errors.New("notification service unavailable")

func Send(user *models.User, message string) error {
	notification := dtos.NotificationDTO{
		Email:   user.Email,
		Message: message,
	}

	reqBytes, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://util.devi.tools/api/v1/notify",
		bytes.NewReader(reqBytes),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return ErrNotificationUnavailable
	}

	return nil
}
