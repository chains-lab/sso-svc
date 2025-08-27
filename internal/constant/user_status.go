package constant

import "fmt"

const (
	UserStatusActive  = "active"
	UserStatusBlocked = "blocked"
)

var userStatuses = []string{
	UserStatusActive,
	UserStatusBlocked,
}

var ErrorUserStatusIsNotSupported = fmt.Errorf("user status is not supported")

func ParseUserStatus(status string) error {
	for _, userStatus := range userStatuses {
		if userStatus == status {
			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrorUserStatusIsNotSupported, status)
}

func GetAllUserStatuses() []string {
	return userStatuses
}
