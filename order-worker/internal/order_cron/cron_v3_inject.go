package order_cron

import (
	"github.com/google/wire"
	"github.com/robfig/cron/v3"
)

var Set = wire.NewSet(
	NewCronInstance,
	NewOrderCompleteCronjob,
)

func NewCronInstance() *cron.Cron {
	return cron.New()
}
