package reporter

import (
	"log"
)

// ReportError acts as a single centralized error-reporting function.
// It logs unexpected errors to the console. If a service like Sentry
// were added later, this is where it would be configured.
func ReportError(err error) {
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
	}
}
