package util

import (
	"github.com/simelo/rextporter/src/core"
)

// MergeStoresInplace to update key / values in destination with those in source
func MergeStoresInplace(dst, src core.RextKeyValueStore) (err error) {
	var value interface{}
	err = nil
	for _, k := range src.GetKeys() {
		if value, err = src.GetObject(k); err == nil {
			if _, err = dst.SetObject(k, value); err != nil {
				return
			}
		} else {
			return
		}
	}
	return
}
