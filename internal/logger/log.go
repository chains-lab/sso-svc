package logger

import (
	"context"

	"github.com/chains-lab/sso-svc/internal/api/grpc/interceptor"
	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func UnaryLogInterceptor(log Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Вместо context.Background() используем входящий ctx,
		// чтобы не потерять таймауты и другую информацию.
		ctxWithLog := context.WithValue(
			ctx,
			interceptor.LogCtxKey,
			log, // ваш интерфейс Logger
		)

		// Далее передаём новый контекст в реальный хэндлер
		return handler(ctxWithLog, req)
	}
}

func Log(ctx context.Context) Logger {
	entry, ok := ctx.Value(interceptor.LogCtxKey).(Logger)
	if !ok {
		logrus.Info("no logger in context")

		entry = NewWithBase(logrus.New())
	}

	requestID := meta.RequestID(ctx)

	return &logger{Entry: entry.WithField("request_id", requestID)}
}

// Logger — это ваш интерфейс: все методы FieldLogger + специальный WithError.
type Logger interface {
	WithError(err error) *logrus.Entry

	logrus.FieldLogger // сюда входят Debug, Info, WithField, WithError и т.д.
}

// logger — реальный тип, который реализует Logger.
type logger struct {
	*logrus.Entry // за счёт встраивания мы уже наследуем все методы FieldLogger
}

// WithError — ваш особый метод.
func (l *logger) WithError(err error) *logrus.Entry {
	ae := ape.Unwrap(err)
	if ae != nil {
		return l.Entry.WithError(ae.Unwrap())
	}

	return l.Entry.WithError(err)
}

func NewWithBase(base *logrus.Logger) Logger {
	log := logger{
		Entry: logrus.NewEntry(base),
	}

	return &log
}
