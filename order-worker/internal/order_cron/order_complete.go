package order_cron

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/robfig/cron/v3"
	"order-worker/config"
	"order-worker/internal/app/orders"
	"sync"
	"time"
)

type OrderCompleteCronjob struct {
	cron      *cron.Cron
	config    *config.Config
	orderServ orders.Usecase
}

func NewOrderCompleteCronjob(cron *cron.Cron, cfg *config.Config, orderServ orders.Usecase) *OrderCompleteCronjob {
	return &OrderCompleteCronjob{
		cron:      cron,
		config:    cfg,
		orderServ: orderServ,
	}
}

func (oc OrderCompleteCronjob) CheckOrderFinishShippingStatus(wg *sync.WaitGroup) {

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Infof("create new checking order status order_cron [schedule :%v]", oc.config.CronJob.OrderCompleteScheduled)
	runCron, err := oc.cron.AddFunc(oc.config.CronJob.OrderCompleteScheduled, func() {
		log.Info("Starting run function checking order status finish shipping")
		if err := oc.orderServ.UpdateCommissionOrderComplete(ctxTimeout); err != nil {
			log.Errorf("[CronJob] error cause:%v", error.Error)
		}
	})

	if err != nil {
		oc.cron.Remove(runCron)
	}
	oc.cron.Run()

	defer wg.Done()

}
