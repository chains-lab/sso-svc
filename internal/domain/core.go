package domain

type Core struct {
	user    UserSvc
	session SessionSvc
}

func NewCore(user UserSvc, session SessionSvc) Core {
	return Core{
		user:    user,
		session: session,
	}
}

func (a *Core) User() UserSvc {
	return a.user
}

func (a *Core) Session() SessionSvc {
	return a.session
}
