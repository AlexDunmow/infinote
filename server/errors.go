package boilerplate

import (
	"boilerplate/db"
	"context"
	"errors"
	"fmt"
	"log"
)

// ErrDataloader when could not retrieve data
var ErrDataloader = errors.New("could not fetch data")

// ErrBadClaims when JWT could not be read
var ErrBadClaims = errors.New("could not read credentials from JWT")

// ErrBlacklisted when a typecast fails
var ErrBlacklisted = errors.New("token has been blacklisted")

// ErrTypeCast when a typecast fails
var ErrTypeCast = errors.New("could not cast interface to type")

// ErrParse when a parse fails
var ErrParse = errors.New("could not parse input")

// ErrBadCredentials when a bad username or password is passed in
var ErrBadCredentials = errors.New("bad credentials")

// ErrNotImplemented for non-implemented funcs
var ErrNotImplemented = errors.New("not implemented")

// ErrUnauthorized for bad permissions
var ErrUnauthorized = errors.New("unauthorized")

// ErrBadContext for missing context values
var ErrBadContext = errors.New("bad context")

// ErrAuthNoEmail during authencation when login failed due to non-existant user email address
var ErrAuthNoEmail = errors.New("user not found")

// ErrAuthWrongPassword during authencation when login failed due to incorrect incorrect password
var ErrAuthWrongPassword = errors.New("wrong password")

// KindInput of error
const KindInput Kind = "input"

// KindSystem of error
const KindSystem Kind = "system"

// Error is the custom error type
type Error struct {
	ID      string            // unique uuid of error so can find and trace it
	Message string            // friendly message to the user
	Err     error             // actual that is refered to
	Kind    Kind              // kind of error
	Meta    map[string]string // any additional information that is useful in debugging error (backend only, do not expose this to user)
}

// Kind of error
type Kind string

// Err returns a new Error
func Err(ID string, Message string, Err error, kind Kind, kvs ...string) *Error {
	meta := map[string]string{}
	if len(kvs)%2 == 0 {
		prev := ""
		for i, val := range kvs {
			if i%2 == 0 {
				meta[val] = ""
			} else {
				meta[prev] = val
			}
			prev = val
		}
	} else {
		fmt.Println("ERROR: Number of KVs not even")
	}

	return &Error{
		ID:      ID,
		Message: Message,
		Err:     Err,
		Kind:    kind,
	}
}

// Unwrap the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}
func (e *Error) Error() string {
	return e.Message
}

// LogExternal will create external log and return GQError itself
func (e *Error) LogExternal() *Error {
	// TODO use sentry instead of log.Printf()
	log.Printf("\033[1;31mERROR\033[0m %+v", e)
	return e
}

// LogExternal will unwrap errors recursively and log to sentry and return an error that can be displayed in the frontend.
// If an error message intended for the user is found in the stack, this message will be returned. Otherwise it will return
// the generic error message.
func LogExternal(ctx context.Context, u *db.User, err error) error {
	// TODO use sentry instead of log.Printf()
	log.Printf("\033[1;31mERROR\033[0m %+v", err)
	return err
}
