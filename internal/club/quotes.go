package club

import (
	"fmt"
	"math"
	"time"
)

// DefaultQuote is shown when the corpus is somehow empty; mirrors
// DEFAULT_QUOTE in data/quotes/index.ts on the site.
var DefaultQuote = Quote{Text: "Talk is cheap. Show me the code.", Author: "Linus Torvalds"}

// quoteStride walks the corpus at the golden ratio of its length, nudged down
// until coprime with the length. Because the stride is coprime, day*STRIDE mod
// N is a permutation: every quote appears exactly once per N days (no repeats,
// no gaps) and consecutive days land ~62% of the list apart instead of walking
// it in file order. Ported verbatim from lib/quotes.ts on the site.
var quoteStride = func() int {
	n := len(Quotes)
	if n < 2 {
		return 1
	}
	s := int(math.Round(float64(n) * 0.618_033_988_7))
	for gcd(s, n) != 1 {
		s--
	}
	return s
}()

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// QuoteForDate picks the quote for a given local calendar day. Pure and
// deterministic: every call on the same day returns the same quote, and the
// pick matches what fossclubkiet.org renders for that day.
func QuoteForDate(t time.Time) Quote {
	n := len(Quotes)
	if n == 0 {
		return DefaultQuote
	}
	// Days since the epoch for the local Y/M/D (UTC-normalized to avoid DST drift).
	y, m, d := t.Date()
	day := time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix() / 86_400
	// ((x % n) + n) % n keeps the index non-negative for pre-1970 dates.
	idx := int(((day*int64(quoteStride))%int64(n) + int64(n)) % int64(n))
	return Quotes[idx]
}

// QuoteOfTheDay is QuoteForDate(time.Now()).
func QuoteOfTheDay() Quote {
	return QuoteForDate(time.Now())
}

// NextWeeklySync returns the next Tuesday 17:00 (5:00 PM) at or after t.
func NextWeeklySync(t time.Time) time.Time {
	next := time.Date(t.Year(), t.Month(), t.Day(), 17, 0, 0, 0, t.Location())
	if t.After(next) || t.Equal(next) {
		// Today's 5 PM has passed — roll to next Tuesday.
		next = next.AddDate(0, 0, 1)
	}
	for next.Weekday() != time.Tuesday {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

// UntilString renders a duration as compact "4d 18h 3m".
func UntilString(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh", days, hours)
	case hours > 0:
		return fmt.Sprintf("%dh %dm", hours, mins)
	default:
		return fmt.Sprintf("%dm", mins)
	}
}
