package meta

type ctxKey int

const (
	LogCtxKey ctxKey = iota
	RequestIDCtxKey
	UserCtxKey
)
