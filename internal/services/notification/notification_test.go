package notification_test

import (
	"testing"

	"github.com/edulustosa/go-pay/internal/database/models"
	"github.com/edulustosa/go-pay/internal/services/notification"
	"github.com/google/uuid"
)

func TestNotificationService(t *testing.T) {
	user := &models.User{
		ID:        uuid.New(),
		FirstName: "john",
		LastName:  "doe",
		Document:  "123456789",
		Email:     "johndoe@email.com",
		Balance:   0,
		Role:      models.RoleCommon,
	}

	if err := notification.Send(user, "test message"); err != nil {
		if err == notification.ErrNotificationUnavailable {
			t.Skip("notification service unavailable")
		}

		t.Fatalf("expected no error, got %v", err)
	}

	t.Log("notification sent successfully")
}
