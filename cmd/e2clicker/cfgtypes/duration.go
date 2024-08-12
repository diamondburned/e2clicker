package cfgtypes

import (
	"encoding/json"
	"fmt"
	"time"
)

// Duration is a type that can be used to unmarshal a time.Duration from JSON.
type Duration time.Duration

func (d Duration) AsDuration() time.Duration {
	return time.Duration(d)
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("failed to parse duration: %w", err)
	}
	*d = Duration(dur)
	return nil
}
