package bodies

import (
	"time"

	"github.com/chains-lab/gatekit/roles"
)

const (
	UsersTopic      = "users"
	UserCreateTopic = "user_create"

	UserCreateType = "USER_CREATE"
)

type UserCreated struct {
	UserID    string     `json:"user_id"`
	Role      roles.Role `json:"role"`
	Timestamp time.Time  `json:"timestamp"`
}
