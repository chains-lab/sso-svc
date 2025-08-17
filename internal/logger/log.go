package logger

import (
	"context"
	"errors"
	"strings"

	"github.com/chains-lab/sso-svc/internal/api/grpc/meta"
	"github.com/chains-lab/sso-svc/internal/config"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func NewLogger(cfg config.Config) *logrus.Logger {
	log := logrus.New()

	lvl, err := logrus.ParseLevel(strings.ToLower(cfg.Server.Log.Level))
	if err != nil {
		log.Warnf("invalid log level '%s', defaulting to 'info'", cfg.Server.Log.Level)
		lvl = logrus.InfoLevel
	}
	log.SetLevel(lvl)

	switch strings.ToLower(cfg.Server.Log.Format) {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		fallthrough
	default:
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return log
}

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
			meta.LogCtxKey,
			log, // ваш интерфейс Logger
		)

		// Далее передаём новый контекст в реальный хэндлер
		return handler(ctxWithLog, req)
	}
}

func Log(ctx context.Context) Logger {
	entry, ok := ctx.Value(meta.LogCtxKey).(Logger)
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
	var ae *ape.Error
	if errors.As(err, &ae) {
		return l.Entry.WithError(ae.Unwrap())
	}
	// для “обычных” ошибок просто стандартный путь
	return l.Entry.WithError(err)
}

func NewWithBase(base *logrus.Logger) Logger {
	log := logger{
		Entry: logrus.NewEntry(base),
	}

	return &log
}
