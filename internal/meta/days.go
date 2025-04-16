package meta

import (
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Days is a number of days. It acts as a duration of time, so 1.5 Days is
// 36 hours.
type Days float64

// DaysFromDuration converts time.Duration to Days.
func DaysFromDuration(d time.Duration) Days {
	return Days(d.Hours() / 24)
}

// DaysFromPostgreSQLInterval converts a PostgreSQL interval to Days.
func DaysFromPostgreSQLInterval(i pgtype.Interval) Days {
	days := float64(i.Days)
	days += float64(i.Microseconds) / 1e6 / (60 * 60 * 24)
	days += float64(i.Months) * 30
	return Days(days)
}

// ToDuration converts Days to time.Duration.
func (d Days) ToDuration() time.Duration {
	return time.Duration(float64(d) * float64(24*time.Hour))
}

// ToPostgreSQLInterval converts Days to a PostgreSQL interval.
func (d Days) ToPostgreSQLInterval() pgtype.Interval {
	int, frac := math.Modf(float64(d))
	return pgtype.Interval{
		Days:         int32(int),
		Microseconds: int64(frac * 24 * 60 * 60 * 1e6),
		Valid:        true,
	}
}
