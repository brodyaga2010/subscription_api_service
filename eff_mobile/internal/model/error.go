package model

import "github.com/pkg/errors"

var (
	ErrSubscriptionNotFound = errors.New("Subscription not found")
	ErrCreateSubscription   = errors.New("Subscription has been no created")
	ErrDateIsNull           = errors.New("From and To dates must be provided")
	ErrDateBefore           = errors.New("'To' date cannot be before 'From' date")
	ErrDateFormat           = errors.New("Date format is incorrect, expected 'MM-YYYY'")
)
