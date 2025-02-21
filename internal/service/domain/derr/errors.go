package derr

import "fmt"

var ErrTokenInvalid = fmt.Errorf("invalid token")

var ErrSessionNotBelongToUser = fmt.Errorf("sessions doesn't belong to account")

var SessionNotFound = fmt.Errorf("session not found")

var ErrAccountNotFound = fmt.Errorf("account not found")
