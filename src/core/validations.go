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

func ValidateResource(r RextResourceDef) (hasError bool) {
	if len(r.GetType()) == 0 {
		hasError = true
		log.Errorln("type is required in metric config")
	}
	if len(r.GetResourcePATH("")) == 0 {
		hasError = true
		log.Errorln("resource path is required in metric config")
	}
	if r.GetDecoder() == nil {
		hasError = true
		log.Errorln("decoder is required in metric config")
	} else if r.GetDecoder().Validate() {
		hasError = true
	}
	if r.GetAuth(nil) != nil {
		if r.GetAuth(nil).Validate() {
			hasError = true
		}
	}
	for _, mtrDef := range r.GetMetricDefs() {
		mtrDef.Validate()
	}
	return hasError
}
