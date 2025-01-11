package handlers

type OperationType string

const (
	RESET_PASSWORD  OperationType = "reset_password"
	CHANGE_PASSWIRD OperationType = "change_password"
	CHANGE_EMAIL    OperationType = "change_email"
	REGISTRATION    OperationType = "registration"
	LOGIN           OperationType = "login"
)

func (op OperationType) IsValid() bool {
	switch op {
	case RESET_PASSWORD, REGISTRATION, LOGIN, CHANGE_PASSWIRD, CHANGE_EMAIL:
		return true
	default:
		return false
	}
}
