package order

import (
	"fmt"
	"time"
)

func InitDateValue() string {
	now := time.Now()
	return fmt.Sprintf("%v-%v-%v", now.Year(), int(now.Month()), "01")
}
