package userlimit

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"golang.org/x/time/rate"
	"e2clicker.app/internal/publicerrors"
)

// LimitExceededError is an error that is returned when the user has reached the
// limit.
type LimitExceededError struct {
	reservation *rate.Reservation
}

func init() {
	publicerrors.MarkTypePublic[*LimitExceededError]()
}

// AsError converts a rate.Reservation to a LimitExceededError if the rate limit
// is exceeded. Otherwise, it returns nil.
func AsError(r *rate.Reservation) error {
	if r == nil {
		return nil
	}
	if r.OK() && r.Delay() < time.Second/2 {
		// sleeping for half a second is ok
		time.Sleep(r.Delay())
		return nil
	}
	// cancel the reservation
	r.Cancel()
	return &LimitExceededError{reservation: r}
}

func (e *LimitExceededError) Error() string {
	if e.reservation == nil {
		return "rate limit reached"
	}
	return fmt.Sprintf("rate limit reached, retry after %s", e.reservation.Delay())
}

func (e *LimitExceededError) MarshalJSON() ([]byte, error) {
	if e.reservation == nil {
		return json.Marshal(map[string]any{
			"error": "rate limit reached",
		})
	}
	return json.Marshal(map[string]any{
		"error": "rate limit reached",
		"delay": e.reservation.Delay(),
	})
}

// Delay returns the delay after which the request can be retried.
func (e *LimitExceededError) Delay() (time.Duration, bool) {
	if e.reservation == nil {
		return 0, false
	}
	return e.reservation.Delay(), true
}

// DelaySeconds returns the delay in seconds after which the request can be
// retried. This is meant for the Retry-After header in HTTP responses.
func (e *LimitExceededError) DelaySeconds() (int, bool) {
	d, ok := e.Delay()
	if !ok {
		return 0, false
	}
	return int(math.Ceil(d.Seconds())), true
}
