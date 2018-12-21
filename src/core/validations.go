package core

import (
	log "github.com/sirupsen/logrus"
)

func ValidateAuth(auth RextAuthDef) (hasError bool) {
	if len(auth.GetAuthType()) == 0 {
		hasError = true
		log.Errorln("type is required in auth config")
	}
	if auth.GetAuthType() == AuthTypeCSRF {
		opts := auth.GetOptions()
		if tkhk, err := opts.GetString(OptKeyRextAuthDefTokenHeaderKey); err != nil || len(tkhk) == 0 {
			hasError = true
			log.Errorln("token header key is required for CSRF auth type")
		}
		if tkge, err := opts.GetString(OptKeyRextAuthDefTokenGenEndpoint); err != nil || len(tkge) == 0 {
			hasError = true
			log.Errorln("token gen endpoint is required for CSRF auth type")
		}
		if tkfe, err := opts.GetString(OptKeyRextAuthDefTokenKeyFromEndpoint); err != nil || len(tkfe) == 0 {
			hasError = true
			log.Errorln("token from endpoint is required for CSRF auth type")
		}
	}
	return hasError
}
