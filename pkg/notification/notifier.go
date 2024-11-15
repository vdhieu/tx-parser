package notification

// Notifier interface of our notification service
type Notifier interface {
	// Notify send notification to subscriber
	Notify(address string, message string, payload interface{}) error
}
