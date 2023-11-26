package orders

import "time"

func IsAfterSevenDays(t time.Time) bool {
	now := time.Now()

	sevenDaysAgo := now.Add(-7 * 24 * time.Hour)
	return t.After(sevenDaysAgo) || t.Equal(sevenDaysAgo)
}
