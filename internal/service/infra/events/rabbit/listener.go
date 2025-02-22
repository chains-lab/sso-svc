package rabbit

import (
	"context"

	"github.com/recovery-flow/rerabbit"
	"github.com/recovery-flow/sso-oauth/internal/service"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/callbacks"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/events/rabbit/amqpconfig"
	"github.com/streadway/amqp"
)

func Listener(ctx context.Context, svc *service.Service) {
	// Создаем подключение к RabbitMQ.
	rabbitWorker, err := rerabbit.NewBroker(svc.Config.Rabbit.URL)
	if err != nil {
		svc.Log.Errorf("Failed to connect to RabbitMQ: %v", err)
		<-ctx.Done()
		return
	}

	go func() {
		<-ctx.Done()
		svc.Log.Info("Shutting down RabbitMQ connection...")
		rabbitWorker.GracefulShutdown(svc.Log)
	}()

	type QueueConfig struct {
		QueueName  string
		RoutingKey string
		Callback   func(context.Context, *service.Service, []byte) error
	}

	queues := []QueueConfig{
		{
			QueueName:  amqpconfig.AccountQueue,
			RoutingKey: amqpconfig.AccountUpdateRoleKey,
			Callback:   callbacks.AccountUpdateRole, // Это функция из domain/callbacks
		},
	}

	// Для каждой очереди запускаем отдельного consumer-а.
	for _, qc := range queues {
		qc := qc // захватываем локальную копию, чтобы избежать гонок
		go func(qc QueueConfig) {
			opts := rerabbit.ConsumeOptions{
				QueueName:   qc.QueueName,
				ConsumerTag: "", // можно задать уникальный тег, если нужно
				AutoAck:     false,
				Exclusive:   false,
				NoLocal:     false,
				NoWait:      false,
				Args:        nil,
			}

			err := rabbitWorker.Consume(ctx, opts, func(ctx context.Context, d amqp.Delivery) {
				// Вызываем callback. Если произошла ошибка, отправляем отрицательное подтверждение (Nack).
				if err := qc.Callback(ctx, svc, d.Body); err != nil {
					svc.Log.Errorf("Error processing message from queue %s: %v", qc.QueueName, err)
					if nackErr := d.Nack(false, true); nackErr != nil {
						svc.Log.Errorf("Failed to Nack message: %v", nackErr)
					}
				} else {
					if ackErr := d.Ack(false); ackErr != nil {
						svc.Log.Errorf("Failed to Ack message: %v", ackErr)
					}
				}
			})
			if err != nil {
				svc.Log.Errorf("Error consuming from queue %s: %v", qc.QueueName, err)
			}
		}(qc)
	}

	<-ctx.Done()
	svc.Log.Info("Listener events [RabbitMQ] stopped")
}
