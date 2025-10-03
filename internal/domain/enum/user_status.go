package enum

import "fmt"

const (
	UserStatusActive  = "active"
	UserStatusBlocked = "blocked"
)

var userStatuses = []string{
	UserStatusActive,
	UserStatusBlocked,
}

var ErrorUserStatusIsNotSupported = fmt.Errorf("user status is not supported, must be one of: %v", GetAllUserStatuses())

func CheckUserStatus(status string) error {
	for _, userStatus := range userStatuses {
		if userStatus == status {
			return nil
		}
	}

	return fmt.Errorf("%s: %w", status, ErrorUserStatusIsNotSupported)
}

func GetAllUserStatuses() []string {
	return userStatuses
}
