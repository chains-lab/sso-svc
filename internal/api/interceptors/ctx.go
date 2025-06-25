package interceptors

type ctxKey int

const (
	LogCtxKey ctxKey = iota
	MetaCtxKey
)
