package notification

import (
	"log"
)

type consoleNotifier struct{}

func NewConsoleNotifier() Notifier {
	return &consoleNotifier{}
}

func (n *consoleNotifier) Notify(address string, message string, payload interface{}) error {
	log.Printf("\n[ConsoleNotifier] Notification for %s: %s\nTransaction: %+v\n", address, message, payload)
	return nil
}
