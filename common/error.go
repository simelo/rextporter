package common

import (
	"errors"
	"log"
)

func ErrorFromThisScope(rootCause, generalScopeErr string) error {
	log.Println("error root cause:", rootCause)
	return errors.New(generalScopeErr)
}
