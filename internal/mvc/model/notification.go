package model

import "encoding/gob"

const (
	// NotificationError represents a notification for actions that failed.
	NotificationError NotificationType = iota

	// NotificationSuccess represents a notification for actions that succeeded.
	NotificationSuccess = iota

	// NotificationKey is the key used to store notifications in the session.
	NotificationKey = "view:notification"
)

// NotificationType determins the type of notification to show.
type NotificationType int

// Notification provides users with feedback on their actions.
type Notification struct {
	Type    NotificationType
	Title   string
	Content string
}

// NewErrorNotification creates a new error Notification.
func NewErrorNotification(title string, content string) *Notification {
	return &Notification{
		Title:   title,
		Content: content,
	}
}

// NewSuccessNotification creates a new error Notification.
func NewSuccessNotification(title string, content string) *Notification {
	return &Notification{
		Title:   title,
		Content: content,
	}
}

func (notificationType NotificationType) String() string {
	switch notificationType {
	case NotificationError:
		return "Error"
	case NotificationSuccess:
		return "Success"
	default:
		return "Success"
	}
}

func init() {
	gob.Register(Notification{})
}
