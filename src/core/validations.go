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

func ValidateService(srv RextServiceDef) (hasError bool) {
	srvOpts := srv.GetOptions()
	jobName, err := srvOpts.GetString(OptKeyRextServiceDefJobName)
	if err != nil {
		hasError = true
		log.WithError(err).Errorln("key for job name not present service config")
	}
	if len(jobName) == 0 {
		hasError = true
		log.Errorln("job name is required in service config")
	}
	var instanceName string
	instanceName, err = srvOpts.GetString(OptKeyRextServiceDefInstanceName)
	if err != nil {
		hasError = true
		log.WithError(err).Errorln("key for job name not present service config")
	}
	if len(instanceName) == 0 {
		hasError = true
		log.Errorln("instance name is required in service config")
	}
	if len(srv.GetProtocol()) == 0 {
		hasError = true
		log.Errorln("protocol should not be null in service config")
	}
	if srv.GetAuthForBaseURL() != nil {
		if srv.GetAuthForBaseURL().Validate() {
			hasError = true
		}
	}
	for _, source := range srv.GetResources() {
		source.Validate()
	}
	return hasError
}
