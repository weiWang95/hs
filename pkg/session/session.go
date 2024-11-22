package session

import "context"

type Session struct {
	UserId uint64
}

func Set(ctx context.Context, userId uint64) context.Context {
	return context.WithValue(ctx, "session", &Session{UserId: userId})
}

func Get(ctx context.Context) *Session {
	return ctx.Value("session").(*Session)
}
