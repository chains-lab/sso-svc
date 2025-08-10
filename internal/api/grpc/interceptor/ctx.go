package interceptor

type ctxKey int

const (
	LogCtxKey ctxKey = iota
	RequestIDCtxKey
)
