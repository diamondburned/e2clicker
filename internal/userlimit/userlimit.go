package userlimit

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/puzpuzpuz/xsync/v3"
	"golang.org/x/time/rate"
)

// UserRateLimiter is a rate limiter for each user.
// It extends x/time/rate.Limiter's capabilities.
type UserRateLimiter[IDType comparable] struct {
	users *xsync.MapOf[IDType, *rateLimiter]
	used  chan struct{}
	rate  rate.Limit
	burst int
}

// NewUserRateLimiter creates a new UserRateLimiter.
func NewUserRateLimiter[IDType comparable](r rate.Limit, burst int) *UserRateLimiter[IDType] {
	return &UserRateLimiter[IDType]{
		users: xsync.NewMapOf[IDType, *rateLimiter](),
		used:  make(chan struct{}, 1),
		rate:  r,
		burst: burst,
	}
}

// BeginCleanup starts a goroutine that cleans up unused rate limiters.
func (l *UserRateLimiter[IDType]) BeginCleanup() (stop func()) {
	stopCh := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var timer <-chan time.Time

		for {
			select {
			case <-stopCh:
				return
			case <-l.used:
				if timer == nil {
					timer = time.After(30 * time.Minute)
				}
			case <-timer:
				l.cleanup()
			}
		}
	}()

	return sync.OnceFunc(func() {
		close(stopCh)
		wg.Wait()
	})
}

func (l *UserRateLimiter[IDType]) markUsed() {
	select {
	case l.used <- struct{}{}:
	default:
	}
}

// Reserve is shorthand for ReserveN(time.Now(), 1).
func (l *UserRateLimiter[IDType]) Reserve(id IDType) *rate.Reservation {
	return l.ReserveN(id, time.Now(), 1)
}

// ReserveN returns a Reservation that indicates how long the caller must wait before n events happen.
// The Limiter takes this Reservation into account when allowing future events.
// The returned Reservation’s OK() method returns false if n exceeds the Limiter's burst size.
// Usage example:
//
//	r := lim.ReserveN(time.Now(), 1)
//	if !r.OK() {
//	  // Not allowed to act! Did you remember to set lim.burst to be > 0 ?
//	  return
//	}
//	time.Sleep(r.Delay())
//	Act()
//
// Use this method if you wish to wait and slow down in accordance with the rate limit without dropping events.
// If you need to respect a deadline or cancel the delay, use Wait instead.
// To drop or skip events exceeding rate limit, use Allow instead.
func (l *UserRateLimiter[IDType]) ReserveN(id IDType, t time.Time, n int) *rate.Reservation {
	l.markUsed()

	limiter, _ := l.users.LoadOrCompute(
		id,
		func() *rateLimiter { return newRateLimiter(l.rate, l.burst) },
	)

	r := limiter.ReserveN(t, n)
	if !r.OK() {
		return r
	}

	limiter.markUsed()
	return r
}

func (l *UserRateLimiter[IDType]) cleanup() {
	l.users.Range(func(key IDType, value *rateLimiter) bool {
		if time.Since(value.lastUsedTime()) > 5*time.Minute {
			l.users.Delete(key)
		}
		return true
	})
}

type rateLimiter struct {
	rate.Limiter
	lastUsed atomic.Int64
}

func newRateLimiter(limit rate.Limit, burst int) *rateLimiter {
	r := &rateLimiter{Limiter: *rate.NewLimiter(limit, burst)}
	r.markUsed()
	return r
}

func (r *rateLimiter) markUsed() {
	r.lastUsed.Store(time.Now().Unix())
}

func (r *rateLimiter) lastUsedTime() time.Time {
	return time.Unix(r.lastUsed.Load(), 0)
}
