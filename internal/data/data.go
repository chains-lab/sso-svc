package data

type Database struct {
	users    Users
	sessions Sessions
}

func NewDatabase(user Users, session Sessions) Database {
	return Database{
		users:    user,
		sessions: session,
	}
}

func (d *Database) Users() Users {
	return d.users
}

func (d *Database) Sessions() Sessions {
	return d.sessions
}
