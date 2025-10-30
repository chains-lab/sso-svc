package listener

import (
	"context"
	"sync"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/sso-svc/internal/events"
	"github.com/chains-lab/sso-svc/internal/events/listener/subscriber"
	"github.com/segmentio/kafka-go"
)

type Callbacks interface {
	UpdateEmployee(ctx context.Context, event kafka.Message) error
	UpdateCityAdmin(ctx context.Context, event kafka.Message) error
}

func Run(ctx context.Context, log logium.Logger, addr string, cb Callbacks) {
	var wg sync.WaitGroup

	employeeSub := subscriber.New(addr, events.TopicCompaniesEmployeeV1, events.GroupID)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := employeeSub.Subscribe(ctx, "employee.update", cb.UpdateEmployee); err != nil {
			log.Printf("employee listener stopped: %v", err)
		}
	}()

	cityAdminSub := subscriber.New(addr, events.TopicCitiesAdminV1, events.GroupID)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := cityAdminSub.Subscribe(ctx, "city_admin.update", cb.UpdateCityAdmin); err != nil {
			log.Printf("city_admin listener stopped: %v", err)
		}
	}()

	<-ctx.Done()
	wg.Wait()
}
