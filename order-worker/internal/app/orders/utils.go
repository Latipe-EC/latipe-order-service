package orders

import "time"

func IsAfterSevenDays(t time.Time) bool {
	now := time.Now()

	sevenDaysAgo := t.Add(7 * 24 * time.Hour)
	return now.After(sevenDaysAgo) || t.Equal(sevenDaysAgo)
}
