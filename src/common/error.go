package common

import (
	"errors"
	"log"
)

// ErrorFromThisScope can be used to get notified about a general error(what do you want from this function)
// and a specific error(`rootCause`), the statamet who causes the error
// this method will log the rootCause and return the general error(`generalScopeErr`)
func ErrorFromThisScope(rootCause, generalScopeErr string) error {
	log.Println("error root cause:", rootCause)
	return errors.New(generalScopeErr)
}
