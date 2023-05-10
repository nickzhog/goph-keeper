package account

import "errors"

var ErrNotFound = errors.New("user not found")
var ErrWrongPassword = errors.New("wrong password")
var ErrInvalidAccessToken = errors.New("invalid auth token")
