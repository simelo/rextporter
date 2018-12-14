package memconfig

import (
	"sort"

	"github.com/simelo/rextporter/src/core"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func eqKvs(assert *assert.Assertions, kvs1 core.RextKeyValueStore, kvs2 core.RextKeyValueStore) bool {
	ks1 := kvs1.GetKeys()
	sort.Strings(ks1)
	ks2 := kvs2.GetKeys()
	sort.Strings(ks2)
	if assert != nil {
		assert.Equal(len(ks1), len(ks2))
	}
	if len(ks1) != len(ks2) {
		return false
	}
	for idxKs := range ks1 {
		if assert != nil {
			assert.Equal(ks1[idxKs], ks2[idxKs])
		}
		if ks1[idxKs] != ks2[idxKs] {
			return false
		}
		v1, err := kvs1.GetString(ks1[idxKs])
		if assert != nil {
			assert.Nil(err)
		}
		if err != nil {
			log.WithError(err).Errorln("Unexpected error")
			return false
		}
		v2, err := kvs2.GetString(ks2[idxKs])
		if assert != nil {
			assert.Nil(err)
		}
		if err != nil {
			log.WithError(err).Errorln("Unexpected error")
			return false
		}
		if assert != nil {
			assert.Equal(v1, v2)
		}
		if v1 != v2 {
			return false
		}
	}
	return true
}

func eqKvss(assert *assert.Assertions, kvss1 []core.RextKeyValueStore, kvss2 []core.RextKeyValueStore) bool {
	if assert != nil {
		assert.Equal(len(kvss1), len(kvss2))
	}
	if len(kvss1) != len(kvss2) {
		return false
	}
	for idxKvs := range kvss1 {
		if !eqKvs(assert, kvss1[idxKvs], kvss2[idxKvs]) {
			return false
		}
	}
	return true
}
